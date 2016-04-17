package loveshow

import (
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"herefriend/crawler/request"
)

var g_bless_regexp_id *regexp.Regexp

type Bless struct {
	id        int
	age       int
	name      string
	district  string
	education string
	bless     string
	timeUTC   time.Time
}

func init() {
	g_bless_regexp_id, _ = regexp.Compile("(?:href=[^=]+=)(\\d+)(?:[^>]+>)([^<]+)(?:[^\\d]+)(\\d+)(?:\\s*岁\\s*)([^\\s]+)(?:\\s*)([^\\s]+)(?:[^\\d]+)([^<]+)")
}

func grepBless(blessSel *goquery.Selection) Bless {
	var bless Bless

	htmlStr, _ := blessSel.Find("dt").Html()
	idArray := g_bless_regexp_id.FindAllStringSubmatch(htmlStr, -1)
	if nil != idArray {
		bless.id, _ = strconv.Atoi(idArray[0][1])
		bless.name = idArray[0][2]
		bless.age, _ = strconv.Atoi(idArray[0][3])
		bless.district = idArray[0][4]
		bless.education = idArray[0][5]
		bless.timeUTC, _ = grepTimeUTC(idArray[0][6])
	}

	bless.bless = blessSel.Find("dd").Text()

	return bless
}

func grepBlessesByPage(url string, page int) []Bless {
	req := request.Request{}
	req.SetUrl(url + "&p=" + strconv.Itoa(page))

	err := req.Download()
	if nil != err {
		return nil
	}

	query, err := req.CreateQuery()
	if nil != err {
		return nil
	}

	/* blesses */
	var blesses []Bless
	sel := query.Find(".blog_publish > dl")
	if sel.Length() > 0 {
		for i, _ := range sel.Nodes {
			blesses = append(blesses, grepBless(sel.Eq(i)))
		}
	}

	return blesses
}
