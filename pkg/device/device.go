package device

import (
	"errors"
	"time"

	"golang.org/x/exp/rand"

	"github.com/aaletov/go-smo/api"
	"github.com/aaletov/go-smo/pkg/events"
	"github.com/aaletov/go-smo/pkg/source"
)

type Request = api.Request
type ReqWGT = api.ReqWGT
type ReqWPT = api.ReqWPT
type ReqWST = api.ReqWST

type Device interface {
	IsFree() bool
	Add(req *ReqWGT) error
	SetIdle()
	GetStartTime() time.Time
	Get() *api.ReqWST
	GetDone() []ReqWPT
	Pop() error
	GetNumber() int
	GetNextEvent() *events.DevFreeEvent
	GetDevA() time.Duration
	GetDevB() time.Duration
}

var (
	DeviceCount = 0
)

func NewDevice(startTime time.Time, a, b time.Duration) Device {
	DeviceCount++
	return &deviceImpl{
		deviceNumber: DeviceCount,
		a:            int64(a),
		b:            int64(b),
		rand:         rand.New(source.RandSource),
		lastStart:    startTime,
		idle:         true,
		doneReqs:     make([]ReqWPT, 0),
	}
}

type deviceImpl struct {
	deviceNumber int
	a, b         int64
	rand         *rand.Rand
	req          *api.Request
	idle         bool
	lastStart    time.Time
	lastDuration time.Duration
	doneReqs     []ReqWPT
	nextEvent    *events.DevFreeEvent
}

func (d deviceImpl) genDuration() time.Duration {
	return time.Duration(float64(d.a) + float64(d.b-d.a)*(rand.Float64()))
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
		d.lastStart = d.lastStart.Add(d.lastDuration)
	}
	d.lastDuration = d.genDuration()
	d.req = &req.Request
	d.idle = false

	d.nextEvent = &events.DevFreeEvent{
		Time:   d.lastStart.Add(d.lastDuration),
		DevNum: d.deviceNumber,
	}

	return nil
}

func (d deviceImpl) Get() *ReqWST {
	return &ReqWST{Request: *d.req, Time: d.lastStart}
}

func (d deviceImpl) GetDone() []ReqWPT {
	return d.doneReqs
}

// Тут lastStart + duration
func (d *deviceImpl) Pop() error {
	if d.req == nil {
		return errors.New("No request in device")
	}
	endTime := d.lastStart.Add(d.lastDuration)
	reqwpt := ReqWPT{
		Request: *d.req,
		Start:   d.lastStart,
		End:     endTime,
	}
	d.doneReqs = append(d.doneReqs, reqwpt)
	d.req = nil
	d.nextEvent = nil
	return nil
}

func (d deviceImpl) GetNumber() int {
	return d.deviceNumber
}

func (d *deviceImpl) GetNextEvent() *events.DevFreeEvent {
	if d.nextEvent != nil {
		event := *d.nextEvent
		d.nextEvent = nil
		return &event
	} else {
		return nil
	}
}

func (d deviceImpl) GetDevA() time.Duration {
	return time.Duration(d.a)
}

func (d deviceImpl) GetDevB() time.Duration {
	return time.Duration(d.b)
}
