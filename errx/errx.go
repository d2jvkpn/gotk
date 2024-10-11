package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"path/filepath"
	"runtime"
)

type ErrX struct {
	Errors []error `json:"errors"`
	Line   int     `json:"line"`
	Fn     string  `json:"fn"`
	File   string  `json:"file"`

	Code string `json:"code"`
	Kind string `json:"kind"`
	Msg  string `json:"msg"`
}

type Option func(*ErrX)

func NewErrX(e error, options ...Option) (errx *ErrX) {
	errx = &ErrX{Errors: make([]error, 0, 1)}

	if e != nil {
		errx.Errors = append(errx.Errors, e)
	}

	for _, opt := range options {
		opt(errx)
	}

	return errx
}

func Trace(skips ...int) Option {
	if len(skips) == 0 {
		skips = []int{1}
	}

	return func(self *ErrX) {
		self.Trace(skips...)
	}
}

func Code(str string) Option {
	return func(self *ErrX) {
		self.WithCode(str)
	}
}

func Kind(str string) Option {
	return func(self *ErrX) {
		self.WithKind(str)
	}
}

func Msg(str string) Option {
	return func(self *ErrX) {
		self.WithMsg(str)
	}
}

// checks if the input is an ErrX
func ErrXFrom(e error) (errx *ErrX) {
	var ok bool

	if errx, ok = e.(*ErrX); ok {
		return errx
	}

	errx = NewErrX(e)

	return errx
}

func (self *ErrX) Trace(skips ...int) *ErrX {
	var (
		skip int
		pc   uintptr
	)

	if skip = 1; len(skips) > 0 {
		skip = skips[0]
	}

	pc, self.File, self.Line, _ = runtime.Caller(skip)
	self.Fn = filepath.Base(runtime.FuncForPC(pc).Name())

	return self
}

func (self *ErrX) IsNil() bool {
	if self == nil {
		return true
	}

	return len(self.Errors) == 0
}

func (self *ErrX) Error() string {

	if self.IsNil() {
		return "<nil>"
	}

	return errors.Join(self.Errors...).Error()
}

func (self *ErrX) WithErr(errs ...error) *ErrX {
	for i := range errs {
		if errs[i] != nil {
			self.Errors = append(self.Errors, errs[i])
		}
	}

	return self
}

func (self *ErrX) WithCode(str string) *ErrX {
	self.Code = str

	return self
}

func (self *ErrX) WithKind(str string) *ErrX {
	self.Kind = str

	return self
}

func (self *ErrX) WithMsg(str string) *ErrX {
	self.Msg = str
	return self
}

func (self *ErrX) MarshalJSON() ([]byte, error) {
	data := struct {
		Errors []string `json:"errors"`
		Line   int      `json:"line,omitempty"`
		Fn     string   `json:"fn,omitempty"`
		File   string   `json:"file,omitempty"`

		Code string `json:"code"`
		Kind string `json:"kind"`
		Msg  string `json:"msg"`
	}{
		Errors: make([]string, 0, len(self.Errors)),
		Line:   self.Line,
		Fn:     self.Fn,
		File:   self.File,

		Code: self.Code,
		Kind: self.Kind,
		Msg:  self.Msg,
	}

	for _, e := range self.Errors {
		data.Errors = append(data.Errors, fmt.Sprintf("%v", e))
	}

	return json.Marshal(data)
}

func (self *ErrX) CodeKind() (string, string) {
	return self.Code, self.Kind
}

func (self *ErrX) Response() (bts json.RawMessage) {
	data := struct {
		Code string `json:"code"`
		Kind string `json:"kind"`
		Msg  string `json:"msg"`
	}{
		Code: self.Code,
		Kind: self.Kind,
		Msg:  self.Msg,
	}

	bts, _ = json.Marshal(data)
	return bts
}

func (self *ErrX) Debug() (bts json.RawMessage) {
	data := struct {
		Errors []string `json:"errors"`
		Line   int      `json:"line"`
		Fn     string   `json:"fn"`
		File   string   `json:"file"`
	}{
		Errors: make([]string, 0, len(self.Errors)),
		Line:   self.Line,
		Fn:     self.Fn,
		File:   self.File,
	}

	for _, e := range self.Errors {
		data.Errors = append(data.Errors, fmt.Sprintf("%+v", e))
	}

	bts, _ = json.Marshal(data)
	return bts
}

func ParallelRun(funcs ...func() error) (errx *ErrX) {
	var (
		hasErr, ok bool
		wg         sync.WaitGroup
		errs       []error
	)

	errs = make([]error, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}
	wg.Wait()

	hasErr = false
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		hasErr = true

		if errx != nil {
			if errx, ok = errs[i].(*ErrX); ok {
				errs[i] = nil
			}
		}
	}

	if !hasErr {
		return nil
	}

	if errx == nil {
		errx = NewErrX(errors.Join(errs...))
	} else {
		errx.WithErr(errs...)
	}

	return errx
}
