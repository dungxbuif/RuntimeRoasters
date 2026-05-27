package usecase_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	paymentdomain "RuntimeRoasters/apps/payment-service/internal/domain"
	"RuntimeRoasters/apps/payment-service/internal/provider"
	paymentusecase "RuntimeRoasters/apps/payment-service/internal/usecase"
	"RuntimeRoasters/pkg/events"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type publishedEvent struct {
	topic   string
	key     string
	payload interface{}
}

type recordingProducer struct {
	events []publishedEvent
}

func (p *recordingProducer) Publish(ctx context.Context, topic string, key string, payload interface{}) error {
	_ = ctx
	p.events = append(p.events, publishedEvent{topic: topic, key: key, payload: payload})
	return nil
}

func (p *recordingProducer) Close() error {
	return nil
}

func TestPaymentSagaAndRefundContract(t *testing.T) {
	ctx := context.Background()
	paymentDB := openPaymentDB(t)

	paymentProducer := &recordingProducer{}
	paymentService := paymentusecase.NewService(
		paymentDB,
		paymentProducer,
		provider.NewFactory("stripe-secret", "vnpay-secret"),
		provider.ProviderStripe,
		"USD",
		true,
		events.TopicPaymentIntentCreated,
		events.TopicPaymentSimulatedCompleted,
		events.TopicPaymentFailed,
		events.TopicPaymentRefunded,
	)

	order := events.RetailOrderCreated{
		EventID:       "event-order-1",
		OrderID:       "11111111-1111-1111-1111-111111111111",
		StoreID:       "22222222-2222-2222-2222-222222222222",
		Items:         []events.RetailOrderItem{{SKU: "SL-ARABICA-ROASTED", Quantity: 999999}},
		TotalAmount:   100,
		PaymentMethod: provider.ProviderStripe,
		OccurredAt:    time.Now(),
	}
	require.NoError(t, paymentService.HandleOrderCreated(ctx, kafkaMessage(t, events.TopicRetailOrderCreated, order.OrderID, 0, order)))
	requireTopic(t, paymentProducer.events, events.TopicPaymentIntentCreated)
	requireNoTopic(t, paymentProducer.events, events.TopicPaymentSimulatedCompleted)

	var pending paymentdomain.Payment
	require.NoError(t, paymentDB.Where("order_id = ?", order.OrderID).Take(&pending).Error)
	require.Equal(t, paymentdomain.PaymentStatusPending, pending.Status)
	successBody := stripePayload(t, "evt_payment_success_1", "payment_intent.succeeded", pending.ProviderRef)
	require.NoError(t, paymentService.ProcessStripeWebhook(ctx, successBody, stripeSignature(t, "whsec_rr_demo_stripe_local", successBody)))
	requireTopic(t, paymentProducer.events, events.TopicPaymentCompleted)

	failed := events.WarehouseStockReservationFailed{
		EventID:    "event-stock-failed-1",
		OrderID:    order.OrderID,
		StoreID:    order.StoreID,
		Reason:     "insufficient stock",
		OccurredAt: time.Now(),
	}
	require.NoError(t, paymentService.HandleCompensationEvent(ctx, kafkaMessage(t, events.TopicWarehouseStockReservationFailed, order.OrderID, 1, failed)))
	requireTopic(t, paymentProducer.events, events.TopicPaymentRefunded)

	var refunded paymentdomain.Payment
	require.NoError(t, paymentDB.Where("order_id = ?", order.OrderID).Take(&refunded).Error)
	require.Equal(t, paymentdomain.PaymentStatusRefunded, refunded.Status)
	require.Equal(t, provider.ProviderStripe, refunded.Provider)
	require.NotEmpty(t, refunded.RefundRef)
}

func openPaymentDB(t *testing.T) *gorm.DB {
	t.Helper()
	name := strings.NewReplacer("/", "-", " ", "-").Replace(t.Name())
	db, err := gorm.Open(sqlite.Open("file:"+name+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(&paymentdomain.Payment{}, &paymentdomain.InboxEvent{}))
	require.NoError(t, db.AutoMigrate(&paymentdomain.WebhookEvent{}, &paymentdomain.PaymentWebhookKey{}))
	return db
}

func TestSimulateWebhookRequiresValidSignatureAndIsIdempotent(t *testing.T) {
	ctx := context.Background()
	db := openPaymentDB(t)
	producer := &recordingProducer{}
	service := paymentusecase.NewService(
		db,
		producer,
		provider.NewFactory("stripe-secret", "vnpay-secret"),
		provider.ProviderStripe,
		"USD",
		true,
		events.TopicPaymentIntentCreated,
		events.TopicPaymentCompleted,
		events.TopicPaymentFailed,
		events.TopicPaymentRefunded,
	)

	payment := paymentdomain.Payment{
		ID:          "33333333-3333-3333-3333-333333333333",
		OrderID:     "11111111-1111-1111-1111-111111111111",
		StoreID:     "22222222-2222-2222-2222-222222222222",
		Provider:    provider.ProviderStripe,
		ProviderRef: "pi_webhook_1",
		Items:       "[]",
		Amount:      100,
		Currency:    "USD",
		Status:      paymentdomain.PaymentStatusPending,
		Simulated:   true,
	}
	require.NoError(t, db.Create(&payment).Error)

	body := []byte(`{"event_id":"evt_1","provider_ref":"pi_webhook_1","status":"SUCCEEDED"}`)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	signature := hmacHex(t, "stripe-secret", []byte(timestamp+"."+string(body)))

	require.Error(t, service.SimulateWebhook(ctx, provider.ProviderStripe, body, timestamp, "bad"))
	require.NoError(t, service.SimulateWebhook(ctx, provider.ProviderStripe, body, timestamp, signature))
	require.NoError(t, service.SimulateWebhook(ctx, provider.ProviderStripe, body, timestamp, signature))

	var updated paymentdomain.Payment
	require.NoError(t, db.Where("id = ?", payment.ID).Take(&updated).Error)
	require.Equal(t, paymentdomain.PaymentStatusSucceeded, updated.Status)

	published := 0
	for _, event := range producer.events {
		if event.topic == events.TopicPaymentCompleted {
			published++
		}
	}
	require.Equal(t, 1, published)
}

func TestProcessStripeWebhookFailureAndSignatureValidation(t *testing.T) {
	ctx := context.Background()
	db := openPaymentDB(t)
	producer := &recordingProducer{}
	service := paymentusecase.NewService(
		db,
		producer,
		provider.NewFactory("stripe-secret", "vnpay-secret"),
		provider.ProviderStripe,
		"USD",
		true,
		events.TopicPaymentIntentCreated,
		events.TopicPaymentCompleted,
		events.TopicPaymentFailed,
		events.TopicPaymentRefunded,
	)

	payment := paymentdomain.Payment{
		ID:          "33333333-3333-3333-3333-333333333333",
		OrderID:     "11111111-1111-1111-1111-111111111111",
		StoreID:     "22222222-2222-2222-2222-222222222222",
		Provider:    provider.ProviderStripe,
		ProviderRef: "pi_webhook_failed_1",
		Items:       "[]",
		Amount:      100,
		Currency:    "USD",
		Status:      paymentdomain.PaymentStatusPending,
		Simulated:   true,
	}
	require.NoError(t, db.Create(&payment).Error)

	body := stripePayload(t, "evt_failed_1", "payment_intent.payment_failed", payment.ProviderRef)
	require.Error(t, service.ProcessStripeWebhook(ctx, body, "t=1,v1=bad"))
	require.NoError(t, service.ProcessStripeWebhook(ctx, body, stripeSignature(t, "whsec_rr_demo_stripe_local", body)))
	require.NoError(t, service.ProcessStripeWebhook(ctx, body, stripeSignature(t, "whsec_rr_demo_stripe_local", body)))

	var updated paymentdomain.Payment
	require.NoError(t, db.Where("id = ?", payment.ID).Take(&updated).Error)
	require.Equal(t, paymentdomain.PaymentStatusFailed, updated.Status)

	published := 0
	for _, event := range producer.events {
		if event.topic == events.TopicPaymentFailed {
			published++
		}
	}
	require.Equal(t, 1, published)
}

func hmacHex(t *testing.T, secret string, payload []byte) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(payload)
	require.NoError(t, err)
	return hex.EncodeToString(mac.Sum(nil))
}

func stripePayload(t *testing.T, eventID string, eventType string, providerRef string) []byte {
	t.Helper()
	payload := map[string]interface{}{
		"id":       eventID,
		"object":   "event",
		"type":     eventType,
		"livemode": false,
		"created":  time.Now().Unix(),
		"data": map[string]interface{}{
			"object": map[string]interface{}{
				"id":       providerRef,
				"object":   "payment_intent",
				"amount":   10000,
				"currency": "usd",
				"status":   "succeeded",
				"metadata": map[string]string{
					"order_id":   "11111111-1111-1111-1111-111111111111",
					"payment_id": "33333333-3333-3333-3333-333333333333",
					"store_id":   "22222222-2222-2222-2222-222222222222",
				},
			},
		},
	}
	if eventType == "payment_intent.payment_failed" {
		object := payload["data"].(map[string]interface{})["object"].(map[string]interface{})
		object["status"] = "requires_payment_method"
		object["last_payment_error"] = map[string]string{
			"type":    "card_error",
			"code":    "card_declined",
			"message": "Demo card declined",
		}
	}
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	return body
}

func stripeSignature(t *testing.T, secret string, body []byte) string {
	t.Helper()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	return fmt.Sprintf("t=%s,v1=%s", timestamp, hmacHex(t, secret, []byte(timestamp+"."+string(body))))
}

func kafkaMessage(t *testing.T, topic string, key string, offset int64, payload interface{}) kafkago.Message {
	t.Helper()
	cloudEvent, err := events.NewCloudEvent(context.Background(), topic, "/tests/payment-saga", fmt.Sprintf("orders/%s", key), payload, testMetadata(payload, key))
	require.NoError(t, err)
	body, err := json.Marshal(cloudEvent)
	require.NoError(t, err)
	return kafkago.Message{Topic: topic, Partition: 0, Offset: offset, Key: []byte(key), Value: body}
}

func testMetadata(payload interface{}, key string) events.Metadata {
	now := time.Now()
	switch event := payload.(type) {
	case events.RetailOrderCreated:
		return events.Metadata{EventID: event.EventID, CorrelationID: event.OrderID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID}
	case events.WarehouseStockReservationFailed:
		return events.Metadata{EventID: event.EventID, CorrelationID: event.OrderID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID}
	default:
		return events.Metadata{EventID: "evt-" + key, CorrelationID: key, OccurredAt: now}
	}
}

func requireTopic(t *testing.T, published []publishedEvent, topic string) publishedEvent {
	t.Helper()
	for _, event := range published {
		if event.topic == topic {
			return event
		}
	}
	t.Fatalf("expected topic %s in published events %#v", topic, published)
	return publishedEvent{}
}

func requireNoTopic(t *testing.T, published []publishedEvent, topic string) {
	t.Helper()
	for _, event := range published {
		require.NotEqual(t, topic, event.topic)
	}
}
