package utils

import (
	"strconv"
	"time"
)

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func StringStampToTime(stampStr string) time.Time {
	timestamp, _ := strconv.ParseInt(stampStr, 10, 64)
	return time.UnixMilli(timestamp)
}

func StringToTimeYYYYMMDD(str string) (time.Time, error) {
	return time.Parse("2006-01-02", str)
}
