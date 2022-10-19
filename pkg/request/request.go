package request

import (
	"time"

	"github.com/aaletov/go-smo/pkg/queue"
)

type Request struct {
	SourceNumber  int
	RequestNumber int
}

type requestWithTime struct {
	Req  *Request
	Time time.Time
}

func (r requestWithTime) Less(other queue.Comparable) bool {
	otherR := other.(ReqWGT)
	return r.Time.Before(otherR.Time)
}

// Request with generation time
type ReqWGT = requestWithTime

// Request with end of processing time
type ReqWPT = requestWithTime
