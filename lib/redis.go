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
	// RedisPrefixUserInfo formats userinfo_${id}
	RedisPrefixUserInfo = "userinfo_%d"
	// RedisPrefixUserGender gender_${id}
	RedisPrefixUserGender = "gender_%d"
	// RedisPrefixDistrict district_${gender}
	RedisPrefixDistrict = "district_%d"
	// RedisPrefixProvCount provcount_${province}_${gender}
	RedisPrefixProvCount = "provcount_%s_%d"
	// RedisPrefixHeartbeatProvCount heartbeat_provcount_${province}_${gender}
	RedisPrefixHeartbeatProvCount = "heartbeat_provcount_%s_%d"
	// RedisPrefixProvAgeCount provcountage_${province}_${gender}_${age}
	RedisPrefixProvAgeCount = "provcountage_%s_%d_%d"
	// RedisPrefixSearchIndex search_${id}_${provicne}_${agemin-max}_${heightmin-max}_${incomemin-max}_${education}_${occupation}_${status}
	RedisPrefixSearchIndex = "search_%d_%s_%d-%d_%d-%d_%d-%d_%s_%s_%s"
	// RedisPrefixSearchBase searchbase_${provicne}_${agemin-max}_${heightmin-max}_${incomemin-max}_${education}_${occupation}_${status}
	RedisPrefixSearchBase = "searchbase_%s_%d-%d_%d-%d_%d-%d_%s_%s_%s"
	// RedisPrefixGoldBeans goldbeans_${id}
	RedisPrefixGoldBeans = "goldbeans_%d"
	// RedisPrefixGiftRecvlist gift_recvlist_${id}
	RedisPrefixGiftRecvlist = "gift_recvlist_%d"
	// RedisPrefixGiftSendlist gift_sendlist_${id}
	RedisPrefixGiftSendlist = "gift_sendlist_%d"
	// RedisPrefixCharmToplist charm_toplist_${gender}_${year}_${month}_${day} (the charm list for last weak)
	RedisPrefixCharmToplist = "charm_toplist_%d_%04d_%02d_%02d"
	// RedisPrefixWealthToplist wealth_toplist_${year}_${month}_${day} (the wealth list for last weak)
	RedisPrefixWealthToplist = "wealth_toplist_%d_%04d_%02d_%02d"
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

func getRedisValue(key string) ([]byte, bool) {
	c := gRedisPool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return nil, false
	}

	return content.([]byte), true
}

func delRedisValue(key string) {
	c := gRedisPool.Get()
	defer c.Close()

	c.Do("Del", key)
}

// SetRedisUserInfo .
func SetRedisUserInfo(id int, info *common.PersonInfo) {
	key := fmt.Sprintf(RedisPrefixUserInfo, id)
	setRedisValue(key, info)
}

// GetRedisUserInfo .
func GetRedisUserInfo(id int) (*common.PersonInfo, bool) {
	key := fmt.Sprintf(RedisPrefixUserInfo, id)
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

// DelRedisUserInfo .
func DelRedisUserInfo(id int) {
	key := fmt.Sprintf(RedisPrefixUserInfo, id)
	delRedisValue(key)
}

// SetRedisGoldBeans .
func SetRedisGoldBeans(id, beans int) {
	key := fmt.Sprintf(RedisPrefixGoldBeans, id)
	setRedisValue(key, beans)
}

// GetRedisGoldBeans .
func GetRedisGoldBeans(id int) (int, bool) {
	key := fmt.Sprintf(RedisPrefixGoldBeans, id)
	value, exist := getRedisValue(key)
	if true == exist {
		beanstr := string(value)
		beans, _ := strconv.Atoi(beanstr)
		return beans, true
	}

	return 0, false
}

// DelRedisGoldBeans .
func DelRedisGoldBeans(id int) {
	key := fmt.Sprintf(RedisPrefixGoldBeans, id)
	delRedisValue(key)
}

// SetRedisGiftRecvList .
func SetRedisGiftRecvList(id int, list *[]common.GiftSendRecvInfo) {
	key := fmt.Sprintf(RedisPrefixGiftRecvlist, id)
	setRedisValue(key, list)
}

// GetRedisGiftRecvList .
func GetRedisGiftRecvList(id int) (*[]common.GiftSendRecvInfo, bool) {
	key := fmt.Sprintf(RedisPrefixGiftRecvlist, id)
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

// DelRedisGiftRecvList .
func DelRedisGiftRecvList(id int) {
	key := fmt.Sprintf(RedisPrefixGiftRecvlist, id)
	delRedisValue(key)
}

// SetRedisGiftSendList .
func SetRedisGiftSendList(id int, list *[]common.GiftSendRecvInfo) {
	key := fmt.Sprintf(RedisPrefixGiftSendlist, id)
	setRedisValue(key, list)
}

// GetRedisGiftSendList .
func GetRedisGiftSendList(id int) (*[]common.GiftSendRecvInfo, bool) {
	key := fmt.Sprintf(RedisPrefixGiftSendlist, id)
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

// DelRedisGiftSendList .
func DelRedisGiftSendList(id int) {
	key := fmt.Sprintf(RedisPrefixGiftSendlist, id)
	delRedisValue(key)
}

// SetRedisUserGender .
func SetRedisUserGender(id, gender int) {
	key := fmt.Sprintf(RedisPrefixUserGender, id)
	setRedisValue(key, gender)
}

// GetRedisUserGender .
func GetRedisUserGender(id int) (int, bool) {
	key := fmt.Sprintf(RedisPrefixUserGender, id)
	value, exist := getRedisValue(key)
	if true == exist {
		genderstr := string(value)
		gender, _ := strconv.Atoi(genderstr)
		return gender, true
	}

	return 0, false
}

// DelRedisUserGender .
func DelRedisUserGender(id int) {
	key := fmt.Sprintf(RedisPrefixUserGender, id)
	delRedisValue(key)
}

// SetRedisDistrict .
func SetRedisDistrict(id int, dist string) {
	if "" == dist {
		return
	}

	key := fmt.Sprintf(RedisPrefixDistrict, id)
	setRedisValue(key, dist)
}

// GetRedisDistrict .
func GetRedisDistrict(id int) (string, bool) {
	key := fmt.Sprintf(RedisPrefixDistrict, id)
	value, exist := getRedisValue(key)
	if true == exist {
		return string(value), true
	}

	return "", false
}

// DelRedisDistrict .
func DelRedisDistrict(id int) {
	key := fmt.Sprintf(RedisPrefixDistrict, id)
	delRedisValue(key)
}

// SetRedisProvCount .
func SetRedisProvCount(province string, gender int, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(RedisPrefixProvCount, province, gender)
	setRedisValue(key, count)
}

// GetRedisProvCount .
func GetRedisProvCount(province string, gender int) (int, bool) {
	key := fmt.Sprintf(RedisPrefixProvCount, province, gender)
	value, exist := getRedisValue(key)
	if true == exist {
		countstr := string(value)
		count, _ := strconv.Atoi(countstr)
		return count, true
	}

	return 0, false
}

// DelRedisProvCount .
func DelRedisProvCount(province string, gender int) {
	key := fmt.Sprintf(RedisPrefixProvCount, province, gender)
	delRedisValue(key)
}

// SetRedisProvAgeCount .
func SetRedisProvAgeCount(province string, gender, age, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(RedisPrefixProvAgeCount, province, gender, age)
	setRedisValue(key, count)
}

// GetRedisProvAgeCount .
func GetRedisProvAgeCount(province string, gender, age int) (int, bool) {
	key := fmt.Sprintf(RedisPrefixProvAgeCount, province, gender, age)
	value, exist := getRedisValue(key)
	if true == exist {
		countstr := string(value)
		count, _ := strconv.Atoi(countstr)
		return count, true
	}

	return 0, false
}

// DelRedisProvAgeCount .
func DelRedisProvAgeCount(province string, gender, age int) {
	key := fmt.Sprintf(RedisPrefixProvAgeCount, province, gender, age)
	delRedisValue(key)
}

// SetRedisHeartbeatProvCount .
func SetRedisHeartbeatProvCount(province string, gender int, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(RedisPrefixHeartbeatProvCount, province, gender)
	setRedisValue(key, count)
}

// GetRedisHeartbeatProvCount .
func GetRedisHeartbeatProvCount(province string, gender int) (int, bool) {
	key := fmt.Sprintf(RedisPrefixHeartbeatProvCount, province, gender)
	value, exist := getRedisValue(key)
	if true == exist {
		countstr := string(value)
		count, _ := strconv.Atoi(countstr)
		return count, true
	}

	return 0, false
}

// DelRedisHeartbeatProvCount .
func DelRedisHeartbeatProvCount(province string, gender int) {
	key := fmt.Sprintf(RedisPrefixHeartbeatProvCount, province, gender)
	delRedisValue(key)
}

// SetRedisSearchIndex .
func SetRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, index int) {
	if 0 == index {
		return
	}

	key := fmt.Sprintf(RedisPrefixSearchIndex, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	setRedisValue(key, index)
}

// GetRedisSearchIndex .
func GetRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) (int, bool) {
	key := fmt.Sprintf(RedisPrefixSearchIndex, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	value, exist := getRedisValue(key)
	if true == exist {
		indexstr := string(value)
		index, _ := strconv.Atoi(indexstr)
		return index, true
	}

	return 0, false
}

// DelRedisSearchIndex .
func DelRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) {
	key := fmt.Sprintf(RedisPrefixSearchIndex, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	delRedisValue(key)
}

// SetRedisSearchBase .
func SetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, base int) {
	if 0 == base {
		return
	}

	key := fmt.Sprintf(RedisPrefixSearchBase, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	setRedisValue(key, base)
}

// GetRedisSearchBase .
func GetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) (int, bool) {
	key := fmt.Sprintf(RedisPrefixSearchBase, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	value, exist := getRedisValue(key)
	if true == exist {
		basestr := string(value)
		base, _ := strconv.Atoi(basestr)
		return base, true
	}

	return 0, false
}

// DelRedisSearchBase .
func DelRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) {
	key := fmt.Sprintf(RedisPrefixSearchBase, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	delRedisValue(key)
}

// SetRedisCharmToplist .
func SetRedisCharmToplist(gender, year int, month time.Month, day int, list *[]common.UserCharmInfo) {
	key := fmt.Sprintf(RedisPrefixCharmToplist, gender, year, month, day)
	setRedisValue(key, list)
}

// GetRedisCharmToplist .
func GetRedisCharmToplist(gender, year int, month time.Month, day int) (*[]common.UserCharmInfo, bool) {
	key := fmt.Sprintf(RedisPrefixCharmToplist, gender, year, month, day)
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

// DelRedisCharmToplist .
func DelRedisCharmToplist(gender, year int, month time.Month, day int) {
	key := fmt.Sprintf(RedisPrefixCharmToplist, gender, year, month, day)
	delRedisValue(key)
}

// SetRedisWealthToplist .
func SetRedisWealthToplist(year int, month time.Month, day int, list *[]common.UserWealthInfo) {
	key := fmt.Sprintf(RedisPrefixWealthToplist, year, month, day)
	setRedisValue(key, list)
}

// GetRedisWealthToplist .
func GetRedisWealthToplist(year int, month time.Month, day int) (*[]common.UserWealthInfo, bool) {
	key := fmt.Sprintf(RedisPrefixWealthToplist, year, month, day)
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

// DelRedisWealthToplist .
func DelRedisWealthToplist(year int, month time.Month, day int) {
	key := fmt.Sprintf(RedisPrefixWealthToplist, year, month, day)
	delRedisValue(key)
}
