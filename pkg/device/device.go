package device

import (
	"time"
	"github.com/aaletov/go-smo/pkg/request"
)

type Device interface {
	Process(req *request.Request) request.ReqWPT
}

type deviceImpl struct {
	pTime time.Time
}