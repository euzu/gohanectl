package utils

import "time"

func NowTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}