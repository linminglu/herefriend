package main

import (
	"fmt"
	"regexp"
	"time"

	"herefriend/lib"
)

var g_districtmap map[string]([]string)
var g_regex *regexp.Regexp

func init() {
	//values
	g_RandValueMax[V_INCOME_GIRLS], g_RandValueMap[V_INCOME_GIRLS] = getRandValueResult(g_RandIncomeGirls)
	g_RandValueMax[V_INCOME_GUYS], g_RandValueMap[V_INCOME_GUYS] = getRandValueResult(g_RandIncomeGuys)
	g_RandValueMax[V_VIPLEVEL_GIRLS], g_RandValueMap[V_VIPLEVEL_GIRLS] = getRandValueResult(g_RandViplevelGirls)
	g_RandValueMax[V_VIPLEVEL_GUYS], g_RandValueMap[V_VIPLEVEL_GUYS] = getRandValueResult(g_RandViplevelGuys)
	g_RandValueMax[V_HEART_VIPLEVEL_GIRLS], g_RandValueMap[V_HEART_VIPLEVEL_GIRLS] = getRandValueResult(g_RandViplevelHeartbeatGirls)
	g_RandValueMax[V_HEART_VIPLEVEL_GUYS], g_RandValueMap[V_HEART_VIPLEVEL_GUYS] = getRandValueResult(g_RandViplevelHeartbeatGuys)

	//options
	g_RandOptionMax[O_DISTRICT_GIRLS], g_RandOptionMap[O_DISTRICT_GIRLS] = getRandOptionResult(g_RandDistrictGirls)
	g_RandOptionMax[O_DISTRICT_GUYS], g_RandOptionMap[O_DISTRICT_GUYS] = getRandOptionResult(g_RandDistrictGuys)
	g_RandOptionMax[O_OCCUPATION_GIRLS], g_RandOptionMap[O_OCCUPATION_GIRLS] = getRandOptionResult(g_RandOccupationGirls)
	g_RandOptionMax[O_OCCUPATION_GUYS], g_RandOptionMap[O_OCCUPATION_GUYS] = getRandOptionResult(g_RandOccupationGuys)
	g_RandOptionMax[O_EDUCATION_GIRLS], g_RandOptionMap[O_EDUCATION_GIRLS] = getRandOptionResult(g_RandEducationGirls)
	g_RandOptionMax[O_EDUCATION_GUYS], g_RandOptionMap[O_EDUCATION_GUYS] = getRandOptionResult(g_RandEducationGuys)
	g_RandOptionMax[O_HOUSING_GIRLS], g_RandOptionMap[O_HOUSING_GIRLS] = getRandOptionResult(g_RandHousingGirls)
	g_RandOptionMax[O_HOUSING_GUYS], g_RandOptionMap[O_HOUSING_GUYS] = getRandOptionResult(g_RandHousingGuys)
	g_RandOptionMax[O_MARRIAGE_GIRLS], g_RandOptionMap[O_MARRIAGE_GIRLS] = getRandOptionResult(g_RandMarriageGirls)
	g_RandOptionMax[O_MARRIAGE_GUYS], g_RandOptionMap[O_MARRIAGE_GUYS] = getRandOptionResult(g_RandMarriageGuys)
	g_RandOptionMax[O_CHARACTOR_GIRLS], g_RandOptionMap[O_CHARACTOR_GIRLS] = getRandOptionResult(g_RandCharactorGirls)
	g_RandOptionMax[O_CHARACTOR_GUYS], g_RandOptionMap[O_CHARACTOR_GUYS] = getRandOptionResult(g_RandCharactorGuys)
	g_RandOptionMax[O_HOBBIES_GIRLS], g_RandOptionMap[O_HOBBIES_GIRLS] = getRandOptionResult(g_RandHobbiesGirls)
	g_RandOptionMax[O_HOBBIES_GUYS], g_RandOptionMap[O_HOBBIES_GUYS] = getRandOptionResult(g_RandHobbiesGuys)

	//districts of prov
	g_districtmap = make(map[string]([]string))
	districtinfo := lib.GetDistrictJsonArray()
	for _, s := range *districtinfo {
		g_districtmap[s.Province] = s.District
	}

	g_regex = regexp.MustCompile("(.*)?(?:(现在常[住|驻]地点)([^，]*)?，)(.*)?")
}

func getRandomProvDist(gender int) (string, string) {
	index := O_DISTRICT_GIRLS
	if 1 == gender {
		index = O_DISTRICT_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	prov := g_RandOptionMap[index][r]

	r = lib.Intn(len(g_districtmap[prov]))
	dist := g_districtmap[prov][r]

	return prov, dist
}

func getRandomIncomeRange(gender int) (int, int) {
	index := V_INCOME_GIRLS
	if 1 == gender {
		index = V_INCOME_GUYS
	}

	r := lib.Intn(g_RandValueMax[index])
	node := g_RandValueMap[index][r]

	return node.value1, node.value2
}

func getRandomViplevel(gender int) (int, int) {
	index := V_VIPLEVEL_GIRLS
	if 1 == gender {
		index = V_VIPLEVEL_GUYS
	}

	r := lib.Intn(g_RandValueMax[index])
	node := g_RandValueMap[index][r]

	return node.value1, node.value2
}

func getRandomOccupation(gender int) string {
	index := O_OCCUPATION_GIRLS
	if 1 == gender {
		index = O_OCCUPATION_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	return g_RandOptionMap[index][r]
}

func getRandomEducation(gender int) string {
	index := O_EDUCATION_GIRLS
	if 1 == gender {
		index = O_EDUCATION_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	return g_RandOptionMap[index][r]
}

func getRandomHoursing(gender int) string {
	index := O_HOUSING_GIRLS
	if 1 == gender {
		index = O_HOUSING_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	return g_RandOptionMap[index][r]
}

func getRandomMarriage(gender int) string {
	index := O_MARRIAGE_GIRLS
	if 1 == gender {
		index = O_MARRIAGE_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	return g_RandOptionMap[index][r]
}

func getRandomCharactor(gender int) string {
	index := O_CHARACTOR_GIRLS
	if 1 == gender {
		index = O_CHARACTOR_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	return g_RandOptionMap[index][r]
}

func getRandomHobbies(gender int) string {
	index := O_HOBBIES_GIRLS
	if 1 == gender {
		index = O_HOBBIES_GUYS
	}

	r := lib.Intn(g_RandOptionMax[index])
	return g_RandOptionMap[index][r]
}

func updateHeartbeatByGender(gender int) {
	var sentenceSelectHeartbeatIds string
	var sentenceSelectIntro string
	var sentenceUpdate string
	var sentenceUpdateSimple string
	var sentenceUpdateHearbeat string
	var id int
	var count int
	var intro string
	var bsimple bool
	var allowincome string

	/*
	 * 遍历heartbeat girls，更新信息
	 */
	sentenceSelectHeartbeatIds = fmt.Sprintf("select id from heartbeat where flag=0 and gender=%d limit 1000", gender)
	sentenceUpdateHearbeat = fmt.Sprintf("update heartbeat set province=?,flag=1 where id=?")

	if 0 == gender {
		sentenceSelectIntro = fmt.Sprintf("select introduction from girls where id=?")
		sentenceUpdateSimple = fmt.Sprintf("update girls set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,flag=1 where id=?")
		sentenceUpdate = fmt.Sprintf("update girls set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,introduction=?,flag=1 where id=?")
	} else {
		sentenceSelectIntro = fmt.Sprintf("select introduction from guys where id=?")
		sentenceUpdateSimple = fmt.Sprintf("update guys set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,flag=1 where id=?")
		sentenceUpdate = fmt.Sprintf("update guys set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,introduction=?,flag=1 where id=?")
	}

	for {
		rows, err := lib.SQLQuery(sentenceSelectHeartbeatIds)
		if nil != err {
			panic(err.Error())
		}

		count = 0
		for rows.Next() {
			err = rows.Scan(&id)
			if nil != err {
				continue
			} else {
				count = count + 1
			}

			bsimple = true
			lib.SQLQueryRow(sentenceSelectIntro, id).Scan(&intro)
			strarray := g_regex.FindAllStringSubmatch(intro, -1)
			if nil != strarray {
				if 5 == len(strarray[0]) {
					intro = strarray[0][1] + strarray[0][4]
					bsimple = false
				}
			}

			prov, dist := getRandomProvDist(gender)
			min, max := getRandomIncomeRange(gender)
			allowmin, allowmax := getRandomIncomeRange(0)
			if 0 != allowmax {
				allowincome = fmt.Sprintf("%d-%d", allowmin, allowmax)
			} else {
				allowincome = ""
			}

			fmt.Println(intro)

			if true == bsimple {
				_, err = lib.SQLExec(sentenceUpdateSimple, prov, dist, min, max, getRandomOccupation(gender), getRandomEducation(gender), getRandomHoursing(gender), getRandomMarriage(gender), getRandomCharactor(gender), getRandomHobbies(gender), prov, getRandomEducation(1-gender), allowincome, id)
			} else {
				_, err = lib.SQLExec(sentenceUpdate, prov, dist, min, max, getRandomOccupation(gender), getRandomEducation(gender), getRandomHoursing(gender), getRandomMarriage(gender), getRandomCharactor(gender), getRandomHobbies(gender), prov, getRandomEducation(1-gender), allowincome, intro, id)
			}
			if nil != err {
				fmt.Println(err)
			}

			_, err = lib.SQLExec(sentenceUpdateHearbeat, prov, id)
			if nil != err {
				fmt.Println(err)
			}
		}

		rows.Close()
		if 0 == count {
			break
		}
	}
}

func updatePersonByGender(gender int) {
	var sentenceSelectIds string
	var sentenceSelectIntro string
	var sentenceUpdate string
	var sentenceUpdateSimple string
	var id int
	var count int
	var intro string
	var bsimple bool
	var allowincome string

	/*
	 * 遍历heartbeat girls，更新信息
	 */
	if 0 == gender {
		sentenceSelectIds = fmt.Sprintf("select id from girls where flag=0 order by id desc limit 1000")
	} else {
		sentenceSelectIds = fmt.Sprintf("select id from guys where flag=0 order by id desc limit 1000")
	}

	if 0 == gender {
		sentenceSelectIntro = fmt.Sprintf("select introduction from girls where id=?")
		sentenceUpdateSimple = fmt.Sprintf("update girls set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,flag=1 where id=?")
		sentenceUpdate = fmt.Sprintf("update girls set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,introduction=?,flag=1 where id=?")
	} else {
		sentenceSelectIntro = fmt.Sprintf("select introduction from guys where id=?")
		sentenceUpdateSimple = fmt.Sprintf("update guys set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,flag=1 where id=?")
		sentenceUpdate = fmt.Sprintf("update guys set province=?,district=?,incomemin=?,incomemax=?,occupation=?,education=?,housing=?,marriage=?,charactor=?,hobbies=?,allow_residence=?,allow_education=?,allow_income=?,introduction=?,flag=1 where id=?")
	}

	for {
		rows, err := lib.SQLQuery(sentenceSelectIds)
		if nil != err {
			panic(err.Error())
		}

		count = 0
		for rows.Next() {
			err = rows.Scan(&id)
			if nil != err {
				continue
			} else {
				count = count + 1
			}

			bsimple = true
			lib.SQLQueryRow(sentenceSelectIntro, id).Scan(&intro)
			strarray := g_regex.FindAllStringSubmatch(intro, -1)
			if nil != strarray {
				if 5 == len(strarray[0]) {
					intro = strarray[0][1] + strarray[0][4]
					bsimple = false
				}
			}

			prov, dist := getRandomProvDist(gender)
			min, max := getRandomIncomeRange(gender)
			allowmin, allowmax := getRandomIncomeRange(0)
			if 0 != allowmax {
				allowincome = fmt.Sprintf("%d-%d", allowmin, allowmax)
			} else {
				allowincome = ""
			}

			fmt.Println(intro)

			if true == bsimple {
				_, err = lib.SQLExec(sentenceUpdateSimple, prov, dist, min, max, getRandomOccupation(gender), getRandomEducation(gender), getRandomHoursing(gender), getRandomMarriage(gender), getRandomCharactor(gender), getRandomHobbies(gender), prov, getRandomEducation(1-gender), allowincome, id)
			} else {
				_, err = lib.SQLExec(sentenceUpdate, prov, dist, min, max, getRandomOccupation(gender), getRandomEducation(gender), getRandomHoursing(gender), getRandomMarriage(gender), getRandomCharactor(gender), getRandomHobbies(gender), prov, getRandomEducation(1-gender), allowincome, intro, id)
			}
			if nil != err {
				fmt.Println(err)
			}
		}

		rows.Close()
		if 0 == count {
			break
		}
	}
}

func updateVipleveByGender(gender int) {
	var sentenceSelectIds string
	var sentenceUpdate string
	var id int
	var count int

	if 0 == gender {
		sentenceSelectIds = fmt.Sprintf("select id from girls where flag=0 order by id desc limit 1000")
	} else {
		sentenceSelectIds = fmt.Sprintf("select id from guys where flag=0 order by id desc limit 1000")
	}

	if 0 == gender {
		sentenceUpdate = fmt.Sprintf("update girls set viplevel=?,vipdays=?,vipexpiretime=?,flag=1 where id=?")
	} else {
		sentenceUpdate = fmt.Sprintf("update guys set viplevel=?,vipdays=?,vipexpiretime=?,flag=1 where id=?")
	}

	for {
		rows, err := lib.SQLQuery(sentenceSelectIds)
		if nil != err {
			fmt.Println(err)
		}

		count = 0
		for rows.Next() {
			err = rows.Scan(&id)
			if nil != err {
				continue
			} else {
				count = count + 1
			}

			level, days := getRandomViplevel(gender)
			expiretime := lib.CurrentTimeUTCInt64() + int64(days)*int64(time.Hour/time.Second)*24
			_, err := lib.SQLExec(sentenceUpdate, level, days, expiretime, id)
			if nil != err {
				fmt.Println(err)
			}

			fmt.Printf("update %d vip: level=%d days=%d expiretime=%v\n", id, level, days, lib.Int64_To_UTCTime(expiretime))
		}

		rows.Close()
		if 0 == count {
			break
		}
	}
}

func main() {
	//updateHeartbeatByGender(0)
	//updateHeartbeatByGender(1)
	//updatePersonByGender(0)
	//updatePersonByGender(1)
	updateVipleveByGender(0)
	updateVipleveByGender(1)
}
