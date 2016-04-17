package loveshow

import (
	"regexp"

	"herefriend/crawler/request"
)

type Spacelist struct {
	listid   int
	crawled  bool
	spaceids []string
	request  *request.Request
}

var g_spacelist_regex *regexp.Regexp

func init() {
	g_spacelist_regex, _ = regexp.Compile("(?:<dt><a href=\"space[.]php[?]id=)([^\"]+)(?:\")")
}

func NewSpacelist(request *request.Request) *Spacelist {
	return &Spacelist{listid: request.GetUid(), request: request}
}

func (this *Spacelist) GetSpaceIds() []string {
	return this.spaceids
}

func (this *Spacelist) Crawl() *Spacelist {
	if true == this.crawled {
		return this
	}

	err := this.request.Download()
	if nil != err {
		return this
	}

	/*
	 * do crawl
	 */
	var bodyStr = string(this.request.GetBody())
	spacelist := g_spacelist_regex.FindAllStringSubmatch(bodyStr, -1)
	if nil != spacelist {
		for _, s := range spacelist {
			this.spaceids = append(this.spaceids, s[1])
		}
	}

	this.crawled = true
	return this
}
