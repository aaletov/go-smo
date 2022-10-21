package source

import (
	"testing"
	"time"

	"github.com/aaletov/go-smo/pkg/clock"
)

func getSource() Source {
	clock.InitClock(time.Now())
	return NewSource(time.Duration(10))
}

func TestNewSource(t *testing.T) {
	_ = getSource()
}

func TestGetRequest(t *testing.T) {
	source := getSource()
	_ = source.Generate()
}
