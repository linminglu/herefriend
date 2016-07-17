package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"herefriend/common"

	"github.com/garyburd/redigo/redis"
)

const (
	//userinfo_${id}
	REDIS_PREFIX_USERINFO = "userinfo_%d"
	//gender_${id}
	REDIS_PREFIX_USERGENDER = "gender_%d"
	//district_${gender}
	REDIS_PREFIX_DISTRICT = "district_%d"
	//provcount_${province}_${gender}
	REDIS_PREFIX_PROVCOUNT = "provcount_%s_%d"
	//heartbeat_provcount_${province}_${gender}
	REDIS_PREFIX_HEARTBEAT_PROVCOUNT = "heartbeat_provcount_%s_%d"
	//provcountage_${province}_${gender}_${age}
	REDIS_PREFIX_PROVAGECOUNT = "provcountage_%s_%d_%d"
	//search_${id}_${provicne}_${agemin-max}_${heightmin-max}_${incomemin-max}_${education}_${occupation}_${status}
	REDIS_PREFIX_SEARCHINDEX = "search_%d_%s_%d-%d_%d-%d_%d-%d_%s_%s_%s"
	//searchbase_${provicne}_${agemin-max}_${heightmin-max}_${incomemin-max}_${education}_${occupation}_${status}
	REDIS_PREFIX_SEARCHBASE = "searchbase_%s_%d-%d_%d-%d_%d-%d_%s_%s_%s"
	//goldbeans_${id}
	REDIS_PREFIX_GOLDBENS = "goldbeans_%d"
	//gift_recvlist_${id}
	REDIS_PREFIX_GIFT_RECVLIST = "gift_recvlist_%d"
	//gift_sendlist_${id}
	REDIS_PREFIX_GIFT_SENDLIST = "gift_sendlist_%d"
	//charm_toplist_${gender}_${year}_${month}_${day} (the charm list for last weak)
	REDIS_PREFIX_CHARM_TOPLIST = "charm_toplist_%d_%4d_%2d_%2d"
	//wealth_toplist_${year}_${month}_${day} (the wealth list for last weak)
	REDIS_PREFIX_WEALTH_TOPLIST = "wealth_toplist_%d_%4d_%2d_%2d"
)

var gRedisPool *redis.Pool

func init() {
	gRedisPool = &redis.Pool{
		MaxIdle: 5,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

/*
 |    Function: setRedisKeyValue
 |      Author: Mr.Sancho
 |        Date: 2016-05-22
 | Description:
 |      Return:
 |
*/
func setRedisValue(key string, v interface{}) {
	var content []byte

	switch v.(type) {
	case int:
		content = []byte(strconv.Itoa(v.(int)))
	case string:
		content = []byte(v.(string))
	case *common.PersonInfo:
		content, _ = json.Marshal(v.(*common.PersonInfo))
	case *[]common.GiftSendRecvInfo:
		content, _ = json.Marshal(v.(*[]common.GiftSendRecvInfo))
	case *[]common.UserCharmInfo:
		content, _ = json.Marshal(v.(*[]common.UserCharmInfo))
	case *[]common.UserWealthInfo:
		content, _ = json.Marshal(v.(*[]common.UserWealthInfo))
	}

	c := gRedisPool.Get()
	defer c.Close()

	c.Do("Set", key, content)
}

/*
 |    Function: getRedisValue
 |      Author: Mr.Sancho
 |        Date: 2016-05-22
 | Description:
 |      Return:
 |
*/
func getRedisValue(key string) ([]byte, bool) {
	c := gRedisPool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return nil, false
	}

	return content.([]byte), true
}

/*
 |    Function: delRedisValue
 |      Author: Mr.Sancho
 |        Date: 2016-05-22
 | Description:
 |      Return:
 |
*/
func delRedisValue(key string) {
	c := gRedisPool.Get()
	defer c.Close()

	c.Do("Del", key)
}

/*
 * user info
 */
func SetRedisUserInfo(id int, info *common.PersonInfo) {
	key := fmt.Sprintf(REDIS_PREFIX_USERINFO, id)
	setRedisValue(key, info)
}

func GetRedisUserInfo(id int) (*common.PersonInfo, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_USERINFO, id)
	value, exist := getRedisValue(key)
	if true == exist {
		var info common.PersonInfo
		err := json.Unmarshal(value, &info)
		if nil == err {
			return &info, true
		}
	}

	return nil, false
}

func DelRedisUserInfo(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_USERINFO, id)
	delRedisValue(key)
}

/*
 * gold beans
 */
func SetRedisGoldBeans(id, beans int) {
	key := fmt.Sprintf(REDIS_PREFIX_GOLDBENS, id)
	setRedisValue(key, beans)
}

func GetRedisGoldBeans(id int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_GOLDBENS, id)
	value, exist := getRedisValue(key)
	if true == exist {
		beanstr := string(value)
		beans, _ := strconv.Atoi(beanstr)
		return beans, true
	}

	return 0, false
}

func DelRedisGoldBeans(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_GOLDBENS, id)
	delRedisValue(key)
}

/*
 * gift recv & send list
 */
func SetRedisGiftRecvList(id int, list *[]common.GiftSendRecvInfo) {
	key := fmt.Sprintf(REDIS_PREFIX_GIFT_RECVLIST, id)
	setRedisValue(key, list)
}

func GetRedisGiftRecvList(id int) (*[]common.GiftSendRecvInfo, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_GIFT_RECVLIST, id)
	value, exist := getRedisValue(key)
	if true == exist {
		var list []common.GiftSendRecvInfo
		err := json.Unmarshal(value, &list)
		if nil == err {
			return &list, true
		}
	}

	return nil, false
}

func DelRedisGiftRecvList(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_GIFT_RECVLIST, id)
	delRedisValue(key)
}

func SetRedisGiftSendList(id int, list *[]common.GiftSendRecvInfo) {
	key := fmt.Sprintf(REDIS_PREFIX_GIFT_SENDLIST, id)
	setRedisValue(key, list)
}

func GetRedisGiftSendList(id int) (*[]common.GiftSendRecvInfo, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_GIFT_SENDLIST, id)
	value, exist := getRedisValue(key)
	if true == exist {
		var list []common.GiftSendRecvInfo
		err := json.Unmarshal(value, &list)
		if nil == err {
			return &list, true
		}
	}

	return nil, false
}

func DelRedisGiftSendList(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_GIFT_SENDLIST, id)
	delRedisValue(key)
}

/*
 * user gender
 */
func SetRedisUserGender(id, gender int) {
	key := fmt.Sprintf(REDIS_PREFIX_USERGENDER, id)
	setRedisValue(key, gender)
}

func GetRedisUserGender(id int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_USERGENDER, id)
	value, exist := getRedisValue(key)
	if true == exist {
		genderstr := string(value)
		gender, _ := strconv.Atoi(genderstr)
		return gender, true
	}

	return 0, false
}

func DelRedisUserGender(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_USERGENDER, id)
	delRedisValue(key)
}

/*
 * district
 */
func SetRedisDistrict(id int, dist string) {
	if "" == dist {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_DISTRICT, id)
	setRedisValue(key, dist)
}

func GetRedisDistrict(id int) (string, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_DISTRICT, id)
	value, exist := getRedisValue(key)
	if true == exist {
		return string(value), true
	}

	return "", false
}

func DelRedisDistrict(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_DISTRICT, id)
	delRedisValue(key)
}

/*
 * count
 */
func SetRedisProvCount(province string, gender int, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_PROVCOUNT, province, gender)
	setRedisValue(key, count)
}

func GetRedisProvCount(province string, gender int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVCOUNT, province, gender)
	value, exist := getRedisValue(key)
	if true == exist {
		countstr := string(value)
		count, _ := strconv.Atoi(countstr)
		return count, true
	}

	return 0, false
}

func DelRedisProvCount(province string, gender int) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVCOUNT, province, gender)
	delRedisValue(key)
}

func SetRedisProvAgeCount(province string, gender, age, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_PROVAGECOUNT, province, gender, age)
	setRedisValue(key, count)
}

func GetRedisProvAgeCount(province string, gender, age int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVAGECOUNT, province, gender, age)
	value, exist := getRedisValue(key)
	if true == exist {
		countstr := string(value)
		count, _ := strconv.Atoi(countstr)
		return count, true
	}

	return 0, false
}

func DelRedisProvAgeCount(province string, gender, age int) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVAGECOUNT, province, gender, age)
	delRedisValue(key)
}

/*
 * heartbeat province count
 */
func SetRedisHeartbeatProvCount(province string, gender int, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_HEARTBEAT_PROVCOUNT, province, gender)
	setRedisValue(key, count)
}

func GetRedisHeartbeatProvCount(province string, gender int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_HEARTBEAT_PROVCOUNT, province, gender)
	value, exist := getRedisValue(key)
	if true == exist {
		countstr := string(value)
		count, _ := strconv.Atoi(countstr)
		return count, true
	}

	return 0, false
}

func DelRedisHeartbeatProvCount(province string, gender int) {
	key := fmt.Sprintf(REDIS_PREFIX_HEARTBEAT_PROVCOUNT, province, gender)
	delRedisValue(key)
}

/*
 * search info
 */
func SetRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, index int) {
	if 0 == index {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_SEARCHINDEX, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	setRedisValue(key, index)
}

func GetRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHINDEX, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	value, exist := getRedisValue(key)
	if true == exist {
		indexstr := string(value)
		index, _ := strconv.Atoi(indexstr)
		return index, true
	}

	return 0, false
}

func DelRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHINDEX, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	delRedisValue(key)
}

//base
func SetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, base int) {
	if 0 == base {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_SEARCHBASE, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	setRedisValue(key, base)
}

func GetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHBASE, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	value, exist := getRedisValue(key)
	if true == exist {
		basestr := string(value)
		base, _ := strconv.Atoi(basestr)
		return base, true
	}

	return 0, false
}

func DelRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHBASE, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	delRedisValue(key)
}

/*
 * charm toplist
 */
func SetRedisCharmToplist(gender, year int, month time.Month, day int, list *[]common.UserCharmInfo) {
	key := fmt.Sprintf(REDIS_PREFIX_CHARM_TOPLIST, gender, year, month, day)
	setRedisValue(key, list)
}

func GetRedisCharmToplist(gender, year int, month time.Month, day int) (*[]common.UserCharmInfo, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_CHARM_TOPLIST, gender, year, month, day)
	value, exist := getRedisValue(key)
	if true == exist {
		var list []common.UserCharmInfo
		err := json.Unmarshal(value, &list)
		if nil == err {
			return &list, true
		}
	}

	return nil, false
}

func DelRedisCharmToplist(gender, year int, month time.Month, day int) {
	key := fmt.Sprintf(REDIS_PREFIX_CHARM_TOPLIST, gender, year, month, day)
	delRedisValue(key)
}

/*
 * wealth toplist
 */
func SetRedisWealthToplist(year int, month time.Month, day int, list *[]common.UserWealthInfo) {
	key := fmt.Sprintf(REDIS_PREFIX_WEALTH_TOPLIST, year, month, day)
	setRedisValue(key, list)
}

func GetRedisWealthToplist(year int, month time.Month, day int) (*[]common.UserWealthInfo, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_WEALTH_TOPLIST, year, month, day)
	value, exist := getRedisValue(key)
	if true == exist {
		var list []common.UserWealthInfo
		err := json.Unmarshal(value, &list)
		if nil == err {
			return &list, true
		}
	}

	return nil, false
}

func DelRedisWealthToplist(year int, month time.Month, day int) {
	key := fmt.Sprintf(REDIS_PREFIX_WEALTH_TOPLIST, year, month, day)
	delRedisValue(key)
}
