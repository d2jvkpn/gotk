package gotk

import (
	"fmt"
	"testing"
)

func TestSender(t *testing.T) {
	sender, err := NewSender("examples/email.yaml", "sender")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", sender)
}
