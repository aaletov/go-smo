package buffer

import (
	"time"
	"github.com/aaletov/go-smo/pkg/request"
	"github.com/aaletov/go-smo/pkg/queue"
)

type Buffer interface {
	IsFree(moment time.Time) bool
	Add(reqwgt request.ReqWGT) error
	Pop() (request.ReqWGT, error)
}

func NewBuffer() Buffer {
	return &bufferImpl{}
}

type bufferImpl struct {
	queue *queue.PriorityQueue[request.ReqWGT]
}

func (b bufferImpl) IsFree(moment time.Time) bool {
	return false
}

func (b *bufferImpl) Add(reqwgt request.ReqWGT) error {
	return nil
}

func (b *bufferImpl) Pop() (request.ReqWGT, error) {
	return *new(request.ReqWGT), nil
}
