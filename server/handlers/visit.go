package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	log "github.com/cihub/seelog"

	"herefriend/lib"
	"herefriend/lib/push"
)

const (
	VISIT_MAX_UNREADNUMBER = 15
)

var gVisitUnreadSentence string
var gVisitInsertSentence string

func init() {
	gVisitUnreadSentence = lib.SQLSentence(lib.SQLMAP_Select_VisitUnreadCount)
	gVisitInsertSentence = lib.SQLSentence(lib.SQLMAP_Insert_Visit)
	go visitRobotRoutine()
}

/*
 |
 |    Function: getVisitAll
 |      Author: sunchao
 |        Date: 15/10/6
 | Description: get visit
 |
*/
func getVisitAll(timeline int64, id, pageid, count int) ([]messageInfo, error) {
	var rows *sql.Rows
	var err error

	sentence := lib.SQLSentence(lib.SQLMAP_Select_VisitByRows)
	rows, err = lib.SQLQuery(sentence, id, timeline, (pageid-1)*count, count)
	if nil != err {
		return nil, err
	}

	defer rows.Close()

	var info messageInfo
	var readtmp int
	var timetmp int64

	infos := make([]messageInfo, 0)
	for rows.Next() {
		err = rows.Scan(&info.MsgId, &info.UserId, &readtmp, &timetmp)
		if nil == err {
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)
			info.Direction = MESSAGE_DIRECTION_TOME

			if 1 == readtmp {
				info.Readed = true
			} else {
				info.Readed = false
			}

			_, gender := getGenderById(info.UserId)
			_, userinfo := GetUserInfo(info.UserId, gender)
			info.UserInfo = &userinfo

			infos = append(infos, info)
		}
	}

	return infos, nil
}

func visitAddVisitor(id, visitor int, timemin, timemax int64) {
	var t int64

	if timemin == timemax {
		if 0 == timemin {
			t = lib.CurrentTimeUTCInt64()
		} else {
			t = timemin
		}
	} else {
		t = timemin + lib.Int63n(timemax-timemin)
	}

	lib.SQLExec(gVisitInsertSentence, visitor, id, t)
	RecommendPushMessage(visitor, id, 0, 1, push.PUSHMSG_TYPE_VISIT, "", t)
}

/*
 |
 |    Function: visitRobotRoutine
 |      Author: sunchao
 |        Date: 15/10/10
 | Description: 后台自动访问的线程
 |
*/
func visitRobotRoutine() {
	var count int
	var fromid int

	needpush := false
	for {
		time.Sleep(lib.SleepTimeDuration(lib.SLEEP_TYPE_ROBOTVISIT))

		needpush = false
		gLiveUsersInfo.lock.RLock()
		for id, user := range gLiveUsersInfo.users {
			count = 0
			err := lib.SQLQueryRow(gVisitUnreadSentence, id, 0).Scan(&count)
			if nil != err || VISIT_MAX_UNREADNUMBER <= count {
				if nil != err {
					lib.SQLError(gVisitUnreadSentence, err, id, 0)
				}
				continue
			}

			// 付费用户每次100%几率，未付费用户（有照片：70%，无照片：40%)
			if false == checkIfUserHaveViplevel(id, user.gender) {
				if true == checkIfUserHavePicture(id, user.gender) {
					if lib.RandomHitPercent(30) {
						continue
					}
				} else {
					if lib.RandomHitPercent(60) {
						continue
					}
				}
			}

			if 0 == user.gender {
				fromid = getRandomUserId(id, 1)
			} else {
				fromid = getRandomHeartbeatId(id, 0)
			}

			if 0 == fromid {
				continue
			}

			visitAddVisitor(id, fromid, 0, 0)
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
 *    Function: ReadVisit
 *      Author: sunchao
 *        Date: 15/10/6
 * Description: 将浏览信息设置为已读
 *
 */
func ReadVisit(req *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	visitidstr := v.Get("visitid")
	if "" == visitidstr {
		return 404, ""
	}

	visitid, _ := strconv.Atoi(visitidstr)
	sentence := lib.SQLSentence(lib.SQLMAP_Update_VisitRead)
	_, err := lib.SQLExec(sentence, visitid)
	if nil != err {
		log.Error(err.Error())
		return 404, ""
	}

	return 200, ""
}

/*
 *
 *    Function: DoVisit
 *      Author: sunchao
 *        Date: 15/10/6
 * Description: post visit
 *
 */
func DoVisit(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	toidstr := v.Get("toid")
	if "" == toidstr {
		return 404, ""
	}

	toid, _ := strconv.Atoi(toidstr)
	exist, togender := getGenderById(toid)
	if true != exist {
		return 404, ""
	}

	t := time.Now().UTC()
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Visit)
	_, err := lib.SQLExec(sentence, id, toid, t.Unix())
	if nil != err {
		log.Error(err.Error())
		return 404, ""
	}

	_, info := GetUserInfo(toid, togender)
	jsonRlt, _ := json.Marshal(info)

	return 200, string(jsonRlt)
}

/*
 *
 *    Function: DeleteVisit
 *      Author: sunchao
 *        Date: 15/10/6
 * Description: delete visit
 *
 */
func DeleteVisit(req *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	visitidstr := v.Get("visitid")
	if "" == visitidstr {
		return 404, ""
	}

	visitid, _ := strconv.Atoi(visitidstr)
	sentence := lib.SQLSentence(lib.SQLMAP_Delete_Visit)
	_, err := lib.SQLExec(sentence, visitid)
	if nil != err {
		return 404, ""
	}

	return 200, ""
}

func visit_GetUnreadNum(id int, timeline int64) int {
	var count int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_VisitUnreadCount)
	lib.SQLQueryRow(sentence, id, timeline).Scan(&count)

	return count
}
