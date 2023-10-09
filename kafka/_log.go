package kafka

import (
	"fmt"
	"strings"
	"time"
)

type Logger interface {
	Trace(string, ...any)
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type ConsoleLogger struct{}

func (lg *ConsoleLogger) Trace(msg string, fields ...any) {
	now := time.Now().Format(RFC3339ms)
	msg = strings.TrimSpace(msg)

	if len(fields) > 0 {
		fmt.Printf("~~~ [%s] TRACE: %s; %#v\n", now, msg, fields)
	} else {
		fmt.Printf("~~~ [%s] TRACE: %s\n", now, msg)
	}
}

func (lg *ConsoleLogger) Debug(msg string, fields ...any) {
	now := time.Now().Format(RFC3339ms)
	msg = strings.TrimSpace(msg)

	if len(fields) > 0 {
		fmt.Printf("~~~ [%s] DEBUG: %s; %#v\n", now, msg, fields)
	} else {
		fmt.Printf("~~~ [%s] DEBUG: %s\n", now, msg)
	}
}

func (lg *ConsoleLogger) Info(msg string, fields ...any) {
	now := time.Now().Format(RFC3339ms)
	msg = strings.TrimSpace(msg)

	if len(fields) > 0 {
		fmt.Printf("~~~ [%s] INFO: %s; %#v\n", now, msg, fields)
	} else {
		fmt.Printf("~~~ [%s] INFO: %s\n", now, msg)
	}
}

func (lg *ConsoleLogger) Warn(msg string, fields ...any) {
	now := time.Now().Format(RFC3339ms)
	msg = strings.TrimSpace(msg)

	if len(fields) > 0 {
		fmt.Printf("~~~ [%s] WARN: %s; %#v\n", now, msg, fields)
	} else {
		fmt.Printf("~~~ [%s] WARN: %s\n", now, msg)
	}
}

func (lg *ConsoleLogger) Error(msg string, fields ...any) {
	now := time.Now().Format(RFC3339ms)
	msg = strings.TrimSpace(msg)

	if len(fields) > 0 {
		fmt.Printf("~~~ [%s] ERROR: %s; %#v\n", now, msg, fields)
	} else {
		fmt.Printf("~~~ [%s] ERROR: %s\n", now, msg)
	}
}
