package errx

import (
	"encoding/json"
	"errors"
	"fmt"

	"path/filepath"
	"runtime"
)

// #### 1. ErrX
type ErrX struct {
	Errors []error `json:"errors"`

	Code string `json:"code,omitempty"`
	Kind string `json:"kind,omitempty"`
	Msg  string `json:"msg,omitempty"`

	Line int    `json:"line,omitempty"`
	Fn   string `json:"fn,omitempty"`
	File string `json:"file,omitempty"`
}

func NewErrX(e error) (errx *ErrX) {
	errx = &ErrX{Errors: make([]error, 0, 1)}

	if e != nil {
		errx.Errors = append(errx.Errors, e)
	}

	return errx
}

// checks if e is an ErrX
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
	return errors.Join(self.Errors...).Error()
}

func (self *ErrX) AddErr(e error) *ErrX {
	if e != nil {
		self.Errors = append(self.Errors, e)
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

		Code string `json:"code,omitempty"`
		Kind string `json:"kind,omitempty"`
		Msg  string `json:"msg,omitempty"`

		Line int    `json:"line,omitempty"`
		Fn   string `json:"fn,omitempty"`
		File string `json:"file,omitempty"`
	}{
		Errors: make([]string, 0, len(self.Errors)),

		Code: self.Code,
		Kind: self.Kind,
		Msg:  self.Msg,

		Line: self.Line,
		Fn:   self.Fn,
		File: self.File,
	}

	for _, e := range self.Errors {
		data.Errors = append(data.Errors, fmt.Sprintf("%v", e))
	}

	return json.Marshal(data)
}

func (self *ErrX) Response() (bts json.RawMessage) {
	data := struct {
		Code string `json:"code,omitempty"`
		Kind string `json:"kind,omitempty"`
		Msg  string `json:"msg,omitempty"`
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

		Line int    `json:"line,omitempty"`
		Fn   string `json:"fn,omitempty"`
		File string `json:"file,omitempty"`
	}{
		Errors: make([]string, 0, len(self.Errors)),

		Line: self.Line,
		Fn:   self.Fn,
		File: self.File,
	}

	for _, e := range self.Errors {
		data.Errors = append(data.Errors, fmt.Sprintf("%v", e))
	}

	bts, _ = json.Marshal(data)
	return bts
}
