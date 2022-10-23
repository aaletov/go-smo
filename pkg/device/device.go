package device

import (
	"errors"
	"time"

	"github.com/aaletov/go-smo/pkg/events"
	"github.com/aaletov/go-smo/pkg/request"
)

type Request = request.Request
type ReqWGT = request.ReqWGT
type ReqWPT = request.ReqWPT
type ReqWST = request.ReqWST

type Device interface {
	IsFree() bool
	Add(req *ReqWGT) error
	SetIdle()
	GetStartTime() time.Time
	Get() *request.ReqWST
	GetDone() []ReqWPT
	Pop() error
	GetNumber() int
	GetNextEvent() *events.DevFreeEvent
}

var (
	deviceCount = 0
)

func NewDevice(startTime time.Time, pTime time.Duration) Device {
	deviceCount++
	return &deviceImpl{
		deviceNumber: deviceCount,
		pTime:        pTime,
		lastStart:    startTime,
		idle:         true,
		doneReqs:     make([]ReqWPT, 0),
	}
}

type deviceImpl struct {
	deviceNumber int
	pTime        time.Duration
	req          *Request
	idle         bool
	lastStart    time.Time
	doneReqs     []ReqWPT
	nextEvent    *events.DevFreeEvent
}

func (d deviceImpl) IsFree() bool {
	return d.req == nil
}

func (d *deviceImpl) SetIdle() {
	d.idle = true
}

func (d deviceImpl) GetStartTime() time.Time {
	return d.lastStart
}

func (d *deviceImpl) Add(req *ReqWGT) error {
	if d.req != nil {
		return errors.New("Device is busy")
	}
	if d.idle {
		d.lastStart = req.Time
	} else {
		d.lastStart = d.lastStart.Add(d.pTime)
	}
	d.req = req.Req
	d.idle = false

	d.nextEvent = &events.DevFreeEvent{
		Time:   d.lastStart.Add(d.pTime),
		DevNum: d.deviceNumber,
	}

	return nil
}

func (d deviceImpl) Get() *ReqWST {
	return &ReqWST{Req: d.req, Time: d.lastStart}
}

func (d deviceImpl) GetDone() []ReqWPT {
	return d.doneReqs
}

// Тут lastStart + duration
func (d *deviceImpl) Pop() error {
	if d.req == nil {
		return errors.New("No request in device")
	}
	endTime := d.lastStart.Add(d.pTime)
	reqwpt := ReqWPT{
		Req:   d.req,
		Start: d.lastStart,
		End:   endTime,
	}
	d.doneReqs = append(d.doneReqs, reqwpt)
	d.req = nil
	d.nextEvent = nil
	return nil
}

func (d deviceImpl) GetNumber() int {
	return d.deviceNumber
}

func (d deviceImpl) GetNextEvent() *events.DevFreeEvent {
	event := d.nextEvent
	if event != nil {
		d.nextEvent = nil
	}
	return event
}
