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
	sourcesLambda := time.Duration(1e12 * 11)
	devDuration := time.Duration(1e10)
	sys := system.NewSystem(3, 4, 3, sourcesLambda, devDuration)

	for i := 0; i < 10; i++ {
		sys.Iterate()
		sys.PrintData()
	}
}
