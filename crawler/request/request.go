package request

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"herefriend/lib"
)

type Request struct {
	uid     int
	url     string
	cookies []*http.Cookie
	body    []byte
	query   *goquery.Document
}

const (
	REQUESTURL_PAGE      = 0
	REQUESTURL_3GPAGE    = 1
	REQUESTURL_3GLOGIN   = 2
	REQUESTURL_SPACELIST = 3
	REQUESTURL_SPACE     = 4
	REQUESTURL_STORY     = 5
)

var g_ReqeustURL = [...]string{
	"http://profile1.baihe.com/?oppID=",
	"http://3g.baihe.com/user/baseinfo?uid=",
	"http://3g.baihe.com/login",
	"http://story.baihe.com/story.php?story_list&dr=&p=",
	"http://story.baihe.com/space.php?id=",
	"http://story.baihe.com/story.php?id="}

var g_3gUserAgent = "Mozilla/5.0 (iPad; CPU OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B554a Safari/9537.53"

func NewRequest(url string, cookies []*http.Cookie) *Request {
	return &Request{url: url, cookies: cookies}
}

func NewRequestBH(urltype, uid int, cookies []*http.Cookie) *Request {
	req := &Request{uid: uid, url: g_ReqeustURL[urltype] + strconv.Itoa(uid), cookies: cookies}

	return req
}

func NewRequestUrl(urltype int, cookies []*http.Cookie) *Request {
	return &Request{url: g_ReqeustURL[urltype], cookies: cookies}
}

func (this *Request) GetUrl() string {
	return this.url
}

func (this *Request) SetUrl(url string) {
	this.url = url
}

func (this *Request) GetUid() int {
	return this.uid
}

func (this *Request) GetCookies() []*http.Cookie {
	return this.cookies
}

func (this *Request) SetCookies(cookies []*http.Cookie) {
	this.cookies = cookies
}

func (this *Request) GetBody() []byte {
	return this.body
}

func (this *Request) SetBody(body []byte) {
	this.body = body
}

func (this *Request) Download(cookies []*http.Cookie) error {
	if 0 == len(this.url) || nil != this.body {
		return nil
	}

	fmt.Println("【Download】" + this.url)
	resp, err := lib.Get(this.url, cookies)
	if nil != err {
		fmt.Println("【Download】failed: " + err.Error())
		return err
	} else {
		defer resp.Body.Close()

		bytebuf, err := lib.ConvertToUtf8(resp.Header.Get("Content-Type"), resp.Body)
		if nil != err {
			fmt.Println("【Download】failed: convtUtf8 failed " + err.Error())
			return err
		}

		bytebuf = bytes.Replace(bytebuf, []byte("\n"), nil, -1)
		bytebuf = bytes.Replace(bytebuf, []byte("\r"), nil, -1)
		bytebuf = bytes.Replace(bytebuf, []byte("\t"), nil, -1)

		this.body = bytebuf
		return nil
	}
}

/*
 *
 *    Function: Download3g
 *      Author: sunchao
 *        Date: 15/10/25
 * Description: download the 3g request
 *
 */
func (this *Request) Download3g() error {
	if 0 == len(this.url) || nil != this.body {
		return nil
	}

	httpreq, err := http.NewRequest("GET", this.url, nil)
	if nil != err {
		return err
	}

	if nil != this.cookies {
		for i := range this.cookies {
			httpreq.AddCookie(this.cookies[i])
		}
	}

	httpreq.Header.Add("User-Agent", g_3gUserAgent)
	httpreq.Header.Add("X-FirePHP-Version", "0.0.6")
	client := &http.Client{
		Transport: &http.Transport{Dial: lib.HTTPTimeoutDial},
	}

	fmt.Println("Downloading: " + this.url)
	resp, err := client.Do(httpreq)
	if nil != err {
		fmt.Println("[Download failed] " + err.Error())
		return err
	} else {
		defer resp.Body.Close()

		bytebuf, err := lib.ConvertToUtf8(resp.Header.Get("Content-Type"), resp.Body)
		if nil != err {
			fmt.Println("Download failed: convtUtf8 failed " + err.Error())
			return err
		}

		bytebuf = bytes.Replace(bytebuf, []byte("\n"), nil, -1)
		bytebuf = bytes.Replace(bytebuf, []byte("\r"), nil, -1)
		bytebuf = bytes.Replace(bytebuf, []byte("\t"), nil, -1)

		this.body = bytebuf
		return nil
	}
}

func (this *Request) Get3gResponse() (*http.Response, error) {
	httpreq, err := http.NewRequest("GET", this.url, nil)
	if nil != err {
		return nil, err
	}

	httpreq.Header.Add("User-Agent", g_3gUserAgent)
	httpreq.Header.Add("X-FirePHP-Version", "0.0.6")
	client := &http.Client{
		Transport: &http.Transport{Dial: lib.HTTPTimeoutDial},
	}

	fmt.Println("Downloading: " + this.url)
	resp, err := client.Do(httpreq)
	if nil != err {
		fmt.Println("[Download failed] " + err.Error())
		return nil, err
	}

	return resp, nil
}

/*
 *
 *    Function: CreateQuery
 *      Author: sunchao
 *        Date: 15/8/22
 * Description: create the Query document
 *
 */
func (this *Request) CreateQuery() (*goquery.Document, error) {
	bodyReader := bytes.NewReader(this.body)
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if nil != err {
		fmt.Println(err)
		return nil, err
	} else {
		this.query = doc
		return doc, nil
	}
}
