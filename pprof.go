package gotk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	pp "runtime/pprof"
	"time"
)

/*
web browser address
  http://localhost:8080/debug/pprof/

get profiles and view in browser
  $ go tool pprof -http=:8081 http://localhost:8080/debug/pprof/allocs?seconds=30
  $ go tool pprof http://localhost:8080/debug/pprof/block?seconds=30
  $ go tool pprof http://localhost:8080/debug/pprof/goroutine?seconds=30
  $ go tool pprof http://localhost:8080/debug/pprof/heap?seconds=30
  $ go tool pprof http://localhost:8080/debug/pprof/mutex?seconds=30
  $ go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
  $ go tool pprof http://localhost:8080/debug/pprof/threadcreate?seconds=30

download profile file and convert to svg image
  $ wget -O profile.out localhost:8080/debug/pprof/profile?seconds=30
  $ go tool pprof -svg profile.out > profile.svg

get pprof in 30 seconds and save to svg image
  $ go tool pprof -svg http://localhost:8080/debug/pprof/allocs?seconds=30 > allocs.svg

get trace in 5 seconds
  $ wget -O trace.out http://localhost:8080/debug/pprof/trace?seconds=5
  $ go tool trace trace.out

get cmdline and symbol binary data
  $ wget -O cmdline.out http://localhost:8080/debug/pprof/cmdline
  $ wget -O symbol.out http://localhost:8080/debug/pprof/symbol
*/

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

/*
hz = 100 is recommended

```bash

	go install github.com/google/pprof@latest
	pprof -eog wk_data/2022-12-01T20-38-26_10s_100hz_cpu.pprof.gz
	go tool pprof -eog wk_data/2022-12-01T20-38-26_10s_100hz_cpu.pprof.gz

```
*/
func PprofCollect(dir string, secs, hz int) (out string, err error) {
	var cf, hf1, hf2 *os.File

	if secs <= 0 || hz <= 0 {
		return "", fmt.Errorf("invalid secs or hz")
	}

	out = filepath.Join(
		dir,
		fmt.Sprintf("pprof_%s_%ds_%dhz", time.Now().Format("2006-01-02T15-04-05"), secs, hz),
	)

	if err = os.MkdirAll(out, 0755); err != nil {
		return "", err
	}

	defer func() {
		if cf != nil {
			_ = cf.Close()
		}
		if hf1 != nil {
			_ = hf1.Close()
		}
		if hf2 != nil {
			_ = hf2.Close()
		}
	}()

	if cf, err = os.Create(filepath.Join(out, "cpu.pprof.gz")); err != nil {
		return out, err
	}
	if hf1, err = os.Create(filepath.Join(out, "heap1.pprof.gz")); err != nil {
		return out, err
	}
	if hf2, err = os.Create(filepath.Join(out, "heap2.pprof.gz")); err != nil {
		return out, err
	}

	if err = pp.WriteHeapProfile(hf1); err != nil {
		return out, nil
	}

	runtime.SetCPUProfileRate(hz)
	if err = pp.StartCPUProfile(cf); err != nil {
		return out, err
	}
	defer pp.StopCPUProfile()

	<-time.After(time.Second * time.Duration(secs))

	if err = pp.WriteHeapProfile(hf2); err != nil {
		return out, nil
	}

	return out, nil
}
