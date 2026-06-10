# Socket GPS Ingestion Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Update `socket-service` to handle incoming GPS updates from drivers via WebSocket and publish them to Kafka.

**Architecture:** `socket-service` acts as an ingress gateway. It receives GPS updates from authenticated drivers, verifies their role, and relays the update to Kafka for other services (like `logistics-service`) to consume.

**Tech Stack:** Go, Gorilla WebSocket, Kafka (Segmentio/kafka-go), CloudEvents.

---

### Task 1: WebSocket GPS Ingestion

**Files:**
- Modify: `src/apps/socket-service/internal/usecase/service.go`
- Test: `src/apps/socket-service/internal/usecase/service_test.go`

- [ ] **Step 1: Write the failing test**

Add a test case to `src/apps/socket-service/internal/usecase/service_test.go` that simulates a driver sending a `gps_update` message. It should mock the Kafka producer and verify that a `LogisticsGPSUpdated` event is published.

```go
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
	// CloudEvent wraps the data
	var ce events.CloudEvent
	err = json.Unmarshal(mockProducer.lastEvent, &ce)
	require.NoError(t, err)
	
	err = json.Unmarshal(ce.Data, &event)
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test -v src/apps/socket-service/internal/usecase/service_test.go`
Expected: FAIL (mockProducer.published is false because `readPump` is empty)

- [ ] **Step 3: Update readPump to handle gps_update**

Implement the logic in `service.go`.

```go
func (s *Service) readPump(client *Client) {
	defer func() {
		client.cancel()
		_ = client.conn.Close()
		s.removeClient(client)
	}()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		var env Envelope
		if err := json.Unmarshal(message, &env); err != nil {
			continue
		}

		if env.Type == "gps_update" {
			if client.scope.Role != "DRIVER" {
				continue
			}

			data, ok := env.Data.(map[string]interface{})
			if !ok {
				continue
			}

			lat, _ := data["lat"].(float64)
			long, _ := data["long"].(float64)
			shipmentID, _ := data["shipment_id"].(string)
			storeID, _ := data["store_id"].(string)

			event := events.LogisticsGPSUpdated{
				EventID:    uuid.NewString(),
				DriverID:   client.scope.Subject,
				ShipmentID: shipmentID,
				StoreID:    storeID,
				Latitude:   lat,
				Longitude:  long,
				OccurredAt: time.Now(),
			}

			ctx := context.Background()
			cloudEvent, err := events.NewCloudEvent(ctx, events.TopicLogisticsGPSUpdated, "/services/socket-service", "drivers/"+client.scope.Subject, event, events.Metadata{
				EventID:       event.EventID,
				CorrelationID: firstNonEmpty(shipmentID, client.scope.Subject),
				OccurredAt:    event.OccurredAt,
				DriverID:      client.scope.Subject,
				ShipmentID:    shipmentID,
				StoreID:       storeID,
			})
			if err == nil {
				_ = s.producer.Publish(ctx, events.TopicLogisticsGPSUpdated, client.scope.Subject, cloudEvent)
			}
		}
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test -v src/apps/socket-service/internal/usecase/service_test.go`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add src/apps/socket-service/internal/usecase/service.go src/apps/socket-service/internal/usecase/service_test.go
git commit -m "feat(socket): add WebSocket GPS ingestion to Kafka"
```
