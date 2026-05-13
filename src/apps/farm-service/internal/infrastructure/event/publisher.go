package event

import (
	"context"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
)

type mockPublisher struct{}

func NewMockPublisher() usecase.EventPublisher {
	return &mockPublisher{}
}

func (p *mockPublisher) Publish(ctx context.Context, event *domain.OutboxEvent) error {
	fmt.Printf("[MOCK KAFKA] Publishing event: %s\n", event.EventType)
	fmt.Printf("  ID: %s\n", event.ID)
	fmt.Printf("  Payload: %s\n", string(event.Payload))
	fmt.Printf("  Metadata: %s\n", string(event.Metadata))
	return nil
}
