package gotk

import (
	"fmt"
	"testing"
	"time"
)

func TestConcRunErrLimit(t *testing.T) {
	var (
		funcs []func() error
		err   error
	)

	funcs = make([]func() error, 10)
	for i := range funcs {
		i := i
		funcs[i] = func() error {
			fmt.Printf("==> FN: %d, %s\n", i, time.Now().Format("15:04:05"))
			time.Sleep(5 * time.Second)
			fmt.Printf("<-- fn: %d\n", i)
			return nil
		}
	}

	err = ConcRunErrLimit(3, funcs...)
	fmt.Println("==> error:", err)
}
