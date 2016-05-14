package cms

import (
	"net/http"
	"strconv"

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

	sentence := lib.SQLSentence(lib.SQLMAP_Insert_PresentGift)
	_, err := lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), "")
	if nil != err {
		return 404, ""
	}

	return 200, ""
}
