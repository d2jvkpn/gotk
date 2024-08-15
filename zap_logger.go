package gotk

import (
	"fmt"
	// "io"

	"go.uber.org/multierr"
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

func NewZapLogger(filename string, level zapcore.LevelEnabler, size_mb int, skips ...int) (
	logger *ZapLogger, err error) {

	if filename == "" || size_mb <= 0 {
		return nil, fmt.Errorf("invalid filename or size_mb")
	}

	logger = new(ZapLogger)

	logger.Writer = &lumberjack.Logger{
		Filename:  filename,
		LocalTime: true,
		MaxSize:   size_mb, // megabytes
		// MaxBackups: 3,
		// MaxAge:     1, // days
		// Compress:   true, // disabled by default
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
		zapcore.AddSync(logger.Writer),
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
	var errors []error

	if logger == nil {
		return
	}

	errors = make([]error, 0, 2)
	if err = logger.Sync(); err != nil {
		errors = append(errors, fmt.Errorf("Logger.Sync: %w", err))
	}

	if err = logger.Writer.Close(); err != nil {
		errors = append(errors, fmt.Errorf("Logger.Writer.Close: %w", err))
	}

	return multierr.Combine(errors...)
}
