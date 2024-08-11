package trace_error

import (
// "fmt"
)

type Err interface {
	// Error() string
	// IsNil() bool
	IsErr() bool

	GetCause() error
	GetCode() string
	GetKind() string
}
