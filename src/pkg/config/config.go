package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// ... BaseConfig remains the same ...

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

// LoadConfig loads configuration and validates that all fields are set
func LoadConfig(path string, name string, out any) error {
	v := viper.New()
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

	if err := v.Unmarshal(out); err != nil {
		return err
	}

	return validateConfig(out)
}

// validateConfig ensures no string fields are empty and no int fields are zero (unless allowed)
// This enforces the "No Fallback / Fail Fast" rule.
func validateConfig(out any) error {
	val := reflect.ValueOf(out)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("mapstructure")

		if tag == "" || tag == ",squash" {
			if field.Kind() == reflect.Struct {
				if err := validateConfig(field.Addr().Interface()); err != nil {
					return err
				}
			}
			continue
		}

		// Fail fast if the field is empty (No fallbacks allowed)
		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				return fmt.Errorf("missing required environment variable for field: %s (tag: %s)", fieldType.Name, tag)
			}
		case reflect.Int:
			if field.Int() == 0 {
				return fmt.Errorf("missing required environment variable (or zero value) for field: %s (tag: %s)", fieldType.Name, tag)
			}
		case reflect.Slice:
			if field.Len() == 0 {
				return fmt.Errorf("missing required environment variable (empty slice) for field: %s (tag: %s)", fieldType.Name, tag)
			}
		}
	}
	return nil
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
