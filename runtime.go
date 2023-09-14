package gotk

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

type RuntimeInfo struct {
	handler func(map[string]string)
	seconds int64
	ch      chan struct{}
	ticker  *time.Ticker
}

func NewRuntimeInfo(handler func(map[string]string), seconds int64) *RuntimeInfo {
	if seconds <= 0 {
		panic("invalid seconds")
	}

	return &RuntimeInfo{
		handler: handler,
		seconds: seconds,
		ch:      make(chan struct{}),
		ticker:  time.NewTicker(time.Duration(seconds) * time.Second),
	}
}

func (item *RuntimeInfo) collect() map[string]string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	data := make(map[string]string)

	uint64ToStr := func(size uint64) string {
		return fmt.Sprintf("%.3fMiB", float64(size/1024.0/1024.0))
	}

	data["num_goroutine"] = strconv.Itoa(runtime.NumGoroutine())

	data["heap_objects"] = strconv.FormatUint(m.HeapObjects, 10)
	data["heap_alloc"] = uint64ToStr(m.HeapAlloc)
	data["total_alloc"] = uint64ToStr(m.TotalAlloc)
	data["stack_inuse"] = uint64ToStr(m.StackInuse)

	return data
}

func (item *RuntimeInfo) Start() {
	go func() {
		ok := true
		for {
			select {
			case <-item.ch:
				ok = false
			case _, ok = <-item.ticker.C:
			}

			if !ok {
				return
			}
			item.handler(item.collect())
		}
	}()
}

func (item *RuntimeInfo) End() {
	item.ticker.Stop()
	item.ch <- struct{}{}
}
