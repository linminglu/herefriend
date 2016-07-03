package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"herefriend/common"
	"herefriend/lib"
	"herefriend/lib/push"
)

func init() {
	go prepare()
}

func prepare() {
	var id int

	sentence := "select distinct toid from giftconsume"
	rows, err := lib.SQLQuery(sentence)
	if nil != err {
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if nil == err && 0 != id {
			PrepareUserRecvGiftList(id)
		}
	}
}

/*
 |    Function: GoldPrice
 |      Author: Mr.Sancho
 |        Date: 2016-04-17
 |   Arguments:
 |      Return:
 | Description: 获取金币价格列表
 |
*/
func GoldPrice(r *http.Request) (int, string) {
	jsonRlt, _ := json.Marshal(gGoldBeansPrices)
	return 200, string(jsonRlt)
}

/*
 |    Function: BuyBeans
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 购买金币
 |      Return:
 |
*/
func BuyBeans(r *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	v := r.URL.Query()
	beansStr := v.Get("beans")
	if "" == beansStr {
		return 404, ""
	}

	beans, _ := strconv.Atoi(beansStr)
	if 0 == beans {
		return 403, ""
	}

	var value int
	var consume int
	selectSentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
	err := lib.SQLQueryRow(selectSentence, id).Scan(&value, &consume)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_GoldBeansById)
			lib.SQLExec(insertSentence, id, gender, beans, 0)
		} else {
			lib.SQLError(selectSentence, err, id)
			return 404, ""
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
		lib.SQLExec(updateSentence, beans+value, consume, id)
	}

	lib.DelRedisGoldBeans(id)
	lib.DelRedisUserInfo(id)
	code, info := GetUserInfo(id, gender)
	info.SendGiftList = GetUserSendGiftList(id)
	jsonRlt, _ := json.Marshal(info)

	return code, string(jsonRlt)
}

/*
 |    Function: getGiftList
 |      Author: Mr.Sancho
 |        Date: 2016-06-09
 | Description:
 |      Return:
 |
*/
func GetGiftList() (error, []GiftInfo) {
	infolist := make([]GiftInfo, 0)

	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftInfo)
	rows, err := lib.SQLQuery(sentence)
	if nil != err {
		return err, infolist
	}
	defer rows.Close()

	var info GiftInfo
	for rows.Next() {
		err = rows.Scan(&info.Id, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageUrl, &info.Effect,
			&info.Price, &info.OriginPrice, &info.DiscountDescription)
		if nil == err {
			info.ImageUrl = lib.GetQiniuGiftImageURL(info.ImageUrl)
			infolist = append(infolist, info)
		}
	}

	return nil, infolist
}

/*
 |    Function: GiftList
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 获取礼物列表
 |      Return:
 |
*/
func GiftList(r *http.Request) (int, string) {
	err, infolist := GetGiftList()
	if nil != err {
		return 404, ""
	}

	jsonRlt, _ := json.Marshal(infolist)
	return 200, string(jsonRlt)
}

/*
 |    Function: PresentGift
 |      Author: Mr.Sancho
 |        Date: 2016-04-24
 | Description: 送出礼物
 |      Return:
 |
*/
func PresentGift(r *http.Request) (int, string) {
	exist, id, gender := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	v := r.URL.Query()
	toidstr := v.Get("toid")
	giftidstr := v.Get("giftid")
	numstr := v.Get("num")
	msg := v.Get("message")
	if "" == toidstr || "" == giftidstr || "" == numstr {
		return 404, ""
	}

	toid, _ := strconv.Atoi(toidstr)
	giftid, _ := strconv.Atoi(giftidstr)
	giftnum, _ := strconv.Atoi(numstr)
	exist, togender, usertype := GetGenderUsertypeById(toid)
	if false == exist || 0 == giftnum {
		return 403, ""
	}

	// check the gifts
	var tmpid int
	var giftname string
	var price int
	var validnum int

	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftById)
	err := lib.SQLQueryRow(sentence, giftid).Scan(&tmpid, &giftname, &price, &validnum)
	if nil != err || giftid != tmpid {
		if nil != err {
			lib.SQLError(sentence, err, giftid)
		}
		return 404, ""
	}

	if validnum > giftnum {
		validnum = validnum - giftnum
	} else {
		validnum = 0
	}

	giftvalue := price * giftnum

	// check if the beans is enough
	var beansValue int
	var consumevalue int

	selectSentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
	err = lib.SQLQueryRow(selectSentence, id).Scan(&beansValue, &consumevalue)
	if nil != err || beansValue < giftvalue {
		if nil != err {
			lib.SQLError(sentence, err, id)
		}
		return 403, "没有足够的金币购买此数量的礼物"
	}

	// present the gifts
	sentence = lib.SQLSentence(lib.SQLMAP_Insert_PresentGift)
	_, err = lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), msg)
	if nil != err {
		return 404, ""
	}

	// consume the gold beans
	updateSentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
	lib.SQLExec(updateSentence, beansValue-giftvalue, consumevalue+giftvalue, id)

	// consume the gift
	sentence = lib.SQLSentence(lib.SQLMAP_Update_ConsumeGift)
	_, err = lib.SQLExec(sentence, validnum, giftid)
	if nil != err {
		return 404, ""
	}

	// updathe the receive value
	var value int
	selectSentence = lib.SQLSentence(lib.SQLMAP_Select_ReceiveValueById)
	err = lib.SQLQueryRow(selectSentence, toid).Scan(&value)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_ReceiveValueById)
			lib.SQLExec(insertSentence, toid, togender, giftvalue)
		} else {
			lib.SQLError(sentence, err, toid)
			return 404, ""
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
		lib.SQLExec(updateSentence, value+giftvalue, toid)
	}

	lib.DelRedisGiftSendList(id)
	lib.DelRedisGiftRecvList(toid)
	lib.DelRedisGoldBeans(id)
	lib.DelRedisUserInfo(id)
	lib.DelRedisUserInfo(toid)

	go func() {
		if common.USERTYPE_RB == usertype {
			return
		}

		_, userinfo := GetUserInfo(id, togender)
		newgiftmsg := fmt.Sprintf("你收到[ %s ]赠送的礼物: %d 个 [ %s ]。你的魅力值又增加了 %d。", userinfo.Name, giftnum, giftname, giftvalue*10)
		clientid := GetClientIdByUserId(toid)

		// 普通通知消息
		push.Add(0, clientid, push.PUSHMSG_TYPE_RECOMMEND, 0, "收到新礼物啦", newgiftmsg)

		// 透传消息
		recvGiftMsg := PushMsgRecvGift{
			SenderId:    id,
			GiftId:      giftid,
			GiftNum:     giftnum,
			GiftName:    giftname,
			ShowMessage: newgiftmsg,
		}

		jsonRlt, _ := json.Marshal(recvGiftMsg)
		notifymsg := PushMessageInfo{Type: push.PUSH_NOTIFYMSG_RECVGIFT, Value: string(jsonRlt)}
		jsonRlt, _ = json.Marshal(notifymsg)
		push.Add(0, clientid, push.PUSHMSG_TYPE_NOTIFYMSG, push.PUSH_NOTIFYMSG_RECVGIFT, "", string(jsonRlt))

		push.DoPush()
	}()

	// refresh personal info
	var presentInfo presentGiftInfo
	_, presentInfo.WhoRecvGift = GetUserInfo(toid, 1-gender)
	_, presentInfo.UserInfo = GetUserInfo(id, gender)
	presentInfo.UserInfo.SendGiftList = GetUserSendGiftList(id)

	jsonRlt, _ := json.Marshal(presentInfo)
	return 200, string(jsonRlt)
}

/*
 |    Function: getRecvGiftListById
 |      Author: Mr.Sancho
 |        Date: 2016-06-09
 | Description:
 |      Return:
 |
*/
func getRecvGiftListById(id int, page, count int) (error, []GiftListVerbose) {
	giftlist := make([]GiftListVerbose, 0)

	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftRecvVerbose)
	rows, err := lib.SQLQuery(sentence, id, (page-1)*count, count)
	if nil != err {
		return err, giftlist
	}
	defer rows.Close()

	var timetmp int64
	var info GiftListVerbose
	var userid int

	for rows.Next() {
		err = rows.Scan(&info.Id, &userid, &info.GiftId, &info.GiftNum, &timetmp, &info.Message)
		if nil == err {
			_, info.Person = GetUserInfoById(userid)
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)
			giftlist = append(giftlist, info)
		}
	}

	return nil, giftlist
}

/*
 |    Function: RecvListVerbose
 |      Author: Mr.Sancho
 |        Date: 2016-04-30
 | Description: 获取收到的礼物详情
 |      Return:
 |
*/
func RecvListVerbose(r *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	v := r.URL.Query()
	queryidstr := v.Get("queryid")
	if "" == queryidstr {
		return 404, ""
	}

	queryid, _ := strconv.Atoi(queryidstr)
	page, count := lib.Get_pageid_count_fromreq(r)
	err, giftlist := getRecvGiftListById(queryid, page, count)
	if nil != err {
		return 404, ""
	}

	jsonRlt, _ := json.Marshal(giftlist)
	return 200, string(jsonRlt)
}

/*
 |    Function: getSendGiftListById
 |      Author: Mr.Sancho
 |        Date: 2016-06-09
 | Description:
 |      Return:
 |
*/
func getSendGiftListById(id int, page, count int) (error, []GiftListVerbose) {
	giftlist := make([]GiftListVerbose, 0)

	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftSendVerbose)
	rows, err := lib.SQLQuery(sentence, id, (page-1)*count, count)
	if nil != err {
		return err, giftlist
	}
	defer rows.Close()

	var timetmp int64
	var info GiftListVerbose
	var userid int
	for rows.Next() {
		err = rows.Scan(&info.Id, &userid, &info.GiftId, &info.GiftNum, &timetmp, &info.Message)
		if nil == err {
			_, info.Person = GetUserInfoById(userid)
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)
			giftlist = append(giftlist, info)
		}
	}

	return nil, giftlist
}

/*
 |    Function: SendListVerbose
 |      Author: Mr.Sancho
 |        Date: 2016-04-30
 | Description: 获取收到的礼物详情
 |      Return:
 |
*/
func SendListVerbose(r *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	v := r.URL.Query()
	queryidstr := v.Get("queryid")
	if "" == queryidstr {
		return 404, ""
	}

	queryid, _ := strconv.Atoi(queryidstr)
	page, count := lib.Get_pageid_count_fromreq(r)

	err, giftlist := getSendGiftListById(queryid, page, count)
	if nil != err {
		return 404, ""
	}

	jsonRlt, _ := json.Marshal(giftlist)
	return 200, string(jsonRlt)
}

/*
 |    Function: CharmTopList
 |      Author: Mr.Sancho
 |        Date: 2016-05-08
 | Description:
 |      Return:
 |
*/
func CharmTopList(r *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	v := r.URL.Query()
	genderstr := v.Get("gender")
	gender, _ := strconv.Atoi(genderstr)
	page, count := lib.Get_pageid_count_fromreq(r)

	sentence := lib.SQLSentence(lib.SQLMAP_Select_CharmToplist)
	rows, err := lib.SQLQuery(sentence, gender, (page-1)*count, count)
	if nil != err {
		return 404, ""
	}
	defer rows.Close()

	charmlist := make([]userCharmInfo, 0)
	for rows.Next() {
		var tempid int
		var code int
		var info userCharmInfo

		err = rows.Scan(&tempid, &info.GiftValue)
		if nil == err {
			code, info.Person = GetUserInfo(tempid, gender)
			if 200 == code {
				info.GiftValue = info.GiftValue * 10
				charmlist = append(charmlist, info)
			}
		}
	}

	jsonRlt, _ := json.Marshal(charmlist)
	return 200, string(jsonRlt)
}

/*
 |    Function: WealthTopList
 |      Author: Mr.Sancho
 |        Date: 2016-05-29
 | Description:
 |      Return:
 |
*/
func WealthTopList(r *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	page, count := lib.Get_pageid_count_fromreq(r)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_WealthToplist)
	rows, err := lib.SQLQuery(sentence, (page-1)*count, count)
	if nil != err {
		return 404, ""
	}
	defer rows.Close()

	wealthlist := make([]userWealthInfo, 0)
	var tempid int
	var code int
	var info userWealthInfo

	for rows.Next() {
		err = rows.Scan(&tempid, &info.ConsumedBeans)
		if nil == err {
			code, info.Person = GetUserInfoById(tempid)
			if 200 == code {
				info.ConsumedBeans = info.ConsumedBeans * 10
				wealthlist = append(wealthlist, info)
			}
		}
	}

	jsonRlt, _ := json.Marshal(wealthlist)
	return 200, string(jsonRlt)
}

/*
 |    Function: DeleteUserWealthAndGiftInfo
 |      Author: Mr.Sancho
 |        Date: 2016-06-09
 | Description:
 |      Return:
 |
*/
func DeleteUserWealthAndGiftInfo(id int) {
	_, giftinfos := GetGiftList()

	//删除收到的礼物信息
	for {
		selectsentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
		updatesentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
		deletesentence := lib.SQLSentence(lib.SQLMAP_Delete_GiftConsumeInfo)

		_, recvlist := getRecvGiftListById(id, 1, 1000)
		if 0 == len(recvlist) {
			break
		}

		for _, info := range recvlist {
			var beans, consumed int
			var value int
			var sender int

			for _, gift := range giftinfos {
				if gift.Id == info.GiftId {
					value = gift.Price * info.GiftNum
					sender = info.Person.Id
					break
				}
			}

			lib.SQLQueryRow(selectsentence, sender).Scan(&beans, consumed)
			if consumed >= value {
				consumed = consumed - value
			}
			lib.SQLExec(updatesentence, sender, beans, consumed)
			lib.SQLExec(deletesentence, info.Id)

			lib.DelRedisGiftSendList(sender)
			lib.DelRedisUserInfo(sender)
		}
	}

	//删除送出的礼物信息
	for {
		selectsentence := lib.SQLSentence(lib.SQLMAP_Select_ReceiveValueById)
		updatesentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
		deletesentence := lib.SQLSentence(lib.SQLMAP_Delete_GiftConsumeInfo)

		_, sendlist := getSendGiftListById(id, 1, 1000)
		if 0 == len(sendlist) {
			break
		}

		for _, info := range sendlist {
			var value int
			var receiver int
			var receivevalue int

			for _, gift := range giftinfos {
				if gift.Id == info.GiftId {
					value = gift.Price * info.GiftNum
					receiver = info.Person.Id
					break
				}
			}

			lib.SQLQueryRow(selectsentence, receiver).Scan(&receivevalue)
			if receivevalue >= value {
				receivevalue = receivevalue - value
			}
			lib.SQLExec(updatesentence, receivevalue, receiver)
			lib.SQLExec(deletesentence, info.Id)

			lib.DelRedisGiftRecvList(receiver)
			lib.DelRedisUserInfo(receiver)
		}
	}

	delwealthSentence := lib.SQLSentence(lib.SQLMAP_Delete_Wealth)
	lib.SQLExec(delwealthSentence, id)
	lib.DelRedisGiftRecvList(id)
	lib.DelRedisGiftSendList(id)
}

func DeleteGiftInfoByUserIdAndGiftId(id, giftid int) {
	_, giftinfos := GetGiftList()

	//删除送出的指定礼物信息
	for {
		selectsentence := lib.SQLSentence(lib.SQLMAP_Select_ReceiveValueById)
		updatesentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
		deletesentence := lib.SQLSentence(lib.SQLMAP_Delete_GiftConsumeInfo)

		_, sendlist := getSendGiftListById(id, 1, 100000)
		if 0 == len(sendlist) {
			break
		}

		for _, info := range sendlist {
			if giftid != info.GiftId {
				continue
			}

			var value int
			var receiver int
			var receivevalue int

			for _, gift := range giftinfos {
				if gift.Id == giftid {
					value = gift.Price * info.GiftNum
					receiver = info.Person.Id
					break
				}
			}

			lib.SQLQueryRow(selectsentence, receiver).Scan(&receivevalue)
			if receivevalue >= value {
				receivevalue = receivevalue - value
			}
			lib.SQLExec(updatesentence, receivevalue, receiver)
			lib.SQLExec(deletesentence, info.Id)

			lib.DelRedisGiftRecvList(receiver)
			lib.DelRedisUserInfo(receiver)
		}
	}

	lib.DelRedisGiftSendList(id)
}
