package metrics

import (
	// "fmt"
	"testing"
)

func TestMetrics(t *testing.T) {
	_, err := PromMetrics("a", "b")
	if err != nil {
		t.Fatal(err)
	}
}
