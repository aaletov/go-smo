package source

import (
	"time"

	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/request"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Source interface {
	GetRequest() (*request.Request, time.Time)
}

var (
	sourcesCount int32 = 0
	randSource         = rand.NewSource(uint64(time.Now().UnixNano()))
)

// Lambda is an expected value of time passed until next req
func NewSource(lambda time.Duration) Source {
	sourceNumber := sourcesCount + 1
	sourcesCount++
	return &sourceImpl{
		sourceNumber:  sourceNumber,
		lastReqNumber: 0,
		lastGenTime:   clock.SMOClock.Time,
		lambda:        lambda,
		gen: distuv.Poisson{
			Lambda: float64(lambda),
			Src:    randSource,
		},
	}
}

type sourceImpl struct {
	sourceNumber  int32
	lastReqNumber int32
	lastGenTime   time.Time
	lambda        time.Duration
	gen           distuv.Poisson
}

func (s *sourceImpl) GetRequest() (*request.Request, time.Time) {
	duration := time.Duration(int64(s.gen.Rand()))
	time := s.lastGenTime.Add(duration)
	s.lastGenTime = time
	return &request.Request{}, time
}
