package generators

import (
	"context"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type CpuLoadGenerator struct {
	wg      sync.WaitGroup
	context context.Context
	cancel  context.CancelFunc
}

func NewCpuLoadGenerator() *CpuLoadGenerator {
	res := &CpuLoadGenerator{}
	res.wg = sync.WaitGroup{}
	res.context, res.cancel = context.WithCancel(context.Background())
	return res
}

func (gen *CpuLoadGenerator) GenerateCPULoad(w http.ResponseWriter, req *http.Request) {
	load := 100
	duration := time.Duration(10 * time.Second)

	if v := req.URL.Query().Get("load"); v != "" {
		if l, err := strconv.Atoi(v); err != nil {
			load = l
		}
	}

	if v := req.URL.Query().Get("duration"); v != "" {
		if d, err := time.ParseDuration(v); err != nil {
			duration = d
		}
	}

	wait := time.After(duration)

	gen.wg.Add(1)
	go func() {
		defer gen.wg.Done()
		for {
			select {
			case <-gen.context.Done():
				return
			case <-wait:
				return
			default:
				math.Floor(3802920393940938382)
				time.Sleep(time.Duration((100-load)/runtime.NumCPU()) * time.Microsecond)
			}
		}
	}()
}

func (gen *CpuLoadGenerator) Wait() {
	gen.cancel()
	gen.wg.Wait()
}
