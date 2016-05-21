package lib

import . "herefriend/common"

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
	for _, dist := range CommonDistrcitInfos {
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
	for i := range CommonDistrcitInfos {
		if "" == districtJson.Province {
			districtJson.Province = CommonDistrcitInfos[i].Province
			districtJson.District = append(districtJson.District, CommonDistrcitInfos[i].District)
		} else if CommonDistrcitInfos[i].Province == districtJson.Province {
			districtJson.District = append(districtJson.District, CommonDistrcitInfos[i].District)
		} else {
			gDistrictJson = append(gDistrictJson, districtJson)
			districtJson.Province = CommonDistrcitInfos[i].Province
			districtJson.District = []string{CommonDistrcitInfos[i].District}
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

func GetDistrictJsonArray() *[]DistrictJson {
	return &gDistrictJson
}
