package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	log "github.com/cihub/seelog"

	"herefriend/common"
	"herefriend/lib"
	"herefriend/lib/push"
)

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
	selectSentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
	err := lib.SQLQueryRow(selectSentence, id).Scan(&value)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_GoldBeansById)
			lib.SQLExec(insertSentence, id, beans)
		} else {
			return 404, ""
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
		lib.SQLExec(updateSentence, beans+value, id)
	}

	code, info := GetUserInfo(id, gender)
	info.SendGiftList = GetUserSendGiftList(id)
	jsonRlt, _ := json.Marshal(info)

	return code, string(jsonRlt)
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
	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftInfo)
	rows, err := lib.SQLQuery(sentence)
	if nil != err {
		return 404, ""
	}

	defer rows.Close()

	infolist := make([]giftInfo, 0)
	var info giftInfo
	for rows.Next() {
		err = rows.Scan(&info.Id, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageUrl, &info.Effect,
			&info.Price, &info.OriginPrice, &info.DiscountDescription)
		if nil == err {
			info.ImageUrl = lib.GetQiniuGiftImageURL(info.ImageUrl)
			infolist = append(infolist, info)
		}
	}

	go log.Tracef("获取礼物列表")
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

	if 0 == giftnum {
		return 403, ""
	}

	// present the gifts
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_PresentGift)
	_, err := lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), msg)
	if nil != err {
		return 404, ""
	}

	// consume the gifts
	var tmpid int
	var giftname string
	var price int
	var validnum int

	sentence = lib.SQLSentence(lib.SQLMAP_Select_GiftById)
	err = lib.SQLQueryRow(sentence, giftid).Scan(&tmpid, &giftname, &price, &validnum)
	if nil != err || giftid != tmpid {
		return 404, ""
	}

	if validnum > giftnum {
		validnum = validnum - giftnum
	} else {
		validnum = 0
	}

	// check if the beans is enough
	var beansValue int
	selectSentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
	err = lib.SQLQueryRow(selectSentence, id).Scan(&beansValue)
	if nil != err || beansValue < price*giftnum {
		return 403, "没有足够的金币购买此数量的礼物"
	}

	updateSentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
	lib.SQLExec(updateSentence, beansValue-price*giftnum, id)

	sentence = lib.SQLSentence(lib.SQLMAP_Update_ConsumeGift)
	_, err = lib.SQLExec(sentence, validnum, giftid)
	if nil != err {
		return 404, ""
	}

	go func() {
		exist, _, usertype := GetGenderUsertypeById(toid)
		if false == exist || common.USERTYPE_RB == usertype {
			return
		}

		_, userinfo := GetUserInfo(id, gender)
		newgiftmsg := fmt.Sprintf("你收到[ %s ]赠送的礼物: %d 个 [ %s ]。你的魅力值又增加了。", userinfo.Name, giftnum, giftname)
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
 |    Function: RecvListVerbose
 |      Author: Mr.Sancho
 |        Date: 2016-04-30
 | Description: 获取收到的礼物详情
 |      Return:
 |
*/
func RecvListVerbose(r *http.Request) (int, string) {
	exist, id, _ := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	page, count := lib.Get_pageid_count_fromreq(r)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftRecvVerbose)
	rows, err := lib.SQLQuery(sentence, id, (page-1)*count, count)
	if nil != err {
		return 404, ""
	}

	defer rows.Close()

	giftlist := make([]giftListVerbose, 0)
	var timetmp int64
	var info giftListVerbose
	var userid int

	for rows.Next() {
		err = rows.Scan(&userid, &info.GiftId, &info.GiftNum, &timetmp, &info.Message)
		if nil == err {
			_, info.Person = GetUserInfoById(userid)
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)
			giftlist = append(giftlist, info)
		}
	}

	jsonRlt, _ := json.Marshal(giftlist)
	return 200, string(jsonRlt)
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
	exist, id, _ := getIdGenderByRequest(r)
	if true != exist {
		return 404, ""
	}

	page, count := lib.Get_pageid_count_fromreq(r)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftSendVerbose)
	rows, err := lib.SQLQuery(sentence, id, (page-1)*count, count)
	if nil != err {
		return 404, ""
	}

	defer rows.Close()

	giftlist := make([]giftListVerbose, 0)
	var timetmp int64
	var info giftListVerbose
	var userid int

	for rows.Next() {
		err = rows.Scan(&userid, &info.GiftId, &info.GiftNum, &timetmp, &info.Message)
		if nil == err {
			_, info.Person = GetUserInfoById(userid)
			info.TimeUTC = lib.Int64_To_UTCTime(timetmp)
			giftlist = append(giftlist, info)
		}
	}

	jsonRlt, _ := json.Marshal(giftlist)
	return 200, string(jsonRlt)
}

type recvValueItem struct {
	id    int
	value int
}

type recvValueSorter []recvValueItem

func (v recvValueSorter) Len() int {
	return len(v)
}

func (v recvValueSorter) Less(i, j int) bool {
	return v[i].value > v[j].value
}

func (v recvValueSorter) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
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

	giftprice := make(map[int]int)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftInfo)
	rows, err := lib.SQLQuery(sentence)
	if nil != err {
		return 404, ""
	} else {
		var info giftInfo
		for rows.Next() {
			err = rows.Scan(&info.Id, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageUrl, &info.Effect,
				&info.Price, &info.OriginPrice, &info.DiscountDescription)
			if nil == err {
				giftprice[info.Id] = info.Price
			}
		}

		rows.Close()
	}

	recvlist := make([]giftRecvListInfo, 0)

	sentence = lib.SQLSentence(lib.SQLMAP_Select_GiftRecvListByGender)
	rows, err = lib.SQLQuery(sentence, 1-gender)
	if nil != err {
		return 404, ""
	} else {
		var info giftRecvListInfo

		for rows.Next() {
			err = rows.Scan(&info.toid, &info.giftid, &info.giftnum)
			if nil == err {
				recvlist = append(recvlist, info)
			}
		}

		rows.Close()
	}

	charmlist := make([]userCharmInfo, 0)

	if 0 != len(recvlist) {
		recvvalue := make(map[int]int)

		for _, info := range recvlist {
			recvvalue[info.toid] = recvvalue[info.toid] + giftprice[info.giftid]*info.giftnum
		}

		valuelist := make(recvValueSorter, 0, len(recvvalue))
		for id, value := range recvvalue {
			valuelist = append(valuelist, recvValueItem{id: id, value: value})
		}
		sort.Sort(valuelist)

		for _, v := range valuelist {
			var charminfo userCharmInfo
			_, charminfo.Person = GetUserInfo(v.id, gender)
			charminfo.GiftValue = v.value

			charmlist = append(charmlist, charminfo)
		}
	}

	jsonRlt, _ := json.Marshal(charmlist)
	return 200, string(jsonRlt)
}
