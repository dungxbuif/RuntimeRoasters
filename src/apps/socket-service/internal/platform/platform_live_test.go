//go:build platform

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"RuntimeRoasters/apps/socket-service/internal/usecase"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarehouseDispatchEventCreatesLiveNotification(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	brokers := envOrDefault("KAFKA_BROKERS", "localhost:9094")
	redisAddr := envOrDefault("VALKEY_ADDR", "localhost:6379")
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer rdb.Close()

	require.NoError(t, rdb.Ping(ctx).Err())

	eventID := uuid.NewString()
	orderID := "platform-order-" + eventID[:8]
	warehouseID := "platform-warehouse-" + eventID[:8]
	storeID := "platform-store-" + eventID[:8]
	event := events.WarehouseDispatchRequested{
		EventID:     eventID,
		DispatchID:  "platform-dispatch-" + eventID[:8],
		OrderID:     orderID,
		StoreID:     storeID,
		WarehouseID: warehouseID,
		OccurredAt:  time.Now().UTC(),
	}
	cloudEvent, err := events.NewCloudEvent(ctx, events.TopicWarehouseDispatchRequested, "/platform/realtime-awareness", "events/"+eventID, event, events.Metadata{
		EventID:       event.EventID,
		CorrelationID: event.OrderID,
		OrderID:       event.OrderID,
		StoreID:       event.StoreID,
		WarehouseID:   event.WarehouseID,
		OccurredAt:    event.OccurredAt,
	})
	require.NoError(t, err)
	payload, err := json.Marshal(cloudEvent)
	require.NoError(t, err)

	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Topic:        events.TopicWarehouseDispatchRequested,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
	}
	defer writer.Close()
	require.NoError(t, writer.WriteMessages(ctx, kafka.Message{Key: []byte(orderID), Value: payload}))

	key := "notification:n-" + eventID
	var notification usecase.Notification
	require.EventuallyWithT(t, func(collect *assertCollect) {
		body, err := rdb.Get(ctx, key).Bytes()
		require.NoError(collect, err)
		require.NoError(collect, json.Unmarshal(body, &notification))
		require.Equal(collect, identity.RoleWarehouseMgr, notification.TargetRole)
		require.Equal(collect, warehouseID, notification.WarehouseID)
		require.Equal(collect, storeID, notification.StoreID)
		require.Equal(collect, orderID, notification.EntityID)
		require.Equal(collect, events.TopicWarehouseDispatchRequested, notification.Type)
		require.Equal(collect, usecase.NotificationStatusUnread, notification.Status)
	}, 15*time.Second, 500*time.Millisecond)

	indexEntries, err := rdb.LRange(ctx, "notifications:warehouse:"+warehouseID, 0, -1).Result()
	require.NoError(t, err)
	require.Contains(t, indexEntries, "n-"+eventID)
}

type assertCollect = assert.CollectT

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func TestMain(m *testing.M) {
	if os.Getenv("RUN_PLATFORM_TESTS") != "1" {
		fmt.Fprintln(os.Stderr, "set RUN_PLATFORM_TESTS=1 to run live platform tests")
		os.Exit(2)
	}
	os.Exit(m.Run())
}
