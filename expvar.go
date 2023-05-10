package gotk

import (
	// "fmt"
	"expvar"
	"runtime"
	// "time"
)

func Expvars() {
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// expvar.Publish("timestamp", expvar.Func(func() any {
	// 	return time.Now().Format(time.RFC3339)
	// }))

	// export memstats and cmdline by default
	//	expvar.Publish("memStats", expvar.Func(func() any {
	//		memStats := new(runtime.MemStats)
	//		runtime.ReadMemStats(memStats)
	//		return memStats
	//	}))
}
