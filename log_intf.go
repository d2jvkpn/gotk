package gotk

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// customize yourself, no concurrency safety guaranteed
type LogIntf interface {
	// Trace(string, ...any)
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	// Fatal, Panic
}

type DefaultLogger struct {
	w      io.Writer // concurrency safety guaranteed
	withTs bool
}

func NewDefaultLogger(writer io.Writer, withTs bool) (logger *DefaultLogger) {
	logger = &DefaultLogger{withTs: withTs}
	// logger.l = log.New(os.Stdout, prefix, log.Lmsgprefix | log.Lshortfile)
	// log.SetOutput(logger.l)
	if writer != nil {
		logger.w = writer
	} else {
		logger.w = os.Stdout
	}

	return logger
}

func (logger *DefaultLogger) Printf(format string, a ...any) (n int, err error) {
	if logger.withTs {
		t := time.Now().Format(RFC3339ms)
		// bytes.TrimSpace(bts)
		n, err = fmt.Fprintf(logger.w, t+" "+strings.TrimSpace(format)+"\n", a...)
	} else {
		n, err = fmt.Fprintf(logger.w, strings.TrimSpace(format)+"\n", a...) // bytes.TrimSpace(bts)
	}

	return
}

func (logger *DefaultLogger) Trace(format string, a ...any) {
	logger.Printf(" [TRACE] "+format, a...)
}

func (logger *DefaultLogger) Debug(format string, a ...any) {
	logger.Printf(" [Debug] "+format, a...)
}

func (logger *DefaultLogger) Info(format string, a ...any) {
	logger.Printf(" [INFO] "+format, a...)
}

func (logger *DefaultLogger) Warn(format string, a ...any) {
	logger.Printf(" [WARN] "+format, a...)
}

func (logger *DefaultLogger) Error(format string, a ...any) {
	logger.Printf(" [ERROR] "+format, a...)
}
