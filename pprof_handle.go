package gotk

import (
	"encoding/json"
	"expvar"
	// "fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
)

// create new Pprof and run server
func LoadPprof(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)

	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.HandleFunc("/debug/runtime/status", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")

		memStats := new(runtime.MemStats)
		runtime.ReadMemStats(memStats)
		num := runtime.NumGoroutine()

		json.NewEncoder(res).Encode(map[string]any{
			"numGoroutine": num,
			"memStats":     memStats,
		})
	})

	return
}

func expvars() {
	// memStats := new(runtime.MemStats)
	// runtime.ReadMemStats(memStats)

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

/*
web browser address

	http://localhost:8080/debug/pprof/

#### get profiles and view in browser

```bash
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/allocs?seconds=30

go tool pprof http://localhost:8080/debug/pprof/block?seconds=30

go tool pprof http://localhost:8080/debug/pprof/goroutine?seconds=30

go tool pprof http://localhost:8080/debug/pprof/heap?seconds=30

go tool pprof http://localhost:8080/debug/pprof/mutex?seconds=30

go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

go tool pprof http://localhost:8080/debug/pprof/threadcreate?seconds=30
````

#### download profile file and convert to svg image
```bash
wget -O profile.out localhost:8080/debug/pprof/profile?seconds=30

go tool pprof -svg profile.out > profile.svg
```

####get pprof in 30 seconds and save to svg image
```bash
go tool pprof -svg http://localhost:8080/debug/pprof/allocs?seconds=30 > allocs.svg
```

#### get trace in 5 seconds
```bash
wget -O trace.out http://localhost:8080/debug/pprof/trace?seconds=5

go tool trace trace.out
```

#### get cmdline and symbol binary data
```bash
wget -O cmdline.out http://localhost:8080/debug/pprof/cmdline

wget -O symbol.out http://localhost:8080/debug/pprof/symbol
```
*/
func PprofHandlerFuncs() map[string]http.HandlerFunc {
	funcs := make(map[string]http.HandlerFunc, 12)

	expvars()
	funcs["expvar"] = expvar.Handler().ServeHTTP

	funcs[""] = pprof.Index
	funcs["profile"] = pprof.Profile
	funcs["trace"] = pprof.Trace
	funcs["cmdline"] = pprof.Cmdline
	funcs["symbol"] = pprof.Symbol

	for _, v := range []string{
		"allocs", "block", "goroutine", "heap", "mutex", "threadcreate",
	} {
		funcs[v] = pprof.Handler(v).ServeHTTP
	}

	return funcs
}

func PprofFuncKeys() []string {
	return []string{
		"", "allocs", "block", "cmdline", "expvar", "goroutine", "heap", "mutex",
		"profile", "symbol", "threadcreate", "trace",
	}
}
