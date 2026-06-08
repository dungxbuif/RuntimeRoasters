package usecase

import (
	"context"
	"encoding/json"
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
	for _, seed := range demoAuditEvents() {
		cloudEvent, err := events.NewCloudEvent(ctx, seed.topic, seed.source, seed.subject, seed.payload, seed.metadata)
		if err != nil {
			return created, err
		}
		payload, err := json.Marshal(cloudEvent)
		if err != nil {
			return created, err
		}
		if err := s.HandleEvent(ctx, kafkago.Message{Topic: seed.topic, Key: []byte(seed.key), Value: payload}); err != nil {
			return created, err
		}
		created++
	}
	return created, nil
}

func (s *Service) DemoAuditCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&struct{}{}).Table("audit_logs").Where("message_id LIKE ?", "%"+demoSeedPrefix+"%").Count(&count).Error
	return count, err
}

type demoAuditEvent struct {
	topic    string
	source   string
	subject  string
	key      string
	payload  interface{}
	metadata events.Metadata
}

func demoAuditEvents() []demoAuditEvent {
	orderID := stableUUID("order", "succeeded")
	paymentID := stableUUID("payment", "succeeded")
	shipmentID := stableUUID("shipment", "completed")
	storeID := "11111111-1111-1111-1111-111111111103"
	traceID := demoSeedPrefix + "-trace-completed"
	failedOrderID := stableUUID("order", "failed")
	failedPaymentID := stableUUID("payment", "failed")
	refundedOrderID := stableUUID("order", "refunded")
	refundedPaymentID := stableUUID("payment", "refunded")
	return []demoAuditEvent{
		{
			topic: events.TopicRetailOrderCreated, source: events.SourceRetailService, subject: "orders/" + orderID, key: orderID,
			payload:  events.RetailOrderCreated{EventID: demoSeedPrefix + "-audit-order-" + orderID, OrderID: orderID, StoreID: storeID, Items: []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 10}}, TotalAmount: 1200000, PaymentMethod: "STRIPE", OccurredAt: demoSeedReferenceTime.Add(-48 * time.Hour)},
			metadata: events.Metadata{EventID: demoSeedPrefix + "-audit-order-" + orderID, CorrelationID: orderID, OrderID: orderID, StoreID: storeID, TraceID: traceID, OccurredAt: demoSeedReferenceTime.Add(-48 * time.Hour)},
		},
		{
			topic: events.TopicPaymentSimulatedCompleted, source: events.SourcePaymentService, subject: "orders/" + orderID, key: orderID,
			payload:  events.PaymentCompleted{EventID: demoSeedPrefix + "-audit-payment-completed-" + orderID, PaymentID: paymentID, OrderID: orderID, StoreID: storeID, Provider: "STRIPE", ProviderRef: demoSeedPrefix + "-pi-succeeded", Items: []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 10}}, Amount: 1200000, Currency: "VND", OccurredAt: demoSeedReferenceTime.Add(-48*time.Hour + 3*time.Minute)},
			metadata: events.Metadata{EventID: demoSeedPrefix + "-audit-payment-completed-" + orderID, CorrelationID: orderID, OrderID: orderID, PaymentID: paymentID, StoreID: storeID, TraceID: traceID, OccurredAt: demoSeedReferenceTime.Add(-48*time.Hour + 3*time.Minute)},
		},
		{
			topic: events.TopicWarehouseStockReserved, source: events.SourceWarehouseService, subject: "orders/" + orderID, key: orderID,
			payload:  events.WarehouseStockReserved{EventID: demoSeedPrefix + "-audit-stock-" + orderID, OrderID: orderID, StoreID: storeID, WarehouseID: "WAREHOUSE-HCM-001", SKU: "SKU-AR-VN-LD-001", Quantity: 10, OccurredAt: demoSeedReferenceTime.Add(-48*time.Hour + 5*time.Minute)},
			metadata: events.Metadata{EventID: demoSeedPrefix + "-audit-stock-" + orderID, CorrelationID: orderID, OrderID: orderID, StoreID: storeID, WarehouseID: "WAREHOUSE-HCM-001", TraceID: traceID, OccurredAt: demoSeedReferenceTime.Add(-48*time.Hour + 5*time.Minute)},
		},
		{
			topic: events.TopicLogisticsDeliveryCompleted, source: events.SourceLogisticsService, subject: "shipments/" + shipmentID, key: orderID,
			payload:  events.LogisticsShipmentDelivered{EventID: demoSeedPrefix + "-audit-delivery-" + orderID, ShipmentID: shipmentID, OrderID: orderID, StoreID: storeID, OccurredAt: demoSeedReferenceTime.Add(-46 * time.Hour)},
			metadata: events.Metadata{EventID: demoSeedPrefix + "-audit-delivery-" + orderID, CorrelationID: orderID, OrderID: orderID, ShipmentID: shipmentID, StoreID: storeID, TraceID: traceID, OccurredAt: demoSeedReferenceTime.Add(-46 * time.Hour)},
		},
		{
			topic: events.TopicPaymentFailed, source: events.SourcePaymentService, subject: "orders/" + failedOrderID, key: failedOrderID,
			payload:  events.PaymentFailed{EventID: demoSeedPrefix + "-audit-payment-failed-" + failedOrderID, PaymentID: failedPaymentID, OrderID: failedOrderID, StoreID: "11111111-1111-1111-1111-111111111102", Provider: "STRIPE", ProviderRef: demoSeedPrefix + "-pi-failed", Reason: "card_declined", OccurredAt: demoSeedReferenceTime.Add(-36*time.Hour + 3*time.Minute)},
			metadata: events.Metadata{EventID: demoSeedPrefix + "-audit-payment-failed-" + failedOrderID, CorrelationID: failedOrderID, OrderID: failedOrderID, PaymentID: failedPaymentID, StoreID: "11111111-1111-1111-1111-111111111102", TraceID: demoSeedPrefix + "-trace-failed", OccurredAt: demoSeedReferenceTime.Add(-36*time.Hour + 3*time.Minute)},
		},
		{
			topic: events.TopicPaymentRefunded, source: events.SourcePaymentService, subject: "orders/" + refundedOrderID, key: refundedOrderID,
			payload:  events.PaymentRefunded{EventID: demoSeedPrefix + "-audit-payment-refunded-" + refundedOrderID, PaymentID: refundedPaymentID, OrderID: refundedOrderID, StoreID: "11111111-1111-1111-1111-111111111101", Provider: "STRIPE", ProviderRef: demoSeedPrefix + "-pi-refunded", RefundRef: demoSeedPrefix + "-re-quality", Reason: "quality_issue", Amount: 980000, Currency: "VND", OccurredAt: demoSeedReferenceTime.Add(-24*time.Hour + 10*time.Minute)},
			metadata: events.Metadata{EventID: demoSeedPrefix + "-audit-payment-refunded-" + refundedOrderID, CorrelationID: refundedOrderID, OrderID: refundedOrderID, PaymentID: refundedPaymentID, StoreID: "11111111-1111-1111-1111-111111111101", TraceID: demoSeedPrefix + "-trace-refunded", OccurredAt: demoSeedReferenceTime.Add(-24*time.Hour + 10*time.Minute)},
		},
	}
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
