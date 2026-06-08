package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
)

const demoSeedPrefix = "rr31-demo"

var demoSeedReferenceTime = time.Date(2026, 6, 8, 8, 0, 0, 0, time.UTC)

func (s *Service) SeedDemoData(ctx context.Context, force bool, usersMap map[string]string) (int64, error) {
	_ = force
	_ = usersMap
	created := int64(0)
	for _, event := range demoTraceEvents() {
		cloudEvent, err := events.NewCloudEvent(ctx, event.topic, event.source, event.subject, event.payload, event.metadata)
		if err != nil {
			return created, err
		}
		if err := s.HandleEvent(ctx, kafkago.Message{
			Topic: event.topic,
			Key:   []byte(event.key),
			Value: mustMarshalCloudEvent(cloudEvent),
		}); err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}

func (s *Service) DemoTraceDocumentCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&struct{}{}).Table("trace_documents").Where("CAST(document AS TEXT) LIKE ?", "%"+demoSeedPrefix+"%").Count(&count).Error
	return count, err
}

type demoTraceEvent struct {
	topic    string
	source   string
	subject  string
	key      string
	payload  interface{}
	metadata events.Metadata
}

func demoTraceEvents() []demoTraceEvent {
	return append(completedTraceEvents(), financeTraceEvents()...)
}

func completedTraceEvents() []demoTraceEvent {
	orderID := stableUUID("order", "succeeded")
	shipmentID := stableUUID("shipment", "completed")
	traceID := demoSeedPrefix + "-trace-completed"
	storeID := "11111111-1111-1111-1111-111111111103"
	return []demoTraceEvent{
		retailOrderEvent(orderID, storeID, traceID, 1200000, demoSeedReferenceTime.Add(-48*time.Hour)),
		paymentIntentEvent(orderID, stableUUID("payment", "succeeded"), storeID, traceID, "rr31-demo-pi-succeeded", 1200000, demoSeedReferenceTime.Add(-48*time.Hour+time.Minute)),
		paymentCompletedEvent(orderID, stableUUID("payment", "succeeded"), storeID, traceID, "rr31-demo-pi-succeeded", 1200000, demoSeedReferenceTime.Add(-48*time.Hour+3*time.Minute)),
		warehouseEvent(events.TopicWarehouseStockReserved, orderID, storeID, traceID, demoSeedReferenceTime.Add(-48*time.Hour+5*time.Minute)),
		logisticsAssignedEvent(orderID, shipmentID, storeID, traceID, demoSeedReferenceTime.Add(-47*time.Hour)),
		logisticsCompletedEvent(orderID, shipmentID, storeID, traceID, demoSeedReferenceTime.Add(-46*time.Hour)),
	}
}

func financeTraceEvents() []demoTraceEvent {
	failedOrderID := stableUUID("order", "failed")
	refundedOrderID := stableUUID("order", "refunded")
	return []demoTraceEvent{
		retailOrderEvent(failedOrderID, "11111111-1111-1111-1111-111111111102", demoSeedPrefix+"-trace-failed", 430000, demoSeedReferenceTime.Add(-36*time.Hour)),
		paymentIntentEvent(failedOrderID, stableUUID("payment", "failed"), "11111111-1111-1111-1111-111111111102", demoSeedPrefix+"-trace-failed", "rr31-demo-pi-failed", 430000, demoSeedReferenceTime.Add(-36*time.Hour+time.Minute)),
		paymentFailedEvent(failedOrderID, stableUUID("payment", "failed"), "11111111-1111-1111-1111-111111111102", demoSeedPrefix+"-trace-failed", "rr31-demo-pi-failed", 430000, demoSeedReferenceTime.Add(-36*time.Hour+3*time.Minute)),
		retailOrderEvent(refundedOrderID, "11111111-1111-1111-1111-111111111101", demoSeedPrefix+"-trace-refunded", 980000, demoSeedReferenceTime.Add(-24*time.Hour)),
		paymentIntentEvent(refundedOrderID, stableUUID("payment", "refunded"), "11111111-1111-1111-1111-111111111101", demoSeedPrefix+"-trace-refunded", "rr31-demo-pi-refunded", 980000, demoSeedReferenceTime.Add(-24*time.Hour+time.Minute)),
		paymentRefundedEvent(refundedOrderID, stableUUID("payment", "refunded"), "11111111-1111-1111-1111-111111111101", demoSeedPrefix+"-trace-refunded", "rr31-demo-pi-refunded", 980000, demoSeedReferenceTime.Add(-24*time.Hour+10*time.Minute)),
	}
}

func retailOrderEvent(orderID string, storeID string, traceID string, amount float64, t time.Time) demoTraceEvent {
	payload := events.RetailOrderCreated{EventID: demoSeedPrefix + "-order-" + orderID, OrderID: orderID, StoreID: storeID, Items: []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 1}}, TotalAmount: amount, PaymentMethod: "STRIPE", OccurredAt: t}
	return demoTraceEvent{topic: events.TopicRetailOrderCreated, source: events.SourceRetailService, subject: "orders/" + orderID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, StoreID: storeID, TraceID: traceID}}
}

func paymentIntentEvent(orderID, paymentID, storeID, traceID, providerRef string, amount float64, t time.Time) demoTraceEvent {
	payload := events.PaymentIntentCreated{EventID: demoSeedPrefix + "-intent-" + orderID, PaymentID: paymentID, OrderID: orderID, StoreID: storeID, Provider: "STRIPE", ProviderRef: providerRef, Items: []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 1}}, Amount: amount, Currency: "VND", CheckoutURL: "https://payments.runtime-roasters.local/demo/" + providerRef, Simulated: true, OccurredAt: t}
	return demoTraceEvent{topic: events.TopicPaymentIntentCreated, source: events.SourcePaymentService, subject: "orders/" + orderID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, PaymentID: paymentID, StoreID: storeID, TraceID: traceID}}
}

func paymentCompletedEvent(orderID, paymentID, storeID, traceID, providerRef string, amount float64, t time.Time) demoTraceEvent {
	payload := events.PaymentCompleted{EventID: demoSeedPrefix + "-completed-" + orderID, PaymentID: paymentID, OrderID: orderID, StoreID: storeID, Provider: "STRIPE", ProviderRef: providerRef, Items: []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 1}}, Amount: amount, Currency: "VND", OccurredAt: t}
	return demoTraceEvent{topic: events.TopicPaymentSimulatedCompleted, source: events.SourcePaymentService, subject: "orders/" + orderID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, PaymentID: paymentID, StoreID: storeID, TraceID: traceID}}
}

func paymentFailedEvent(orderID, paymentID, storeID, traceID, providerRef string, amount float64, t time.Time) demoTraceEvent {
	payload := events.PaymentFailed{EventID: demoSeedPrefix + "-failed-" + orderID, PaymentID: paymentID, OrderID: orderID, StoreID: storeID, Provider: "STRIPE", ProviderRef: providerRef, Reason: "card_declined", OccurredAt: t}
	return demoTraceEvent{topic: events.TopicPaymentFailed, source: events.SourcePaymentService, subject: "orders/" + orderID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, PaymentID: paymentID, StoreID: storeID, TraceID: traceID}}
}

func paymentRefundedEvent(orderID, paymentID, storeID, traceID, providerRef string, amount float64, t time.Time) demoTraceEvent {
	payload := events.PaymentRefunded{EventID: demoSeedPrefix + "-refunded-" + orderID, PaymentID: paymentID, OrderID: orderID, StoreID: storeID, Provider: "STRIPE", ProviderRef: providerRef, RefundRef: demoSeedPrefix + "-re-quality", Reason: "quality_issue", Amount: amount, Currency: "VND", OccurredAt: t}
	return demoTraceEvent{topic: events.TopicPaymentRefunded, source: events.SourcePaymentService, subject: "orders/" + orderID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, PaymentID: paymentID, StoreID: storeID, TraceID: traceID}}
}

func warehouseEvent(topic, orderID, storeID, traceID string, t time.Time) demoTraceEvent {
	payload := events.WarehouseStockReserved{EventID: demoSeedPrefix + "-stock-" + orderID, OrderID: orderID, StoreID: storeID, WarehouseID: "WAREHOUSE-HCM-001", SKU: "SKU-AR-VN-LD-001", Quantity: 10, OccurredAt: t}
	return demoTraceEvent{topic: topic, source: events.SourceWarehouseService, subject: "orders/" + orderID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, StoreID: storeID, WarehouseID: payload.WarehouseID, TraceID: traceID}}
}

func logisticsAssignedEvent(orderID, shipmentID, storeID, traceID string, t time.Time) demoTraceEvent {
	payload := events.LogisticsShipmentAssigned{EventID: demoSeedPrefix + "-delivery-assigned-" + orderID, ShipmentID: shipmentID, OrderID: orderID, StoreID: storeID, DriverID: "22222222-2222-2222-2222-222222222203", VehicleID: "VEHICLE-DEMO-003", OccurredAt: t}
	return demoTraceEvent{topic: events.TopicLogisticsDeliveryAssigned, source: events.SourceLogisticsService, subject: "shipments/" + shipmentID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, ShipmentID: shipmentID, StoreID: storeID, DriverID: payload.DriverID, VehicleID: payload.VehicleID, TraceID: traceID}}
}

func logisticsCompletedEvent(orderID, shipmentID, storeID, traceID string, t time.Time) demoTraceEvent {
	payload := events.LogisticsShipmentDelivered{EventID: demoSeedPrefix + "-delivery-completed-" + orderID, ShipmentID: shipmentID, OrderID: orderID, StoreID: storeID, OccurredAt: t}
	return demoTraceEvent{topic: events.TopicLogisticsDeliveryCompleted, source: events.SourceLogisticsService, subject: "shipments/" + shipmentID, key: orderID, payload: payload, metadata: events.Metadata{EventID: payload.EventID, CorrelationID: orderID, OccurredAt: t, OrderID: orderID, ShipmentID: shipmentID, StoreID: storeID, TraceID: traceID}}
}

func stableUUID(parts ...string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(demoSeedPrefix+"|"+joinSeedParts(parts))).String()
}

func joinSeedParts(parts []string) string {
	joined := ""
	for i, part := range parts {
		if i > 0 {
			joined += "|"
		}
		joined += part
	}
	return joined
}

func mustMarshalCloudEvent(event interface{}) []byte {
	payload, err := json.Marshal(event)
	if err != nil {
		panic(fmt.Sprintf("marshal demo cloudevent: %v", err))
	}
	return payload
}
