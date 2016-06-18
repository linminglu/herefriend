package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"herefriend/crawler/image"
	"herefriend/lib"
)

type cmsUserInfo struct {
	Id       int
	Age      int
	Name     string
	Img      string
	Province string
	Selected bool
	Usertype int
	VipLevel int
}

var gDistrictMap map[string]([]string)
var gRegexp *regexp.Regexp

func init() {
	//values
	gRandValueMax[V_INCOME_GIRLS], gRandValueMap[V_INCOME_GIRLS] = getRandValueResult(gRandIncomeGirls)
	gRandValueMax[V_INCOME_GUYS], gRandValueMap[V_INCOME_GUYS] = getRandValueResult(gRandIncomeGuys)
	gRandValueMax[V_VIPLEVEL_GIRLS], gRandValueMap[V_VIPLEVEL_GIRLS] = getRandValueResult(gRandViplevelGirls)
	gRandValueMax[V_VIPLEVEL_GUYS], gRandValueMap[V_VIPLEVEL_GUYS] = getRandValueResult(gRandViplevelGuys)
	gRandValueMax[V_HEART_VIPLEVEL_GIRLS], gRandValueMap[V_HEART_VIPLEVEL_GIRLS] = getRandValueResult(gRandViplevelHeartbeatGirls)
	gRandValueMax[V_HEART_VIPLEVEL_GUYS], gRandValueMap[V_HEART_VIPLEVEL_GUYS] = getRandValueResult(gRandViplevelHeartbeatGuys)

	//options
	gRandOptionMax[O_DISTRICT_GIRLS], gRandOptionMap[O_DISTRICT_GIRLS] = getRandOptionResult(gRandDistrictGirls)
	gRandOptionMax[O_DISTRICT_GUYS], gRandOptionMap[O_DISTRICT_GUYS] = getRandOptionResult(gRandDistrictGuys)
	gRandOptionMax[O_OCCUPATION_GIRLS], gRandOptionMap[O_OCCUPATION_GIRLS] = getRandOptionResult(gRandOccupationGirls)
	gRandOptionMax[O_OCCUPATION_GUYS], gRandOptionMap[O_OCCUPATION_GUYS] = getRandOptionResult(gRandOccupationGuys)
	gRandOptionMax[O_EDUCATION_GIRLS], gRandOptionMap[O_EDUCATION_GIRLS] = getRandOptionResult(gRandEducationGirls)
	gRandOptionMax[O_EDUCATION_GUYS], gRandOptionMap[O_EDUCATION_GUYS] = getRandOptionResult(gRandEducationGuys)
	gRandOptionMax[O_HOUSING_GIRLS], gRandOptionMap[O_HOUSING_GIRLS] = getRandOptionResult(gRandHousingGirls)
	gRandOptionMax[O_HOUSING_GUYS], gRandOptionMap[O_HOUSING_GUYS] = getRandOptionResult(gRandHousingGuys)
	gRandOptionMax[O_MARRIAGE_GIRLS], gRandOptionMap[O_MARRIAGE_GIRLS] = getRandOptionResult(gRandMarriageGirls)
	gRandOptionMax[O_MARRIAGE_GUYS], gRandOptionMap[O_MARRIAGE_GUYS] = getRandOptionResult(gRandMarriageGuys)
	gRandOptionMax[O_CHARACTOR_GIRLS], gRandOptionMap[O_CHARACTOR_GIRLS] = getRandOptionResult(gRandCharactorGirls)
	gRandOptionMax[O_CHARACTOR_GUYS], gRandOptionMap[O_CHARACTOR_GUYS] = getRandOptionResult(gRandCharactorGuys)
	gRandOptionMax[O_HOBBIES_GIRLS], gRandOptionMap[O_HOBBIES_GIRLS] = getRandOptionResult(gRandHobbiesGirls)
	gRandOptionMax[O_HOBBIES_GUYS], gRandOptionMap[O_HOBBIES_GUYS] = getRandOptionResult(gRandHobbiesGuys)

	//districts of prov
	gDistrictMap = make(map[string]([]string))
	districtinfo := lib.GetDistrictJsonArray()
	for _, s := range *districtinfo {
		gDistrictMap[s.Province] = s.District
	}

	gRegexp = regexp.MustCompile("(.*)?(?:(现在常[住|驻]地点)([^，]*)?，)(.*)?")
}

func getRandomProvDist(gender int) (string, string) {
	index := O_DISTRICT_GIRLS
	if 1 == gender {
		index = O_DISTRICT_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	prov := gRandOptionMap[index][r]

	r = lib.Intn(len(gDistrictMap[prov]))
	dist := gDistrictMap[prov][r]

	return prov, dist
}

func getRandomIncomeRange(gender int) (int, int) {
	index := V_INCOME_GIRLS
	if 1 == gender {
		index = V_INCOME_GUYS
	}

	r := lib.Intn(gRandValueMax[index])
	node := gRandValueMap[index][r]

	return node.value1, node.value2
}

func getRandomViplevel(gender int) (int, int) {
	index := V_VIPLEVEL_GIRLS
	if 1 == gender {
		index = V_VIPLEVEL_GUYS
	}

	r := lib.Intn(gRandValueMax[index])
	node := gRandValueMap[index][r]

	return node.value1, node.value2
}

func getRandomOccupation(gender int) string {
	index := O_OCCUPATION_GIRLS
	if 1 == gender {
		index = O_OCCUPATION_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	return gRandOptionMap[index][r]
}

func getRandomEducation(gender int) string {
	index := O_EDUCATION_GIRLS
	if 1 == gender {
		index = O_EDUCATION_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	return gRandOptionMap[index][r]
}

func getRandomHoursing(gender int) string {
	index := O_HOUSING_GIRLS
	if 1 == gender {
		index = O_HOUSING_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	return gRandOptionMap[index][r]
}

func getRandomMarriage(gender int) string {
	index := O_MARRIAGE_GIRLS
	if 1 == gender {
		index = O_MARRIAGE_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	return gRandOptionMap[index][r]
}

func getRandomCharactor(gender int) string {
	index := O_CHARACTOR_GIRLS
	if 1 == gender {
		index = O_CHARACTOR_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	return gRandOptionMap[index][r]
}

func getRandomHobbies(gender int) string {
	index := O_HOBBIES_GIRLS
	if 1 == gender {
		index = O_HOBBIES_GUYS
	}

	r := lib.Intn(gRandOptionMax[index])
	return gRandOptionMap[index][r]
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
			strarray := gRegexp.FindAllStringSubmatch(intro, -1)
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
			strarray := gRegexp.FindAllStringSubmatch(intro, -1)
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

func modifygifttime() {
	curtime := lib.CurrentTimeUTCInt64()
	basetimelength := int64(3600 * 24 * 365)

	curid := 0
	sentence := "select id from giftconsume where id>? limit 10000"
	updatesentence := "update giftconsume set time=? where id=?"
	for {
		rows, err := lib.SQLQuery(sentence, curid)
		if nil != err {
			continue
		}

		count := 0
		for rows.Next() {
			count = count + 1
			rows.Scan(&curid)
			randtime := curtime - lib.Int63n(basetimelength)
			fmt.Printf("%s %d %d\r\n", updatesentence, randtime, curid)
			lib.SQLExec(updatesentence, randtime, curid)
		}

		rows.Close()

		if 0 == count {
			break
		}
	}
}

func deleteHeadPic(gender int) {
	var sentence string
	var updatesentence string
	if 0 == gender {
		sentence = "select id from girls where usertype=3 and flag=0 and id>? limit 1000"
		updatesentence = "update girls set flag=1 where id=?"
	} else {
		sentence = "select id from guys where usertype=3 and flag=0 and id>? limit 1000"
		updatesentence = "update guys set flag=1 where id=?"
	}

	curid := 0
	for {
		rows, err := lib.SQLQuery(sentence, curid)
		if nil != err {
			continue
		}

		count := 0
		for rows.Next() {
			count = count + 1
			rows.Scan(&curid)

			fmt.Printf("http://localhost:8080/cms/GetSingleUserInfo?id=%d&gender=%d\n", curid, gender)
			bytes, err := lib.GetResultByMethod("", fmt.Sprintf("http://localhost:8080/cms/GetSingleUserInfo?id=%d&gender=%d", curid, gender), nil)
			if nil != err {
				fmt.Println(err)
				continue
			}

			var info cmsUserInfo
			err = json.Unmarshal(bytes, &info)
			if "" == info.Img {
				lib.GetResultByMethod("", fmt.Sprintf("http://localhost:8080/cms/AddBlacklist?id=%d&gender=%d", curid, gender), nil)
				continue
			}

			err, _, h := image.GetImageWidthHight(info.Img)
			if nil != err {
				fmt.Println(err)
				continue
			} else if h >= 80 {
				continue
			}

			lib.GetResultByMethod("", fmt.Sprintf("http://localhost:8080/cms/DeleteHeadImage?id=%d&gender=%d", curid, gender), nil)
			bytes, err = lib.GetResultByMethod("", fmt.Sprintf("http://localhost:8080/cms/GetSingleUserInfo?id=%d&gender=%d", curid, gender), nil)
			if nil == err {
				lib.SQLExec(updatesentence, curid)
				err = json.Unmarshal(bytes, &info)
				if nil == err {
					if "" == info.Img {
						lib.GetResultByMethod("", fmt.Sprintf("http://localhost:8080/cms/AddBlacklist?id=%d&gender=%d", curid, gender), nil)
					}
				} else {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
			}
		}

		rows.Close()
		if 0 == count {
			break
		}
	}
}

func main() {
	deleteHeadPic(0)
	deleteHeadPic(1)
}
