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

type LoadGenerator struct {
	wg      sync.WaitGroup
	context context.Context
	cancel  context.CancelFunc
}

func NewLoadGenerator() *LoadGenerator {
	res := &LoadGenerator{}
	res.wg = sync.WaitGroup{}
	res.context, res.cancel = context.WithCancel(context.Background())
	return res
}

func (gen *LoadGenerator) GenerateCPULoad(w http.ResponseWriter, req *http.Request) {
	load := 100
	duration := time.Duration(10 * time.Second)

	if v := req.URL.Query().Get("load"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			load = l
		}
	}

	if v := req.URL.Query().Get("duration"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
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

func (gen *LoadGenerator) GenerateMemLoad(w http.ResponseWriter, req *http.Request) {
	var size uint64 = 80
	duration := time.Duration(10 * time.Second)

	if v := req.URL.Query().Get("size"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			size = uint64(l)
		}
	}

	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	total := (memStats.Sys - memStats.Mallocs) * size / 100

	if v := req.URL.Query().Get("duration"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			duration = d
		}
	}

	wait := time.After(duration)

	gen.wg.Add(1)
	go func() {
		defer gen.wg.Done()
		buffer := make([]byte, total*1024*1024)
		for i := 0; i < len(buffer); i++ {
			break
		}
		for {
			select {
			case <-gen.context.Done():
				buffer = nil
				return
			case <-wait:
				buffer = nil
				return
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
}

func (gen *LoadGenerator) Wait() {
	gen.cancel()
	gen.wg.Wait()
}
