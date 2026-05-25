package config

import "RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig    `mapstructure:",squash"`
	KafkaBrokers         []string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID         string   `mapstructure:"KAFKA_GROUP_ID"`
	TraceTopics          []string `mapstructure:"TRACE_TOPICS"`
	ElasticsearchEnabled bool     `mapstructure:"ELASTICSEARCH_ENABLED"`
	ElasticsearchURL     string   `mapstructure:"ELASTICSEARCH_URL"`
	ElasticsearchIndex   string   `mapstructure:"ELASTICSEARCH_INDEX"`
	InternalSecret       string   `mapstructure:"INTERNAL_SECRET"`
	JWKSURL              string   `mapstructure:"JWKS_URL"`
	JWKSCacheTTL         string   `mapstructure:"JWKS_CACHE_TTL"`
	ExpectedIssuer       string   `mapstructure:"EXPECTED_ISSUER"`
	AuthServiceAddr      string   `mapstructure:"AUTH_SERVICE_ADDR"`
	KafkaPolicyTopic     string   `mapstructure:"KAFKA_AUTH_POLICY_TOPIC"`
}

func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	if err := config.LoadFirstConfig(paths, ".env", &cfg, func() bool { return cfg.DatabaseURL != "" }); err != nil {
		panic(err)
	}
	return cfg
}
