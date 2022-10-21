package smgr

import (
	"github.com/aaletov/go-smo/pkg/buffer"
	"github.com/aaletov/go-smo/pkg/queue"
	"github.com/aaletov/go-smo/pkg/request"
	"github.com/aaletov/go-smo/pkg/source"
)

type ReqWGT = request.ReqWGT
type ReqWRT = request.ReqWRT
type Queue = queue.PriorityQueue[ReqWGT]
type QueueEl = queue.QueueElement[ReqWGT]

type SetManager interface {
	Collect()
	ToBuffer()
}

func NewSetManager(sources []source.Source, buffers []buffer.Buffer) SetManager {
	return &setManagerImpl{
		sources:    sources,
		buffers:    buffers,
		bufPtr:     0,
		reqQueue:   queue.NewPriorityQueue[ReqWGT](),
		rejectList: make([]ReqWGT, 0),
	}
}

type setManagerImpl struct {
	sources    []source.Source
	buffers    []buffer.Buffer
	bufPtr     int
	reqQueue   Queue
	rejectList []ReqWRT
}

func (s *setManagerImpl) Collect() {
	for _, src := range s.sources {
		req, time := src.GetRequest()
		s.reqQueue.Add(ReqWGT{Req: req, Time: time})
	}
}

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

func (s *setManagerImpl) reject(req *ReqWGT) {
	prevBufPtr := s.bufPtr
	for {
		s.bufPtr = (s.bufPtr + 1) % len(s.sources)
		inBufReq := s.buffers[s.bufPtr].Get()
		if inBufReq.Req.SourceNumber < req.Req.SourceNumber {
			s.buffers[s.bufPtr].Pop(req.Time)
			inBufReq.Time = req.Time
			s.rejectList = append(s.rejectList, *inBufReq)
			s.buffers[s.bufPtr].Add(req)
		}
		if s.bufPtr == prevBufPtr {
			s.rejectList = append(s.rejectList, *req)
			break
		}

	}
}

func (s *setManagerImpl) ToBuffer() {
	for el := s.reqQueue.Front(); el != nil; el = el.Next() {
		if !isNewest(s.reqQueue, el, len(s.sources)) {
			break
		}
		prevBufPtr := s.bufPtr
		for !s.buffers[s.bufPtr].IsFree() {
			s.bufPtr = (s.bufPtr + 1) % len(s.sources)
			if s.bufPtr == prevBufPtr {
				reqwgt := el.Get()
				s.reject(&reqwgt)
				return
			}
		}
		reqwgt := el.Get()
		s.buffers[s.bufPtr].Add(&reqwgt)
		s.bufPtr = (s.bufPtr + 1) % len(s.sources)
	}
}
