package gotk

import (
	"testing"
)

func TestLogger(t *testing.T) {
	var logger Logger

	logger = NewDefaultLogger(nil, true)

	logger.Debug("d2jvkpn called func1")
	logger.Info("hello, world!")
	logger.Warn("a warning message received")
	logger.Error("an errror occured")
}
