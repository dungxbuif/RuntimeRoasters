package kafka

import (
	"context"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/segmentio/kafka-go"
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
			StartOffset:           kafka.FirstOffset,
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
		zap.String("topic", c.reader.Config().Topic),
	)
	log.Info("consumer listening")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Warn("failed to read message", zap.Error(err))
			continue
		}

		log.Info("message received",
			zap.String("message_topic", m.Topic),
			zap.Int("partition", m.Partition),
			zap.Int64("offset", m.Offset),
			zap.ByteString("key", m.Key),
		)

		if err := handler(ctx, m); err != nil {
			log.Error("failed to handle message", zap.Error(err))
			continue
		}
	}
}

func (c *consumer) Close() error {
	return c.reader.Close()
}
