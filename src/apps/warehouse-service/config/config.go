package config

import "RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig       `mapstructure:",squash"`
	KafkaBrokers            []string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID            string   `mapstructure:"KAFKA_GROUP_ID"`
	KafkaHarvestTopic       string   `mapstructure:"KAFKA_HARVEST_TOPIC"`
	KafkaStockTopic         string   `mapstructure:"KAFKA_STOCK_TOPIC"`
	KafkaPickupTopic        string   `mapstructure:"KAFKA_PICKUP_REQUESTED_TOPIC"`
	KafkaPickupArrivedTopic string   `mapstructure:"KAFKA_PICKUP_ARRIVED_TOPIC"`
	KafkaIntakeCreatedTopic string   `mapstructure:"KAFKA_INTAKE_CREATED_TOPIC"`
	KafkaNotificationTopic  string   `mapstructure:"KAFKA_NOTIFICATION_CREATED_TOPIC"`
	KafkaOrderTopic         string   `mapstructure:"KAFKA_ORDER_CREATED_TOPIC"`
	KafkaStockReservedTopic string   `mapstructure:"KAFKA_STOCK_RESERVED_TOPIC"`
	KafkaStockFailedTopic   string   `mapstructure:"KAFKA_STOCK_FAILED_TOPIC"`
	DefaultWarehouseID      string   `mapstructure:"DEFAULT_WAREHOUSE_ID"`
	InternalSecret         string   `mapstructure:"INTERNAL_SECRET"`
	JWKSURL                string   `mapstructure:"JWKS_URL"`
	JWKSCacheTTL           string   `mapstructure:"JWKS_CACHE_TTL"`
	ExpectedIssuer         string   `mapstructure:"EXPECTED_ISSUER"`
	AuthServiceAddr        string   `mapstructure:"AUTH_SERVICE_ADDR"`
	KafkaPolicyTopic       string   `mapstructure:"KAFKA_AUTH_POLICY_TOPIC"`
}

func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	if err := config.LoadFirstConfig(paths, ".env", &cfg, func() bool { return cfg.DatabaseURL != "" }); err != nil {
		panic(err)
	}
	return cfg
}
