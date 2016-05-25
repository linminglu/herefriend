package cms

import (
	"database/sql"
	"net/http"
	"strconv"

	log "github.com/cihub/seelog"

	"herefriend/lib"
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
		log.Errorf("SQLQueryRow Error: %s %v\n", sentence, err)
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
			log.Errorf("SQLQueryRow Error: %s %v\n", selectSentence, err)
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
			log.Errorf("SQLQueryRow Error: %s %v\n", selectSentence, err)
			return 404, ""
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
		lib.SQLExec(updateSentence, value+giftvalue, toid)
	}
	return 200, ""
}

func RefreshGiftConsume(r *http.Request) {
	sentence := "select fromid, fromgender, toid, giftid, giftnum from giftconsume"
	rows, err := lib.SQLQuery(sentence)
	if nil != err {
		return
	}
	defer rows.Close()

	var id, gender, toid, giftid, giftnum int
	var tmpid int
	var giftname string
	var price int
	var validnum int
	var value int
	var consume int

	giftsentence := lib.SQLSentence(lib.SQLMAP_Select_GiftById)
	consumeSentence := lib.SQLSentence(lib.SQLMAP_Select_GoldBeansById)
	receiveSentence := lib.SQLSentence(lib.SQLMAP_Select_ReceiveValueById)

	for rows.Next() {
		err = rows.Scan(&id, &gender, &toid, &giftid, &giftnum)
		if nil != err {
			continue
		}

		err = lib.SQLQueryRow(giftsentence, giftid).Scan(&tmpid, &giftname, &price, &validnum)
		if nil != err || giftid != tmpid {
			log.Errorf("SQLQueryRow Error: %s %v\n", giftsentence, err)
			return
		}

		giftvalue := price * giftnum

		// consume the gold beans
		err = lib.SQLQueryRow(consumeSentence, id).Scan(&value, &consume)
		if nil != err {
			if sql.ErrNoRows == err {
				insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_GoldBeansById)
				lib.SQLExec(insertSentence, id, gender, 0, giftvalue)
			} else {
				log.Errorf("SQLQueryRow Error: %s %v\n", consumeSentence, err)
			}
		} else {
			updateSentence := lib.SQLSentence(lib.SQLMAP_Update_GoldBeansById)
			lib.SQLExec(updateSentence, value, consume+giftvalue, id)
		}

		// updathe the receive value
		err = lib.SQLQueryRow(receiveSentence, toid).Scan(&value)
		if nil != err {
			if sql.ErrNoRows == err {
				insertSentence := lib.SQLSentence(lib.SQLMAP_Insert_ReceiveValueById)
				lib.SQLExec(insertSentence, toid, 1-gender, giftvalue)
			} else {
				log.Errorf("SQLQueryRow Error: %s %v\n", receiveSentence, err)
			}
		} else {
			updateSentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
			lib.SQLExec(updateSentence, value+giftvalue, toid)
		}
	}
}
