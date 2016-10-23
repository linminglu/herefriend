package handlers

import (
	"encoding/json"
	"strconv"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"herefriend/common"
	"herefriend/lib"
)

func getHeartbeatBaseCountByProvinceGender(province string, gender int) int {
	count, exist := lib.GetRedisHeartbeatProvCount(province, gender)
	if true != exist {
		sentence := lib.SQLSentence(lib.SQLMapSelectHeartbeatProvinceCount, gender)
		lib.SQLQueryRow(sentence, province).Scan(&count)
		if 0 != count {
			lib.SetRedisHeartbeatProvCount(province, gender, count)
		}
	}

	return count
}

func doReqHeartbeat(id, gender, count int) (int, string) {
	var info common.PersonInfo

	_, info = GetUserInfo(id, gender)
	gender = 1 - gender

	/*
	 * get the persons' infos, with random page
	 */
	baseline := getHeartbeatBaseCountByProvinceGender(info.Province, gender)
	sentence := lib.SQLSentence(lib.SQLMapSelectHeartbeatInfoByRows, gender)
	rows, err := lib.SQLQuery(sentence, info.Province, lib.Intn(baseline-count), count)
	if nil != err {
		log.Error(err.Error())
		return 404, ""
	}
	defer rows.Close()

	var idtmp int
	var code int

	//init size with 0, if there is no data, http response body will be []
	var infos []common.PersonInfo
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

	jsonRlt, _ := json.Marshal(infos)
	return 200, string(jsonRlt)
}

// Heartbeat 心动女生
func Heartbeat(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	_, gender := getGenderByID(id)

	count := lib.GetCountRequestArgument(c)
	code, content := doReqHeartbeat(id, gender, count)
	c.String(code, content)
}
