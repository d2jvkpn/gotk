package gotk

import (
	// "fmt"
	"time"
)

type Ticker struct {
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
		funcs:    funcs,
		duration: duration,
		ch:       make(chan struct{}),
		ticker:   time.NewTicker(duration),
	}
}

func (self *Ticker) Start() {
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
}

func (self *Ticker) End() {
	self.ticker.Stop()
	self.ch <- struct{}{}
}
