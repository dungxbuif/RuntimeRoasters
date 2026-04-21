package config

import (
	"github.com/dungxbuif/RuntimeRoasters/pkg/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
}

// Load loads the farm-service configuration
func Load() Config {
	var cfg Config
	// Attempt to load from parent directory assuming we run from src
	if err := config.LoadConfig("..", ".env", &cfg); err != nil {
		// Log errors gracefully if required, typically falling back to env variables works fine.
	}
	return cfg
}
