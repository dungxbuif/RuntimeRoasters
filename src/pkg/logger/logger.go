package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

type contextKey string

const TraceIDKey contextKey = "trace_id"

// InitLogger initializes a global zap logger based on the environment
func InitLogger(env string, levelStr string) {
	var cfg zap.Config

	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	level, err := zapcore.ParseLevel(levelStr)
	if err == nil {
		cfg.Level.SetLevel(level)
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	Log = logger
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if Log == nil {
		InitLogger("development", "debug")
	}
	return Log
}

// FromContext extracts the TraceID from context and attaches it to the logger
func FromContext(ctx context.Context) *zap.Logger {
	logger := GetLogger()
	traceID, _ := ctx.Value(TraceIDKey).(string)
	if traceID != "" {
		return logger.With(zap.String("trace_id", traceID))
	}
	return logger
}
