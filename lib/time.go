package lib

import (
	"math/rand"
	"time"
)

/*
 * APIs about convert between time and int64
 *
 */
func CurrentTimeUTCInt64() int64 {
	return time.Now().UTC().Unix()
}

func Int64_To_UTCTime(value int64) time.Time {
	return time.Unix(value, 0).UTC()
}

func Time_To_UTCInt64(timeUTC time.Time) int64 {
	return timeUTC.UTC().Unix()
}

func TimeStr_To_UTCInt64(timestr string) int64 {
	var utcvalue int64

	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(time.RFC3339, timestr, loc)
	if nil == err {
		utcvalue = t.Unix()
	}

	return utcvalue
}

/*
 * Get random sleep duration with particular type
 *
 */
const (
	SLEEP_DURATION_VIPSTATUS      = time.Hour
	SLEEP_DURATION_LIVESTATUS     = time.Minute
	SLEEP_DURATION_NOTIFYMSG      = time.Second * 15
	SLEEP_DURATION_PUSH_QUEUEMSG  = time.Second
	SLEEP_BASELINE_ROBOTRECOMMEND = int64(time.Minute * 30)
	SLEEP_BASELINE_ROBOTVISIT     = int64(time.Minute * 60)
)

const (
	SLEEP_TYPE_ROBOTREPLY     = 0
	SLEEP_TYPE_ROBOTRECOMMEND = 1
	SLEEP_TYPE_ROBOTVISIT     = 2
)

var gHourDuration = [24]time.Duration{
	time.Minute * 10, //0
	time.Hour * 4,    //1
	time.Hour * 3,    //2
	time.Hour * 2,    //3
	time.Hour,        //4
	time.Hour,        //5
	time.Minute * 35, //6
	time.Minute * 2,  //7
	time.Second * 35, //8
	time.Second * 20, //9
	time.Second * 30, //10
	time.Second * 30, //11
	time.Second * 20, //12
	time.Minute * 1,  //13
	time.Minute * 1,  //14
	time.Minute * 1,  //15
	time.Minute * 1,  //16
	time.Minute * 1,  //17
	time.Minute * 2,  //18
	time.Second * 40, //19
	time.Second * 40, //20
	time.Second * 20, //21
	time.Second * 20, //22
	time.Minute * 2,  //23
}

/*
 |    Function: SleepTimeDuration
 |      Author: Mr.Sancho
 |        Date: 2016-01-13
 |   Arguments:
 |      Return:
 | Description: get the sleep time duration by the current time
 |
*/
func SleepTimeDuration(sleeptype int) time.Duration {
	base := int64(gHourDuration[time.Now().Hour()%24])

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	switch sleeptype {
	case SLEEP_TYPE_ROBOTRECOMMEND:
		if base < SLEEP_BASELINE_ROBOTRECOMMEND {
			base = SLEEP_BASELINE_ROBOTRECOMMEND
		}
		return time.Duration(r.Int63n(base) + r.Int63n(base))
	case SLEEP_TYPE_ROBOTVISIT:
		if base < SLEEP_BASELINE_ROBOTVISIT {
			base = SLEEP_BASELINE_ROBOTVISIT
		}
		return time.Duration(r.Int63n(base) + r.Int63n(base))
	case SLEEP_TYPE_ROBOTREPLY:
		return time.Duration(r.Int63n(base) + r.Int63n(base))
	default:
		return time.Second
	}
}
