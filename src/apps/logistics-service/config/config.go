package config

import "github.com/dungxbuif/RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig         `mapstructure:",squash"`
	ValkeyAddr                string   `mapstructure:"VALKEY_ADDR"`
	KafkaBrokers              []string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID              string   `mapstructure:"KAFKA_GROUP_ID"`
	StockReservedTopic        string   `mapstructure:"KAFKA_STOCK_RESERVED_TOPIC"`
	StockUpdatedTopic         string   `mapstructure:"KAFKA_STOCK_TOPIC"`
	PickupRequestedTopic      string   `mapstructure:"KAFKA_PICKUP_REQUESTED_TOPIC"`
	ShipmentAssignedTopic     string   `mapstructure:"KAFKA_SHIPMENT_ASSIGNED_TOPIC"`
	ShipmentDeliveredTopic    string   `mapstructure:"KAFKA_SHIPMENT_DELIVERED_TOPIC"`
	GPSUpdatedTopic           string   `mapstructure:"KAFKA_GPS_UPDATED_TOPIC"`
	DefaultDestinationStoreID string   `mapstructure:"DEFAULT_DESTINATION_STORE_ID"`
	InternalSecret            string   `mapstructure:"INTERNAL_SECRET"`
	JWKSURL                   string   `mapstructure:"JWKS_URL"`
	JWKSCacheTTL              string   `mapstructure:"JWKS_CACHE_TTL"`
	ExpectedIssuer            string   `mapstructure:"EXPECTED_ISSUER"`
	AuthServiceAddr           string   `mapstructure:"AUTH_SERVICE_ADDR"`
	KafkaPolicyTopic          string   `mapstructure:"KAFKA_AUTH_POLICY_TOPIC"`
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
