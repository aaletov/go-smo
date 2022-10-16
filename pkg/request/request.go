package request

import (
	"time"
	"github.com/aaletov/go-smo/pkg/queue"
)

type Request struct {
	SourceNumber int32
	RequestNumber int32
}

type requestWithTime struct {
	req *Request
	time time.Time
}

func (r requestWithTime) Less(other queue.Comparable) bool {
	otherR := other.(ReqWGT)
	return r.time.Before(otherR.time)
}

// Request with generation time
type ReqWGT = requestWithTime
// Request with end of processing time
type ReqWPT = requestWithTime