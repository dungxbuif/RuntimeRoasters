package kafka

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	kafkago "github.com/segmentio/kafka-go"
)

func TestMessageIDUsesCloudEventIDWhenPresent(t *testing.T) {
	cloudEvent, err := events.NewCloudEvent(t.Context(), events.TopicRetailOrderCreated, events.SourceRetailService, "orders/order-1", map[string]string{"order_id": "order-1"}, events.Metadata{
		EventID:       "evt-ce-1",
		CorrelationID: "order-1",
		OrderID:       "order-1",
		OccurredAt:    time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(cloudEvent)
	if err != nil {
		t.Fatal(err)
	}
	msg := kafkago.Message{
		Topic: events.TopicRetailOrderCreated,
		Key:   []byte("order-1"),
		Value: body,
	}

	if got := MessageID(msg); got != "retail.order.created:evt-ce-1" {
		t.Fatalf("expected CloudEvent id based message id, got %q", got)
	}
}

func TestMessageIDUsesEventIDWhenPresent(t *testing.T) {
	msg := kafkago.Message{
		Topic: "retail.order.created",
		Key:   []byte("order-1"),
		Value: []byte(`{"event_id":"evt-1","order_id":"order-1"}`),
	}

	if got := MessageID(msg); got != "retail.order.created:evt-1" {
		t.Fatalf("expected event_id based message id, got %q", got)
	}
}

func TestMessageIDFallsBackToTopicKey(t *testing.T) {
	msg := kafkago.Message{
		Topic: "farm.harvest.created",
		Key:   []byte("harvest-1"),
		Value: []byte(`{"harvest_id":"1"}`),
	}

	if got := MessageID(msg); got != "farm.harvest.created:harvest-1" {
		t.Fatalf("expected topic key based message id, got %q", got)
	}
}

func TestMessageIDFallsBackToOffset(t *testing.T) {
	msg := kafkago.Message{
		Topic:     "legacy.topic",
		Partition: 2,
		Offset:    3,
		Value:     []byte(`{}`),
	}

	if got := MessageID(msg); got != "legacy.topic-2-3" {
		t.Fatalf("expected offset based message id, got %q", got)
	}
}
