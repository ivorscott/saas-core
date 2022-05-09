package log

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ErrInvalidLogPath represents an invalid log destination
var ErrInvalidLogPath = errors.New("invalid log path")

func NewLoggerOrPanic(logPath string) (*zap.Logger, func() error) {
	if logPath == "" {
		panic(ErrInvalidLogPath)
	}
	// fail immediately if we cannot log to file
	if _, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0600); err != nil {
		panic(err)
	}

	// log retention policy
	lj := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1000,
		MaxBackups: 100,
		MaxAge:     30,
		Compress:   true,
	}

	// integrate lumberjack logger with zap
	ws := zapcore.AddSync(lj)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, ws, zapcore.DebugLevel),        // log to file
		zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel), // log to stdout
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, logger.Sync
}
