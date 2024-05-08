package trace_error

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type Error struct {
	Cause string `json:"cause"`
	Code  string `json:"code"`
	Kind  string `json:"kind"`
	Msg   string `json:"msg"`

	Fn    string `json:"fn,omitempty"`
	File  string `json:"file,omitempty"`
	Line  int    `json:"line,omitempty"`
	cause error
	skip  int
}

type ErrorOption func(*Error)

func Msg(msg string) ErrorOption {
	return func(self *Error) {
		self.Msg = msg
	}
}

func Skip(skip int) ErrorOption {
	return func(self *Error) {
		self.skip = skip
	}
}

func NoTrace() ErrorOption {
	return func(self *Error) {
		self.skip = -1
	}
}

func NewError(cause error, code, kind string, opts ...ErrorOption) (self *Error) {
	if cause == nil {
		return nil
	}

	self = &Error{Cause: cause.Error(), Code: code, Kind: kind, Msg: "", cause: cause, skip: 1}
	for _, opt := range opts {
		opt(self)
	}

	if self.skip <= 0 {
		return self
	}

	fn, file, line, ok := runtime.Caller(self.skip)
	if !ok {
		return self
	}

	self.Line = line
	self.Fn = runtime.FuncForPC(fn).Name()
	self.File = filepath.Base(file)

	return self
}

func (self *Error) Retrace() *Error {
	fn, file, line, ok := runtime.Caller(1)
	if !ok {
		return self
	}

	self.skip = 1
	self.Line = line
	self.Fn = runtime.FuncForPC(fn).Name()
	self.File = filepath.Base(file)

	return self
}

func (self *Error) String() string {
	return fmt.Sprintf(
		"cause=%q, code=%q, kind=%q, msg=%q",
		self.Cause, self.Code, self.Kind, self.Msg,
	)
}

/* don't make Error implments error interface
func (self *Error) Error() string {
	return fmt.Sprintf(
		"cause: %q, code: %q, kind: %q, msg: %q",
		self.Cause, self.Code, self.Kind, self.Msg,
	)
}
*/

func (self *Error) XCause(e error) *Error {
	if e == nil {
		return self
	}

	self.Cause = e.Error()
	self.cause = e
	return self
}

func (err *Error) XMsg(msg string) *Error {
	err.Msg = msg
	return err
}

func (self *Error) XKind(kind string) *Error {
	self.Kind = kind
	return self
}

/*
func (self *Error) String() string {
	return fmt.Sprintf(
		"cause=%q, code=%q, kind=%q, msg=%q",
		self.Cause.Error(), self.Code, self.Kind, self.Msg,
	)
}
*/

func (self *Error) Trace() string {
	if self.Fn == "" {
		return ""
	}

	return fmt.Sprintf(
		"fn=%q, file=%q, line=%d, skip=%d",
		self.Fn, self.File, self.Line, self.skip,
	)
}

func (self *Error) Describe() string {
	str := fmt.Sprintf(
		"cause=%q, code=%q, kind=%q, msg=%q",
		self.Cause, self.Code, self.Kind, self.Msg,
	)

	trace := self.Trace()

	if trace == "" {
		return str
	}
	return fmt.Sprintf("%s, %s", str, trace)
}

func (self *Error) IsNil() bool {
	return self == nil
}

func (self *Error) IsErr() bool {
	return self != nil
}

func (self *Error) GetCause() error {
	return self.cause
}

func (self *Error) GetCode() string {
	return self.Code
}

func (self *Error) GetKind() string {
	return self.Kind
}
