package clock

import (
	"time"
)

// No thread safety (at all)
type Clock struct {
	Time time.Time
}

var (
	SMOClock *Clock = nil
)

func InitClock(time time.Time) {
	SMOClock = &Clock{time}
}
