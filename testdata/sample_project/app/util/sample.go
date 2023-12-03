package util

import "time"

func GetNowStringInJst() string {
	return time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("2006-01-02 15:04:05")
}
