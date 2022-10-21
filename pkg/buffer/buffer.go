package buffer

import (
	"errors"
	"time"

	"github.com/aaletov/go-smo/pkg/request"
)

type Request = request.Request
type ReqWGT = request.ReqWGT
type ReqSE = request.ReqSE

type Buffer interface {
	IsFree() bool
	Get() *ReqWGT
	Add(reqwgt *ReqWGT) error
	Pop(popTime time.Time) error
}

var (
	bufCount int = 0
)

func NewBuffer() Buffer {
	bufCount++
	return &bufferImpl{
		bufNumber:    bufCount,
		allProcessed: make([]ReqSE, 0),
	}
}

type bufferImpl struct {
	bufNumber    int
	reqwgt       *ReqWGT
	allProcessed []ReqSE
}

func (b bufferImpl) IsFree() bool {
	return b.reqwgt == nil
}

func (b bufferImpl) Get() *ReqWGT {
	return b.reqwgt
}

func (b *bufferImpl) Add(reqwgt *ReqWGT) error {
	if b.reqwgt != nil {
		return errors.New("Buffer is busy")
	}
	b.reqwgt = reqwgt
	return nil
}

func (b *bufferImpl) Pop(popTime time.Time) error {
	if b.reqwgt != nil {
		return errors.New("Buffer is empty")
	}
	b.allProcessed = append(b.allProcessed, ReqSE{
		Start: b.reqwgt.Time,
		End:   popTime,
		Req:   b.reqwgt.Req,
	})
	b.reqwgt = nil
	return nil
}
