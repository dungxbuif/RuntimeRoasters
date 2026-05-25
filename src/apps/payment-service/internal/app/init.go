package app

import (
	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/payment-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/provider"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/usecase"
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
	providers := provider.NewFactory(cfg.StripeWebhookSecret, cfg.VNPayWebhookSecret)
	service := usecase.NewService(
		db.DB,
		producer,
		providers,
		cfg.DefaultProvider,
		cfg.DefaultCurrency,
		cfg.BackfillOrders,
		cfg.IntentCreatedTopic,
		cfg.PaymentCompletedTopic,
		cfg.PaymentFailedTopic,
		cfg.PaymentRefundedTopic,
	)
	orderConsumers := []kafka.Consumer{
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-orders", cfg.OrderCreatedTopic),
	}
	compensationConsumers := []kafka.Consumer{
		kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-stock-failed", cfg.StockFailedTopic),
	}
	if cfg.ShipmentFailedTopic != "" {
		compensationConsumers = append(compensationConsumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-shipment-failed", cfg.ShipmentFailedTopic))
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
	app := NewApp(baseApp, &cfg, db, service, orderConsumers, compensationConsumers, guards)
	cleanup := func() {
		guards.Close()
		_ = producer.Close()
		for _, c := range orderConsumers {
			_ = c.Close()
		}
		for _, c := range compensationConsumers {
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
		cfg.AppName = "payment-service"
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
		cfg.AppPort = 8086
	}
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = 50056
	}
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "payment-service-group"
	}
	if cfg.OrderCreatedTopic == "" {
		cfg.OrderCreatedTopic = events.TopicRetailOrderCreated
	}
	if cfg.StockFailedTopic == "" {
		cfg.StockFailedTopic = events.TopicWarehouseStockReservationFailed
	}
	if cfg.IntentCreatedTopic == "" {
		cfg.IntentCreatedTopic = events.TopicPaymentIntentCreated
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
	if cfg.DefaultProvider == "" {
		cfg.DefaultProvider = provider.ProviderStripe
	}
	if cfg.DefaultCurrency == "" {
		cfg.DefaultCurrency = "USD"
	}
	return cfg
}
