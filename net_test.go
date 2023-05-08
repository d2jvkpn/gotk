package gotk

import (
	"fmt"
	"testing"
)

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("IP:", ip)
}

func TestPortFromAddr(t *testing.T) {
	var (
		port int
		err  error
	)

	if port, err = PortFromAddr("abc.ddd:90"); err != nil {
		t.Fatal(err)
	}
	fmt.Println("port:", port)

	if port, err = PortFromAddr("abc.ddd"); err == nil {
		t.Fatal("expected an error") // missing port in address
	}
	fmt.Println("port:", port)
}
