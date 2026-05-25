package config

import "RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig      `mapstructure:",squash"`
	KafkaBrokers           []string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID           string   `mapstructure:"KAFKA_GROUP_ID"`
	OrderCreatedTopic      string   `mapstructure:"KAFKA_ORDER_CREATED_TOPIC"`
	PaymentIntentTopic     string   `mapstructure:"KAFKA_PAYMENT_INTENT_CREATED_TOPIC"`
	PaymentCompletedTopic  string   `mapstructure:"KAFKA_PAYMENT_COMPLETED_TOPIC"`
	PaymentFailedTopic     string   `mapstructure:"KAFKA_PAYMENT_FAILED_TOPIC"`
	PaymentRefundedTopic   string   `mapstructure:"KAFKA_PAYMENT_REFUNDED_TOPIC"`
	StockReservedTopic     string   `mapstructure:"KAFKA_STOCK_RESERVED_TOPIC"`
	StockFailedTopic       string   `mapstructure:"KAFKA_STOCK_FAILED_TOPIC"`
	ShipmentAssignedTopic  string   `mapstructure:"KAFKA_SHIPMENT_ASSIGNED_TOPIC"`
	ShipmentDeliveredTopic string   `mapstructure:"KAFKA_SHIPMENT_DELIVERED_TOPIC"`
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
