package cmgr

import (
	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/device"
)

type ChoiceManager interface {
	Collect()
	ToDevices()
}

type SetManager struct {
	buffers []*buffer.Buffer
	devices []*device.Device
}