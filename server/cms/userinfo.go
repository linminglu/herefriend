package cms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/cihub/seelog"

	"herefriend/common"
	"herefriend/lib"
	"herefriend/server/handlers"
)

/*
 *
 *    Function: SetHeartbeat
 *      Author: sunchao
 *        Date: 15/7/12
 * Description: change the heartbeat status
 *
 */
func SetHeartbeat(req *http.Request) (int, string) {
	v := req.URL.Query()
	idStr := v.Get("id")
	acttionStr := v.Get("action")
	genderStr := v.Get("gender")

	if "" == idStr || "" == acttionStr || "" == genderStr {
		return 404, ""
	}

	id, _ := strconv.Atoi(idStr)
	gender, _ := strconv.Atoi(genderStr)

	var err error
	if "0" == acttionStr {
		sentence := lib.SQLSentence(lib.SQLMAP_Delete_Heartbeat)
		_, err = lib.SQLExec(sentence, id)
	} else {
		sentence := lib.SQLSentence(lib.SQLMAP_Insert_Heartbeat)
		_, userinfo := handlers.GetUserInfo(id, gender)
		_, err = lib.SQLExec(sentence, id, gender, userinfo.Province)
	}

	if nil != err {
		return 404, ""
	}

	return 200, ""
}

/*
 |    Function: GetUserInfos
 |      Author: Mr.Sancho
 |        Date: 2016-01-21
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func GetUserInfos(req *http.Request) (int, string) {
	v := req.URL.Query()
	genderStr := v.Get("gender")

	if "" == genderStr {
		return 404, ""
	}

	gender, _ := strconv.Atoi(genderStr)

	/*
	 * Second get the persons' infos
	 */
	page, count := lib.Get_pageid_count_fromreq(req)
	sentence := lib.SQLSentence(lib.SQLMAP_CMS_Select_BriefInfoByRows, gender)
	rows, err := lib.SQLQuery(sentence, (page-1)*count, count)
	if nil != err {
		return 404, ""
	}
	defer rows.Close()

	var infos []cmsUserInfo
	for rows.Next() {
		var info cmsUserInfo
		var idChk int

		rows.Scan(&info.Id)
		code, userinfo := handlers.GetUserInfo(info.Id, gender)
		if 200 == code {
			info.Name = userinfo.Name
			info.Age = userinfo.Age
			info.Img = userinfo.IconUrl
			info.Province = userinfo.Province

			/* check if is heartbeat selected */
			checkSql := lib.SQLSentence(lib.SQLMAP_CMS_Select_CheckHeatbeatValid)
			lib.SQLQueryRow(checkSql, info.Id).Scan(&idChk)
			if idChk == info.Id {
				info.Selected = true
			}

			infos = append(infos, info)
		}
	}

	jsonRlt, _ := json.Marshal(infos)
	return 200, string(jsonRlt)
}

/*
 |    Function: GetSingleUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-01-25
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func GetSingleUserInfo(req *http.Request) (int, string) {
	v := req.URL.Query()
	idStr := v.Get("id")
	genderStr := v.Get("gender")

	if "" == idStr || "" == genderStr {
		return 404, ""
	}

	id, _ := strconv.Atoi(idStr)
	gender, _ := strconv.Atoi(genderStr)

	/*
	 * Second get the persons' infos
	 */
	var info cmsUserInfo
	var idChk int

	_, userinfo := handlers.GetUserInfo(id, gender)
	info.Id = id
	info.Age = userinfo.Age
	info.Img = userinfo.IconUrl
	info.Name = userinfo.Name
	info.Province = userinfo.Province

	checkSql := lib.SQLSentence(lib.SQLMAP_CMS_Select_CheckHeatbeatValid)
	lib.SQLQueryRow(checkSql, info.Id).Scan(&idChk)
	if idChk == info.Id {
		info.Selected = true
	}

	jsonRlt, _ := json.Marshal(info)
	return 200, string(jsonRlt)
}

/*
 |    Function: SetSingleUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-03-05
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func SetSingleUserInfo(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	idStr := v.Get("id")
	genderStr := v.Get("gender")

	if "" == idStr || "" == genderStr {
		w.WriteHeader(404)
		return
	}

	id, _ := strconv.Atoi(idStr)
	gender, _ := strconv.Atoi(genderStr)

	w.WriteHeader(handlers.UpdateProfile(req, id, gender))
	return
}

/*
 |    Function: ChangeHeadImage
 |      Author: Mr.Sancho
 |        Date: 2016-01-25
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func ChangeHeadImage(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	idStr := v.Get("id")
	genderStr := v.Get("gender")

	if "" == idStr || "" == genderStr {
		w.WriteHeader(404)
		return
	}

	id, _ := strconv.Atoi(idStr)
	_, _, usertype := handlers.GetGenderUsertypeById(id)
	if common.USERTYPE_USER == usertype {
		w.WriteHeader(403)
		return
	}

	gender, _ := strconv.Atoi(genderStr)

	sentence := lib.SQLSentence(lib.SQLMAP_CMS_Select_Pictures, gender)
	rows, err := lib.SQLQuery(sentence, id)
	if nil != err {
		w.WriteHeader(404)
		return
	}
	defer rows.Close()

	updateSentence := lib.SQLSentence(lib.SQLMAP_Update_SetPictureTag, gender)
	lib.DelRedisUserInfo(id)

	var infos []cmsImageInfo
	var info cmsImageInfo
	headindex := -1

	for rows.Next() {
		rows.Scan(&info.filename, &info.tag)
		if 1 == info.tag {
			headindex = len(infos)
		}

		infos = append(infos, info)
	}

	if -1 == headindex {
		if 0 != len(infos) {
			lib.SQLExec(updateSentence, 1, id, infos[0].filename)
		}
	} else {
		if 0 != len(infos) {
			lib.SQLExec(updateSentence, 0, id, infos[headindex].filename)
			lib.SQLExec(updateSentence, 1, id, infos[(headindex+1)%len(infos)].filename)
		}
	}

	return
}

/*
 |    Function: DeleteHeadImage
 |      Author: Mr.Sancho
 |        Date: 2016-02-13
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func DeleteHeadImage(w http.ResponseWriter, req *http.Request) {
	v := req.URL.Query()
	idStr := v.Get("id")
	genderStr := v.Get("gender")

	if "" == idStr || "" == genderStr {
		w.WriteHeader(404)
		return
	}

	id, _ := strconv.Atoi(idStr)
	if 1 == id {
		w.WriteHeader(403)
		return
	}

	_, _, usertype := handlers.GetGenderUsertypeById(id)
	if common.USERTYPE_USER == usertype {
		w.WriteHeader(403)
		return
	}

	gender, _ := strconv.Atoi(genderStr)
	sentence := lib.SQLSentence(lib.SQLMAP_CMS_Select_Pictures, gender)
	rows, err := lib.SQLQuery(sentence, id)
	if nil != err {
		w.WriteHeader(404)
		return
	}
	defer rows.Close()

	updateSentence := lib.SQLSentence(lib.SQLMAP_Update_SetPictureTag, gender)
	lib.DelRedisUserInfo(id)

	var infos []cmsImageInfo
	var info cmsImageInfo
	headindex := -1

	for rows.Next() {
		rows.Scan(&info.filename, &info.tag)
		if 1 == info.tag {
			headindex = len(infos)
		}

		infos = append(infos, info)
	}

	if -1 == headindex {
		if 0 != len(infos) {
			lib.SQLExec(updateSentence, 1, id, infos[0].filename)
		}
	} else {
		if 0 != len(infos) {
			deletesentence := lib.SQLSentence(lib.SQLMAP_Delete_HeadPicture, gender)
			lib.SQLExec(deletesentence, id)
			lib.DeleteImageFromQiniu(id, infos[headindex].filename)

			if 1 != len(infos) {
				lib.SQLExec(updateSentence, 1, id, infos[(headindex+1)%len(infos)].filename)
			}
		}
	}

	return
}

/*
 |    Function: AddBlacklist
 |      Author: Mr.Sancho
 |        Date: 2016-01-24
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func AddBlacklist(req *http.Request) (int, string) {
	v := req.URL.Query()
	idStr := v.Get("id")
	genderStr := v.Get("gender")

	if "" == idStr || "" == genderStr {
		return 404, ""
	}

	id, _ := strconv.Atoi(idStr)
	if 1 == id {
		return 403, ""
	}

	_, _, usertype := handlers.GetGenderUsertypeById(id)
	if common.USERTYPE_USER == usertype {
		return 403, ""
	}

	gender, _ := strconv.Atoi(genderStr)
	/* delete from live user queue */
	handlers.DeleteLiveUser(id)
	handlers.OfflineProc(id, gender)

	/* move to blacklist */
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Blacklist, gender)
	lib.SQLExec(sentence, id)

	delSql := lib.SQLSentence(lib.SQLMAP_Delete_UserId, gender)
	lib.SQLExec(delSql, id)
	lib.DelRedisUserInfo(id)

	var idChk int
	checkSql := lib.SQLSentence(lib.SQLMAP_CMS_Select_CheckHeatbeatValid)
	lib.SQLQueryRow(checkSql, id).Scan(&idChk)
	if idChk == id {
		delSql := lib.SQLSentence(lib.SQLMAP_Delete_Heartbeat)
		lib.SQLExec(delSql, id)
	}

	handlers.SubUserCount(gender)
	return 200, ""
}

/*
 |    Function: SearchUserInfos
 |      Author: Mr.Sancho
 |        Date: 2016-03-04
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func SearchUserInfos(req *http.Request) (int, string) {
	v := req.URL.Query()
	genderStr := v.Get("gender")
	fieldStr := v.Get("field")
	keyStr := v.Get("key")

	if "" == genderStr || "" == fieldStr || "" == keyStr {
		return 404, ""
	}

	gender, _ := strconv.Atoi(genderStr)
	fieldid, _ := strconv.Atoi(fieldStr)

	/*
	 * check the field
	 */
	if 2 < fieldid {
		return 404, ""
	}

	field := []string{"name", "introduction", "id"}[fieldid]

	countsentence := lib.SQLSentence(lib.SQLMAP_Select_UserCount, gender)
	countsentence += fmt.Sprintf(" where position('%s' in %s)", keyStr, field)

	var searchInfo cmsSearchInfo
	err := lib.SQLQueryRow(countsentence).Scan(&searchInfo.Count)
	if nil == err && 0 != searchInfo.Count {
		sentence := lib.SQLSentence(lib.SQLMAP_CMS_Select_BriefInfo, gender)
		sentence += fmt.Sprintf(" where position('%s' in %s) order by id desc limit ?,?", keyStr, field)

		page, count := lib.Get_pageid_count_fromreq(req)
		rows, err := lib.SQLQuery(sentence, (page-1)*count, count)
		if nil != err {
			return 404, ""
		}
		defer rows.Close()

		var info cmsUserInfo
		var idChk int

		for rows.Next() {
			rows.Scan(&info.Id)
			code, userinfo := handlers.GetUserInfo(info.Id, gender)
			if 200 == code {
				info.Name = userinfo.Name
				info.Age = userinfo.Age
				info.Img = userinfo.IconUrl
				info.Province = userinfo.Province

				/* check if is heartbeat selected */
				checkSql := lib.SQLSentence(lib.SQLMAP_CMS_Select_CheckHeatbeatValid)
				lib.SQLQueryRow(checkSql, info.Id).Scan(&idChk)
				if idChk == info.Id {
					info.Selected = true
				}

				searchInfo.Users = append(searchInfo.Users, info)
			}
		}
	} else {
		log.Errorf("SQLQueryRow Error: %s %v\n", countsentence, err)
	}

	jsonRlt, _ := json.Marshal(searchInfo)
	return 200, string(jsonRlt)
}

/*
 |    Function: SystemUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-01-30
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func SystemUserInfo(r *http.Request) string {
	info := cmsSystemUsersSummary{
		GirlsNum:  handlers.GetUserCountByGender(0),
		GuysNum:   handlers.GetUserCountByGender(1),
		ActiveNum: handlers.GetActiveUserNumber(),
		OnlineNum: handlers.GetLiveUserNumber(),
		RegistNum: handlers.GetRegistUserNumber(),
	}

	jsonRlt, _ := json.Marshal(info)
	return string(jsonRlt)
}

/*
 |    Function: RefreshUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-03-16
 |   Arguments:
 |      Return:
 | Description: 刷新用户信息
 |
*/
func RefreshUserInfo(w http.ResponseWriter, req *http.Request) {
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

	lib.DelRedisUserInfo(id)
	return
}

/*
 |    Function: RegistUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-03-31
 |   Arguments:
 |      Return:
 | Description: 获取注册用户信息
 |
*/
func RegistUserInfo(req *http.Request) (int, string) {
	v := req.URL.Query()
	genderStr := v.Get("gender")

	if "" == genderStr {
		return 404, ""
	}

	gender, _ := strconv.Atoi(genderStr)
	countsentence := lib.SQLSentence(lib.SQLMAP_Select_UserCount, gender) + " where usertype=1"

	var searchInfo cmsSearchInfo
	err := lib.SQLQueryRow(countsentence).Scan(&searchInfo.Count)
	if nil == err && 0 != searchInfo.Count {
		sentence := lib.SQLSentence(lib.SQLMAP_CMS_Select_BriefInfo, gender) + " where usertype=1 order by id desc limit ?,?"

		page, count := lib.Get_pageid_count_fromreq(req)
		rows, err := lib.SQLQuery(sentence, (page-1)*count, count)
		if nil != err {
			return 404, ""
		}
		defer rows.Close()

		var info cmsUserInfo

		for rows.Next() {
			rows.Scan(&info.Id)
			code, userinfo := handlers.GetUserInfo(info.Id, gender)
			if 200 == code {
				info.Name = userinfo.Name
				info.Age = userinfo.Age
				info.Img = userinfo.IconUrl
				info.Province = userinfo.Province
				info.VipLevel = userinfo.VipLevel

				searchInfo.Users = append(searchInfo.Users, info)
			}
		}
	} else {
		log.Errorf("SQLQueryRow Error: %s %v\n", countsentence, err)
	}

	jsonRlt, _ := json.Marshal(searchInfo)
	return 200, string(jsonRlt)
}
