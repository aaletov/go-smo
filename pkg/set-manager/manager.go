package smgr

import (
	"github.com/aaletov/go-smo/pkg/source"
	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/queue"
	"github.com/aaletov/go-smo/pkg/request"
)

type SetManager interface {
	Collect()
	ToBuffer()
}

type setManagerImpl struct {
	sources []*source.Source
	buffers []*buffer.Buffer
	requests *queue.PriorityQueue[request.ReqWGT]
}