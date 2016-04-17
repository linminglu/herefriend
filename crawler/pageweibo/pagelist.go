package pageweibo

import (
	"herefriend/crawler/request"
)

type PageList struct {
	crawled bool
	ids     []string
	request *request.Request
}

type VGirlAreaInfo struct {
	Area    string
	PageNum int
}

var G_VGirls = [...]VGirlAreaInfo{ //about 120000 vgirls
	{"11", 343},
	{"12", 67},
	{"13", 77},
	{"14", 41},
	{"15", 24},
	{"21", 162},
	{"22", 49},
	{"23", 72},
	{"31", 347},
	{"32", 255},
	{"33", 339},
	{"34", 120},
	{"35", 229},
	{"36", 106},
	{"37", 113},
	{"41", 115},
	{"42", 306},
	{"43", 138},
	{"44", 1264},
	{"45", 119},
	{"46", 38},
	{"50", 121},
	{"51", 227},
	{"52", 56},
	{"53", 66},
	{"54", 1},
	{"61", 90},
	{"62", 22},
	{"63", 3},
	{"64", 8},
	{"65", 40},
}

func NewPageList(request *request.Request) *PageList {
	return &PageList{ids: make([]string, 0), request: request}
}

func (this *PageList) GetvgirlIds() []string {
	return this.ids
}

func (this *PageList) Crawl() *PageList {
	if true == this.crawled {
		return this
	}

	err := this.request.Download(nil)
	if nil != err {
		return this
	}

	/*
	 * do crawl
	 */
	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	/* vgirls info */
	sel := query.Find(".vg_short > li")
	if sel.Length() > 0 {
		for i, _ := range sel.Nodes {
			idsel := sel.Eq(i).Find(".show_check")
			if idsel.Length() > 0 {
				id, exist := idsel.Attr("value")
				if true == exist {
					this.ids = append(this.ids, id)
				}
			}
		}
	}

	this.crawled = true
	return this
}
