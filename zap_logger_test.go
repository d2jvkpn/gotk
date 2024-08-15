package gotk

import (
	// "fmt"
	"testing"

	"go.uber.org/zap"
)

func TestZapLogger(t *testing.T) {
	logger, _ := NewZapLogger("logs/api.log", zap.DebugLevel, 1)
	logger.Info("this is a test", zap.String("hello", "42"))
}
