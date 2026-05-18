package config

import "github.com/dungxbuif/RuntimeRoasters/pkg/config"

type Config struct {
	config.BaseConfig `mapstructure:",squash"`
	KafkaBrokers      []string `mapstructure:"KAFKA_BROKERS"`
	KafkaGroupID      string   `mapstructure:"KAFKA_GROUP_ID"`
	KafkaHarvestTopic string   `mapstructure:"KAFKA_HARVEST_TOPIC"`
	KafkaStockTopic   string   `mapstructure:"KAFKA_STOCK_TOPIC"`
}

func Load() Config {
	var cfg Config
	paths := []string{".", "..", "../..", "../../.."}
	if err := config.LoadFirstConfig(paths, ".env", &cfg, func() bool { return cfg.DatabaseURL != "" }); err != nil {
		panic(err)
	}
	return cfg
}
