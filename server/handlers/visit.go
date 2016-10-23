package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"herefriend/lib"
	"herefriend/lib/push"
)

const (
	// VisitMaxUnreadNum .
	VisitMaxUnreadNum = 15
)

var gVisitUnreadSentence string
var gVisitInsertSentence string

func init() {
	gVisitUnreadSentence = lib.SQLSentence(lib.SQLMapSelectVisitUnreadCount)
	gVisitInsertSentence = lib.SQLSentence(lib.SQLMapInsertVisit)
	go visitRobotRoutine()
}

func getVisitAll(timeline int64, id, pageid, count int) ([]messageInfo, error) {
	var rows *sql.Rows
	var err error

	sentence := lib.SQLSentence(lib.SQLMapSelectVisitByRows)
	rows, err = lib.SQLQuery(sentence, id, timeline, (pageid-1)*count, count)
	if nil != err {
		return nil, err
	}

	defer rows.Close()

	var info messageInfo
	var readtmp int
	var timetmp int64

	var infos []messageInfo
	for rows.Next() {
		err = rows.Scan(&info.MsgID, &info.UserID, &readtmp, &timetmp)
		if nil == err {
			info.TimeUTC = lib.Int64ToUTCTime(timetmp)
			info.Direction = MessageDirectionToMe

			if 1 == readtmp {
				info.Readed = true
			} else {
				info.Readed = false
			}

			_, gender := getGenderByID(info.UserID)
			_, userinfo := GetUserInfo(info.UserID, gender)
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
	RecommendPushMessage(visitor, id, 0, 1, push.PushMsgVisit, "", t)
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
		time.Sleep(lib.SleepTimeDuration(lib.SleepTypeRobotVisit))

		needpush = false
		gLiveUsersInfo.lock.RLock()
		for id, user := range gLiveUsersInfo.users {
			count = 0
			err := lib.SQLQueryRow(gVisitUnreadSentence, id, 0).Scan(&count)
			if nil != err || VisitMaxUnreadNum <= count {
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
				fromid = getRandomUserID(id, 1)
			} else {
				fromid = getRandomHeartbeatID(id, 0)
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

// ReadVisit 将浏览信息设置为已读
func ReadVisit(c *gin.Context) {
	exist, _, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	visitidstr := c.Query("visitid")
	if visitidstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	visitid, _ := strconv.Atoi(visitidstr)
	sentence := lib.SQLSentence(lib.SQLMapUpdateVisitRead)
	_, err := lib.SQLExec(sentence, visitid)
	if err != nil {
		log.Error(err.Error())
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
	return
}

// DoVisit .
func DoVisit(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	toidstr := c.Query("toid")
	if toidstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	toid, _ := strconv.Atoi(toidstr)
	exist, togender := getGenderByID(toid)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	t := time.Now().UTC()
	sentence := lib.SQLSentence(lib.SQLMapInsertVisit)
	_, err := lib.SQLExec(sentence, id, toid, t.Unix())
	if err != nil {
		log.Error(err.Error())
		c.Status(http.StatusNotFound)
		return
	}

	_, info := GetUserInfo(toid, togender)
	c.JSON(http.StatusOK, info)
}

// DeleteVisit .
func DeleteVisit(c *gin.Context) {
	exist, _, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	visitidstr := c.Query("visitid")
	if visitidstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	visitid, _ := strconv.Atoi(visitidstr)
	sentence := lib.SQLSentence(lib.SQLMapDeleteVisit)
	_, err := lib.SQLExec(sentence, visitid)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}

func visitGetUnreadNum(id int, timeline int64) int {
	var count int

	sentence := lib.SQLSentence(lib.SQLMapSelectVisitUnreadCount)
	lib.SQLQueryRow(sentence, id, timeline).Scan(&count)

	return count
}
