package buffer

import (
	"errors"
	"fmt"
	"time"

	"github.com/aaletov/go-smo/pkg/request"
	"github.com/sirupsen/logrus"
)

type Request = request.Request
type ReqWGT = request.ReqWGT
type ReqSE = request.ReqSE

type Buffer interface {
	IsFree() bool
	Get() *ReqWGT
	Add(reqwgt *ReqWGT) error
	Pop(popTime time.Time) error
	GetAllProcessed() []ReqSE
	GetNumber() int
}

var (
	bufCount int = 0
)

func NewBuffer(logger *logrus.Logger) Buffer {
	bufCount++
	ll := logger.WithFields(logrus.Fields{
		"component": fmt.Sprintf("Buffer #%v", bufCount),
	})

	return &bufferImpl{
		logger:       ll,
		bufNumber:    bufCount,
		allProcessed: make([]ReqSE, 0),
	}
}

type bufferImpl struct {
	logger       *logrus.Entry
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
	b.logger.Infof("Added %v", reqwgt.Req.String())
	return nil
}

func (b *bufferImpl) Pop(popTime time.Time) error {
	if b.reqwgt == nil {
		return errors.New("Buffer is empty")
	}
	b.allProcessed = append(b.allProcessed, ReqSE{
		Start: b.reqwgt.Time,
		End:   popTime,
		Req:   b.reqwgt.Req,
	})
	b.reqwgt = nil
	b.logger.Infof("Popped %v", b.allProcessed[len(b.allProcessed)-1].Req.String())
	return nil
}

func (b bufferImpl) GetAllProcessed() []ReqSE {
	return b.allProcessed
}

func (b bufferImpl) GetNumber() int {
	return b.bufNumber
}
