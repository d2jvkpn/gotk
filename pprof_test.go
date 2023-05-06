package gox

import (
	"testing"
)

func TestPprofCollect(t *testing.T) {
	if _, err := PprofCollect("wk_data", 5, 100); err != nil {
		t.Fatal(err)
	}
}
