package gotk

import (
	// "fmt"
	"errors"
	"sync"
)

func ConcRun(funcs ...func()) {
	var wg sync.WaitGroup

	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			funcs[i]()
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func ConcRunErr(funcs ...func() error) (err error) {
	var (
		wg   sync.WaitGroup
		errs []error
	)

	errs = make([]error, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range errs {
		err = errors.Join(err, errs[i])
	}

	return err
}

// https://go.dev/talks/2013/distsys.slide#38
func ConcRunErrLimit(limit int, funcs ...func() error) (err error) {
	if limit <= 0 || limit > len(funcs) {
		return ConcRunErr(funcs...)
	}

	var (
		n    int
		ch   chan struct{}
		errs []error
	)

	errs = make([]error, len(funcs))
	ch = make(chan struct{}, limit)

	for i := range funcs {
		if n += 1; n > limit {
			_ = <-ch
			n -= 1
		}

		go func(i int) {
			errs[i] = funcs[i]()
			ch <- struct{}{}
		}(i)
	}

	for ; n > 0; n-- {
		_ = <-ch
	}

	for i := range errs {
		err = errors.Join(err, errs[i])
	}

	return err
}
