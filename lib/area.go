package lib

import "herefriend/common"

// DistrictJSON .
type DistrictJSON struct {
	Province string
	District []string
}

var gProvMap map[string]string
var gDistMap map[string]string
var gDistrictJSON []DistrictJSON

func init() {
	/* init maps */
	gProvMap = make(map[string]string)
	gDistMap = make(map[string]string)

	var ok bool
	for _, dist := range common.CommonDistrcitInfos {
		_, ok = gProvMap[dist.Provcode]
		if true != ok {
			gProvMap[dist.Provcode] = dist.Province
		}

		_, ok = gDistMap[dist.Distcode]
		if true != ok {
			gDistMap[dist.Distcode] = dist.District
		}
	}

	var districtJSON DistrictJSON
	for i := range common.CommonDistrcitInfos {
		if "" == districtJSON.Province {
			districtJSON.Province = common.CommonDistrcitInfos[i].Province
			districtJSON.District = append(districtJSON.District, common.CommonDistrcitInfos[i].District)
		} else if common.CommonDistrcitInfos[i].Province == districtJSON.Province {
			districtJSON.District = append(districtJSON.District, common.CommonDistrcitInfos[i].District)
		} else {
			gDistrictJSON = append(gDistrictJSON, districtJSON)
			districtJSON.Province = common.CommonDistrcitInfos[i].Province
			districtJSON.District = []string{common.CommonDistrcitInfos[i].District}
		}
	}

	gDistrictJSON = append(gDistrictJSON, districtJSON)
}

func getProvinceByCode(code string) string {
	s, _ := gProvMap[code]

	return s
}

func getDistrictByCode(code string) string {
	s, _ := gDistMap[code]

	return s
}

// GetDistrictJSONArray returns district json struct
func GetDistrictJSONArray() *[]DistrictJSON {
	return &gDistrictJSON
}
