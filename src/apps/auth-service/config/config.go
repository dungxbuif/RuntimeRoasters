package config

import (
	"fmt"
	"github.com/dungxbuif/RuntimeRoasters/pkg/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	// Kafka Config
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic   string   `mapstructure:"KAFKA_AUTH_TOPIC"`
}

// Load loads the auth-service configuration
func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	for _, p := range paths {
		fmt.Printf("Trying config path: %s\n", p)
		if err := config.LoadConfig(p, ".env", &cfg); err == nil {
			if cfg.DatabaseURL != "" {
				fmt.Printf("Found valid config at: %s\n", p)
				break
			}
		}
	}
	return cfg
}
