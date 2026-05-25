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

	paymentdomain "github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/provider"
	paymentusecase "github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
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
		PaymentMethod: provider.ProviderVNPay,
		OccurredAt:    time.Now(),
	}
	require.NoError(t, paymentService.HandleOrderCreated(ctx, kafkaMessage(t, events.TopicRetailOrderCreated, order.OrderID, 0, order)))
	requireTopic(t, paymentProducer.events, events.TopicPaymentIntentCreated)
	requireTopic(t, paymentProducer.events, events.TopicPaymentSimulatedCompleted)

	failed := events.WarehouseStockReservationFailed{
		EventID:    "event-stock-failed-1",
		OrderID:    order.OrderID,
		StoreID:    order.StoreID,
		Reason:     "insufficient stock",
		OccurredAt: time.Now(),
	}
	require.NoError(t, paymentService.HandleCompensationEvent(ctx, kafkaMessage(t, events.TopicWarehouseStockReservationFailed, order.OrderID, 1, failed)))
	requireTopic(t, paymentProducer.events, events.TopicPaymentRefunded)

	var payment paymentdomain.Payment
	require.NoError(t, paymentDB.Where("order_id = ?", order.OrderID).Take(&payment).Error)
	require.Equal(t, paymentdomain.PaymentStatusRefunded, payment.Status)
	require.Equal(t, provider.ProviderVNPay, payment.Provider)
	require.NotEmpty(t, payment.RefundRef)
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
	require.NoError(t, db.AutoMigrate(&paymentdomain.WebhookEvent{}))
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

func hmacHex(t *testing.T, secret string, payload []byte) string {
	t.Helper()
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(payload)
	require.NoError(t, err)
	return hex.EncodeToString(mac.Sum(nil))
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
