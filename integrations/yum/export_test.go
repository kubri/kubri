package yum

import "time"

func SetTime(t time.Time) {
	timeNow = t.Unix
}
