package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := otel.Tracer("github.com/dungxbuif/RuntimeRoasters/pkg/kafka").Start(ctx, "kafka.produce "+topic,
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.destination.name", topic),
			attribute.String("messaging.kafka.message.key", key),
		),
	)
	defer span.End()

	value, err := json.Marshal(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	headers := make(kafkaHeadersCarrier, 0, 2)
	otel.GetTextMapPropagator().Inject(ctx, &headers)

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic:   topic,
		Key:     []byte(key),
		Value:   value,
		Headers: headers,
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("failed to publish kafka message: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
