package server

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/aaletov/go-smo/api"
	"github.com/aaletov/go-smo/pkg/clock"
	"github.com/aaletov/go-smo/pkg/system"
	"github.com/labstack/gommon/log"
)

type goSmoServer struct {
}

// var (
// 	sourcesLambda = time.Duration(1e8 * 7)
// 	devA          = time.Duration(1e9)
// 	devB          = time.Duration(1.4 * 1e9)
// )

func NewServer() api.ServerInterface {
	return goSmoServer{}
}

func (g goSmoServer) DoStep(w http.ResponseWriter, r *http.Request) {
	system.SysLock.Lock()
	defer system.SysLock.Unlock()

	system.GlobalSystem.Iterate()
	system.GlobalSystem.PrintData()

	return
}

func (g goSmoServer) InitSystem(w http.ResponseWriter, r *http.Request) {
	//(sourcesCount, buffersCount, devicesCount int,
	//	sourcesLambda, devA, devB time.Duration, stepMode bool) {
	var (
		err           error
		resp          api.SystemParameters
		sourcesLambda time.Duration
		devA          time.Duration
		devB          time.Duration
	)

	if err = json.NewDecoder(r.Body).Decode(&resp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("unable to decode: %v", err)
		return
	}

	if sourcesLambda, err = time.ParseDuration(resp.SourcesLambda); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("unable to parse sourcesLambda: %v", err)
		return
	}

	if devA, err = time.ParseDuration(resp.DevA); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("unable to parse devA: %v", err)
		return
	}

	if devB, err = time.ParseDuration(resp.DevB); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("unable to parse devB: %v", err)
		return
	}

	system.InitSystem(
		resp.SourcesCount,
		resp.BuffersCount,
		resp.DevicesCount,
		sourcesLambda,
		devA,
		devB,
	)

	if resp.StepMode {
		return
	}

	if resp.IterationsCount == nil {
		w.Write([]byte("Iterations count not provided for automatic mode"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i := 0; i < *resp.IterationsCount; i++ {
		func() {
			system.SysLock.Lock()
			defer system.SysLock.Unlock()

			system.GlobalSystem.Iterate()
			system.GlobalSystem.PrintData()
		}()
	}

	return
}

func (g goSmoServer) GetWaveNumber(w http.ResponseWriter, r *http.Request) {
	if system.GlobalSystem == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	apiSources := make([]api.APISource, 0)
	for _, s := range system.GlobalSystem.Sources {
		generated := s.GetGenerated()

		num := s.GetNumber()
		apiSources = append(apiSources, api.APISource{
			SourceNum: num,
			Generated: generated,
		})
	}

	apiBuffers := make([]api.APIBuffer, 0)
	for _, b := range system.GlobalSystem.Buffers {
		processed := b.GetAllProcessed()

		apiBuffers = append(apiBuffers, api.APIBuffer{
			BufNum:    b.GetNumber(),
			Processed: processed,
			Current:   b.Get(),
		})
	}

	allDone := make([]api.ReqWPT, 0)
	apiDevices := make([]api.APIDevice, 0)
	for _, d := range system.GlobalSystem.Devices {
		done := d.GetDone()
		allDone = append(allDone, done...)

		apiDevices = append(apiDevices, api.APIDevice{
			DevNum:  d.GetNumber(),
			Done:    done,
			Current: d.Get(),
		})
	}

	allRejected := system.GlobalSystem.SetMgr.GetRejectList()

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

func (g goSmoServer) GetPivotInfo(w http.ResponseWriter, r *http.Request) {
	if system.GlobalSystem == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	devices := system.GlobalSystem.Devices
	devicesPivotInfo := make([]api.DevicePivotInfo, len(devices))

	start := clock.SMOClock.StartTime
	end := clock.SMOClock.Time
	for i, dev := range devices {
		devicesPivotInfo[i] = api.DevicePivotInfo{
			Name:      fmt.Sprintf("Device #%v", dev.GetNumber()),
			UsageCoef: g.getUsageCoef(dev.GetNumber(), start, end),
		}
	}

	sources := system.GlobalSystem.Sources
	sourcesPivotInfo := make([]api.SourcePivotInfo, len(sources))

	devA := system.GlobalSystem.Devices[0].GetDevA()
	devB := system.GlobalSystem.Devices[0].GetDevB()

	for i, source := range sources {
		sourcesPivotInfo[i] = api.SourcePivotInfo{
			Name:               fmt.Sprintf("Source #%v", source.GetNumber()),
			ReqCount:           len(source.GetGenerated()),
			RejChance:          g.getRejChance(source.GetNumber()),
			ProcTime:           g.getAvgProcTime(source.GetNumber()).String(),
			WaitTime:           g.getAvgWaitTime(source.GetNumber()).String(),
			SysTime:            g.getAvgSysTime(source.GetNumber()).String(),
			WaitTimeDispertion: g.getWaitTimeDispertion(source.GetLambda(), source.GetNumber()).String(),
			ProcTimeDispertion: g.getProcTimeDispertion(devA, devB, source.GetNumber()).String(),
		}
	}

	rv := api.PivotInfo{
		DevicesPivotInfo: devicesPivotInfo,
		SourcesPivotInfo: sourcesPivotInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rv)
	w.WriteHeader(http.StatusOK)
}

func (g goSmoServer) getRejChance(sourceNum int) float64 {
	doneCount := len(getAllBySource(g.getAllDone(), sourceNum))
	rejCount := len(getAllBySource(system.GlobalSystem.SetMgr.GetRejectList(), sourceNum))

	return float64(rejCount) / float64(rejCount+doneCount)
}

func (g goSmoServer) getAvgWaitTime(sourceNum int) time.Duration {
	rwseArr := getAllBySource(g.getAllProcessed(), sourceNum)
	sumDuration := time.Duration(0)
	for _, r := range rwseArr {
		sumDuration += r.End.Sub(r.Start)
	}

	return time.Duration(int64(sumDuration) / int64(len(rwseArr)))
}

func (g goSmoServer) getAvgSysTime(sourceNum int) time.Duration {
	bufArr := getAllBySource(g.getAllProcessed(), sourceNum)
	devArr := getAllBySource(g.getAllDone(), sourceNum)

	getBRForDR := func(dr api.ReqWSE) *api.ReqWSE {
		for _, br := range bufArr {
			if br.Equals(dr) {
				return &br
			}
		}
		return nil
	}
	sumDuration := time.Duration(0)
	for _, dr := range devArr {
		br := getBRForDR(dr)
		if br != nil {
			sumDuration += br.End.Sub(br.Start)
		}
	}

	return time.Duration(int64(sumDuration) / int64(len(devArr)))
}

func (g goSmoServer) getAvgProcTime(sourceNum int) time.Duration {
	rwseArr := getAllBySource(g.getAllDone(), sourceNum)
	sumDuration := time.Duration(0)
	for _, r := range rwseArr {
		sumDuration += r.End.Sub(r.Start)
	}

	return time.Duration(int64(sumDuration) / int64(len(rwseArr)))
}

func (g goSmoServer) getWaitTimeDispertion(lambda time.Duration, sourceNum int) time.Duration {
	var sum float64
	sourceArr := getAllBySource(g.getAllProcessed(), sourceNum)
	for _, rwse := range sourceArr {
		sum += math.Pow(float64(int64(rwse.End.Sub(rwse.Start))-int64(lambda)), 2)
	}

	return time.Duration((sum / float64(len(sourceArr))))
}

func (g goSmoServer) getProcTimeDispertion(a, b time.Duration, sourceNum int) time.Duration {
	var sum float64
	devArr := getAllBySource(g.getAllDone(), sourceNum)
	exp := (b + a) / 2
	for _, rwse := range devArr {
		sum += math.Pow(float64(int64(rwse.End.Sub(rwse.Start))-int64(exp)), 2)
	}

	return time.Duration((sum / float64(len(devArr))))
}

func getAllBySource(rwseArr []api.ReqWSE, sourceNum int) []api.ReqWSE {
	filtered := make([]api.ReqWSE, 0)
	for _, r := range rwseArr {
		if r.Request.SourceNumber == sourceNum {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func (g goSmoServer) getAllProcessed() []api.ReqWSE {
	rwseArr := make([]api.ReqWSE, 0)
	buffers := system.GlobalSystem.Buffers
	for _, b := range buffers {
		rwseArr = append(rwseArr, b.GetAllProcessed()...)
	}

	return rwseArr
}

func (g goSmoServer) getAllDone() []api.ReqWSE {
	rwseArr := make([]api.ReqWSE, 0)
	devices := system.GlobalSystem.Devices
	for _, d := range devices {
		rwseArr = append(rwseArr, d.GetDone()...)
	}

	return rwseArr
}

func (g goSmoServer) getUsageCoef(devNum int, start, end time.Time) float64 {
	rwseArr := system.GlobalSystem.Devices[devNum-1].GetDone()
	var sumDuration time.Duration
	for _, rwse := range rwseArr {
		sumDuration += rwse.End.Sub(rwse.Start)
	}

	return float64(sumDuration) / float64(end.Sub(start))
}
