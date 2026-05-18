package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

// Common interfaces and types for the Kafka base package

type Message struct {
	Key     string
	Value   interface{}
	Headers map[string]string
	Topic   string
}

type Handler func(ctx context.Context, msg kafka.Message) error

// Producer defines the interface for publishing messages
type Producer interface {
	Publish(ctx context.Context, topic string, key string, payload interface{}) error
	Close() error
}

// Consumer defines the interface for subscribing to messages
type Consumer interface {
	Listen(ctx context.Context, handler Handler) error
	Close() error
}

// Config holds common Kafka connection settings
type Config struct {
	Brokers []string
	GroupID string
	Topic   string
}

func normalizeBrokers(brokers []string) []string {
	normalized := make([]string, len(brokers))
	for i, broker := range brokers {
		if strings.HasPrefix(broker, "localhost:") {
			normalized[i] = strings.Replace(broker, "localhost:", "127.0.0.1:", 1)
			continue
		}
		normalized[i] = broker
	}
	return normalized
}
