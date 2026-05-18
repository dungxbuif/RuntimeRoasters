package config

import "github.com/dungxbuif/RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	InternalSecret    string   `mapstructure:"INTERNAL_SECRET"`
	JWKSURL           string   `mapstructure:"JWKS_URL"`
	JWKSCacheTTL      string   `mapstructure:"JWKS_CACHE_TTL"`
	ExpectedIssuer    string   `mapstructure:"EXPECTED_ISSUER"`
	AuthServiceAddr   string   `mapstructure:"AUTH_SERVICE_ADDR"`
	KafkaBrokers      []string `mapstructure:"KAFKA_BROKERS"`
	KafkaHarvestTopic string   `mapstructure:"KAFKA_HARVEST_TOPIC"`
	KafkaPolicyTopic  string   `mapstructure:"KAFKA_AUTH_POLICY_TOPIC"`
	ValkeyAddr        string   `mapstructure:"VALKEY_ADDR"`
}

// Load loads the farm-service configuration
func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	if err := config.LoadFirstConfig(paths, ".env", &cfg, func() bool { return cfg.DatabaseURL != "" }); err != nil {
		panic(err)
	}
	return cfg
}
