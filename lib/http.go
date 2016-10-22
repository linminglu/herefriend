package lib

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/charset"

	"herefriend/common"
)

/*
 *
 *    Function: ConvertToUtf8
 *      Author: sunchao
 *        Date: 15/8/16
 * Description: Charset auto determine. Use golang.org/x/net/html/charset. Get html body and change it to utf-8
 *
 */
func ConvertToUtf8(contentTypeStr string, r io.ReadCloser) ([]byte, error) {
	destReader, err := charset.NewReader(r, contentTypeStr)
	if err != nil {
		destReader = r
	}

	return ioutil.ReadAll(destReader)
}

/*
 *
 *    Function: DownloadUrl
 *      Author: sunchao
 *        Date: 15/8/23
 * Description: Download the content of the url and return the response
 *
 */
func Get(url string, cookies []*http.Cookie) (*http.Response, error) {
	httpreq, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return nil, err
	}

	if nil != cookies {
		for i := range cookies {
			httpreq.AddCookie(cookies[i])
		}
	}

	client := &http.Client{
		Transport: &http.Transport{Dial: HttpTimeoutDial, DisableKeepAlives: true},
	}

	httpreq.Header.Add("User-Agent", common.ClientAgent)
	resp, err := client.Do(httpreq)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}

func Post(url string, cookies []*http.Cookie) (*http.Response, error) {
	httpreq, err := http.NewRequest("POST", url, nil)
	if nil != err {
		return nil, err
	}

	if nil != cookies {
		for i := range cookies {
			httpreq.AddCookie(cookies[i])
		}
	}

	client := &http.Client{
		Transport: &http.Transport{Dial: HttpTimeoutDial, DisableKeepAlives: true},
	}

	httpreq.Header.Add("User-Agent", common.ClientAgent)
	resp, err := client.Do(httpreq)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	return resp, nil
}

func GetResultByMethod(method, url string, cookies []*http.Cookie) ([]byte, error) {
	httpreq, err := http.NewRequest(method, url, nil)
	if nil != err {
		return nil, err
	}

	if nil != cookies {
		for i := range cookies {
			httpreq.AddCookie(cookies[i])
		}
	}

	client := &http.Client{
		Transport: &http.Transport{Dial: HttpTimeoutDial, DisableKeepAlives: true},
	}

	httpreq.Header.Add("User-Agent", common.ClientAgent)
	resp, err := client.Do(httpreq)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	defer resp.Body.Close()

	bytebuf, err := ConvertToUtf8(resp.Header.Get("Content-Type"), resp.Body)
	if nil != err {
		log.Error(err)
		return nil, err
	}

	return bytebuf, nil
}

/*
 *
 *    Function: get_pageid_count_fromreq
 *      Author: sunchao
 *        Date: 15/8/16
 * Description: 从请求中获取page编号和每页数目
 *
 */
func Get_pageid_count_fromreq(c *gin.Context) (int, int) {
	pageStr := c.Query("page")
	countStr := c.Query("count")
	if pageStr == "" || pageStr == "0" {
		pageStr = "1"
	}

	if countStr == "" {
		countStr = "10"
	}

	pageid, _ := strconv.Atoi(pageStr)
	count, _ := strconv.Atoi(countStr)

	return pageid, count
}

func GetCountRequestArgument(c *gin.Context) int {
	countStr := c.Query("count")
	if "" == countStr {
		countStr = "10"
	}

	count, _ := strconv.Atoi(countStr)
	return count
}

/*
 *
 *    Function: HttpTimeoutDial
 *      Author: sunchao
 *        Date: 15/8/23
 * Description: dial function
 *
 */
func HttpTimeoutDial(netw, addr string) (net.Conn, error) {
	deadline := time.Now().Add(time.Duration(60) * time.Minute)
	c, err := net.DialTimeout(netw, addr, time.Duration(60)*time.Minute)
	if err != nil {
		return nil, err
	}

	c.SetDeadline(deadline)
	return c, nil
}
