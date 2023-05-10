package gotk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"
)

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

	if err = pprof.WriteHeapProfile(hf1); err != nil {
		return out, nil
	}

	runtime.SetCPUProfileRate(hz)
	if err = pprof.StartCPUProfile(cf); err != nil {
		return out, err
	}
	defer pprof.StopCPUProfile()

	<-time.After(time.Second * time.Duration(secs))

	if err = pprof.WriteHeapProfile(hf2); err != nil {
		return out, nil
	}

	return out, nil
}
