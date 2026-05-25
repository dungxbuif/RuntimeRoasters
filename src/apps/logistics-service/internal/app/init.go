package app

import (
	"RuntimeRoasters/apps/logistics-service/config"
	logisticsgrpc "RuntimeRoasters/apps/logistics-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/logistics-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/valkey"
	"github.com/redis/go-redis/v9"
)

func connectDB(cfg *config.Config) (*database.DB, error) {
	return database.NewPostgres(database.PostgresConfig{
		URL:      cfg.DatabaseURL,
		LogLevel: cfg.DBLogLevel,
	})
}

func connectValkey(cfg *config.Config) *redis.Client {
	return valkey.NewClient(valkey.Config{
		Addr: cfg.ValkeyAddr,
	})
}

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()
	if cfg.DBLogLevel == "" {
		cfg.DBLogLevel = "warn"
	}

	db, err := connectDB(&cfg)
	if err != nil {
		return nil, nil, err
	}

	vdb := connectValkey(&cfg)
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "logistics-service-group"
	}
	if cfg.StockReservedTopic == "" {
		cfg.StockReservedTopic = events.TopicWarehouseStockReserved
	}
	if cfg.StockUpdatedTopic == "" {
		cfg.StockUpdatedTopic = events.TopicWarehouseInventoryUpdated
	}
	if cfg.PickupRequestedTopic == "" {
		cfg.PickupRequestedTopic = events.TopicWarehousePickupRequested
	}
	if cfg.ShipmentAssignedTopic == "" {
		cfg.ShipmentAssignedTopic = events.TopicLogisticsDeliveryAssigned
	}
	if cfg.ShipmentDeliveredTopic == "" {
		cfg.ShipmentDeliveredTopic = events.TopicLogisticsDeliveryCompleted
	}
	if cfg.GPSUpdatedTopic == "" {
		cfg.GPSUpdatedTopic = events.TopicLogisticsGPSUpdated
	}
	if cfg.DefaultDestinationStoreID == "" {
		cfg.DefaultDestinationStoreID = "11111111-1111-1111-1111-111111111101"
	}
	producer := kafka.NewProducer(cfg.KafkaBrokers)
	service := usecase.NewService(db.DB, vdb, producer, cfg.ShipmentAssignedTopic, cfg.ShipmentDeliveredTopic, cfg.GPSUpdatedTopic, cfg.DefaultDestinationStoreID)
	systemHandler := logisticsgrpc.NewSystemHandler(service)

	consumers := []kafka.Consumer{
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-reserved", cfg.StockReservedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-updated", cfg.StockUpdatedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-pickup-requested", cfg.PickupRequestedTopic),
	}

	guards, err := security.NewHTTPGuards(security.HTTPOptions{
		ServiceName:      "logistics-service",
		JWKSURL:          cfg.JWKSURL,
		InternalSecret:   cfg.InternalSecret,
		JWKSCacheTTL:     cfg.JWKSCacheTTL,
		ExpectedIssuer:   cfg.ExpectedIssuer,
		AuthServiceAddr:  cfg.AuthServiceAddr,
		KafkaBrokers:     cfg.KafkaBrokers,
		KafkaPolicyTopic: cfg.KafkaPolicyTopic,
	})
	if err != nil {
		return nil, nil, err
	}

	baseApp := base.NewApp(base.Options{
		Name:   "logistics-service",
		Config: cfg.BaseConfig,
	})

	app := NewApp(baseApp, &cfg, db, vdb, service, systemHandler, consumers, guards)

	cleanup := func() {
		guards.Close()
		_ = producer.Close()
		for _, consumer := range consumers {
			_ = consumer.Close()
		}
		_ = vdb.Close()
		if sqlDB, sqlErr := db.DB.DB(); sqlErr == nil {
			_ = sqlDB.Close()
		}
	}

	return app, cleanup, nil
}
