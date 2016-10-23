package pageweibo

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"herefriend/crawler/page"
	"herefriend/crawler/request"
	"herefriend/lib"
)

const (
	g_CRAW_WORKCHAN_BUFFER = 100
)

var g_workChan chan string
var g_activeChan chan int
var g_stopChan chan int
var g_activenum int
var g_activelock sync.Mutex
var g_randomidlock sync.Mutex

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

		go crawlVGirlInfo(idstr)
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
 |    Function: checkVGirlIdExist
 |      Author: Mr.Sancho
 |        Date: 2016-01-09
 |   Arguments:
 |      Return:
 | Description: 检查vgirl的ID是否存在
 |
*/
func checkVGirlIdExist(id string) bool {
	var tmpid string
	sentence := lib.SQLSentence(lib.SQLMapSelectCheckVGirlID)
	err := lib.SQLQueryRow(sentence, id).Scan(&tmpid)
	if nil != err || "" == tmpid {
		return false
	}

	return true
}

/*
 |    Function: crawlVGirlInfo
 |      Author: Mr.Sancho
 |        Date: 2016-01-10
 |   Arguments:
 |      Return:
 | Description:
 |
*/
func crawlVGirlInfo(idstr string) {
	var crawid int

	//fensi
	url := "http://vgirl.weibo.com/weibo/index.php?uid=" + idstr
	req := request.NewRequest(url, nil)
	pfensi := NewPageInfo(req)
	pfensi.Crawl()

	//pictures
	req = request.NewRequest("http://vgirl.weibo.com/"+idstr, nil)
	ppic := NewPagePic(req)
	ppic.Crawl()

	//get a random girl from girlsid
	g_randomidlock.Lock()

	sentence := lib.SQLSentence(lib.SQLMapSelectRandomUncrawlGirlsID)
	randomlimit := rand.Intn(2000000)
	lib.SQLQueryRow(sentence, randomlimit).Scan(&crawid)

	sentence = lib.SQLSentence(lib.SQLMapDeleteUnCrawledGirlsID)
	lib.SQLExec(sentence, crawid)

	g_randomidlock.Unlock()

	//craw the girl info and use the pictures instead original pcitures
	fmt.Printf("[Start to craw bh_id] crawid: %d\r\n", crawid)
	req = request.NewRequestBH(request.REQUESTURL_PAGE, crawid, nil)
	pageuser := page.NewPage(req)

	for {
		pageuser.Crawl(false)

		if true != pageuser.IsCrawled() {
			time.Sleep(time.Second)
			continue
		}

		break
	}

	pageuser.SetHeadImg(ppic.headimg)
	pageuser.SetImages(ppic.imgs)
	pageuser.Save()

	sentence = lib.SQLSentence(lib.SQLMapInsertVGirlID)
	_, err := lib.SQLExec(sentence, idstr, pfensi.fensi, pageuser.GetUsrId())
	if nil != err {
		fmt.Println(err)
	}
	sentence = lib.SQLSentence(lib.SQLMapInsertHeartbeat)
	_, err = lib.SQLExec(sentence, pageuser.GetUsrId(), 0)
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

	sentence = lib.SQLSentence(lib.SQLMapSelectVGirlProcess)
	err := lib.SQLQueryRow(sentence).Scan(&areaindex, &areapage)
	if nil != err {
		panic(err.Error())
	}

	for ; areaindex < len(G_VGirls); areaindex = areaindex + 1 {
		for ; areapage < G_VGirls[areaindex].PageNum; areapage = areapage + 1 {
			fmt.Printf("[area] area:%s page:%d\r\n", G_VGirls[areaindex].Area, areapage+1)

			sentence = lib.SQLSentence(lib.SQLMapUpdateVGirlProcess)
			_, err := lib.SQLExec(sentence, areaindex, areapage)
			if nil != err {
				panic(err.Error())
			}

			//list
			url = "http://vgirl.weibo.com/area.php?" + fmt.Sprintf("p=%s&page=%d", G_VGirls[areaindex].Area, areapage+1)
			req := request.NewRequest(url, nil)
			plist := NewPageList(req)
			plist.Crawl()

			vGirlIds := plist.GetvgirlIds()
			for _, idstr := range vGirlIds {
				if true == checkVGirlIdExist(idstr) {
					continue
				}

				g_workChan <- idstr
			}
		}

		areapage = 1
	}

	g_workChan <- "<END>"
}

func Start() {
	DoCrawl()
	WaitStop()
}
