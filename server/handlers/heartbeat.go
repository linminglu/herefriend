package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/cihub/seelog"

	"herefriend/common"
	"herefriend/lib"
)

/*
 |    Function: getHeartbeatBaseCountByProvinceGender
 |      Author: Mr.Sancho
 |        Date: 2016-02-21
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func getHeartbeatBaseCountByProvinceGender(province string, gender int) int {
	count, exist := lib.GetRedisHeartbeatProvCount(province, gender)
	if true != exist {
		sentence := lib.SQLSentence(lib.SQLMAP_Select_HeartbeatProvinceCount, gender)
		lib.SQLQueryRow(sentence, province).Scan(&count)
		if 0 != count {
			lib.SetRedisHeartbeatProvCount(province, gender, count)
		}
	}

	return count
}

/*
 |    Function: doReqHeartbeat
 |      Author: Mr.Sancho
 |        Date: 2016-02-21
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func doReqHeartbeat(id, gender, count int) (int, string) {
	var info common.PersonInfo

	_, info = GetUserInfo(id, gender)
	gender = 1 - gender

	/*
	 * get the persons' infos, with random page
	 */
	baseline := getHeartbeatBaseCountByProvinceGender(info.Province, gender)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_HeartbeatInfoByRows, gender)
	rows, err := lib.SQLQuery(sentence, info.Province, lib.Intn(baseline-count), count)
	if nil != err {
		log.Error(err.Error())
		return 404, ""
	}
	defer rows.Close()

	var idtmp int
	var code int

	//init size with 0, if there is no data, http response body will be []
	infos := make([]common.PersonInfo, 0)
	for rows.Next() {
		err = rows.Scan(&idtmp)
		if nil != err {
			continue
		}

		code, info = GetUserInfo(idtmp, gender)
		if 200 == code {
			infos = append(infos, info)
		}
	}

	go log.Tracef("获取推荐列表: Id=%d gender=%d", id, gender)
	jsonRlt, _ := json.Marshal(infos)
	return 200, string(jsonRlt)
}

/*
 *
 *    Function: Heartbeat
 *      Author: sunchao
 *        Date: 15/8/15
 * Description: 心动女生
 *
 */
func Heartbeat(req *http.Request) (int, string) {
	v := req.URL.Query()
	idStr := v.Get("id")

	id, _ := strconv.Atoi(idStr)
	_, gender := getGenderById(id)

	count := lib.GetCountRequestArgument(req)
	return doReqHeartbeat(id, gender, count)
}
