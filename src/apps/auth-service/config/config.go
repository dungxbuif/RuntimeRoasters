package config

import "RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	// Kafka Config
	KafkaBrokers     []string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic       string   `mapstructure:"KAFKA_AUTH_TOPIC"`
	KafkaPolicyTopic string   `mapstructure:"KAFKA_AUTH_POLICY_TOPIC"`

	// Ory Kratos Admin
	KratosAdminURL string `mapstructure:"KRATOS_ADMIN_URL"`
	HydraAdminURL  string `mapstructure:"HYDRA_ADMIN_URL"`

	// Security
	InternalSecret string `mapstructure:"INTERNAL_SECRET"`
	JWKSURL        string `mapstructure:"JWKS_URL"`
	JWKSCacheTTL   string `mapstructure:"JWKS_CACHE_TTL"`
	ExpectedIssuer string `mapstructure:"EXPECTED_ISSUER"`

	// Downstream Services for Seeding
	FarmServiceURL      string `mapstructure:"FARM_SERVICE_URL"`
	RetailServiceURL    string `mapstructure:"RETAIL_SERVICE_URL"`
	LogisticsServiceURL string `mapstructure:"LOGISTICS_SERVICE_URL"`
	WarehouseServiceURL string `mapstructure:"WAREHOUSE_SERVICE_URL"`
}

// Load loads the auth-service configuration
func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	if err := config.LoadFirstConfig(paths, ".env", &cfg, func() bool { return cfg.DatabaseURL != "" }); err != nil {
		panic(err)
	}
	return cfg
}
