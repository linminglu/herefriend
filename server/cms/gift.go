package cms

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"herefriend/lib"
	"herefriend/server/handlers"
)

// PresentGift .
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

	sentence := lib.SQLSentence(lib.SQLMapSelectGiftByID)
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
	sentence = lib.SQLSentence(lib.SQLMapInsertPresentGift)
	_, err = lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), "")
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	var value int
	var consume int

	// consume the gold beans
	selectSentence := lib.SQLSentence(lib.SQLMapSelectGoldBeansByID)
	err = lib.SQLQueryRow(selectSentence, id).Scan(&value, &consume)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMapInsertGoldBeansByID)
			lib.SQLExec(insertSentence, id, gender, 0, giftvalue)
		} else {
			lib.SQLError(selectSentence, err, id)
			c.Status(http.StatusNotFound)
			return
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMapUpdateGoldBeansByID)
		lib.SQLExec(updateSentence, value, consume+giftvalue, id)
	}

	// updathe the receive value
	selectSentence = lib.SQLSentence(lib.SQLMapSelectReceiveValueByID)
	err = lib.SQLQueryRow(selectSentence, toid).Scan(&value)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMapInsertReceiveValueByID)
			lib.SQLExec(insertSentence, toid, 1-gender, giftvalue)
		} else {
			lib.SQLError(selectSentence, err, toid)
			c.Status(http.StatusNotFound)
			return
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMapUpdateReceiveValueByID)
		lib.SQLExec(updateSentence, value+giftvalue, toid)
	}

	c.Status(http.StatusOK)
	return
}

// GetGiftList .
func GetGiftList(c *gin.Context) {
	infolist, err := handlers.GetGiftList()
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, infolist)
}

// GetGiftVerbose .
func GetGiftVerbose(c *gin.Context) {
	idstr := c.Query("id")
	id, _ := strconv.Atoi(idstr)
	if id <= 0 {
		c.Status(http.StatusNotFound)
		return
	}

	var info handlers.GiftInfo
	sentence := lib.SQLSentence(lib.SQLMapSelectGiftInfoByID)
	err := lib.SQLQueryRow(sentence, id).Scan(&info.ID, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageURL, &info.Effect,
		&info.Price, &info.OriginPrice, &info.DiscountDescription)
	if nil != err {
		lib.SQLError(sentence, err, id)
		c.Status(http.StatusNotFound)
		return
	}

	info.ImageURL = lib.GetQiniuGiftImageURL(info.ImageURL)
	c.JSON(http.StatusOK, info)
}
