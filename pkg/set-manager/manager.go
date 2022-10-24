package smgr

import (
	"errors"

	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/request"
	"github.com/aaletov/go-smo/pkg/source"
)

type ReqWGT = request.ReqWGT
type ReqSE = request.ReqSE

type SetManager interface {
	GetRejectList() []ReqSE
	ProcessSource(sourceNum int)
}

func NewSetManager(sources []source.Source, buffers []buffer.Buffer) SetManager {
	return &setManagerImpl{
		sources:    sources,
		buffers:    buffers,
		bufPtr:     0,
		rejectList: make([]ReqSE, 0),
	}
}

type setManagerImpl struct {
	sources    []source.Source
	buffers    []buffer.Buffer
	bufPtr     int
	rejectList []ReqSE
}

func (s *setManagerImpl) movePtr() {
	s.bufPtr = (s.bufPtr + 1) % len(s.buffers)
}

func (s setManagerImpl) currentBuf() buffer.Buffer {
	return s.buffers[s.bufPtr]
}

func (s *setManagerImpl) handleReject(rwgt *ReqWGT) {
	rwgtToBufPtr := make(map[*ReqWGT]int)
	rwgtSlice := make([]*ReqWGT, 0)
	for i, b := range s.buffers {
		if !b.IsFree() {
			rwgtToBufPtr[b.Get()] = i
			rwgtSlice = append(rwgtSlice, b.Get())
		}
	}
	if len(rwgtSlice) == 0 {
		s.rejectList = append(s.rejectList, ReqSE{
			Req:   rwgt.Req,
			Start: rwgt.Time,
			End:   rwgt.Time,
		})
		return
	}
	minPriorRwgt := rwgtSlice[0]
	for _, currRwgt := range rwgtSlice {
		if currRwgt.Req.SourceNumber > minPriorRwgt.Req.SourceNumber {
			minPriorRwgt = currRwgt
		}
	}
	if minPriorRwgt.Req.SourceNumber > rwgt.Req.SourceNumber {
		// reject minPrior
		bufPtr := rwgtToBufPtr[minPriorRwgt]
		bufRwgt := s.buffers[bufPtr].Get()
		s.buffers[bufPtr].Pop(rwgt.Time)
		s.rejectList = append(s.rejectList, ReqSE{
			Req:   bufRwgt.Req,
			Start: bufRwgt.Time,
			End:   rwgt.Time,
		})
		s.buffers[bufPtr].Add(rwgt)
	} else {
		// reject rwgt
		s.rejectList = append(s.rejectList, ReqSE{
			Req:   rwgt.Req,
			Start: rwgt.Time,
			End:   rwgt.Time,
		})
	}
}

func (s *setManagerImpl) movePtrToFreeBuf() error {
	if s.currentBuf().IsFree() {
		return nil
	}
	prevBufPtr := s.bufPtr
	s.movePtr()
	for ; prevBufPtr != s.bufPtr; s.movePtr() {
		if s.currentBuf().IsFree() {
			return nil
		}
	}
	return errors.New("No free buffers in system")
}

func (s setManagerImpl) GetRejectList() []ReqSE {
	return s.rejectList
}

func (s *setManagerImpl) ProcessSource(sourceNum int) {
	rwgt := s.sources[sourceNum-1].Generate()
	err := s.movePtrToFreeBuf()
	if err != nil {
		s.handleReject(rwgt)
	} else {
		s.currentBuf().Add(rwgt)
		s.movePtr()
	}
}
