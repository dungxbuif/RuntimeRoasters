package kafka

import (
	"encoding/json"
	"fmt"

	"RuntimeRoasters/pkg/events"
	kafkago "github.com/segmentio/kafka-go"
)

func MessageID(msg kafkago.Message) string {
	if event, err := events.ParseCloudEvent(msg.Value); err == nil && event.ID() != "" {
		return fmt.Sprintf("%s:%s", msg.Topic, event.ID())
	}
	var payload struct {
		EventID string `json:"event_id"`
	}
	if err := json.Unmarshal(msg.Value, &payload); err == nil && payload.EventID != "" {
		return fmt.Sprintf("%s:%s", msg.Topic, payload.EventID)
	}
	if len(msg.Key) > 0 {
		return fmt.Sprintf("%s:%s", msg.Topic, string(msg.Key))
	}
	return fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset)
}
