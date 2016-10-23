package filter

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"herefriend/common"
	"herefriend/config"
	"herefriend/crawler/page"
	"herefriend/crawler/request"
	"herefriend/lib"
)

type filterInfo struct {
	field   string
	content string
}

var g_tables = []string{"girls", "guys"}

var g_filters = []filterInfo{
	{"name", "%会员%"},
	{"introduction", "%我是3G这里%"},
	{"introduction", "%大家好，我是%，我平时的%"},
	{"introduction", "%等待，只为与你相遇%"},
	{"introduction", "%这里会员%"},
}

func eascapString(s string) string {
	ss := strings.Replace(s, "百合网", "这里", -1)
	ss = strings.Replace(ss, "百合", "这里", -1)
	return ss
}

func checkPersonInfo(info *common.PersonInfo) (string, string, bool) {
	name := eascapString(info.Name)
	intro := eascapString(info.Introduction)
	checkresult := true

	name = strings.TrimSpace(name)
	intro = strings.TrimSpace(intro)

	buf, _ := json.Marshal(info)
	fmt.Println(string(buf))

	if 0 == len(name) || strings.Contains(name, "会员") {
		checkresult = false
	} else if strings.Contains(intro, `我是3G这里`) {
		checkresult = false
	} else if strings.Contains(intro, `大家好，我是`) && strings.Contains(intro, `，我平时的`) {
		checkresult = false
	} else if strings.Contains(intro, `等待，只为与你相遇`) {
		checkresult = false
	}

	return name, intro, checkresult
}

func getAgeRange(age, scan int) (min, max int) {
	min = age - scan
	max = age + scan

	if config.ConfAgeMin > min {
		min = config.ConfAgeMin
	}

	if config.ConfAgeMax < max {
		max = config.ConfAgeMax
	}

	return
}

func doReplace(id, gender, age, scan int, province string) bool {
	var crawid int
	var crawcount int
	var err error
	var sentence string

	if "" == province {
		fmt.Printf("\n【Replace】id=%d age=%d gender=%d\n", id, age, gender)
		sentence = fmt.Sprintf("select id from %sid where age>=? and age<=? limit 1", g_tables[gender])
	} else {
		fmt.Printf("\n【Replace】id=%d age=%d gender=%d province=%s\n", id, age, gender, province)
		sentence = fmt.Sprintf("select id from %sid where age>=? and age<=? and province=? limit 1", g_tables[gender])
	}

	min, max := getAgeRange(age, scan)
	delsentence := fmt.Sprintf("delete from %sid where id=?", g_tables[gender])
	updatesentence := fmt.Sprintf("update %s set name=?,introduction=? where id=?", g_tables[gender])

	for {
		if "" == province {
			err = lib.SQLQueryRow(sentence, min, max).Scan(&crawid)
		} else {
			err = lib.SQLQueryRow(sentence, min, max, province).Scan(&crawid)
		}
		if nil != err {
			fmt.Println(err)
			return false
		} else {
			lib.SQLExec(delsentence, crawid)
		}

		//craw the girl info and use the pictures instead original pcitures
		fmt.Printf("[Start to craw user] crawid: %d\n", crawid)
		req := request.NewRequestBH(request.REQUESTURL_PAGE, crawid, nil)
		pageuser := page.NewPage(req)

		crawcount = 0
		for {
			crawcount++
			pageuser.CrawlSimple()
			if true != pageuser.IsCrawled() {
				if 20 > crawcount {
					time.Sleep(time.Millisecond * 100)
					continue
				}
			}

			break
		}

		if 20 <= crawcount {
			continue
		}

		// check if this user is ok
		info := pageuser.GetPersonInfo()
		name, introduction, ok := checkPersonInfo(info)
		if true != ok {
			continue
		} else {
			lib.SQLExec(updatesentence, name, introduction, id)
			lib.Get(fmt.Sprintf("http://localhost:8080/cms/RefreshUserInfo?id=%d", id), nil)
			break
		}
	}

	return true
}

func replaceInfoByGender(gender int) {
	sentence := "select id, age, province from " + g_tables[gender] + " where usertype !=1 and ("
	for i, f := range g_filters {
		if 0 != i {
			sentence += " or "
		}
		sentence += fmt.Sprintf("(%s like \"%s\")", f.field, f.content)
	}
	sentence += ") "

	lastid := 999999999
	count := 0
	id := 0
	age := 0
	province := ""

	for {
		count = 0
		rows, err := lib.SQLQuery(sentence + fmt.Sprintf("and id<%d order by id desc limit 1000", lastid))
		if nil == err {
			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&id, &age, &province)
				if nil == err {
					count++
					lastid = id

					if false == doReplace(id, gender, age, 5, province) {
						if false == doReplace(id, gender, age, 10, province) {
							doReplace(id, gender, age, 5, "")
						}
					}
				}
			}
		}

		if 0 == count {
			break
		}
	}
}

func Start() {
	replaceInfoByGender(0)
	replaceInfoByGender(1)
}
