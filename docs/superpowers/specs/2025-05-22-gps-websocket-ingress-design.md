# Design Spec: Socket Service Ingress (GPS WebSocket)

**Status:** Draft
**Date:** 2025-05-22
**Topic:** GPS WebSocket Ingress for Socket Service

## 1. Overview
The goal is to enable the `socket-service` to ingest GPS updates from drivers via WebSocket and publish them to Kafka as CloudEvents. This is part of the Valkey Geo Tracking implementation.

## 2. Architecture
The `socket-service` maintains WebSocket connections with clients (drivers, admins, managers). The `readPump` goroutine handles incoming messages from each client.

### Data Flow
1. **Driver** sends a JSON message over WebSocket.
2. `readPump` reads the message and unmarshals it into an `Envelope`.
3. If the message type is `gps_update` and the client's role is `DRIVER`:
    - The message is validated.
    - A `LogisticsGPSUpdated` event is created.
    - The event is wrapped in a `CloudEvent`.
    - The `CloudEvent` is published to Kafka on the `logistics.gps.updated` topic.

## 3. Components

### Socket Service (`readPump`)
- **Validation**: Only clients with the `DRIVER` role can send GPS updates.
- **Payload Extraction**: extracts `lat`, `long`, `shipment_id`, and `store_id` from the `Envelope.Data`.
- **Event Construction**: Uses `pkg/events` to create a `LogisticsGPSUpdated` struct and then a `CloudEvent`.
- **Publishing**: Uses the service's Kafka producer to publish the event.

## 4. Error Handling
- Invalid JSON in WebSocket messages will be ignored (logged if necessary, but loop continues).
- Messages from non-driver roles with `gps_update` type will be ignored.
- Errors in publishing to Kafka will be logged but won't crash the connection.

## 5. Testing
- `TestGPSIngestion` in `src/apps/socket-service/internal/usecase/service_test.go` verifies the end-to-end flow:
    - Dialing the WebSocket.
    - Sending a `gps_update` message.
    - Verifying the mock Kafka producer receives the expected event.
