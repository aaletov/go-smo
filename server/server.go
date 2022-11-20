package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aaletov/go-smo/api"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/system"
)

type goSmoServer struct {
	system *system.System
}

func NewServer() api.ServerInterface {
	clock.InitClock(time.Unix(0, 0))
	sourcesLambda := time.Duration(1e9 * 11)
	devA := time.Duration(1e10)
	devB := time.Duration(1.4 * 1e10)
	sys := system.NewSystem(3, 4, 3, sourcesLambda, devA, devB)

	for i := 0; i < 20; i++ {
		sys.Iterate()
		sys.PrintData()
	}

	return goSmoServer{sys}
}

func (g goSmoServer) GetWaveNumber(w http.ResponseWriter, r *http.Request) {
	apiSources := make([]api.APISource, 0)
	for _, s := range g.system.Sources {
		generated := s.GetGenerated()

		for _, g := range generated {
			fmt.Println("Source " + g.String())
		}

		num := s.GetNumber()
		apiSources = append(apiSources, api.APISource{
			SourceNum: num,
			Generated: generated,
		})
	}

	apiBuffers := make([]api.APIBuffer, 0)
	for _, b := range g.system.Buffers {
		processed := b.GetAllProcessed()

		for _, g := range processed {
			fmt.Println("Buf " + g.String())
		}

		apiBuffers = append(apiBuffers, api.APIBuffer{
			BufNum:    b.GetNumber(),
			Processed: processed,
			Current:   b.Get(),
		})
	}

	allDone := make([]api.ReqWPT, 0)
	apiDevices := make([]api.APIDevice, 0)
	for _, d := range g.system.Devices {
		done := d.GetDone()
		allDone = append(allDone, done...)

		for _, g := range allDone {
			fmt.Println("Dev " + g.String())
		}

		apiDevices = append(apiDevices, api.APIDevice{
			DevNum:  d.GetNumber(),
			Done:    done,
			Current: d.Get(),
		})
	}

	allRejected := g.system.SetMgr.GetRejectList()

	rv := api.WaveInfo{
		Sources:   apiSources,
		Buffers:   apiBuffers,
		Devices:   apiDevices,
		Done:      allDone,
		Rejected:  allRejected,
		StartTime: clock.SMOClock.StartTime,
		EndTime:   clock.SMOClock.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rv)
	w.WriteHeader(http.StatusOK)
}
