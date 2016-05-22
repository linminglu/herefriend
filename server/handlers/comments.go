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

	"herefriend/common"
	"herefriend/lib"
	"herefriend/lib/push"
)

/*
 * 聊天消息类型,存放到数据库中
 */
const (
	RECOMMEND_MSGTYPE_GREET  = 1 //打招呼
	RECOMMEND_MSGTYPE_TALK   = 2 //普通聊天信息
	RECOMMEND_MSGTYPE_HEART  = 3 //心动消息
	RECOMMEND_MSGTYPE_ASKMSG = 4 //请求上传图片...
)

const (
	RECOMMEND_MAX_UNREADNUMBER = 20
)

var gRecommendReg *regexp.Regexp
var gRecommendRegUser *regexp.Regexp
var gRecommendQueue *list.List
var gRecommendQueueLock sync.RWMutex
var gRecommendInUnreadCount string
var gRecommendNumber int
var gCountApiRecommend int
var gMsgTemplates []string
var gRobotUrl string

/*
 * 定期推送消息
 */
var gEvaluationMsgContent = "好评送免费VIP 3个月哦，赶紧去评价一下吧。"
var gEnableEvaluation = true

func init() {
	gRecommendReg, _ = regexp.Compile("(?:#)([^#]+)(?:#)")
	gRecommendRegUser, _ = regexp.Compile("(?:#USER_)([^#]+)(?:#)")
	gRecommendInUnreadCount = lib.SQLSentence(lib.SQLMAP_Select_UnreadMessageCount)
	gRecommendQueue = list.New()

	sentence := lib.SQLSentence(lib.SQLMAP_Select_AllRecommendCount)
	lib.SQLQueryRow(sentence).Scan(&gRecommendNumber)

	ReloadRecommendTemplates()
	InitRobotUrl()

	// start push workroutine here
	push.InitPush()
	go recommendPushRoutine()
	go recommendRobotRoutine()
}

func GetRecommendNumber() int {
	return gRecommendNumber
}

func InitRobotUrl() {
	lib.SQLQueryRow("select url from robotURL where id=1").Scan(&gRobotUrl)
}

func ReloadRecommendTemplates() {
	sentence := lib.SQLSentence(lib.SQLMAP_Select_AllMsgTemplate)
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

/*
 |    Function: RecommendInsertMessageToDB
 |      Author: Mr.Sancho
 |        Date: 2016-04-13
 |   Arguments:
 |      Return:
 | Description: insert new comment to database
 |
*/
func RecommendInsertMessageToDB(fromid, toid int, msgtype int, msg string, timevalue int64) (int, error) {
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Recomment)
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

/*
 *
 *    Function: RecommendPushMessage
 *      Author: sunchao
 *        Date: 15/11/2
 * Description: 添加一条消息
 *
 */
func RecommendPushMessage(fromid, toid int, fromusertyp, tousertyp int, pushtype int, msg string, timevalue int64) {
	if 1 == tousertyp {
		clientid := GetClientIdByUserId(toid)

		var title string

		if 1 == fromid {
			title = "来自管理员的消息"
		} else {
			_, gender := getGenderById(fromid)
			_, userinfo := GetUserInfo(fromid, gender)

			switch pushtype {
			case push.PUSHMSG_TYPE_VISIT:
				title = "有新的来访者"
				msg = "[" + userinfo.Name + "] 查看了你的资料"
			case push.PUSHMSG_TYPE_RECOMMEND:
				title = "收到 [" + userinfo.Name + "] 的消息"
			}
		}

		recommendCount := recommend_GetUnreadNum(toid, 0)
		visitCount := visit_GetUnreadNum(toid, 0)
		push.Add(recommendCount+visitCount, clientid, pushtype, 0, title, msg)
	}

	//虚拟用户给该用户发普通信(打招呼、回招呼、回消息)后,25%概率出现在此,访问时间为发信时间±1分钟随机
	if push.PUSHMSG_TYPE_VISIT != pushtype && common.USERTYPE_RB == fromusertyp && 1 != fromid && true == lib.RandomHitPercent(25) {
		visitAddVisitor(toid, fromid, timevalue-60, timevalue+60)
	}
}

type reobotResponse struct {
	Response string `json:"response"`
	Result   int    `json:"result"`
}

func getPostMessageBySessionId(sessionid string, msg string) (int, string) {
	msg = url.QueryEscape(msg)
	buffer, err := lib.GetResultByMethod("GET", gRobotUrl+msg, nil)
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

		for word, _ := range keywords {
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
	sentence := lib.SQLSentence(lib.SQLMAP_Select_HaveSameReply)
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
	code, message := getPostMessageBySessionId(sessionid, node.message)

	if 200 == code {
		if true != tulingResponseCheck(message, node.fromid) {
			message = tulingResponseChange(node.toid)
		}
	} else {
		message = tulingResponseChange(node.toid)
	}

	return message
}

/*
 |    Function: RemoveCommentToPush
 |      Author: Mr.Sancho
 |        Date: 2016-04-15
 |   Arguments:
 |      Return:
 | Description: 删除要推送的消息
 |
*/
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
		time.Sleep(lib.SLEEP_DURATION_PUSH_QUEUEMSG)
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
					sentence := lib.SQLSentence(lib.SQLMAP_Select_RecommendCount)
					lib.SQLQueryRow(sentence, 1, n.fromid, n.toid).Scan(&count)
					if 0 < count && true == lib.RandomHitPercent(25) {
						needreplay = false
					}
				}

				if needreplay {
					timevalue := lib.CurrentTimeUTCInt64()
					msg := getResponseMsg(n)
					if "" != msg {
						RecommendInsertMessageToDB(n.toid, n.fromid, n.msgtype, msg, timevalue)
						RecommendPushMessage(n.toid, n.fromid, n.tousertype, n.fromusertype, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
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
		time.Sleep(lib.SleepTimeDuration(lib.SLEEP_TYPE_ROBOTRECOMMEND))
		needpush = false

		gLiveUsersInfo.lock.RLock()
		for id, user := range gLiveUsersInfo.users {
			count = 0
			err := lib.SQLQueryRow(gRecommendInUnreadCount, id).Scan(&count)
			if nil != err || RECOMMEND_MAX_UNREADNUMBER <= count {
				continue
			}

			// 付费用户每次100%几率，未付费用户（有照片：70%，无照片：40%)
			if false == checkIfUserHaveViplevel(id, user.gender) {
				if true == checkIfUserHavePicture(id, user.gender) {
					if lib.RandomHitPercent(10) {
						continue
					}
				} else {
					if lib.RandomHitPercent(40) {
						continue
					}
				}
			}

			/* get the random id */
			var fromid int
			if 0 == user.gender {
				fromid = getRandomUserId(id, 1-user.gender)
			} else {
				fromid = getRandomHeartbeatId(id, 1-user.gender)
			}

			if 0 == fromid {
				continue
			}

			ok, msg := getFormatMsg(fromid)
			if !ok {
				msg = gHelloArray[lib.Intn(len(gHelloArray))]
			}

			timevalue := lib.CurrentTimeUTCInt64()
			RecommendInsertMessageToDB(fromid, id, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
			RecommendPushMessage(fromid, id, 0, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
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
	sentence := lib.SQLSentence(lib.SQLMAP_Select_CheckCommentDailyLock)
	err := lib.SQLQueryRow(sentence, fromid, toid, msgtype).Scan(&timevalue)
	if nil != err || 0 == timevalue {
		return false
	}

	timeUTC := lib.Int64_To_UTCTime(timevalue)
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

	if true == checkAlreadySendSameCommentToday(id, toid, RECOMMEND_MSGTYPE_GREET) {
		return 403, -1, "抱歉,一天只能向同一个人打招呼一次.", false
	}

	RecommendInsertMessageToDB(id, toid, RECOMMEND_MSGTYPE_GREET, "", timevalue)
	ok, replymsg := getFormatMsg(id)
	if !ok {
		replymsg = gHelloArray[lib.Intn(len(gHelloArray))]
	}

	lastid, err := RecommendInsertMessageToDB(id, toid, RECOMMEND_MSGTYPE_TALK, replymsg, timevalue)
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

	if true == checkAlreadySendSameCommentToday(id, toid, RECOMMEND_MSGTYPE_HEART) {
		return 403, -1, "抱歉,一天只能向同一个人发送一次心动消息.", false
	}

	lastid, err := RecommendInsertMessageToDB(id, toid, RECOMMEND_MSGTYPE_HEART, "", timevalue)
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

	if true == checkAlreadySendSameCommentToday(id, toid, RECOMMEND_MSGTYPE_ASKMSG) {
		return 403, -1, "抱歉,一天只能向同一个人邀请一次.", false
	}

	RecommendInsertMessageToDB(id, toid, RECOMMEND_MSGTYPE_ASKMSG, "", timevalue)
	lastid, err := RecommendInsertMessageToDB(id, toid, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
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
	if gender == togender {
		return 403, -1, "抱歉,同性之间不能发消息.", false
	}

	lastid, err := RecommendInsertMessageToDB(id, toid, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
	if nil != err {
		return 404, -1, err.Error(), false
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

	_, _, tousertype := GetGenderUsertypeById(toid)
	if common.USERTYPE_USER == tousertype {
		/*
		 * 注册用户直接发送
		 * 这里发送者一定是注册用户
		 */
		RecommendPushMessage(id, toid, common.USERTYPE_USER, common.USERTYPE_USER, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
		push.DoPush()
	} else {
		/*
		 * 自动回复缓存, 例外如下
		 * 1.索要照片等信息除外
		 * 2.如果发消息用户为VIP用户,也不会缓存
		 */
		if RECOMMEND_MSGTYPE_ASKMSG == msgtype {
			return
		}

		if true == checkIfUserHaveViplevel(id, gender) {
			return
		}

		gRecommendQueueLock.Lock()
		gRecommendQueue.PushBack(&recommendQueueNode{
			timewait:     int64(lib.SleepTimeDuration(lib.SLEEP_TYPE_ROBOTREPLY) / time.Second),
			fromid:       id,
			toid:         toid,
			fromusertype: common.USERTYPE_USER,
			tousertype:   common.USERTYPE_RB,
			msgtype:      RECOMMEND_MSGTYPE_TALK,
			message:      msg,
			timevalue:    timevalue})
		gRecommendQueueLock.Unlock()
	}
}

/*
 *
 *    Function: ActionRecommend
 *      Author: sunchao
 *        Date: 15/8/16
 * Description: 打招呼
 *
 */
func ActionRecommend(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, http.ErrNotSupported.Error()
	}

	v := req.URL.Query()
	toidStr := v.Get("toid")
	typeStr := v.Get("type")
	if "" == toidStr || "" == typeStr {
		return 404, ""
	}

	msgtype, _ := strconv.Atoi(typeStr)
	toid, _ := strconv.Atoi(toidStr)

	/* the dest id not exist */
	exist, togender := getGenderById(toid)
	if true != exist {
		return 404, http.ErrNotSupported.Error()
	}

	msg := v.Get("msg")
	if 0 == strings.Compare("我对你感兴趣，方便聊一下吗？", msg) {
		msgtype = RECOMMEND_MSGTYPE_GREET
	}

	t := time.Now()
	timevalue := lib.Time_To_UTCInt64(t)

	var code, lastid int
	var replymsg string
	var bpush bool

	switch msgtype {
	case RECOMMEND_MSGTYPE_GREET:
		code, lastid, replymsg, bpush = incomeGreetCommentProc(id, gender, toid, togender, msg, timevalue)
	case RECOMMEND_MSGTYPE_HEART:
		code, lastid, replymsg, bpush = incomeHeartCommentProc(id, gender, toid, togender, msg, timevalue)
	case RECOMMEND_MSGTYPE_ASKMSG:
		code, lastid, replymsg, bpush = incomeAskCommentProc(id, gender, toid, togender, msg, timevalue)
	case RECOMMEND_MSGTYPE_TALK:
		code, lastid, replymsg, bpush = incomeTalkCommentProc(id, gender, toid, togender, msg, timevalue)
	default:
		return 404, ""
	}

	//不继续处理
	if 200 != code {
		return code, replymsg
	}

	if true == bpush {
		go incomeCommentPushMsgProc(id, gender, toid, msgtype, replymsg, timevalue)
	}

	gCountApiRecommend = gCountApiRecommend + 1
	jsonRlt, _ := json.Marshal(messageInfo{
		MsgId:     lastid,
		MsgText:   replymsg,
		UserId:    toid,
		Direction: MESSAGE_DIRECTION_FROMME,
		Readed:    false,
		TimeUTC:   t,
	})

	return 200, string(jsonRlt)
}

/*
 *
 *    Function: DelRecommend
 *      Author: sunchao
 *        Date: 15/8/16
 * Description: delete the recommend
 *
 */
func DelRecommend(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, http.ErrNotSupported.Error()
	}

	v := req.URL.Query()
	msgidStr := v.Get("msgid")
	if "" == msgidStr {
		return 404, http.ErrNotSupported.Error()
	}

	talkidStr := v.Get("talkid")
	if "" == talkidStr {
		return 404, http.ErrNotSupported.Error()
	}

	/* the dest id not exist */
	msgid, _ := strconv.Atoi(msgidStr)
	talkid, _ := strconv.Atoi(talkidStr)
	sentence := lib.SQLSentence(lib.SQLMAP_Delete_Recommend)
	lib.SQLExec(sentence, msgid, talkid, id, id, talkid)
	gRecommendNumber = gRecommendNumber - 1

	return 200, ""
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

	infos := make([]messageInfo, 0)
	for rows.Next() {
		err = rows.Scan(&info.MsgId, &fromid, &toid, &readtmp, &timetmp, &info.MsgText)
		if nil == err {
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)

			if fromid == id || 1 == readtmp {
				info.Readed = true
			} else {
				info.Readed = false
			}

			if fromid == id {
				info.UserId = toid
				info.Direction = MESSAGE_DIRECTION_FROMME
			} else {
				info.UserId = fromid
				info.Direction = MESSAGE_DIRECTION_TOME
			}

			infos = append(infos, info)
		} else {
			log.Error(err.Error())
		}
	}

	return infos
}

/*
 *
 *    Function: GetWaterFlow
 *      Author: sunchao
 *        Date: 15/8/16
 * Description: 获取聊天流水
 *
 */
func GetWaterFlow(req *http.Request) (int, string) {
	v := req.URL.Query()
	talkidStr := v.Get("talkid")
	if "" == talkidStr {
		return 404, ""
	}

	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	var lastMsgId int
	lastMsgIdStr := v.Get("lastmsgid")
	if "" != lastMsgIdStr {
		lastMsgId, _ = strconv.Atoi(lastMsgIdStr)
		if 0 > lastMsgId {
			lastMsgId = 0
		}
	}

	pageid, count := lib.Get_pageid_count_fromreq(req)
	talkid, _ := strconv.Atoi(talkidStr)
	exist, _ = getGenderById(talkid)
	if true != exist {
		return 404, ""
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Select_MessageHistory)
	rows, err := lib.SQLQuery(sentence, RECOMMEND_MSGTYPE_TALK, lastMsgId, id, talkid, talkid, id, (pageid-1)*count, count)
	if nil != err {
		return 404, ""
	}

	defer rows.Close()

	infos := getRecommendByRows(id, rows)
	if len(infos) > 0 {
		sentence := lib.SQLSentence(lib.SQLMAP_Update_RecommendRead)
		lib.SQLExec(sentence, talkid, id, infos[0].MsgId)
	}
	jsonRlt, _ := json.Marshal(infos)

	return 200, string(jsonRlt)
}

/*
 *
 *    Function: GetRecommendAll
 *      Author: sunchao
 *        Date: 15/10/5
 * Description: 获取所有的聊天信息
 *
 */
func GetRecommendAll(timeline int64, id, pageid, count int) ([]messageInfo, error) {
	sentence := lib.SQLSentence(lib.SQLMAP_Select_DistinctRecommend)
	rows, err := lib.SQLQuery(sentence, id, RECOMMEND_MSGTYPE_TALK, timeline, id, RECOMMEND_MSGTYPE_TALK, timeline, (pageid-1)*count, count)
	if nil == err {
		defer rows.Close()
		var gender int

		infos := getRecommendByRows(id, rows)
		for i, _ := range infos {
			_, gender = getGenderById(infos[i].UserId)
			_, userinfo := GetUserInfo(infos[i].UserId, gender)
			infos[i].UserInfo = &userinfo
		}

		return infos, nil
	}

	return nil, err
}

/*
 *
 *    Function: GetAllMessage
 *      Author: sunchao
 *        Date: 15/10/5
 * Description: 获取所有的聊天信息和谁看过我信息
 *
 */
func GetAllMessage(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	var timeline int64

	v := req.URL.Query()

	timelinestr := v.Get("lasttime")
	if "" != timelinestr {
		timeline = lib.TimeStr_To_UTCInt64(timelinestr)
	}

	var allmessage allMessageInfo
	pageid, count := lib.Get_pageid_count_fromreq(req)
	recommendAlls, err := GetRecommendAll(timeline, id, pageid, count)
	if nil == err {
		allmessage.RecommendArray = recommendAlls
	}

	visitAlls, err := getVisitAll(timeline, id, pageid, count)
	if nil == err {
		allmessage.VisitArray = visitAlls
	}

	go log.Tracef("获取所有聊天信息: Id=%d gender=%d", id, gender)
	jsonRlt, _ := json.Marshal(allmessage)
	return 200, string(jsonRlt)
}

/*
 |    Function: GetUnreadMessage
 |      Author: Mr.Sancho
 |        Date: 2016-05-06
 | Description: get the unread message
 |      Return:
 |
*/
func GetUnreadMessage(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	var timeline int64

	v := req.URL.Query()
	timelinestr := v.Get("lasttime")
	if "" != timelinestr {
		timeline = lib.TimeStr_To_UTCInt64(timelinestr)
	}

	recommendCount := recommend_GetUnreadNum(id, timeline)
	visitCount := visit_GetUnreadNum(id, timeline)
	unreadmsg := unreadMessageInfo{UnreadRecommend: recommendCount, UnreadVisit: visitCount, Badge: recommendCount + visitCount}
	jsonRlt, _ := json.Marshal(unreadmsg)
	return 200, string(jsonRlt)
}

func recommend_GetUnreadNum(id int, timeline int64) int {
	var count int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_UnreadMessageCount)
	lib.SQLQueryRow(sentence, id, timeline).Scan(&count)

	return count
}

func GetApiRecommendCount() int {
	return gCountApiRecommend
}

/*
 |    Function: PeriodOnlineCommentSet
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 定期推送消息设置
 |      Return:
 |
*/
func PeriodOnlineCommentSet(enable bool, msg string) {
	gEnableEvaluation = enable
	gEvaluationMsgContent = msg
}

/*
 |    Function: PeriodOnlineCommentPush
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 定期推送消息处理
 |      Return:
 |
*/
func PeriodOnlineCommentPush(id, gender int, lastEvaluationTime int64) {
	curTime := lib.CurrentTimeUTCInt64()
	if 0 == lastEvaluationTime || 43200 <= (curTime-lastEvaluationTime) {
		evaluationMsg := PushMsgEvaluation{Enable: true, ShowMessage: gEvaluationMsgContent}
		jsonRlt, _ := json.Marshal(evaluationMsg)
		notifymsg := PushMessageInfo{Type: push.PUSH_NOTIFYMSG_EVALUATION, Value: string(jsonRlt)}
		jsonRlt, _ = json.Marshal(notifymsg)

		push.Add(0, GetClientIdByUserId(id), push.PUSHMSG_TYPE_NOTIFYMSG, push.PUSH_NOTIFYMSG_EVALUATION, "", string(jsonRlt))
		push.DoPush()

		sentence := lib.SQLSentence(lib.SQLMAP_Update_EvaluationTime, gender)
		lib.SQLExec(sentence, curTime, id)
	}

	return
}
