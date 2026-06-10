package usecase

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"github.com/alicebob/miniredis/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTicketFlow(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})

	service := NewService(rdb, nil, "test-topic", 1*time.Minute, "")
	ctx := context.Background()

	claims := identity.Claims{
		Subject: "user-123",
		Role:    "ADMIN",
	}

	t.Run("Issue and Exchange Ticket", func(t *testing.T) {
		ticket, err := service.IssueTicket(ctx, claims)
		assert.NoError(t, err)
		assert.NotEmpty(t, ticket)

		exchangedClaims, err := service.ExchangeTicket(ctx, ticket)
		assert.NoError(t, err)
		assert.Equal(t, claims.Subject, exchangedClaims.Subject)
		assert.Equal(t, claims.Role, exchangedClaims.Role)

		// One-time use: second exchange should fail
		_, err = service.ExchangeTicket(ctx, ticket)
		assert.Error(t, err)
		assert.Equal(t, "invalid or expired ticket", err.Error())
	})

	t.Run("Invalid Ticket", func(t *testing.T) {
		_, err := service.ExchangeTicket(ctx, "non-existent")
		assert.Error(t, err)
	})
}

func TestNotificationScopeAndAck(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	service := NewService(rdb, nil, "test-topic", 1*time.Minute, "")
	ctx := context.Background()
	notification := Notification{
		ID:          "notification-1",
		TargetRole:  identity.RoleWarehouseMgr,
		WarehouseID: "warehouse-1",
		EntityType:  "order",
		EntityID:    "order-1",
		Type:        events.TopicWarehouseDispatchRequested,
		Title:       "Dispatch requested",
		Message:     "Order ready for dispatch",
		Severity:    "info",
		Status:      NotificationStatusUnread,
		CreatedAt:   time.Now(),
	}
	require.NoError(t, service.CreateNotification(ctx, notification))

	warehouseClaims := identity.Claims{Subject: "mgr-1", Role: identity.RoleWarehouseMgr, WarehouseIDs: []string{"warehouse-1"}}
	visible, err := service.ListNotifications(ctx, warehouseClaims)
	require.NoError(t, err)
	require.Len(t, visible, 1)
	assert.Equal(t, notification.ID, visible[0].ID)

	storeClaims := identity.Claims{Subject: "store-1", Role: identity.RoleStoreMgr, StoreIDs: []string{"store-1"}}
	hidden, err := service.ListNotifications(ctx, storeClaims)
	require.NoError(t, err)
	assert.Empty(t, hidden)

	acked, err := service.AcknowledgeNotification(ctx, warehouseClaims, notification.ID)
	require.NoError(t, err)
	assert.Equal(t, NotificationStatusAcked, acked.Status)
	assert.NotNil(t, acked.AcknowledgedAt)
}

func TestKafkaEventCreatesScopedNotification(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	service := NewService(rdb, nil, "test-topic", 1*time.Minute, "")
	ctx := context.Background()
	event := events.WarehouseDispatchRequested{
		EventID:     uuid.NewString(),
		DispatchID:  "dispatch-1",
		OrderID:     "order-1",
		StoreID:     "store-1",
		WarehouseID: "warehouse-1",
		OccurredAt:  time.Now(),
	}
	payload := socketCloudEvent(t, events.TopicWarehouseDispatchRequested, event, events.Metadata{
		EventID:       event.EventID,
		CorrelationID: event.OrderID,
		OrderID:       event.OrderID,
		StoreID:       event.StoreID,
		WarehouseID:   event.WarehouseID,
		OccurredAt:    event.OccurredAt,
	})

	require.NoError(t, service.HandleKafkaMessage(ctx, kafkago.Message{Topic: events.TopicWarehouseDispatchRequested, Key: []byte(event.OrderID), Value: payload}))

	claims := identity.Claims{Subject: "mgr-1", Role: identity.RoleWarehouseMgr, WarehouseIDs: []string{event.WarehouseID}}
	notifications, err := service.ListNotifications(ctx, claims)
	require.NoError(t, err)
	require.Len(t, notifications, 1)
	assert.Equal(t, events.TopicWarehouseDispatchRequested, notifications[0].Type)
	assert.Equal(t, event.OrderID, notifications[0].EntityID)
	assert.Equal(t, event.WarehouseID, notifications[0].WarehouseID)
}

func TestGPSIngestion(t *testing.T) {
	s := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	mockProducer := &mockProducer{}
	service := NewService(rdb, mockProducer, "broadcast", 1*time.Minute, "")

	// Create a dummy WebSocket connection using a local server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = service.ServeWS(w, r, ClientScope{
			Role:    "DRIVER",
			Subject: "driver-123",
		})
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	gpsMsg := Envelope{
		Type: "gps_update",
		Data: map[string]interface{}{
			"lat":         10.5,
			"long":        20.5,
			"shipment_id": "ship-1",
		},
	}
	body, _ := json.Marshal(gpsMsg)
	err = conn.WriteMessage(websocket.TextMessage, body)
	require.NoError(t, err)

	// Wait for processing
	time.Sleep(100 * time.Millisecond)

	assert.True(t, mockProducer.published)
	assert.Equal(t, events.TopicLogisticsGPSUpdated, mockProducer.topic)

	var event events.LogisticsGPSUpdated
	// cloudevents.Event wraps the data
	var ce cloudevents.Event
	err = json.Unmarshal(mockProducer.lastEvent, &ce)
	require.NoError(t, err)

	err = ce.DataAs(&event)
	require.NoError(t, err)

	assert.Equal(t, "driver-123", event.DriverID)
	assert.Equal(t, 10.5, event.Latitude)
}

type mockProducer struct {
	published bool
	topic     string
	lastEvent []byte
}

func (m *mockProducer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	m.published = true
	m.topic = topic
	m.lastEvent, _ = json.Marshal(value)
	return nil
}

func (m *mockProducer) Close() error { return nil }

func socketCloudEvent(t *testing.T, topic string, data interface{}, metadata events.Metadata) []byte {
	t.Helper()
	cloudEvent, err := events.NewCloudEvent(context.Background(), topic, "/tests/socket-service", "events/"+metadata.EventID, data, metadata)
	require.NoError(t, err)
	payload, err := json.Marshal(cloudEvent)
	require.NoError(t, err)
	return payload
}
