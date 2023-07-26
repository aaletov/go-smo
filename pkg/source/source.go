package source

import (
	"fmt"
	"time"

	"github.com/aaletov/go-smo/api"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/events"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Source interface {
	Generate() *api.ReqWGT
	GetNumber() int
	GetGenerated() []api.ReqWGT
	GetNextEvent() *events.GenReqEvent
	GetLambda() time.Duration
}

var (
	SourcesCount int = 0
	RandSource       = rand.NewSource(uint64(time.Now().UnixNano()))
)

// Lambda is an expected value of time passed until next req
func NewSource(logger *logrus.Logger, lambda time.Duration) Source {
	SourcesCount++
	ll := logger.WithFields(logrus.Fields{
		"component": fmt.Sprintf("Source #%v", SourcesCount),
	})

	return &sourceImpl{
		logger:        ll,
		sourceNumber:  SourcesCount,
		lastReqNumber: 0,
		nextGenTime:   clock.SMOClock.StartTime,
		lambda:        lambda,
		gen: distuv.Poisson{
			Lambda: float64(lambda),
			Src:    RandSource,
		},
		allGenerated: make([]api.ReqWGT, 0),
	}
}

type sourceImpl struct {
	logger        *logrus.Entry
	sourceNumber  int
	lastReqNumber int
	nextGenTime   time.Time
	lambda        time.Duration
	gen           distuv.Poisson
	allGenerated  []api.ReqWGT
}

func (s *sourceImpl) Generate() *api.ReqWGT {
	ll := s.logger.WithField("method", "Generate")
	s.lastReqNumber++
	req := &api.Request{
		SourceNumber:  s.sourceNumber,
		RequestNumber: s.lastReqNumber,
	}
	genTime := s.nextGenTime
	s.allGenerated = append(s.allGenerated, api.ReqWGT{
		Request: *req,
		Time:    genTime,
	})
	ll.Info("Generated " + req.String())
	return &api.ReqWGT{Request: *req, Time: s.nextGenTime}
}

func (s sourceImpl) GetNumber() int {
	return s.sourceNumber
}

func (s sourceImpl) GetGenerated() []api.ReqWGT {
	return s.allGenerated
}

func (s *sourceImpl) GetNextEvent() *events.GenReqEvent {
	duration := time.Duration(int64(s.gen.Rand()))
	time := s.nextGenTime.Add(duration)
	s.nextGenTime = time

	return &events.GenReqEvent{
		Time:      s.nextGenTime,
		SourceNum: s.sourceNumber,
	}
}

func (s sourceImpl) GetLambda() time.Duration {
	return s.lambda
}
