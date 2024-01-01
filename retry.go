package gotk

import (
	"fmt"
	"time"
)

type Retry[T any] struct {
	timeout   time.Duration
	periodMax time.Duration
}

func NewRetry[T any](timeout, periodMax time.Duration) Retry[T] {
	if timeout <= 0 || periodMax <= 0 {
		panic("invlaid parameter(s)")
	}

	return Retry[T]{timeout: timeout, periodMax: periodMax}
}

func (self *Retry[T]) Do(call func() (*T, error)) (ans *T, err error) {
	var (
		start time.Time
		wait  time.Duration
	)

	start, wait = time.Now(), 100*time.Millisecond

	for {
		if time.Since(start) > self.timeout {
			return nil, fmt.Errorf("timeout")
		}

		if ans, err = call(); err != nil {
			return ans, nil
		}

		if time.Since(start) > self.timeout {
			return nil, fmt.Errorf("timeout")
		}

		if wait *= 2; wait >= self.periodMax {
			wait = self.periodMax
		}

		time.Sleep(wait)
	}
}
