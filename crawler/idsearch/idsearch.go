package idsearch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"

	"herefriend/lib"
)

const AGE_MIN = 18
const AGE_MAX = 85
const ANIMAL_MIN = 1
const ANIMAL_MAX = 12
const ASTROLOGY_MIN = 1
const ASTROLOGY_MAX = 12
const HEIGHT_MIN = 140
const HEIGHT_MAX = 190
const MAXSEARCHNUM = 300

var g_IdStrRegex *regexp.Regexp
var g_IdsRegex *regexp.Regexp

func init() {
	g_IdStrRegex, _ = regexp.Compile("(?:name=\"userIds\" value=\")([^\"]+)(?:\"/>)")
	g_IdsRegex, _ = regexp.Compile("(\\d+)(?::[^,]+,?)")
}

func clearCheckMap(check map[string]bool) {
	for s, _ := range check {
		delete(check, s)
	}

	return
}

func searchByUrlValues(checkMap *map[string]bool, values url.Values) ([]string, int) {
	//Get Ids
	resp, err := http.PostForm("http://search.baihe.com/solrAdvanceSearch", values)
	if err != nil {
		return nil, 0
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	/* find the ids */
	idstr := g_IdStrRegex.FindAllStringSubmatch(string(body), -1)

	var ids []string
	var iNum int = 0

	if nil != idstr {
		idstr := g_IdsRegex.FindAllStringSubmatch(idstr[0][1], -1)
		if nil != idstr {
			for _, id := range idstr {
				_, ok := (*checkMap)[id[1]]
				if true != ok {
					ids = append(ids, id[1])
					(*checkMap)[id[1]] = true
					iNum++
				}
			}
		}

		return ids, iNum
	}

	return nil, 0
}

func doSearch(checkMap *map[string]bool, provcode, distcode, animalsign, astrology string, gender, age int) ([]string, int) {
	agestr := strconv.Itoa(age)
	genderStr := strconv.Itoa(gender)

	values := make(url.Values)
	values.Set("search.se1.countryId", "86")
	values.Set("search.se1.provinceId", provcode)
	values.Set("search.se1.districtId", distcode)
	values.Set("search.sb.animalSign", animalsign)
	values.Set("search.sb.astrology", astrology)
	values.Set("search.sb.gender", genderStr)
	values.Set("search.sb.minage", agestr)
	values.Set("search.sb.maxage", agestr)

	return searchByUrlValues(checkMap, values)
}

func doSearchWithHeight(checkMap *map[string]bool, provcode, distcode, animalsign, astrology string, gender, age, height int) ([]string, int) {
	agestr := strconv.Itoa(age)
	genderStr := strconv.Itoa(gender)
	minheight := strconv.Itoa(height)
	maxheight := strconv.Itoa(height + 10)

	values := make(url.Values)
	values.Set("search.se1.countryId", "86")
	values.Set("search.se1.provinceId", provcode)
	values.Set("search.se1.districtId", distcode)
	values.Set("search.sb.animalSign", animalsign)
	values.Set("search.sb.astrology", astrology)
	values.Set("search.sb.gender", genderStr)
	values.Set("search.sb.minage", agestr)
	values.Set("search.sb.maxage", agestr)
	values.Set("search.sb.minheight", minheight)
	values.Set("search.sb.maxheight", maxheight)

	return searchByUrlValues(checkMap, values)
}

/*
 * the function search ids deeper to astrology and height
 */
func searchByCodeDeeper(checkMap *map[string]bool, provcode, distcode string, gender, age, animal int) {
	var ids []string
	var num int
	var animalstr string
	var astrologystr string

	animalstr = strconv.Itoa(animal)

	//search ids by astrology
	for astrology := ASTROLOGY_MIN; astrology <= ASTROLOGY_MAX; astrology++ {
		astrologystr = strconv.Itoa(astrology)

		ids, num = doSearch(checkMap, provcode, distcode, animalstr, astrologystr, gender, age)
		if 0 != num {
			dbopt.SavePersonIds(ids, provcode, distcode, gender, age)

			//need search deeper to height
			if MAXSEARCHNUM == num {
				for height := HEIGHT_MIN; height <= HEIGHT_MAX; height += 10 {
					ids, num = doSearchWithHeight(checkMap, provcode, distcode, animalstr, astrologystr, gender, age, height)
					if 0 != num {
						dbopt.SavePersonIds(ids, provcode, distcode, gender, age)
					}
				}
			}
		}
	}

	return
}

/*
 * the function search ids by province code and district code
 */
func searchByCode(checkMap *map[string]bool, provcode, distcode string, gender int) {
	var ids []string
	var num int

	for age := AGE_MIN; age <= AGE_MAX; age++ {
		/* search ids by age and animal */
		for animal := ANIMAL_MIN; animal <= ANIMAL_MAX; animal++ {
			ids, num = doSearch(checkMap, provcode, distcode, strconv.Itoa(animal), "0", gender, age)
			if 0 != num {
				dbopt.SavePersonIds(ids, provcode, distcode, gender, age)

				/* need search deep to astroloy */
				if MAXSEARCHNUM == num {
					searchByCodeDeeper(checkMap, provcode, distcode, gender, age, animal)
				}
			}
		}

		clearCheckMap(*checkMap)
	}

	return
}

var g_Mutex sync.Mutex

/*
 * the function to search all the ids
 */
func idsSearch(bSearchGirl bool, c chan int) {
	var provcode, distcode string
	var err error
	var status dbopt.SearchStatus
	var gender int = 0

	if false == bSearchGirl {
		gender = 1
	}

	checkMap := make(map[string]bool)

	for {
		provcode, distcode, err = dbopt.GetNextDistrict(provcode, distcode)
		if nil == err {
			g_Mutex.Lock()
			status, err = dbopt.GetSearchStatus(provcode, distcode)
			if nil == err && dbopt.StatusNone == status {
				fmt.Printf("Searching [%s %s]\r\n", dbopt.GetProvinceByCode(provcode), dbopt.GetDistrictByCode(distcode))
				dbopt.SetSearchStatus(provcode, distcode, dbopt.StatusSearching)

				g_Mutex.Unlock()
				searchByCode(&checkMap, provcode, distcode, gender)
				g_Mutex.Lock()

				dbopt.SetSearchStatus(provcode, distcode, dbopt.StatusFinish)
			}

			g_Mutex.Unlock()
		} else {
			break
		}
	}

	c <- 1
	return
}

const g_RoutineNum = 5

func SearchGirls() {
	c := make(chan int, g_RoutineNum)

	for i := 0; i < g_RoutineNum; i++ {
		go idsSearch(true, c)
	}

	for i := 0; i < g_RoutineNum; i++ {
		<-c
	}

	fmt.Println("Girls ids search finished.")

	return
}

func SearchGuyes() {
	c := make(chan int, g_RoutineNum)

	for i := 0; i < g_RoutineNum; i++ {
		go idsSearch(false, c)
	}

	for i := 0; i < g_RoutineNum; i++ {
		<-c
	}

	fmt.Println("Guys ids search finished.")
}
