package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

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

var gGiftRandomLimit int
var gGiftRandomBuf []giftInfo
var gGiftNumLimit int
var gGiftNumBuf []int
var gGiftNumSource = []giftNum{
	{1, 1, 5000},
	{2, 10, 1000},
	{11, 20, 50},
	{21, 60, 1},
	{61, 90, 1},
	{91, 100, 1},
}

func init() {
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMAP_Select_FakeCount, 0)).Scan(&gFakeGirls)
	lib.SQLQueryRow(lib.SQLSentence(lib.SQLMAP_Select_FakeCount, 1)).Scan(&gFakeGuys)

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

	sentence := lib.SQLSentence(lib.SQLMAP_Select_RandomId, gender)
	othersentence := lib.SQLSentence(lib.SQLMAP_Select_RandomId, 1-gender)
	for {
		// get random id
		err := lib.SQLQueryRow(sentence, lib.Intn(baselimit)).Scan(&id)
		if nil != err || 1 >= id {
			fmt.Println(err)
			continue
		}

		err = lib.SQLQueryRow(othersentence, lib.Intn(otherlimit)).Scan(&toid)
		if nil != err || 1 >= toid {
			fmt.Println(err)
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

		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	sendGiftByGender(0)
}
