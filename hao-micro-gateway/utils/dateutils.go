package utils

import (
	"time"
)

var Format_YMD = "2006-01-02"
var Format_YMDHS = "2006-01-02 15:04:05"

func TimeFormat(fromat string, currentTime time.Time) string {
	return currentTime.Format(fromat)
}

func TimeFormatNow(fromat string) string {
	currentTime := time.Now() //获取当前时间
	return currentTime.Format(fromat)
}

func DiffTime(startTime time.Time, endTime time.Time) int64 {
	duration := endTime.Sub(startTime)
	seconds := int64(duration.Seconds())
	// fmt.Println("时间差（秒）：", seconds)
	return seconds
}
