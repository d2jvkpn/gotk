package gotk

import (
	"fmt"
	"os"
	// "io"
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	Writer *lumberjack.Logger
	config zapcore.EncoderConfig
	core   zapcore.Core
	*zap.Logger
}

// if filename is empty, send log to stdout
func NewZapLogger(filename string, level zapcore.LevelEnabler, size_mb int, skips ...int) (
	logger *ZapLogger, err error) {

	var ws zapcore.WriteSyncer

	if filename != "" && size_mb <= 0 {
		return nil, fmt.Errorf("invalid size_mb")
	}

	logger = new(ZapLogger)

	if filename != "" {
		logger.Writer = &lumberjack.Logger{
			Filename:  filename,
			LocalTime: true,
			MaxSize:   size_mb, // megabytes
			// MaxBackups: 3,
			// MaxAge:     1, // days
			// Compress:   true, // disabled by default
		}

		ws = zapcore.AddSync(logger.Writer)
	} else {
		ws = zapcore.AddSync(os.Stdout)
	}

	logger.config = zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		TimeKey:     "time",
		NameKey:     "name",
		CallerKey:   "caller",
		FunctionKey: "func",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		// EncodeTime:   zapcore.RFC3339NanoTimeEncoder,
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000-07:00"),
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// zap.InfoLevel, zapcore.BufferedWriteSyncer
	logger.core = zapcore.NewCore(
		zapcore.NewJSONEncoder(logger.config),
		ws,
		level,
	)

	/*
		// w: io.Writer
		if w != nil {
			consoleEncoder := zapcore.NewConsoleEncoder(logger.config)
			core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(w), level)
			logger.core = zapcore.NewTee(logger.core, core)
		}
	*/

	if len(skips) > 0 {
		logger.Logger = zap.New(logger.core, zap.AddCaller(), zap.AddCallerSkip(skips[0]))
	} else {
		logger.Logger = zap.New(logger.core)
	}

	return logger, nil
}

func (logger *ZapLogger) Down() (err error) {
	var e error

	if logger == nil || logger.Writer == nil {
		return
	}

	if e = logger.Sync(); e != nil {
		err = errors.Join(err, fmt.Errorf("sync: %w", e))
	}

	if e = logger.Writer.Close(); e != nil {
		err = errors.Join(err, fmt.Errorf("close: %w", e))
	}

	return err
}
