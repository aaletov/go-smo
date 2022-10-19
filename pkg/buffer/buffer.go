package buffer

import (
	"errors"
	"time"

	"github.com/aaletov/go-smo/pkg/request"
)

type Request = request.Request
type ReqWGT = request.ReqWGT

type Buffer interface {
	IsFree() bool
	Add(req *Request) error
	Pop() (*Request, error)
}

var (
	bufCount int = 0
)

func NewBuffer(procTime time.Duration) Buffer {
	bufCount++
	return &bufferImpl{
		bufNumber: bufCount,
	}
}

type bufferImpl struct {
	bufNumber int
	req       *Request
}

func (b bufferImpl) IsFree() bool {
	return b.req == nil
}

func (b *bufferImpl) Add(req *Request) error {
	if req != nil {
		return errors.New("Buffer is busy")
	}
	b.req = req
	return nil
}

func (b *bufferImpl) Pop() (*Request, error) {
	if b.req != nil {
		return nil, errors.New("Buffer is empty")
	}
	tmp := b.req
	b.req = nil
	return tmp, nil
}
