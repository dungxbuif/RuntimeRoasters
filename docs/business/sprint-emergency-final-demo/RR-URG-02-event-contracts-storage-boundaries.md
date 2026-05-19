# RR-URG-02: Event Contracts, IDs, And Storage Boundaries

## Priority

P0 after RR-URG-01. This blocks reliable implementation across services.

## Problem

Earlier sprint slices use incompatible event names and payloads. The final demo needs stable event contracts carrying business IDs, correlation IDs, trace metadata, and store/farm/warehouse scope fields.

## Scope

- Standardize event names and payload fields.
- Define business IDs vs `trace_id`.
- Enforce storage responsibilities:
  - Postgres for source-of-truth state.
  - JSONB for metadata/payload snapshots only.
  - Elasticsearch for traceability read model.
  - Cassandra for audit and trace-service live history.
  - Valkey for GEO/liveness/locks.

## Implementation Details

### 1. Event Contract Package

Update or create contracts in `src/pkg/events/contracts.go`.

Required common fields:

- `event_id`
- `event_type`
- `occurred_at`
- `producer`
- `correlation_id`
- `causation_id`
- `trace_id`
- `org_id`
- scope fields when relevant:
  - `farm_id`
  - `warehouse_id`
  - `store_id`
  - `driver_id`
  - `vehicle_id`
  - `shipment_id`
  - `order_id`
  - `harvest_id`

Checklist:

- [ ] Common event envelope defined.
- [ ] Payload structs defined for farm pickup events.
- [ ] Payload structs defined for paid order events.
- [ ] Payload structs defined for logistics status/GPS events.
- [ ] Payload structs defined for warehouse intake/inventory events.
- [ ] Payload structs defined for notification/realtime events.

### 2. Canonical Event Names

Farm/Warehouse:

- `farm.harvest.created`
- `warehouse.pickup.requested`
- `warehouse.pickup.received`
- `warehouse.intake.created`

Logistics inbound:

- `logistics.pickup.assigned`
- `logistics.pickup.departed`
- `logistics.pickup.arrived_at_farm`
- `logistics.pickup.loading_confirmed`
- `logistics.pickup.return_started`
- `logistics.pickup.arrived_at_warehouse`

Retail/Paid order:

- `retail.order.created`
- `payment.completed` or `payment.simulated_completed`
- `warehouse.stock.reserved`
- `warehouse.dispatch.requested`

Logistics outbound:

- `logistics.delivery.assigned`
- `logistics.delivery.departed`
- `logistics.delivery.arrived_at_store`
- `logistics.delivery.driver_confirmed`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

Common:

- `notification.created`
- `notification.acknowledged`
- `logistics.gps.updated`
- `logistics.shipment.status_changed`
- `socket.broadcast.requested`

Checklist:

- [ ] Existing producers migrated to canonical names.
- [ ] Existing consumers accept canonical names.
- [ ] Old names removed or mapped explicitly.
- [ ] Trace/audit topic lists updated.

### 3. Storage Boundary Enforcement

Postgres typed columns required for:

- `store_id`
- `farm_id`
- `warehouse_id`
- `driver_id`
- `vehicle_id`
- `shipment_id`
- `order_id`
- `harvest_id`
- `status`
- timestamps used in workflow

JSONB allowed for:

- outbox event payload
- optional attributes
- identity traits
- request snapshots

Checklist:

- [ ] No auth scope depends only on JSONB keys.
- [ ] No workflow transition depends only on JSONB keys.
- [ ] Elasticsearch model is projection only.
- [ ] Cassandra model is append-only.

## Acceptance Criteria

1. Every event emitted in the final demo has event ID, correlation ID, trace ID if available, and relevant business IDs.
2. Trace-service can project a business journey without parsing undocumented payload shapes.
3. Audit-service can record immutable actions without depending on service DB joins.
4. UI can subscribe to stable event names for realtime updates.
5. Storage boundary rules are reflected in migrations/models.

## Test Checklist

- [ ] Unit: event envelope validation.
- [ ] Unit: event payload serialization/deserialization.
- [ ] Contract test: producer payload matches consumer expectation.
- [ ] Contract test: trace-service handles every canonical event.
- [ ] Regression: old SAGA events either migrated or mapped.
