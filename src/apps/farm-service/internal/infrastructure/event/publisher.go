package event

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

type kafkaPublisher struct {
	producer kafka.Producer
	topic    string
}

func NewKafkaPublisher(producer kafka.Producer, topic string) usecase.EventPublisher {
	return &kafkaPublisher{producer: producer, topic: topic}
}

func (p *kafkaPublisher) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	key := event.ID
	logger.FromContext(ctx).Info("publishing kafka event",
		zap.String("topic", p.topic),
		zap.String("event_type", event.EventType),
		zap.String("event_id", event.ID),
	)
	return p.producer.Publish(ctx, p.topic, key, event.Payload)
}

type mockPublisher struct{}

func NewMockPublisher() usecase.EventPublisher {
	return &mockPublisher{}
}

func (m *mockPublisher) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	logger.FromContext(ctx).Info("mock publisher emitted event",
		zap.String("event_type", event.EventType),
		zap.String("event_id", event.ID),
	)
	return nil
}
