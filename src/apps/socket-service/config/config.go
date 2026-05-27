package config

import (
	"strings"

	baseconfig "RuntimeRoasters/pkg/config"
	"github.com/spf13/viper"
)

type Config struct {
	BaseConfig           baseconfig.BaseConfig
	KafkaBrokers         []string
	KafkaGroupID         string
	TraceTopics          []string
	ValkeyAddr           string
	InternalAPIKeys      string
	SessionTTL           string
	InternalSecret       string
	JWKSURL              string
	JWKSCacheTTL         string
	ExpectedIssuer       string
	AuthServiceAddr      string
	KafkaPolicyTopic     string
	SocketBroadcastTopic string
}

func Load() Config {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	for _, path := range []string{".", "..", "../..", "../../.."} {
		v.AddConfigPath(path)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = v.ReadInConfig()
	cfg := Config{
		BaseConfig: baseconfig.BaseConfig{
			AppName:      v.GetString("APP_NAME"),
			AppEnv:       v.GetString("APP_ENV"),
			AppPort:      v.GetInt("APP_PORT"),
			GRPCPort:     v.GetInt("GRPC_PORT"),
			LogLevel:     v.GetString("LOG_LEVEL"),
			DBLogLevel:   v.GetString("DB_LOG_LEVEL"),
			InternalHost: v.GetString("INTERNAL_HOST"),
			ValkeyAddr:   v.GetString("VALKEY_ADDR"),
			OTLPEndpoint: v.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
		},
		KafkaBrokers:         csv(v.GetString("KAFKA_BROKERS")),
		KafkaGroupID:         v.GetString("KAFKA_GROUP_ID"),
		TraceTopics:          csv(v.GetString("TRACE_TOPICS")),
		ValkeyAddr:           v.GetString("VALKEY_ADDR"),
		InternalAPIKeys:      v.GetString("SOCKET_INTERNAL_API_KEYS"),
		SessionTTL:           v.GetString("SOCKET_SESSION_TTL"),
		InternalSecret:       v.GetString("INTERNAL_SECRET"),
		JWKSURL:              v.GetString("JWKS_URL"),
		JWKSCacheTTL:         v.GetString("JWKS_CACHE_TTL"),
		ExpectedIssuer:       v.GetString("EXPECTED_ISSUER"),
		AuthServiceAddr:      v.GetString("AUTH_SERVICE_ADDR"),
		KafkaPolicyTopic:     v.GetString("KAFKA_AUTH_POLICY_TOPIC"),
		SocketBroadcastTopic: v.GetString("KAFKA_SOCKET_BROADCAST_TOPIC"),
	}
	return cfg
}

func csv(raw string) []string {
	values := []string{}
	for _, value := range strings.Split(raw, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}
