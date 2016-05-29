package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/cihub/seelog"

	"herefriend/lib"
)

type giftInfo struct {
	Id    int //礼物固定id
	Price int //价格(beans)
}

type giftNum struct {
	Start int
	End   int
	Num   int
}

var gFakeGirls int
var gFakeGuys int
var gGiftList []giftInfo

//随机礼物Id用
var gGiftRandomLimit int
var gGiftRandomBuf []giftInfo

//随机礼物数量用
var gGiftNumLimit int
var gGiftNumBuf []int
var gGiftNumSource = []giftNum{
	{1, 1, 50000},
	{2, 10, 10000},
	{11, 20, 500},
	{21, 100, 1},
}

func init() {
	lib.SQLQueryRow("select count(*) from heartbeat where gender=0 and id in (select id from girls where usertype!=1)").Scan(&gFakeGirls)
	lib.SQLQueryRow("select count(*) from heartbeat where gender=1 and id in (select id from guys where usertype!=1)").Scan(&gFakeGuys)

	resp, err := lib.Get("http://localhost:8080/Gift/GiftList", nil)
	if nil != err {
		panic(err)
	}

	bytebuf, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}
	resp.Body.Close()

	err = json.Unmarshal(bytebuf, &gGiftList)
	if nil != err {
		panic(err)
	}

	var max, min, sum int
	var dot int

	min = 999999
	for _, info := range gGiftList {
		if max < info.Price {
			max = info.Price
		}

		if min > info.Price {
			min = info.Price
		}
	}

	sum = max + min
	for i, info := range gGiftList {
		dot = sum - gGiftList[i].Price

		//按照礼物价值修改概率
		if gGiftList[i].Price < 100 {
			dot = dot * 100
		} else if gGiftList[i].Price < 1000 {
			dot = dot * 2
		} else if dot > 10 {
			dot = dot / 10
		}

		for i := 0; i < dot; i = i + 1 {
			gGiftRandomBuf = append(gGiftRandomBuf, info)
		}
	}

	for _, s := range gGiftNumSource {
		for i := s.Start; i <= s.End; i = i + 1 {
			for j := 0; j < s.Num; j = j + 1 {
				gGiftNumBuf = append(gGiftNumBuf, i)
			}
		}
	}

	gGiftRandomLimit = len(gGiftRandomBuf)
	gGiftNumLimit = len(gGiftNumBuf)
}

func sendGiftByGender(gender int) {
	var id int
	var toid int

	var baselimit, otherlimit = func() (int, int) {
		if 0 == gender {
			return gFakeGirls, gFakeGuys
		} else {
			return gFakeGuys, gFakeGirls
		}
	}()

	sentence, othersentence := func() (a, b string) {
		if 0 == gender {
			a = "select id from heartbeat where gender=0 and id in (select id from girls where usertype!=1) limit ?,1"
			b = "select id from heartbeat where gender=1 and id in (select id from guys where usertype!=1) limit ?,1"
		} else {
			a = "select id from heartbeat where gender=1 and id in (select id from guys where usertype!=1) limit ?,1"
			b = "select id from heartbeat where gender=0 and id in (select id from girls where usertype!=1) limit ?,1"
		}
		return
	}()

	for {
		// get random id
		err := lib.SQLQueryRow(sentence, lib.Intn(baselimit)).Scan(&id)
		if nil != err || 1 >= id {
			log.Errorf("SQLQueryRow Error: %s %v\n", sentence, err)
			continue
		}

		err = lib.SQLQueryRow(othersentence, lib.Intn(otherlimit)).Scan(&toid)
		if nil != err || 1 >= toid {
			log.Errorf("SQLQueryRow Error: %s %v\n", othersentence, err)
			continue
		}

		index := lib.Intn(gGiftRandomLimit)
		giftid := gGiftRandomBuf[index].Id

		numindex := lib.Intn(gGiftNumLimit)
		num := gGiftNumBuf[numindex]
		fmt.Printf("http://localhost:8080/cms/PresentGift?id=%d&gender=%d&toid=%d&giftid=%d&num=%d\r\n", id, gender, toid, giftid, num)
		resp, err := lib.Get(fmt.Sprintf("http://localhost:8080/cms/PresentGift?id=%d&gender=%d&toid=%d&giftid=%d&num=%d", id, gender, toid, giftid, num), nil)
		if nil != err {
			fmt.Println(err)
		} else {
			resp.Body.Close()
		}

		lib.DelRedisGiftSendList(id)
		lib.DelRedisGiftRecvList(toid)
		lib.DelRedisUserInfo(id)
		lib.DelRedisUserInfo(toid)

		time.Sleep(time.Second)
	}
}

func main() {
	go sendGiftByGender(0)
	sendGiftByGender(1)
}
