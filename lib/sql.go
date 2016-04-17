package lib

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"herefriend/config"
)

const (
	SQLMAP_Select_CheckIdValid = iota
	SQLMAP_Select_CheckIdPasswordValid
	SQLMAP_Select_LastId
	SQLMAP_Select_BlacklistLastId
	SQLMAP_Select_OnlineIds
	SQLMAP_Select_UserType
	SQLMAP_Select_SearchBriefInfoHead
	SQLMAP_Select_BriefInfoByRows
	SQLMAP_Select_PersonInfoById
	SQLMAP_Select_IdByRows
	SQLMAP_Select_Pictures
	SQLMAP_Select_SearchPictures
	SQLMAP_Select_SearchPicturesByFlag
	SQLMAP_Select_CheckHeatbeatValid
	SQLMAP_Select_UserCount
	SQLMAP_Select_UserActive
	SQLMAP_Select_HeartbeatCount
	SQLMAP_Select_HeartbeatProvinceCount
	SQLMAP_Select_HeartbeatInfoByRows
	SQLMAP_Select_Heartbeat_RandRows
	SQLMAP_Select_LoveshowByRows
	SQLMAP_Select_LoveshowWithHideByRows
	SQLMAP_Select_LoveshowPicture
	SQLMAP_Select_LoveshowBless
	SQLMAP_Select_SearchHead
	SQLMAP_Select_SearchCountHead
	SQLMAP_Select_InRecommendUnreadCount
	SQLMAP_Select_TalkRecommend
	SQLMAP_Select_RecommendCount
	SQLMAP_Select_AllRecommendCount
	SQLMAP_Select_DistinctRecommend
	SQLMAP_Select_CheckGreetReommend
	SQLMAP_Select_CheckHeartbeatReommend
	SQLMAP_Select_VisitByRows
	SQLMAP_Select_VisitUnreadCount
	SQLMAP_Select_RandomId
	SQLMAP_Select_RandomProvId
	SQLMAP_Select_HeartbeatRandomProvId
	SQLMAP_Select_RandomProvAgeId
	SQLMAP_Select_LastLoginTime
	SQLMAP_Select_ClientID
	SQLMAP_Select_VIPRows
	SQLMAP_Select_VipLevelByID
	SQLMAP_Select_VGirlProcess
	SQLMAP_Select_CheckVGirlId
	SQLMAP_Select_RandomUncrawlGirlsId
	SQLMAP_Select_AllMsgTemplate
	SQLMAP_Select_UserBlacklist
	SQLMAP_Select_CheckUserBlacklist
	SQLMAP_Insert_Info
	SQLMAP_Select_CountByProv
	SQLMAP_Select_CountByProvAge
	SQLMAP_Insert_Picture
	SQLMAP_Insert_Heartbeat
	SQLMAP_Insert_LoveshowBless
	SQLMAP_Insert_Recomment
	SQLMAP_Insert_Visit
	SQLMAP_Insert_VGirlId
	SQLMAP_Insert_Report
	SQLMAP_Insert_Blacklist
	SQLMAP_Insert_UserBlacklist
	SQLMAP_Update_Info
	SQLMAP_Update_InfoPictureFlag
	SQLMAP_Update_RandomInfo
	SQLMAP_Update_Online
	SQLMAP_Update_Background
	SQLMAP_Update_Offline
	SQLMAP_Update_LoginInfo
	SQLMAP_Update_UserActive
	SQLMAP_Update_LoveshowGirl
	SQLMAP_Update_LoveshowGuy
	SQLMAP_Update_HideLoveshow
	SQLMAP_Update_RecommendRead
	SQLMAP_Update_VisitRead
	SQLMAP_Update_Password
	SQLMAP_Update_VIPById
	SQLMAP_Update_VGirlProcess
	SQLMAP_Update_VGirlId
	SQLMAP_Update_SetPictureFlag
	SQLMAP_Update_SetPictureTag
	SQLMAP_Delete_UserId
	SQLMAP_Delete_Picture
	SQLMAP_Delete_HeadPicture
	SQLMAP_Delete_Heartbeat
	SQLMAP_Delete_Recommend
	SQLMAP_Delete_Visit
	SQLMAP_Delete_UserBlacklist
	SQLMAP_Delete_MultiClientID
	SQLMAP_Delete_UnCrawledGirlsId
)

type sqlmapnode struct {
	format_type string
	sentence    string
}

var g_sqlmap = map[int]sqlmapnode{
	SQLMAP_Select_CheckIdValid:         {"s", "select id from %s where id=?"},
	SQLMAP_Select_CheckIdPasswordValid: {"s", "select id from %s where id=? and password=?"},
	SQLMAP_Select_LastId:               {"s", "select id from %s order by id desc limit 1"},
	SQLMAP_Select_BlacklistLastId:      {"", "select id from blacklist order by id desc limit 1"},
	SQLMAP_Select_OnlineIds:            {"s", "select id, onlineStatus from %s where onlineStatus!=0"},
	SQLMAP_Select_UserType:             {"s", "select id, usertype from %s where id=?"},
	SQLMAP_Select_SearchBriefInfoHead:  {"s", "select id from %s"},
	SQLMAP_Select_BriefInfoByRows:      {"s", "select id from %s order by id desc limit ?,?"},
	SQLMAP_Select_PersonInfoById: {"s", "select id, name, age, gender, onlineStatus, viplevel, vipexpiretime, height, weight, lovetype, bodytype, marriage, province, district, " +
		"native, education, income, incomemin, incomemax, occupation, housing, carstatus, introduction, school, speciality, animal, astrology, lang, " +
		"bloodtype, selfjudge, companytype, companyindustry, nationnality, religion, charactor, hobbies, citylove, naken, " +
		"allow_age, allow_residence, allow_height, allow_marriage, allow_education, allow_housing, allow_income, allow_kidstatus from %s where id=?"},
	SQLMAP_Select_IdByRows:               {"s", "select id from %s where pictureflag=1 order by id desc limit ?,?"},
	SQLMAP_Select_Pictures:               {"s", "select filename, tag from %s_picture where id=?"},
	SQLMAP_Select_SearchPictures:         {"s", "select filename from %s_picture where id=? and tag=?"},
	SQLMAP_Select_SearchPicturesByFlag:   {"s", "select id, filename, tag from %s_picture where flag=0 order by id desc limit ?,?"},
	SQLMAP_Select_CheckHeatbeatValid:     {"", "select id from heartbeat where id=?"},
	SQLMAP_Select_UserCount:              {"s", "select count(*) from %s"},
	SQLMAP_Select_UserActive:             {"", "select regist,buy from useractive where id=1"},
	SQLMAP_Select_HeartbeatCount:         {"d", "select count(*) from heartbeat where gender=%d"},
	SQLMAP_Select_HeartbeatProvinceCount: {"d", "select count(*) from heartbeat where gender=%d and province=?"},
	SQLMAP_Select_HeartbeatInfoByRows:    {"d", "select id from heartbeat where gender=%d and province=? order by id desc limit ?,?"},
	SQLMAP_Select_Heartbeat_RandRows:     {"d", "select id from heartbeat where gender=%d and flag=0 limit 1000"},
	SQLMAP_Select_LoveshowByRows: {"", "select loveshowid, time, blessnum, falldays, girl_id, guy_id, girl_age," +
		"guy_age, girl_name, guy_name, girl_headimg, guy_headimg, girl_district, guy_district, status, title, story " +
		"from loveshow where hide=0 order by time desc limit ?,?"},
	SQLMAP_Select_LoveshowWithHideByRows: {"", "select loveshowid, time, blessnum, falldays, girl_id, guy_id, girl_age," +
		"guy_age, girl_name, guy_name, girl_headimg, guy_headimg, girl_district, guy_district, status, title, story, hide " +
		"from loveshow order by time desc limit ?,?"},
	SQLMAP_Select_LoveshowPicture:        {"", "select filename from loveshowpicture where loveshowid=?"},
	SQLMAP_Select_LoveshowBless:          {"", "select uid, age, time, name, district, education, bless from loveshowbless where loveshowid=? order by time desc"},
	SQLMAP_Select_SearchHead:             {"s", "select id from %s where usertype!=1 and "},
	SQLMAP_Select_SearchCountHead:        {"s", "select count(*) from %s where "},
	SQLMAP_Select_InRecommendUnreadCount: {"", "select count(distinct fromid) from recommend where type=2 and toid=? and readed=0"},
	SQLMAP_Select_TalkRecommend:          {"", "select id, fromid, toid, readed, time, msg from recommend where type=? and id>? and ((fromid=? and toid=?) or (fromid=? and toid=?)) order by id desc limit ?,?"},
	SQLMAP_Select_RecommendCount:         {"", "select count(*) from recommend where type=? and fromid=? and toid=?"},
	SQLMAP_Select_AllRecommendCount:      {"", "select count(*) from recommend where type=1 or type=2"},
	SQLMAP_Select_DistinctRecommend: {"", "select id, fromid, toid, readed, time, msg from recommend where id in (select * from (select id from (" +
		"select id, toid as talker from recommend where fromid=? and type=? and time>? union select id, fromid as talker from recommend where toid=? and type=? and time>? order by id desc) as A" +
		" group by talker) as B) order by id desc limit ?, ?"},
	SQLMAP_Select_CheckGreetReommend:     {"", "select time from recommend where fromid=? and toid=? and type=1 order by time desc limit 1"},
	SQLMAP_Select_CheckHeartbeatReommend: {"", "select time from recommend where fromid=? and toid=? and type=3 order by time desc limit 1"},
	SQLMAP_Select_VisitByRows:            {"", "select id, fromid, readed, time from visit where toid=? and time>? order by id desc limit ?,?"},
	SQLMAP_Select_VisitUnreadCount:       {"", "select count(*) from visit where toid=? and readed=0"},
	SQLMAP_Select_RandomId:               {"s", "select id from %s where usertype!=1 limit ?,1"},
	SQLMAP_Select_RandomProvId:           {"s", "select id from %s where usertype!=1 and province=? limit ?,1"},
	SQLMAP_Select_HeartbeatRandomProvId:  {"d", "select id from heartbeat where gender=%d and province=? limit ?,1"},
	SQLMAP_Select_RandomProvAgeId:        {"s", "select id from %s where usertype!=1 and province=? and age>=? and age<=? limit ?,1"},
	SQLMAP_Select_LastLoginTime:          {"s", "select logintime from %s where id=?"},
	SQLMAP_Select_ClientID:               {"s", "select clientid from %s where id=?"},
	SQLMAP_Select_VIPRows:                {"s", "select id, viplevel, vipdays, vipexpiretime from %s where usertype=1 and viplevel!=0"},
	SQLMAP_Select_VipLevelByID:           {"s", "select viplevel, vipdays from %s where id=?"},
	SQLMAP_Select_CountByProv:            {"s", "select count(*) from %s where usertype!=1 and province=?"},
	SQLMAP_Select_CountByProvAge:         {"s", "select count(*) from %s where usertype!=1 and province=? and age=?"},
	SQLMAP_Select_VGirlProcess:           {"", "select areaindex, page from vgirlprocess where base=0"},
	SQLMAP_Select_CheckVGirlId:           {"", "select id from vgirlsid where id=?"},
	SQLMAP_Select_RandomUncrawlGirlsId:   {"", "select id from girlsid where age>=18 and age<=28 limit ?,1"},
	SQLMAP_Select_AllMsgTemplate:         {"", "select msg from msgtemplate where type=? and gender=?"},
	SQLMAP_Select_UserBlacklist:          {"", "select blackid from userblacklist where fromid=?"},
	SQLMAP_Select_CheckUserBlacklist:     {"", "select blackid from userblacklist where fromid=? and blackid=?"},
	SQLMAP_Insert_Info: {"s", "insert into %s (id, password, name, gender, logintime, age, usertype, clientid, height, weight, " +
		"province, district, citylove, naken) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"},
	SQLMAP_Insert_Picture:   {"s", "insert into %s_picture (id, filename, tag, flag) value (?,?,?,1)"},
	SQLMAP_Insert_Heartbeat: {"", "insert into heartbeat (id, gender, province) values (?,?,?)"},
	SQLMAP_Insert_LoveshowBless: {"", "insert into loveshowbless (loveshowid, uid, age, time, name, district, education, bless) values " +
		"(?, ?, ?, ?, ?, ?, ?, ?)"},
	SQLMAP_Insert_Recomment:     {"", "insert into recommend (fromid, toid, time, type, msg, readed) value (?,?,?,?,?,0)"},
	SQLMAP_Insert_Visit:         {"", "insert into visit (fromid, toid, time, readed) value (?,?,?,0)"},
	SQLMAP_Insert_VGirlId:       {"", "insert into vgirlsid (id, fensi, flag) value (?,?,?)"},
	SQLMAP_Insert_Report:        {"", "insert into report (fromid, reportedid, reason) values (?,?,?)"},
	SQLMAP_Insert_Blacklist:     {"s", "insert into blacklist (select * from %s where id=?)"},
	SQLMAP_Insert_UserBlacklist: {"", "insert into userblacklist (fromid, blackid) value (?,?)"},
	SQLMAP_Update_Info: {"s", "update %s set lovetype=?, bodytype=?, marriage=?, province=?, district=?, native=?, education=?, " +
		"occupation=?, housing=?, carstatus=?, introduction=?, school=?, speciality=?, animal=?, astrology=?, lang=?, " +
		"bloodtype=?, selfjudge=?, companytype=?, companyindustry=?, nationnality=?, religion=?, charactor=?, hobbies=?, " +
		"allow_age=?, allow_residence=?, allow_height=?, allow_marriage=?, allow_education=?, allow_housing=?, allow_income=?, " +
		"allow_kidstatus=? where id=?"},
	SQLMAP_Update_InfoPictureFlag:  {"s", "update %s set pictureflag=1 where id=?"},
	SQLMAP_Update_RandomInfo:       {"s", "update %s set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income where id=?"},
	SQLMAP_Update_Online:           {"s", "update %s set onlineStatus=1, logintime=? where id=?"},
	SQLMAP_Update_Background:       {"s", "update %s set onlineStatus=2 where id=?"},
	SQLMAP_Update_Offline:          {"s", "update %s set onlineStatus=0 where id=?"},
	SQLMAP_Update_LoginInfo:        {"s", "update %s set clientid=?, onlineStatus=1, logintime=? where id=?"},
	SQLMAP_Update_UserActive:       {"", "update useractive set regist=?,buy=? where id=1"},
	SQLMAP_Update_LoveshowGirl:     {"", "update loveshow set girl_id=?, girl_age=?, girl_name=?, girl_headimg=?, girl_district=? where loveshowid=?"},
	SQLMAP_Update_LoveshowGuy:      {"", "update loveshow set guy_id=?, guy_age=?, guy_name=?, guy_headimg=?, guy_district=? where loveshowid=?"},
	SQLMAP_Update_HideLoveshow:     {"", "update loveshow set hide=? where loveshowid=?"},
	SQLMAP_Update_RecommendRead:    {"", "update recommend set readed=1 where fromid=? and toid=? and type=2 and id<=?"},
	SQLMAP_Update_VisitRead:        {"", "update visit set readed=1 where id=?"},
	SQLMAP_Update_Password:         {"s", "update %s set password=? where id=?"},
	SQLMAP_Update_VIPById:          {"s", "update %s set viplevel=?, vipdays=?, vipexpiretime=? where id=?"},
	SQLMAP_Update_VGirlProcess:     {"", "update vgirlprocess set areaindex=?, page=? where base=0"},
	SQLMAP_Update_VGirlId:          {"", "update vgirlsid set flag=1 where id=?"},
	SQLMAP_Update_SetPictureFlag:   {"s", "update %s_picture set flag=1 where id=? and filename=? and tag=?"},
	SQLMAP_Update_SetPictureTag:    {"s", "update %s_picture set tag=? where id=? and filename=?"},
	SQLMAP_Delete_UserId:           {"s", "delete from %s where id=?"},
	SQLMAP_Delete_Picture:          {"s", "delete from %s_picture where id=? and filename=?"},
	SQLMAP_Delete_HeadPicture:      {"s", "delete from %s_picture where id=? and tag=1"},
	SQLMAP_Delete_Heartbeat:        {"", "delete from heartbeat where id=?"},
	SQLMAP_Delete_Recommend:        {"", "delete from recommend where id<=? and ((fromid=? and toid=?) or (fromid=? and toid=?))"},
	SQLMAP_Delete_Visit:            {"", "delete from visit where id=?"},
	SQLMAP_Delete_UserBlacklist:    {"", "delete from userblacklist where fromid=? and blackid=?"},
	SQLMAP_Delete_MultiClientID:    {"s", "update %s set clientid='' where clientid=?"},
	SQLMAP_Delete_UnCrawledGirlsId: {"", "delete from girlsid where id=?"},
}

var g_DBHandle *sql.DB

/*
 *
 *    Function: init
 *      Author: sunchao
 *        Date: 15/6/20
 * Description: init the connection to db
 *
 */
func init() {
	g_DBHandle, _ = sql.Open(config.Conf_Driver, config.Conf_Dns)
	err := g_DBHandle.Ping()
	if nil != err {
		panic(err.Error())
	}
}

/*
 *
 *    Function: Fini
 *      Author: sunchao
 *        Date: 15/6/20
 * Description: release the resources
 *
 */
func CloseSQL() {
	g_DBHandle.Close()
}

/*
 |    Function: SQLSentence
 |      Author: Mr.Sancho
 |        Date: 2016-01-09
 |   Arguments:
 |      Return:
 | Description: 生成一条SQL语句
 |
*/
func SQLSentence(key int, args ...interface{}) string {
	mapnode, ok := g_sqlmap[key]

	if true != ok {
		return ""
	}

	switch mapnode.format_type {
	case "s":
		return fmt.Sprintf(mapnode.sentence, [2]string{"girls", "guys"}[args[0].(int)])
	case "d":
		return fmt.Sprintf(mapnode.sentence, args[0].(int))
	default:
		return mapnode.sentence
	}
}

func SQLExec(query string, args ...interface{}) (sql.Result, error) {
	return g_DBHandle.Exec(query, args...)
}

func SQLQueryRow(query string, args ...interface{}) *sql.Row {
	return g_DBHandle.QueryRow(query, args...)
}

func SQLQuery(query string, args ...interface{}) (*sql.Rows, error) {
	return g_DBHandle.Query(query, args...)
}