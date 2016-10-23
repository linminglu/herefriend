package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"herefriend/common"
	"herefriend/config"
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

// GoldPrice 获取金币价格列表
func GoldPrice(c *gin.Context) {
	c.JSON(http.StatusOK, gGoldBeansPrices)
}

// BuyBeans 购买金币
func BuyBeans(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	beansStr := c.Query("beans")
	if "" == beansStr {
		c.Status(http.StatusNotFound)
		return
	}

	beans, _ := strconv.Atoi(beansStr)
	if 0 == beans {
		c.Status(http.StatusForbidden)
		return
	}

	var value int
	var consume int
	selectSentence := lib.SQLSentence(lib.SQLMapSelectGoldBeansByID)
	err := lib.SQLQueryRow(selectSentence, id).Scan(&value, &consume)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMapInsertGoldBeansByID)
			lib.SQLExec(insertSentence, id, gender, beans, 0)
		} else {
			lib.SQLError(selectSentence, err, id)
			c.Status(http.StatusNotFound)
			return
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMapUpdateGoldBeansByID)
		lib.SQLExec(updateSentence, beans+value, consume, id)
	}

	lib.DelRedisGoldBeans(id)
	lib.DelRedisUserInfo(id)
	code, info := GetUserInfo(id, gender)
	info.SendGiftList = GetUserSendGiftList(id)

	c.JSON(code, info)
}

// GetGiftList .
func GetGiftList() ([]GiftInfo, error) {
	var infolist []GiftInfo

	sentence := lib.SQLSentence(lib.SQLMapSelectGiftInfo)
	rows, err := lib.SQLQuery(sentence)
	if nil != err {
		return infolist, err
	}
	defer rows.Close()

	var info GiftInfo
	for rows.Next() {
		err = rows.Scan(&info.ID, &info.Type, &info.Name, &info.Description, &info.ValidNum, &info.ImageURL, &info.Effect,
			&info.Price, &info.OriginPrice, &info.DiscountDescription)
		if nil == err {
			info.ImageURL = lib.GetQiniuGiftImageURL(info.ImageURL)
			infolist = append(infolist, info)
		}
	}

	return infolist, nil
}

// GiftList 获取礼物列表
func GiftList(c *gin.Context) {
	infolist, err := GetGiftList()
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, infolist)
}

// PresentGift 送出礼物
func PresentGift(c *gin.Context) {
	exist, id, gender := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	toidstr := c.Query("toid")
	giftidstr := c.Query("giftid")
	numstr := c.Query("num")
	msg := c.Query("message")
	if toidstr == "" || giftidstr == "" || numstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	toid, _ := strconv.Atoi(toidstr)
	giftid, _ := strconv.Atoi(giftidstr)
	giftnum, _ := strconv.Atoi(numstr)
	exist, togender, usertype := GetGenderUsertypeByID(toid)
	if !exist || giftnum == 0 {
		c.Status(http.StatusForbidden)
		return
	}

	// check the gifts
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

	if validnum > giftnum {
		validnum = validnum - giftnum
	} else {
		validnum = 0
	}

	giftvalue := price * giftnum

	// check if the beans is enough
	var beansValue int
	var consumevalue int

	selectSentence := lib.SQLSentence(lib.SQLMapSelectGoldBeansByID)
	err = lib.SQLQueryRow(selectSentence, id).Scan(&beansValue, &consumevalue)
	if nil != err || beansValue < giftvalue {
		if nil != err {
			lib.SQLError(sentence, err, id)
		}

		c.String(http.StatusForbidden, "没有足够的金币购买此数量的礼物")
		return
	}

	// present the gifts
	sentence = lib.SQLSentence(lib.SQLMapInsertPresentGift)
	_, err = lib.SQLExec(sentence, id, gender, toid, giftid, giftnum, lib.CurrentTimeUTCInt64(), msg)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	// consume the gold beans
	updateSentence := lib.SQLSentence(lib.SQLMapUpdateGoldBeansByID)
	lib.SQLExec(updateSentence, beansValue-giftvalue, consumevalue+giftvalue, id)

	// consume the gift
	sentence = lib.SQLSentence(lib.SQLMapUpdateConsumeGift)
	_, err = lib.SQLExec(sentence, validnum, giftid)
	if nil != err {
		c.Status(http.StatusNotFound)
		return
	}

	// updathe the receive value
	var value int
	selectSentence = lib.SQLSentence(lib.SQLMapSelectReceiveValueByID)
	err = lib.SQLQueryRow(selectSentence, toid).Scan(&value)
	if nil != err {
		if sql.ErrNoRows == err {
			insertSentence := lib.SQLSentence(lib.SQLMapInsertReceiveValueByID)
			lib.SQLExec(insertSentence, toid, togender, giftvalue)
		} else {
			lib.SQLError(sentence, err, toid)
			c.Status(http.StatusNotFound)
			return
		}
	} else {
		updateSentence := lib.SQLSentence(lib.SQLMapUpdateReceiveValueByID)
		lib.SQLExec(updateSentence, value+giftvalue, toid)
	}

	lib.DelRedisGiftSendList(id)
	lib.DelRedisGiftRecvList(toid)
	lib.DelRedisGoldBeans(id)
	lib.DelRedisUserInfo(id)
	lib.DelRedisUserInfo(toid)

	go func() {
		if common.UserTypeRobot == usertype {
			return
		}

		_, userinfo := GetUserInfo(id, togender)
		newgiftmsg := fmt.Sprintf("你收到[ %s ]赠送的礼物: %d 个 [ %s ]。你的魅力值又增加了 %d。", userinfo.Name, giftnum, giftname, giftvalue*10)
		clientid := GetClientIDByUserID(toid)

		// 普通通知消息
		push.Add(0, clientid, push.PushMsgComment, 0, "收到新礼物啦", newgiftmsg)

		// 透传消息
		recvGiftMsg := PushMsgRecvGift{
			SenderID:    id,
			GiftID:      giftid,
			GiftNum:     giftnum,
			GiftName:    giftname,
			ShowMessage: newgiftmsg,
		}

		jsonRlt, _ := json.Marshal(recvGiftMsg)
		notifymsg := PushMessageInfo{Type: push.NotifyMsgRecvGift, Value: string(jsonRlt)}
		jsonRlt, _ = json.Marshal(notifymsg)
		push.Add(0, clientid, push.PushMsgNotify, push.NotifyMsgRecvGift, "", string(jsonRlt))

		push.DoPush()
	}()

	// refresh personal info
	var presentInfo presentGiftInfo
	_, presentInfo.WhoRecvGift = GetUserInfo(toid, 1-gender)
	_, presentInfo.UserInfo = GetUserInfo(id, gender)
	presentInfo.UserInfo.SendGiftList = GetUserSendGiftList(id)

	c.JSON(http.StatusOK, presentInfo)
}

func getRecvGiftListByID(id int, page, count int) ([]GiftListVerbose, error) {
	var giftlist []GiftListVerbose

	sentence := lib.SQLSentence(lib.SQLMapSelectGiftRecvVerbose)
	rows, err := lib.SQLQuery(sentence, id, (page-1)*count, count)
	if nil != err {
		return giftlist, err
	}
	defer rows.Close()

	var timetmp int64
	var info GiftListVerbose
	var userid int

	for rows.Next() {
		err = rows.Scan(&info.ID, &userid, &info.GiftID, &info.GiftNum, &timetmp, &info.Message)
		if nil == err {
			_, info.Person = GetUserInfoByID(userid)
			info.TimeUTC = lib.Int64ToUTCTime(timetmp)
			giftlist = append(giftlist, info)
		}
	}

	return giftlist, nil
}

// RecvListVerbose 获取收到的礼物详情
func RecvListVerbose(c *gin.Context) {
	exist, _, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	queryidstr := c.Query("queryid")
	if queryidstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	queryid, _ := strconv.Atoi(queryidstr)
	page, count := lib.GetPageidCount(c)
	giftlist, err := getRecvGiftListByID(queryid, page, count)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, giftlist)
}

// getSendGiftListByID .
func getSendGiftListByID(id int, page, count int) ([]GiftListVerbose, error) {
	var giftlist []GiftListVerbose

	sentence := lib.SQLSentence(lib.SQLMapSelectGiftSendVerbose)
	rows, err := lib.SQLQuery(sentence, id, (page-1)*count, count)
	if nil != err {
		return giftlist, err
	}
	defer rows.Close()

	var timetmp int64
	var info GiftListVerbose
	var userid int
	for rows.Next() {
		err = rows.Scan(&info.ID, &userid, &info.GiftID, &info.GiftNum, &timetmp, &info.Message)
		if nil == err {
			_, info.Person = GetUserInfoByID(userid)
			info.TimeUTC = lib.Int64ToUTCTime(timetmp)
			giftlist = append(giftlist, info)
		}
	}

	return giftlist, nil
}

// SendListVerbose 获取收到的礼物详情
func SendListVerbose(c *gin.Context) {
	exist, _, _ := getIDGenderByRequest(c)
	if !exist {
		c.Status(http.StatusNotFound)
		return
	}

	queryidstr := c.Query("queryid")
	if queryidstr == "" {
		c.Status(http.StatusNotFound)
		return
	}

	queryid, _ := strconv.Atoi(queryidstr)
	page, count := lib.GetPageidCount(c)

	giftlist, err := getSendGiftListByID(queryid, page, count)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, giftlist)
}

// CharmTopList .
func CharmTopList(c *gin.Context) {
	exist, _, _ := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	genderstr := c.Query("gender")
	gender, _ := strconv.Atoi(genderstr)

	var charmlist *[]common.UserCharmInfo
	var exists bool

	nowtime := time.Now().UTC()
	year, month, day := nowtime.Date()
	charmlist, exists = lib.GetRedisCharmToplist(gender, year, month, day)
	if true != exists {
		infolist, err := GetGiftList()
		if nil != err {
			c.Status(http.StatusNotFound)
			return
		}

		giftmap := make(map[int]*GiftInfo)
		for _, g := range infolist {
			giftmap[g.ID] = &g
		}

		until := lib.TimeToUTCInt64(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
		from := until - config.ConfToplistDuration
		sentence := lib.SQLSentence(lib.SQLMapSelectCharmToplist)
		rows, err := lib.SQLQuery(sentence, 1-gender, from, until)
		if nil != err {
			c.Status(http.StatusNotFound)
			return
		}
		defer rows.Close()

		charminfomap := make(map[int]common.UserCharmInfo)
		for rows.Next() {
			var toid int
			var giftid int
			var giftnum int
			err = rows.Scan(&toid, &giftid, &giftnum)
			if nil == err {
				value := giftmap[giftid].Price * giftnum * 10

				charminfo, ok := charminfomap[toid]
				if true == ok {
					charminfo.GiftValue += value
					charminfomap[toid] = charminfo
				} else {
					var info common.UserCharmInfo

					var code int
					code, info.Person = GetUserInfo(toid, gender)
					if 200 == code {
						info.GiftValue = value
						charminfomap[toid] = info
					}
				}
			}
		}

		var newlist []common.UserCharmInfo
		for _, info := range charminfomap {
			newlist = append(newlist, info)
		}

		sortcharmlist := common.UserCharmInfoList(newlist)
		sort.Sort(sortcharmlist)
		charmlist = &newlist
		lib.SetRedisCharmToplist(gender, year, month, day, &newlist)
	}

	maxlen := len(*charmlist)
	page, count := lib.GetPageidCount(c)
	start := (page - 1) * count

	if start >= maxlen {
		c.String(http.StatusOK, "[]")
		return
	}

	end := start + count
	if end >= maxlen {
		end = maxlen
	}

	resultlist := (*charmlist)[start:end]
	c.JSON(http.StatusOK, resultlist)
}

// WealthTopList .
func WealthTopList(c *gin.Context) {
	exist, _, _ := getIDGenderByRequest(c)
	if true != exist {
		c.Status(http.StatusNotFound)
		return
	}

	var wealthlist *[]common.UserWealthInfo
	var exists bool

	nowtime := time.Now().UTC()
	year, month, day := nowtime.Date()
	wealthlist, exists = lib.GetRedisWealthToplist(year, month, day)
	if true != exists {
		infolist, err := GetGiftList()
		if nil != err {
			c.Status(http.StatusNotFound)
			return
		}

		giftmap := make(map[int]*GiftInfo)
		for _, g := range infolist {
			giftmap[g.ID] = &g
		}

		until := lib.TimeToUTCInt64(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
		from := until - config.ConfToplistDuration
		sentence := lib.SQLSentence(lib.SQLMapSelectWealthToplist)
		rows, err := lib.SQLQuery(sentence, from, until)
		if nil != err {
			c.Status(http.StatusNotFound)
			return
		}
		defer rows.Close()

		wealthinfomap := make(map[int]common.UserWealthInfo)
		for rows.Next() {
			var fromid int
			var giftid int
			var giftnum int
			err = rows.Scan(&fromid, &giftid, &giftnum)
			if nil == err {
				value := giftmap[giftid].Price * giftnum * 10

				wealthinfo, ok := wealthinfomap[fromid]
				if true == ok {
					wealthinfo.ConsumedBeans += value
					wealthinfomap[fromid] = wealthinfo
				} else {
					var info common.UserWealthInfo

					var code int
					code, info.Person = GetUserInfoByID(fromid)
					if 200 == code {
						info.ConsumedBeans = value
						wealthinfomap[fromid] = info
					}
				}
			}
		}

		var newlist []common.UserWealthInfo
		for _, info := range wealthinfomap {
			newlist = append(newlist, info)
		}

		sortwealthlist := common.UserWealthInfoList(newlist)
		sort.Sort(sortwealthlist)
		wealthlist = &newlist
		lib.SetRedisWealthToplist(year, month, day, &newlist)
	}

	maxlen := len(*wealthlist)
	page, count := lib.GetPageidCount(c)
	start := (page - 1) * count

	if start >= maxlen {
		c.String(http.StatusOK, "[]")
		return
	}

	end := start + count
	if end >= maxlen {
		end = maxlen
	}

	resultlist := (*wealthlist)[start:end]
	c.JSON(http.StatusOK, resultlist)
}

// DeleteUserWealthAndGiftInfo .
func DeleteUserWealthAndGiftInfo(id int) {
	giftinfos, _ := GetGiftList()

	//删除收到的礼物信息
	for {
		selectsentence := lib.SQLSentence(lib.SQLMapSelectGoldBeansByID)
		updatesentence := lib.SQLSentence(lib.SQLMapUpdateGoldBeansByID)
		deletesentence := lib.SQLSentence(lib.SQLMapDeleteGiftConsumeInfo)

		recvlist, _ := getRecvGiftListByID(id, 1, 1000)
		if 0 == len(recvlist) {
			break
		}

		for _, info := range recvlist {
			var beans, consumed int
			var value int
			var sender int

			for _, gift := range giftinfos {
				if gift.ID == info.GiftID {
					value = gift.Price * info.GiftNum
					sender = info.Person.ID
					break
				}
			}

			lib.SQLQueryRow(selectsentence, sender).Scan(&beans, consumed)
			if consumed >= value {
				consumed = consumed - value
			}
			lib.SQLExec(updatesentence, sender, beans, consumed)
			lib.SQLExec(deletesentence, info.ID)

			lib.DelRedisGiftSendList(sender)
			lib.DelRedisUserInfo(sender)
		}
	}

	//删除送出的礼物信息
	for {
		selectsentence := lib.SQLSentence(lib.SQLMapSelectReceiveValueByID)
		updatesentence := lib.SQLSentence(lib.SQLMapUpdateReceiveValueByID)
		deletesentence := lib.SQLSentence(lib.SQLMapDeleteGiftConsumeInfo)

		sendlist, _ := getSendGiftListByID(id, 1, 1000)
		if 0 == len(sendlist) {
			break
		}

		for _, info := range sendlist {
			var value int
			var receiver int
			var receivevalue int

			for _, gift := range giftinfos {
				if gift.ID == info.GiftID {
					value = gift.Price * info.GiftNum
					receiver = info.Person.ID
					break
				}
			}

			lib.SQLQueryRow(selectsentence, receiver).Scan(&receivevalue)
			if receivevalue >= value {
				receivevalue = receivevalue - value
			}
			lib.SQLExec(updatesentence, receivevalue, receiver)
			lib.SQLExec(deletesentence, info.ID)

			lib.DelRedisGiftRecvList(receiver)
			lib.DelRedisUserInfo(receiver)
		}
	}

	delwealthSentence := lib.SQLSentence(lib.SQLMapDeleteWealth)
	lib.SQLExec(delwealthSentence, id)
	lib.DelRedisGiftRecvList(id)
	lib.DelRedisGiftSendList(id)
}

// DeleteGiftInfoByUserIDAndGiftID .
func DeleteGiftInfoByUserIDAndGiftID(id, giftid int) {
	giftinfos, _ := GetGiftList()

	//删除送出的指定礼物信息
	for {
		selectsentence := lib.SQLSentence(lib.SQLMapSelectReceiveValueByID)
		updatesentence := lib.SQLSentence(lib.SQLMapUpdateReceiveValueByID)
		deletesentence := lib.SQLSentence(lib.SQLMapDeleteGiftConsumeInfo)

		sendlist, _ := getSendGiftListByID(id, 1, 100000)
		if 0 == len(sendlist) {
			break
		}

		for _, info := range sendlist {
			if giftid != info.GiftID {
				continue
			}

			var value int
			var receiver int
			var receivevalue int

			for _, gift := range giftinfos {
				if gift.ID == giftid {
					value = gift.Price * info.GiftNum
					receiver = info.Person.ID
					break
				}
			}

			lib.SQLQueryRow(selectsentence, receiver).Scan(&receivevalue)
			if receivevalue >= value {
				receivevalue = receivevalue - value
			}
			lib.SQLExec(updatesentence, receivevalue, receiver)
			lib.SQLExec(deletesentence, info.ID)

			lib.DelRedisGiftRecvList(receiver)
			lib.DelRedisUserInfo(receiver)
		}
	}

	lib.DelRedisGiftSendList(id)
}
