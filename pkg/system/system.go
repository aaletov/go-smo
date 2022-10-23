package system

import (
	"fmt"
	"os"
	"strings"
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

	"github.com/jedib0t/go-pretty/v6/table"
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

func NewSystem(sourcesCount, buffersCount, devicesCount int, sourcesLambda, devDuration time.Duration) *System {
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
	setManager := smgr.NewSetManager(sources, buffers)
	devices := make([]device.Device, devicesCount)
	for i := 0; i < devicesCount; i++ {
		devDuration := time.Duration(1e9 * (10 + i))
		devices[i] = device.NewDevice(clock.SMOClock.Time, devDuration)
	}
	choiceManager := cmgr.NewChoiceManager(buffers, devices)

	ll := logger.WithFields(logrus.Fields{
		"component": "System",
	})
	return &System{
		ll,
		0,
		sources,
		setManager,
		buffers,
		choiceManager,
		devices,
		queue.NewPriorityQueue[events.Event](),
	}
}

func (s *System) GetEvents() {
	ll := s.Logger.WithField("method", "GetEvents")

	sb := strings.Builder{}
	sb.WriteString("Got events: [ ")
	for _, source := range s.Sources {
		event := source.GetNextEvent()
		s.Events.Add(event)
		sb.WriteString(event.String() + " ")
	}
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
	s.SetMgr.ProcessSource(event.SourceNum)
	s.ChoiceMgr.HandleGenReqEvent(event)
}

func (s *System) HandleDevFreeEvent(event *events.DevFreeEvent) {
	s.ChoiceMgr.HandleDevFreeEvent(event)
}

func (s System) printSourcesTable() {
	sourceTable := table.NewWriter()
	sourceTable.SetOutputMirror(os.Stdout)
	for _, s := range s.Sources {
		sourceRow := []any{fmt.Sprintf("Source #%v", s.GetNumber())}
		for _, r := range s.GetGenerated() {
			sourceRow = append(sourceRow, r.Req.String())
		}
		sourceTable.AppendRow(sourceRow)
		sourceTable.AppendSeparator()
	}
	sourceTable.Render()
}

func (s System) printBuffersTable() {
	bufferTable := table.NewWriter()
	bufferTable.SetOutputMirror(os.Stdout)
	for _, b := range s.Buffers {
		bufRow := []any{fmt.Sprintf("Buffer #%v", b.GetNumber())}
		for _, r := range b.GetAllProcessed() {
			bufRow = append(bufRow, r.Req.String())
		}
		if b.Get() != nil {
			bufRow = append(bufRow, "-> "+b.Get().Req.String())
		}
		bufferTable.AppendRow(bufRow)
		bufferTable.AppendSeparator()
	}
	bufferTable.Render()
}

func (s System) printDevTable() {
	devTable := table.NewWriter()
	devTable.SetOutputMirror(os.Stdout)
	for _, d := range s.Devices {
		devRow := []any{fmt.Sprintf("Device #%v", d.GetNumber())}
		for _, rwpt := range d.GetDone() {
			devRow = append(devRow, rwpt.Req.String())
		}
		if !d.IsFree() {
			devRow = append(devRow, "-> "+d.Get().Req.String())
		}
		devTable.AppendRow(devRow)
		devTable.AppendSeparator()
	}
	devTable.Render()
}

func (s System) PrintData() {
	s.printSourcesTable()
	s.printBuffersTable()
	s.printDevTable()
}
