# RR-URG-06: Realtime Notifications And Socket/SSE Broadcasts

## Priority

P2. Depends on stable event contracts and enough backend events to broadcast.

## Problem

The demo needs two-screen visibility. Dashboards must update while the actor performs actions, but socket/SSE must not become the source of truth.

## Scope

- Add persisted notifications or notification module.
- Add socket/SSE realtime service/module.
- Support public sanitized topology stream.
- Support authenticated role-scoped private streams.
- Add display delay/pacing for demo.

## Implementation Details

### 1. Notification Model

Fields:

- `id`
- `scope`
- `target_role`
- `target_user_id`
- `farm_id`
- `warehouse_id`
- `store_id`
- `shipment_id`
- `entity_type`
- `entity_id`
- `type`
- `title`
- `message`
- `status`: `UNREAD`, `READ`, `ACKED`, `RESOLVED`
- `created_at`
- `acknowledged_at`

Endpoints:

- `GET /v1/notifications`
- `POST /v1/notifications/:id/ack`
- `POST /v1/notifications/:id/resolve`

Checklist:

- [x] Persisted notifications.
- [x] Scoped notifications.
- [x] Pickup request notification.
- [x] Outbound dispatch notification.
- [x] Driver assignment notification.
- [x] Delivery/return notification.

### 2. Socket/SSE Stream

Private stream:

- requires JWT.
- uses role/entity scope.
- sends only events the user can see.

Public stream:

- no auth allowed only for sanitized root ArchitectureTopology/public trace visualization.
- must not expose private IDs, payment secrets, user data, or operational action payloads.

Candidate endpoints:

- `GET /v1/realtime/stream`
- `GET /v1/realtime/public/topology`

Checklist:

- [x] Private stream validates JWT via one-time ticket exchange.
- [x] Private stream scopes by role/entity.
- [x] Public stream sanitized.
- [x] Disconnects handled.
- [x] Backpressure simple but safe.

### 3. Events To Broadcast

- `notification.created`
- `warehouse.pickup.requested`
- `logistics.pickup.assigned`
- `logistics.gps.updated`
- `logistics.shipment.status_changed`
- `logistics.pickup.arrived_at_warehouse`
- `warehouse.intake.created`
- `warehouse.dispatch.requested`
- `logistics.delivery.assigned`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

### 4. Demo Pacing

Config:

- `REALTIME_DEMO_DELAY_MS=500`
- `LOGISTICS_SIM_TICK_SECONDS=3`

Checklist:

- [ ] Delay applies only to broadcast/display, not persisted state.
- [ ] Demo viewer can see each transition.
- [ ] Fast-forward action emits audited event.

## Acceptance Criteria

1. Warehouse dashboard updates when harvest creates pickup request.
2. Driver Client receives assigned shipment.
3. Logistics map receives GPS/status updates.
4. Retail dashboard receives incoming delivery updates.
5. Public ArchitectureTopology receives sanitized service/event updates.
6. Private streams require auth.
7. Socket/SSE failure does not corrupt backend state.

## Test Checklist

- [x] Unit: notification scoping.
- [x] Unit: socket event sanitizer.
- [x] Integration: authenticated stream receives assigned event.
- [x] Integration: unauthorized stream denied.
- [x] E2E: Operational realtime awareness flows cover warehouse outbound notifications, store incoming delivery notifications, and role-scope isolation through the production Next runtime with mocked gateway contracts.
- [x] Platform: Live Kafka `warehouse.dispatch.requested` event is consumed by running socket-service and persisted as a scoped Redis/Valkey notification.
- [ ] Manual: two-browser demo updates in order.

## Implementation Evidence

- `socket-service` persists notifications in Valkey and exposes role/entity-scoped APIs under `/v1/realtime/notifications`.
- Kafka events create notification records for pickup, dispatch, driver assignment, delivery/return, and fulfillment failure events.
- Warehouse and Store dashboards now poll API-backed realtime notifications instead of static-only notification lists.
- Verification run:
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/socket-service/... ./apps/auth-service/internal/infrastructure/casbin/...`
  - `npm test`
  - `npm run lint`
  - `npm run build`
  - `npm run test:e2e:realtime-notifications`
  - KrakenD JSON parse check.
- Added fail-fast, flow-oriented E2E guardrails:
  - `start:e2e` binds Next production runtime to `127.0.0.1:3000`.
  - Playwright `webServer` builds/starts the app, waits for `http://127.0.0.1:3000/`, and times out instead of letting tests hang.
  - Playwright default timeout is 30s, navigation timeout is 10s, action/expect timeout is 5s.
  - Flow script: `npm run test:e2e:realtime-notifications` passed 3 Playwright tests.
- Platform live flow:
  - `RUN_PLATFORM_TESTS=1 VALKEY_ADDR=127.0.0.1:6379 KAFKA_BROKERS=127.0.0.1:9094 go test -tags=platform ./apps/socket-service/internal/platform -run TestWarehouseDispatchEventCreatesLiveNotification -count=1 -v` passed.
  - This verifies Kafka -> running socket-service consumer -> Redis/Valkey notification persistence and warehouse index creation with a unique dispatch event.
- Remaining manual demo gap: full two-browser visual review with real logged-in role sessions has not been run yet.
