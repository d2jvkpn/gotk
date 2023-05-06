package gox

import (
	"fmt"
	"testing"
)

func TestCaller(t *testing.T) {
	at, fn := Caller()
	fmt.Println("~~~", at, fn)
}
