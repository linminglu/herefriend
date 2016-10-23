package handlers

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"herefriend/common"
	"herefriend/config"
	"herefriend/lib"
	"herefriend/lib/push"
)

/*
 * 聊天消息类型,存放到数据库中
 */
const (
	// CommentMsgTypeGreet 打招呼
	CommentMsgTypeGreet = 1
	// CommentMsgTypeTalk 普通聊天信息
	CommentMsgTypeTalk = 2
	// CommentMsgTypeHeart 心动消息
	CommentMsgTypeHeart = 3
	// CommentMsgTypeAskMsg 请求上传图片...
	CommentMsgTypeAskMsg = 4
	// CommentMaxUnreadNum .
	CommentMaxUnreadNum = 20
)

var gRecommendReg *regexp.Regexp
var gRecommendRegUser *regexp.Regexp
var gRecommendQueue *list.List
var gRecommendQueueLock sync.RWMutex
var gRecommendInUnreadCount string
var gRecommendNumber int
var gCountAPIRecommend int
var gMsgTemplates []string
var gRobotURL string

func init() {
	gRecommendReg, _ = regexp.Compile("(?:#)([^#]+)(?:#)")
	gRecommendRegUser, _ = regexp.Compile("(?:#USER_)([^#]+)(?:#)")
	gRecommendInUnreadCount = lib.SQLSentence(lib.SQLMapSelectUnreadMessageCount)
	gRecommendQueue = list.New()

	sentence := lib.SQLSentence(lib.SQLMapSelectAllRecommendCount)
	lib.SQLQueryRow(sentence).Scan(&gRecommendNumber)

	ReloadRecommendTemplates()
	InitRobotURL()

	// start push workroutine here
	push.InitPush()
	go recommendPushRoutine()
	go recommendRobotRoutine()
}

// GetRecommendNumber .
func GetRecommendNumber() int {
	return gRecommendNumber
}

// InitRobotURL .
func InitRobotURL() {
	lib.SQLQueryRow("select url from robotURL where id=1").Scan(&gRobotURL)
}

// ReloadRecommendTemplates .
func ReloadRecommendTemplates() {
	sentence := lib.SQLSentence(lib.SQLMapSelectAllMsgTemplate)
	rows, err := lib.SQLQuery(sentence, 0, -1)
	if nil != err {
		return
	}

	defer rows.Close()

	gMsgTemplates = make([]string, 0)

	template := ""
	for rows.Next() {
		err = rows.Scan(&template)
		if nil == err {
			gMsgTemplates = append(gMsgTemplates, template)
		}
	}
}

// RecommendInsertMessageToDB insert new comment to database
func RecommendInsertMessageToDB(fromid, toid int, msgtype int, msg string, timevalue int64) (int, error) {
	sentence := lib.SQLSentence(lib.SQLMapInsertRecomment)
	result, err := lib.SQLExec(sentence, fromid, toid, timevalue, msgtype, msg)
	if nil != err {
		return -1, err
	}

	lastid, err := result.LastInsertId()
	if nil != err {
		return -1, err
	}

	gRecommendNumber = gRecommendNumber + 1

	return int(lastid), nil
}

// RecommendPushMessage 添加一条消息
func RecommendPushMessage(fromid, toid int, fromusertyp, tousertyp int, pushtype int, msg string, timevalue int64) {
	if 1 == tousertyp {
		clientid := GetClientIDByUserID(toid)

		var title string

		if 1 == fromid {
			title = "来自管理员的消息"
		} else {
			_, gender := getGenderByID(fromid)
			_, userinfo := GetUserInfo(fromid, gender)

			switch pushtype {
			case push.PushMsgVisit:
				title = "有新的来访者"
				msg = "[" + userinfo.Name + "] 查看了你的资料"
			case push.PushMsgComment:
				title = "收到 [" + userinfo.Name + "] 的消息"
			}
		}

		recommendCount := recommendGetUnreadNum(toid, 0)
		visitCount := visitGetUnreadNum(toid, 0)
		push.Add(recommendCount+visitCount, clientid, pushtype, 0, title, msg)
	}

	//虚拟用户给该用户发普通信(打招呼、回招呼、回消息)后,25%概率出现在此,访问时间为发信时间±1分钟随机
	if push.PushMsgVisit != pushtype && common.UserTypeRobot == fromusertyp && 1 != fromid && true == lib.RandomHitPercent(25) {
		visitAddVisitor(toid, fromid, timevalue-60, timevalue+60)
	}
}

type reobotResponse struct {
	Response string `json:"response"`
	Result   int    `json:"result"`
}

func getPostMessageBySessionID(sessionid string, msg string) (int, string) {
	msg = url.QueryEscape(msg)
	buffer, err := lib.GetResultByMethod("GET", gRobotURL+msg, nil)
	if nil != err {
		return 404, ""
	}

	var robotResp reobotResponse
	err = json.Unmarshal(buffer, &robotResp)
	if nil == err && 100 == robotResp.Result {
		return 200, robotResp.Response
	}

	return 404, ""
}

func getFormatMsg(id int) (bool, string) {
	if 0 == len(gMsgTemplates) {
		return false, ""
	}

	replyStr := gMsgTemplates[lib.Intn(len(gMsgTemplates))]

	kws := gRecommendReg.FindAllStringSubmatch(replyStr, -1)
	if nil != kws {
		keywords := make(map[string]bool)

		for _, kw := range kws {
			keywords[kw[1]] = true
		}

		_, info := GetUserInfo(id, 0)
		table := reflect.ValueOf(info)

		for word := range keywords {
			el := table.FieldByName(word)

			switch el.Kind() {
			case reflect.Int:
				tmp := el.Int()
				if 0 == tmp {
					return false, ""
				}

				replyStr = strings.Replace(replyStr, "#"+word+"#", strconv.FormatInt(tmp, 10), -1)
			case reflect.String:
				s := el.String()
				if "" == s || true == strings.Contains(s, "以后告诉你") {
					return false, ""
				}

				replyStr = strings.Replace(replyStr, "#"+word+"#", s, -1)
			default:
				break
			}
		}
	}

	return true, replyStr
}

func tulingResponseCheck(str string, id int) bool {
	//长度太长返回失败
	if len(str) > 80 {
		return false
	}

	for _, s := range gRobotResponseCheckList {
		if true == strings.Contains(str, s) {
			return false
		}
	}

	//如果回复过相同的话，不再回复
	var count int
	sentence := lib.SQLSentence(lib.SQLMapSelectHaveSameReply)
	lib.SQLQueryRow(sentence, id, str).Scan(&count)
	if 0 != count {
		return false
	}

	return true
}

func tulingResponseChange(id int) string {
	ok, replacestr := getFormatMsg(id)
	if !ok {
		ok, replacestr = getFormatMsg(id)
	}

	if !ok {
		replacestr = gHelloArray[lib.Intn(len(gHelloArray))]
	}

	return replacestr
}

/*
 *
 *    Function: getResponseJson
 *      Author: sunchao
 *        Date: 15/9/20
 * Description: 获取回复信息
 *
 */
func getResponseMsg(node *recommendQueueNode) string {
	sessionid := strconv.Itoa(node.fromid) + strconv.Itoa(node.toid)
	code, message := getPostMessageBySessionID(sessionid, node.message)

	if 200 == code {
		if true != tulingResponseCheck(message, node.fromid) {
			message = tulingResponseChange(node.toid)
		}
	} else {
		message = tulingResponseChange(node.toid)
	}

	return message
}

// RemoveCommentToPush 删除要推送的消息
func RemoveCommentToPush(fromid, toid int) {
	gRecommendQueueLock.Lock()
	for e := gRecommendQueue.Front(); e != nil; {
		n := e.Value.(*recommendQueueNode)
		if n.fromid == toid && n.toid == fromid {
			next := e.Next()
			gRecommendQueue.Remove(e)
			e = next
		} else {
			e = e.Next()
		}
	}
	gRecommendQueueLock.Unlock()
}

/*
 *
 *    Function: recommendPushRoutine
 *      Author: sunchao
 *        Date: 15/11/2
 * Description: routine to push the messages on the queue
 *
 */
func recommendPushRoutine() {
	needpush := false
	for {
		time.Sleep(lib.SleepDurationPushQueuMsg)
		needpush = false

		gRecommendQueueLock.Lock()
		for e := gRecommendQueue.Front(); e != nil; {
			n := e.Value.(*recommendQueueNode)
			n.timewait = n.timewait - 1

			// 可以回复消息了
			if 0 >= n.timewait {
				needreplay := true

				// 通过几率过滤
				if 1 == n.msgtype {
					count := 0
					sentence := lib.SQLSentence(lib.SQLMapSelectRecommendCount)
					lib.SQLQueryRow(sentence, 1, n.fromid, n.toid).Scan(&count)
					if 0 < count && true == lib.RandomHitPercent(25) {
						needreplay = false
					}
				}

				if needreplay {
					timevalue := lib.CurrentTimeUTCInt64()
					msg := getResponseMsg(n)
					if "" != msg {
						// 多次发送消息，则自动设置为VIP
						if true == checkAlreadySendSameCommentToday(n.toid, n.fromid, CommentMsgTypeTalk) {
							_, gender := getGenderByID(n.toid)
							if true != checkIfUserHaveViplevel(n.toid, gender) {
								changevipsentence := func() string {
									if 0 == gender {
										return "update girls set viplevel=1 where id=?"
									}
									return "update guys set viplevel=1 where id=?"
								}()
								lib.SQLExec(changevipsentence, n.toid)
							}
						}

						RecommendInsertMessageToDB(n.toid, n.fromid, n.msgtype, msg, timevalue)
						RecommendPushMessage(n.toid, n.fromid, n.tousertype, n.fromusertype, push.PushMsgComment, msg, timevalue)
						needpush = true
					}
				}

				next := e.Next()
				gRecommendQueue.Remove(e)
				e = next
			} else {
				e = e.Next()
			}
		}
		gRecommendQueueLock.Unlock()

		if true == needpush {
			push.DoPush()
		}
	}
}

/*
 *
 *    Function: recommendRobotRoutine
 *      Author: sunchao
 *        Date: 15/11/4
 * Description: 发送消息线程
 *
 */
func recommendRobotRoutine() {
	var count int

	needpush := false
	for {
		time.Sleep(lib.SleepTimeDuration(lib.SleepTypeRobotComment))
		needpush = false

		gLiveUsersInfo.lock.RLock()
		for id, user := range gLiveUsersInfo.users {
			count = 0
			err := lib.SQLQueryRow(gRecommendInUnreadCount, id, 0).Scan(&count)
			if nil != err || CommentMaxUnreadNum <= count {
				if nil != err {
					lib.SQLError(gRecommendInUnreadCount, err, id, 0)
				}
				continue
			}

			// 付费用户每次100%几率，未付费用户（有照片：70%，无照片：40%)
			if false == checkIfUserHaveViplevel(id, user.gender) {
				if true == checkIfUserHavePicture(id, user.gender) {
					if lib.RandomHitPercent(35) {
						continue
					}
				} else {
					if lib.RandomHitPercent(50) {
						continue
					}
				}
			}

			/* get the random id */
			fromid := getRandomHeartbeatID(id, 1-user.gender)
			if 0 == fromid {
				continue
			} else {
				_, usertype := GetUsertypeByIDGender(id, 1-user.gender)
				if common.UserTypeUser == usertype {
					continue
				}
			}

			ok, msg := getFormatMsg(fromid)
			if !ok {
				msg = gHelloArray[lib.Intn(len(gHelloArray))]
			}

			// 多次发送消息，则自动设置为VIP
			if true == checkAlreadySendSameCommentToday(fromid, id, CommentMsgTypeTalk) {
				if true != checkIfUserHaveViplevel(fromid, 1-user.gender) {
					changevipsentence := func() string {
						if 0 == 1-user.gender {
							return "update girls set viplevel=1 where id=?"
						}
						return "update guys set viplevel=1 where id=?"
					}()

					lib.SQLExec(changevipsentence, fromid)
				}
			}

			timevalue := lib.CurrentTimeUTCInt64()
			RecommendInsertMessageToDB(fromid, id, CommentMsgTypeTalk, msg, timevalue)
			RecommendPushMessage(fromid, id, 0, 1, push.PushMsgComment, msg, timevalue)
			needpush = true
		}
		gLiveUsersInfo.lock.RUnlock()

		if true == needpush {
			push.DoPush()
		}
	}
}

/*
 *
 *    Function: checkAlreadySendSameCommentToday
 *      Author: sunchao
 *        Date: 15/11/15
 * Description: 检查今天是否打过招呼
 *
 */
func checkAlreadySendSameCommentToday(fromid, toid, msgtype int) bool {
	var timevalue int64
	sentence := lib.SQLSentence(lib.SQLMapSelectCheckCommentDailyLock)
	err := lib.SQLQueryRow(sentence, fromid, toid, msgtype).Scan(&timevalue)
	if nil != err || 0 == timevalue {
		return false
	}

	timeUTC := lib.Int64ToUTCTime(timevalue)
	todyUTC := time.Now().UTC()
	if timeUTC.Year() == todyUTC.Year() && timeUTC.Month() == todyUTC.Month() && timeUTC.Day() == todyUTC.Day() {
		return true
	}

	return false
}

/*
 |    Function: incomeGreetCommentProc
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 打招呼消息处理
 |      Return: int->httpcode, string->返回消息, bool->是否推送消息
 |
*/
func incomeGreetCommentProc(id, gender, toid, togender int, msg string, timevalue int64) (int, int, string, bool) {
	if gender == togender {
		return 403, -1, "抱歉,同性之间不能打招呼.", false
	}

	if true == checkAlreadySendSameCommentToday(id, toid, CommentMsgTypeGreet) {
		return 403, -1, "抱歉,一天只能向同一个人打招呼一次.", false
	}

	RecommendInsertMessageToDB(id, toid, CommentMsgTypeGreet, "", timevalue)
	ok, replymsg := getFormatMsg(id)
	if !ok {
		replymsg = gHelloArray[lib.Intn(len(gHelloArray))]
	}

	lastid, err := RecommendInsertMessageToDB(id, toid, CommentMsgTypeTalk, replymsg, timevalue)
	if nil != err {
		return 404, -1, err.Error(), false
	}

	//后续推送自动选择的消息
	return 200, lastid, replymsg, true
}

/*
 |    Function: incomeHeartCommentProc
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 心动消息处理
 |      Return: int->httpcode, string->返回消息, bool->是否推送消息
 |
*/
func incomeHeartCommentProc(id, gender, toid, togender int, msg string, timevalue int64) (int, int, string, bool) {
	if gender == togender {
		return 403, -1, "抱歉,同性之间不能发心动消息.", false
	}

	if true == checkAlreadySendSameCommentToday(id, toid, CommentMsgTypeHeart) {
		return 403, -1, "抱歉,一天只能向同一个人发送一次心动消息.", false
	}

	lastid, err := RecommendInsertMessageToDB(id, toid, CommentMsgTypeHeart, "", timevalue)
	if nil != err {
		return 404, -1, err.Error(), false
	}

	return 200, lastid, "", false
}

/*
 |    Function: incomeAskCommentProc
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 邀请(索要)消息处理
 |      Return: int->httpcode, string->返回消息, bool->是否推送消息
 |
*/
func incomeAskCommentProc(id, gender, toid, togender int, msg string, timevalue int64) (int, int, string, bool) {
	if gender == togender {
		return 403, -1, "抱歉,同性之间不能发消息.", false
	}

	if true == checkAlreadySendSameCommentToday(id, toid, CommentMsgTypeAskMsg) {
		return 403, -1, "抱歉,一天只能向同一个人邀请一次.", false
	}

	RecommendInsertMessageToDB(id, toid, CommentMsgTypeAskMsg, "", timevalue)
	lastid, err := RecommendInsertMessageToDB(id, toid, CommentMsgTypeTalk, msg, timevalue)
	if nil != err {
		return 404, -1, err.Error(), false
	}

	//后续推送自动选择的消息
	return 200, lastid, msg, true
}

/*
 |    Function: incomeTalkCommentProc
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 聊天消息处理
 |      Return: int->httpcode, string->返回消息, bool->是否推送消息
 |
*/
func incomeTalkCommentProc(id, gender, toid, togender int, msg string, timevalue int64) (int, int, string, bool) {
	if 1 != toid && gender == togender {
		return 403, -1, "抱歉,同性之间不能发消息.", false
	}

	lastid, err := RecommendInsertMessageToDB(id, toid, CommentMsgTypeTalk, msg, timevalue)
	if nil != err {
		return 404, -1, err.Error(), false
	}

	if 1 == toid {
		return 200, lastid, msg, false
	}
	return 200, lastid, msg, true
}

/*
 |    Function: incomeCommentPushMsgProc
 |      Author: Mr.Sancho
 |        Date: 2016-04-25
 | Description: 推送消息处理
 |      Return:
 |
*/
func incomeCommentPushMsgProc(id, gender, toid, msgtype int, msg string, timevalue int64) {
	if 1 == toid {
		return
	}

	_, _, tousertype := GetGenderUsertypeByID(toid)
	if common.UserTypeUser == tousertype {
		/*
		 * 注册用户直接发送
		 * 这里发送者一定是注册用户
		 */
		RecommendPushMessage(id, toid, common.UserTypeUser, common.UserTypeUser, push.PushMsgComment, msg, timevalue)
		push.DoPush()
	} else {
		/*
		 * 自动回复缓存, 例外如下
		 * 1.索要照片等信息除外
		 * 2.如果发消息用户为VIP用户,也不会缓存
		 */
		if CommentMsgTypeAskMsg == msgtype {
			return
		}

		if true == checkIfUserHaveViplevel(id, gender) {
			return
		}

		gRecommendQueueLock.Lock()
		gRecommendQueue.PushBack(&recommendQueueNode{
			timewait:     int64(lib.SleepTimeDuration(lib.SleepTypeRobotReply) / time.Second),
			fromid:       id,
			toid:         toid,
			fromusertype: common.UserTypeUser,
			tousertype:   common.UserTypeRobot,
			msgtype:      CommentMsgTypeTalk,
			message:      msg,
			timevalue:    timevalue})
		gRecommendQueueLock.Unlock()
	}
}

// ActionRecommend 打招呼
func ActionRecommend(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusForbidden)
		return
	}

	toidStr := c.Query("toid")
	typeStr := c.Query("type")
	if "" == toidStr || "" == typeStr {
		c.Status(http.StatusNotFound)
		return
	}

	msgtype, _ := strconv.Atoi(typeStr)
	toid, _ := strconv.Atoi(toidStr)

	/* the dest id not exist */
	exist, togender := getGenderByID(toid)
	if true != exist {
		c.Status(http.StatusForbidden)
		return
	}

	msg := c.Query("msg")
	t := time.Now()
	timevalue := lib.TimeToUTCInt64(t)

	var code, lastid int
	var replymsg string
	var bpush bool

	switch msgtype {
	case CommentMsgTypeGreet:
		code, lastid, replymsg, bpush = incomeGreetCommentProc(id, gender, toid, togender, msg, timevalue)
	case CommentMsgTypeHeart:
		code, lastid, replymsg, bpush = incomeHeartCommentProc(id, gender, toid, togender, msg, timevalue)
	case CommentMsgTypeAskMsg:
		code, lastid, replymsg, bpush = incomeAskCommentProc(id, gender, toid, togender, msg, timevalue)
	case CommentMsgTypeTalk:
		code, lastid, replymsg, bpush = incomeTalkCommentProc(id, gender, toid, togender, msg, timevalue)
	default:
		c.Status(http.StatusNotFound)
		return
	}

	//不继续处理
	if 200 != code {
		c.String(code, replymsg)
		return
	}

	if true == bpush {
		go incomeCommentPushMsgProc(id, gender, toid, msgtype, replymsg, timevalue)
	}

	gCountAPIRecommend = gCountAPIRecommend + 1

	c.JSON(http.StatusOK, messageInfo{
		MsgID:     lastid,
		MsgText:   replymsg,
		UserID:    toid,
		Direction: MessageDirectionFromMe,
		Readed:    false,
		TimeUTC:   t,
	})
}

// DelRecommend delete the recommend
func DelRecommend(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusForbidden)
		return
	}

	msgidStr := c.Query("msgid")
	if msgidStr == "" {
		c.Status(http.StatusForbidden)
		return
	}

	talkidStr := c.Query("talkid")
	if "" == talkidStr {
		c.Status(http.StatusForbidden)
		return
	}

	/* the dest id not exist */
	msgid, _ := strconv.Atoi(msgidStr)
	talkid, _ := strconv.Atoi(talkidStr)
	sentence := lib.SQLSentence(lib.SQLMapDeleteRecommend)
	lib.SQLExec(sentence, msgid, talkid, id, id, talkid)
	gRecommendNumber = gRecommendNumber - 1

	c.Status(http.StatusOK)
}

/*
 |    Function: getRecommendByRows
 |      Author: Mr.Sancho
 |        Date: 2016-01-02
 |   Arguments:
 |      Return:
 | Description: 根据sql行数据获取所有的聊天信息
 |
*/
func getRecommendByRows(id int, rows *sql.Rows) []messageInfo {
	var info messageInfo
	var fromid int
	var toid int
	var readtmp int
	var timetmp int64
	var err error

	var infos []messageInfo
	for rows.Next() {
		err = rows.Scan(&info.MsgID, &fromid, &toid, &readtmp, &timetmp, &info.MsgText)
		if nil == err {
			info.TimeUTC = lib.Int64ToUTCTime(timetmp)

			if fromid == id || 1 == readtmp {
				info.Readed = true
			} else {
				info.Readed = false
			}

			if fromid == id {
				info.UserID = toid
				info.Direction = MessageDirectionFromMe
			} else {
				info.UserID = fromid
				info.Direction = MessageDirectionToMe
			}

			infos = append(infos, info)
		} else {
			log.Error(err.Error())
		}
	}

	return infos
}

// GetWaterFlow 获取聊天流水
func GetWaterFlow(c *gin.Context) {
	talkidStr := c.Query("talkid")
	if talkidStr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	exist, id, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	var lastMsgID int
	lastMsgIDStr := c.Query("lastmsgid")
	if "" != lastMsgIDStr {
		lastMsgID, _ = strconv.Atoi(lastMsgIDStr)
		if 0 > lastMsgID {
			lastMsgID = 0
		}
	}

	pageid, count := lib.GetPageidCount(c)
	talkid, _ := strconv.Atoi(talkidStr)
	exist, _ = getGenderByID(talkid)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMapSelectMessageHistory)
	rows, err := lib.SQLQuery(sentence, CommentMsgTypeTalk, lastMsgID, id, talkid, talkid, id, (pageid-1)*count, count)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	defer rows.Close()

	infos := getRecommendByRows(id, rows)
	if len(infos) > 0 {
		sentence := lib.SQLSentence(lib.SQLMapUpdateRecommendRead)
		lib.SQLExec(sentence, talkid, id, infos[0].MsgID)
	}

	c.JSON(http.StatusOK, infos)
}

// GetRecommendAll 获取所有的聊天信息
func GetRecommendAll(timeline int64, id, pageid, count int) ([]messageInfo, error) {
	sentence := lib.SQLSentence(lib.SQLMapSelectDistinctRecommend)
	rows, err := lib.SQLQuery(sentence, id, CommentMsgTypeTalk, timeline, id, CommentMsgTypeTalk, timeline, (pageid-1)*count, count)
	if nil == err {
		defer rows.Close()
		var gender int

		infos := getRecommendByRows(id, rows)
		for i := range infos {
			_, gender = getGenderByID(infos[i].UserID)
			_, userinfo := GetUserInfo(infos[i].UserID, gender)
			infos[i].UserInfo = &userinfo
		}

		return infos, nil
	}

	return nil, err
}

// GetAllMessage 获取所有的聊天信息和谁看过我信息
func GetAllMessage(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var timeline int64
	timelinestr := c.Query("lasttime")
	if "" != timelinestr {
		timeline = lib.TimeStrToUTCInt64(timelinestr)
	}

	var allmessage allMessageInfo
	pageid, count := lib.GetPageidCount(c)

	recommendAlls, err := GetRecommendAll(timeline, id, pageid, count)
	if nil == err {
		allmessage.RecommendArray = recommendAlls
	}

	visitAlls, err := getVisitAll(timeline, id, pageid, count)
	if nil == err {
		allmessage.VisitArray = visitAlls
	}

	c.JSON(http.StatusOK, allmessage)
}

// GetComments .
func GetComments(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	var timeline int64
	timelinestr := c.Query("lasttime")
	if timelinestr != "" {
		timeline = lib.TimeStrToUTCInt64(timelinestr)
	}

	pageid, count := lib.GetPageidCount(c)
	recommendAlls, err := GetRecommendAll(timeline, id, pageid, count)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, recommendAlls)
}

// GetVisits .
func GetVisits(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var timeline int64
	timelinestr := c.Query("lasttime")
	if "" != timelinestr {
		timeline = lib.TimeStrToUTCInt64(timelinestr)
	}

	pageid, count := lib.GetPageidCount(c)
	visitAlls, err := getVisitAll(timeline, id, pageid, count)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, visitAlls)
}

// GetUnreadMessage get the unread message
func GetUnreadMessage(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var timeline int64

	timelinestr := c.Query("lasttime")
	if "" != timelinestr {
		timeline = lib.TimeStrToUTCInt64(timelinestr)
	}

	recommendCount := recommendGetUnreadNum(id, timeline)
	visitCount := visitGetUnreadNum(id, timeline)
	unreadmsg := unreadMessageInfo{UnreadRecommend: recommendCount, UnreadVisit: visitCount, Badge: recommendCount + visitCount}

	c.JSON(http.StatusOK, unreadmsg)
}

func recommendGetUnreadNum(id int, timeline int64) int {
	var count int

	sentence := lib.SQLSentence(lib.SQLMapSelectUnreadMessageCount)
	lib.SQLQueryRow(sentence, id, timeline).Scan(&count)

	return count
}

// GetAPIRecommendCount .
func GetAPIRecommendCount() int {
	return gCountAPIRecommend
}

// PeriodOnlineCommentSet 定期推送消息设置
func PeriodOnlineCommentSet(enable bool, msg string) {
	config.ConfEvaluationSwitch = enable
	config.ConfEvaluationMsgContent = msg
}

// PeriodOnlineCommentPush 定期推送消息处理
func PeriodOnlineCommentPush(id, gender int, lastEvaluationTime int64) {
	curTime := lib.CurrentTimeUTCInt64()
	if 0 == lastEvaluationTime || 43200 <= (curTime-lastEvaluationTime) {
		evaluationMsg := PushMsgEvaluation{Enable: true, ShowMessage: config.ConfEvaluationMsgContent}
		jsonRlt, _ := json.Marshal(evaluationMsg)
		notifymsg := PushMessageInfo{Type: push.NotifyMsgEvaluation, Value: string(jsonRlt)}
		jsonRlt, _ = json.Marshal(notifymsg)

		push.Add(0, GetClientIDByUserID(id), push.PushMsgNotify, push.NotifyMsgEvaluation, "", string(jsonRlt))
		push.DoPush()

		sentence := lib.SQLSentence(lib.SQLMapUpdateEvaluationTime, gender)
		lib.SQLExec(sentence, curTime, id)
	}

	return
}
