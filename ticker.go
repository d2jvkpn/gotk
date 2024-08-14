package gotk

import (
	// "fmt"
	"time"
)

type Ticker struct {
	status   int // 0=init, 1=running, 2=stopped
	funcs    []func()
	duration time.Duration
	ch       chan struct{}
	ticker   *time.Ticker
}

func NewTicker(funcs []func(), duration time.Duration) *Ticker {
	if len(funcs) <= 0 {
		panic("invalid funcs")
	}
	if duration <= 0 {
		panic("invalid duration")
	}

	return &Ticker{
		status:   0,
		funcs:    funcs,
		duration: duration,
		ch:       make(chan struct{}),
		ticker:   time.NewTicker(duration),
	}
}

func (self *Ticker) Status() int {
	return self.status
}

func (self *Ticker) Start() {
	self.status = 1
	go func() {
		ok := true
		for {
			select {
			case <-self.ch:
				ok = false
			case _, ok = <-self.ticker.C:
			}

			if !ok {
				return
			}
			for _, fn := range self.funcs {
				fn()
			}
		}
	}()

	self.status = 1
}

func (self *Ticker) End() {
	if self.status != 1 {
		return
	}

	self.ticker.Stop()
	self.ch <- struct{}{}
	self.status = 2
}
