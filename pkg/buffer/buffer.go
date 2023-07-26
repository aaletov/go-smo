package buffer

import (
	"errors"
	"fmt"
	"time"

	"github.com/aaletov/go-smo/api"
	"github.com/sirupsen/logrus"
)

type Request = api.Request
type ReqWGT = api.ReqWGT
type ReqWSE = api.ReqWSE

type Buffer interface {
	IsFree() bool
	Get() *ReqWGT
	Add(reqwgt *ReqWGT) error
	Pop(popTime time.Time) error
	GetAllProcessed() []ReqWSE
	GetNumber() int
}

var (
	BufCount int = 0
)

func NewBuffer(logger *logrus.Logger) Buffer {
	BufCount++
	ll := logger.WithFields(logrus.Fields{
		"component": fmt.Sprintf("Buffer #%v", BufCount),
	})

	return &bufferImpl{
		logger:       ll,
		bufNumber:    BufCount,
		allProcessed: make([]ReqWSE, 0),
	}
}

type bufferImpl struct {
	logger       *logrus.Entry
	bufNumber    int
	reqwgt       *ReqWGT
	allProcessed []ReqWSE
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
	b.logger.Infof("Added %v", reqwgt.Request.String())
	return nil
}

func (b *bufferImpl) Pop(popTime time.Time) error {
	if b.reqwgt == nil {
		return errors.New("Buffer is empty")
	}
	b.allProcessed = append(b.allProcessed, ReqWSE{
		Start:   b.reqwgt.Time,
		End:     popTime,
		Request: b.reqwgt.Request,
	})
	b.reqwgt = nil
	b.logger.Infof("Popped %v", b.allProcessed[len(b.allProcessed)-1].Request.String())
	return nil
}

func (b bufferImpl) GetAllProcessed() []ReqWSE {
	return b.allProcessed
}

func (b bufferImpl) GetNumber() int {
	return b.bufNumber
}
