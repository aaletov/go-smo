package source

import (
	"fmt"
	"time"

	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/events"
	"github.com/aaletov/go-smo/pkg/request"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Source interface {
	Generate() *request.ReqWGT
	GetNumber() int
	GetGenerated() []request.ReqWGT
	GetNextEvent() *events.GenReqEvent
}

var (
	sourcesCount int = 0
	randSource       = rand.NewSource(uint64(time.Now().UnixNano()))
)

// Lambda is an expected value of time passed until next req
func NewSource(logger *logrus.Logger, lambda time.Duration) Source {
	sourcesCount++
	ll := logger.WithFields(logrus.Fields{
		"component": fmt.Sprintf("Source #%v", sourcesCount),
	})

	return &sourceImpl{
		logger:        ll,
		sourceNumber:  sourcesCount,
		lastReqNumber: 0,
		nextGenTime:   clock.SMOClock.Time,
		lambda:        lambda,
		gen: distuv.Poisson{
			Lambda: float64(lambda),
			Src:    randSource,
		},
		allGenerated: make([]request.ReqWGT, 0),
	}
}

type sourceImpl struct {
	logger        *logrus.Entry
	sourceNumber  int
	lastReqNumber int
	nextGenTime   time.Time
	lambda        time.Duration
	gen           distuv.Poisson
	allGenerated  []request.ReqWGT
}

func (s *sourceImpl) Generate() *request.ReqWGT {
	ll := s.logger.WithField("method", "Generate")
	s.lastReqNumber++
	req := request.Request{
		SourceNumber:  s.sourceNumber,
		RequestNumber: s.lastReqNumber,
	}
	s.allGenerated = append(s.allGenerated, request.ReqWGT{
		Req:  &req,
		Time: s.nextGenTime,
	})
	ll.Info("Generated " + req.String())
	return &request.ReqWGT{Req: &req, Time: s.nextGenTime}
}

func (s sourceImpl) GetNumber() int {
	return s.sourceNumber
}

func (s sourceImpl) GetGenerated() []request.ReqWGT {
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
