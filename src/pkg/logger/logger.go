package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
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
		cfg.EncoderConfig.StacktraceKey = "stacktrace"
		cfg.DisableStacktrace = true
	}

	level, err := zapcore.ParseLevel(levelStr)
	if err == nil {
		cfg.Level.SetLevel(level)
	}

	opts := []zap.Option{}
	if env != "production" {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	logger, err := cfg.Build(opts...)
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
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		return logger.With(zap.String("trace_id", spanContext.TraceID().String()))
	}
	traceID, _ := ctx.Value(TraceIDKey).(string)
	if traceID != "" {
		return logger.With(zap.String("trace_id", traceID))
	}
	return logger
}
