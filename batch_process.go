package gotk

import (
	"fmt"
	"sync"
	"time"
)

const (
	BP_StatusReady   = 0
	BP_StatusRunning = 1
	BP_StatusClosed  = -1
)

type BatchProcess[T any] struct {
	count    int
	duration time.Duration
	status   int8
	process  func([]T)

	wg     *sync.WaitGroup
	mutex  *sync.RWMutex
	ticker *time.Ticker
	recv   chan T
	data   []T
}

func NewBatchProcess[T any](count int, duration time.Duration, process func([]T)) (
	bp *BatchProcess[T], err error) {

	if count <= 0 || duration <= 0 {
		return nil, fmt.Errorf("invalid count or duration")
	}

	bp = &BatchProcess[T]{
		count:    count,
		duration: duration,
		status:   BP_StatusReady,
		process:  process,

		wg:     new(sync.WaitGroup),
		mutex:  new(sync.RWMutex),
		ticker: time.NewTicker(duration),
		recv:   make(chan T, count),
		data:   make([]T, 0, count),
	}

	bp.status = BP_StatusRunning
	go bp.run()
	return bp, nil
}

func (bp *BatchProcess[T]) run() {
	process := func() {
		// TODO: ?? sort bp.data
		bp.process(bp.data) // check len(bp.data) == 0 in bp.process
		bp.data = bp.data[:0]
	}

LOOP:
	for {
		select {
		case item := <-bp.recv:
			bp.data = append(bp.data, item)
			if len(bp.data) >= bp.count {
				process()
			}
		case _, isOpen := <-bp.ticker.C:
			if !isOpen {
				break LOOP
			}
			process()
		}
	}

	process() // always process when exit
}

func (bp *BatchProcess[T]) Recv(v T) (err error) {
	bp.mutex.RLock()
	if status := bp.status; status != BP_StatusRunning {
		bp.mutex.RUnlock()
		return fmt.Errorf("unexpected status: %d", status)
	}
	bp.mutex.RUnlock()

	bp.wg.Add(1)
	// go func() { ... } // will cause the app to take up too much memory
	bp.recv <- v
	bp.wg.Done()

	return nil
}

func (bp *BatchProcess[T]) Down() {
	bp.mutex.Lock() // stop recv any data
	bp.status = BP_StatusClosed
	bp.mutex.Unlock()

	bp.wg.Wait()     // sync all goroutines
	bp.ticker.Stop() // close goroutine bp.run
}
