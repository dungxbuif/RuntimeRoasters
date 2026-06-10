# Socket Service GPS Ingress Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Update `readPump` in `socket-service` to handle `gps_update` messages and publish them to Kafka.

**Architecture:**
- WebSocket `readPump` loop unmarshals messages into `Envelope`.
- If type is `gps_update`, it verifies the user is a `DRIVER`.
- Extracts GPS data and publishes to Kafka as a CloudEvent.

**Tech Stack:**
- Go
- Gorilla WebSocket
- Kafka
- CloudEvents

---

### Task 1: WebSocket GPS Ingestion

**Files:**
- Modify: `src/apps/socket-service/internal/usecase/service.go`
- Test: `src/apps/socket-service/internal/usecase/service_test.go`

- [ ] **Step 1: Run existing tests to verify failure**

Run: `go test -v src/apps/socket-service/internal/usecase/service_test.go src/apps/socket-service/internal/usecase/service.go`
Expected: `TestGPSIngestion` should FAIL.

- [ ] **Step 2: Implement GPS update logic in readPump**

Modify `src/apps/socket-service/internal/usecase/service.go`:

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

			cloudEvent, err := events.NewCloudEvent(context.Background(), events.TopicLogisticsGPSUpdated, "/services/socket-service", "drivers/"+client.scope.Subject, event, events.Metadata{
				EventID:    event.EventID,
				OccurredAt: event.OccurredAt,
				DriverID:   event.DriverID,
				ShipmentID: event.ShipmentID,
			})
			if err != nil {
				continue
			}

			_ = s.producer.Publish(context.Background(), events.TopicLogisticsGPSUpdated, event.DriverID, cloudEvent)
		}
	}
}
```

- [ ] **Step 3: Run tests to verify success**

Run: `go test -v src/apps/socket-service/internal/usecase/service_test.go src/apps/socket-service/internal/usecase/service.go`
Expected: `TestGPSIngestion` should PASS.

- [ ] **Step 4: Commit changes**

```bash
git add src/apps/socket-service/internal/usecase/service.go
git commit -m "feat(socket): add WebSocket GPS ingestion to Kafka"
```
