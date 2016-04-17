package page3g

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"herefriend/crawler/dbtables"
	"herefriend/crawler/request"
)

type Page3g struct {
	usrid   int
	crawled bool
	gender  int
	headimg string
	request *request.Request
}

var g_regexSessionId *regexp.Regexp
var g_Cookies []*http.Cookie

func init() {
	g_regexSessionId, _ = regexp.Compile("(?:PHPSESSID=)([^;]+)")
}

func NewPage(request *request.Request, gender int) *Page3g {
	request.SetCookies(g_Cookies)
	return &Page3g{usrid: request.GetUid(), request: request, gender: gender}
}

func generateCookies() {
	request := request.NewRequestUrl(request.REQUESTURL_3GLOGIN, nil)
	resp, err := request.Get3gResponse()
	if nil != err {
		return
	}

	defer resp.Body.Close()
	phpsession_str := resp.Header.Get("Set-Cookie")
	matches := g_regexSessionId.FindAllStringSubmatch(phpsession_str, -1)
	if nil != matches {
		g_Cookies = []*http.Cookie{
			&http.Cookie{Name: "PHPSESSID", Value: matches[0][1], Path: "/"},
			&http.Cookie{Name: "SESSION_COOKIE", Value: "100", Path: "/"},
		}
	}
}

/*
 * the crawl action, get all the <a> and <img> from the page
 */
func (this *Page3g) Crawl() *Page3g {
	if true == this.crawled {
		return this
	}

	err := this.request.Download3g()
	if nil != err {
		return this
	}

	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	/* check if need get a new session id */
	sel := query.Find("#firstLogin")
	if nil != sel && sel.Length() > 0 {
		generateCookies()
		this.request = request.NewRequest(request.REQUESTURL_3GLOGIN, this.usrid, nil)
		this.request.SetCookies(g_Cookies)

		this.request.SetBody(nil)
		this.request.Download3g()
		query, _ = this.request.CreateQuery()
	}

	sel = query.Find(".personal > div > dl > dt > a > img")
	if nil != sel && sel.Length() > 0 {
		/* crawl the picture */
		this.headimg, _ = sel.Attr("src")
	}

	this.crawled = true
	return this
}

/*
 * download url as image and put it to QiNiu
 */
func (this *Page3g) ReplaceHeadImage(id, gender int, oldfilename string) {
	if true != this.crawled {
		return
	}

	/*
	 * delete old head file
	 */
	dbopt.DeleteOldHeadPictureById(id, gender)
	request.DeleteImageFromQiniu("images/" + strconv.Itoa(id) + "/" + oldfilename)

	/*
	 * get new head file
	 */
	/* check if the new head is valid */
	if true == strings.Contains(this.headimg, "images/default_pictures/80_100/havepic_") {
		return
	}

	this.headimg = strings.Replace(this.headimg, "80_100", "120_150", 1)
	resp, err := request.DownloadImage(this.headimg)
	if nil != err {
		return
	}

	defer resp.Body.Close()
	strslice := strings.Split(this.headimg, "/")
	imgname := "icon-" + strslice[len(strslice)-1]

	/* upload to qiniu */
	err = request.PostImageToQiniu("images/"+strconv.Itoa(id)+"/"+imgname, resp.Body)
	if nil == err {
		dbopt.InsertPictureById(id, gender, imgname, true, true)
	}

	return
}

func repleaceHeadCallBack(id, gender int, oldfilename string) {
	page := NewPage(request.NewRequest(request.REQUESTURL_3GPAGE, id, nil), gender).Crawl()
	ReplaceHeadImage(id, gender, oldfilename)
}

func Start() {
	dbtables.DoRenewHeadPicture(0, repleaceHeadCallBack)
	dbtables.DoRenewHeadPicture(1, repleaceHeadCallBack)
}
