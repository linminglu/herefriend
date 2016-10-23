package lib

import (
	"math/rand"
	"time"
)

// CurrentTimeUTCInt64 .
func CurrentTimeUTCInt64() int64 {
	return time.Now().UTC().Unix()
}

// Int64ToUTCTime .
func Int64ToUTCTime(value int64) time.Time {
	return time.Unix(value, 0).UTC()
}

// TimeToUTCInt64 .
func TimeToUTCInt64(timeUTC time.Time) int64 {
	return timeUTC.UTC().Unix()
}

// TimeStrToUTCInt64 .
func TimeStrToUTCInt64(timestr string) int64 {
	var utcvalue int64

	loc, _ := time.LoadLocation("UTC")
	t, err := time.ParseInLocation(time.RFC3339, timestr, loc)
	if nil == err {
		utcvalue = t.Unix()
	}

	return utcvalue
}

const (
	// SleepDurationVIPStatus  .
	SleepDurationVIPStatus = time.Hour
	// SleepDurationLiveStatus .
	SleepDurationLiveStatus = time.Minute
	// SleepDurationNotifyMsg .
	SleepDurationNotifyMsg = time.Second * 15
	// SleepDurationPushQueuMsg .
	SleepDurationPushQueuMsg = time.Second
	// SleepBaseLineRobotComment .
	SleepBaseLineRobotComment = int64(time.Minute * 30)
	// SleepBaseLineRobotVisit .
	SleepBaseLineRobotVisit = int64(time.Minute * 90)
)

const (
	// SleepTypeRobotReply .
	SleepTypeRobotReply = 0
	// SleepTypeRobotComment .
	SleepTypeRobotComment = 1
	// SleepTypeRobotVisit .
	SleepTypeRobotVisit = 2
)

var gHourDuration = [24]time.Duration{
	time.Minute * 40, //0
	time.Hour * 5,    //1
	time.Hour * 4,    //2
	time.Hour * 3,    //3
	time.Hour * 2,    //4
	time.Hour,        //5
	time.Minute * 45, //6
	time.Minute * 20, //7
	time.Second * 45, //8
	time.Second * 50, //9
	time.Second * 40, //10
	time.Second * 40, //11
	time.Second * 40, //12
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

// SleepTimeDuration get the sleep time duration by the current time
func SleepTimeDuration(sleeptype int) time.Duration {
	base := int64(gHourDuration[time.Now().Hour()%24])

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	switch sleeptype {
	case SleepTypeRobotComment:
		if base < SleepBaseLineRobotComment {
			base = SleepBaseLineRobotComment
		}
		return time.Duration(r.Int63n(base) + r.Int63n(base))
	case SleepTypeRobotVisit:
		if base < SleepBaseLineRobotVisit {
			base = SleepBaseLineRobotVisit
		}
		return time.Duration(r.Int63n(base) + r.Int63n(base))
	case SleepTypeRobotReply:
		return time.Duration(r.Int63n(base) + r.Int63n(base))
	default:
		return time.Second
	}
}
