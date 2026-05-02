package config

import (
	"fmt"
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
	viper.Reset()
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// Ignore if file doesn't exist, rely on env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(out)
	if err == nil {
		// We can't easily print 'out' generically without reflection, but we can check a known field
		if cfg, ok := out.(*BaseConfig); ok {
			fmt.Printf("Loaded DATABASE_URL: %s\n", cfg.DatabaseURL)
		}
	}
	return err
}
