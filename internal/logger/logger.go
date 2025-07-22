// Package logger provides structured logging functionality using zap.
package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// log is the private logger instance
var log *zap.Logger

// Initialize sets up the global logger with the specified log level and stores it globally.
func Initialize(level string) error {
	if log != nil {
		return nil
	}

	logger, err := Create(level)
	if err != nil {
		return err
	}

	log = logger
	return nil
}

// Create sets up a logger with the specified log level and returns it without storing globally.
func Create(level string) (*zap.Logger, error) {
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// WithPackage returns a logger with package context.
// If the global logger is not initialized, it creates a default one.
func WithPackage(packageName string) *zap.Logger {
	if log == nil {
		// Fallback: create a default logger if not initialized
		defaultLogger, err := Create("debug")
		if err != nil {
			// Last resort: use zap's no-op logger
			return zap.NewNop()
		}
		log = defaultLogger
	}
	return log.With(zap.String("package", packageName))
}
