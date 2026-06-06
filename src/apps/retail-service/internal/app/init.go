package app

import (
	svcconfig "RuntimeRoasters/apps/retail-service/config"
	retailgrpc "RuntimeRoasters/apps/retail-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/retail-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
)

func InitializeApp() (*App, func(), error) {
	cfg := defaults(svcconfig.Load())
	logger.InitLogger(cfg.AppEnv, cfg.LogLevel)

	db, err := database.NewPostgres(database.PostgresConfig{URL: cfg.DatabaseURL, LogLevel: cfg.DBLogLevel})
	if err != nil {
		return nil, nil, err
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	svc := usecase.NewService(db.DB, producer, cfg.OrderCreatedTopic)
	systemHandler := retailgrpc.NewSystemHandler(svc)

	consumers := []kafka.Consumer{
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-intent", cfg.PaymentIntentTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-failed", cfg.PaymentFailedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-refunded", cfg.PaymentRefundedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-reserved", cfg.StockReservedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-failed", cfg.StockFailedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-shipment-assigned", cfg.ShipmentAssignedTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-shipment-delivered", cfg.ShipmentDeliveredTopic),
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-driver-returned-to-base", events.TopicLogisticsDriverReturnedToBase),
	}
	for _, topic := range uniqueTopics(cfg.PaymentCompletedTopic, events.TopicPaymentCompleted) {
		consumers = append(consumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-payment-completed-"+topic, topic))
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
	app := NewApp(baseApp, &cfg, db, svc, systemHandler, consumers, guards)
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
		cfg.PaymentCompletedTopic = events.TopicPaymentSimulatedCompleted
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
		cfg.ShipmentAssignedTopic = events.TopicLogisticsDeliveryAssigned
	}
	if cfg.ShipmentDeliveredTopic == "" {
		cfg.ShipmentDeliveredTopic = events.TopicLogisticsDeliveryCompleted
	}
	return cfg
}

func uniqueTopics(topics ...string) []string {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(topics))
	for _, topic := range topics {
		if topic == "" {
			continue
		}
		if _, ok := seen[topic]; ok {
			continue
		}
		seen[topic] = struct{}{}
		unique = append(unique, topic)
	}
	return unique
}
