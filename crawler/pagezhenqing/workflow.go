package pagezhenqing

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"herefriend/crawler/page"
	"herefriend/crawler/request"
	"herefriend/lib"
)

const (
	gCurrentGender         = 0
	g_CRAW_WORKCHAN_BUFFER = 1
)

var g_workChan chan string
var g_activeChan chan int
var g_stopChan chan int
var g_activenum int
var g_activelock sync.Mutex

func init() {
	g_workChan = make(chan string)
	g_activeChan = make(chan int, g_CRAW_WORKCHAN_BUFFER)
	g_stopChan = make(chan int)

	for i := 0; i < g_CRAW_WORKCHAN_BUFFER; i++ {
		g_activeChan <- 1
	}

	go crawWorkRoutine()
}

func crawWorkRoutine() {
	for idstr := range g_workChan {
		if "<END>" == idstr {
			break
		}

		<-g_activeChan
		g_activelock.Lock()
		g_activenum = g_activenum + 1
		g_activelock.Unlock()

		go crawlUserInfo(idstr)
	}

	var numtmp int
	for {
		g_activelock.Lock()
		numtmp = g_activenum
		g_activelock.Unlock()

		if 0 == numtmp {
			break
		}

		time.Sleep(time.Second)
	}

	g_stopChan <- 1
}

func WaitStop() {
	<-g_stopChan
}

/*
 |    Function: checkUserIdExist
 |      Author: Mr.Sancho
 |        Date: 2016-01-09
 |   Arguments:
 |      Return:
 | Description: 检查user的ID是否存在
 |
*/
func checkUserIdExist(id int) bool {
	var tmpid int
	sentence := lib.SQLSentence(lib.SQLMapSelectCheckZQUserID)
	err := lib.SQLQueryRow(sentence, id).Scan(&tmpid)
	if nil != err || 0 == tmpid {
		return false
	}

	return true
}

/*
 |    Function: crawlUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-01-10
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func crawlUserInfo(idstr string) {
	url := fmt.Sprintf(gVerboseFmt, idstr)
	req := request.NewRequest(url, nil)
	pinfo := NewPageInfo(req)
	pinfo.Crawl()

	url = fmt.Sprintf(gAlbumFmt, idstr)
	req = request.NewRequest(url, nil)
	palbum := NewPageAlbum(req)
	palbum.Crawl()

	pageuser := page.NewPage(nil)
	pageuser.SetHeadImg(pinfo.Headimg)
	pageuser.SetImages(palbum.Pictures)
	pageuser.SetPersonInfo(pinfo.Info)
	pageuser.Save()

	sentence := lib.SQLSentence(lib.SQLMapInsertZQID)
	id, _ := strconv.Atoi(idstr)
	_, err := lib.SQLExec(sentence, id)
	if nil != err {
		fmt.Println(err)
	}
	sentence = lib.SQLSentence(lib.SQLMapInsertHeartbeat)
	_, err = lib.SQLExec(sentence, pageuser.GetUsrId(), pageuser.GetGender(), pageuser.GetProvince())
	if nil != err {
		fmt.Println(err)
	}

	g_activelock.Lock()
	g_activenum = g_activenum - 1
	g_activelock.Unlock()

	g_activeChan <- 1
}

/*
 |    Function: DoCrawl
 |      Author: Mr.Sancho
 |        Date: 2016-01-09
 |   Arguments:
 |      Return:
 | Description: the workflow of crawling
 |
*/
func DoCrawl() {
	var areaindex, areapage int
	var url string
	var sentence string

	sentence = lib.SQLSentence(lib.SQLMapSelectZQProcess)
	err := lib.SQLQueryRow(sentence, gCurrentGender).Scan(&areaindex, &areapage)
	if nil != err {
		panic(err)
	}

	for ; areaindex < len(gProvinceList); areaindex = areaindex + 1 {
		for {
			fmt.Printf("[area] area:%s page:%d\r\n", gProvinceList[areaindex].Desc, areapage)

			sentence = lib.SQLSentence(lib.SQLMapUpdateZQProcess)
			_, err := lib.SQLExec(sentence, areaindex, areapage, gCurrentGender)
			if nil != err {
				panic(err.Error())
			}

			//zhenqing's gender is oppose
			url = fmt.Sprintf(gSearchFmt, 1-gCurrentGender, gProvinceList[areaindex].Number, areapage)
			req := request.NewRequest(url, nil)
			plist := NewPageList(req)
			plist.Crawl()

			ids := plist.GetIds()
			if 0 == len(ids) {
				break
			}

			for _, idstr := range ids {
				id, _ := strconv.Atoi(idstr)
				if true == checkUserIdExist(id) {
					continue
				}
				g_workChan <- idstr
			}

			areapage = areapage + 1
		}

		areapage = 1
	}

	g_workChan <- "<END>"
}

func Start() {
	DoCrawl()
	WaitStop()
}
