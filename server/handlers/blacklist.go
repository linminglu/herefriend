package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"herefriend/lib"
)

/*
 |    Function: Report
 |      Author: Mr.Sancho
 |        Date: 2016-01-24
 |   Arguments:
 |      Return:
 | Description: report with reason
 |
*/
func Report(c *gin.Context) {
	exist, id, _ := getIdGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	reportedidstr := c.Query("reportedid")
	reportedid, err := strconv.Atoi(reportedidstr)
	if nil != err || 0 == reportedid {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	reason := c.Query("reason")
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Report)
	_, err = lib.SQLExec(sentence, id, reportedid, reason)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

/*
 |    Function: UserAddBlacklist
 |      Author: Mr.Sancho
 |        Date: 2016-02-29
 |   Arguments:
 |      Return:
 | Description: 用户黑名单
 |
*/
func UserAddBlacklist(c *gin.Context) {
	exist, id, _ := getIdGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	blackidstr := c.Query("blacklistid")
	blackid, err := strconv.Atoi(blackidstr)
	if nil != err || 0 == blackid {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	var idcheck int
	checksentence := lib.SQLSentence(lib.SQLMAP_Select_CheckUserBlacklist)
	lib.SQLQueryRow(checksentence, id, blackid).Scan(&idcheck)
	if idcheck == blackid {
		c.Status(http.StatusOK)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Insert_UserBlacklist)
	_, err = lib.SQLExec(sentence, id, blackid)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func UserDelBlacklist(c *gin.Context) {
	exist, id, _ := getIdGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	blackidstr := c.Query("blacklistid")
	blackid, err := strconv.Atoi(blackidstr)
	if nil != err || 0 == blackid {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	var idcheck int
	checksentence := lib.SQLSentence(lib.SQLMAP_Select_CheckUserBlacklist)
	lib.SQLQueryRow(checksentence, id, blackid).Scan(&idcheck)
	if 0 == idcheck {
		c.Status(http.StatusOK)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Delete_UserBlacklist)
	_, err = lib.SQLExec(sentence, id, blackid)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func UserGetBlacklist(c *gin.Context) {
	exist, id, _ := getIdGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var blist userBlacklist
	var blackid int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_UserBlacklist)
	rows, err := lib.SQLQuery(sentence, id)
	if nil == err {
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&blackid)
			if nil == err {
				blist.Blacklist = append(blist.Blacklist, blackid)
			}
		}
	}

	blist.Id = id
	c.JSON(http.StatusOK, blist)
}
