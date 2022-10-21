package main

import (
	"time"

	"github.com/aaletov/go-smo/pkg/buffer"
	cmgr "github.com/aaletov/go-smo/pkg/choice-manager"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/device"
	smgr "github.com/aaletov/go-smo/pkg/set-manager"
	"github.com/aaletov/go-smo/pkg/source"
)

const (
	sourcesLambda = 13
	sourcesCount  = 3
	bufferCount   = 4
	deviceCount   = 3
)

func main() {
	clock.InitClock(time.Now())
	sources := make([]source.Source, sourcesCount)
	for i := 0; i < sourcesCount; i++ {
		sources[i] = source.NewSource(sourcesLambda)
	}
	buffers := make([]buffer.Buffer, bufferCount)
	for i := 0; i < bufferCount; i++ {
		buffers[i] = buffer.NewBuffer()
	}
	setManager := smgr.NewSetManager(sources, buffers)
	devices := make([]device.Device, deviceCount)
	for i := 0; i < deviceCount; i++ {
		devDuration := time.Duration(1e10 * (10 + i))
		devices[i] = device.NewDevice(clock.SMOClock.Time, devDuration)
	}
	choiceManager := cmgr.NewChoiceManager(buffers, devices)

	for i := 0; i < 100; i++ {
		setManager.Collect()
		setManager.ToBuffer()
		choiceManager.ToDevices()
	}
}
