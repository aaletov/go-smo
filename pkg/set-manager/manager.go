package smgr

import (
	"errors"

	"github.com/aaletov/go-smo/api"
	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/events"
	"github.com/aaletov/go-smo/pkg/queue"
	"github.com/aaletov/go-smo/pkg/source"
	"github.com/sirupsen/logrus"
)

type ReqWGT = api.ReqWGT
type ReqWSE = api.ReqWSE

type SetManager interface {
	GetEventFromSource(sourceNum int)
	GetRejectList() []ReqWSE
	ProcessSource(sourceNum int)
}

type Queue = queue.PriorityQueue[events.Event]

func NewSetManager(logger *logrus.Logger, sources []source.Source, buffers []buffer.Buffer, events Queue) SetManager {
	ll := logger.WithFields(logrus.Fields{
		"component": "SetManager",
	})

	return &setManagerImpl{
		logger:     ll,
		sources:    sources,
		buffers:    buffers,
		bufPtr:     0,
		rejectList: make([]ReqWSE, 0),
		events:     events,
	}
}

type setManagerImpl struct {
	logger     *logrus.Entry
	sources    []source.Source
	buffers    []buffer.Buffer
	bufPtr     int
	rejectList []ReqWSE
	events     Queue
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
		s.rejectList = append(s.rejectList, ReqWSE{
			Request: rwgt.Request,
			Start:   rwgt.Time,
			End:     rwgt.Time,
		})
		return
	}
	minPriorRwgt := rwgtSlice[0]
	for _, currRwgt := range rwgtSlice {
		if *currRwgt.Request.SourceNumber > *minPriorRwgt.Request.SourceNumber {
			minPriorRwgt = currRwgt
		}
	}
	if *minPriorRwgt.Request.SourceNumber > *rwgt.Request.SourceNumber {
		// reject minPrior
		bufPtr := rwgtToBufPtr[minPriorRwgt]
		bufRwgt := s.buffers[bufPtr].Get()
		s.buffers[bufPtr].Pop(*rwgt.Time)
		s.rejectList = append(s.rejectList, ReqWSE{
			Request: bufRwgt.Request,
			Start:   bufRwgt.Time,
			End:     rwgt.Time,
		})
		s.buffers[bufPtr].Add(rwgt)
	} else {
		// reject rwgt
		s.rejectList = append(s.rejectList, ReqWSE{
			Request: rwgt.Request,
			Start:   rwgt.Time,
			End:     rwgt.Time,
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

func (s *setManagerImpl) GetEventFromSource(sourceNum int) {
	ll := s.logger.WithField("method", "GetEventFromSource")

	event := s.sources[sourceNum-1].GetNextEvent()
	s.events.Add(event)

	ll.Infof("Got event: %v", event.String())
	ll.Infof("Front queue element: %v", s.events.Front().Get().String())
}

func (s setManagerImpl) GetRejectList() []ReqWSE {
	return s.rejectList
}

func (s *setManagerImpl) ProcessSource(sourceNum int) {
	s.GetEventFromSource(sourceNum)
	rwgt := s.sources[sourceNum-1].Generate()
	err := s.movePtrToFreeBuf()
	if err != nil {
		s.handleReject(rwgt)
	} else {
		s.currentBuf().Add(rwgt)
		s.movePtr()
	}
}
