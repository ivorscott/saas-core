package log

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ErrInvalidLogPath represents an invalid log destination.
var ErrInvalidLogPath = errors.New("invalid log path")

// NewProductionLogger creates a zap production logger with lumberjack logger.
func NewProductionLogger(logPath string) (*zap.Logger, error) {
	if logPath == "" {
		return nil, ErrInvalidLogPath
	}

	if _, err := os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0600); err != nil {
		return nil, err
	}

	// Create log retention policy with lumberjack logger.
	lj := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1000,
		MaxBackups: 100,
		MaxAge:     30,
		Compress:   true,
	}

	ws := zapcore.AddSync(lj)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, ws, zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	return logger, nil
}
