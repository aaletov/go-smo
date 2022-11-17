package server

import (
	"encoding/json"
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
	clock.InitClock(time.Now())
	sourcesLambda := time.Duration(1e9 * 11)
	devA := time.Duration(1e11)
	devB := time.Duration(1e12)
	sys := system.NewSystem(3, 4, 3, sourcesLambda, devA, devB)

	for i := 0; i < 10; i++ {
		sys.Iterate()
		sys.PrintData()
	}

	return goSmoServer{sys}
}

func (g goSmoServer) GetAllBuffers(w http.ResponseWriter,
	r *http.Request) {
	buffers := g.system.Buffers
	apiBuffers := make([]api.APIBuffer, len(buffers))
	for i, b := range buffers {
		num := b.GetNumber()
		apiBuffers[i].BufNum = &num
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiBuffers)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetAllBufProcessedRequests(w http.ResponseWriter,
	r *http.Request,
	params api.GetAllBufProcessedRequestsParams) {
	if (params.BufNum < 0) || (params.BufNum > len(g.system.Buffers)) {
		panic("Do nothing")
	}

	buffer := g.system.Buffers[params.BufNum-1]
	requests := buffer.GetAllProcessed()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetAllDevices(w http.ResponseWriter,
	r *http.Request) {
	devices := g.system.Buffers
	apiDevices := make([]api.APIDevice, len(devices))
	for i, b := range devices {
		num := b.GetNumber()
		apiDevices[i].DevNum = &num
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiDevices)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetDeviceDoneRequests(w http.ResponseWriter,
	r *http.Request,
	params api.GetDeviceDoneRequestsParams) {
	if (params.DevNum < 0) || (params.DevNum > len(g.system.Devices)) {
		panic("Do nothing")
	}

	device := g.system.Devices[params.DevNum-1]
	requests := device.GetDone()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetAllDoneRequests(w http.ResponseWriter, r *http.Request) {
	requests := make([]api.ReqWSE, 0)
	for _, d := range g.system.Devices {
		for _, reqWSE := range d.GetDone() {
			requests = append(requests, reqWSE)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetAllRejectedRequests(w http.ResponseWriter, r *http.Request) {
	rejectList := g.system.SetMgr.GetRejectList()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rejectList)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetAllSources(w http.ResponseWriter, r *http.Request) {
	sources := g.system.Sources
	apiSources := make([]api.APISource, len(sources))
	for i, b := range sources {
		num := b.GetNumber()
		apiSources[i].SourceNum = &num
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiSources)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) GetAllGenRequests(w http.ResponseWriter,
	r *http.Request,
	params api.GetAllGenRequestsParams) {
	if (params.SourceNum < 0) || (params.SourceNum > len(g.system.Sources)) {
		panic("Do nothing")
	}

	source := g.system.Sources[params.SourceNum-1]
	requests := source.GetGenerated()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
	w.WriteHeader(http.StatusOK)
}
