package trace_error

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Don't alter field CodeStr for predefined Error
type Error struct {
	Cause error `json:"cause"`

	CodeInt int    `json:"codeInt"`
	CodeStr string `json:"codeStr"`
	Msg     string `json:"msg"`

	Skip int    `json:"skip"`
	Fn   string `json:"fn"`
	File string `json:"file"`
	Line int    `json:"line"`
}

// type ErrorOption func(*Error) bool
type ErrorOption func(*Error)

func Msg(msg string) ErrorOption {
	return func(self *Error) {
		self.Msg = msg
	}
}

func Skip(skip int) ErrorOption {
	return func(self *Error) {
		self.Skip = skip
	}
}

func NoTrace() ErrorOption {
	return func(self *Error) {
		self.Skip = -1
	}
}

func NewError(cause error, codeInt int, codeStr string, opts ...ErrorOption) (self *Error) {
	if cause == nil {
		return nil
	}

	self = &Error{Cause: cause, CodeInt: codeInt, CodeStr: codeStr, Msg: "...", Skip: 1}
	for _, opt := range opts {
		opt(self)
	}

	if self.Skip <= 0 {
		return self
	}

	fn, file, line, ok := runtime.Caller(self.Skip)
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

	self.Skip = 1
	self.Line = line
	self.Fn = runtime.FuncForPC(fn).Name()
	self.File = filepath.Base(file)

	return self
}

func (self *Error) Error() string {
	return fmt.Sprintf(
		"cause: %q, code_int: %d, code_str: %q, msg: %q",
		self.Cause.Error(), self.CodeInt, self.CodeStr, self.Msg,
	)
}

func (self *Error) XCause(e error) *Error {
	if e == nil {
		return self
	}
	self.Cause = e
	return self
}

func (err *Error) XMsg(msg string) *Error {
	err.Msg = msg
	return err
}

func (self *Error) XCode(codeInt int) *Error {
	self.CodeInt = codeInt
	return self
}

func (self *Error) String() string {
	return fmt.Sprintf(
		"cause=%q, code_int=%d, code_str=%q, msg=%q",
		self.Cause.Error(), self.CodeInt, self.CodeStr, self.Msg,
	)
}

func (self *Error) Trace() string {
	if self.Fn == "" {
		return ""
	}

	return fmt.Sprintf(
		"fn=%q, file=%q, line=%d, skip=%d",
		self.Fn, self.File, self.Line, self.Skip,
	)
}

func (self *Error) Describe() string {
	str := self.String()
	trace := self.Trace()

	if trace == "" {
		return str
	}
	return fmt.Sprintf("%s; %s", str, trace)
}

func (self *Error) IsNil() bool {
	return self == nil
}

func (self *Error) IsErr() bool {
	return self != nil
}

func (self *Error) GetCause() error {
	return self.Cause
}

func (self *Error) GetCode() string {
	return self.CodeStr
}
