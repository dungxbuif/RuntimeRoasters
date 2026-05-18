package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) Producer {
	brokers = normalizeBrokers(brokers)
	return &producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			Async:                  false, // Sync for reliability in this demo
			RequiredAcks:           kafka.RequireAll,
			WriteTimeout:           10 * time.Second,
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *producer) Publish(ctx context.Context, topic string, key string, payload interface{}) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})

	if err != nil {
		return fmt.Errorf("failed to publish kafka message: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
