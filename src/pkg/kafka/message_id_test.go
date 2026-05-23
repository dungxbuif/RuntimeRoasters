package kafka

import (
	"testing"

	kafkago "github.com/segmentio/kafka-go"
)

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
		Topic: "farm.harvest.events",
		Key:   []byte("harvest-1"),
		Value: []byte(`{"harvest_id":"1"}`),
	}

	if got := MessageID(msg); got != "farm.harvest.events:harvest-1" {
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
