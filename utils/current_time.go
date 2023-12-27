package utils

import "time"

func CurrentTime() int64 {
	currentTime := time.Now()
	unixTimestampSeconds := currentTime.Unix()

	return unixTimestampSeconds
}