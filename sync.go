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
	var wg sync.WaitGroup

	errs := make([]error, len(funcs))
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
