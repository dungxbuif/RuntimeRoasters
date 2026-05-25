package events

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewCloudEventRequiresStrictMetadata(t *testing.T) {
	_, err := NewCloudEvent(t.Context(), TopicRetailOrderCreated, SourceRetailService, "orders/order-1", map[string]string{"order_id": "order-1"}, Metadata{
		EventID:    "evt-1",
		OccurredAt: time.Now(),
		OrderID:    "order-1",
	})
	if err == nil {
		t.Fatal("expected missing correlationid to fail")
	}
}

func TestCloudEventRoundTripExtractsExtensionsAndData(t *testing.T) {
	event, err := NewCloudEvent(t.Context(), TopicRetailOrderCreated, SourceRetailService, "orders/order-1", RetailOrderCreated{
		EventID:    "evt-1",
		OrderID:    "order-1",
		StoreID:    "store-1",
		OccurredAt: time.Now(),
	}, Metadata{
		EventID:       "evt-1",
		CorrelationID: "order-1",
		OrderID:       "order-1",
		StoreID:       "store-1",
		TraceID:       "trace-1",
		OccurredAt:    time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseCloudEvent(body)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.ID() != "evt-1" {
		t.Fatalf("expected event id evt-1, got %q", parsed.ID())
	}
	if got := ExtensionString(parsed, "orderid"); got != "order-1" {
		t.Fatalf("expected orderid extension, got %q", got)
	}
	data, err := DataAs[RetailOrderCreated](parsed)
	if err != nil {
		t.Fatal(err)
	}
	if data.OrderID != "order-1" || data.StoreID != "store-1" {
		t.Fatalf("unexpected data: %#v", data)
	}
}

func TestAllTraceableTopicsRoundTripAsCloudEvents(t *testing.T) {
	now := time.Now()
	for _, topic := range TraceableTopics {
		t.Run(topic, func(t *testing.T) {
			payload := map[string]interface{}{
				"event_id":    "evt-" + topic,
				"order_id":    "order-1",
				"store_id":    "store-1",
				"shipment_id": "shipment-1",
				"status":      "TEST",
				"occurred_at": now,
			}
			event, err := NewCloudEvent(t.Context(), topic, SourceLogisticsService, "contracts/"+topic, payload, Metadata{
				EventID:       "evt-" + topic,
				CorrelationID: "order-1",
				CausationID:   "evt-root",
				TraceID:       "11111111111111111111111111111111",
				OccurredAt:    now,
				OrderID:       "order-1",
				StoreID:       "store-1",
				ShipmentID:    "shipment-1",
				FarmID:        "farm-1",
				WarehouseID:   "warehouse-1",
				DriverID:      "driver-1",
				VehicleID:     "vehicle-1",
			})
			if err != nil {
				t.Fatal(err)
			}
			body, err := json.Marshal(event)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := ParseCloudEvent(body)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.Type() != topic {
				t.Fatalf("expected type %q, got %q", topic, parsed.Type())
			}
			if ExtensionString(parsed, "traceid") == "" {
				t.Fatal("expected traceid extension")
			}
			if _, err := DataAs[map[string]interface{}](parsed); err != nil {
				t.Fatal(err)
			}
		})
	}
}
