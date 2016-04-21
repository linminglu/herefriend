package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"herefriend/common"
	"herefriend/lib"
)

var gRegexp *regexp.Regexp

func init() {
	gRegexp, _ = regexp.Compile("(?:</code>&nbsp;来自：)([^<]+)(?:</p>)")
}

func getDistrictString(addStr string) (string, string) {
	var provcode string
	var distcode string

	if "" != addStr {
		for _, s := range common.G_DistrictA {
			if strings.Contains(addStr, s.Province) {
				provcode = s.Provcode
				break
			}
		}

		if "" != provcode {
			for _, s := range common.G_DistrictA {
				if (provcode == s.Provcode) && (strings.Contains(addStr, s.District)) {
					distcode = s.Distcode
					break
				}
			}
		}
	}

	var province string
	var district string

	if "" != provcode {
		for _, s := range common.G_DistrictB {
			if provcode == s.Provcode {
				province = s.Province
				break
			}
		}
	}

	if "" != distcode {
		for _, s := range common.G_DistrictB {
			if distcode == s.Distcode {
				district = s.District
				break
			}
		}
	}

	if "" == province {
		tmp := strings.Split(addStr, " ")

		province = tmp[0]
		if len(tmp) > 1 {
			district = tmp[1]
		}
	}

	return province, district
}

/*
 *
 *    Function: GetIpAddress
 *      Author: sunchao
 *        Date: 15/7/11
 * Description: 根据IP地址获取地址信息
 *
 */
func GetIpAddress(r *http.Request) (string, string) {
	ipaddr := r.Header.Get("x-forwarded-for")
	if "" == ipaddr || "unknown" == strings.ToLower(ipaddr) {
		ipaddr = r.Header.Get("Proxy-Client-IP")
		if "" == ipaddr || "unknown" == strings.ToLower(ipaddr) {
			ipaddr = r.Header.Get("WL-Proxy-Client-IP")
			if "" == ipaddr || "unknown" == strings.ToLower(ipaddr) {
				ipaddr = r.RemoteAddr
			}
		}
	}

	ipstrs := strings.Split(ipaddr, ":")
	if 2 != len(ipstrs) {
		return "广东省", "东莞市"
	}

	addStr := ""
	buf, _ := lib.GetResultByMethod("GET", "http://www.ip.cn/index.php?ip="+ipstrs[0], nil)
	descStr := gRegexp.FindStringSubmatch(string(buf))
	if nil != descStr {
		addStr = descStr[1]
	}

	if "" == addStr {
		return "广东省", "东莞市"
	}

	return getDistrictString(addStr)
}

/*
 *
 *    Function: GetDistrict
 *      Author: sunchao
 *        Date: 15/7/9
 * Description: get the district information
 *
 */
func GetDistrict() (int, string) {
	jsonRlt, _ := json.Marshal(lib.GetDistrictJsonArray())
	return 200, string(jsonRlt)
}
