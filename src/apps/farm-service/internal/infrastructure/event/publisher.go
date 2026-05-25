package event

import (
	"context"
	"encoding/json"

	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/apps/farm-service/internal/usecase"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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
	ctx = contextFromOutboxMetadata(ctx, event.Metadata)
	key := event.ID
	logger.FromContext(ctx).Info("publishing kafka event",
		zap.String("topic", p.topic),
		zap.String("event_type", event.EventType),
		zap.String("event_id", event.ID),
	)
	return p.producer.Publish(ctx, p.topic, key, json.RawMessage(event.Payload))
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

func contextFromOutboxMetadata(ctx context.Context, metadata []byte) context.Context {
	if len(metadata) == 0 {
		return ctx
	}
	values := map[string]string{}
	if err := json.Unmarshal(metadata, &values); err != nil {
		return ctx
	}
	carrier := propagation.MapCarrier{}
	if traceparent := values["traceparent"]; traceparent != "" {
		carrier.Set("traceparent", traceparent)
	}
	if tracestate := values["tracestate"]; tracestate != "" {
		carrier.Set("tracestate", tracestate)
	}
	if len(carrier) == 0 {
		return ctx
	}
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}
