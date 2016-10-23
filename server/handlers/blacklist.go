package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"herefriend/lib"
)

// Report .
func Report(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
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
	sentence := lib.SQLSentence(lib.SQLMapInsertReport)
	_, err = lib.SQLExec(sentence, id, reportedid, reason)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// UserAddBlacklist 用户黑名单
func UserAddBlacklist(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
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
	checksentence := lib.SQLSentence(lib.SQLMapSelectCheckUserBlacklist)
	lib.SQLQueryRow(checksentence, id, blackid).Scan(&idcheck)
	if idcheck == blackid {
		c.Status(http.StatusOK)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMapInsertUserBlacklist)
	_, err = lib.SQLExec(sentence, id, blackid)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// UserDelBlacklist .
func UserDelBlacklist(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
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
	checksentence := lib.SQLSentence(lib.SQLMapSelectCheckUserBlacklist)
	lib.SQLQueryRow(checksentence, id, blackid).Scan(&idcheck)
	if 0 == idcheck {
		c.Status(http.StatusOK)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMapDeleteUserBlacklist)
	_, err = lib.SQLExec(sentence, id, blackid)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// UserGetBlacklist .
func UserGetBlacklist(c *gin.Context) {
	exist, id, _ := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var blist userBlacklist
	var blackid int

	sentence := lib.SQLSentence(lib.SQLMapSelectUserBlacklist)
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

	blist.ID = id
	c.JSON(http.StatusOK, blist)
}
