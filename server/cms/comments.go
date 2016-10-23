package cms

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"herefriend/lib"
	"herefriend/lib/push"
	"herefriend/server/handlers"
)

const (
	// CmsPushMsgTypeNormalMsg .
	CmsPushMsgTypeNormalMsg = 1
	// CmsPushMsgTypeEvaluation .
	CmsPushMsgTypeEvaluation = 2
	// CmsLittleImageView .
	CmsLittleImageView = "?imageView2/5/w/50/h/50"
)

// CommentInfo .
func CommentInfo(c *gin.Context) {
	info := cmsCommentSummary{
		TalkNum:   handlers.GetAPIRecommendCount(),
		PushNum:   push.GetPushNum(),
		BuyVIPNum: handlers.GetBuyVIPCount(),
	}

	c.JSON(http.StatusOK, info)
}

// RecentComments 获取最新的消息
func RecentComments(c *gin.Context) {
	var lastmsgid int
	var fromid int
	var toid int

	lastmsgidstr := c.Query("lastmsgid")
	if "" == lastmsgidstr {
		lastmsgid = 0
	} else {
		lastmsgid, _ = strconv.Atoi(lastmsgidstr)
	}

	sentence := "select id,fromid,toid,time,type,msg from recommend where id>? order by id desc limit 20"
	rows, err := lib.SQLQuery(sentence, lastmsgid)
	if nil != err {
		c.String(http.StatusOK, "[]")
		return
	}
	defer rows.Close()

	var info cmsCommentInfo
	var timevalue int64
	var infos []cmsCommentInfo

	for rows.Next() {
		err = rows.Scan(&info.MsgID, &fromid, &toid, &timevalue, &info.MsgType, &info.MsgText)
		if nil != err {
			continue
		}

		code, userinfo := handlers.GetUserInfoByID(fromid)
		if 200 == code && "" != userinfo.Name {
			info.From = "[" + userinfo.Province + "]" + userinfo.Name
		} else {
			info.From = "[" + userinfo.Province + "]" + strconv.Itoa(fromid)
		}

		info.FromPic = userinfo.IconURL + CmsLittleImageView

		code, userinfo = handlers.GetUserInfoByID(toid)
		if 200 == code && "" != userinfo.Name {
			info.To = "[" + userinfo.Province + "]" + userinfo.Name
		} else {
			info.To = "[" + userinfo.Province + "]" + strconv.Itoa(toid)
		}
		info.ToPic = userinfo.IconURL + CmsLittleImageView

		info.TimeUTC = lib.Int64ToUTCTime(timevalue)
		infos = append(infos, info)
	}

	if len(infos) > 0 {
		len := len(infos)
		for i := 0; i < len/2; i++ {
			infos[i], infos[len-i-1] = infos[len-i-1], infos[i]
		}
	}

	c.JSON(http.StatusOK, infos)
}

// MsgTemplate .
func MsgTemplate(c *gin.Context) {
	var msgtype int
	var gender int

	msgtypestr := c.Query("type")
	genderstr := c.Query("gender")

	if "" == msgtypestr {
		msgtype = 0
	} else {
		msgtype, _ = strconv.Atoi(msgtypestr)
	}

	if "" == genderstr {
		gender = -1
	} else {
		gender, _ = strconv.Atoi(genderstr)
	}

	sentence := "select id,msg from msgtemplate where type=? and gender=?"

	rows, err := lib.SQLQuery(sentence, msgtype, gender)
	if nil != err {
		c.String(http.StatusOK, "[]")
	}
	defer rows.Close()

	var info cmsMessageTempalte
	var infos []cmsMessageTempalte

	for rows.Next() {
		err = rows.Scan(&info.ID, &info.Template)
		if nil != err {
			continue
		}

		infos = append(infos, info)
	}

	c.JSON(http.StatusOK, infos)
}

// MsgTemplateAdd .
func MsgTemplateAdd(c *gin.Context) {
	var msgtype int
	var gender int

	msgtypestr := c.Query("type")
	genderstr := c.Query("gender")
	templatestr := c.Query("template")

	if "" == templatestr {
		c.Status(http.StatusNotFound)
		return
	}

	if "" == msgtypestr {
		msgtype = 0
	} else {
		msgtype, _ = strconv.Atoi(msgtypestr)
	}

	if "" == genderstr {
		gender = -1
	} else {
		gender, _ = strconv.Atoi(genderstr)
	}

	sentence := "insert into msgtemplate (msg,type,gender) values (?,?,?)"
	result, err := lib.SQLExec(sentence, templatestr, msgtype, gender)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	lastid, err := result.LastInsertId()
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	handlers.ReloadRecommendTemplates()

	var info cmsMessageTempalte
	info.ID = int(lastid)
	info.Template = templatestr

	c.JSON(http.StatusOK, info)
}

// MsgTemplateDel .
func MsgTemplateDel(c *gin.Context) {
	idstr := c.Query("id")

	if "" == idstr {
		c.Status(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idstr)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	sentence := "delete from msgtemplate where id=?"
	_, err = lib.SQLExec(sentence, id)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	handlers.ReloadRecommendTemplates()
	return
}

// MsgTemplateModify .
func MsgTemplateModify(c *gin.Context) {
	idstr := c.Query("id")
	templatestr := c.Query("template")

	if "" == idstr || "" == templatestr {
		c.Status(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idstr)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	sentence := "update msgtemplate set msg=? where id=?"
	_, err = lib.SQLExec(sentence, templatestr, id)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	handlers.ReloadRecommendTemplates()
	return
}

// GetChartsList get the charts list
func GetChartsList(c *gin.Context) {
	idstr := c.Query("id")
	if idstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	id, _ := strconv.Atoi(idstr)

	var commentsInfo []cmsCommentInfo
	pageid, count := lib.GetPageidCount(c)
	recommendAlls, err := handlers.GetRecommendAll(0, id, pageid, count)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	_, userinfo := handlers.GetUserInfoByID(id)
	for _, r := range recommendAlls {
		c := cmsCommentInfo{
			MsgID:     r.MsgID,
			MsgText:   r.MsgText,
			TimeUTC:   r.TimeUTC,
			Direction: r.Direction,
		}

		_, _, usertype := handlers.GetGenderUsertypeByID(r.UserID)
		if 1 == usertype {
			continue
		}

		if 1 == c.Direction {
			c.FromID = id
			c.ToID = r.UserID
			c.From = "[" + userinfo.Province + "]" + userinfo.Name
			if "" != userinfo.IconURL {
				c.FromPic = userinfo.IconURL + CmsLittleImageView
			}
			c.To = "[" + r.UserInfo.Province + "]" + r.UserInfo.Name
			if "" != r.UserInfo.IconURL {
				c.ToPic = r.UserInfo.IconURL + CmsLittleImageView
			}
		} else {
			c.FromID = r.UserID
			c.ToID = id
			c.To = "[" + userinfo.Province + "]" + userinfo.Name
			if "" != userinfo.IconURL {
				c.ToPic = userinfo.IconURL + CmsLittleImageView
			}
			c.From = "[" + r.UserInfo.Province + "]" + r.UserInfo.Name
			if "" != r.UserInfo.IconURL {
				c.FromPic = r.UserInfo.IconURL + CmsLittleImageView
			}
		}

		commentsInfo = append(commentsInfo, c)
	}

	c.JSON(http.StatusOK, commentsInfo)
}

// GetTalkHistory .
func GetTalkHistory(c *gin.Context) {
	idStr := c.Query("id")
	talkidStr := c.Query("talkid")
	if talkidStr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	id, _ := strconv.Atoi(idStr)
	talkid, _ := strconv.Atoi(talkidStr)

	var lastMsgID int
	lastMsgIDStr := c.Query("lastmsgid")
	if "" != lastMsgIDStr {
		lastMsgID, _ = strconv.Atoi(lastMsgIDStr)
		if 0 > lastMsgID {
			lastMsgID = 0
		}
	}

	sentence := lib.SQLSentence(lib.SQLMapSelectMessageHistory)
	pageid, count := lib.GetPageidCount(c)
	rows, err := lib.SQLQuery(sentence, handlers.CommentMsgTypeTalk, lastMsgID, id, talkid, talkid, id, (pageid-1)*count, count)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	defer rows.Close()

	history := cmsTalkHistoryInfo{}
	var info cmsTalkCommentInfo
	var readtmp int
	var timetmp int64

	_, userinfo := handlers.GetUserInfoByID(id)
	if "" != userinfo.IconURL {
		history.UserPic = userinfo.IconURL + CmsLittleImageView
	}

	history.UserName = userinfo.Name

	_, userinfo = handlers.GetUserInfoByID(talkid)
	if "" != userinfo.IconURL {
		history.TalkerPic = userinfo.IconURL + CmsLittleImageView
	}

	history.TalkerName = userinfo.Name

	for rows.Next() {
		err = rows.Scan(&info.MsgID, &info.FromID, &info.ToID, &readtmp, &timetmp, &info.MsgText)
		if nil == err {
			info.TimeUTC = lib.Int64ToUTCTime(timetmp)
			history.Comments = append(history.Comments, info)
		}
	}

	c.JSON(http.StatusOK, history)
}

// DoTalk .
func DoTalk(c *gin.Context) {
	fromidstr := c.Query("fromid")
	toidstr := c.Query("toid")
	msg := c.Query("msg")

	if fromidstr == "" || toidstr == "" || msg == "" {
		c.Status(http.StatusNotFound)
		return
	}

	fromid, _ := strconv.Atoi(fromidstr)
	toid, _ := strconv.Atoi(toidstr)

	// 首先要删掉要自动推送的消息
	handlers.RemoveCommentToPush(fromid, toid)

	// 执行推送
	timevalue := lib.CurrentTimeUTCInt64()
	lastid, _ := handlers.RecommendInsertMessageToDB(fromid, toid, handlers.CommentMsgTypeTalk, msg, timevalue)
	handlers.RecommendPushMessage(fromid, toid, 0, 1, push.PushMsgComment, msg, timevalue)
	push.DoPush()

	var info cmsTalkCommentInfo
	info.MsgID = lastid
	info.FromID = fromid
	info.ToID = toid
	info.MsgText = msg
	info.TimeUTC = lib.Int64ToUTCTime(timevalue)

	c.JSON(http.StatusOK, info)
}

// MessagePushSet 消息推送
func MessagePushSet(c *gin.Context) {
	typestr := c.Query("type")
	msg := c.Query("msg")

	if typestr == "" || msg == "" {
		c.Status(http.StatusNotFound)
		return
	}

	t, _ := strconv.Atoi(typestr)
	if CmsPushMsgTypeEvaluation == t {
		// 配置评价消息推送
		enableStr := c.Query("enable")
		enable := func() bool {
			if enableStr == "1" {
				return true
			}

			return false
		}()

		handlers.PeriodOnlineCommentSet(enable, msg)
	} else {
		// 推送普通消息
		genderstr := c.Query("gender")
		if "" == genderstr {
			c.Status(http.StatusNotFound)
			return
		}
	}

	c.Status(http.StatusOK)
	return
}

// AdminChartsList .
func AdminChartsList(c *gin.Context) {
	var searchInfo cmsSearchInfo
	countsentence := "select count(distinct fromid) from recommend where toid=1"
	err := lib.SQLQueryRow(countsentence).Scan(&searchInfo.Count)
	if nil == err && 0 != searchInfo.Count {
		sentence := "select distinct fromid from recommend where toid=1 order by fromid desc limit ?,?"
		page, count := lib.GetPageidCount(c)
		rows, err := lib.SQLQuery(sentence, (page-1)*count, count)
		if nil != err {
			c.Status(http.StatusNotFound)
			return
		}
		defer rows.Close()

		var info cmsUserInfo
		for rows.Next() {
			rows.Scan(&info.ID)
			code, userinfo := handlers.GetUserInfoByID(info.ID)
			if 200 == code {
				info.Name = userinfo.Name
				info.Age = userinfo.Age
				info.Img = userinfo.IconURL
				info.Province = userinfo.Province
				info.VipLevel = userinfo.VipLevel

				searchInfo.Users = append(searchInfo.Users, info)
			}
		}
	} else if nil != err {
		lib.SQLError(countsentence, err, nil)
	}

	c.JSON(http.StatusOK, searchInfo)
}
