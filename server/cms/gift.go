package cms

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"herefriend/lib"
	"herefriend/server/handlers"
)

/*
 |    Function: PresentGift
 |      Author: Mr.Sancho
 |        Date: 2016-05-14
 | Description:
 |      Return:
 |
*/
func PresentGift(r *http.Request) (int, string) {
	v := r.URL.Query()
	idstr := v.Get("id")
	genderstr := v.Get("gender")
	toidstr := v.Get("toid")
	giftidstr := v.Get("giftid")
	numstr := v.Get("num")

	id, _ := strconv.Atoi(idstr)
	gender, _ := strconv.Atoi(genderstr)
	toid, _ := strconv.Atoi(toidstr)
	giftid, _ := strconv.Atoi(giftidstr)
	giftnum, _ := strconv.Atoi(numstr)

	if 0 == giftnum {
		return 403, ""
	}

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

	giftvalue := price * giftnum

	// present the gifts
	sentence = lib.SQLSentence(lib.SQLMAP_Insert_PresentGift)
	_, err = lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), "")
	if nil != err {
		return 404, ""
	}

	var value int
	var consume int

	// consume the gold beans
	selectSentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
	err = lib.SQLQueryRow(selectSentence, id).Scan(&value, &consume)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_GoldBeansById)
			lib.SQLExec(insertSentence, id, gender, 0, giftvalue)
		} else {
			lib.SQLError(selectSentence, err, id)
			return 404, ""
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
		lib.SQLExec(updateSentence, value, consume+giftvalue, id)
	}

	// updathe the receive value
	selectSentence = lib.SQLSentence(lib.SQLMAP_Select_ReceiveValueById)
	err = lib.SQLQueryRow(selectSentence, toid).Scan(&value)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_ReceiveValueById)
			lib.SQLExec(insertSentence, toid, 1-gender, giftvalue)
		} else {
			lib.SQLError(selectSentence, err, toid)
			return 404, ""
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
		lib.SQLExec(updateSentence, value+giftvalue, toid)
	}
	return 200, ""
}

/*
 |    Function: GetGiftList
 |      Author: Mr.Sancho
 |        Date: 2016-07-02
 | Description:
 |      Return:
 |
*/
func GetGiftList(r *http.Request) (int, string) {
	err, infolist := handlers.GetGiftList()
	if nil != err {
		return 404, ""
	}

	jsonRlt, _ := json.Marshal(infolist)
	return 200, string(jsonRlt)
}

/*
 |    Function: GetGiftVerbose
 |      Author: Mr.Sancho
 |        Date: 2016-07-02
 | Description:
 |      Return:
 |
*/
func GetGiftVerbose(r *http.Request) (int, string) {
	v := r.URL.Query()
	idstr := v.Get("id")
	id, _ := strconv.Atoi(idstr)
	if id <= 0 {
		return 404, ""
	}

	var info handlers.GiftInfo
	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftInfoById)
	err := lib.SQLQueryRow(sentence, id).Scan(&info.Id, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageUrl, &info.Effect,
		&info.Price, &info.OriginPrice, &info.DiscountDescription)
	if nil != err {
		lib.SQLError(sentence, err, id)
		return 404, ""
	}

	info.ImageUrl = lib.GetQiniuGiftImageURL(info.ImageUrl)
	jsonRlt, _ := json.Marshal(info)
	return 200, string(jsonRlt)
}
