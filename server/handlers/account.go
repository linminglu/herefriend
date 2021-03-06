package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"

	"herefriend/common"
	"herefriend/config"
	"herefriend/lib"
	"herefriend/lib/push"
)

const (
	// LiveUserTickOnline .
	LiveUserTickOnline = 15 // minutes
	// LiveUserTickBackGround .
	LiveUserTickBackGround = 2880 // minutes of 2 days
	// VIPUserSleepDuration .
	VIPUserSleepDuration = int64(time.Hour)
	// LiveUserStatusOnline .
	LiveUserStatusOnline = 1
	// LiveUserStatusBackGround .
	LiveUserStatusBackGround = 2
)

var gRegLock sync.Mutex
var gLiveUsersInfo *liveUsersInfo
var gVipUsersInfo *vipUsersInfo
var gAddAgeNum = []int{7, 3}
var gSubAgeNum = []int{3, 7}

var gCountGirls int
var gCountGuys int
var gCountRegist int
var gCountBuyVIP int

func init() {
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMapSelectUserCount, 0)).Scan(&gCountGirls)
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMapSelectUserCount, 1)).Scan(&gCountGuys)
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMapSelectUserActive)).Scan(&gCountRegist, &gCountBuyVIP)

	gLiveUsersInfo = &liveUsersInfo{users: make(map[int]*liveUser), lock: sync.RWMutex{}}
	go liveUserGoRoute()

	gVipUsersInfo = &vipUsersInfo{users: make(map[int]*vipUser), lock: sync.RWMutex{}}
	go vipUserGoRoute()
}

func updateUserActive() {
	lib.SQLExec(lib.SQLSentence(lib.SQLMapUpdateUserActive), gCountRegist, gCountBuyVIP)
}

// GetUserCountByGender .
func GetUserCountByGender(gender int) int {
	if 0 == gender {
		return gCountGirls
	}

	return gCountGuys
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
		if LiveUserStatusBackGround == info.status {
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

// DeleteLiveUser .
func DeleteLiveUser(id int) {
	gCountRegist = gCountRegist - 1
	updateUserActive()

	deleteLiveUserInfo(gLiveUsersInfo, id)
}

// GetLiveUserNumber .
func GetLiveUserNumber() int {
	return len(gLiveUsersInfo.users)
}

// GetActiveUserNumber .
func GetActiveUserNumber() int {
	var num int

	for _, user := range gLiveUsersInfo.users {
		if LiveUserStatusOnline == user.status {
			num = num + 1
		}
	}

	return num
}

// GetRegistUserNumber .
func GetRegistUserNumber() int {
	return gCountRegist
}

// UpdateVipUserInfo .
func UpdateVipUserInfo(id, gender, level, days int, expiretime int64) bool {
	gVipUsersInfo.lock.Lock()
	info, ok := gVipUsersInfo.users[id]
	if true == ok {
		info.level = level
		info.days = days
		info.expiretime = expiretime
	} else {
		gVipUsersInfo.users[id] = &vipUser{gender: gender, level: level, days: days, expiretime: expiretime}
	}

	gVipUsersInfo.lock.Unlock()

	return ok
}

func onlineProc(id, gender int) {
	sentence := lib.SQLSentence(lib.SQLMapUpdateOnline, gender)
	lib.SQLExec(sentence, lib.CurrentTimeUTCInt64(), id)
}

func backgroundProc(id, gender int) {
	sentence := lib.SQLSentence(lib.SQLMapUpdateBackground, gender)
	lib.SQLExec(sentence, id)
}

// OfflineProc .
func OfflineProc(id, gender int) {
	sentence := lib.SQLSentence(lib.SQLMapUpdateOffline, gender)
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
	if config.ConfAgeMin <= age && age <= config.ConfAgeMax {
		if 0 == gender {
			min = age - gSubAgeNum[gender]
			max = age + gAddAgeNum[gender]

			if min < config.ConfAgeMin {
				min = config.ConfAgeMin
			}

			if max > config.ConfAgeMax {
				max = config.ConfAgeMax
			}
		} else {
			min = 18
			max = 28
		}
	}

	return min, max
}

func getRandomHeartbeatID(id, gender int) int {
	var tmpid int
	var baselimit int
	var sentence string

	province, exist := lib.GetRedisDistrict(id)
	if true == exist {
		sentence = lib.SQLSentence(lib.SQLMapSelectHeartbeatRandomProvID, gender)

		baselimit = getHeartbeatBaseCountByProvinceGender(province, gender)
		randomvalue := lib.Intn(baselimit)
		err := lib.SQLQueryRow(sentence, province, randomvalue).Scan(&tmpid)
		if nil != err {
			lib.SQLError(sentence, err, province, randomvalue)
		} else if 0 != tmpid {
			return tmpid
		}
	}

	return getRandomUserID(id, gender)
}

// getRandomUserID 根据用户的id与性别,获取一个合适的随机id
func getRandomUserID(id, gender int) int {
	var tmpid int
	var baselimit int

	province, exist := lib.GetRedisDistrict(id)
	if true == exist {
		/* the first step: get random id by province and age range
		 */
		getcount := false
		rangecount := 0
		_, info := GetUserInfo(id, gender)
		if config.ConfAgeMin <= info.Age && info.Age <= config.ConfAgeMax {
			min, max := getSearchAgeRange(info.Gender, info.Age)
			if 0 != min {
				for tmpage := min; tmpage <= max; tmpage = tmpage + 1 {
					count, exist := lib.GetRedisProvAgeCount(province, gender, tmpage)
					if true == exist {
						rangecount = rangecount + count
					} else {
						sentence := lib.SQLSentence(lib.SQLMapSelectCountByProvAge, gender)
						err := lib.SQLQueryRow(sentence, province, tmpage).Scan(&count)
						if nil == err && 0 != count {
							lib.SetRedisProvAgeCount(province, gender, tmpage, count)
							rangecount = rangecount + count
						} else if nil != err {
							lib.SQLError(sentence, err, province, tmpage)
						}
					}
				}

				if 0 != rangecount {
					sentence := lib.SQLSentence(lib.SQLMapSelectRandomProvAgeID, gender)
					baselimit = lib.Intn(rangecount)
					err := lib.SQLQueryRow(sentence, province, min, max, baselimit).Scan(&tmpid)
					if nil == err && 1 < tmpid {
						return tmpid
					} else if nil != err {
						lib.SQLError(sentence, err, province, min, max, baselimit)
					}
				}
			}
		}

		/* the second step: get random id by province
		 */
		getcount = false
		count, exist := lib.GetRedisProvCount(province, gender)
		if false == exist {
			sentence := lib.SQLSentence(lib.SQLMapSelectCountByProv, gender)
			err := lib.SQLQueryRow(sentence, province).Scan(&count)
			if nil == err && 0 != count {
				lib.SetRedisProvCount(province, gender, count)
				getcount = true
			} else if nil != err {
				lib.SQLError(sentence, err, province)
			}
		} else {
			getcount = true
		}

		if true == getcount {
			sentence := lib.SQLSentence(lib.SQLMapSelectRandomProvID, gender)
			baselimit = lib.Intn(count)
			err := lib.SQLQueryRow(sentence, province, baselimit).Scan(&tmpid)
			if nil == err && 1 < tmpid {
				return tmpid
			} else if nil != err {
				lib.SQLError(sentence, err, province, baselimit)
			}
		}
	}

	/* the third step: get random id from all data
	 */
	if 0 == gender {
		baselimit = lib.Intn(gCountGirls)
	} else {
		baselimit = lib.Intn(gCountGuys)
	}

	sentence := lib.SQLSentence(lib.SQLMapSelectRandomID, gender)
	err := lib.SQLQueryRow(sentence, baselimit).Scan(&tmpid)
	if nil != err || 1 >= tmpid {
		lib.SQLQueryRow(sentence, 2).Scan(&tmpid)
	}

	return tmpid
}

// GetGenderUsertypeByID get gender usertype by id, if return false means there is no such id in tables
func GetGenderUsertypeByID(id int) (bool, int, int) {
	if 0 != id {
		var idtmp int
		var usertype int

		sentence := lib.SQLSentence(lib.SQLMapSelectUserType, 0)
		err := lib.SQLQueryRow(sentence, id).Scan(&idtmp, &usertype)
		if nil == err && id == idtmp {
			return true, 0, usertype
		}

		sentence = lib.SQLSentence(lib.SQLMapSelectUserType, 1)
		err = lib.SQLQueryRow(sentence, id).Scan(&idtmp, &usertype)
		if nil == err && id == idtmp {
			return true, 1, usertype
		}
	}

	return false, 1, 0
}

// GetUsertypeByIDGender 根据id和性别获取用户类型
func GetUsertypeByIDGender(id, gender int) (bool, int) {
	var idtmp int
	var usertype int

	sentence := lib.SQLSentence(lib.SQLMapSelectUserType, gender)
	err := lib.SQLQueryRow(sentence, id).Scan(&idtmp, &usertype)
	if nil == err && id == idtmp {
		return true, usertype
	}

	return false, 0
}

// getGenderByID get gender by id, if return false means there is no such id in tables
func getGenderByID(id int) (bool, int) {
	gender, exist := lib.GetRedisUserGender(id)
	if true == exist {
		return true, gender
	}

	var idstr int

	sentence := lib.SQLSentence(lib.SQLMapSelectCheckIsValidID, 0)
	err := lib.SQLQueryRow(sentence, id).Scan(&idstr)
	if nil == err {
		lib.SetRedisUserGender(id, 0)
		return true, 0
	}

	sentence = lib.SQLSentence(lib.SQLMapSelectCheckIsValidID, 1)
	err = lib.SQLQueryRow(sentence, id).Scan(&idstr)
	if nil == err {
		lib.SetRedisUserGender(id, 1)
		return true, 1
	}

	return false, 1
}

// getGenderByIDPw get gender by id and password, if return false means there is no such id in tables
func getGenderByIDPw(id int, pw string) (bool, int) {
	if 0 != id {
		var idScan int

		sentence := lib.SQLSentence(lib.SQLMapSelectCheckIsValidPasswd, 0)
		err := lib.SQLQueryRow(sentence, id, pw).Scan(&idScan)
		if nil == err {
			return true, 0
		}

		sentence = lib.SQLSentence(lib.SQLMapSelectCheckIsValidPasswd, 1)
		err = lib.SQLQueryRow(sentence, id, pw).Scan(&idScan)
		if nil == err {
			return true, 1
		}
	}

	return false, 1
}

// getIDGenderByRequest 根据用户请求获取id和性别
func getIDGenderByRequest(c *gin.Context) (bool, int, int) {
	idStr := c.Query("id")
	pwStr := c.Query("password")
	if "" == idStr || "" == pwStr {
		return false, 0, 0
	}

	id, _ := strconv.Atoi(idStr)
	bExist, gender := getGenderByIDPw(id, pwStr)
	if true != bExist {
		return false, 0, 0
	}

	return true, id, gender
}

func checkIfUserHavePicture(id, gender int) bool {
	var filename string

	sentence := lib.SQLSentence(lib.SQLMapSelectSearchPictures, gender)

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
	var expiretime int64

	sentence := lib.SQLSentence(lib.SQLMapSelectVipLevelByID, gender)
	err := lib.SQLQueryRow(sentence, id).Scan(&viplevel, &vipdays, &expiretime)
	if nil == err && 0 != viplevel {
		return true
	} else if nil != err {
		lib.SQLError(sentence, err, id)
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
func getUserPictrues(id, gender int, info *common.PersonInfo) {
	var filename string

	sentence := lib.SQLSentence(lib.SQLMapSelectSearchPictures, gender)

	/* 获取头像 */
	err := lib.SQLQueryRow(sentence, id, 1).Scan(&filename)
	if nil == err && "" != filename {
		info.IconURL = lib.GetQiniuUserImageURL(id, filename)
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

// GetUserInfoByID .
func GetUserInfoByID(id int) (int, common.PersonInfo) {
	var info common.PersonInfo

	exist, gender := getGenderByID(id)
	if false == exist {
		return 404, info
	}

	code, info := GetUserInfo(id, gender)
	if 200 != code {
		return 404, info
	}

	return 200, info
}

// GetUserGoldBeans 获取用户的金币数量
func GetUserGoldBeans(id int) int {
	var beans int
	var consumed int

	beans, exist := lib.GetRedisGoldBeans(id)
	if true == exist {
		return beans
	}

	sentence := lib.SQLSentence(lib.SQLMapSelectGoldBeansByID)
	lib.SQLQueryRow(sentence, id).Scan(&beans, &consumed)

	lib.SetRedisGoldBeans(id, beans)
	return beans
}

// GetUserRecvGiftList 获取收到的礼物列表
func GetUserRecvGiftList(id int) []common.GiftSendRecvInfo {
	redislist, exist := lib.GetRedisGiftRecvList(id)
	if true == exist {
		return *redislist
	}

	var infolist []common.GiftSendRecvInfo
	sentence := lib.SQLSentence(lib.SQLMapSelectGiftRecvSum)
	rows, err := lib.SQLQuery(sentence, id)
	if nil != err {
		return infolist
	}

	defer rows.Close()

	var giftid, giftnum int
	giftnuminfo := make(map[int]int)
	for rows.Next() {
		err = rows.Scan(&giftid, &giftnum)
		if nil == err {
			giftnuminfo[giftid] = giftnuminfo[giftid] + giftnum
		}
	}

	for k, v := range giftnuminfo {
		infolist = append(infolist, common.GiftSendRecvInfo{GiftID: k, Number: v})
	}

	lib.SetRedisGiftRecvList(id, &infolist)
	return infolist
}

// PrepareUserRecvGiftList .
func PrepareUserRecvGiftList(id int) {
	_, exist := lib.GetRedisGiftRecvList(id)
	if true == exist {
		return
	}

	var infolist []common.GiftSendRecvInfo
	sentence := lib.SQLSentence(lib.SQLMapSelectGiftRecvSum)
	rows, err := lib.SQLQuery(sentence, id)
	if nil != err {
		return
	}

	defer rows.Close()

	var giftid, giftnum int
	giftnuminfo := make(map[int]int)
	for rows.Next() {
		err = rows.Scan(&giftid, &giftnum)
		if nil == err {
			giftnuminfo[giftid] = giftnuminfo[giftid] + giftnum
		}
	}

	for k, v := range giftnuminfo {
		infolist = append(infolist, common.GiftSendRecvInfo{GiftID: k, Number: v})
	}

	lib.SetRedisGiftRecvList(id, &infolist)
	return
}

// GetUserSendGiftList 获取送出的礼物列表
func GetUserSendGiftList(id int) []common.GiftSendRecvInfo {
	redislist, exist := lib.GetRedisGiftSendList(id)
	if true == exist {
		return *redislist
	}

	var infolist []common.GiftSendRecvInfo
	sentence := lib.SQLSentence(lib.SQLMapSelectGiftSendSum)
	rows, err := lib.SQLQuery(sentence, id)
	if nil != err {
		return infolist
	}

	defer rows.Close()

	var giftid, giftnum int
	giftnuminfo := make(map[int]int)
	for rows.Next() {
		err = rows.Scan(&giftid, &giftnum)
		if nil == err {
			giftnuminfo[giftid] = giftnuminfo[giftid] + giftnum
		}
	}

	for k, v := range giftnuminfo {
		infolist = append(infolist, common.GiftSendRecvInfo{GiftID: k, Number: v})
	}

	lib.SetRedisGiftSendList(id, &infolist)
	return infolist
}

// GetUserInfo get the user information by id and gender
func GetUserInfo(id int, gender int) (int, common.PersonInfo) {
	redisinfo, exist := lib.GetRedisUserInfo(id)
	if true == exist {
		return 200, *redisinfo
	}

	var info common.PersonInfo
	var timeValue int64
	sentence := lib.SQLSentence(lib.SQLMapSelectPersonInfo, gender)
	err := lib.SQLQueryRow(sentence, id).Scan(&info.ID, &info.Name, &info.Age, &info.Gender, &info.OnlineStatus, &info.VipLevel, &timeValue, &info.Height, &info.Weight,
		&info.LoveType, &info.BodyType, &info.Marriage, &info.Province, &info.District, &info.Native, &info.Education, &info.Income, &info.IncomeMin, &info.IncomeMax,
		&info.Occupation, &info.Housing, &info.Carstatus, &info.Introduction, &info.School, &info.Speciality,
		&info.Animal, &info.Constellation, &info.Lang, &info.BloodType, &info.Selfjudge, &info.Companytype, &info.Companyindustry,
		&info.Nationnality, &info.Religion, &info.Charactor, &info.Hobbies, &info.CityLove, &info.Naken, &info.AllowAge, &info.AllowResidence,
		&info.AllowHeight, &info.AllowMarriage, &info.AllowEducation, &info.AllowHousing, &info.AllowIncome, &info.AllowKidStatus)
	if nil == err {
		if 0 != info.VipLevel {
			info.VipExpireTime = lib.Int64ToUTCTime(timeValue)
		}

		getUserPictrues(id, gender, &info)
	} else {
		lib.SQLError(sentence, err, id)
		return 404, info
	}

	info.GoldBeans = GetUserGoldBeans(id)
	info.RecvGiftList = GetUserRecvGiftList(id)
	lib.SetRedisUserInfo(id, &info)

	return 200, info
}

// GetClientIDByUserID .
func GetClientIDByUserID(id int) string {
	var clientid string

	sentence := lib.SQLSentence(lib.SQLMapSelectClientID, 0)
	lib.SQLQueryRow(sentence, id).Scan(&clientid)
	if "" == clientid {
		sentence = lib.SQLSentence(lib.SQLMapSelectClientID, 1)
		lib.SQLQueryRow(sentence, id).Scan(&clientid)
	}

	return clientid
}

// UpdateProfile .
func UpdateProfile(req *http.Request, id, gender int) int {
	v := req.URL.Query()

	for key, values := range v {
		if key != "id" && key != "password" && key != "newpassword" && key != "_" {
			sqlStr := func() string {
				if key == "province" {
					lib.SQLExec("update heartbeat set province=? where id=?", values[0], id)

					if gender == 0 {
						return "update girls set province=?,district='' where id=?"
					}
					return "update guys set province=?,district='' where id=?"
				}

				if gender == 0 {
					return "update girls set " + key + "=? where id=?"
				}
				return "update guys set " + key + "=? where id=?"
			}()

			lib.SQLExec(sqlStr, values[0], id)
		}
	}

	newpassword := v.Get("newpassword")
	if "" != newpassword {
		sentense := lib.SQLSentence(lib.SQLMapUpdatePassword, gender)
		_, err := lib.SQLExec(sentense, newpassword, id)
		if nil != err {
			log.Error(err.Error())
			return http.StatusNotFound
		}
	}

	lib.DelRedisUserInfo(id)
	return http.StatusOK
}

// SetProfile 配置用户属性
func SetProfile(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	result := UpdateProfile(c.Request, id, gender)
	if result != http.StatusOK {
		c.Status(http.StatusNotFound)
		return
	}

	code, info := GetUserInfo(id, gender)
	c.JSON(code, info)
}

// GetPersonInfo 获取人物信息
func GetPersonInfo(c *gin.Context) {
	idstr := c.Query("id")
	id, _ := strconv.Atoi(idstr)
	exist, gender := getGenderByID(id)
	if false == exist {
		c.Status(http.StatusNotFound)
		return
	}

	code, info := GetUserInfo(id, gender)
	if http.StatusOK != code {
		c.Status(http.StatusNotFound)
		return
	}

	info.SendGiftList = GetUserSendGiftList(id)
	c.JSON(http.StatusOK, info)
}

// DeleteUserImage .
func DeleteUserImage(id, gender int, imagename string) {
	sentence := lib.SQLSentence(lib.SQLMapDeletePicture, gender)
	lib.SQLExec(sentence, id, imagename)

	lib.DeleteImageFromQiniu(id, imagename)
}

// PostImage 上传图片处理
func PostImage(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	pictypestr := c.Query("pictype")
	if "" == pictypestr {
		log.Error("Failed to get picture type")
		c.Status(http.StatusNotFound)
		return
	}

	file, handle, err := c.Request.FormFile("file")
	if nil != err {
		log.Error(err)
		c.Status(http.StatusNotFound)
		return
	}

	defer file.Close()

	strslice := strings.Split(handle.Filename, ".")
	subfix := strslice[len(strslice)-1]
	if "" == subfix {
		c.Status(http.StatusNotFound)
		return
	}

	imagename := lib.RandStringBytesMaskImprSrc(32) + "." + subfix
	err = lib.PutImageToQiniu(id, imagename, file)
	if nil != err {
		log.Error(err)
		c.Status(http.StatusNotFound)
		return
	}

	pictype, _ := strconv.Atoi(pictypestr)

	//头像只有一个,需要先删除旧的
	if 0 == pictype {
		var oldfilename string

		//delete picture from Qiniu
		sentence := lib.SQLSentence(lib.SQLMapSelectSearchPictures, gender)
		lib.SQLQueryRow(sentence, id, 1).Scan(&oldfilename)
		lib.DeleteImageFromQiniu(id, oldfilename)

		//delete database
		sentence = lib.SQLSentence(lib.SQLMapDeletePicture, gender)
		lib.SQLExec(sentence, id, oldfilename)
	}

	sentence := lib.SQLSentence(lib.SQLMapInsertPicture, gender)
	_, err = lib.SQLExec(sentence, id, imagename, ([2]int{1, 0})[pictype])
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	sentence = lib.SQLSentence(lib.SQLMapUpdateInfoPictureFlag, gender)
	lib.SQLExec(sentence, id)

	lib.DelRedisUserInfo(id)
	_, info := GetUserInfo(id, gender)

	c.JSON(http.StatusOK, info)
}

// DeleteImage delete image
func DeleteImage(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	filename := c.Query("filename")
	if "" == filename {
		c.Status(http.StatusNotFound)
		return
	}

	strslice := strings.Split(filename, "/")
	imagename := strslice[len(strslice)-1]
	DeleteUserImage(id, gender, imagename)

	lib.DelRedisUserInfo(id)
	_, info := GetUserInfo(id, gender)
	c.JSON(http.StatusOK, info)
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
	sentence := lib.SQLSentence(lib.SQLMapSelectCount, 1-gender)

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
	sentence := lib.SQLSentence(lib.SQLMapSelectSearch, 1-gender)

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

// Search handler the search request
func Search(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	var agemin, agemax int
	var heightmin, heightmax int
	var incomemin, incomemax int
	var province, education, occupation, status string

	useheartbeat := true
	queries, _ := url.ParseQuery(c.Request.URL.RawQuery)
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

	page, count := lib.GetPageidCount(c)
	if true == useheartbeat && page <= 2 {
		code, content := doReqHeartbeat(id, gender, count)
		c.String(code, content)
		return
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
		c.Status(http.StatusNotFound)
		return
	}
	defer rows.Close()

	var info common.PersonInfo
	var idtmp int
	var code int

	var infos []common.PersonInfo
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

	c.JSON(http.StatusOK, infos)
}

// Register New user register
func Register(c *gin.Context) {
	ageStr := c.Query("age")
	genderStr := c.Query("gender")
	cid := c.Query("cid")

	if ageStr == "" || genderStr == "" || cid == "" {
		c.Status(http.StatusNotFound)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMapDeleteMultiClientID, 0)
	lib.SQLExec(sentence, cid)

	sentence = lib.SQLSentence(lib.SQLMapDeleteMultiClientID, 1)
	lib.SQLExec(sentence, cid)

	girlLastIDSentence := lib.SQLSentence(lib.SQLMapSelectLastID, 0)
	guylLastIDSentence := lib.SQLSentence(lib.SQLMapSelectLastID, 1)
	gender, _ := strconv.Atoi(genderStr)
	insertSentence := lib.SQLSentence(lib.SQLMapInsertInfo, gender)

	gRegLock.Lock()
	defer gRegLock.Unlock()

	/* First get the girls last id */
	var girlsLastID int
	lib.SQLQueryRow(girlLastIDSentence).Scan(&girlsLastID)
	if 0 == girlsLastID {
		c.Status(http.StatusNotFound)
		return
	}

	var guysLastID int
	lib.SQLQueryRow(guylLastIDSentence).Scan(&guysLastID)
	if 0 == guysLastID {
		c.Status(http.StatusNotFound)
		return
	}

	var lastID = func() int {
		if girlsLastID > guysLastID {
			return girlsLastID + 1
		}
		return guysLastID + 1
	}()

	var blacklistlastid int
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMapSelectBlacklistLastID)).Scan(&blacklistlastid)
	if blacklistlastid == lastID {
		lastID = lastID + 1
	}

	password := strconv.Itoa((lib.Intn(1000000) + lib.Intn(1000000)) % 1000000)
	name := strconv.Itoa(lastID)
	gender, _ = strconv.Atoi(genderStr)
	age, _ := strconv.Atoi(ageStr)
	height := [2]int{160, 175}[gender]
	weight := [2]int{45, 65}[gender]

	province, district := GetIPAddress(c.Request)
	_, err := lib.SQLExec(insertSentence, lastID, password, name, gender, lib.CurrentTimeUTCInt64(), age, common.UserTypeUser, cid, height, weight, province, district, 0, 0)
	if nil == err {
		var info registerInfo
		info.ID = lastID
		info.PassWord = password
		info.Member.ID = lastID
		info.Member.Name = strconv.Itoa(lastID)
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
			msg := config.ConfWelcomeMessage
			timevalue := lib.CurrentTimeUTCInt64()
			RecommendInsertMessageToDB(1, lastID, CommentMsgTypeTalk, msg, timevalue)
			lib.SetRedisDistrict(lastID, province)
			RecommendPushMessage(1, lastID, 1, 1, push.PushMsgComment, msg, timevalue)
			push.DoPush()
		}()

		//add the count of user
		gCountRegist = gCountRegist + 1
		updateUserActive()
		if 0 == gender {
			gCountGirls = gCountGirls + 1
		} else {
			gCountGuys = gCountGuys + 1
		}

		go updateLiveUserInfo(gLiveUsersInfo, lastID, gender, LiveUserStatusOnline, LiveUserTickOnline)

		c.JSON(http.StatusOK, info)
	} else {
		log.Error(err.Error())
		c.Status(http.StatusNotFound)
	}
}

// Login user login
func Login(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	/* set cid */
	cid := c.Query("cid")
	if "" == cid {
		c.Status(http.StatusNotFound)
		return
	}

	_, usertype := GetUsertypeByIDGender(id, gender)
	if common.UserTypeUser != usertype {
		c.Status(http.StatusForbidden)
		return
	}

	sentence := lib.SQLSentence(lib.SQLMapDeleteMultiClientID, 0)
	lib.SQLExec(sentence, cid)

	sentence = lib.SQLSentence(lib.SQLMapDeleteMultiClientID, 1)
	lib.SQLExec(sentence, cid)

	sentence = lib.SQLSentence(lib.SQLMapUpdateLoginInfo, gender)
	_, err := lib.SQLExec(sentence, cid, lib.CurrentTimeUTCInt64(), id)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	_, exist = lib.GetRedisDistrict(id)
	if true != exist {
		province, _ := GetIPAddress(c.Request)
		lib.SetRedisDistrict(id, province)
	}

	go updateLiveUserInfo(gLiveUsersInfo, id, gender, LiveUserStatusOnline, LiveUserTickOnline)

	lib.DelRedisUserInfo(id)
	code, info := GetUserInfo(id, gender)
	info.SendGiftList = GetUserSendGiftList(id)

	c.JSON(code, info)
}

// WatchDog 狗叫服务
func WatchDog(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	if true == config.ConfEvaluationSwitch {
		var lastEvaluationTime int64
		sentence := lib.SQLSentence(lib.SQLMapSelectLastEvaluationTime, gender)
		lib.SQLQueryRow(sentence, id).Scan(&lastEvaluationTime)

		go PeriodOnlineCommentPush(id, gender, lastEvaluationTime)
	}

	ok := updateLiveUserInfo(gLiveUsersInfo, id, gender, LiveUserStatusOnline, LiveUserTickOnline)
	if true != ok {
		go func() {
			if _, exist := lib.GetRedisDistrict(id); false == exist {
				province, _ := GetIPAddress(c.Request)
				lib.SetRedisDistrict(id, province)
			}
		}()
	}

	onlineProc(id, gender)
	c.Status(http.StatusOK)
}

// Logout 用户退出
func Logout(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	deleteLiveUserInfo(gLiveUsersInfo, id)
	OfflineProc(id, gender)
	lib.DelRedisDistrict(id)

	c.Status(http.StatusOK)
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

	sentence := lib.SQLSentence(lib.SQLMapSelectOnlineIDs, 0)
	rows, err := lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&onlineid, &onlineStatus)
			if nil == err {
				if LiveUserStatusOnline == onlineStatus {
					updateLiveUserInfo(gLiveUsersInfo, onlineid, 0, LiveUserStatusOnline, LiveUserTickOnline)
				} else {
					updateLiveUserInfo(gLiveUsersInfo, onlineid, 0, LiveUserStatusBackGround, LiveUserTickBackGround)
				}
			}
		}

		rows.Close()
	}

	sentence = lib.SQLSentence(lib.SQLMapSelectOnlineIDs, 1)
	rows, err = lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&onlineid, &onlineStatus)
			if nil == err {
				if LiveUserStatusOnline == onlineStatus {
					updateLiveUserInfo(gLiveUsersInfo, onlineid, 1, LiveUserStatusOnline, LiveUserTickOnline)
				} else {
					updateLiveUserInfo(gLiveUsersInfo, onlineid, 1, LiveUserStatusBackGround, LiveUserTickBackGround)
				}
			}
		}

		rows.Close()
	}

	for {
		time.Sleep(lib.SleepDurationLiveStatus)

		gLiveUsersInfo.lock.Lock()
		for id, user := range gLiveUsersInfo.users {
			user.livetick = user.livetick - 1
			if 0 == user.livetick {
				if LiveUserStatusBackGround == user.status {
					delete(gLiveUsersInfo.users, id)
					OfflineProc(id, user.gender)
					lib.DelRedisDistrict(id)
				} else {
					user.status = LiveUserStatusBackGround
					user.livetick = LiveUserTickBackGround

					backgroundProc(id, user.gender)
				}
			}

			user.livetime = user.livetime + 1
			if 5 == user.livetime && LiveUserStatusOnline == user.status {
				if true != checkIfUserHavePicture(id, user.gender) {
					//管理员发送第二封信
					msg := "您还没有更新照片哦,上传照片获得更高的推荐机会!"
					timevalue := lib.CurrentTimeUTCInt64()
					RecommendInsertMessageToDB(1, id, CommentMsgTypeTalk, msg, timevalue)
					RecommendPushMessage(1, id, 1, 1, push.PushMsgComment, msg, timevalue)
					push.DoPush()
				}
			}
		}
		gLiveUsersInfo.lock.Unlock()
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

	sentence := lib.SQLSentence(lib.SQLMapSelectVIPRows, 0)
	rows, err := lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&userid, &level, &days, &expiretime)
			if nil == err {
				if curtime >= expiretime {
					detachVipFromUser(userid, 0, level)
				} else {
					UpdateVipUserInfo(userid, 0, level, days, expiretime)
				}
			}
		}

		rows.Close()
	}

	sentence = lib.SQLSentence(lib.SQLMapSelectVIPRows, 1)
	rows, err = lib.SQLQuery(sentence)
	if nil == err {
		for rows.Next() {
			err = rows.Scan(&userid, &level, &days, &expiretime)
			if nil == err {
				if curtime >= expiretime {
					detachVipFromUser(userid, 1, level)
				} else {
					UpdateVipUserInfo(userid, 1, level, days, expiretime)
				}
			}
		}

		rows.Close()
	}

	needpush := false
	for {
		time.Sleep(lib.SleepDurationVIPStatus)
		needpush = false

		gVipUsersInfo.lock.Lock()
		curtime = lib.CurrentTimeUTCInt64()
		for id, user := range gVipUsersInfo.users {
			if curtime >= user.expiretime {
				detachVipFromUser(id, user.gender, user.level)
				needpush = true
			} else {
				difhours := int(time.Duration(user.expiretime-curtime) / time.Hour)
				if 24 == difhours || 120 == difhours || 240 == difhours {
					//发送信息, 提醒到期
					msg := "您的 " + [...]string{"1级会员", "2级会员", "3级会员"}[user.level] + " 将在 " + strconv.Itoa(difhours/24) + " 天后到期."
					RecommendInsertMessageToDB(1, id, CommentMsgTypeTalk, msg, curtime)
					RecommendPushMessage(1, id, 1, 1, push.PushMsgComment, msg, curtime)
					needpush = true
				}
			}
		}
		gVipUsersInfo.lock.Unlock()

		if true == needpush {
			push.DoPush()
		}
	}
}

func detachVipFromUser(id, gender, level int) {
	delete(gVipUsersInfo.users, id)

	sentence := lib.SQLSentence(lib.SQLMapUpdateVIPByID, gender)
	_, err := lib.SQLExec(sentence, 0, 0, 0, id)
	if nil == err {
		//发送信息, VIP权限已经取消
		msg := "非常抱歉通知您, 您的 " + []string{"初始会员", "1级会员", "2级会员", "3级会员"}[level] + " 已经到期!"
		timevalue := lib.CurrentTimeUTCInt64()
		RecommendInsertMessageToDB(1, id, CommentMsgTypeTalk, msg, timevalue)
		RecommendPushMessage(1, id, 1, 1, push.PushMsgComment, msg, timevalue)

		lib.DelRedisUserInfo(id)
	}
}

// BuyVip 购买服务后更新后台数据
func BuyVip(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	levelstr := c.Query("level")
	daysstr := c.Query("days")
	if levelstr == "" || daysstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	level, _ := strconv.Atoi(levelstr)
	days, _ := strconv.Atoi(daysstr)

	/* check if ther user already buy VIP */
	var oldlevel int
	var olddays int
	var expiretime int64

	sentence := lib.SQLSentence(lib.SQLMapSelectVipLevelByID, gender)
	lib.SQLQueryRow(sentence, id).Scan(&oldlevel, &olddays, &expiretime)
	if 0 != oldlevel {
		if oldlevel > level {
			level = oldlevel
		}

		days = days + olddays
	}

	//秒为单位
	if 0 == expiretime {
		expiretime = lib.CurrentTimeUTCInt64()
	}

	expiretime += int64(days) * int64(time.Hour/time.Second) * 24
	sentence = lib.SQLSentence(lib.SQLMapUpdateVIPByID, gender)
	_, err := lib.SQLExec(sentence, level, days, expiretime, id)
	if nil != err {
		c.String(http.StatusNotFound, err.Error())
	}

	gCountBuyVIP = gCountBuyVIP + 1
	updateUserActive()

	//更新到线程
	go UpdateVipUserInfo(id, gender, level, days, expiretime)

	//发送信息, VIP已经开通
	expireUTC := lib.Int64ToUTCTime(expiretime)
	msg := "您的 " + []string{"初始会员", "写信会员", "钻石会员", "至尊会员"}[level] + " 已经开通啦! 会员到期时间：" +
		fmt.Sprintf("%d年%d月%d日", expireUTC.Year(), expireUTC.Month(), expireUTC.Day())
	timevalue := lib.CurrentTimeUTCInt64()
	RecommendInsertMessageToDB(1, id, CommentMsgTypeTalk, msg, timevalue)
	RecommendPushMessage(1, id, 1, 1, push.PushMsgComment, msg, timevalue)
	push.DoPush()

	lib.DelRedisUserInfo(id)
	code, info := GetUserInfo(id, gender)
	c.JSON(code, info)
}

// VipPrice get the vip price information
func VipPrice(c *gin.Context) {
	c.JSON(http.StatusOK, gVipLevels)
}

// GetBuyVIPCount .
func GetBuyVIPCount() int {
	return gCountBuyVIP
}

// SubUserCount .
func SubUserCount(gender int) {
	if 0 == gender {
		if gCountGirls > 0 {
			gCountGirls = gCountGirls - 1
		}
	} else {
		if gCountGuys > 0 {
			gCountGuys = gCountGuys - 1
		}
	}
}

// GetAppConfig .
func GetAppConfig(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var appconfig AppConfig
	var code int
	code, appconfig.Person = GetUserInfo(id, gender)
	appconfig.Person.SendGiftList = GetUserSendGiftList(id)
	appconfig.StartupView.ImageURL = "http://7xjwto.com1.z0.glb.clouddn.com/images/startup/c44eb332f28cfe1b8d067d7da68ffc1e.png"
	appconfig.StartupView.Duration = 4
	appconfig.StartupView.LinkEnable = false
	appconfig.VersionInfo.VersionStr = ""

	c.JSON(code, appconfig)
}
