package trace_error

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func fna() (err *Error) {
	e := fmt.Errorf("invalid page number")
	return NewError(e, "bad_request", "invalid_parameter")
}

func fnb() (err *Error) {
	e := fmt.Errorf("unmarshal failed")
	return NewError(e, "bad_request", "unmarshal_failed", Skip(2))
}

func fnc1() (err *Error) {
	e := fmt.Errorf("no user")
	return NewError(e, "bad_request", "not_found", Skip(2))
}

func fnc2() (err *Error) {
	return fnc1()
}

func fnc3() (err *Error) {
	err = fnc1()
	err.Retrace()
	return err
}

func fnc4() (err *Error) {
	return fnc2()
}

func fnc5() (err *Error) {
	return fnc2()
}

func func4() (err *Error) {
	e := fmt.Errorf("an error")
	return NewError(e, "service_unavailable", "", Skip(-1))
}

func Test01(t *testing.T) {
	var err *Error

	err = fna()
	fmt.Printf("==> fna\n%s\n", err.Describe())

	err = fnb()
	fmt.Printf("==> fnb\n%s\n", err.Describe())

	err = fnc2()
	fmt.Printf("==> fnc2\n%s\n", err.Describe())

	err = fnc3()
	fmt.Printf("==> fna3\n%s\n", err.Describe())

	err = func4()
	fmt.Printf("==> fna4\n%s\n", err.Describe())
}

func Test02(t *testing.T) {
	check := func(d any) {
		fmt.Println(d == nil)
	}

	check(nil)
}

func Test03(t *testing.T) {
	var (
		e   error
		err *Error
	)

	e = fmt.Errorf("an error")
	err = NewError(e, "service_unavailable", "")
	fmt.Println(err.Describe())

	err = NewError(e, "service_unavailable", "", Skip(-4))
	fmt.Println(err.Describe())

	fmt.Println(">>> func5")
	err = fnc5()
	fmt.Println(err.Describe())
	err.Retrace()
	fmt.Println(err.Describe())
}

func Test04_AsError(t *testing.T) {
	var err error

	err = anError()
	if err == nil {
		t.Fatal(fmt.Errorf("shouldn't be nil"))
	}

	if e, ok := err.(*Error); !ok {
		t.Fatal(fmt.Errorf("assert as *Error failed"))
	} else {
		fmt.Printf(
			"==> type *Error:\n    string=%s\n    trace=%s\n    is_error=%t, cause=%v, code=%q\n",
			e, e.Trace(), e.IsErr(), e.GetCause(), e.GetCode(),
		)
	}

	if e, ok := err.(Err); !ok {
		t.Fatal(fmt.Errorf("assert as Err failed"))
	} else {
		fmt.Printf(
			"==> interface Err:\n    string=%s\n    is_error=%t, cause=%v, code=%q\n",
			e, e.IsErr(), e.GetCause(), e.GetCode(),
		)
	}
}

func anError() *Error {
	return NewError(errors.New("wrong"), "e0001", "an_error")
}

func TestErrMarshal(t *testing.T) {
	var (
		err error
		bts []byte
		e   error
	)

	err = errors.New("xxxx")
	bts, e = json.Marshal(err)
	if e != nil {
		t.Fatal(e)
	}

	fmt.Printf("==> %s\n", bts)

	err = fmt.Errorf("xxxx")
	bts, e = json.Marshal(err)
	if e != nil {
		t.Fatal(e)
	}

	fmt.Printf("==> %s\n", bts)
}
