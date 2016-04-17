package loveshow

import (
	"strings"
	"time"

	"herefriend/crawler/dbtables"
	"herefriend/crawler/request"
)

type Story struct {
	storyid int
	crawled bool
	title   string
	story   string
	timeUTC time.Time
	blesses []Bless
	request *request.Request
}

func NewStory(request *request.Request) *Story {
	return &Story{storyid: request.GetUid(), request: request}
}

func (this *Story) Crawl() *Story {
	if true == this.crawled {
		return this
	}

	err := this.request.Download()
	if nil != err {
		return this
	}

	body := this.request.GetBody()
	strbuf := string(body)
	strbuf = strings.Replace(strbuf, "<div class=\"clear\">", "</div><div><div class=\"clear\">", 1)
	this.request.SetBody([]byte(strbuf))

	/*
	 * do crawl
	 */
	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	/* story info */
	sel := query.Find(".story_colfl")
	if sel.Length() > 0 {
		//title
		titleSel := sel.Find(".title")
		this.title = eascapString(titleSel.Children().Eq(0).Text())
		if len([]rune(this.title)) > 50 {
			this.title = string([]rune(this.title)[:2048])
		}

		//time
		timeUTC, result := grepTimeUTC(titleSel.Find("p").Text())
		if true == result {
			this.timeUTC = timeUTC
		}

		//txt
		storySel := sel.Find(".txt")
		this.story = eascapString(strings.TrimSpace(storySel.Text()))
		if len([]rune(this.story)) > 2048 {
			this.story = string([]rune(this.story)[:2048])
		}
	}

	/* blesses */
	sel = query.Find(".blog_publish > dl")
	if sel.Length() > 0 {
		for i, _ := range sel.Nodes {
			this.blesses = append(this.blesses, grepBless(sel.Eq(i)))
		}

		/* following pages blesses */
		end := false
		page := 2
		for true != end {
			blesses := grepBlessesByPage(this.request.GetUrl(), page)
			if nil == blesses || 0 == len(blesses) {
				end = true
			} else {
				for _, b := range blesses {
					this.blesses = append(this.blesses, b)
				}

				page = page + 1
			}
		}
	}

	this.crawled = true
	return this
}

func (this *Story) Save(loveshowid int) {
	if true != this.crawled {
		return
	}

	for _, b := range this.blesses {
		checkAndGrepPersonInfo(b.id, b.district)

		dbopt.InsertToLoveshowBless(loveshowid, b.id, b.age, b.timeUTC, b.name, b.district, b.education, b.bless)
	}
}
