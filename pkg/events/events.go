package events

import "time"

type Event interface {
	GetTime() time.Time
}

type GenReqEvent struct {
	Time      time.Time
	SourceNum int
}

func (g GenReqEvent) GetTime() time.Time {
	return g.Time
}

type DevFreeEvent struct {
	Time   time.Time
	DevNum int
}

func (d DevFreeEvent) GetTime() time.Time {
	return d.Time
}
