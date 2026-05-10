package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key string, payload interface{}) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) Producer {
	return &producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *producer) Publish(ctx context.Context, topic string, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal kafka payload: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	})

	if err != nil {
		return fmt.Errorf("failed to publish kafka message: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
