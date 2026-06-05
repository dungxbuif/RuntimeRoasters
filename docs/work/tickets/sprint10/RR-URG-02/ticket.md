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

## Current Status

RR-URG-02 is implemented as the event-contract and storage-boundary foundation for the final demo:

- CloudEvents envelope and canonical topic constants are defined for the final demo event surface.
- Payload structs exist for paid order, warehouse, pickup, delivery, return, GPS, notification, and socket broadcast events.
- Trace-service and audit-service subscribe to the full canonical topic list.
- Runtime SAGA path works for Retail paid order -> Payment simulated completion -> Warehouse stock reservation -> Logistics delivery assignment.
- Client trace timeline topic names are updated for current SAGA events.

Out-of-scope for RR-URG-02 and implemented by later tickets:

- Full API/state-machine implementation for farm pickup driver lifecycle.
- Full API/state-machine implementation for retail delivery driver lifecycle.
- Notification/realtime/socket service fanout implementation.
- Cassandra-backed trace-service live history beyond audit append-only use.
- Legacy topic cleanup in historical sprint/reference docs. Runtime code/config/tests must stay canonical.

Manual verification guide: `docs/work/tickets/sprint10/RR-URG-02-manual-test-guide.md`.

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

- [x] Common CloudEvents envelope defined in `pkg/events.NewCloudEvent`.
- [x] Farm pickup driver lifecycle payload structs defined.
- [x] Payload structs defined for paid order events.
- [x] Payload structs defined for logistics assignment/delivery/GPS events.
- [x] Payload structs defined for warehouse intake/inventory/reservation events.
- [x] Notification/realtime payload structs defined.

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

- [x] Existing runtime producers migrated to canonical names.
- [x] Existing runtime consumers accept canonical names.
- [x] Current runtime old-name aliases removed. Legacy names may appear only in deprecated/reference docs.
- [x] Trace/audit topic lists updated for current canonical topics.
- [x] Future pickup/driver-return/socket topics have canonical contracts; runtime workflows are later tickets.

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

- [x] Current retail/payment/logistics/trace/audit store scoping uses typed IDs where implemented.
- [x] Current workflow transition IDs use typed payload/extensions, not undocumented JSONB-only keys.
- [x] Elasticsearch is documented and used as traceability projection only.
- [x] Cassandra append-only audit exists; trace-service live-history remains outside RR-URG-02.

## Acceptance Criteria

1. Every event emitted in the final demo has event ID, correlation ID, trace ID if available, and relevant business IDs.
2. Trace-service can project a business journey without parsing undocumented payload shapes.
3. Audit-service can record immutable actions without depending on service DB joins.
4. UI can subscribe to stable event names for realtime updates.
5. Storage boundary rules are reflected in migrations/models.

## Test Checklist

- [x] Unit: event envelope validation.
- [x] Unit: event payload serialization/deserialization.
- [x] Contract: paid-order producer payload matches payment/warehouse/logistics consumer expectations.
- [x] Contract: trace-service handles current canonical events.
- [x] Regression: current old SAGA runtime topics migrated; no runtime legacy consumers remain for the current flow.
- [x] Manual contract guide covers current paid-order flow plus future pickup/socket topic projection through trace/audit.

## Manual Test Data

Use `docs/work/tickets/sprint10/RR-URG-02-manual-test-guide.md` as the single manual test script for this ticket.

The guide includes:

- exact accounts for optional UI observation;
- terminal-first setup commands;
- fixed business IDs;
- dynamic current timestamp handling for payment-service backfill protection;
- paid-order SAGA verification commands;
- future pickup/socket topic projection commands;
- expected output for every DB query;
- final pass checklist.
