package lib

import (
	"database/sql"
	"fmt"
	"runtime/debug"

	log "github.com/cihub/seelog"
	// mysql just use like this
	_ "github.com/go-sql-driver/mysql"

	"herefriend/config"
)

const (
	// SQLMapSelectCheckIsValidID .
	SQLMapSelectCheckIsValidID = iota
	// SQLMapSelectCheckIsValidPasswd .
	SQLMapSelectCheckIsValidPasswd
	// SQLMapSelectLastID .
	SQLMapSelectLastID
	// SQLMapSelectBlacklistLastID .
	SQLMapSelectBlacklistLastID
	// SQLMapSelectOnlineIDs .
	SQLMapSelectOnlineIDs
	// SQLMapSelectUserType .
	SQLMapSelectUserType
	// SQLMapSelectPersonInfo .
	SQLMapSelectPersonInfo
	// SQLMapSelectSearchPictures .
	SQLMapSelectSearchPictures
	// SQLMapSelectSearchPicturesByFlag .
	SQLMapSelectSearchPicturesByFlag
	// SQLMapSelectUserCount .
	SQLMapSelectUserCount
	// SQLMapSelectUserActive .
	SQLMapSelectUserActive
	// SQLMapSelectHeartbeatCount .
	SQLMapSelectHeartbeatCount
	// SQLMapSelectHeartbeatProvinceCount .
	SQLMapSelectHeartbeatProvinceCount
	// SQLMapSelectHeartbeatInfoByRows .
	SQLMapSelectHeartbeatInfoByRows
	// SQLMapSelectHeartbeatRandRows .
	SQLMapSelectHeartbeatRandRows
	// SQLMapSelectSearch .
	SQLMapSelectSearch
	// SQLMapSelectCount .
	SQLMapSelectCount
	// SQLMapSelectUnreadMessageCount .
	SQLMapSelectUnreadMessageCount
	// SQLMapSelectMessageHistory .
	SQLMapSelectMessageHistory
	// SQLMapSelectHaveSameReply .
	SQLMapSelectHaveSameReply
	// SQLMapSelectRecommendCount .
	SQLMapSelectRecommendCount
	// SQLMapSelectAllRecommendCount .
	SQLMapSelectAllRecommendCount
	// SQLMapSelectDistinctRecommend .
	SQLMapSelectDistinctRecommend
	// SQLMapSelectCheckCommentDailyLock .
	SQLMapSelectCheckCommentDailyLock
	// SQLMapSelectVisitByRows .
	SQLMapSelectVisitByRows
	// SQLMapSelectVisitUnreadCount .
	SQLMapSelectVisitUnreadCount
	// SQLMapSelectRandomID .
	SQLMapSelectRandomID
	// SQLMapSelectRandomProvID .
	SQLMapSelectRandomProvID
	// SQLMapSelectHeartbeatRandomProvID .
	SQLMapSelectHeartbeatRandomProvID
	// SQLMapSelectRandomProvAgeID .
	SQLMapSelectRandomProvAgeID
	// SQLMapSelectLastLoginTime .
	SQLMapSelectLastLoginTime
	// SQLMapSelectLastEvaluationTime .
	SQLMapSelectLastEvaluationTime
	// SQLMapSelectClientID .
	SQLMapSelectClientID
	// SQLMapSelectVIPRows .
	SQLMapSelectVIPRows
	// SQLMapSelectVipLevelByID .
	SQLMapSelectVipLevelByID
	// SQLMapSelectVGirlProcess .
	SQLMapSelectVGirlProcess
	// SQLMapSelectZQProcess .
	SQLMapSelectZQProcess
	// SQLMapSelectCheckVGirlID .
	SQLMapSelectCheckVGirlID
	// SQLMapSelectCheckZQUserID .
	SQLMapSelectCheckZQUserID
	// SQLMapSelectRandomUncrawlGirlsID .
	SQLMapSelectRandomUncrawlGirlsID
	// SQLMapSelectAllMsgTemplate .
	SQLMapSelectAllMsgTemplate
	// SQLMapSelectUserBlacklist .
	SQLMapSelectUserBlacklist
	// SQLMapSelectCheckUserBlacklist .
	SQLMapSelectCheckUserBlacklist
	// SQLMapSelectCountByProv .
	SQLMapSelectCountByProv
	// SQLMapSelectCountByProvAge .
	SQLMapSelectCountByProvAge
	// SQLMapSelectGiftInfo .
	SQLMapSelectGiftInfo
	// SQLMapSelectGiftInfoByID .
	SQLMapSelectGiftInfoByID
	// SQLMapSelectGiftByID .
	SQLMapSelectGiftByID
	// SQLMapSelectGiftRecvSum .
	SQLMapSelectGiftRecvSum
	// SQLMapSelectGiftSendSum .
	SQLMapSelectGiftSendSum
	// SQLMapSelectGiftRecvVerbose .
	SQLMapSelectGiftRecvVerbose
	// SQLMapSelectGiftSendVerbose .
	SQLMapSelectGiftSendVerbose
	// SQLMapSelectGiftRecvListByGender .
	SQLMapSelectGiftRecvListByGender
	// SQLMapSelectGoldBeansByID .
	SQLMapSelectGoldBeansByID
	// SQLMapSelectReceiveValueByID .
	SQLMapSelectReceiveValueByID
	// SQLMapSelectCharmToplist .
	SQLMapSelectCharmToplist
	// SQLMapSelectWealthToplist .
	SQLMapSelectWealthToplist
	// SQLMapInsertInfo .
	SQLMapInsertInfo
	// SQLMapInsertPicture .
	SQLMapInsertPicture
	// SQLMapInsertHeartbeat .
	SQLMapInsertHeartbeat
	// SQLMapInsertRecomment .
	SQLMapInsertRecomment
	// SQLMapInsertVisit .
	SQLMapInsertVisit
	// SQLMapInsertVGirlID .
	SQLMapInsertVGirlID
	// SQLMapInsertZQID .
	SQLMapInsertZQID
	// SQLMapInsertReport .
	SQLMapInsertReport
	// SQLMapInsertBlacklist .
	SQLMapInsertBlacklist
	// SQLMapInsertUserBlacklist .
	SQLMapInsertUserBlacklist
	// SQLMapInsertPresentGift .
	SQLMapInsertPresentGift
	// SQLMapInsertGoldBeansByID .
	SQLMapInsertGoldBeansByID
	// SQLMapInsertReceiveValueByID .
	SQLMapInsertReceiveValueByID
	// SQLMapUpdateInfo .
	SQLMapUpdateInfo
	// SQLMapUpdateInfoPictureFlag .
	SQLMapUpdateInfoPictureFlag
	// SQLMapUpdateRandomInfo .
	SQLMapUpdateRandomInfo
	// SQLMapUpdateOnline .
	SQLMapUpdateOnline
	// SQLMapUpdateBackground .
	SQLMapUpdateBackground
	// SQLMapUpdateOffline .
	SQLMapUpdateOffline
	// SQLMapUpdateLoginInfo .
	SQLMapUpdateLoginInfo
	// SQLMapUpdateUserActive .
	SQLMapUpdateUserActive
	// SQLMapUpdateRecommendRead .
	SQLMapUpdateRecommendRead
	// SQLMapUpdateVisitRead .
	SQLMapUpdateVisitRead
	// SQLMapUpdatePassword .
	SQLMapUpdatePassword
	// SQLMapUpdateVIPByID .
	SQLMapUpdateVIPByID
	// SQLMapUpdateVGirlProcess .
	SQLMapUpdateVGirlProcess
	// SQLMapUpdateVGirlID .
	SQLMapUpdateVGirlID
	// SQLMapUpdateZQProcess .
	SQLMapUpdateZQProcess
	// SQLMapUpdateSetPictureFlag .
	SQLMapUpdateSetPictureFlag
	// SQLMapUpdateSetPictureTag .
	SQLMapUpdateSetPictureTag
	// SQLMapUpdateConsumeGift .
	SQLMapUpdateConsumeGift
	// SQLMapUpdateGoldBeansByID .
	SQLMapUpdateGoldBeansByID
	// SQLMapUpdateReceiveValueByID .
	SQLMapUpdateReceiveValueByID
	// SQLMapUpdateEvaluationTime .
	SQLMapUpdateEvaluationTime
	// SQLMapDeleteWealth .
	SQLMapDeleteWealth
	// SQLMapDeleteGiftConsumeInfo .
	SQLMapDeleteGiftConsumeInfo
	// SQLMapDeleteUserID .
	SQLMapDeleteUserID
	// SQLMapDeletePicture .
	SQLMapDeletePicture
	// SQLMapDeleteHeadPicture .
	SQLMapDeleteHeadPicture
	// SQLMapDeleteHeartbeat .
	SQLMapDeleteHeartbeat
	// SQLMapDeleteRecommend .
	SQLMapDeleteRecommend
	// SQLMapDeleteVisit .
	SQLMapDeleteVisit
	// SQLMapDeleteRecommendByUserID .
	SQLMapDeleteRecommendByUserID
	// SQLMapDeleteVisitByUserID .
	SQLMapDeleteVisitByUserID
	// SQLMapDeleteUserBlacklist .
	SQLMapDeleteUserBlacklist
	// SQLMapDeleteMultiClientID .
	SQLMapDeleteMultiClientID
	// SQLMapDeleteUnCrawledGirlsID .
	SQLMapDeleteUnCrawledGirlsID
	// SQLMapCMSSelectBriefInfo .
	SQLMapCMSSelectBriefInfo
	// SQLMapCMSSelectSetVipAppVersion .
	SQLMapCMSSelectSetVipAppVersion
	// SQLMapCMSSelectPictures .
	SQLMapCMSSelectPictures
	// SQLMapCMSSelectBriefInfoByRows .
	SQLMapCMSSelectBriefInfoByRows
	// SQLMapCMSSelectCheckHeatbeatValid .
	SQLMapCMSSelectCheckHeatbeatValid
)

type sqlmapnode struct {
	formatType string
	sentence   string
}

var gSQLMap = map[int]sqlmapnode{
	SQLMapSelectCheckIsValidID:     {"s", "select id from %s where id=?"},
	SQLMapSelectCheckIsValidPasswd: {"s", "select id from %s where id=? and password=?"},
	SQLMapSelectLastID:             {"s", "select id from %s order by id desc limit 1"},
	SQLMapSelectBlacklistLastID:    {"", "select id from blacklist order by id desc limit 1"},
	SQLMapSelectOnlineIDs:          {"s", "select id, onlineStatus from %s where onlineStatus!=0"},
	SQLMapSelectUserType:           {"s", "select id, usertype from %s where id=?"},
	SQLMapSelectPersonInfo: {"s", "select id, name, age, gender, onlineStatus, viplevel, vipexpiretime, " +
		"height, weight, lovetype, bodytype, marriage, province, district, native, education, income, incomemin, incomemax, " +
		"occupation, housing, carstatus, introduction, school, speciality, animal, astrology, lang, bloodtype, selfjudge, " +
		"companytype, companyindustry, nationnality, religion, charactor, hobbies, citylove, naken, allow_age, allow_residence, " +
		"allow_height, allow_marriage, allow_education, allow_housing, allow_income, allow_kidstatus from %s where id=?"},
	SQLMapSelectSearchPictures:         {"s", "select filename from %s_picture where id=? and tag=? order by numid"},
	SQLMapSelectSearchPicturesByFlag:   {"s", "select id, filename, tag from %s_picture where flag=0 and id>=? order by id desc limit ?,?"},
	SQLMapSelectUserCount:              {"s", "select count(*) from %s"},
	SQLMapSelectUserActive:             {"", "select regist,buy from useractive where id=1"},
	SQLMapSelectHeartbeatCount:         {"d", "select count(*) from heartbeat where gender=%d"},
	SQLMapSelectHeartbeatProvinceCount: {"d", "select count(*) from heartbeat where gender=%d and province=?"},
	SQLMapSelectHeartbeatInfoByRows:    {"d", "select id from heartbeat where gender=%d and province=? order by id desc limit ?,?"},
	SQLMapSelectHeartbeatRandRows:      {"d", "select id from heartbeat where gender=%d and flag=0 limit 1000"},
	SQLMapSelectSearch:                 {"s", "select id from %s where usertype!=1 and "},
	SQLMapSelectCount:                  {"s", "select count(*) from %s where "},
	SQLMapSelectUnreadMessageCount:     {"", "select count(distinct fromid) from recommend where type=2 and toid=? and readed=0 and time>?"},
	SQLMapSelectMessageHistory: {"", "select id, fromid, toid, readed, time, msg from recommend where type=? and id>? and " +
		"((fromid=? and toid=?) or (fromid=? and toid=?)) order by id desc limit ?,?"},
	SQLMapSelectHaveSameReply:     {"", "select count(*) from recommend where type=2 and toid=? and msg=?"},
	SQLMapSelectRecommendCount:    {"", "select count(*) from recommend where type=? and fromid=? and toid=?"},
	SQLMapSelectAllRecommendCount: {"", "select count(*) from recommend where type=1 or type=2"},
	SQLMapSelectDistinctRecommend: {"", "select id, fromid, toid, readed, time, msg from recommend where id in (select * from (select id from (" +
		"select id, toid as talker from recommend where fromid=? and type=? and time>? union select id, fromid as talker from recommend where toid=? " +
		"and type=? and time>? order by id desc) as A group by talker) as B) order by id desc limit ?,?"},
	SQLMapSelectCheckCommentDailyLock: {"", "select time from recommend where fromid=? and toid=? and type=? order by time desc limit 1"},
	SQLMapSelectVisitByRows:           {"", "select id, fromid, readed, time from visit where toid=? and time>? order by id desc limit ?,?"},
	SQLMapSelectVisitUnreadCount:      {"", "select count(*) from visit where toid=? and readed=0 and time>?"},
	SQLMapSelectRandomID:              {"s", "select id from %s where usertype!=1 limit ?,1"},
	SQLMapSelectRandomProvID:          {"s", "select id from %s where usertype!=1 and province=? limit ?,1"},
	SQLMapSelectHeartbeatRandomProvID: {"d", "select id from heartbeat where gender=%d and province=? limit ?,1"},
	SQLMapSelectRandomProvAgeID:       {"s", "select id from %s where usertype!=1 and province=? and age>=? and age<=? limit ?,1"},
	SQLMapSelectLastLoginTime:         {"s", "select logintime from %s where id=?"},
	SQLMapSelectLastEvaluationTime:    {"s", "select evaluationtime from %s where id=?"},
	SQLMapSelectClientID:              {"s", "select clientid from %s where id=?"},
	SQLMapSelectVIPRows:               {"s", "select id, viplevel, vipdays, vipexpiretime from %s where usertype=1 and viplevel!=0"},
	SQLMapSelectVipLevelByID:          {"s", "select viplevel, vipdays, vipexpiretime from %s where id=?"},
	SQLMapSelectVGirlProcess:          {"", "select areaindex, page from vgirlprocess where base=0"},
	SQLMapSelectZQProcess:             {"", "select areaindex, page from zhenqingprocess where base=?"},
	SQLMapSelectCheckVGirlID:          {"", "select id from vgirlsid where id=?"},
	SQLMapSelectCheckZQUserID:         {"", "select id from zhenqingids where id=?"},
	SQLMapSelectRandomUncrawlGirlsID:  {"", "select id from girlsid where age>=18 and age<=28 limit ?,1"},
	SQLMapSelectAllMsgTemplate:        {"", "select msg from msgtemplate where type=? and gender=?"},
	SQLMapSelectUserBlacklist:         {"", "select blackid from userblacklist where fromid=?"},
	SQLMapSelectCheckUserBlacklist:    {"", "select blackid from userblacklist where fromid=? and blackid=?"},
	SQLMapSelectCountByProv:           {"s", "select count(*) from %s where usertype!=1 and province=?"},
	SQLMapSelectCountByProvAge:        {"s", "select count(*) from %s where usertype!=1 and province=? and age=?"},
	SQLMapSelectGiftInfo:              {"", "select id,type,name,description,validnum,imageurl,effect,price,origin_price,discount_desciption from gift"},
	SQLMapSelectGiftInfoByID:          {"", "select id,type,name,description,validnum,imageurl,effect,price,origin_price,discount_desciption from gift where id=?"},
	SQLMapSelectGiftByID:              {"", "select id, name, price, validnum from gift where id=?"},
	SQLMapSelectGiftRecvSum:           {"", "select giftid, giftnum from giftconsume where toid=?"},
	SQLMapSelectGiftSendSum:           {"", "select giftid, giftnum from giftconsume where fromid=?"},
	SQLMapSelectGiftRecvVerbose:       {"", "select id, fromid, giftid, giftnum, time, message from giftconsume where toid=? order by time desc limit ?,?"},
	SQLMapSelectGiftSendVerbose:       {"", "select id, toid, giftid, giftnum, time, message from giftconsume where fromid=? order by time desc limit ?,?"},
	SQLMapSelectGiftRecvListByGender:  {"", "select toid, giftid, giftnum from giftconsume where fromgender=? order by toid"},
	SQLMapSelectGoldBeansByID:         {"", "select beans,consumed from wealth where id=?"},
	SQLMapSelectReceiveValueByID:      {"", "select receive from wealth where id=?"},
	SQLMapSelectCharmToplist:          {"", "select toid, giftid, giftnum from giftconsume where fromgender=? and time>=? and time<?"},
	SQLMapSelectWealthToplist:         {"", "select fromid, giftid, giftnum from giftconsume where time>=? and time<?"},
	SQLMapInsertInfo: {"s", "insert into %s (id, password, name, gender, logintime, age, usertype, clientid, height, weight, " +
		"province, district, citylove, naken) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"},
	SQLMapInsertPicture:          {"s", "insert into %s_picture (id, filename, tag, flag) value (?,?,?,1)"},
	SQLMapInsertHeartbeat:        {"", "insert into heartbeat (id, gender, province) values (?,?,?)"},
	SQLMapInsertRecomment:        {"", "insert into recommend (fromid, toid, time, type, msg, readed) value (?,?,?,?,?,0)"},
	SQLMapInsertVisit:            {"", "insert into visit (fromid, toid, time, readed) value (?,?,?,0)"},
	SQLMapInsertVGirlID:          {"", "insert into vgirlsid (id, fensi, flag) value (?,?,?)"},
	SQLMapInsertZQID:             {"", "insert into zhenqingids (id) value (?)"},
	SQLMapInsertReport:           {"", "insert into report (fromid, reportedid, reason) values (?,?,?)"},
	SQLMapInsertBlacklist:        {"s", "insert into blacklist (select * from %s where id=?)"},
	SQLMapInsertUserBlacklist:    {"", "insert into userblacklist (fromid, blackid) value (?,?)"},
	SQLMapInsertPresentGift:      {"", "insert into giftconsume (fromid, fromgender, toid, giftid, giftnum, time, message) values (?,?,?,?,?,?,?)"},
	SQLMapInsertGoldBeansByID:    {"", "insert into wealth (id, gender, beans, consumed, receive) values (?,?,?,?,0)"},
	SQLMapInsertReceiveValueByID: {"", "insert into wealth (id, gender, beans, receive) values (?,?,0,?)"},
	SQLMapUpdateInfo: {"s", "update %s set lovetype=?, bodytype=?, marriage=?, province=?, district=?, native=?, education=?, " +
		"occupation=?, housing=?, carstatus=?, introduction=?, school=?, speciality=?, animal=?, astrology=?, lang=?, " +
		"bloodtype=?, selfjudge=?, companytype=?, companyindustry=?, nationnality=?, religion=?, charactor=?, hobbies=?, " +
		"allow_age=?, allow_residence=?, allow_height=?, allow_marriage=?, allow_education=?, allow_housing=?, allow_income=?, " +
		"allow_kidstatus=? where id=?"},
	SQLMapUpdateInfoPictureFlag: {"s", "update %s set pictureflag=1 where id=?"},
	SQLMapUpdateRandomInfo: {"s", "update %s set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?," +
		"housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income where id=?"},
	SQLMapUpdateOnline:                {"s", "update %s set onlineStatus=1, logintime=? where id=?"},
	SQLMapUpdateBackground:            {"s", "update %s set onlineStatus=2 where id=?"},
	SQLMapUpdateOffline:               {"s", "update %s set onlineStatus=0 where id=?"},
	SQLMapUpdateLoginInfo:             {"s", "update %s set clientid=?, onlineStatus=1, logintime=? where id=?"},
	SQLMapUpdateUserActive:            {"", "update useractive set regist=?,buy=? where id=1"},
	SQLMapUpdateRecommendRead:         {"", "update recommend set readed=1 where fromid=? and toid=? and type=2 and id<=?"},
	SQLMapUpdateVisitRead:             {"", "update visit set readed=1 where id=?"},
	SQLMapUpdatePassword:              {"s", "update %s set password=? where id=?"},
	SQLMapUpdateVIPByID:               {"s", "update %s set viplevel=?, vipdays=?, vipexpiretime=? where id=?"},
	SQLMapUpdateVGirlProcess:          {"", "update vgirlprocess set areaindex=?, page=? where base=0"},
	SQLMapUpdateVGirlID:               {"", "update vgirlsid set flag=1 where id=?"},
	SQLMapUpdateZQProcess:             {"", "update zhenqingprocess set areaindex=?, page=? where base=?"},
	SQLMapUpdateSetPictureFlag:        {"s", "update %s_picture set flag=1 where id=? and filename=? and tag=?"},
	SQLMapUpdateSetPictureTag:         {"s", "update %s_picture set tag=? where id=? and filename=?"},
	SQLMapUpdateConsumeGift:           {"", "update gift set validnum=? where id=?"},
	SQLMapUpdateGoldBeansByID:         {"", "update wealth set beans=?,consumed=? where id=?"},
	SQLMapUpdateReceiveValueByID:      {"", "update wealth set receive=? where id=?"},
	SQLMapUpdateEvaluationTime:        {"s", "update %s set evaluationtime=? where id=?"},
	SQLMapDeleteWealth:                {"", "delete from wealth where id=?"},
	SQLMapDeleteGiftConsumeInfo:       {"", "delete from giftconsume where id=?"},
	SQLMapDeleteUserID:                {"s", "delete from %s where id=?"},
	SQLMapDeletePicture:               {"s", "delete from %s_picture where id=? and filename=?"},
	SQLMapDeleteHeadPicture:           {"s", "delete from %s_picture where id=? and tag=1"},
	SQLMapDeleteHeartbeat:             {"", "delete from heartbeat where id=?"},
	SQLMapDeleteRecommend:             {"", "delete from recommend where id<=? and ((fromid=? and toid=?) or (fromid=? and toid=?))"},
	SQLMapDeleteVisit:                 {"", "delete from visit where id=?"},
	SQLMapDeleteRecommendByUserID:     {"", "delete from recommend where fromid=? or toid=?"},
	SQLMapDeleteVisitByUserID:         {"", "delete from visit where fromid=? or toid=?"},
	SQLMapDeleteUserBlacklist:         {"", "delete from userblacklist where fromid=? and blackid=?"},
	SQLMapDeleteMultiClientID:         {"s", "update %s set clientid='' where clientid=?"},
	SQLMapDeleteUnCrawledGirlsID:      {"", "delete from girlsid where id=?"},
	SQLMapCMSSelectBriefInfo:          {"s", "select id from %s"},
	SQLMapCMSSelectSetVipAppVersion:   {"s", "select setvip_appversion from %s where id=?"},
	SQLMapCMSSelectPictures:           {"s", "select filename, tag from %s_picture where id=?"},
	SQLMapCMSSelectBriefInfoByRows:    {"s", "select id from %s order by id desc limit ?,?"},
	SQLMapCMSSelectCheckHeatbeatValid: {"", "select id from heartbeat where id=?"},
}

var gDBHandle *sql.DB

func init() {
	gDBHandle, _ = sql.Open(config.ConfDBDriver, config.ConfDBDns)
	err := gDBHandle.Ping()
	if nil != err {
		panic(err.Error())
	}
}

// CloseSQL .
func CloseSQL() {
	gDBHandle.Close()
}

// SQLSentence 生成一条SQL语句
func SQLSentence(key int, args ...interface{}) string {
	mapnode, ok := gSQLMap[key]

	if true != ok {
		return ""
	}

	switch mapnode.formatType {
	case "s":
		return fmt.Sprintf(mapnode.sentence, [2]string{"girls", "guys"}[args[0].(int)])
	case "d":
		return fmt.Sprintf(mapnode.sentence, args[0].(int))
	default:
		return mapnode.sentence
	}
}

// SQLExec .
func SQLExec(query string, args ...interface{}) (sql.Result, error) {
	result, err := gDBHandle.Exec(query, args...)
	if nil != err {
		SQLError(query, err, args...)
	}

	return result, err
}

// SQLQueryRow .
func SQLQueryRow(query string, args ...interface{}) *sql.Row {
	return gDBHandle.QueryRow(query, args...)
}

// SQLQuery .
func SQLQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := gDBHandle.Query(query, args...)
	if nil != err {
		SQLError(query, err, args...)
	}

	return rows, err
}

// SQLError .
func SQLError(query string, err error, args ...interface{}) {
	if nil == args {
		log.Errorf("SQL:[%s] error:[%v]", query, err)
	} else {
		log.Errorf("SQL:[%s] args:%v error:[%v]", query, args, err)
	}

	log.Error(string(debug.Stack()))
}
