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

type numRangeInfo struct {
	Start int
	End   int
	Num   int
}

var gFakeGirls int
var gFakeGuys int
var gGiftList []giftInfo
var gGiftListNum int

//随机礼物Id用
var gGiftRandomLimit int
var gGiftRandomBuf []giftInfo

//随机礼物数量用
var gGiftNumLimit int
var gGiftNumBuf []int
var gGiftNumSource = []numRangeInfo{
	{1, 1, 50000},
	{2, 10, 10000},
	{11, 20, 500},
	{21, 100, 1},
}

//随机送礼物的人
var gUserNumLimit int
var gUserNumBuf []int
var gUserNumSource = []numRangeInfo{
	{1, 2, 5},
	{3, 20, 500},
	{21, 50, 10},
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
	gGiftListNum = len(gGiftList)
	fmt.Printf("list num = %d\n", gGiftListNum)

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
		if i == len(gGiftList)-1 {
			break
		}

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

	for _, n := range gUserNumSource {
		for i := n.Start; i <= n.End; i = i + 1 {
			for j := 0; j < n.Num; j = j + 1 {
				gUserNumBuf = append(gUserNumBuf, i)
			}
		}
	}

	gGiftRandomLimit = len(gGiftRandomBuf)
	gGiftNumLimit = len(gGiftNumBuf)
	gUserNumLimit = len(gUserNumBuf)
}

func getRandomUserIdByGender(gender int) (error, int) {
	baselimit := func() int {
		if 0 == gender {
			return gFakeGirls
		} else {
			return gFakeGuys
		}
	}()

	sentence := func() string {
		if 0 == gender {
			return "select id from heartbeat where gender=0 and id in (select id from girls where usertype!=1) limit ?,1"
		} else {
			return "select id from heartbeat where gender=1 and id in (select id from guys where usertype!=1) limit ?,1"
		}
	}()

	var id int
	randomvalue := lib.Intn(baselimit)
	err := lib.SQLQueryRow(sentence, randomvalue).Scan(&id)
	if nil != err || 1 >= id {
		if nil != err {
			lib.SQLError(sentence, err, randomvalue)
		}
		return err, 0
	}

	return nil, id
}

func getRandomGiftNumberByGiftId(id int) int {
	var price int
	for _, info := range gGiftList {
		if info.Id == id {
			price = info.Price
			break
		}
	}

	if 0 == price {
		return 0
	}

	if price > 1000 {
		return lib.Intn(5)
	}

	var index int
	var num int

	//每次不能超过10000
	for {
		index = lib.Intn(gGiftNumLimit)
		num = gGiftNumBuf[index]

		if num*price < 2000 {
			break
		}
	}

	return num
}

func getRandomUserNumber() int {
	index := lib.Intn(gUserNumLimit)
	return gUserNumBuf[index]
}

func resetfulSendGift(id, gender, toid, giftid, num int) {
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
}

func randomSendGiftBySender(id, gender int) {
	sendernum := getRandomUserNumber()
	for i := 0; i < sendernum; i++ {
		err, toid := getRandomUserIdByGender(1 - gender)
		if nil != err {
			fmt.Println(err)
			continue
		}

		index := lib.Intn(gGiftRandomLimit)
		giftid := gGiftRandomBuf[index].Id
		resetfulSendGift(id, gender, toid, giftid, getRandomGiftNumberByGiftId(giftid))
		time.Sleep(time.Minute * 20)
	}
}

func DoSend(gender int) {
	for {
		sendermap := make(map[int]bool)

		//随机一个receiver
		err, toid := getRandomUserIdByGender(1 - gender)
		if nil != err {
			fmt.Println(err)
			continue
		}

		//随机收到礼物数量
		giftnum := lib.Intn(gGiftListNum)
		if 0 == giftnum {
			giftnum = gGiftListNum
		}

		if 4 > giftnum {
			giftnum = 4
		}

		fmt.Printf("user %d will receive %d gifts\n", toid, giftnum)
		for n := 0; n < giftnum; n++ {
			giftid := gGiftList[n].Id

			//随机送礼物的人数
			sendernum := getRandomUserNumber()
			fmt.Printf("user %d will receive gift %d from %d users\n", toid, giftid, sendernum)
			for i := 0; i < sendernum; i++ {
				err, id := getRandomUserIdByGender(gender)
				if nil != err {
					fmt.Println(err)
					continue
				}
				sendermap[id] = true

				resetfulSendGift(id, gender, toid, giftid, getRandomGiftNumberByGiftId(giftid))
				time.Sleep(time.Minute * 30)
			}
		}

		//每个人随机送其他人礼物
		for k, _ := range sendermap {
			randomSendGiftBySender(k, gender)
		}
	}
}

func main() {
	go DoSend(0)
	DoSend(1)
}
