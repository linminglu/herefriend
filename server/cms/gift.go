package cms

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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
func PresentGift(c *gin.Context) {
	idstr := c.Query("id")
	genderstr := c.Query("gender")
	toidstr := c.Query("toid")
	giftidstr := c.Query("giftid")
	numstr := c.Query("num")

	id, _ := strconv.Atoi(idstr)
	gender, _ := strconv.Atoi(genderstr)
	toid, _ := strconv.Atoi(toidstr)
	giftid, _ := strconv.Atoi(giftidstr)
	giftnum, _ := strconv.Atoi(numstr)

	if 0 == giftnum {
		c.Status(http.StatusForbidden)
		return
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
		c.Status(http.StatusNotFound)
		return
	}

	giftvalue := price * giftnum

	// present the gifts
	sentence = lib.SQLSentence(lib.SQLMAP_Insert_PresentGift)
	_, err = lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), "")
	if nil != err {
		c.Status(http.StatusNotFound)
		return
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
			c.Status(http.StatusNotFound)
			return
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
			c.Status(http.StatusNotFound)
			return
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMAP_Update_ReceiveValueById)
		lib.SQLExec(updateSentence, value+giftvalue, toid)
	}

	c.Status(http.StatusOK)
	return
}

/*
 |    Function: GetGiftList
 |      Author: Mr.Sancho
 |        Date: 2016-07-02
 | Description:
 |      Return:
 |
*/
func GetGiftList(c *gin.Context) {
	err, infolist := handlers.GetGiftList()
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, infolist)
}

/*
 |    Function: GetGiftVerbose
 |      Author: Mr.Sancho
 |        Date: 2016-07-02
 | Description:
 |      Return:
 |
*/
func GetGiftVerbose(c *gin.Context) {
	idstr := c.Query("id")
	id, _ := strconv.Atoi(idstr)
	if id <= 0 {
		c.Status(http.StatusNotFound)
		return
	}

	var info handlers.GiftInfo
	sentence := lib.SQLSentence(lib.SQLMAP_Select_GiftInfoById)
	err := lib.SQLQueryRow(sentence, id).Scan(&info.Id, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageUrl, &info.Effect,
		&info.Price, &info.OriginPrice, &info.DiscountDescription)
	if nil != err {
		lib.SQLError(sentence, err, id)
		c.Status(http.StatusNotFound)
		return
	}

	info.ImageUrl = lib.GetQiniuGiftImageURL(info.ImageUrl)
	c.JSON(http.StatusOK, info)
}
