package system

import (
	"strings"
	"sync"
	"time"

	"github.com/aaletov/go-smo/pkg/buffer"
	cmgr "github.com/aaletov/go-smo/pkg/choice-manager"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/device"
	"github.com/aaletov/go-smo/pkg/events"
	"github.com/aaletov/go-smo/pkg/queue"
	smgr "github.com/aaletov/go-smo/pkg/set-manager"
	"github.com/aaletov/go-smo/pkg/source"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

type Queue = queue.PriorityQueue[events.Event]

type System struct {
	Logger    *logrus.Entry
	Iteration int
	Sources   []source.Source
	SetMgr    smgr.SetManager
	Buffers   []buffer.Buffer
	ChoiceMgr cmgr.ChoiceManager
	Devices   []device.Device
	Events    Queue
}

var (
	GlobalSystem *System
	SysLock      *sync.Mutex
)

func InitSystem(sourcesCount, buffersCount, devicesCount int, sourcesLambda, devA, devB time.Duration) {
	source.SourcesCount = 0
	buffer.BufCount = 0
	device.DeviceCount = 0

	logger := logrus.New()
	logger.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"component", "method"},
	})
	sources := make([]source.Source, sourcesCount)
	for i := 0; i < sourcesCount; i++ {
		sources[i] = source.NewSource(logger, sourcesLambda)
	}
	buffers := make([]buffer.Buffer, buffersCount)
	for i := 0; i < buffersCount; i++ {
		buffers[i] = buffer.NewBuffer(logger)
	}
	events := queue.NewPriorityQueue[events.Event]()
	setManager := smgr.NewSetManager(logger, sources, buffers, events)

	for _, s := range sources {
		setManager.GetEventFromSource(s.GetNumber())
	}

	devices := make([]device.Device, devicesCount)
	for i := 0; i < devicesCount; i++ {
		devices[i] = device.NewDevice(clock.SMOClock.StartTime, devA, devB)
	}
	choiceManager := cmgr.NewChoiceManager(buffers, devices)

	ll := logger.WithFields(logrus.Fields{
		"component": "System",
	})

	GlobalSystem = &System{
		ll,
		0,
		sources,
		setManager,
		buffers,
		choiceManager,
		devices,
		events,
	}
	SysLock = &sync.Mutex{}

	return
}

func (s *System) GetEvents() {
	ll := s.Logger.WithField("method", "GetEvents")

	sb := strings.Builder{}
	sb.WriteString("Got events: [ ")
	for _, d := range s.Devices {
		event := d.GetNextEvent()
		if event != nil {
			s.Events.Add(event)
			sb.WriteString(event.String() + " ")
		}
	}
	sb.WriteString("]")
	ll.Info(sb.String())
	ll.Infof("Front queue element: %v", s.Events.Front().Get().String())
}

func (s *System) Iterate() {
	ll := s.Logger.WithField("method", "Iterate")
	s.GetEvents()
	nextEvent := s.Events.Front().Get()
	s.Events.Pop()
	switch e := nextEvent.(type) {
	case *events.GenReqEvent:
		s.HandleGenReqEvent(e)
	case *events.DevFreeEvent:
		s.HandleDevFreeEvent(e)
	}
	ll.Infof("Processed event %v\n", nextEvent.String())
}

func (s *System) HandleGenReqEvent(event *events.GenReqEvent) {
	clock.SMOClock.Time = event.Time
	s.SetMgr.ProcessSource(event.SourceNum)
	s.ChoiceMgr.HandleGenReqEvent(event)
}

func (s *System) HandleDevFreeEvent(event *events.DevFreeEvent) {
	clock.SMOClock.Time = event.Time
	s.ChoiceMgr.HandleDevFreeEvent(event)
}
