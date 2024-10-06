package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestErrx(t *testing.T) {
	// 1.
	var err error

	err = fmt.Errorf("hello")
	err = fmt.Errorf("an error: %w", err)
	fmt.Printf("==> a1. %v\n", err)
	fmt.Printf("==> a2. %v\n", errors.Unwrap(err))

	e1, e2 := errors.New("hello"), errors.New("world")
	err = errors.Join(e1, e2)
	fmt.Printf("==> a3. %v\n", err)
	fmt.Printf("==> a4. %v\n", errors.Unwrap(err))

	// 2.
	var errx *ErrX

	// errx = new(ErrX)
	errx = NewErrX(errors.New("wrong"))
	errx.WithCode("code42").WithKind("kind42").WithKind("kind_xx")

	fmt.Printf("==> b1. ErrX=%+#v\n", errx)

	fmt.Printf("==> b2. errors=%v\n", errx.Errors)

	// 3.
	errx = Fn01ErrX()
	fmt.Printf("==> c1. errx is nil: %t\n", errx == nil)

	var e error = Fn02ErrX()
	errx, _ = e.(*ErrX)
	fmt.Printf("==> c2. %t, %t\n", e == nil, errx.IsNil())
	// false, true, true

	errx = NewErrX(nil)
	e = errx
	fmt.Printf("==> c3. is_nil=%t, e=%v\n", errx.IsNil(), e)

	// 4.
	var bts []byte

	errx = NewErrX(errors.New("e1"))
	errx.WithErr(errors.New("e2")).WithKind("kind01").Trace()

	fmt.Printf("==> d3. errx=%v\n", errx)

	bts, _ = json.Marshal(errx)
	fmt.Printf("==> d3. json=%s\n", bts)

	err = testBizError(errors.New("account not found")).WithMsg("account not exists")
	errx = ErrXFrom(err)
	errx.WithErr(errors.New("sorry")).WithErr(nil)
	bts, _ = json.Marshal(errx)
	fmt.Printf("==> d4. json=%s\n", bts)

	fmt.Printf("==> d5. respone=%s\n", errx.Response())
	fmt.Printf("==> d5. debug=%s\n", errx.Debug())
}

func Fn01ErrX() (errx *ErrX) {
	return nil
}

func Fn02ErrX() (err error) {
	return Fn01ErrX()
}

func testBizError(e error) (errx *ErrX) {
	return NewErrX(e).Trace(2).WithCode("Biz").WithKind("NotFound")
}
