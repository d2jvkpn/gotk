package impls

import (
	"fmt"
	"testing"
	"time"
)

func TestParallel(t *testing.T) {
	p := NewParallel(4)

	for i := 0; i < 16; i++ {
		i := i
		p.Do(func() error {
			fmt.Printf(">>> %d, %s\n", i, time.Now().Format(time.RFC3339))
			time.Sleep(time.Second)
			return nil
		}, nil)
	}

	p.Wait()
}
