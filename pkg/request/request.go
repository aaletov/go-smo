package request

import (
	"strconv"
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

type requestWithStartEnd struct {
	Req   *Request
	Start time.Time
	End   time.Time
}

func (r requestWithTime) Less(other queue.Comparable) bool {
	otherR := other.(ReqWGT)
	return r.Time.Before(otherR.Time)
}

// Request with generation time
type ReqWGT = requestWithTime

// Request with start time
type ReqWST = requestWithTime

// Request with end of processing time
type ReqWPT requestWithStartEnd

// Request with reject time
type ReqWRT = requestWithTime

type ReqSE = requestWithStartEnd

func (r Request) String() string {
	return "Request[" + strconv.Itoa(r.SourceNumber) + "." + strconv.Itoa(r.RequestNumber) + "]"
}
