package app

import (
	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/audit-service/config"
	cassandrastore "github.com/dungxbuif/RuntimeRoasters/apps/audit-service/internal/cassandra"
	"github.com/dungxbuif/RuntimeRoasters/apps/audit-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
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
	var cassandra *cassandrastore.Store
	if cfg.CassandraEnabled {
		store, err := cassandrastore.NewStore(cfg.CassandraHosts, cfg.CassandraKeyspace)
		if err != nil {
			logger.GetLogger().Warn("cassandra audit store unavailable, using postgres fallback", zap.Error(err))
		} else {
			cassandra = store
		}
	}
	service := usecase.NewService(db.DB, cassandra)
	consumers := make([]kafka.Consumer, 0, len(cfg.AuditTopics))
	for _, topic := range cfg.AuditTopics {
		consumers = append(consumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-"+topic, topic))
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
	app := NewApp(baseApp, &cfg, db, service, consumers, guards)
	cleanup := func() {
		guards.Close()
		for _, c := range consumers {
			_ = c.Close()
		}
		if cassandra != nil {
			cassandra.Close()
		}
		if sqlDB, err := db.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	return app, cleanup, nil
}

func defaults(cfg svcconfig.Config) svcconfig.Config {
	if cfg.AppName == "" {
		cfg.AppName = "audit-service"
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
		cfg.AppPort = 8088
	}
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = 50058
	}
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "audit-service-group"
	}
	if cfg.CassandraKeyspace == "" {
		cfg.CassandraKeyspace = "runtime_roasters_audit"
	}
	if len(cfg.AuditTopics) == 0 {
		cfg.AuditTopics = []string{
			events.TopicRetailOrderCreated,
			events.TopicPaymentIntentCreated,
			events.TopicPaymentCompleted,
			events.TopicPaymentFailed,
			events.TopicPaymentRefunded,
			events.TopicWarehouseStockReserved,
			events.TopicWarehouseStockReservationFailed,
			events.TopicWarehouseStockUpdated,
			events.TopicLogisticsShipmentAssigned,
			events.TopicLogisticsShipmentDelivered,
			events.TopicLogisticsGPSUpdated,
		}
	}
	return cfg
}
