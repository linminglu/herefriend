package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
func Report(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	reportedidstr := v.Get("reportedid")
	reportedid, err := strconv.Atoi(reportedidstr)
	if nil != err || 0 == reportedid {
		return 404, err.Error()
	}

	reason := v.Get("reason")
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Report)
	_, err = lib.SQLExec(sentence, id, reportedid, reason)
	if nil != err {
		return 404, err.Error()
	}

	return 200, ""
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
func UserAddBlacklist(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	blackidstr := v.Get("blacklistid")
	blackid, err := strconv.Atoi(blackidstr)
	if nil != err || 0 == blackid {
		return 404, err.Error()
	}

	var idcheck int
	checksentence := lib.SQLSentence(lib.SQLMAP_Select_CheckUserBlacklist)
	lib.SQLQueryRow(checksentence, id, blackid).Scan(&idcheck)
	if idcheck == blackid {
		return 200, ""
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Insert_UserBlacklist)
	_, err = lib.SQLExec(sentence, id, blackid)
	if nil != err {
		return 404, err.Error()
	}

	return 200, ""
}

func UserDelBlacklist(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	blackidstr := v.Get("blacklistid")
	blackid, err := strconv.Atoi(blackidstr)
	if nil != err || 0 == blackid {
		return 404, err.Error()
	}

	var idcheck int
	checksentence := lib.SQLSentence(lib.SQLMAP_Select_CheckUserBlacklist)
	lib.SQLQueryRow(checksentence, id, blackid).Scan(&idcheck)
	if 0 == idcheck {
		return 200, ""
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Delete_UserBlacklist)
	_, err = lib.SQLExec(sentence, id, blackid)
	if nil != err {
		return 404, err.Error()
	}

	return 200, ""
}

func UserGetBlacklist(req *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
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

	jsonRlt, _ := json.Marshal(blist)
	return 200, string(jsonRlt)
}
