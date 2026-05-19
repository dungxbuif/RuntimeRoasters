package config

import "github.com/dungxbuif/RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	ValkeyAddr        string `mapstructure:"VALKEY_ADDR"`
}

// Load loads the logistics-service configuration
func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	if err := config.LoadFirstConfig(paths, ".env", &cfg, func() bool { return cfg.DatabaseURL != "" }); err != nil {
		panic(err)
	}
	return cfg
}
