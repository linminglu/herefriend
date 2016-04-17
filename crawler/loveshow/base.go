package loveshow

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"herefriend/crawler/dbtables"
	"herefriend/crawler/page"
	"herefriend/crawler/request"
)

var g_base_regex_time *regexp.Regexp

func init() {
	g_base_regex_time, _ = regexp.Compile("(\\d+)(?:-)(\\d+)(?:-)(\\d+)(?: )(\\d+)(?::)(\\d+)(?::)(\\d+)")
}

func eascapString(s string) string {
	ss := strings.Replace(s, "百合网", "这里", -1)
	ss = strings.Replace(ss, "百合", "这里", -1)

	return ss
}

/* analyze the string and get the UTC time */
func grepTimeUTC(timeStr string) (time.Time, bool) {
	timeArray := g_base_regex_time.FindAllStringSubmatch(timeStr, -1)
	if nil != timeArray {
		var year, mon, mday, hour, min, sec int

		year, _ = strconv.Atoi(timeArray[0][1])
		mon, _ = strconv.Atoi(timeArray[0][2])
		mday, _ = strconv.Atoi(timeArray[0][3])
		hour, _ = strconv.Atoi(timeArray[0][4])
		min, _ = strconv.Atoi(timeArray[0][5])
		sec, _ = strconv.Atoi(timeArray[0][6])

		return time.Date(year, time.Month(mon), mday, hour, min, sec, 0, time.UTC), true
	} else {
		return time.Now(), false
	}
}

/* download and save it to QiNiu, return the image name */
func downloadImgAndSave(url string, tag string) (string, error) {
	if "" == url {
		return "", nil
	}

	resp, err := request.DownloadUrl(url, nil)
	if nil != err {
		return "", err
	}

	defer resp.Body.Close()

	strslice := strings.Split(url, "/")
	imgname := strslice[len(strslice)-1]
	if "" != imgname {
		err = request.PostImageToQiniu(tag+"/"+imgname, resp.Body)
	}

	return imgname, err
}

/*
 *
 *    Function: checkAndGrepPersonInfo
 *      Author: sunchao
 *        Date: 15/8/23
 * Description: grep the information of person and return the gender
 *
 */
func checkAndGrepPersonInfo(id int, addStr string) int {
	if 0 == id {
		return 3
	}

	gender, bexist := lib.CheckIsPersonExist(id)
	if true == bexist {
		return gender
	}

	province, district := dbopt.GetDistrictString(addStr)
	p := page.NewPage(request.NewRequest(request.REQUESTURL_PAGE, id, nil)).Crawl()
	dbopt.SavePersonIdInfo(id, p.GetGender(), p.GetAge(), province, district)
	p.Save()

	return p.GetGender()
}
