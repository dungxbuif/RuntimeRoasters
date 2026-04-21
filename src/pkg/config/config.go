package config

import (
	"github.com/spf13/viper"
)

// BaseConfig contains fields that every service needs
type BaseConfig struct {
	AppName  string `mapstructure:"APP_NAME"`
	AppEnv   string `mapstructure:"APP_ENV"`
	AppPort  int    `mapstructure:"APP_PORT"`
	GRPCPort int    `mapstructure:"GRPC_PORT"`
	LogLevel string `mapstructure:"LOG_LEVEL"`

	DatabaseURL string `mapstructure:"DATABASE_URL"`
	RedisAddr   string `mapstructure:"REDIS_ADDR"`

	OTLPEndpoint      string  `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	TracingSampleRate float64 `mapstructure:"OTEL_TRACES_SAMPLE_RATE"`
}

// LoadConfig loads configuration from a path into the provided out struct
func LoadConfig(path string, name string, out any) error {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Ignore if file doesn't exist, rely on env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return viper.Unmarshal(out)
}
