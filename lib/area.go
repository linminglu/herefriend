package lib

import (
	"strings"

	. "herefriend/common"
)

type DistrictJson struct {
	Province string
	District []string
}

var gProvMap map[string]string
var gDistMap map[string]string
var gDistrictJson []DistrictJson

func init() {
	/* init maps */
	gProvMap = make(map[string]string)
	gDistMap = make(map[string]string)

	var ok bool
	for _, dist := range G_DistrictB {
		_, ok = gProvMap[dist.Provcode]
		if true != ok {
			gProvMap[dist.Provcode] = dist.Province
		}

		_, ok = gDistMap[dist.Distcode]
		if true != ok {
			gDistMap[dist.Distcode] = dist.District
		}
	}

	var districtJson DistrictJson
	for i := range G_DistrictB {
		if "" == districtJson.Province {
			districtJson.Province = G_DistrictB[i].Province
			districtJson.District = append(districtJson.District, G_DistrictB[i].District)
		} else if G_DistrictB[i].Province == districtJson.Province {
			districtJson.District = append(districtJson.District, G_DistrictB[i].District)
		} else {
			gDistrictJson = append(gDistrictJson, districtJson)
			districtJson.Province = G_DistrictB[i].Province
			districtJson.District = []string{G_DistrictB[i].District}
		}
	}

	gDistrictJson = append(gDistrictJson, districtJson)
}

func getProvinceByCode(code string) string {
	s, _ := gProvMap[code]

	return s
}

func getDistrictByCode(code string) string {
	s, _ := gDistMap[code]

	return s
}

func GetDistrictString(addStr string) (string, string) {
	var provcode string
	var distcode string

	if "" != addStr {
		for _, s := range G_DistrictA {
			if strings.Contains(addStr, s.Province) {
				provcode = s.Provcode
				break
			}
		}

		if "" != provcode {
			for _, s := range G_DistrictA {
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
		province = getProvinceByCode(provcode)
	}

	if "" != distcode {
		district = getDistrictByCode(distcode)
	}

	return province, district
}

func GetDistrictJsonArray() *[]DistrictJson {
	return &gDistrictJson
}
