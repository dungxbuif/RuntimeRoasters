# Missing Implementation Spec: End-to-End Supply Chain Continuity

Date: 2026-05-20

## 1. Purpose

This document captures the missing cross-sprint implementation context before opening the final production-readiness ticket.

This file is the current review working spec for the final production-demo implementation ticket. It now merges the useful findings from the external detailed research and is the file to review next.

The current codebase has several vertical slices implemented independently:

- Farm harvest creation publishes a harvest event.
- Warehouse can consume harvest events into raw intakes and run processing/stock-in logic.
- Retail order SAGA can reserve finished stock.
- Logistics can create shipments from warehouse events and track driver locations.
- Trace/Audit can project events.
- OTel tracing is partially configured through the shared base app, but end-to-end propagation through gateway, Kafka, database spans, and UI-facing trace streams still needs verification.

However, the full operational flow is not yet continuous. The missing part is the physical movement workflow and its dashboard/notification lifecycle across Farm -> Warehouse and Warehouse -> Retail.

Production note: this project production environment is still demo-oriented. Several logistics, processing, payment, and operational flows are expected to run as controlled simulations in production. The implementation should optimize for a complete, reviewable demo flow over preserving incomplete legacy shortcuts.

RR-URG-02 implementation note:
- Kafka domain messages are now required to use CloudEvents JSON format.
- Business IDs stay separate from OpenTelemetry `trace_id`.
- `trace_id` is propagated for observability/correlation only; durable business lookup uses IDs such as `order_id`, `shipment_id`, `harvest_id`, `batch_id`, `store_id`, `farm_id`, `warehouse_id`, `driver_id`, and `vehicle_id`.
- Derived CloudEvents must preserve the incoming `traceid` extension, especially for production-demo simulator flows that may not carry W3C `traceparent` transport headers.
- Demo paid orders use `payment.simulated_completed`; real provider/webhook success may still emit `payment.completed`.
- Canonical topic rename is immediate; runtime code should not consume legacy topics such as `farm.harvest.events`, `warehouse.stock.updated`, `logistics.shipment.assigned`, or `logistics.shipment.delivered`.
- Legacy topic names may remain only in docs where they are explicitly marked as deprecated/forbidden. They must not appear in runtime code, config, tests, seed scripts, or demo evidence expectations.

## 2. Investigation Summary

### 2.1 What Exists

Farm:
- `farm-service` creates harvest records and writes outbox events.
- Harvest event is published as `farm.harvest.created` CloudEvent. Its `data` contains `harvest_id`, `coffee_type`, `origin_code`, `quantity`; its extensions carry durable IDs such as `correlationid`, `traceid`, `harvestid`, and `farmid`.
- Relevant files:
  - `src/apps/farm-service/internal/usecase/harvest_usecase.go`
  - `src/apps/farm-service/internal/infrastructure/event/outbox_relay.go`

Warehouse:
- `warehouse-service` has a Kafka `HarvestWorker`.
- Harvest events are consumed into `Intake` records with status `UNASSIGNED`.
- Processing flow exists after raw material is already inside warehouse:
  - Create production batch from intakes.
  - Simulate roasting.
  - Finalize batch into inventory.
  - Publish stock update.
- Retail order reservation exists against finished inventory.
- Relevant files:
  - `src/apps/warehouse-service/internal/worker/harvest_worker.go`
  - `src/apps/warehouse-service/internal/usecase/intake.go`
  - `src/apps/warehouse-service/internal/usecase/aggregation.go`
  - `src/apps/warehouse-service/internal/usecase/processing.go`
  - `src/apps/warehouse-service/internal/usecase/inventory.go`
  - `src/apps/warehouse-service/internal/usecase/order_reservation.go`

Retail:
- `retail-service` creates orders and publishes `retail.order.created`.
- It consumes SAGA events to update order status.
- Relevant file:
  - `src/apps/retail-service/internal/usecase/service.go`

Payment:
- Payment events exist and currently drive the order reservation sequence.
- Relevant file:
  - `src/apps/payment-service/internal/usecase/service.go`

Logistics:
- `logistics-service` can create shipments from:
  - `warehouse.stock.reserved`
  - `warehouse.inventory.updated`
- It assigns a nearest driver, stores driver GEO coordinates in Valkey, and publishes:
  - `logistics.delivery.assigned`
  - `logistics.delivery.completed`
  - `logistics.gps.updated`
- Relevant files:
  - `src/apps/logistics-service/internal/usecase/service.go`
  - `src/apps/logistics-service/internal/domain/shipment.go`

Frontend:
- Warehouse dashboard exists visually for intake/processing/inventory.
- Logistics dashboard exists visually with a map, route replay, and driver location updates.
- Retail and traceability dashboards exist.
- Relevant files:
  - `src/apps/client-app/src/app/(dashboard)/dashboard/warehouse/page.tsx`
  - `src/apps/client-app/src/app/(dashboard)/dashboard/logistics/page.tsx`
  - `src/apps/client-app/src/services/warehouse.service.ts`
  - `src/apps/client-app/src/services/logistics.service.ts`

Docs:
- `docs/technical/FLOW_SEQUENCES.md` describes a high-level Farm-to-Cup journey.
- `docs/technical/LOGISTICS_SIMULATION.md` describes route generation and a driver simulator concept.
- `docs/technical/KAFKA.md` already calls out Valkey/Redlock for inventory-style shared resource locking.
- `deployments/docker-compose.dev.yaml` already includes Valkey and Cassandra infrastructure.
- `deployments/docker-compose.dev.yaml` now defines SigNoz + ClickHouse + SigNoz OTel Collector for OpenTelemetry waterfall evidence.
- SigNoz UI is exposed at `http://localhost:3301`; OTLP gRPC/HTTP remain `localhost:4317` and `localhost:4318`.
- Kafka remains Apache Kafka without ZooKeeper. The SigNoz ClickHouse coordination service is separate observability infrastructure and must not be used by Kafka.

### 2.2 What Is Missing Or Inconsistent

The user's concern is confirmed.

The current code skips from "harvest created" directly to "warehouse intake exists". It does not model:

- Warehouse receiving a harvest notification and dispatching a vehicle to the farm.
- A physical pickup shipment from farm to warehouse.
- Driver movement from warehouse to farm.
- Driver-owned confirmation for loading, unloading, delivery, and return milestones.
- Driver movement back to warehouse.
- Warehouse receiving confirmation and converting the arrived shipment into intake.
- Configurable simulation timing per segment.
- Equivalent delivery workflow from warehouse to retail after paid order.
- Driver return-to-warehouse lifecycle after retail delivery.

There are also additional implementation gaps:

- Warehouse frontend calls `/v1/warehouse/batches`, `/v1/warehouse/inventory`, etc., but `warehouse-service` currently has no HTTP app/routes and KrakenD has no `/v1/warehouse/*` routes.
- Frontend logistics service uses `/v1/logistics/shipments`, `/v1/logistics/locations`, `/v1/logistics/drivers/location`, but KrakenD currently exposes `/v1/shipments`, `/v1/drivers/location` without `/v1/logistics` prefix, and `logistics-service` does not expose `/v1/locations`.
- `docs/technical/LOGISTICS_SIMULATION.md` references `src/scripts/simulate_drivers.go`, but the repo currently only has `generate_routes.go` and `simulate_harvests.go`.
- Logistics shipment model is destination-centric and does not model route legs, origin location, pickup confirmation, return trip, cargo status, or paired confirmations.
- Notification concepts exist in docs/UI copy, but there is no durable notification/event inbox API for dashboards.
- Trace/Audit can project Kafka events, but missing physical logistics events means the journey timeline is incomplete.
- Existing docs describe the lifecycle at a high level, but sprint-specific implementation drift caused earlier flows to remain disconnected.
- Driver simulation was researched in docs, but implementation status is unclear. Treat `simulate_drivers.go` as a concern to research first; if absent, implement the demo driver simulation intentionally in the Driver Client flow or as a backend/script fallback.
- Current logistics frontend contains mock fallback movement. For the final demo, the Driver Client may intentionally simulate the truck movement on the UI and post GPS coordinates to backend like a real driver device, but only after authenticated login as `DRIVER` and only for assigned shipments.
- Trace IDs and business IDs are not consistently defined. Prefer keeping W3C `trace_id` separate from business IDs, but research should also determine how to reuse OTel/event data to build UI flow diagrams and component/service visualization.
- Kafka trace propagation must be verified explicitly. Standard Kafka producers/consumers do not guarantee W3C `traceparent` propagation unless the code injects/extracts headers.
- Existing Cassandra usage is audit-oriented. A separate trace-history/read-model decision is needed if the UI needs fast replay of full journey spans/events.

## 3. Target End-to-End Flow

### 3.1 Farm Harvest To Warehouse Intake

1. Farm manager records a harvest.
2. Farm service publishes `farm.harvest.created`.
3. Warehouse receives the event and creates an inbound pickup request, not a final intake yet.
4. Warehouse dashboard shows a notification: new harvest needs pickup.
5. Warehouse operator dispatches a vehicle/driver from dashboard.
6. Logistics creates a farm pickup shipment.
7. Driver travels from warehouse to farm using simulated route progress.
8. Dashboard and logistics map show status and location updates.
9. Driver arrives at farm.
10. Driver confirms arrival.
11. Driver confirms pickup/loading status for demo.
12. Backend records the pickup/loading milestone and broadcasts it to subscribed clients.
13. Driver travels back to warehouse.
14. Driver confirms arrival at warehouse; Warehouse dashboard receives the event.
15. Warehouse converts the pickup shipment into an `Intake`.
16. Existing processing, roasting, stock-in, and inventory logic continues.

### 3.2 Warehouse To Retail Store Delivery

1. Retail creates an order or stock need.
2. Payment/order SAGA triggers warehouse reservation.
3. Warehouse reserves finished stock.
4. Warehouse dashboard shows outbound dispatch notification.
5. Warehouse operator assigns vehicle/driver.
6. Logistics creates store delivery shipment.
7. Driver travels from warehouse to retail store.
8. Retail manager sees incoming shipment notification.
9. Driver confirms arrival.
10. Driver confirms delivery status for demo.
11. Shipment status becomes delivered.
12. Driver returns to warehouse/base as a required leg.
13. Retail order is marked completed after delivery confirmation and return leg is recorded according to final status policy.

### 3.3 Lightweight Fleet Management

Add a small vehicle/driver management capability inside warehouse/logistics operations:

- Vehicles:
  - ID, plate number, capacity, status, current warehouse/home base.
- Drivers:
  - ID, name, phone, status, current location, assigned vehicle.
- Driver assignment:
  - Available, assigned, en route to pickup, loading, returning, en route to retail, delivered, returning to base.
- Demo controls:
  - Start/pause simulation.
  - Configure segment durations: e.g. 30s outbound, 45s return.
  - Manually fast-forward status for demo recovery.

## 4. Proposed Domain Additions

### 4.1 New Event Contracts

Add explicit event contracts to `src/pkg/events/contracts.go`.

Farm/Warehouse inbound:
- `farm.harvest.created`
- `warehouse.pickup.requested`
- `logistics.pickup.assigned`
- `logistics.pickup.departed`
- `logistics.pickup.arrived_at_farm`
- `logistics.pickup.driver_confirmed_arrival`
- `logistics.pickup.loading_confirmed`
- `logistics.pickup.cargo_loaded`
- `logistics.pickup.return_started`
- `logistics.pickup.arrived_at_warehouse`
- `warehouse.pickup.received`
- `warehouse.intake.created`

Warehouse/Retail outbound:
- `retail.order.created`
- `payment.completed` or `payment.simulated_completed`
- `warehouse.stock.reserved`
- `warehouse.dispatch.requested`
- `logistics.delivery.assigned`
- `logistics.delivery.departed`
- `logistics.delivery.arrived_at_store`
- `logistics.delivery.driver_confirmed`
- `retail.delivery.received`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

Common:
- `notification.created`
- `notification.acknowledged`
- `logistics.gps.updated`
- `logistics.shipment.status_changed`
- `socket.broadcast.requested`
- `socket.client.connected`

### 4.2 Shipment Model Changes

Current `Shipment` has:
- `order_id`
- `batch_id`
- `origin_warehouse_id`
- `destination_store_id`
- `driver_id`
- `status`
- `destination_address`

Needed:
- `type`: `FARM_PICKUP`, `RETAIL_DELIVERY`, `RETURN_TO_BASE`
- `origin_location_id`
- `destination_location_id`
- `farm_id`
- `harvest_id`
- `warehouse_id`
- `store_id`
- `vehicle_id`
- `route_id`
- `current_leg`
- `status`
- `requires_driver_confirmation`
- `requires_origin_confirmation`: default `false` for production-demo unless stricter business mode is enabled.
- `requires_destination_confirmation`: default `false` for production-demo unless stricter business mode is enabled.
- `confirmation_mode`: `DRIVER_ONLY`, `DUAL_CONFIRMATION`
- timestamps:
  - `assigned_at`
  - `departed_at`
  - `arrived_at_origin_at`
  - `loaded_at`
  - `departed_origin_at`
  - `arrived_destination_at`
  - `received_at`
  - `returned_at`

### 4.3 Notification Model

Add a lightweight notification model for dashboards:

- `id`
- `scope`: `WAREHOUSE`, `FARM`, `RETAIL`, `LOGISTICS`
- `target_role`
- `target_user_id`
- `store_id`
- `farm_id`
- `shipment_id`
- `entity_type`
- `entity_id`
- `type`
- `title`
- `message`
- `status`: `UNREAD`, `READ`, `ACKED`, `RESOLVED`
- `created_at`
- `acknowledged_at`

Notification read/ack endpoints can use polling APIs, but live demo visuals should use socket/SSE broadcasts where practical.

### 4.4 Socket/Realtime Gateway Model

The final demo needs realtime visual delay so viewers can watch state changes across two browser windows.

Add a lightweight socket/realtime service or module:

- Accept authenticated WebSocket/SSE connections for private dashboards.
- Allow public sanitized streams only for root `ArchitectureTopology` showcase.
- Subscribe to logistics, notification, and trace events.
- Broadcast role-scoped updates to connected clients:
  - Warehouse dashboard: pickup requests, dispatch changes, inbound arrival, intake created.
  - Logistics dashboard: shipment status, GPS updates, vehicle availability.
  - Farm dashboard: assigned pickup status.
  - Retail dashboard: incoming delivery status.
  - Driver dashboard/client: assigned shipment and next action.
- Support configurable demo pacing so updates are not too fast for humans watching the demo.

Authentication rule:

- Driver simulation clients must authenticate as `DRIVER`.
- Private dashboard streams should be authenticated when they carry operational data.
- Root architecture/topology stream can remain public only if data is sanitized.

### 4.5 Vehicle And Driver Model

Add a first-class lightweight fleet model owned by Logistics, surfaced in Warehouse/Logistics dashboards:

- `vehicles`
  - `id`
  - `plate_number`
  - `type`
  - `capacity_kg`
  - `home_warehouse_id`
  - `status`: `IDLE`, `ASSIGNED`, `EN_ROUTE_PICKUP`, `LOADING`, `EN_ROUTE_DROPOFF`, `UNLOADING`, `RETURNING_TO_BASE`, `MAINTENANCE`
  - `current_latitude`
  - `current_longitude`
  - `current_location_label`
  - `updated_at`
- `drivers`
  - `id`
  - `user_id`
  - `name`
  - `phone`
  - `vehicle_id`
  - `status`: `AVAILABLE`, `ASSIGNED`, `DRIVING`, `WAITING_CONFIRMATION`, `OFFLINE`
  - `current_shipment_id`

The external research used the term `Trucks`. Prefer `vehicles` in APIs and DB naming so the demo can support vans/trucks later without renaming.

### 4.6 Business IDs, Trace IDs, And Correlation

Prefer not to map Origin Batch ID directly to W3C `trace_id` unless the final traceability design intentionally chooses that tradeoff.

Use separate identifiers:

- Business identifiers:
  - `origin_batch_id`
  - `harvest_id`
  - `intake_id`
  - `production_batch_id`
  - `inventory_lot_id`
  - `order_id`
  - `shipment_id`
- Observability identifiers:
  - `trace_id`
  - `span_id`
  - `parent_span_id`
- Cross-domain correlation:
  - `correlation_id`
  - `causation_id`

Every event should carry both the relevant business IDs and the observability/correlation IDs. The trace timeline UI should query by business ID first, then show the linked trace IDs/spans as technical detail.

UI research note:

- The final UI can use backend domain events as the primary source for business flow rendering.
- OTel data can be reused to discover and visualize which services/components participated in a flow.
- Do not force the user-facing journey to be a raw OTel trace waterfall. Use OTel to enrich service/component topology, latency, and connectivity; use business events to explain the operational story.

### 4.7 Trace History In Trace Service

The external research adds a useful control-plane visualization requirement. Keep it inside `trace-service` instead of defining another monitor microservice.

Extend `trace-service`:

- Ingest telemetry from OTel Collector through an internal endpoint if live span replay is needed.
- Optionally consume selected domain events from Kafka for business milestone visualization.
- Persist high-volume append-only trace history in Cassandra.
- Broadcast live updates to the client through SSE.

Candidate internal endpoint:

- `POST /v1/internal/telemetry/ingest`

Candidate public/demo endpoint:

- `GET /v1/trace/stream?trace_id=...`
- `GET /v1/trace/stream?correlation_id=...`

Auth note:

- The root Client App page is intentionally public because it showcases `ArchitectureTopology` for learning and portfolio/demo purposes.
- That public root page may consume unauthenticated architecture/health/demo visualization data.
- Private operational dashboards and Driver Client streams should use authenticated socket/SSE access when realtime data is exposed.
- Public/no-auth socket/SSE is allowed only for sanitized root architecture/topology showcase data.

Cassandra trace-history candidate schema:

```sql
CREATE TABLE trace_history (
    trace_id text,
    span_id text,
    parent_span_id text,
    service_name text,
    operation_name text,
    start_time timestamp,
    attributes map<text, text>,
    PRIMARY KEY ((trace_id), start_time, span_id)
) WITH CLUSTERING ORDER BY (start_time ASC);
```

Recommended operational defaults:

- TTL: 48 hours for showcase/demo trace history.
- Batch writes: buffer for 1 second or 100 spans, whichever comes first.
- Keep audit logs separate from trace-history. Audit is compliance/immutability; trace-history is replay/visualization.

### 4.8 Telemetry Propagation Requirements

Verify and implement W3C trace propagation end-to-end:

- KrakenD must preserve/generate `traceparent` and forward it downstream.
- HTTP/Gin and gRPC already use shared OTel hooks; verify every service uses the shared base app consistently.
- Kafka producer must inject current trace context into message headers.
- Kafka consumer must extract trace context from message headers before running business handlers.
- GORM/Postgres should use an OTel plugin so DB spans appear under the business request trace.
- Logs should include `trace_id` where available.

Kafka header propagation is a must-have for the showcase because most SAGA transitions cross Kafka boundaries.

### 4.9 Inventory Guard

The external research correctly calls out inventory locking. The codebase already has Valkey infrastructure and docs mention Redlock.

For final implementation:

- Warehouse stock reservation and release must use a distributed lock keyed by SKU or inventory lot.
- Recommended key format: `lock:inventory:{sku}` or `lock:inventory:{warehouse_id}:{sku}`.
- Reservation writes must remain idempotent with event IDs/idempotency keys.
- Concurrent paid retail orders must not oversell available stock.
- Tests should cover concurrent reservations for the same SKU.

### 4.10 Persistence Boundary Decisions

Clarify storage roles before implementing trace/monitor work:

- PostgreSQL is the transactional source of truth for service-owned state:
  - harvests
  - pickup requests
  - shipments
  - vehicles/drivers/assignments
  - orders
  - inventory/reservations
  - workflow statuses
- PostgreSQL JSONB is allowed for flexible metadata and payload snapshots:
  - outbox payloads
  - identity traits
  - small order item snapshots
  - optional trace/event attributes
- JSONB must not hide fields required for authorization, filtering, workflow transitions, or dashboard performance. `store_id`, `farm_id`, `warehouse_id`, `driver_id`, `shipment_id`, `order_id`, `status`, and timestamps should be typed/indexed columns when they affect behavior.
- Elasticsearch is the CQRS read model for user-facing traceability search and business journey timelines.
- Cassandra is append-only history:
  - immutable audit logs
  - optional short-lived trace-service history if implemented with separate TTL/schema
- Valkey is realtime/coordination:
  - driver GEO/liveness
  - idempotency cache
  - inventory distributed locks

Do not make Elasticsearch, Cassandra, or JSONB the source of truth for operational state.

## 5. Backend Implementation Plan

### 5.1 Warehouse Service

Add HTTP API and KrakenD routes for existing frontend needs:

- `GET /v1/warehouse/intakes`
- `GET /v1/warehouse/batches`
- `GET /v1/warehouse/batches/:id`
- `POST /v1/warehouse/batches`
- `PATCH /v1/warehouse/batches/:id/intake`
- `POST /v1/warehouse/batches/:id/process`
- `POST /v1/warehouse/batches/:id/finalize`
- `GET /v1/warehouse/inventory`

Add pickup/dispatch workflow:

- `GET /v1/warehouse/pickup-requests`
- `POST /v1/warehouse/pickup-requests/:id/dispatch`
- `POST /v1/warehouse/pickup-requests/:id/receive`
- `GET /v1/warehouse/dispatch-requests`
- `POST /v1/warehouse/dispatch-requests/:id/dispatch`

Change harvest consumption behavior:

- Current behavior: harvest event directly creates `Intake`.
- Required behavior: harvest event creates `PickupRequest` and notification.
- Intake is created only after warehouse receipt confirmation.

Compatibility decision:

- The old direct `HarvestCreated -> Intake` shortcut can be removed or broken if needed.
- The canonical production-demo flow is `HarvestCreated -> PickupRequest -> Logistics Pickup -> Warehouse Receipt -> Intake`.
- Tests should be updated to the new canonical flow instead of preserving the shortcut as a default behavior.
- A one-off seed/test helper may create intakes directly, but it must not be exposed as the normal product flow.

### 5.2 Logistics Service

Add route/location APIs:

- `GET /v1/logistics/locations`
- `GET /v1/logistics/shipments`
- `GET /v1/logistics/shipments/:id`
- `POST /v1/logistics/shipments/:id/assign`
- `POST /v1/logistics/shipments/:id/depart`
- `POST /v1/logistics/shipments/:id/arrive`
- `POST /v1/logistics/shipments/:id/confirm-load`
- `POST /v1/logistics/shipments/:id/confirm-delivery`
- `POST /v1/logistics/shipments/:id/return`
- `POST /v1/logistics/drivers/location`
- `GET /v1/logistics/drivers`
- `GET /v1/logistics/vehicles`
- `GET /v1/logistics/vehicles/available`
- `POST /v1/logistics/vehicles`
- `PATCH /v1/logistics/vehicles/:id/status`

Add simulation engine:

- Config:
  - `LOGISTICS_SIM_OUTBOUND_SECONDS=30`
  - `LOGISTICS_SIM_RETURN_SECONDS=45`
  - `LOGISTICS_SIM_TICK_SECONDS=3`
  - `LOGISTICS_SIM_AUTO_ADVANCE=true`
- Replay route points from seeded `routes.json` or DB route table.
- Support a Driver Client simulation mode:
  - Driver logs in with a `DRIVER` account.
  - Driver opens assigned shipment.
  - Driver clicks Start.
  - Client animates truck movement on the route for 30-45 seconds.
  - Client posts GPS coordinates/status updates to backend at each tick.
  - Backend validates driver assignment, persists state, publishes GPS/status events, and broadcasts realtime updates.
- Support backend/script simulation as a fallback for unattended demos.
- Publish status and GPS events at each accepted tick.
- Support manual driver confirmation gates.
- Pause movement while status is `LOADING`, `UNLOADING`, or `WAITING_CONFIRMATION`.
- Keep driver location update API available for the `DRIVER` role, but validate the driver is assigned to the shipment.
- Research first whether `src/scripts/simulate_drivers.go` exists in another branch or was only documented. If absent, implement Driver Client simulation intentionally and keep a backend/script fallback only if needed.

### 5.3 Farm Service

Add pickup confirmation endpoint:

- `GET /v1/farms/pickups`
- `POST /v1/farms/pickups/:shipment_id/confirm-loading`

For production-demo, farm-side confirmation is optional. The endpoint can be implemented later for stricter dual-confirmation mode, but the canonical demo flow only requires authenticated driver confirmation.

### 5.4 Retail Service

Add retail receiving endpoint:

- `GET /v1/retail/deliveries`
- `POST /v1/retail/deliveries/:shipment_id/confirm-receipt`

For production-demo, store-side receipt confirmation is optional. The endpoint can be implemented later for stricter dual-confirmation mode, but the canonical demo flow only requires authenticated driver delivery confirmation. Use existing `store_ids` scope if the store endpoint is implemented.

### 5.5 Notification Service Or Module

MVP option:

- Implement notification tables per service where notifications are generated.
- Expose service-specific notification endpoints.

Better option:

- Add `notification-service`.
- Consume domain events and write unified notifications.
- Expose:
  - `GET /v1/notifications`
  - `POST /v1/notifications/:id/ack`
  - `POST /v1/notifications/:id/resolve`

Recommendation for final ticket:

- Use a lightweight unified notification service or module to avoid duplicating notification logic across Farm/Warehouse/Retail.

### 5.6 Trace Service Live Stream

Extend trace-service if the final demo needs live topology and trace playback beyond existing Trace/Audit screens. Do not define a separate monitor microservice for this ticket.

Current infrastructure status:

- Application services already initialize OpenTelemetry and export to `localhost:4317` by default.
- `deployments/docker-compose.dev.yaml` includes SigNoz, ClickHouse, the SigNoz OTel Collector, and a ClickHouse coordination service.
- `deployments/otel-collector-config.yaml` is the collector pipeline used by the SigNoz collector.
- SigNoz UI is available at `http://localhost:3301`.
- RR-URG-01 final evidence must include the same `trace_id` in trace-service and SigNoz.

Responsibilities:

- Receive OTel spans from the collector through an internal ingest endpoint.
- Consume selected business milestone events from Kafka if span-only visualization is not business-friendly enough.
- Persist trace history in Cassandra with short TTL.
- Broadcast live topology/trace updates through SSE.

Implementation notes:

- Extend `deployments/otel-collector-config.yaml` with an additional exporter to trace-service internal ingest if live trace-history from spans is implemented.
- Do not store operational trace-history in the audit Cassandra table unless the schema and retention policy are intentionally shared.
- Root/public architecture visualization may use unauthenticated public/demo data.
- Role-specific operations dashboards should remain behind normal auth. SSE/socket auth is required when those dashboards receive private realtime operational streams.

### 5.7 Kafka, Gateway, And Database Trace Propagation

Add platform-level tasks before claiming end-to-end trace continuity:

- KrakenD: propagate W3C `traceparent` headers.
- Kafka producer: inject trace context into message headers.
- Kafka consumer: extract trace context before handler execution.
- Postgres/GORM: add OTel DB instrumentation.
- Logs: keep `trace_id` in structured logs where available.

Acceptance check:

- A single retail order or harvest journey should produce one connected trace across gateway, service handler, database operations, Kafka publish, Kafka consume, downstream service handler, and final state transition.

### 5.8 Socket/Realtime Service

Add a dedicated socket/realtime service or a clearly isolated module if that is faster.

Responsibilities:

- WebSocket/SSE endpoint for private dashboards.
- Optional public stream for sanitized root architecture showcase.
- Event fan-out from Kafka/domain events to connected clients.
- Role/entity scoping based on JWT claims:
  - `DRIVER`: assigned shipments only.
  - `WAREHOUSE_MGR`: assigned warehouse events.
  - `FARM_MANAGER`: assigned farm events.
  - `STORE_MGR`: assigned store events.
  - `ADMIN`: aggregate/demo overview only by default.
- Configurable artificial display delay:
  - `REALTIME_DEMO_DELAY_MS`, default `500` or `1000`.
  - This delay helps reviewers see each transition during a live demo.
- Backpressure and disconnect handling can be simple for demo; correctness of persisted backend state is more important than guaranteed socket delivery.

## 6. Frontend Implementation Plan

### 6.1 Warehouse Dashboard

Add sections:

- Inbound pickup notifications.
- Dispatch panel for farm pickup.
- Active inbound shipments.
- Receipt confirmation.
- Existing intake/processing/inventory panels.
- Outbound retail dispatch queue.

### 6.2 Logistics Dashboard

Add:

- Real shipment data only; replace generic mock fallback with intentional authenticated Driver Client simulation.
- Driver/vehicle roster.
- Shipment detail drawer with current leg and required confirmation.
- Simulation controls and duration config visibility.
- Driver Client route simulation:
  - Only visible after login as `DRIVER`.
  - Shows assigned shipment, route, current milestone, and next action.
  - Start button animates truck movement and posts GPS/status updates to backend.
  - Buttons for driver milestones: confirm departure, arrived at pickup/store, loaded/picked up, delivered, returned.
- Status timeline:
  - assigned
  - departed
  - arrived
  - confirmed
  - loaded
  - returning
  - received/delivered

### 6.3 Farm Dashboard

Add:

- Pickup arrival notification.
- Confirm loading action.
- View active pickup shipment.

### 6.4 Retail Dashboard

Add:

- Incoming shipment notification.
- Confirm receipt action.
- Show link between retail order and shipment.

### 6.5 Public Root Architecture Showcase

The Client App root page is intentionally public for learning/showcase.

Requirements:

- Root page can show `ArchitectureTopology` without login.
- Public root page can show sanitized architecture, service topology, and demo/health style data.
- Public root page must not expose private user, shipment, payment, store, farm, or operational action data.
- Private realtime dashboard streams should be authenticated and role-scoped.

### 6.6 Demo Guide And Account Flow Documentation

Create a final user-facing demo guide before production review.

The guide must include:

- Existing demo accounts and passwords from seeded identity docs.
- Which browser/window should log in with which account.
- Which buttons to click in each step.
- Which dashboard or topology screen to watch after each action.
- Recommended two-screen setup:
  - Screen A: actor dashboard, e.g. Warehouse Manager or Driver.
- Screen B: Logistics map, Trace page, or public ArchitectureTopology.
- Expected realtime visual changes and approximate delays.
- Recovery steps if a simulation is stuck.
- Known demo shortcuts, such as driver-only confirmation.

Candidate demo scripts:

- Script 1: Farm harvest -> Warehouse pickup -> Intake -> Processing.
- Script 2: Paid retail order -> Warehouse reserve -> Driver delivery -> Driver return -> Retail completed.
- Script 3: Public architecture/topology showcase with live service/event highlights.
- Script 4: Public QR trace demo using seeded products and real trace documents.

## 7. Auth And Authorization Notes

The recent auth work added record-level store scoping for SAGA services. This must be preserved.

Additional authorization required:

- `ADMIN` is setup/assignment/overview, not detailed operator by default.
- `ADMIN` can create manager accounts, farms, warehouses, and retail stores.
- `ADMIN` can assign managers to farms, warehouses, retail stores, and fleet resources.
- `ADMIN` can view aggregate operational health and counts.
- `FARM_MANAGER` can view assigned farm pickup status. Farm-side loading confirmation is optional/future strict mode for this demo.
- `WAREHOUSE_MGR` can dispatch and receive warehouse shipments.
- `DRIVER` can update/confirm only assigned shipments.
- `STORE_MGR` can view incoming delivery status for stores in `store_ids`. Store-side receipt confirmation is optional/future strict mode for this demo.
- Support override for `ADMIN`, if implemented, must be explicit, visually separated, and audited.

Known auth gap to fix before final implementation:

- Casbin matcher currently may not allow direct role policies reliably. It should support direct role match and role inheritance:
  - direct: `r.sub == p.sub`
  - inherited: `g(r.sub, p.sub)`

## 8. Trace And Audit Requirements

Trace timeline must include physical logistics events, not only business state events.

Add projection support for:

- pickup requested
- pickup assigned
- departed
- arrived at farm
- loading confirmed
- cargo loaded
- returned to warehouse
- warehouse received
- delivery assigned
- arrived at store
- retail received
- driver returned

Audit must record:

- dispatch decisions
- manual confirmation actions
- simulation fast-forward actions
- status changes

Telemetry/control-plane requirements:

- Business traceability and OTel tracing must remain distinct but linked.
- Events should carry `correlation_id` and relevant business IDs.
- Kafka messages should preserve `traceparent` headers.
- Derived CloudEvents should also preserve the incoming `traceid` extension so trace-service, audit evidence, and UI flow visualizers can present one coherent chain.
- Trace service may store short-lived trace history in Cassandra for public/demo visualization.
- Root Client App `ArchitectureTopology` can remain public/no-auth with sanitized data.
- Socket/realtime service should broadcast role-scoped dashboard updates, with optional public sanitized stream for root topology.

## 9. Compatibility And Production-Demo Strategy

The old flow can be broken if it blocks the full canonical flow. Backward compatibility is not a product requirement for incomplete shortcuts from earlier sprints.

Production is still a demo environment, so simulation is not a temporary development hack. Simulation is part of the production-demo behavior and must be reliable, configurable, and visible on dashboards.

Canonical behavior:

- Harvest does not directly create intake.
- Warehouse pickup and retail delivery must go through logistics.
- Driver movement is simulated with configurable timing in production-demo.
- The primary demo simulation can run in the authenticated Driver Client: the browser animates seeded route points and posts GPS/status updates to backend like a real driver device.
- Backend remains the persisted source of truth: it validates assignment, stores accepted coordinates/status, emits events, and broadcasts realtime updates.
- Manual driver confirmation gates remain real user actions unless explicitly configured for demo auto-advance.
- Farm/store dual confirmation is optional strict mode, not required for the main demo.
- Existing tests must be migrated to the new flow.

Allowed config flags:

- `LOGISTICS_SIM_OUTBOUND_SECONDS`, default `30`.
- `LOGISTICS_SIM_RETURN_SECONDS`, default `45`.
- `LOGISTICS_SIM_TICK_SECONDS`, default `3`.
- `LOGISTICS_SIM_AUTO_ADVANCE`, default `true` for production-demo movement.
- `LOGISTICS_REQUIRE_MANUAL_CONFIRMATIONS`, default `true` for driver milestones.
- `PAYMENT_SIMULATE_PROVIDERS`, default can remain enabled for demo.

Disallowed as normal product behavior:

- `WAREHOUSE_AUTO_INTAKE_ON_HARVEST=true` as a default.
- Frontend-only fake shipment state that bypasses backend persistence/events.
- Hidden admin bypasses for normal operations.
- Using W3C `trace_id` as the only business batch identifier without an explicit traceability decision.
- Exposing private operational streams without auth outside the public root architecture showcase.

Existing processing and stock-in logic should remain usable once an intake exists, but intake creation must be moved behind warehouse receipt confirmation.

## 10. Candidate Final Ticket Scope

Title:

Finalize End-to-End Farm-to-Retail Logistics Continuity

Acceptance criteria by flow and role:

### 10.1 Admin Setup And Oversight

1. `ADMIN` can create manager user accounts for `FARM_MANAGER`, `WAREHOUSE_MGR`, `STORE_MGR`, `DRIVER`, and optional `PROCESSOR`.
2. `ADMIN` can create farms, warehouses, retail stores, vehicles, and driver records.
3. `ADMIN` can assign managers to farms, warehouses, stores, and drivers/vehicles.
4. `ADMIN` can view aggregate counts and health: farms, warehouses, stores, active harvests, active shipments, inventory summary, failed/blocked processes.
5. `ADMIN` does not see normal operational buttons by default: create harvest, dispatch shipment, confirm loading, confirm receipt, update GPS, finalize batch.
6. Any `ADMIN` support override, if implemented, is explicit and audited.

### 10.2 Farm Harvest And Pickup

1. `FARM_MANAGER` can view only assigned farms.
2. `FARM_MANAGER` can create harvests only for assigned farms.
3. Creating a harvest creates a warehouse pickup request and notification, not immediate intake.
4. `FARM_MANAGER` can view pickup status for assigned harvests.
5. `FARM_MANAGER` can watch pickup status for assigned harvests; farm-side confirmation is optional/future strict mode.
6. `FARM_MANAGER` cannot act on other farms.
7. `WAREHOUSE_MGR` receives pickup notification for the assigned warehouse.
8. `WAREHOUSE_MGR` can dispatch vehicle/driver to farm.
9. `DRIVER` can see only assigned pickup shipment.
10. `DRIVER` can update current GPS/location only for assigned pickup shipment.
11. `DRIVER` can confirm departure, farm arrival, cargo loaded/picked up, return progress, and warehouse arrival for assigned shipment.
12. Authenticated Driver Client can simulate warehouse -> farm movement by replaying seeded route points and posting GPS/status updates to backend.
13. Authenticated Driver Client can simulate farm -> warehouse return with configurable timing.
14. Warehouse dashboard receives realtime return/arrival notification.
15. Driver or Warehouse Manager finalizes warehouse arrival according to the chosen demo button placement.
16. Warehouse receipt creates intake and unlocks existing processing flow.

### 10.3 Warehouse Processing And Inventory

1. `WAREHOUSE_MGR` can view inbound queue, intakes, processing batches, and inventory for assigned warehouse.
2. `WAREHOUSE_MGR` can create production batch from received intakes.
3. `WAREHOUSE_MGR` or `PROCESSOR`, depending on final role split, can start processing simulation.
4. Processing simulation runs with persisted backend state, not only frontend visuals.
5. `WAREHOUSE_MGR` can finalize ready batches into finished inventory.
6. Finished inventory publishes stock update events for trace/audit and downstream demand.
7. `ADMIN` can view warehouse aggregate status but cannot finalize batches by default.

### 10.4 Retail Demand And Store Delivery

1. `STORE_MGR` can view only assigned stores using `store_ids`.
2. `STORE_MGR` can create paid retail orders only for assigned stores.
3. Payment can remain simulated in production-demo, but the flow should still look like a paid order before warehouse reservation.
4. Warehouse reservation runs against finished inventory.
5. `WAREHOUSE_MGR` sees outbound dispatch request after reservation.
6. `WAREHOUSE_MGR` can assign vehicle/driver and dispatch shipment.
7. `DRIVER` can see only assigned retail delivery shipment.
8. `DRIVER` can update current GPS/location only for assigned retail delivery shipment.
9. Authenticated Driver Client simulates warehouse -> retail movement with configurable timing.
10. `DRIVER` confirms arrival at retail store.
11. `STORE_MGR` can watch incoming delivery status for assigned store; store-side receipt confirmation is optional/future strict mode.
12. `DRIVER` confirms delivery completion.
13. Driver return-to-base is mandatory and visible.
14. Retail order reaches completed state after driver delivery confirmation and required return leg recording according to final status policy.

### 10.5 Notifications, Trace, Audit, And UI

1. Warehouse, logistics, farm, and retail dashboards show relevant notifications and statuses.
2. Notifications are persisted and scoped by role/account/entity assignment.
3. Logistics dashboard uses backend shipment state as source of truth; Driver Client route replay must post accepted updates to backend.
4. Trace timeline shows full physical journey from harvest to warehouse to retail.
5. Audit logs capture dispatch, confirmation, simulation fast-forward, admin override, and status changes.
6. Auth scoping is enforced for all new actions.
7. Old processing/inventory/order SAGA tests are updated to pass under the new canonical flow.
8. Socket/realtime updates are delayed or paced enough for a reviewer to see visual state changes during a live demo.

### 10.6 Observability, Trace Stream, And Public Showcase

1. Client App root page can remain public/no-auth to showcase `ArchitectureTopology`.
2. Public root page shows sanitized architecture/topology/demo status only.
3. Private operational dashboards remain authenticated and role-scoped.
4. End-to-end traces preserve `traceparent` across KrakenD, service handlers, DB calls, Kafka publish/consume, and downstream handlers.
5. Kafka messages include trace context headers.
6. Business events include business IDs plus `correlation_id`; using `trace_id` as the only Origin Batch ID requires an explicit traceability decision.
7. Trace-service can stream live topology/trace updates over SSE.
8. If trace SSE is public, it must only expose sanitized demo/control-plane data.
9. Cassandra trace-history, if implemented, uses short TTL and remains separate from immutable audit logs.
10. Research should determine how to combine backend domain events and OTel data to render flow UI: business events explain what happened, OTel shows which services/components participated.

### 10.8 Public QR Trace Demo

1. Public Client App page can generate/show a list of QR codes for pre-seeded products.
2. QR codes point to public trace URLs, for example `/trace/public/:trace_code` or equivalent.
3. Each QR opens a real trace document projected by trace-service/Elasticsearch from prepared seed events/data.
4. Public trace page must show the real journey: farm, pickup, warehouse intake, processing, inventory, paid order, delivery, driver return.
5. Public trace page may include service/component participation from OTel as supporting technical detail.
6. Public trace page must not expose private user/account/payment secrets.

### 10.7 Inventory Concurrency Guard

1. Warehouse stock reservation uses Valkey/Redlock or equivalent distributed locking.
2. Lock scope prevents concurrent oversell for the same warehouse/SKU or inventory lot.
3. Reservation and release are idempotent.
4. Tests cover concurrent retail orders competing for limited stock.

## 11. Open Product Decisions To Confirm

1. Retail demand trigger:
   - Decided: use paid order.
   - Payment may be simulated, but warehouse reservation/delivery starts from a paid or simulated-paid retail order.
   - Replenishment is not part of the final demo ticket unless explicitly reopened later.

2. Confirmation rules:
   - Main demo uses driver-only confirmation for pickup, delivery, and return milestones.
   - Decide whether strict dual-confirmation mode is needed later for farm/store actors.
   - Should explicit audited `ADMIN` support override exist for stuck demo flows?

3. Driver return:
   - Decided: return-to-base is a required shipment leg.
   - Open detail: should retail order be marked completed immediately after delivery confirmation or only after return-to-base is recorded?

4. Notification/realtime implementation:
   - Polling API can handle notification read/ack.
   - Socket/SSE should handle live demo visual updates where practical.

5. Simulation controls:
   - Backend is source of truth.
   - Frontend may start/pause/fast-forward only through backend APIs.

6. Driver simulation implementation:
   - Confirm whether `src/scripts/simulate_drivers.go` exists in another branch/history or was only documented.
   - Primary proposal: authenticated Driver Client owns route animation and posts GPS/status updates to backend.
   - Backend/script simulation can remain fallback for unattended demo runs.

7. Trace/public realtime:
   - Root architecture showcase can be public.
   - Decide later whether private dashboard SSE/WebSocket requires auth when private realtime pages are added.

8. Trace storage:
   - Decided: trace-history belongs to trace-service, not a new monitor-service.
   - Keep Elasticsearch as business traceability read model.
   - Keep Cassandra as append-only audit/monitor history, not source-of-truth state.

9. User demo guide:
   - Prepare final guide after implementation stabilizes.
   - Include current seeded accounts, login sequence, button-by-button demo scripts, expected realtime visuals, and recommended multi-window setup.
   - Include public QR trace demo steps and seeded trace product list.

## 12. Recommendation

Before implementation, update high-level and business docs to make this physical logistics lifecycle the canonical end-to-end flow.

Do not implement new isolated sprint features until this flow is accepted, because it changes the meaning of harvest intake, warehouse readiness, logistics shipments, notifications, and traceability.

## 13. Domain Docs Update Plan

The domain docs must become UI-aware and role-aware, not only backend-flow-aware.

Completed planning updates should cover:

- `docs/domain/01-FARM_OPERATIONS.md`
  - Farm Manager harvest flow.
  - Farm pickup notification.
  - Farm-side loading status visibility; confirmation is optional strict mode.
  - Assigned-farm-only UI behavior.

- `docs/domain/02-PROCESSING_INVENTORY.md`
  - Warehouse inbound notification.
  - Dispatch, receipt, intake, processing, stock-in, outbound dispatch.
  - Warehouse Manager ownership.

- `docs/domain/03-ORDER_FULFILLMENT.md`
  - Store Manager paid order flow.
  - Payment/reservation/delivery receipt.
  - Assigned-store-only UI behavior.

- `docs/domain/04-LOGISTICS_TRACKING.md`
  - Shipment types, driver/vehicle lifecycle, route simulation, confirmation gates, return-to-base.

- `docs/domain/05-ROLE_UI_MATRIX.md`
  - Canonical role boundaries and dashboard visibility.
  - Explicit rule: `ADMIN` is setup/assignment/overview, not detailed operator.

Role correction:

- `ADMIN` creates manager accounts and top-level business nodes.
- `ADMIN` assigns managers to farms, warehouses, and retail stores.
- `ADMIN` views aggregate counts and operational health.
- Detailed data and actions belong to assigned managers.
- Support override, if implemented, must be explicit and audited.
