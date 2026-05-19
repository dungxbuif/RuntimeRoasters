package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/retail-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/retail-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
)

func InitializeApp() (*App, func(), error) {
	cfg := defaults(svcconfig.Load())
	logger.InitLogger(cfg.AppEnv, cfg.LogLevel)

	db, err := database.NewPostgres(database.PostgresConfig{URL: cfg.DatabaseURL, LogLevel: cfg.DBLogLevel})
	if err != nil {
		return nil, nil, err
	}
	if err := AutoMigrate(db); err != nil {
		return nil, nil, err
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	svc := usecase.NewService(db.DB, producer, cfg.OrderCreatedTopic)
	if err := svc.SeedStores(context.Background()); err != nil {
		return nil, nil, err
	}

	consumers := []kafka.Consumer{
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-intent", cfg.PaymentIntentTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-completed", cfg.PaymentCompletedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-failed", cfg.PaymentFailedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-refunded", cfg.PaymentRefundedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-reserved", cfg.StockReservedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-failed", cfg.StockFailedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-shipment-assigned", cfg.ShipmentAssignedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-shipment-delivered", cfg.ShipmentDeliveredTopic),
	}

	guards, err := security.NewHTTPGuards(security.HTTPOptions{
		ServiceName:      cfg.AppName,
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

	baseApp := base.NewApp(base.Options{Name: cfg.AppName, Config: cfg.BaseConfig})
	app := NewApp(baseApp, &cfg, db, svc, consumers, guards)
	cleanup := func() {
		guards.Close()
		_ = producer.Close()
		for _, c := range consumers {
			_ = c.Close()
		}
		if sqlDB, err := db.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	return app, cleanup, nil
}

func defaults(cfg svcconfig.Config) svcconfig.Config {
	if cfg.AppName == "" {
		cfg.AppName = "retail-service"
	}
	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.DBLogLevel == "" {
		cfg.DBLogLevel = "warn"
	}
	if cfg.AppPort == 0 {
		cfg.AppPort = 8084
	}
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = 50054
	}
	if cfg.InternalHost == "" {
		cfg.InternalHost = "127.0.0.1"
	}
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "retail-service-group"
	}
	if cfg.OrderCreatedTopic == "" {
		cfg.OrderCreatedTopic = events.TopicRetailOrderCreated
	}
	if cfg.PaymentIntentTopic == "" {
		cfg.PaymentIntentTopic = events.TopicPaymentIntentCreated
	}
	if cfg.PaymentCompletedTopic == "" {
		cfg.PaymentCompletedTopic = events.TopicPaymentCompleted
	}
	if cfg.PaymentFailedTopic == "" {
		cfg.PaymentFailedTopic = events.TopicPaymentFailed
	}
	if cfg.PaymentRefundedTopic == "" {
		cfg.PaymentRefundedTopic = events.TopicPaymentRefunded
	}
	if cfg.StockReservedTopic == "" {
		cfg.StockReservedTopic = events.TopicWarehouseStockReserved
	}
	if cfg.StockFailedTopic == "" {
		cfg.StockFailedTopic = events.TopicWarehouseStockReservationFailed
	}
	if cfg.ShipmentAssignedTopic == "" {
		cfg.ShipmentAssignedTopic = events.TopicLogisticsShipmentAssigned
	}
	if cfg.ShipmentDeliveredTopic == "" {
		cfg.ShipmentDeliveredTopic = events.TopicLogisticsShipmentDelivered
	}
	return cfg
}
