package cmgr

import (
	"sort"

	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/device"
	"github.com/aaletov/go-smo/pkg/events"
	"github.com/aaletov/go-smo/pkg/request"
)

type ChoiceManager interface {
	HandleGenReqEvent(*events.GenReqEvent)
	HandleDevFreeEvent(*events.DevFreeEvent)
}

func NewChoiceManager(buffers []buffer.Buffer, devices []device.Device) ChoiceManager {
	return &choiceManagerImpl{
		buffers: buffers,
		devices: devices,
		bufPtr:  0,
	}
}

type choiceManagerImpl struct {
	buffers []buffer.Buffer
	devices []device.Device
	bufPtr  int
}

func (c *choiceManagerImpl) toDevices() {
	reqToBuf := make(map[*request.ReqWGT]int, 0)
	reqwgtSlice := make([]*buffer.ReqWGT, 0)
	for i, b := range c.buffers {
		if reqwgt := b.Get(); reqwgt != nil {
			reqwgtSlice = append(reqwgtSlice, reqwgt)
			reqToBuf[reqwgt] = i
		}
	}
	sort.Slice(reqwgtSlice, func(i, j int) bool {
		iSource := reqwgtSlice[i].Req.SourceNumber
		jSource := reqwgtSlice[j].Req.SourceNumber
		return (iSource < jSource) || ((iSource == jSource) &&
			(reqwgtSlice[i].Req.RequestNumber < reqwgtSlice[j].Req.RequestNumber))
	})

	for _, device := range c.devices {
		if device.IsFree() {
			if len(reqwgtSlice) != 0 {
				reqwgt := reqwgtSlice[0]
				device.Add(reqwgt)
				bufNum := reqToBuf[reqwgt]
				c.buffers[bufNum].Pop(device.GetStartTime())
				reqwgtSlice = reqwgtSlice[1:]
			} else {
				device.SetIdle()
			}
		}
	}
}

func (c *choiceManagerImpl) HandleGenReqEvent(event *events.GenReqEvent) {
	c.toDevices()
}

func (c *choiceManagerImpl) HandleDevFreeEvent(event *events.DevFreeEvent) {
	c.devices[event.DevNum-1].Pop()
	c.toDevices()
}
