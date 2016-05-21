package pagezhenqing

import (
	"regexp"

	"herefriend/crawler/request"
)

type PageList struct {
	crawled bool
	ids     []string
	request *request.Request
}

type ProvinceInfo struct {
	Number int
	Desc   string
}

var gGenderList = []int{
	0, //man
	1, //woman
}

var gSearchFmt = "http://www.lol99.com/search/?type=2&icon_updated=1&gender=%d&province=%d&page=%d"
var gVerboseFmt = "http://www.lol99.com/member/view.php?uid=%s"
var gAlbumFmt = "http://www.lol99.com/member/photo_list.php?uid=%s"
var gRegex *regexp.Regexp

var gProvinceList = [...]ProvinceInfo{
	{11, "北京"},
	{12, "天津"},
	{13, "河北"},
	{14, "山西"},
	{15, "内蒙古"},
	{21, "辽宁"},
	{22, "吉林"},
	{23, "黑龙江"},
	{31, "上海"},
	{32, "江苏"},
	{33, "浙江"},
	{34, "安徽"},
	{35, "福建"},
	{36, "江西"},
	{37, "山东"},
	{41, "河南"},
	{42, "湖北"},
	{43, "湖南"},
	{44, "广东"},
	{45, "广西"},
	{46, "海南"},
	{50, "重庆"},
	{51, "四川"},
	{52, "贵州"},
	{53, "云南"},
	{54, "西藏"},
	{61, "陕西"},
	{62, "甘肃"},
	{63, "青海"},
	{64, "宁夏"},
	{65, "新疆"},
	{71, "台湾"},
	{81, "香港"},
	{82, "澳门"},
	{99, "国外"},
}

func init() {
	gRegex = regexp.MustCompile("(?:[^?]*?uid=)(\\d+)")
}

func NewPageList(request *request.Request) *PageList {
	return &PageList{ids: make([]string, 0), request: request}
}

func (this *PageList) GetIds() []string {
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

	/* ids info */
	sel := query.Find(".just_photo > ul > li > .pic > a")
	if sel.Length() > 0 {
		for i, _ := range sel.Nodes {
			href, exist := sel.Eq(i).Attr("href")
			if true == exist {
				idstrs := gRegex.FindAllStringSubmatch(href, -1)
				if nil != idstrs {
					this.ids = append(this.ids, idstrs[0][1])
				}
			}
		}
	}

	this.crawled = true
	return this
}
