package gotk

import (
	"fmt"
	"testing"
)

func TestSigningMd5(t *testing.T) {
	sign, err := NewSigningMd5("xxxxxxxx", "sign", true)
	if err != nil {
		t.Fatal(err)
	}
	query := sign.SignQuery(map[string]string{"a": "1", "b": "zzzz"})

	fmt.Println(">>> query:", query)

	if err := sign.VerifyQuery(query); err != nil {
		t.Fatal(err)
	}
}
