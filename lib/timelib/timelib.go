package timelib

import "time"

func GetEpoch() float64 {
	return float64(time.Now().Unix())
}
