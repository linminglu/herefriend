package pagezhenqing

import (
	"strings"

	"herefriend/crawler/request"
)

type PageAlbum struct {
	Pictures []string
	crawled  bool
	request  *request.Request
}

func NewPageAlbum(request *request.Request) *PageAlbum {
	return &PageAlbum{request: request}
}

func (this *PageAlbum) Crawl() *PageAlbum {
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

	/* pictures */
	sel := query.Find(".photo_list > .yc > ul > li > a > img")
	if sel.Length() > 0 {
		for i, _ := range sel.Nodes {
			href, exist := sel.Eq(i).Attr("src")
			if true == exist {
				href = strings.TrimSpace(strings.Replace(href, "_slt.", "_y.", -1))
				this.Pictures = append(this.Pictures, href)
			}
		}
	}

	sel = query.Find(".photo_list > .yc > ul > li > .phototm > a > img")
	if sel.Length() > 0 {
		for i, _ := range sel.Nodes {
			href, exist := sel.Eq(i).Attr("src")
			if true == exist {
				href = strings.TrimSpace(strings.Replace(href, "_slt.", "_y.", -1))
				this.Pictures = append(this.Pictures, href)
			}
		}
	}

	this.crawled = true
	return this
}
