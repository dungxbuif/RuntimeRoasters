package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger khởi tạo logger chuẩn (JSON cho Production, Console cho Development)
func InitLogger(env string, level string) {
	var config zapcore.EncoderConfig
	if env == "production" {
		config = zap.NewProductionEncoderConfig()
	} else {
		config = zap.NewDevelopmentEncoderConfig()
	}

	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.TimeKey = "timestamp"

	var encoder zapcore.Encoder
	if env == "production" {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	logLevel := zap.NewAtomicLevel()
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		logLevel.SetLevel(zap.InfoLevel)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), logLevel)
	Log = zap.New(core, zap.AddCaller())
}

// GetLogger trả về instance logger hiện tại
func GetLogger() *zap.Logger {
	if Log == nil {
		InitLogger("development", "debug")
	}
	return Log
}
