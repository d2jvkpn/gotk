package impls

import (
	// "fmt"
	"runtime"
)

type Parallel struct {
	i, n int
	done chan struct{}
}

func NewParallel(n int) (p *Parallel) {
	if n <= 0 {
		n = runtime.NumCPU()
	}

	return &Parallel{
		n:    n,
		done: make(chan struct{}, uint(n)),
	}
}

// func (p *Parallel) Do(fn func() error, handle func(error)) {
func (p *Parallel) Do(fn func()) {
	if p.i += 1; p.i > p.n { // avoid creating too many goroutines
		<-p.done
		p.i -= 1
	}

	go func() {
		/*
			if err := fn(); err != nil && handle != nil {
				handle(err)
			}
		*/
		fn()
		p.done <- struct{}{}
	}()
}

func (p *Parallel) Wait() {
	for ; p.i > 0; p.i-- {
		<-p.done
	}
}
