package loveshow

import (
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"herefriend/crawler/dbtables"
	"herefriend/crawler/request"
)

type spacePersonInfo struct {
	id       int
	age      int
	name     string
	imgurl   string
	district string
}

type Space struct {
	spaceid   int //空间id
	stordyid  int
	loaded    bool
	crawled   bool
	falldays  int             //见面多少天之后就恋爱了
	blessnum  int             //收到的祝福数量
	status    string          //目前的恋爱状态
	girl      spacePersonInfo //女生信息
	guy       spacePersonInfo //男生信息
	storyinfo *Story          //sotry
	imgs      []string        //恋爱秀照片
	body      []byte
	request   *request.Request
	querydoc  *goquery.Document
}

var g_space_regexp_id *regexp.Regexp
var g_space_regexp_info *regexp.Regexp
var g_space_regexp_status *regexp.Regexp
var g_space_regexp_story *regexp.Regexp

func init() {
	g_space_regexp_id, _ = regexp.Compile("(?:[^=]+=)(\\d+)")
	g_space_regexp_info, _ = regexp.Compile("(?:\\s*)([^\\s]+)(?:\\s*年龄：)([\\d]+)(?:\\s*岁\\s*地区：)([^\\s]+)")
	g_space_regexp_status, _ = regexp.Compile("(?:<strong>在百合 )(\\d+)(?:[^<]+<\\/strong><br/>\\s*会员状态：)([^<]+)(?:<br\\/>收到祝福：<span id=\"blessing_num\">)(\\d+)")
	g_space_regexp_story, _ = regexp.Compile("(?:story[.]php[?]id=)(\\d+)")
}

func NewSpace(request *request.Request) *Space {
	return &Space{spaceid: request.GetUid(), request: request}
}

func (this *Space) GetStoryId() int {
	return this.stordyid
}

func grepPersonInfo(sel *goquery.Selection, info *spacePersonInfo) bool {
	selBody := sel.Find("a")
	if nil == selBody {
		return false
	}

	var idStr string
	var imgStr string
	var infoStr string

	idStr, _ = selBody.Attr("href")
	bufArray := g_space_regexp_id.FindAllStringSubmatch(idStr, -1)
	if nil == bufArray {
		return false
	} else {
		info.id, _ = strconv.Atoi(bufArray[0][1])
	}

	imgStr, _ = selBody.Children().First().Attr("src")
	if "" == imgStr {
		return false
	} else {
		info.imgurl = imgStr
	}

	infoStr = selBody.Text()
	infoArray := g_space_regexp_info.FindAllStringSubmatch(infoStr, -1)
	if nil == infoArray {
		return false
	} else {
		info.name = infoArray[0][1]
		info.age, _ = strconv.Atoi(infoArray[0][2])
		info.district = infoArray[0][3]
	}

	return true
}

func (this *Space) Crawl() *Space {
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
	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	/* person info */
	sel := query.Find(".story_user > li")
	if sel.Length() > 0 {
		grepPersonInfo(sel.First(), &this.girl)
		grepPersonInfo(sel.Next(), &this.guy)

		genderA := checkAndGrepPersonInfo(this.girl.id, this.girl.district)
		genderB := checkAndGrepPersonInfo(this.guy.id, this.guy.district)

		if 1 == genderA && 0 == genderB {
			this.girl, this.guy = this.guy, this.girl
		} else if 1 == genderA && 3 == genderB {
			this.girl, this.guy = this.guy, this.girl
		} else if 3 == genderA && 0 == genderB {
			this.girl, this.guy = this.guy, this.girl
		}
	}

	/* the love status */
	sel = query.Find(".story_state")
	if sel.Length() > 0 {
		htmlStr, _ := sel.Html()
		infoArray := g_space_regexp_status.FindAllStringSubmatch(htmlStr, -1)
		if nil != infoArray {
			this.falldays, _ = strconv.Atoi(infoArray[0][1])
			this.status = infoArray[0][2]
			this.blessnum, _ = strconv.Atoi(infoArray[0][3])
		}
	}

	/* images */
	sel = query.Find("#viewerFrame input")
	if sel.Length() > 0 {
		var imgStr string

		for i, _ := range sel.Nodes {
			imgStr, _ = sel.Eq(i).Attr("value")
			this.imgs = append(this.imgs, imgStr)
		}
	}

	/* story id */
	sel = query.Find(".blog_text > dl > dd > p > a")
	if sel.Length() > 0 {
		hrefStr, _ := sel.Attr("href")
		idArray := g_space_regexp_story.FindAllStringSubmatch(hrefStr, -1)
		if nil != idArray {
			this.stordyid, _ = strconv.Atoi(idArray[0][1])
		}
	}

	this.storyinfo = NewStory(request.NewRequest(request.REQUESTURL_STORY, this.stordyid, nil)).Crawl()

	this.crawled = true
	return this
}

func (this *Space) Save() {
	if true != this.crawled {
		return
	}

	if true == dbopt.IsLoveshowExist(this.spaceid) {
		return
	}

	tag := "loveshow/" + strconv.Itoa(this.spaceid)

	/* love show basic information */
	girlImg, _ := downloadImgAndSave(this.girl.imgurl, tag)
	guyImg, _ := downloadImgAndSave(this.guy.imgurl, tag)
	dbopt.InsertToLoveshow(this.spaceid, this.storyinfo.timeUTC, this.blessnum, this.falldays, this.girl.id, this.guy.id,
		this.girl.age, this.guy.age, this.girl.name, this.guy.name, girlImg, guyImg, this.girl.district,
		this.guy.district, this.status, this.storyinfo.title, this.storyinfo.story)

	/* pictures */
	for _, p := range this.imgs {
		imgname, err := downloadImgAndSave(p, tag)
		if nil == err {
			dbopt.InsertToLoveshowpicture(this.spaceid, imgname)
		}
	}

	this.storyinfo.Save(this.spaceid)
}
