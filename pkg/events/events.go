package events

import (
	"strconv"
	"time"

	"github.com/aaletov/go-smo/pkg/queue"
)

type Event interface {
	GetTime() time.Time
	Less(other queue.Comparable) bool
	String() string
}

type GenReqEvent struct {
	Time      time.Time
	SourceNum int
}

func (g GenReqEvent) GetTime() time.Time {
	return g.Time
}

func (g GenReqEvent) Less(other queue.Comparable) bool {
	return g.Time.Before(other.(Event).GetTime())
}

func (g GenReqEvent) String() string {
	return "GR[" + strconv.Itoa(g.SourceNum) + "]"
}

type DevFreeEvent struct {
	Time   time.Time
	DevNum int
}

func (d DevFreeEvent) GetTime() time.Time {
	return d.Time
}

func (d DevFreeEvent) Less(other queue.Comparable) bool {
	return d.Time.Before(other.(Event).GetTime())
}

func (d DevFreeEvent) String() string {
	return "DF[" + strconv.Itoa(d.DevNum) + "]"
}
