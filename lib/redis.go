package lib

import (
	"fmt"
	"strconv"

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
)

var g_redispool *redis.Pool

func init() {
	g_redispool = &redis.Pool{
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
 * search info
 */
func SetRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, index int) {
	if 0 == index {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_SEARCHINDEX, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(strconv.Itoa(index)))
}

func GetRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHINDEX, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	c := g_redispool.Get()
	defer c.Close()
	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return 0, false
	}

	indexstr := string(content.([]byte))
	index, _ := strconv.Atoi(indexstr)
	return index, true
}

func DelRedisSearchIndex(id int, agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHINDEX, id, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

//base
func SetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string, base int) {
	if 0 == base {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_SEARCHBASE, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(strconv.Itoa(base)))
}

func GetRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHBASE, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return 0, false
	}

	basestr := string(content.([]byte))
	base, _ := strconv.Atoi(basestr)
	return base, true
}

func DelRedisSearchBase(agemin, agemax, heightmin, heightmax, incomemin, incomemax int, province, education, occupation, status string) {
	key := fmt.Sprintf(REDIS_PREFIX_SEARCHBASE, province, agemin, agemax, heightmin, heightmax, incomemin, incomemax, education, occupation, status)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

/*
 * user info
 */
func SetRedisUserInfo(id int, content []byte) {
	key := fmt.Sprintf(REDIS_PREFIX_USERINFO, id)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, content)
}

func GetRedisUserInfo(id int) ([]byte, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_USERINFO, id)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return nil, false
	}

	return content.([]byte), true
}

func DelRedisUserInfo(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_USERINFO, id)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

/*
 * user gender
 */
func SetRedisUserGender(id, gender int) {
	key := fmt.Sprintf(REDIS_PREFIX_USERGENDER, id)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(strconv.Itoa(gender)))
}

func GetRedisUserGender(id int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_USERGENDER, id)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return 0, false
	}

	genderstr := string(content.([]byte))
	gender, _ := strconv.Atoi(genderstr)
	return gender, true
}

func DelRedisUserGender(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_USERGENDER, id)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

/*
 * district
 */
func SetRedisDistrict(id int, dist string) {
	if "" == dist {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_DISTRICT, id)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(dist))
}

func GetRedisDistrict(id int) (string, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_DISTRICT, id)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return "", false
	}

	return string(content.([]byte)), true
}

func DelRedisDistrict(id int) {
	key := fmt.Sprintf(REDIS_PREFIX_DISTRICT, id)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

/*
 * count
 */
func SetRedisProvCount(province string, gender int, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_PROVCOUNT, province, gender)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(strconv.Itoa(count)))
}

func GetRedisProvCount(province string, gender int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVCOUNT, province, gender)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return 0, false
	}

	countstr := string(content.([]byte))
	count, _ := strconv.Atoi(countstr)
	return count, true
}

func DelRedisProvCount(province string, gender int) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVCOUNT, province, gender)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

func SetRedisProvAgeCount(province string, gender, age, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_PROVAGECOUNT, province, gender, age)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(strconv.Itoa(count)))
}

func GetRedisProvAgeCount(province string, gender, age int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVAGECOUNT, province, gender, age)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return 0, false
	}

	countstr := string(content.([]byte))
	count, _ := strconv.Atoi(countstr)
	return count, true
}

func DelRedisProvAgeCount(province string, gender, age int) {
	key := fmt.Sprintf(REDIS_PREFIX_PROVAGECOUNT, province, gender, age)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}

/*
 * heartbeat province count
 */
func SetRedisHeartbeatProvCount(province string, gender int, count int) {
	if 0 == count {
		return
	}

	key := fmt.Sprintf(REDIS_PREFIX_HEARTBEAT_PROVCOUNT, province, gender)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Set", key, []byte(strconv.Itoa(count)))
}

func GetRedisHeartbeatProvCount(province string, gender int) (int, bool) {
	key := fmt.Sprintf(REDIS_PREFIX_HEARTBEAT_PROVCOUNT, province, gender)
	c := g_redispool.Get()
	defer c.Close()

	content, err := c.Do("Get", key)
	if nil == content || nil != err {
		return 0, false
	}

	countstr := string(content.([]byte))
	count, _ := strconv.Atoi(countstr)
	return count, true
}

func DelRedisHeartbeatProvCount(province string, gender int) {
	key := fmt.Sprintf(REDIS_PREFIX_HEARTBEAT_PROVCOUNT, province, gender)
	c := g_redispool.Get()
	defer c.Close()

	c.Do("Del", key)
}
