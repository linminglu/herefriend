package lib

import (
	"math"
	"math/rand"
	"time"
)

func Intn(v int) int {
	if 0 >= v {
		return 0
	}

	src := rand.NewSource(time.Now().UnixNano() + rand.Int63())
	return rand.New(src).Intn(v)
}

func Int63n(v int64) int64 {
	if 0 >= v {
		return 0
	}

	src := rand.NewSource(time.Now().UnixNano())
	return rand.New(src).Int63n(v)
}

func TruncFloat(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func RandomHitPercent(percent int) bool {
	src := rand.NewSource(time.Now().UnixNano())
	key := rand.New(src).Intn(100)

	if key < percent {
		return true
	} else {
		return false
	}
}
