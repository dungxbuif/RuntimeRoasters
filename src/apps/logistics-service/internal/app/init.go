package app

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/dungxbuif/RuntimeRoasters/pkg/valkey"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
	logger.GetLogger().Info("running database migrations", zap.String("service", "logistics-service"))
	if err := AutoMigrate(db); err != nil {
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
		cfg.StockUpdatedTopic = events.TopicWarehouseStockUpdated
	}
	if cfg.ShipmentAssignedTopic == "" {
		cfg.ShipmentAssignedTopic = events.TopicLogisticsShipmentAssigned
	}
	if cfg.ShipmentDeliveredTopic == "" {
		cfg.ShipmentDeliveredTopic = events.TopicLogisticsShipmentDelivered
	}
	if cfg.GPSUpdatedTopic == "" {
		cfg.GPSUpdatedTopic = events.TopicLogisticsGPSUpdated
	}
	if cfg.DefaultDestinationStoreID == "" {
		cfg.DefaultDestinationStoreID = "11111111-1111-1111-1111-111111111101"
	}
	producer := kafka.NewProducer(cfg.KafkaBrokers)
	service := usecase.NewService(db.DB, vdb, producer, cfg.ShipmentAssignedTopic, cfg.ShipmentDeliveredTopic, cfg.GPSUpdatedTopic, cfg.DefaultDestinationStoreID)
	if err := service.SeedDrivers(context.Background()); err != nil {
		return nil, nil, err
	}
	consumers := []kafka.Consumer{
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-reserved", cfg.StockReservedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-updated", cfg.StockUpdatedTopic),
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

	app := NewApp(baseApp, &cfg, db, vdb, service, consumers, guards)

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
