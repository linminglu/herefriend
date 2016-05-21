package page

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"herefriend/common"
	"herefriend/crawler/image"
	"herefriend/crawler/request"
	"herefriend/lib"
)

type Page struct {
	usrid   int
	crawled bool
	gender  int
	age     int
	headimg string
	imgs    []string
	info    common.PersonInfo
	request *request.Request
}

var g_regex *regexp.Regexp
var g_authCookie []*http.Cookie

func init() {
	g_regex, _ = regexp.Compile("(?:\"defaultUrl\":\")([^\"]+)(?:\",)")
	g_authCookie = []*http.Cookie{
		&http.Cookie{Name: "AuthCookie", Value: "4BFFD62B611D896EC112EEBC1A9A06D86A33C350B2C6EA9BF4BE75777AA316ED0A241934D8D4622C904AF051C2A81FF4D7A6E2425F8EA8580F79DE7244C250078EF4F6523209E8F26BFB0F813A653583", Path: "/"},
	}
}

func NewPage(request *request.Request) *Page {
	if nil == request {
		return &Page{crawled: true}
	} else {
		return &Page{usrid: request.GetUid(), request: request}
	}
}

func (this *Page) IsCrawled() bool {
	return this.crawled
}

func (this *Page) GetGender() int {
	return this.gender
}

func (this *Page) GetProvince() string {
	return this.info.Province
}

func (this *Page) GetUsrId() int {
	return this.usrid
}

func (this *Page) GetAge() int {
	return this.age
}

func (this *Page) SetHeadImg(url string) {
	this.headimg = url
}

func (this *Page) SetImages(urls []string) {
	this.imgs = urls
}

func (this *Page) GetPersonInfo() *common.PersonInfo {
	return &this.info
}

func (this *Page) SetPersonInfo(info common.PersonInfo) {
	this.info = info
	this.gender = info.Gender
	this.age = info.Age
}

func eascapString(s string) string {
	ss := strings.Replace(s, "百合网", "这里", -1)
	ss = strings.Replace(ss, "百合", "这里", -1)

	return ss
}

func (this *Page) CrawlSimple() *Page {
	if true == this.crawled {
		return this
	}

	err := this.request.Download(nil)
	if nil != err {
		return this
	}

	fmt.Println("[Craw] Start")
	body := this.request.GetBody()
	if nil == body {
		return this
	}

	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	dataSel := query.Find("div.profile")
	//Name
	sel := dataSel.Find(".name > span")
	if sel.Length() > 0 {
		this.info.Name = sel.Eq(1).Text()
	}

	profileSel := query.Find(".perData")
	intrSel := profileSel.Find(".intr")
	if intrSel.Length() > 0 {
		this.info.Introduction = intrSel.Text()
	}

	this.crawled = true
	fmt.Println("[Craw] Finished")
	return this
}

/*
 |    Function: Crawl
 |      Author: Mr.Sancho
 |        Date: 2016-01-10
 |   Arguments:
 |      Return:
 | Description: the crawl action, get all the <a> and <img> from the page
 |
*/
func (this *Page) Crawl(needpic bool) *Page {
	if true == this.crawled {
		return this
	}

	err := this.request.Download(g_authCookie)
	if nil != err {
		return this
	}

	fmt.Println("【Craw】start")
	body := this.request.GetBody()
	if nil == body {
		return this
	}

	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	/* Gender */
	sel := query.Find(".tabMenu > a")
	if sel.Length() > 0 {
		if true == strings.Contains(sel.Text(), "她") {
			this.gender = 0
		} else {
			this.gender = 1
		}
	}

	this.info.Gender = this.gender

	dataSel := query.Find("div.profile")
	//Name
	sel = dataSel.Find(".name > span")
	if sel.Length() > 0 {
		this.info.Name = sel.Eq(1).Text()
	}

	//LoveStyle
	sel = dataSel.Find(".data > dl > dt")
	if sel.Length() > 0 {
		this.info.LoveType = strings.Replace(sel.Eq(1).Text(), "恋爱类型：", "", -1)
	}

	matchSel := dataSel.Find(".matching > .cont > dl")
	if matchSel.Length() > 0 {
		for i := 1; i < matchSel.Length(); i = i + 1 {
			sel = matchSel.Eq(i).Children()
			switch sel.Eq(0).Text() {
			case "年　　龄：":
				this.info.Age, _ = strconv.Atoi(strings.Replace(sel.Eq(1).Text(), "岁", "", -1))
				this.info.Allow_age = sel.Eq(3).Text()
			case "身　　高：":
				this.info.Height, _ = strconv.Atoi(strings.Replace(sel.Eq(1).Text(), "cm", "", -1))
				this.info.Allow_height = sel.Eq(3).Text()
			case "学　　历：":
				this.info.Education = sel.Eq(1).Text()
				this.info.Allow_education = sel.Eq(3).Text()
			case "月 收  入：":
				this.info.Income = sel.Eq(1).Text()
				this.info.Allow_income = sel.Eq(3).Text()
			case "婚姻状况：":
				this.info.Marriage = sel.Eq(1).Text()
				this.info.Allow_marriage = sel.Eq(3).Text()
			case "购房情况：":
				this.info.Housing = sel.Eq(1).Text()
				this.info.Allow_housing = sel.Eq(3).Text()
			case "所在地区：":
				this.info.Province, this.info.District = common.GetDistrictByString(sel.Eq(1).Text())
				this.info.Allow_residence = sel.Eq(3).Text()
			case "有无子女：":
				this.info.Allow_kidstatus = sel.Eq(3).Text()
			}
		}
	}

	profileSel := query.Find(".perData")
	intrSel := profileSel.Find(".intr")
	if intrSel.Length() > 0 {
		this.info.Introduction = intrSel.Text()
	}

	infoSel := profileSel.Find("dl")
	if infoSel.Length() > 0 {
		dlSel := infoSel.Eq(0).Children()
		for i := 0; i < dlSel.Length(); i = i + 2 {
			switch dlSel.Eq(i).Text() {
			case "民　　族：":
				this.info.Nationnality = dlSel.Eq(i + 1).Text()
			case "家　　乡：":
				this.info.Native = dlSel.Eq(i + 1).Text()
			case "属　　相：":
				this.info.Animal = dlSel.Eq(i + 1).Text()
			case "星　　座：":
				this.info.Constellation = dlSel.Eq(i + 1).Text()
			case "血　　型：":
				this.info.BloodType = dlSel.Eq(i + 1).Text()
			case "体　　型：":
				this.info.BodyType = dlSel.Eq(i + 1).Text()
			case "体　　重：":
				this.info.Weight, _ = strconv.Atoi(strings.Replace(dlSel.Eq(i+1).Text(), "公斤", "", -1))
			case "职　　业：":
				this.info.Occupation = dlSel.Eq(i + 1).Text()
			case "月　　薪：":
				this.info.Income = dlSel.Eq(i + 1).Text()
			case "购　　房：":
				this.info.Housing = dlSel.Eq(i + 1).Text()
			case "购　　车：":
				this.info.Carstatus = dlSel.Eq(i + 1).Text()
			}
		}

		dlSel = infoSel.Eq(1).Children()
		for i := 0; i < dlSel.Length(); i = i + 2 {
			switch dlSel.Eq(i).Text() {
			case "掌握语言：":
				this.info.Lang = dlSel.Eq(i + 1).Text()
			case "毕业学校：":
				this.info.School = dlSel.Eq(i + 1).Text()
			case "所学专业：":
				this.info.Speciality = dlSel.Eq(i + 1).Text()
			case "宗教信仰：":
				this.info.Religion = dlSel.Eq(i + 1).Text()
			case "相貌自评：":
				this.info.Selfjudge = dlSel.Eq(i + 1).Text()
			case "公司性质：":
				this.info.Companytype = dlSel.Eq(i + 1).Text()
			case "公司行业：":
				this.info.Companyindustry = dlSel.Eq(i + 1).Text()
			case "职　　业：":
				this.info.Occupation = dlSel.Eq(i + 1).Text()
			case "月　　薪：":
				this.info.Income = dlSel.Eq(i + 1).Text()
			case "购　　房：":
				this.info.Housing = dlSel.Eq(i + 1).Text()
			case "购　　车：":
				this.info.Carstatus = dlSel.Eq(i + 1).Text()
			}
		}
	}

	if true == needpic {
		/* crawl the pictures */
		var bodyStr = string(body)
		picnumstr := g_regex.FindAllStringSubmatch(bodyStr, -1)
		if nil != picnumstr {
			var picurl string
			for i, s := range picnumstr {
				if 0 == i {
					this.headimg = strings.Replace(s[1], "120_150", "290_290", 1)
				}

				picurl = strings.Replace(s[1], "120_150", "640_480", 1)
				this.imgs = append(this.imgs, picurl)
			}
		}
	}

	this.crawled = true
	fmt.Println("【Craw】finished")
	return this
}

/*
 |    Function: insertNewId
 |      Author: Mr.Sancho
 |        Date: 2016-01-10
 |   Arguments:
 |      Return:
 | Description: 插入新行，保证id唯一
 |
*/
var g_lastidLock sync.Mutex
var g_lastidtmp = 0

func insertNewId(gender, usertype int, info *common.PersonInfo) int {
	var lastId int

	g_lastidLock.Lock()
	if 0 == g_lastidtmp {
		girlLastIdSentence := lib.SQLSentence(lib.SQLMAP_Select_LastId, 0)
		guylLastIdSentence := lib.SQLSentence(lib.SQLMAP_Select_LastId, 1)

		var girlsLastId int
		var guysLastId int
		lib.SQLQueryRow(girlLastIdSentence).Scan(&girlsLastId)
		lib.SQLQueryRow(guylLastIdSentence).Scan(&guysLastId)

		g_lastidtmp = func() int {
			if girlsLastId > guysLastId {
				return girlsLastId + 1
			} else {
				return guysLastId + 1
			}
		}()
	} else {
		g_lastidtmp = g_lastidtmp + 1
	}

	lastId = g_lastidtmp
	g_lastidLock.Unlock()

	insertsentence := lib.SQLSentence(lib.SQLMAP_Insert_Info, gender)

	for {
		_, err := lib.SQLExec(insertsentence, lastId, "", info.Name, gender, 0, info.Age, usertype, "", info.Height, info.Weight, info.Province, info.District, info.CityLove, info.Naken)
		if nil == err {
			return lastId
		} else {
			fmt.Println(err)
		}

		g_lastidLock.Lock()
		g_lastidtmp = g_lastidtmp + 1
		lastId = g_lastidtmp
		g_lastidLock.Unlock()
	}
}

/*
 |    Function: updateUserInfo
 |      Author: Mr.Sancho
 |        Date: 2016-01-10
 |   Arguments:
 |      Return:
 | Description: 更新用户数据
 |
*/
func updateUserInfo(id, gender int, info *common.PersonInfo) error {
	sentence := lib.SQLSentence(lib.SQLMAP_Update_Info, gender)
	_, err := lib.SQLExec(sentence, info.LoveType, info.BodyType, info.Marriage, info.Province, info.District, info.Native,
		info.Education, info.Occupation, info.Housing, info.Carstatus, info.Introduction, info.School, info.Speciality, info.Animal,
		info.Constellation, info.Lang, info.BloodType, info.Selfjudge, info.Companytype, info.Companyindustry, info.Nationnality,
		info.Religion, info.Charactor, info.Hobbies, info.Allow_age, info.Allow_residence, info.Allow_height, info.Allow_marriage,
		info.Allow_education, info.Allow_housing, info.Allow_income, info.Allow_kidstatus, id)
	return err
}

/*
 |    Function: insertPictureById
 |      Author: Mr.Sancho
 |        Date: 2016-01-10
 |   Arguments:
 |      Return:
 | Description: 插入用户照片
 |
*/
func insertPictureById(id, gender int, filename string, bHead bool) {
	sentence := lib.SQLSentence(lib.SQLMAP_Insert_Picture, gender)

	if true == bHead {
		lib.SQLExec(sentence, id, filename, 1)
	} else {
		lib.SQLExec(sentence, id, filename, 0)
	}
}

/*
 * the function to output the crawl results to database
 */
func (this *Page) Save() {
	if true != this.crawled {
		return
	}

	this.info.Name = eascapString(this.info.Name)
	this.info.Introduction = eascapString(this.info.Introduction)
	if len([]rune(this.info.Introduction)) > 1024 {
		this.info.Introduction = string([]rune(this.info.Introduction)[:1024])
	}
	this.info.Selfjudge = eascapString(this.info.Selfjudge)
	if len([]rune(this.info.Selfjudge)) > 512 {
		this.info.Selfjudge = string([]rune(this.info.Selfjudge)[:510])
	}

	this.usrid = insertNewId(this.gender, common.USERTYPE_RB, &this.info)
	this.info.Id = this.usrid

	/* head image */
	if "" != this.headimg {
		imagename := lib.RandStringBytesMaskImprSrc(32) + ".jpg"
		err := image.DownloadImageAndPutToQiniu(this.headimg, true, this.usrid, imagename)
		if nil == err {
			insertPictureById(this.usrid, this.gender, imagename, true)
			this.info.IconUrl = imagename
		}
	}

	/* images */
	if 0 != len(this.imgs) {
		for _, s := range this.imgs {
			imagename := lib.RandStringBytesMaskImprSrc(32) + ".jpg"
			err := image.DownloadImageAndPutToQiniu(s, false, this.usrid, imagename)
			if nil == err {
				insertPictureById(this.usrid, this.gender, imagename, false)
				this.info.Pics = append(this.info.Pics, imagename)
			}
		}
	}

	updateUserInfo(this.usrid, this.gender, &this.info)
	jsonRlt, _ := json.Marshal(this.info)
	fmt.Println(string(jsonRlt))

	return
}
