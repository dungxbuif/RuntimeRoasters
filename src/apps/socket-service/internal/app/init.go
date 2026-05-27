package app

import (
	"time"

	svcconfig "RuntimeRoasters/apps/socket-service/config"
	"RuntimeRoasters/apps/socket-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/valkey"
)

func InitializeApp() (*App, func(), error) {
	cfg := defaults(svcconfig.Load())
	rdb := valkey.NewClient(valkey.Config{Addr: cfg.ValkeyAddr})
	producer := kafka.NewProducer(cfg.KafkaBrokers)
	ttl, _ := time.ParseDuration(cfg.SessionTTL)
	service := usecase.NewService(rdb, producer, cfg.SocketBroadcastTopic, ttl, cfg.InternalAPIKeys)

	consumers := make([]kafka.Consumer, 0, len(cfg.TraceTopics))
	for _, topic := range cfg.TraceTopics {
		consumers = append(consumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-"+topic, topic))
	}

	guards, err := security.NewHTTPGuards(security.HTTPOptions{
		ServiceName:      cfg.BaseConfig.AppName,
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

	baseApp := base.NewApp(base.Options{Name: cfg.BaseConfig.AppName, Config: cfg.BaseConfig})
	app := NewApp(baseApp, &cfg, rdb, service, producer, consumers, guards)
	cleanup := func() { _ = app.Shutdown() }
	return app, cleanup, nil
}

func defaults(cfg svcconfig.Config) svcconfig.Config {
	if cfg.BaseConfig.AppName == "" {
		cfg.BaseConfig.AppName = "socket-service"
	}
	if cfg.BaseConfig.AppEnv == "" {
		cfg.BaseConfig.AppEnv = "development"
	}
	if cfg.BaseConfig.LogLevel == "" {
		cfg.BaseConfig.LogLevel = "info"
	}
	if cfg.BaseConfig.AppPort == 0 {
		cfg.BaseConfig.AppPort = 8091
	}
	if cfg.BaseConfig.GRPCPort == 0 {
		cfg.BaseConfig.GRPCPort = 50060
	}
	if cfg.BaseConfig.InternalHost == "" {
		cfg.BaseConfig.InternalHost = "127.0.0.1"
	}
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "socket-service-group"
	}
	if cfg.ValkeyAddr == "" {
		cfg.ValkeyAddr = "localhost:6379"
	}
	cfg.BaseConfig.ValkeyAddr = cfg.ValkeyAddr
	if cfg.SessionTTL == "" {
		cfg.SessionTTL = "90s"
	}
	if cfg.SocketBroadcastTopic == "" {
		cfg.SocketBroadcastTopic = events.TopicSocketBroadcastRequested
	}
	if len(cfg.TraceTopics) == 0 {
		cfg.TraceTopics = events.TraceableTopics
	}
	return cfg
}
