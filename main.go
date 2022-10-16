package main

import (
	"github.com/aaletov/go-smo/pkg/source"
	"github.com/aaletov/go-smo/pkg/buffer"
	smgr "github.com/aaletov/go-smo/pkg/set-manager"
	"github.com/aaletov/go-smo/pkg/request"
	"github.com/aaletov/go-smo/pkg/device"
	"github.com/aaletov/go-smo/pkg/queue"
)

func main() {
	_ = source.NewSource(0)
	_ = buffer.NewBuffer()
	_ = new(smgr.SetManager)
	_ = new(request.Request)
	_ = new(device.Device)
	_ = new(queue.PriorityQueue[request.ReqWGT])
}