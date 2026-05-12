package config

import (
	"fmt"
	"strings"

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
	ValkeyAddr  string `mapstructure:"VALKEY_ADDR"`

	OTLPEndpoint      string  `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	TracingSampleRate float64 `mapstructure:"OTEL_TRACES_SAMPLE_RATE"`
}

// LoadConfig loads configuration from a path into the provided out struct
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
	} else {
		fmt.Printf("Using config file: %s\n", v.ConfigFileUsed())
	}

	err := v.Unmarshal(out)
	if err == nil {
		if cfg, ok := out.(*BaseConfig); ok {
			fmt.Printf("Loaded App: %s, DB: %s, Port: %d, gRPC: %d\n", cfg.AppName, cfg.DatabaseURL, cfg.AppPort, cfg.GRPCPort)
		}
	}
	return err
}
