package handlers

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"herefriend/common"
	"herefriend/lib"
)

var gRegexp *regexp.Regexp

func init() {
	gRegexp, _ = regexp.Compile("(?:</code>&nbsp;来自：)([^<]+)(?:</p>)")
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

	return common.GetDistrictByString(addStr)
}

/*
 *
 *    Function: GetDistrict
 *      Author: sunchao
 *        Date: 15/7/9
 * Description: get the district information
 *
 */
func GetDistrict(c *gin.Context) {
	c.JSON(http.StatusOK, lib.GetDistrictJsonArray())
}
