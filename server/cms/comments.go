package cms

import (
	"encoding/json"
	"net/http"
	"strconv"

	"herefriend/lib"
	"herefriend/lib/push"
	"herefriend/server/handlers"
)

const (
	CMS_PUSHMSG_TYPE_NORMALMSG  = 1
	CMS_PUSHMSG_TYPE_EVALUATION = 2
)

const CMS_LittleImgView = "?imageView2/5/w/50/h/50"

/*
 |    Function: CommentInfo
 |      Author: Mr.Sancho
 |        Date: 2016-02-28
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func CommentInfo(req *http.Request) string {
	info := cmsCommentSummary{
		TalkNum:   handlers.GetApiRecommendCount(),
		PushNum:   push.GetPushNum(),
		BuyVIPNum: handlers.GetBuyVIPCount(),
	}

	jsonRlt, _ := json.Marshal(info)
	return string(jsonRlt)
}

/*
 |    Function: Recommendhistory
 |      Author: Mr.Sancho
 |        Date: 2016-02-12
 |   Arguments:
 |      Return:
 | Description: 获取最新的消息
 |
*/
func Recommendhistory(req *http.Request) string {
	var lastmsgid int
	var fromid int
	var toid int

	v := req.URL.Query()
	lastmsgidstr := v.Get("lastmsgid")
	if "" == lastmsgidstr {
		lastmsgid = 0
	} else {
		lastmsgid, _ = strconv.Atoi(lastmsgidstr)
	}

	sentence := "select id,fromid,toid,time,type,msg from recommend where id>? order by id desc limit 20"
	rows, err := lib.SQLQuery(sentence, lastmsgid)
	if nil != err {
		return "[]"
	}
	defer rows.Close()

	var info cmsCommentInfo
	var timevalue int64
	infos := make([]cmsCommentInfo, 0)

	for rows.Next() {
		err = rows.Scan(&info.MsgId, &fromid, &toid, &timevalue, &info.MsgType, &info.MsgText)
		if nil != err {
			continue
		}

		code, userinfo := handlers.GetUserInfoById(fromid)
		if 200 == code && "" != userinfo.Name {
			info.From = "[" + userinfo.Province + "]" + userinfo.Name
		} else {
			info.From = "[" + userinfo.Province + "]" + strconv.Itoa(fromid)
		}

		info.FromPic = userinfo.IconUrl + CMS_LittleImgView

		code, userinfo = handlers.GetUserInfoById(toid)
		if 200 == code && "" != userinfo.Name {
			info.To = "[" + userinfo.Province + "]" + userinfo.Name
		} else {
			info.To = "[" + userinfo.Province + "]" + strconv.Itoa(toid)
		}
		info.ToPic = userinfo.IconUrl + CMS_LittleImgView

		info.TimeUTC = lib.Int64_To_UTCTime(timevalue)
		infos = append(infos, info)
	}

	if len(infos) > 0 {
		len := len(infos)
		for i := 0; i < len/2; i++ {
			infos[i], infos[len-i-1] = infos[len-i-1], infos[i]
		}
	}

	jsonRlt, _ := json.Marshal(infos)
	return string(jsonRlt)
}

/*
 |    Function: MsgTemplate
 |      Author: Mr.Sancho
 |        Date: 2016-02-17
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func MsgTemplate(req *http.Request) string {
	var msgtype int
	var gender int

	v := req.URL.Query()
	msgtypestr := v.Get("type")
	genderstr := v.Get("gender")

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
		return "[]"
	}
	defer rows.Close()

	var info cmsMessageTempalte
	infos := make([]cmsMessageTempalte, 0)

	for rows.Next() {
		err = rows.Scan(&info.Id, &info.Template)
		if nil != err {
			continue
		}

		infos = append(infos, info)
	}

	jsonRlt, _ := json.Marshal(infos)
	return string(jsonRlt)
}

/*
 |    Function: MsgTemplateAdd
 |      Author: Mr.Sancho
 |        Date: 2016-02-17
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func MsgTemplateAdd(req *http.Request) (int, string) {
	var msgtype int
	var gender int

	v := req.URL.Query()
	msgtypestr := v.Get("type")
	genderstr := v.Get("gender")
	templatestr := v.Get("template")

	if "" == templatestr {
		return 404, ""
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
		return 404, ""
	}

	lastid, err := result.LastInsertId()
	if nil != err {
		return 404, ""
	}

	handlers.ReloadRecommendTemplates()

	var info cmsMessageTempalte
	info.Id = int(lastid)
	info.Template = templatestr

	jsonRlt, _ := json.Marshal(info)
	return 200, string(jsonRlt)
}

/*
 |    Function: MsgTemplateDel
 |      Author: Mr.Sancho
 |        Date: 2016-02-17
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func MsgTemplateDel(req *http.Request, w http.ResponseWriter) {
	v := req.URL.Query()
	idstr := v.Get("id")

	if "" == idstr {
		w.WriteHeader(404)
		return
	}

	id, err := strconv.Atoi(idstr)
	if nil != err {
		w.WriteHeader(404)
		return
	}

	sentence := "delete from msgtemplate where id=?"
	_, err = lib.SQLExec(sentence, id)
	if nil != err {
		w.WriteHeader(404)
		return
	}

	handlers.ReloadRecommendTemplates()
	return
}

/*
 |    Function: MsgTemplateModify
 |      Author: Mr.Sancho
 |        Date: 2016-02-17
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func MsgTemplateModify(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	idstr := v.Get("id")
	templatestr := v.Get("template")

	if "" == idstr || "" == templatestr {
		w.WriteHeader(404)
		return
	}

	id, err := strconv.Atoi(idstr)
	if nil != err {
		w.WriteHeader(404)
		return
	}

	sentence := "update msgtemplate set msg=? where id=?"
	_, err = lib.SQLExec(sentence, templatestr, id)
	if nil != err {
		w.WriteHeader(404)
		return
	}

	handlers.ReloadRecommendTemplates()
	return
}

/*
 |
 |    Function: GetChartsList
 |      Author: sunchao
 |        Date: 16/4/4
 | Description: get the charts list
 |
*/
func GetChartsList(w http.ResponseWriter, r *http.Request) (int, string) {
	v := r.URL.Query()

	idstr := v.Get("id")
	if "" == idstr {
		return 404, ""
	}

	id, _ := strconv.Atoi(idstr)

	var commentsInfo []cmsCommentInfo
	pageid, count := lib.Get_pageid_count_fromreq(r)
	recommendAlls, err := handlers.GetRecommendAll(0, id, pageid, count)
	if nil != err {
		return 404, err.Error()
	}

	_, userinfo := handlers.GetUserInfoById(id)
	for _, r := range recommendAlls {
		c := cmsCommentInfo{
			MsgId:     r.MsgId,
			MsgText:   r.MsgText,
			TimeUTC:   r.TimeUTC,
			Direction: r.Direction,
		}

		_, _, usertype := handlers.GetGenderUsertypeById(r.UserId)
		if 1 == usertype {
			continue
		}

		if 1 == c.Direction {
			c.FromId = id
			c.ToId = r.UserId
			c.From = "[" + userinfo.Province + "]" + userinfo.Name
			if "" != userinfo.IconUrl {
				c.FromPic = userinfo.IconUrl + CMS_LittleImgView
			}
			c.To = "[" + r.UserInfo.Province + "]" + r.UserInfo.Name
			if "" != r.UserInfo.IconUrl {
				c.ToPic = r.UserInfo.IconUrl + CMS_LittleImgView
			}
		} else {
			c.FromId = r.UserId
			c.ToId = id
			c.To = "[" + userinfo.Province + "]" + userinfo.Name
			if "" != userinfo.IconUrl {
				c.ToPic = userinfo.IconUrl + CMS_LittleImgView
			}
			c.From = "[" + r.UserInfo.Province + "]" + r.UserInfo.Name
			if "" != r.UserInfo.IconUrl {
				c.FromPic = r.UserInfo.IconUrl + CMS_LittleImgView
			}
		}

		commentsInfo = append(commentsInfo, c)
	}

	jsonRlt, _ := json.Marshal(commentsInfo)
	return 200, string(jsonRlt)
}

/*
 |    Function: GetTalkHistory
 |      Author: Mr.Sancho
 |        Date: 2016-04-12
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func GetTalkHistory(r *http.Request) (int, string) {
	v := r.URL.Query()
	idStr := v.Get("id")
	talkidStr := v.Get("talkid")
	if "" == talkidStr {
		return 404, ""
	}

	id, _ := strconv.Atoi(idStr)
	talkid, _ := strconv.Atoi(talkidStr)

	var lastMsgId int
	lastMsgIdStr := v.Get("lastmsgid")
	if "" != lastMsgIdStr {
		lastMsgId, _ = strconv.Atoi(lastMsgIdStr)
		if 0 > lastMsgId {
			lastMsgId = 0
		}
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Select_MessageHistory)
	pageid, count := lib.Get_pageid_count_fromreq(r)
	rows, err := lib.SQLQuery(sentence, handlers.RECOMMEND_MSGTYPE_TALK, lastMsgId, id, talkid, talkid, id, (pageid-1)*count, count)
	if nil != err {
		return 404, ""
	}

	defer rows.Close()

	history := cmsTalkHistoryInfo{}
	var info cmsTalkCommentInfo
	var readtmp int
	var timetmp int64

	_, userinfo := handlers.GetUserInfoById(id)
	if "" != userinfo.IconUrl {
		history.UserPic = userinfo.IconUrl + CMS_LittleImgView
	}

	history.UserName = userinfo.Name

	_, userinfo = handlers.GetUserInfoById(talkid)
	if "" != userinfo.IconUrl {
		history.TalkerPic = userinfo.IconUrl + CMS_LittleImgView
	}

	history.TalkerName = userinfo.Name

	for rows.Next() {
		err = rows.Scan(&info.MsgId, &info.FromId, &info.ToId, &readtmp, &timetmp, &info.MsgText)
		if nil == err {
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)
			history.Comments = append(history.Comments, info)
		}
	}

	jsonRlt, _ := json.Marshal(history)
	return 200, string(jsonRlt)
}

/*
 |    Function: DoTalk
 |      Author: Mr.Sancho
 |        Date: 2016-04-13
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func DoTalk(w http.ResponseWriter, r *http.Request) (int, string) {
	v := r.URL.Query()
	fromidstr := v.Get("fromid")
	toidstr := v.Get("toid")
	msg := v.Get("msg")

	if "" == fromidstr || "" == toidstr || "" == msg {
		return 404, ""
	}

	fromid, _ := strconv.Atoi(fromidstr)
	toid, _ := strconv.Atoi(toidstr)

	// 首先要删掉要自动推送的消息
	handlers.RemoveCommentToPush(fromid, toid)

	// 执行推送
	timevalue := lib.CurrentTimeUTCInt64()
	lastid, _ := handlers.RecommendInsertMessageToDB(fromid, toid, handlers.RECOMMEND_MSGTYPE_TALK, msg, timevalue)
	handlers.RecommendPushMessage(fromid, toid, 0, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
	push.DoPush()

	var info cmsTalkCommentInfo
	info.MsgId = lastid
	info.FromId = fromid
	info.ToId = toid
	info.MsgText = msg
	info.TimeUTC = lib.Int64_To_UTCTime(timevalue)

	jsonRlt, _ := json.Marshal(info)
	return 200, string(jsonRlt)
}

/*
 |    Function: MessagePushSet
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 消息推送
 |      Return:
 |
*/
func MessagePushSet(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	typestr := v.Get("type")
	msg := v.Get("msg")

	if "" == typestr || "" == msg {
		w.WriteHeader(404)
		return
	}

	t, _ := strconv.Atoi(typestr)
	if CMS_PUSHMSG_TYPE_EVALUATION == t {
		// 配置评价消息推送
		enableStr := v.Get("enable")
		enable := func() bool {
			if "1" == enableStr {
				return true
			} else {
				return false
			}
		}()

		handlers.PeriodOnlineCommentSet(enable, msg)
	} else {
		// 推送普通消息
		genderstr := v.Get("gender")
		if "" == genderstr {
			w.WriteHeader(404)
			return
		}
	}

	w.WriteHeader(200)
	return
}
