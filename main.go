package main

import (
	"time"

	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/system"
)

const (
	sourcesLambda = 13
	sourcesCount  = 3
	bufferCount   = 4
	devicesCount  = 3
)

func main() {
	clock.InitClock(time.Now())
	sourcesLambda := time.Duration(1e9 * 11)
	devA := time.Duration(1e11)
	devB := time.Duration(1e12)
	sys := system.NewSystem(3, 4, 3, sourcesLambda, devA, devB)

	for i := 0; i < 10; i++ {
		sys.Iterate()
		sys.PrintData()
	}
}
