package system

import (
	"fmt"
	"os"
	"time"

	"github.com/aaletov/go-smo/pkg/buffer"
	cmgr "github.com/aaletov/go-smo/pkg/choice-manager"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/device"
	smgr "github.com/aaletov/go-smo/pkg/set-manager"
	"github.com/aaletov/go-smo/pkg/source"

	"github.com/jedib0t/go-pretty/v6/table"
)

type System struct {
	Iteration int
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

	return &System{0, sources, setManager, buffers, choiceManager, devices}
}

func (s *System) Iterate() {
	s.SetMgr.Iterate()
	s.ChoiceMgr.Iterate()
	for _, d := range s.Devices {
		d.Pop()
	}
}

func (s System) PrintData() {
	fmt.Println("Sources")
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

	fmt.Println("Buffers")
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

	for _, d := range s.Devices {
		if !d.IsFree() {
			fmt.Printf("Device #%v processes %v\n", d.GetNumber(), d.Get().Req.String())
		}
	}
}
