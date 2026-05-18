package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// BaseConfig contains fields that every service needs
type BaseConfig struct {
	AppName      string `mapstructure:"APP_NAME"`
	AppEnv       string `mapstructure:"APP_ENV"`
	AppPort      int    `mapstructure:"APP_PORT"`
	GRPCPort     int    `mapstructure:"GRPC_PORT"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
	DBLogLevel   string `mapstructure:"DB_LOG_LEVEL"`
	InternalHost string `mapstructure:"INTERNAL_HOST"`

	DatabaseURL string `mapstructure:"DATABASE_URL"`
	ValkeyAddr  string `mapstructure:"VALKEY_ADDR"`

	OTLPEndpoint      string  `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	TracingSampleRate float64 `mapstructure:"OTEL_TRACES_SAMPLE_RATE"`
}

// LoadConfig loads configuration from a path into the provided out struct.
func LoadConfig(path string, name string, out any) error {
	v := viper.New() // Use a new instance to avoid global state issues
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType("env")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	err := v.Unmarshal(out)
	return err
}

// LoadFirstConfig loads configuration from the first path that satisfies the validator.
func LoadFirstConfig(paths []string, name string, out any, isValid func() bool) error {
	var lastErr error
	for _, path := range paths {
		if err := LoadConfig(path, name, out); err != nil {
			lastErr = err
			continue
		}
		if isValid == nil || isValid() {
			return nil
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return errors.New("no valid configuration found")
}
