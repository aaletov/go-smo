package system

import (
	"time"

	"github.com/aaletov/go-smo/pkg/buffer"
	cmgr "github.com/aaletov/go-smo/pkg/choice-manager"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/device"
	smgr "github.com/aaletov/go-smo/pkg/set-manager"
	"github.com/aaletov/go-smo/pkg/source"
)

type System struct {
	Sources   []source.Source
	SetMgr    smgr.SetManager
	Buffers   []buffer.Buffer
	ChoiceMgr cmgr.ChoiceManager
	Devices   []device.Device
}

func NewSystem(sourcesCount, buffersCount, devicesCount int, sourcesLambda, devDuration time.Duration) *System {
	sources := make([]source.Source, sourcesCount)
	for i := 0; i < sourcesCount; i++ {
		sources[i] = source.NewSource(sourcesLambda)
	}
	buffers := make([]buffer.Buffer, buffersCount)
	for i := 0; i < buffersCount; i++ {
		buffers[i] = buffer.NewBuffer()
	}
	setManager := smgr.NewSetManager(sources, buffers)
	devices := make([]device.Device, devicesCount)
	for i := 0; i < devicesCount; i++ {
		devDuration := time.Duration(1e10 * (10 + i))
		devices[i] = device.NewDevice(clock.SMOClock.Time, devDuration)
	}
	choiceManager := cmgr.NewChoiceManager(buffers, devices)

	return &System{sources, setManager, buffers, choiceManager, devices}
}

func (s *System) Iterate() {
	s.SetMgr.Collect()
	s.SetMgr.ToBuffer()
	s.ChoiceMgr.ToDevices()
}
