package yum

import "time"

func SetTime(t time.Time) {
	timeNow = func() int { return int(t.Unix()) }
}
