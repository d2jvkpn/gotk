/*
#### Usage

```go
package errx

import (

	"github.com/d2jvkpn/gotk/trace_error/v2"

)

type Errorx = trace_error.Error
type ErrorOption = trace_error.ErrorOption

var (

	Msg     = trace_error.Msg
	Skip    = trace_error.Skip
	NoTrace = trace_error.NoTrace

)

```
*/
package trace_error

type ErrorKind struct {
	Err  error  `json:"err"`
	Kind string `json:"kind"`
}

func (self *ErrorKind) IsNil() bool {
	if self == nil || self.Err == nil {
		return true
	}

	return false
}
