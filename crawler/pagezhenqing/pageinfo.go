package pagezhenqing

import (
	"herefriend/common"
	"herefriend/crawler/request"
	"strconv"
	"strings"
)

type PageInfo struct {
	Headimg string
	crawled bool
	Info    common.PersonInfo
	request *request.Request
}

func NewPageInfo(request *request.Request) *PageInfo {
	return &PageInfo{request: request}
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

	/* headpic */
	sel := query.Find(".left_content > .pic > img")
	if sel.Length() > 0 {
		href, exist := sel.Attr("src")
		if true == exist {
			this.Headimg = strings.TrimSpace(href)
		}
	}

	// name
	sel = query.Find(".right_content > h2")
	if sel.Length() > 0 {
		namestr := sel.Text()
		namestr = strings.Split(namestr, "[ID:")[0]
		this.Info.Name = strings.TrimSpace(namestr)
	}

	// gender age ...
	sel = query.Find(".right_content > p > span")
	if sel.Length() > 0 {
		for i := 0; i < sel.Length(); i = i + 1 {
			texts := strings.SplitN(sel.Eq(i).Text(), "：", 2)
			AnalysisVerbose(&this.Info, texts[0], texts[1])
		}
	}

	sels := query.Find(".box_bottom > .box_content")
	if sels.Length() > 0 {
		for i := 0; i < sels.Length(); i = i + 1 {
			psel := sels.Eq(i).Find("h2")
			if psel.Length() > 0 {
				switch psel.Eq(0).Text() {
				case "择偶条件":
					sel = sels.Eq(i).Find("p > span")
					if sel.Length() > 0 {
						for i := 0; i < sel.Length(); i = i + 1 {
							texts := strings.SplitN(sel.Eq(i).Text(), "：", 2)
							AnalysisAllowVerbose(&this.Info, texts[0], texts[1])
						}
					}
				case "Ta的内心独白":
					sel = sels.Eq(i).Find("p")
					if sel.Length() > 0 {
						introcution := strings.TrimSpace(sel.Eq(0).Text())
						this.Info.Introduction = strings.Replace(introcution, "&nbsp;", "", -1)
					}
				case "Ta的详细资料":
					sel = sels.Eq(i).Find("p")
					if sel.Length() > 0 {
						hobbies := strings.SplitN(sel.Eq(0).Text(), "：", 2)[1]
						if "保密" != hobbies {
							this.Info.Hobbies = strings.TrimSpace(hobbies)
						}
					}

					sel = sels.Eq(i).Find("ul > li")
					if sel.Length() > 0 {
						for i := 0; i < sel.Length(); i = i + 1 {
							texts := strings.SplitN(sel.Eq(i).Text(), "：", 2)
							AnalysisVerbose(&this.Info, texts[0], texts[1])
						}
					}
				}
			}
		}

	}

	this.crawled = true
	return this
}

func AnalysisVerbose(info *common.PersonInfo, title, content string) {
	if "保密" == content || "不限" == content {
		return
	}

	switch title {
	case "性别":
		info.Gender = func() int {
			if "女" == content {
				return 0
			} else {
				return 1
			}
		}()
	case "年龄":
		info.Age, _ = strconv.Atoi(strings.TrimRight(content, "岁"))
	case "学历":
		info.Education = content
	case "身高":
		info.Height, _ = strconv.Atoi(strings.TrimRight(content, "cm"))
	case "月薪":
		info.Income = content
		incomes := strings.Split(content, "-")
		if 2 == len(incomes) {
			info.IncomeMin, _ = strconv.Atoi(incomes[0])
			info.IncomeMax, _ = strconv.Atoi(incomes[1])
		}
	case "籍贯":
		info.Native = content
	case "婚姻状况":
		info.Marriage = content
	case "所在地":
		addrs := strings.Split(content, "-")
		info.Province, info.District = common.GetDistrictByString(strings.TrimSpace(addrs[0]) + " " + strings.TrimSpace(addrs[1]))
	case "体    重":
		info.Weight, _ = strconv.Atoi(strings.TrimRight(content, "kg"))
	case "体    型":
		info.BodyType = content
	case "相    貌":
		info.Selfjudge = content
	case "生    肖":
		info.Animal = content
	case "星    座":
		info.Constellation = strings.Split(content, "(")[0]
	case "血    型":
		info.BloodType = content
	case "民    族":
		info.Nationnality = content
	case "信    仰":
		info.Religion = content
	case "个    性":
		info.Charactor = content
	case "毕业院校":
		info.School = content
	case "公司类型":
		info.Companytype = content
	case "公司名称":
		info.Companyindustry = content
	case "职业类别":
		info.Occupation = content
	}
}

func AnalysisAllowVerbose(info *common.PersonInfo, title, content string) {
	if "保密" == content || "不限" == content {
		return
	}

	switch title {
	case "年    龄":
		info.Allow_age = content
	case "身    高":
		info.Allow_height = content
	case "学    历":
		info.Allow_education = content
	case "婚姻状况":
		info.Allow_marriage = content
	case "是否有孩子":
		info.Allow_kidstatus = content
	case "购房情况":
		info.Allow_housing = content
	case "所在地":
		info.Allow_residence = content
	}
}
