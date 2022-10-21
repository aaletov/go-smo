package source

import (
	"time"

	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/request"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Source interface {
	Generate() *request.ReqWGT
	GetGenerated() []request.ReqWGT
}

var (
	sourcesCount int = 0
	randSource       = rand.NewSource(uint64(time.Now().UnixNano()))
)

// Lambda is an expected value of time passed until next req
func NewSource(lambda time.Duration) Source {
	sourcesCount++
	return &sourceImpl{
		sourceNumber:  sourcesCount,
		lastReqNumber: 0,
		lastGenTime:   clock.SMOClock.Time,
		lambda:        lambda,
		gen: distuv.Poisson{
			Lambda: float64(lambda),
			Src:    randSource,
		},
		allGenerated: make([]request.ReqWGT, 0),
	}
}

type sourceImpl struct {
	sourceNumber  int
	lastReqNumber int
	lastGenTime   time.Time
	lambda        time.Duration
	gen           distuv.Poisson
	allGenerated  []request.ReqWGT
}

func (s *sourceImpl) Generate() *request.ReqWGT {
	duration := time.Duration(int64(s.gen.Rand()))
	time := s.lastGenTime.Add(duration)
	s.lastGenTime = time
	s.lastReqNumber++
	req := request.Request{
		SourceNumber:  s.sourceNumber,
		RequestNumber: s.lastReqNumber,
	}
	s.allGenerated = append(s.allGenerated, request.ReqWGT{
		Req:  &req,
		Time: time,
	})
	return &request.ReqWGT{Req: &req, Time: time}
}

func (s sourceImpl) GetGenerated() []request.ReqWGT {
	return s.allGenerated
}
