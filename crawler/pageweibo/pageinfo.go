package pageweibo

import (
	"strconv"

	"herefriend/crawler/request"
)

type PageInfo struct {
	crawled bool
	fensi   int
	request *request.Request
}

func NewPageInfo(request *request.Request) *PageInfo {
	return &PageInfo{request: request}
}

func (this *PageInfo) GetFensi() int {
	return this.fensi
}

/*
 * the crawl action, get all the <a> and <img> from the page
 */
func (this *PageInfo) Crawl() *PageInfo {
	if true == this.crawled {
		return this
	}

	err := this.request.Download(nil)
	if nil != err {
		return this
	}

	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	/* gender */
	sel := query.Find(".user_info_list > li")
	if sel.Length() > 0 {
		fensiSel := sel.Eq(1).Find("a > em")
		if fensiSel.Length() > 0 {
			s, _ := fensiSel.Html()
			this.fensi, _ = strconv.Atoi(s)
		}
	}

	this.crawled = true
	return this
}
