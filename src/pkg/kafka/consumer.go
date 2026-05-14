package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string, topic string) Consumer {
	return &consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:                brokers,
			GroupID:                groupID,
			Topic:                  topic,
			StartOffset:            kafka.FirstOffset,
			MaxWait:                1 * time.Second,
			ReadBackoffMin:         100 * time.Millisecond,
			ReadBackoffMax:         1 * time.Second,
			WatchPartitionChanges: true,
		}),
	}
}

func (c *consumer) Listen(ctx context.Context, handler Handler) error {
	fmt.Printf("[KAFKA] Consumer listening on topic: %s\n", c.reader.Config().Topic)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			fmt.Printf("failed to read message: %v\n", err)
			continue
		}

		fmt.Printf("[KAFKA] Received message: topic=%s partition=%d offset=%d key=%s\n", m.Topic, m.Partition, m.Offset, string(m.Key))

		if err := handler(ctx, m); err != nil {
			fmt.Printf("failed to handle message: %v\n", err)
			continue
		}
	}
}

func (c *consumer) Close() error {
	return c.reader.Close()
}
