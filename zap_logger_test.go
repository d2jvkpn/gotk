package gotk

import (
	// "fmt"
	"testing"

	"go.uber.org/zap"
)

func TestZapLogger(t *testing.T) {
	var (
		err    error
		logger *ZapLogger
	)

	//
	if logger, err = NewZapLogger("logs/api.log", zap.DebugLevel, 1); err != nil {
		t.Fatal(err)
	}
	logger.Info("this is a test", zap.String("hello", "42"))

	//
	if logger, err = NewZapLogger("", zap.DebugLevel, 1); err != nil {
		t.Fatal(err)
	}
	logger.Info("this is a test", zap.String("hello", "42"))
}
