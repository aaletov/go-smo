package system

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

func (s System) printSourcesTable() {
	sourceTable := table.NewWriter()
	sourceTable.SetOutputMirror(os.Stdout)
	for _, s := range s.Sources {
		sourceRow := []any{fmt.Sprintf("Source #%v", s.GetNumber())}
		for _, r := range s.GetGenerated() {
			sourceRow = append(sourceRow, r.Request.String())
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
			bufRow = append(bufRow, r.Request.String())
		}
		if b.Get() != nil {
			bufRow = append(bufRow, "-> "+b.Get().Request.String())
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
			devRow = append(devRow, rwpt.Request.String())
		}
		if !d.IsFree() {
			devRow = append(devRow, "-> "+d.Get().Request.String())
		}
		devTable.AppendRow(devRow)
		devTable.AppendSeparator()
	}
	devTable.Render()
}

func (s System) printReject() {
	rejTable := table.NewWriter()
	rejTable.SetOutputMirror(os.Stdout)
	rejRow := []any{"Reject"}
	for _, rwse := range s.SetMgr.GetRejectList() {
		rejRow = append(rejRow, rwse.Request.String())
	}
	rejTable.AppendRow(rejRow)
	rejTable.Render()
}

func (s System) PrintData() {
	s.printSourcesTable()
	s.printBuffersTable()
	s.printDevTable()
	s.printReject()
}
