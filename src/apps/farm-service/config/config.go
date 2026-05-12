package config

import (
	"fmt"
	"github.com/dungxbuif/RuntimeRoasters/pkg/config"
)

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	InternalSecret    string `mapstructure:"INTERNAL_SECRET"`
	JWKSURL           string `mapstructure:"JWKS_URL"`
	JWKSCacheTTL      string `mapstructure:"JWKS_CACHE_TTL"`
	ExpectedIssuer    string `mapstructure:"EXPECTED_ISSUER"`
	AuthServiceAddr   string `mapstructure:"AUTH_SERVICE_ADDR"`
	ValkeyAddr        string `mapstructure:"VALKEY_ADDR"`
}

// Load loads the farm-service configuration
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
