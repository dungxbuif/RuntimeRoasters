package kafka

import (
	"context"
	"errors"
	"io"
	"time"

	"RuntimeRoasters/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string, topic string) Consumer {
	brokers = normalizeBrokers(brokers)
	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:               brokers,
			GroupID:               groupID,
			Topic:                 topic,
			StartOffset:           kafka.LastOffset,
			MaxWait:               1 * time.Second,
			ReadBackoffMin:        100 * time.Millisecond,
			ReadBackoffMax:        1 * time.Second,
			WatchPartitionChanges: true,
		}),
	}
}

func (c *consumer) Listen(ctx context.Context, handler Handler) error {
	log := logger.GetLogger().With(
		zap.String("component", "kafka-consumer"),
		zap.String("topic", c.Topic()),
	)
	log.Info("consumer listening")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, io.EOF) {
				return nil
			}
			log.Warn("failed to read message", zap.Error(err))
			continue
		}

		headers := kafkaHeadersCarrier(m.Headers)
		msgCtx := otel.GetTextMapPropagator().Extract(ctx, &headers)
		msgCtx, span := otel.Tracer("RuntimeRoasters/pkg/kafka").Start(msgCtx, "kafka.consume "+m.Topic,
			trace.WithSpanKind(trace.SpanKindConsumer),
			trace.WithAttributes(
				attribute.String("messaging.system", "kafka"),
				attribute.String("messaging.destination.name", m.Topic),
				attribute.Int("messaging.kafka.partition", m.Partition),
				attribute.Int64("messaging.kafka.offset", m.Offset),
				attribute.String("messaging.kafka.message.key", string(m.Key)),
			),
		)
		msgLog := logger.FromContext(msgCtx)

		msgLog.Info("message received",
			zap.String("message_topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
			zap.ByteString("key", m.Key),
		)

		if err := handler(msgCtx, m); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			msgLog.Error("failed to handle message", zap.Error(err))
			span.End()
			continue
		}
		span.End()
	}
}

func (c *consumer) Topic() string {
	return c.reader.Config().Topic
}

func (c *consumer) Close() error {
	return c.reader.Close()
}
