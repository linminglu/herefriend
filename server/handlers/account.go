package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"herefriend/common"
	"herefriend/config"
	"herefriend/lib"
	"herefriend/lib/push"
)

const (
	LIVEUSER_TICK_ONLINE     = 15   // minutes
	LIVEUSER_TICK_BACKGROUND = 2880 // minutes of 2 days
	VIPUSER_SLEEPDURATION    = int64(time.Hour)
)

const (
	LIVEUSER_STATUS_ONLINE     = 1
	LIVEUSER_STATUS_BACKGROUND = 2
)

var g_regLock sync.Mutex
var g_liveUsersInfo *liveUsersInfo
var g_vipUsersInfo *vipUsersInfo
var g_AddAgeNum = []int{7, 3}
var g_SubAgeNum = []int{3, 7}

var g_Count_Girls int
var g_Count_Guys int
var g_Count_Regist int
var g_Count_BuyVIP int

func init() {
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMAP_Select_UserCount, 0)).Scan(&g_Count_Girls)
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMAP_Select_UserCount, 1)).Scan(&g_Count_Guys)
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMAP_Select_UserActive)).Scan(&g_Count_Regist, &g_Count_BuyVIP)

	g_liveUsersInfo = &liveUsersInfo{users: make(map[int]*liveUser), lock: sync.RWMutex{}}
	go liveUserGoRoute()
	go liveUserNotifyMsgRoutine()

	g_vipUsersInfo = &vipUsersInfo{users: make(map[int]*vipUser), lock: sync.RWMutex{}}
	go vipUserGoRoute()
}

func updateUserActive() {
	lib.SQLExec(lib.SQLSentence(lib.SQLMAP_Update_UserActive), g_Count_Regist, g_Count_BuyVIP)
}

func GetUserCountByGender(gender int) int {
	if 0 == gender {
		return g_Count_Girls
	} else {
		return g_Count_Guys
	}
}

func updateLiveUserInfo(usersinfo *liveUsersInfo, id, gender, status, tick int) bool {
	usersinfo.lock.Lock()
	info, ok := usersinfo.users[id]
	if false == ok {
		info = &liveUser{gender: gender}
		usersinfo.users[id] = info
	}

	info.status = status
	info.livetick = tick

	usersinfo.lock.Unlock()
	return ok
}

func checkLiveUserBackground(usersinfo *liveUsersInfo, id int) bool {
	bBackground := false

	usersinfo.lock.Lock()
	info, ok := usersinfo.users[id]
	if true == ok {
		if LIVEUSER_STATUS_BACKGROUND == info.status {
			bBackground = true
		}
	}
	usersinfo.lock.Unlock()

	return bBackground
}

func deleteLiveUserInfo(usersinfo *liveUsersInfo, id int) {
	usersinfo.lock.Lock()
	delete(usersinfo.users, id)
	usersinfo.lock.Unlock()
}

func DeleteLiveUser(id int) {
	g_Count_Regist = g_Count_Regist - 1
	updateUserActive()

	deleteLiveUserInfo(g_liveUsersInfo, id)
}

func GetLiveUserNumber() int {
	return len(g_liveUsersInfo.users)
}

func GetRegistUserNumber() int {
	return g_Count_Regist
}

func updateVipUserInfo(usersinfo *vipUsersInfo, id, gender, level, days int, expiretime int64) bool {
	usersinfo.lock.Lock()
	info, ok := usersinfo.users[id]
	if true == ok {
		info.level = level
		info.days = days
		info.expiretime = expiretime
	} else {
		usersinfo.users[id] = &vipUser{gender: gender, level: level, days: days, expiretime: expiretime}
	}

	usersinfo.lock.Unlock()

	return ok
}

func deleteVipUserInfo(usersinfo *vipUsersInfo, id int) {
	usersinfo.lock.Lock()
	delete(usersinfo.users, id)
	usersinfo.lock.Unlock()
}

func onlineProc(id, gender int) {
	sentence := lib.SQLSentence(lib.SQLMAP_Update_Online, gender)
	lib.SQLExec(sentence, lib.CurrentTimeUTCInt64(), id)
}

func backgroundProc(id, gender int) {
	sentence := lib.SQLSentence(lib.SQLMAP_Update_Background, gender)
	lib.SQLExec(sentence, id)
}

func OfflineProc(id, gender int) {
	sentence := lib.SQLSentence(lib.SQLMAP_Update_Offline, gender)
	lib.SQLExec(sentence, id)
}

/*
 |    Function: getSearchAgeRange
 |      Author: Mr.Sancho
 |        Date: 2016-01-31
 |   Arguments:
 |      Return:
 | Description: 默认年龄段,女生为上7下3。男生为18~28
*/
func getSearchAgeRange(gender int, age int) (int, int) {
	min, max := 0, 0
	if config.Conf_AgeMin <= age && age <= config.Conf_AgeMax {
		if 0 == gender {
			min = age - g_SubAgeNum[gender]
			max = age + g_AddAgeNum[gender]

			if min < config.Conf_AgeMin {
				min = config.Conf_AgeMin
			}

			if max > config.Conf_AgeMax {
				max = config.Conf_AgeMax
			}
		} else {
			min = 18
			max = 28
		}
	}

	return min, max
}

/*
 |    Function: getRandomHeartbeatId
 |      Author: Mr.Sancho
 |        Date: 2016-02-21
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func getRandomHeartbeatId(id, gender int) int {
	var tmpid int
	var baselimit int
	var sentence string

	province, exist := lib.GetRedisDistrict(id)
	if true == exist {
		sentence = lib.SQLSentence(lib.SQLMAP_Select_HeartbeatRandomProvId, gender)

		baselimit = getHeartbeatBaseCountByProvinceGender(province, gender)
		lib.SQLQueryRow(sentence, province, lib.Intn(baselimit)).Scan(&tmpid)
		return tmpid
	}

	return getRandomUserId(id, gender)
}

/*
 |    Function: getRandomUserId
 |      Author: Mr.Sancho
 |        Date: 2016-01-01
 |   Arguments:
 |      Return:
 | Description: 根据用户的id与性别,获取一个合适的随机id
 |
*/
func getRandomUserId(id, gender int) int {
	var tmpid int
	var baselimit int

	province, exist := lib.GetRedisDistrict(id)
	if true == exist {
		/* the first step: get random id by province and age range
		 */
		getcount := false
		rangecount := 0
		_, info := GetUserInfo(id, gender)
		if config.Conf_AgeMin <= info.Age && info.Age <= config.Conf_AgeMax {
			min, max := getSearchAgeRange(info.Gender, info.Age)
			if 0 != min {
				for tmpage := min; tmpage <= max; tmpage = tmpage + 1 {
					count, exist := lib.GetRedisProvAgeCount(province, gender, tmpage)
					if true == exist {
						rangecount = rangecount + count
					} else {
						sentence := lib.SQLSentence(lib.SQLMAP_Select_CountByProvAge, gender)
						err := lib.SQLQueryRow(sentence, province, tmpage).Scan(&count)
						if nil == err && 0 != count {
							lib.SetRedisProvAgeCount(province, gender, tmpage, count)
							rangecount = rangecount + count
						}
					}
				}

				if 0 != rangecount {
					sentence := lib.SQLSentence(lib.SQLMAP_Select_RandomProvAgeId, gender)
					baselimit = lib.Intn(rangecount)
					err := lib.SQLQueryRow(sentence, province, min, max, baselimit).Scan(&tmpid)
					if nil == err && 1 < tmpid {
						return tmpid
					}
				}
			}
		}

		/* the second step: get random id by province
		 */
		getcount = false
		count, exist := lib.GetRedisProvCount(province, gender)
		if false == exist {
			sentence := lib.SQLSentence(lib.SQLMAP_Select_CountByProv, gender)
			err := lib.SQLQueryRow(sentence, province).Scan(&count)
			if nil == err && 0 != count {
				lib.SetRedisProvCount(province, gender, count)
				getcount = true
			}
		} else {
			getcount = true
		}

		if true == getcount {
			sentence := lib.SQLSentence(lib.SQLMAP_Select_RandomProvId, gender)
			baselimit = lib.Intn(count)
			err := lib.SQLQueryRow(sentence, province, baselimit).Scan(&tmpid)
			if nil == err && 1 < tmpid {
				return tmpid
			}
		}
	}

	/* the third step: get random id from all data
	 */
	if 0 == gender {
		baselimit = lib.Intn(g_Count_Girls)
	} else {
		baselimit = lib.Intn(g_Count_Guys)
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Select_RandomId, gender)
	err := lib.SQLQueryRow(sentence, baselimit).Scan(&tmpid)
	if nil != err || 1 >= tmpid {
		lib.SQLQueryRow(sentence, 2).Scan(&tmpid)
	}

	return tmpid
}

/*
 *
 *    Function: GetGenderUsertypeById
 *      Author: sunchao
 *        Date: 15/6/22
 * Description: get gender usertype by id, if return false means there is no such id in tables
 *
 */
func GetGenderUsertypeById(id int) (bool, int, int) {
	if 0 != id {
		var idtmp int
		var usertype int

		sentence := lib.SQLSentence(lib.SQLMAP_Select_UserType, 0)
		err := lib.SQLQueryRow(sentence, id).Scan(&idtmp, &usertype)
		if nil == err && id == idtmp {
			return true, 0, usertype
		}

		sentence = lib.SQLSentence(lib.SQLMAP_Select_UserType, 1)
		err = lib.SQLQueryRow(sentence, id).Scan(&idtmp, &usertype)
		if nil == err {
			return true, 1, usertype
		}
	}

	return false, 1, 0
}

/*
 |    Function: getUsertypeByIdGender
 |      Author: Mr.Sancho
 |        Date: 2016-01-09
 |   Arguments:
 |      Return:
 | Description: 根据id和性别获取用户类型
 |
*/
func getUsertypeByIdGender(id, gender int) (bool, int) {
	var idtmp int
	var usertype int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_UserType, gender)
	err := lib.SQLQueryRow(sentence, id).Scan(&idtmp, &usertype)
	if nil == err && id == idtmp {
		return true, usertype
	}

	return false, 0
}

/*
 *
 *    Function: getGenderById
 *      Author: sunchao
 *        Date: 15/6/22
 * Description: get gender by id, if return false means there is no such id in tables
 *
 */
func getGenderById(id int) (bool, int) {
	gender, exist := lib.GetRedisUserGender(id)
	if true == exist {
		return true, gender
	}

	var idstr int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_CheckIdValid, 0)
	err := lib.SQLQueryRow(sentence, id).Scan(&idstr)
	if nil == err {
		lib.SetRedisUserGender(id, 0)
		return true, 0
	}

	sentence = lib.SQLSentence(lib.SQLMAP_Select_CheckIdValid, 1)
	err = lib.SQLQueryRow(sentence, id).Scan(&idstr)
	if nil == err {
		lib.SetRedisUserGender(id, 1)
		return true, 1
	}

	return false, 1
}

/*
 *
 *    Function: GetGenderByIdPw
 *      Author: sunchao
 *        Date: 15/6/23
 * Description: get gender by id and password, if return false means there is no such id in tables
 *
 */
func getGenderByIdPw(id int, pw string) (bool, int) {
	if 0 != id {
		var idScan int

		sentence := lib.SQLSentence(lib.SQLMAP_Select_CheckIdPasswordValid, 0)
		err := lib.SQLQueryRow(sentence, id, pw).Scan(&idScan)
		if nil == err {
			return true, 0
		}

		sentence = lib.SQLSentence(lib.SQLMAP_Select_CheckIdPasswordValid, 1)
		err = lib.SQLQueryRow(sentence, id, pw).Scan(&idScan)
		if nil == err {
			return true, 1
		}
	}

	return false, 1
}

/*
 *
 *    Function: GetIdGenderByRequest
 *      Author: sunchao
 *        Date: 15/7/10
 * Description: 根据用户请求获取id和性别
 *
 */
func getIdGenderByRequest(req *http.Request) (bool, int, int) {
	v := req.URL.Query()
	idStr := v.Get("id")
	pwStr := v.Get("password")
	if "" == idStr || "" == pwStr {
		return false, 0, 0
	}

	id, _ := strconv.Atoi(idStr)
	bExist, gender := getGenderByIdPw(id, pwStr)
	if true != bExist {
		return false, 0, 0
	}

	return true, id, gender
}

func checkIfUserHavePicture(id, gender int) bool {
	var filename string

	sentence := lib.SQLSentence(lib.SQLMAP_Select_SearchPictures, gender)

	/* 获取头像 */
	err := lib.SQLQueryRow(sentence, id, 1).Scan(&filename)
	if nil == err && "" != filename {
		return true
	}

	/* 获取相册图片 */
	rows, err := lib.SQLQuery(sentence, id, 0)
	if nil == err {
		defer rows.Close()
		if rows.Next() {
			return true
		}
	}

	return false
}

func checkIfUserHaveViplevel(id, gender int) bool {
	var viplevel int
	var vipdays int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_VipLevelByID, gender)
	err := lib.SQLQueryRow(sentence, id).Scan(&viplevel, &vipdays)
	if nil == err && 0 != viplevel {
		return true
	}

	return false
}

/*
 *
 *    Function: GetUserPictrues
 *      Author: sunchao
 *        Date: 15/6/20
 * Description: get pictures by the id number
 *
 */
func getUserPictrues(id, gender int, info *personInfo) {
	var filename string

	sentence := lib.SQLSentence(lib.SQLMAP_Select_SearchPictures, gender)

	/* 获取头像 */
	err := lib.SQLQueryRow(sentence, id, 1).Scan(&filename)
	if nil == err && "" != filename {
		info.IconUrl = lib.GetQiniuUserImageURL(id, filename)
	}

	info.Pics = make([]string, 0)
	/* 获取相册图片 */
	rows, err := lib.SQLQuery(sentence, id, 0)
	if nil != err {
		return
	}

	defer rows.Close()
	for rows.Next() {
		filename = ""
		err = rows.Scan(&filename)
		if "" != filename {
			filename = lib.GetQiniuUserImageURL(id, filename)
			info.Pics = append(info.Pics, filename)
		}
	}

	return
}

/*
 |    Function: GetUserInfoById
 |      Author: Mr.Sancho
 |        Date: 2016-02-23
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func GetUserInfoById(id int) (int, personInfo) {
	var info personInfo

	exist, gender := getGenderById(id)
	if false == exist {
		return 404, info
	}

	code, info := GetUserInfo(id, gender)
	if 200 != code {
		return 404, info
	}

	return 200, info
}

/*
 *
 *    Function: GetUserInfo
 *      Author: sunchao
 *        Date: 15/7/11
 * Description: get the user information by id and gender
 *
 */
func GetUserInfo(id int, gender int) (int, personInfo) {
	var info personInfo

	content, exist := lib.GetRedisUserInfo(id)
	if true == exist {
		err := json.Unmarshal(content, &info)
		if nil == err {
			return 200, info
		}
	}

	var timeValue int64
	sentence := lib.SQLSentence(lib.SQLMAP_Select_PersonInfoById, gender)
	err := lib.SQLQueryRow(sentence, id).Scan(&info.Id, &info.Name, &info.Age, &info.Gender, &info.OnlineStatus, &info.VipLevel, &timeValue, &info.Height, &info.Weight,
		&info.LoveType, &info.BodyType, &info.Marriage, &info.Province, &info.District, &info.Native, &info.Education, &info.Income, &info.IncomeMin, &info.IncomeMax,
		&info.Occupation, &info.Housing, &info.Carstatus, &info.Introduction, &info.School, &info.Speciality,
		&info.Animal, &info.Constellation, &info.Lang, &info.BloodType, &info.Selfjudge, &info.Companytype, &info.Companyindustry,
		&info.Nationnality, &info.Religion, &info.Charactor, &info.Hobbies, &info.CityLove, &info.Naken, &info.Allow_age, &info.Allow_residence,
		&info.Allow_height, &info.Allow_marriage, &info.Allow_education, &info.Allow_housing, &info.Allow_income, &info.Allow_kidstatus)
	if nil == err {
		if 0 != info.VipLevel {
			info.VipExpireTime = lib.Int64_To_UTCTime(timeValue)
		}

		getUserPictrues(id, gender, &info)
	} else {
		return 404, info
	}

	jsonRlt, _ := json.Marshal(info)
	lib.SetRedisUserInfo(id, jsonRlt)

	return 200, info
}

func GetClientIdByUserId(id int) string {
	var clientid string

	sentence := lib.SQLSentence(lib.SQLMAP_Select_ClientID, 0)
	lib.SQLQueryRow(sentence, id).Scan(&clientid)
	if "" == clientid {
		sentence = lib.SQLSentence(lib.SQLMAP_Select_ClientID, 1)
		lib.SQLQueryRow(sentence, id).Scan(&clientid)
	}

	return clientid
}

func UpdateProfile(req *http.Request, id, gender int) int {
	v := req.URL.Query()

	for key, values := range v {
		if "id" != key && "password" != key && "newpassword" != key {
			sqlStr := func() string {
				if "province" == key {
					lib.SQLExec("update heartbeat set province=? where id=?", values[0], id)

					if 0 == gender {
						return "update girls set province=?,district='' where id=?"
					} else {
						return "update guys set province=?,district='' where id=?"
					}
				} else {
					if 0 == gender {
						return "update girls set " + key + "=? where id=?"
					} else {
						return "update guys set " + key + "=? where id=?"
					}
				}
			}()

			lib.SQLExec(sqlStr, values[0], id)
		}
	}

	newpassword := v.Get("newpassword")
	if "" != newpassword {
		sentense := lib.SQLSentence(lib.SQLMAP_Update_Password, gender)
		_, err := lib.SQLExec(sentense, newpassword, id)
		if nil != err {
			fmt.Println(err.Error())
			return 404
		}
	}

	lib.DelRedisUserInfo(id)
	return 200
}

/*
 *
 *    Function: SetProfile
 *      Author: sunchao
 *        Date: 15/7/10
 * Description: 配置用户属性
 *
 */
func SetProfile(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	result := UpdateProfile(req, id, gender)
	if 200 != result {
		return 404, ""
	}

	code, info := GetUserInfo(id, gender)
	jsonRlt, _ := json.Marshal(info)

	return code, string(jsonRlt)
}

/*
 *
 *    Function: GetPersonInfo
 *      Author: sunchao
 *        Date: 15/10/5
 * Description: 获取人物信息
 *
 */
func GetPersonInfo(req *http.Request) (int, string) {
	v := req.URL.Query()
	idstr := v.Get("id")

	id, _ := strconv.Atoi(idstr)
	exist, gender := getGenderById(id)
	if false == exist {
		return 404, ""
	}

	code, info := GetUserInfo(id, gender)
	if 200 != code {
		return 404, ""
	}

	jsonRlt, _ := json.Marshal(info)
	return 200, string(jsonRlt)
}

func DeleteUserImage(id, gender int, imagename string) {
	sentence := lib.SQLSentence(lib.SQLMAP_Delete_Picture, gender)
	lib.SQLExec(sentence, id, imagename)

	lib.DeleteImageFromQiniu(id, imagename)
}

/*
 *
 *    Function: PostImage
 *      Author: sunchao
 *        Date: 15/10/6
 * Description: 上传图片处理
 *
 */
func PostImage(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	pictypestr := v.Get("pictype")
	if "" == pictypestr {
		return 404, ""
	}

	file, handle, err := req.FormFile("file")
	if nil != err {
		return 404, ""
	}

	defer file.Close()

	strslice := strings.Split(handle.Filename, ".")
	subfix := strslice[len(strslice)-1]
	if "" == subfix {
		return 404, ""
	}

	imagename := lib.RandStringBytesMaskImprSrc(32) + "." + subfix
	err = lib.PutImageToQiniu(id, imagename, file)
	if nil != err {
		return 404, ""
	}

	pictype, _ := strconv.Atoi(pictypestr)

	//头像只有一个,需要先删除旧的
	if 0 == pictype {
		var oldfilename string

		//delete picture from Qiniu
		sentence := lib.SQLSentence(lib.SQLMAP_Select_SearchPictures, gender)
		lib.SQLQueryRow(sentence, id, 1).Scan(&oldfilename)
		lib.DeleteImageFromQiniu(id, oldfilename)

		//delete database
		sentence = lib.SQLSentence(lib.SQLMAP_Delete_Picture, gender)
		lib.SQLExec(sentence, id, oldfilename)
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Picture, gender)
	_, err = lib.SQLExec(sentence, id, imagename, ([2]int{1, 0})[pictype])
	if nil != err {
		return 404, ""
	}

	sentence = lib.SQLSentence(lib.SQLMAP_Update_InfoPictureFlag, gender)
	lib.SQLExec(sentence, id)

	lib.DelRedisUserInfo(id)
	_, info := GetUserInfo(id, gender)
	jsonRlt, _ := json.Marshal(info)

	return 200, string(jsonRlt)
}

/*
 *
 *    Function: DeleteImage
 *      Author: sunchao
 *        Date: 15/10/7
 * Description: delete image
 *
 */
func DeleteImage(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	filename := v.Get("filename")
	if "" == filename {
		return 404, ""
	}

	strslice := strings.Split(filename, "/")
	imagename := strslice[len(strslice)-1]
	DeleteUserImage(id, gender, imagename)

	lib.DelRedisUserInfo(id)
	_, info := GetUserInfo(id, gender)
	jsonRlt, _ := json.Marshal(info)
	return 200, string(jsonRlt)
}

/*
 |    Function: getSearchBaselineSQLSentence
 |      Author: Mr.Sancho
 |        Date: 2016-02-07
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func getSearchBaselineSQLSentence(gender int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) string {
	sentence := lib.SQLSentence(lib.SQLMAP_Select_SearchCountHead, 1-gender)

	if "在线" == status {
		sentence += "((usertype=1 and onlineStatus=1) or usertype!=1) and "
	}
	if 0 != agemax {
		sentence += "age<=" + strconv.Itoa(agemax) + " and "
	}
	if 0 != agemin {
		sentence += "age>=" + strconv.Itoa(agemin) + " and "
	}
	if 0 != heightmax {
		sentence += "height<=" + strconv.Itoa(heightmax) + " and "
	}
	if 0 != heightmin {
		sentence += "height>=" + strconv.Itoa(heightmin) + " and "
	}
	if 0 != incomemax {
		sentence += "incomemax<=" + strconv.Itoa(incomemax) + " and "
	}
	if 0 != incomemin {
		sentence += "incomemin>=" + strconv.Itoa(incomemin) + " and "
	}
	if "" != education {
		sentence += "education='" + education + "' and "
	}
	if "" != occupation {
		sentence += "occupation='" + occupation + "' and "
	}
	if "" != province {
		sentence += "province='" + province + "' and "
	}

	sentence += "pictureflag=1"
	return sentence
}

func getSearchBaseline(gender, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) int {
	var baseline int

	baseline, exist := lib.GetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status)
	if true != exist {
		countsentence := getSearchBaselineSQLSentence(gender, agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status)
		lib.SQLQueryRow(countsentence).Scan(&baseline)
		lib.SetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status, baseline)
	}

	return baseline
}

/*
 |    Function: getUserSearchIndexSQLSentence
 |      Author: Mr.Sancho
 |        Date: 2016-02-07
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func getUserSearchIndexSQLSentence(gender int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) string {
	sentence := lib.SQLSentence(lib.SQLMAP_Select_SearchHead, 1-gender)

	if 0 != agemax {
		sentence += "age<=" + strconv.Itoa(agemax) + " and "
	}
	if 0 != agemin {
		sentence += "age>=" + strconv.Itoa(agemin) + " and "
	}
	if 0 != heightmax {
		sentence += "height<=" + strconv.Itoa(heightmax) + " and "
	}
	if 0 != heightmin {
		sentence += "height>=" + strconv.Itoa(heightmin) + " and "
	}
	if 0 != incomemax {
		sentence += "incomemax<=" + strconv.Itoa(incomemax) + " and "
	}
	if 0 != incomemin {
		sentence += "incomemin>=" + strconv.Itoa(incomemin) + " and "
	}
	if "" != education {
		sentence += "education='" + education + "' and "
	}
	if "" != occupation {
		sentence += "occupation='" + occupation + "' and "
	}
	if "" != province {
		sentence += "province='" + province + "' and "
	}

	sentence += "pictureflag=1 order by id desc limit ?,?"
	return sentence
}

func getUserSearchIndex(id, gender int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, count int) int {
	var baseline int
	var index int

	baseline = getSearchBaseline(gender, agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status)
	index, exist := lib.GetRedisSearchIndex(id, agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status)
	if true != exist {
		index = lib.Intn(baseline - count)
	} else {
		if (index + count) > baseline {
			index = 0
		}
	}

	lib.SetRedisSearchIndex(id, agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status, index+count)
	return index
}

/*
 *
 *    Function: Search
 *      Author: sunchao
 *        Date: 15/7/26
 * Description: handler the search request
 *
 */
func Search(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	var agemin, agemax int
	var heightmin, heightmax int
	var incomemin, incomemax int
	var province, education, occupation, status string

	useheartbeat := true

	queries, _ := url.ParseQuery(req.URL.RawQuery)
	for key, values := range queries {
		v := values[0]
		if "" != v && "0" != v && "不限" != v {
			switch key {
			case "Status":
				{
					useheartbeat = false
					status = v
				}
			case "AgeMax":
				{
					useheartbeat = false
					agemax, _ = strconv.Atoi(v)
				}
			case "AgeMin":
				{
					useheartbeat = false
					agemin, _ = strconv.Atoi(v)
				}
			case "HeightMax":
				{
					useheartbeat = false
					heightmax, _ = strconv.Atoi(v)
				}
			case "HeightMin":
				{
					useheartbeat = false
					heightmin, _ = strconv.Atoi(v)
				}
			case "IncomeMax":
				{
					useheartbeat = false
					incomemax, _ = strconv.Atoi(v)
				}
			case "IncomeMin":
				{
					useheartbeat = false
					incomemin, _ = strconv.Atoi(v)
				}
			case "Study":
				{
					useheartbeat = false
					education = v
				}
			case "Work":
				{
					useheartbeat = false
					occupation = v
				}
			case "Province":
				{
					useheartbeat = false
					province = v
				}
			}
		}
	}

	page, count := lib.Get_pageid_count_fromreq(req)
	if 1 == gender && true == useheartbeat && page <= 2 {
		return doReqHeartbeat(id, gender, count)
	}

	/*
	 * 男搜索女，默认年龄为上3下7. 女搜索男，默认年龄为上7下3.
	 */
	if 0 == agemin && 0 == agemax {
		_, userinfo := GetUserInfo(id, gender)
		agemin, agemax = getSearchAgeRange(userinfo.Gender, userinfo.Age)
	}

	sentence := getUserSearchIndexSQLSentence(gender, agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status)
	index := getUserSearchIndex(id, gender, agemin, agemax, heightmin, heightmax, incomemin, incomemax, province, education, occupation, status, count)
	rows, err := lib.SQLQuery(sentence, index, count)
	if nil != err {
		return 404, ""
	}
	defer rows.Close()

	var info personInfo
	var idtmp int
	var code int

	infos := make([]personInfo, 0)
	gender = 1 - gender
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

/*
 *
 *    Function: Register
 *      Author: sunchao
 *        Date: 15/6/22
 * Description: New user register
 *
 */
func Register(req *http.Request) (int, string) {
	v := req.URL.Query()
	ageStr := v.Get("age")
	genderStr := v.Get("gender")
	cid := v.Get("cid")

	if "" == ageStr || "" == genderStr || "" == cid {
		return 404, ""
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Delete_MultiClientID, 0)
	lib.SQLExec(sentence, cid)

	sentence = lib.SQLSentence(lib.SQLMAP_Delete_MultiClientID, 1)
	lib.SQLExec(sentence, cid)

	girlLastIdSentence := lib.SQLSentence(lib.SQLMAP_Select_LastId, 0)
	guylLastIdSentence := lib.SQLSentence(lib.SQLMAP_Select_LastId, 1)
	gender, _ := strconv.Atoi(genderStr)
	insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_Info, gender)

	g_regLock.Lock()
	defer g_regLock.Unlock()

	/* First get the girls last id */
	var girlsLastId int
	lib.SQLQueryRow(girlLastIdSentence).Scan(&girlsLastId)
	if 0 == girlsLastId {
		return 404, ""
	}

	var guysLastId int
	lib.SQLQueryRow(guylLastIdSentence).Scan(&guysLastId)
	if 0 == guysLastId {
		return 404, ""
	}

	var lastId = func() int {
		if girlsLastId > guysLastId {
			return girlsLastId + 1
		} else {
			return guysLastId + 1
		}
	}()

	var blacklistlastid int
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMAP_Select_BlacklistLastId)).Scan(&blacklistlastid)
	if blacklistlastid == lastId {
		lastId = lastId + 1
	}

	password := strconv.Itoa((lib.Intn(1000000) + lib.Intn(1000000)) % 1000000)
	name := strconv.Itoa(lastId)
	gender, _ = strconv.Atoi(genderStr)
	age, _ := strconv.Atoi(ageStr)
	height := [2]int{160, 175}[gender]
	weight := [2]int{45, 65}[gender]

	province, district := GetIpAddress(req)
	_, err := lib.SQLExec(insertSentence, lastId, password, name, gender, lib.CurrentTimeUTCInt64(), age, common.USERTYPE_USER, cid, height, weight, province, district, 0, 0)
	if nil == err {
		var info registerInfo
		info.Id = lastId
		info.PassWord = password
		info.Member.Id = lastId
		info.Member.Name = strconv.Itoa(lastId)
		info.Member.Age = age
		info.Member.OnlineStatus = 1
		info.Member.Gender = gender
		info.Member.Height = height
		info.Member.Weight = weight
		info.Member.Province = province
		info.Member.District = district
		info.ReviewAlertInfo.ShowReviewAlert = false
		info.ReviewAlertInfo.ReviewAlertMsg = "主人，赏赐个评价吧~"
		info.ReviewAlertInfo.ReviewAlertCancel = "残忍拒绝"
		info.ReviewAlertInfo.ReviewAlertGo = "赏个5星"

		//发送欢迎信息
		go func() {
			msg := "欢迎你来到寂寞同城搭讪"
			timevalue := lib.CurrentTimeUTCInt64()
			RecommendInsertMessageToDB(1, lastId, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
			lib.SetRedisDistrict(lastId, province)
			RecommendPushMessage(1, lastId, 1, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
			push.DoPush()
		}()

		//add the count of user
		g_Count_Regist = g_Count_Regist + 1
		updateUserActive()
		if 0 == gender {
			g_Count_Girls = g_Count_Girls + 1
		} else {
			g_Count_Guys = g_Count_Guys + 1
		}

		go updateLiveUserInfo(g_liveUsersInfo, lastId, gender, LIVEUSER_STATUS_ONLINE, LIVEUSER_TICK_ONLINE)

		jsonRlt, _ := json.Marshal(info)
		return 200, string(jsonRlt)
	} else {
		fmt.Println(err.Error())
		return 404, ""
	}
}

/*
 *
 *    Function: Login
 *      Author: sunchao
 *        Date: 15/6/22
 * Description: user login
 *
 */
func Login(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	/* set cid */
	v := req.URL.Query()
	cid := v.Get("cid")
	if "" == cid {
		return 404, ""
	}

	_, usertype := getUsertypeByIdGender(id, gender)
	if common.USERTYPE_USER != usertype {
		return 403, ""
	}

	sentence := lib.SQLSentence(lib.SQLMAP_Delete_MultiClientID, 0)
	lib.SQLExec(sentence, cid)

	sentence = lib.SQLSentence(lib.SQLMAP_Delete_MultiClientID, 1)
	lib.SQLExec(sentence, cid)

	var lastlogintime, curlogintime int64

	sentence = lib.SQLSentence(lib.SQLMAP_Select_LastLoginTime, gender)
	lib.SQLQueryRow(sentence, id).Scan(&lastlogintime)

	curlogintime = lib.CurrentTimeUTCInt64()
	sentence = lib.SQLSentence(lib.SQLMAP_Update_LoginInfo, gender)
	_, err := lib.SQLExec(sentence, cid, curlogintime, id)
	if nil != err {
		fmt.Println(err)
		return 404, ""
	}

	_, exist = lib.GetRedisDistrict(id)
	if true != exist {
		province, _ := GetIpAddress(req)
		lib.SetRedisDistrict(id, province)
	}

	go visitLoginProc(id, gender, lastlogintime, curlogintime)
	go updateLiveUserInfo(g_liveUsersInfo, id, gender, LIVEUSER_STATUS_ONLINE, LIVEUSER_TICK_ONLINE)

	lib.DelRedisUserInfo(id)
	code, info := GetUserInfo(id, gender)
	jsonRlt, _ := json.Marshal(info)

	return code, string(jsonRlt)
}

/*
 *
 *    Function: WatchDog
 *      Author: sunchao
 *        Date: 15/10/18
 * Description: 狗叫服务
 *
 */
func WatchDog(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	ok := updateLiveUserInfo(g_liveUsersInfo, id, gender, LIVEUSER_STATUS_ONLINE, LIVEUSER_TICK_ONLINE)
	if true != ok {
		if _, exist := lib.GetRedisDistrict(id); false == exist {
			province, _ := GetIpAddress(req)
			lib.SetRedisDistrict(id, province)
		}
	}

	onlineProc(id, gender)
	return 200, ""
}

/*
 *
 *    Function: Logout
 *      Author: sunchao
 *        Date: 15/10/11
 * Description: 用户退出
 *
 */
func Logout(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	deleteLiveUserInfo(g_liveUsersInfo, id)
	OfflineProc(id, gender)
	lib.DelRedisDistrict(id)

	return 200, ""
}

/*
 *
 *    Function: liveUserGoRoute
 *      Author: sunchao
 *        Date: 15/11/4
 * Description: 在线用户管理线程
 *
 */
func liveUserGoRoute() {
	var err error
	var onlineid int
	var onlineStatus int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_OnlineIds, 0)
	rows, err := lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&onlineid, &onlineStatus)
			if nil == err {
				if LIVEUSER_STATUS_ONLINE == onlineStatus {
					updateLiveUserInfo(g_liveUsersInfo, onlineid, 0, LIVEUSER_STATUS_ONLINE, LIVEUSER_TICK_ONLINE)
				} else {
					updateLiveUserInfo(g_liveUsersInfo, onlineid, 0, LIVEUSER_STATUS_BACKGROUND, LIVEUSER_TICK_BACKGROUND)
				}
			}
		}

		rows.Close()
	}

	sentence = lib.SQLSentence(lib.SQLMAP_Select_OnlineIds, 1)
	rows, err = lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&onlineid, &onlineStatus)
			if nil == err {
				if LIVEUSER_STATUS_ONLINE == onlineStatus {
					updateLiveUserInfo(g_liveUsersInfo, onlineid, 1, LIVEUSER_STATUS_ONLINE, LIVEUSER_TICK_ONLINE)
				} else {
					updateLiveUserInfo(g_liveUsersInfo, onlineid, 1, LIVEUSER_STATUS_BACKGROUND, LIVEUSER_TICK_BACKGROUND)
				}
			}
		}

		rows.Close()
	}

	for {
		time.Sleep(lib.SLEEP_DURATION_LIVESTATUS)

		g_liveUsersInfo.lock.Lock()
		for id, user := range g_liveUsersInfo.users {
			user.livetick = user.livetick - 1
			if 0 == user.livetick {
				if LIVEUSER_STATUS_BACKGROUND == user.status {
					delete(g_liveUsersInfo.users, id)
					OfflineProc(id, user.gender)
					lib.DelRedisDistrict(id)
				} else {
					user.status = LIVEUSER_STATUS_BACKGROUND
					user.livetick = LIVEUSER_TICK_BACKGROUND

					backgroundProc(id, user.gender)
				}
			}

			user.livetime = user.livetime + 1
			if 1 == user.livetime && LIVEUSER_STATUS_ONLINE == user.status {
				if true != checkIfUserHavePicture(id, user.gender) {
					//管理员发送第二封信
					msg := "您还没有更新照片哦,上传照片获得更高的推荐机会!"
					timevalue := lib.CurrentTimeUTCInt64()
					RecommendInsertMessageToDB(1, id, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
					RecommendPushMessage(1, id, 1, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
					push.DoPush()
				}
			}
		}
		g_liveUsersInfo.lock.Unlock()
	}
}

/*
 |    Function: liveUserNotifyMsgRoutine
 |      Author: Mr.Sancho
 |        Date: 2016-01-13
 |   Arguments:
 |      Return:
 | Description: notify live user the messages unreaded
 |
*/
func liveUserNotifyMsgRoutine() {
	needpush := false

	for {
		time.Sleep(lib.SLEEP_DURATION_NOTIFYMSG)
		needpush = false

		g_liveUsersInfo.lock.Lock()
		for id, user := range g_liveUsersInfo.users {
			if LIVEUSER_STATUS_ONLINE == user.status {
				recommendCount := recommend_GetUnreadNum(id)
				visitCount := visit_GetUnreadNum(id)

				unreadmsg := unreadMessage{UnreadRecommend: recommendCount, UnreadVisit: visitCount, Badge: recommendCount + visitCount}
				jsonRlt, _ := json.Marshal(unreadmsg)
				notifymsg := notifyMessageInfo{Type: push.PUSH_NOTIFYMSG_UNREAD, Value: string(jsonRlt)}
				jsonRlt, _ = json.Marshal(notifymsg)

				clientid := GetClientIdByUserId(id)
				push.Add(recommendCount+visitCount, clientid, push.PUSHMSG_TYPE_NOTIFYMSG, push.PUSH_NOTIFYMSG_UNREAD, "", string(jsonRlt))
				needpush = true
			}
		}
		g_liveUsersInfo.lock.Unlock()

		if true == needpush {
			push.DoPush()
		}
	}
}

/*
 *
 *    Function: vipGoRoute
 *      Author: sunchao
 *        Date: 15/11/12
 * Description: vip用户处理
 *
 */
func vipUserGoRoute() {
	var err error
	var userid int
	var level int
	var days int
	var expiretime int64

	curtime := lib.CurrentTimeUTCInt64()

	sentence := lib.SQLSentence(lib.SQLMAP_Select_VIPRows, 0)
	rows, err := lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&userid, &level, &days, &expiretime)
			if nil == err {
				if curtime >= expiretime {
					detachVipFromUser(userid, 0, level)
				} else {
					updateVipUserInfo(g_vipUsersInfo, userid, 0, level, days, expiretime)
				}
			}
		}

		rows.Close()
	}

	sentence = lib.SQLSentence(lib.SQLMAP_Select_VIPRows, 1)
	rows, err = lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&userid, &level, &days, &expiretime)
			if nil == err {
				if curtime >= expiretime {
					detachVipFromUser(userid, 1, level)
				} else {
					updateVipUserInfo(g_vipUsersInfo, userid, 1, level, days, expiretime)
				}
			}
		}

		rows.Close()
	}

	needpush := false
	for {
		time.Sleep(lib.SLEEP_DURATION_VIPSTATUS)
		needpush = false

		g_vipUsersInfo.lock.Lock()
		for id, user := range g_vipUsersInfo.users {
			curtime = lib.CurrentTimeUTCInt64()

			if curtime >= user.expiretime {
				detachVipFromUser(id, user.gender, user.level)
				needpush = true
			} else {
				difhours := int(time.Duration(user.expiretime-curtime) / time.Hour)
				if 24 == difhours || 120 == difhours || 240 == difhours {
					//发送信息, 提醒到期
					msg := "您的 " + [...]string{"1级会员", "2级会员", "3级会员"}[user.level] + " 将在 " + strconv.Itoa(difhours/24) + " 天后到期."
					RecommendInsertMessageToDB(1, id, RECOMMEND_MSGTYPE_TALK, msg, curtime)
					RecommendPushMessage(1, id, 1, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, curtime)
					needpush = true
				}
			}
		}
		g_vipUsersInfo.lock.Unlock()

		if true == needpush {
			push.DoPush()
		}
	}
}

func detachVipFromUser(id, gender, level int) {
	deleteVipUserInfo(g_vipUsersInfo, id)

	sentence := lib.SQLSentence(lib.SQLMAP_Update_VIPById, gender)
	_, err := lib.SQLExec(sentence, 0, 0, 0, id)
	if nil == err {
		//发送信息, VIP权限已经取消
		msg := "非常抱歉通知您, 您的 " + []string{"初始会员", "1级会员", "2级会员", "3级会员"}[level] + " 已经到期!"
		timevalue := lib.CurrentTimeUTCInt64()
		RecommendInsertMessageToDB(1, id, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
		RecommendPushMessage(1, id, 1, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
	}
}

/*
 *
 *    Function: BuyVip
 *      Author: sunchao
 *        Date: 15/11/8
 * Description: 购买服务后更新后台数据
 *
 */
func BuyVip(req *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	v := req.URL.Query()
	levelstr := v.Get("level")
	daysstr := v.Get("days")

	if "" == levelstr || "" == daysstr {
		return 404, ""
	}

	level, _ := strconv.Atoi(levelstr)
	days, _ := strconv.Atoi(daysstr)

	/* check if ther user already buy VIP */
	var oldlevel int
	var olddays int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_VipLevelByID, gender)
	lib.SQLQueryRow(sentence, id).Scan(&oldlevel, &olddays)
	if 0 != oldlevel {
		if oldlevel > level {
			level = oldlevel
		}

		days = days + olddays
	}

	//秒为单位
	expiretime := lib.CurrentTimeUTCInt64() + int64(days)*int64(time.Hour/time.Second)*24
	sentence = lib.SQLSentence(lib.SQLMAP_Update_VIPById, gender)
	_, err := lib.SQLExec(sentence, level, days, expiretime, id)
	if nil != err {
		return 404, err.Error()
	}

	g_Count_BuyVIP = g_Count_BuyVIP + 1
	updateUserActive()

	//更新到线程
	go updateVipUserInfo(g_vipUsersInfo, id, gender, level, days, expiretime)

	//发送信息, VIP已经开通
	expireUTC := lib.Int64_To_UTCTime(expiretime)
	msg := "您的 " + []string{"初始会员", "写信会员", "钻石会员", "至尊会员"}[level] + " 已经开通啦! 会员到期时间：" +
		fmt.Sprintf("%d年%d月%d日", expireUTC.Year(), expireUTC.Month(), expireUTC.Day())
	timevalue := lib.CurrentTimeUTCInt64()
	RecommendInsertMessageToDB(1, id, RECOMMEND_MSGTYPE_TALK, msg, timevalue)
	RecommendPushMessage(1, id, 1, 1, push.PUSHMSG_TYPE_RECOMMEND, msg, timevalue)
	push.DoPush()

	lib.DelRedisUserInfo(id)
	code, info := GetUserInfo(id, gender)
	jsonRlt, _ := json.Marshal(info)

	return code, string(jsonRlt)
}

/*
 *
 *    Function: VipPrice
 *      Author: sunchao
 *        Date: 15/7/11
 * Description: get the vip price information
 *
 */
func VipPrice() (int, string) {
	jsonRlt, _ := json.Marshal(g_VipLevels)
	return 200, string(jsonRlt)
}

func GetBuyVIPCount() int {
	return g_Count_BuyVIP
}

func SubUserCount(gender int) {
	if 0 == gender {
		if g_Count_Girls > 0 {
			g_Count_Girls = g_Count_Girls - 1
		}
	} else {
		if g_Count_Guys > 0 {
			g_Count_Guys = g_Count_Guys - 1
		}
	}
}