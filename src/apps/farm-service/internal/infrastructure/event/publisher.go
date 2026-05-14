package event

import (
	"context"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
)

type kafkaPublisher struct {
	producer kafka.Producer
}

func NewKafkaPublisher(producer kafka.Producer) usecase.EventPublisher {
	return &kafkaPublisher{producer: producer}
}

func (p *kafkaPublisher) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	// In a real system, the aggregate ID (e.g., harvest_id) should be in the metadata or payload
	// For this demo, we'll try to extract it or use the event ID as key for partitioning
	key := event.ID

	topic := "rr.farm.events"
	fmt.Printf("[KAFKA] Publishing real event: %s to topic %s\n", event.EventType, topic)
	
	return p.producer.Publish(ctx, topic, key, event.Payload)
}

type mockPublisher struct{}

func NewMockPublisher() usecase.EventPublisher {
	return &mockPublisher{}
}

func (m *mockPublisher) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	fmt.Printf("[MOCK PUBLISHER] Published event: %s\n", event.EventType)
	return nil
}
