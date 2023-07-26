package api

import (
	"strconv"

	"github.com/aaletov/go-smo/pkg/queue"
)

// type requestWithTime api.ReqWT
// type requestWithStartEnd api.ReqWSE
// type Request api.Request

func (r ReqWT) Less(other queue.Comparable) bool {
	otherR := other.(*ReqWGT)
	return r.Time.Before(otherR.Time)
}

// Request with generation time
type ReqWGT = ReqWT

// Request with start time
type ReqWST = ReqWT

// Request with end of processing time
type ReqWPT = ReqWSE

// Request with reject time
type ReqWRT = ReqWT

func (r ReqWT) String() string {
	return r.Request.String() + "WithTime[" + r.Time.String() + "]"
}

func (r ReqWSE) String() string {
	return r.Request.String() + "Start[" + r.Start.String() + "]" +
		"End[" + r.End.String() + "]"
}

func (r Request) String() string {
	return "Req[" + strconv.Itoa(r.SourceNumber) + "." + strconv.Itoa(r.RequestNumber) + "]"
}

func (r Request) Equals(other Request) bool {
	return (r.SourceNumber == other.SourceNumber) && (r.RequestNumber == other.RequestNumber)
}

func (r ReqWSE) Equals(other ReqWSE) bool {
	return r.Request.Equals(other.Request)
}
