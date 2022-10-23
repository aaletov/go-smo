package smgr

import (
	"errors"

	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/queue"
	"github.com/aaletov/go-smo/pkg/request"
	"github.com/aaletov/go-smo/pkg/source"
)

type ReqWGT = request.ReqWGT
type ReqSE = request.ReqSE
type Queue = queue.PriorityQueue[*ReqWGT]
type QueueEl = queue.QueueElement[*ReqWGT]

type SetManager interface {
	Iterate()
	GetRejectList() []ReqSE
}

func NewSetManager(sources []source.Source, buffers []buffer.Buffer) SetManager {
	return &setManagerImpl{
		sources:    sources,
		buffers:    buffers,
		bufPtr:     0,
		reqQueue:   queue.NewPriorityQueue[*ReqWGT](),
		rejectList: make([]ReqSE, 0),
	}
}

type setManagerImpl struct {
	sources    []source.Source
	buffers    []buffer.Buffer
	bufPtr     int
	reqQueue   Queue
	rejectList []ReqSE
}

func (s *setManagerImpl) collectReqs() {
	for _, src := range s.sources {
		s.reqQueue.Add(src.Generate())
	}
}

// isNewest() determines are all other sources have generated their reqs after
// passed req. Returnes false if not all of other sources have generated next
// request
func isNewest(queue Queue, el *QueueEl, sourcesCount int) bool {
	matchedSrc := make([]bool, sourcesCount)
	matchedSrcCount := 0
	for ; el != nil; el = el.Next() {
		src := el.Get().Req.SourceNumber - 1
		if !matchedSrc[src] {
			matchedSrcCount++
			matchedSrc[src] = true
		}
		if matchedSrcCount == sourcesCount {
			return true
		}
	}
	return false
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
	for prevBufPtr := s.bufPtr; prevBufPtr != s.bufPtr; s.movePtr() {
		if s.currentBuf().IsFree() {
			return nil
		}
	}
	return errors.New("No free buffers in system")
}

func (s *setManagerImpl) Iterate() {
	s.collectReqs()

	for el := s.reqQueue.Front(); el != nil; el = el.Next() {
		if !isNewest(s.reqQueue, el, len(s.sources)) {
			break
		}
		rwgt := el.Get()
		err := s.movePtrToFreeBuf()
		if err != nil {
			s.handleReject(rwgt)
		} else {
			s.buffers[s.bufPtr].Add(rwgt)
			s.buffers[s.bufPtr].Pop(rwgt.Time)
			s.movePtr()
		}
	}
}

func (s setManagerImpl) GetRejectList() []ReqSE {
	return s.rejectList
}
