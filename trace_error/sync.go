package trace_error

import (
	// "fmt"
	"sync"
)

func ConcRun(funcs ...func() *Error) (err *Error) {
	var wg sync.WaitGroup

	errs := make([]*Error, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range errs {
		err = err.Join(errs[i])
	}

	return err
}
