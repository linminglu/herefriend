package lib

import (
	"strings"

	. "herefriend/common"
)

type DistrictJson struct {
	Province string
	District []string
}

var g_provmap map[string]string
var g_distmap map[string]string
var g_DistrictJson []DistrictJson

func init() {
	/* init maps */
	g_provmap = make(map[string]string)
	g_distmap = make(map[string]string)

	var ok bool
	for _, dist := range G_DistrictB {
		_, ok = g_provmap[dist.Provcode]
		if true != ok {
			g_provmap[dist.Provcode] = dist.Province
		}

		_, ok = g_distmap[dist.Distcode]
		if true != ok {
			g_distmap[dist.Distcode] = dist.District
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
			g_DistrictJson = append(g_DistrictJson, districtJson)
			districtJson.Province = G_DistrictB[i].Province
			districtJson.District = []string{G_DistrictB[i].District}
		}
	}

	g_DistrictJson = append(g_DistrictJson, districtJson)
}

func getProvinceByCode(code string) string {
	s, _ := g_provmap[code]

	return s
}

func getDistrictByCode(code string) string {
	s, _ := g_distmap[code]

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
	return &g_DistrictJson
}
