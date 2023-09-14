package gotk

import (
	"fmt"
	"testing"
	"time"
)

func TestRuntimeInfo(t *testing.T) {
	item := NewRuntimeInfo(func(data map[string]string) {
		fmt.Printf("%#v\n", data)
	}, 5)

	item.Start()
	time.Sleep(30 * time.Second)
	item.End()
}
