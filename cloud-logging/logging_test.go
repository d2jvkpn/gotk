package logging

import (
	// "fmt"
	"testing"

	"go.uber.org/zap"
)

func TestLogging(t *testing.T) {
	logger, _ := NewLogger("logs/api.log", zap.DebugLevel, 1)
	logger.Info("this is a test", zap.String("hello", "42"))
}
