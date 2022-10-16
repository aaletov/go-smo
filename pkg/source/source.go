package source

import (
	"time"
	"github.com/aaletov/go-smo/pkg/request"
)

type Source interface {
	GetRequest() (*request.Request, time.Time)
}

func NewSource(lambda time.Duration) Source {
	return sourceImpl{
		lambda: lambda,
	}
}

type sourceImpl struct {
	lambda time.Duration
}

func (s sourceImpl) GetRequest() (*request.Request, time.Time) {
	return nil, time.Time{}
}