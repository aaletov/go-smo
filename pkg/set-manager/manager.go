package smgr

import (
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

func (s *setManagerImpl) reject(rwgt *ReqWGT) {
	prevBufPtr := s.bufPtr
	s.movePtr()
	for ; s.bufPtr != prevBufPtr; s.movePtr() {
		bufRwgt := s.currentBuf().Get()
		if rwgt.Req.SourceNumber < bufRwgt.Req.SourceNumber {
			s.currentBuf().Pop(rwgt.Time)
			s.rejectList = append(s.rejectList, ReqSE{
				Req:   bufRwgt.Req,
				Start: bufRwgt.Time,
				End:   rwgt.Time,
			})
			s.currentBuf().Add(rwgt)
		}
	}
	if s.bufPtr == prevBufPtr {
		s.rejectList = append(s.rejectList, ReqSE{
			Req:   rwgt.Req,
			Start: rwgt.Time,
			End:   rwgt.Time,
		})
	}
}

func (s *setManagerImpl) Iterate() {
	s.collectReqs()

	for el := s.reqQueue.Front(); el != nil; el = el.Next() {
		if !isNewest(s.reqQueue, el, len(s.sources)) {
			break
		}
		prevBufPtr := s.bufPtr
		for !s.buffers[s.bufPtr].IsFree() {
			s.bufPtr = (s.bufPtr + 1) % len(s.sources)
			if s.bufPtr == prevBufPtr {
				reqwgt := el.Get()
				s.reject(reqwgt)
				return
			}
		}
		reqwgt := el.Get()
		s.buffers[s.bufPtr].Add(reqwgt)
		s.buffers[s.bufPtr].Pop(reqwgt.Time)
		s.bufPtr = (s.bufPtr + 1) % len(s.sources)
	}
}

func (s setManagerImpl) GetRejectList() []ReqWRT {
	return s.rejectList
}
