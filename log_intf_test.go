package gox

import (
	"testing"
)

func TestLogIntf(t *testing.T) {
	var logger LogIntf

	logger = NewDefaultLogger(nil, true)

	logger.Debug("d2jvkpn called func1")
	logger.Info("hello, world!")
	logger.Warn("a warning message received")
	logger.Error("an errror occured")
}
