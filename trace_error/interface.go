package trace_error

import (
// "fmt"
)

type Err interface {
	Error() string
	// IsNil() bool
	IsErr() bool

	GetCasue() error
	GetCode() string
}
