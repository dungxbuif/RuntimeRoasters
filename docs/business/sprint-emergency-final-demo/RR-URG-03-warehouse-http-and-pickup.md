# RR-URG-03: Warehouse HTTP APIs And Harvest-To-Pickup Flow

## Priority

P1. Depends on RR-URG-01 and RR-URG-02.

## Problem

Warehouse frontend expects HTTP APIs, but warehouse-service currently behaves primarily as a worker/usecase service. Also, harvest events currently become intakes too early. The final demo requires Warehouse to create pickup requests first, then intake only after driver return.

## Scope

- Add Warehouse HTTP app/routes and KrakenD routes.
- Replace normal `HarvestCreated -> Intake` shortcut with `HarvestCreated -> PickupRequest`.
- Add inbound pickup queue and dispatch endpoints.
- Create intake after returned pickup/receipt milestone.

## Implementation Details

### 1. Warehouse HTTP API

Add endpoints:

- `GET /v1/warehouse/intakes`
- `GET /v1/warehouse/batches`
- `GET /v1/warehouse/batches/:id`
- `POST /v1/warehouse/batches`
- `PATCH /v1/warehouse/batches/:id/intake`
- `POST /v1/warehouse/batches/:id/process`
- `POST /v1/warehouse/batches/:id/finalize`
- `GET /v1/warehouse/inventory`

Pickup/dispatch endpoints:

- `GET /v1/warehouse/pickup-requests`
- `POST /v1/warehouse/pickup-requests/:id/dispatch`
- `POST /v1/warehouse/pickup-requests/:id/receive`
- `GET /v1/warehouse/dispatch-requests`
- `POST /v1/warehouse/dispatch-requests/:id/dispatch`

Checklist:

- [ ] Warehouse Gin/HTTP app exists.
- [ ] Routes registered.
- [ ] KrakenD exposes `/v1/warehouse/*`.
- [ ] Frontend service paths match gateway paths.
- [ ] Auth middleware applied.

### 2. Pickup Request Model

Add model/table:

- `id`
- `harvest_id`
- `farm_id`
- `warehouse_id`
- `origin_location_id`
- `quantity`
- `coffee_type`
- `status`
- `notification_id`
- `shipment_id`
- `created_at`
- `dispatched_at`
- `received_at`

Statuses:

- `REQUESTED`
- `DISPATCHED`
- `IN_TRANSIT_TO_FARM`
- `PICKED_UP`
- `RETURNING`
- `ARRIVED_WAREHOUSE`
- `RECEIVED`
- `CANCELLED`

Checklist:

- [ ] Migration added.
- [ ] Repository added.
- [ ] Usecase added.
- [ ] Events emitted on status changes.

### 3. Harvest Worker Change

Current behavior:

- consume harvest event
- create `Intake`

Required behavior:

- consume harvest event
- create `PickupRequest`
- create notification
- publish `warehouse.pickup.requested`

Checklist:

- [ ] Normal flow no longer creates intake immediately.
- [ ] Tests updated to assert pickup request creation.
- [ ] Test helper can create intake directly only for setup.

### 4. Receipt To Intake

On pickup return:

- Logistics emits `logistics.pickup.arrived_at_warehouse`.
- Warehouse marks pickup request `ARRIVED_WAREHOUSE`.
- Warehouse Manager or driver-arrival action triggers receipt.
- Warehouse creates `Intake`.
- Warehouse publishes `warehouse.intake.created`.

Checklist:

- [ ] Receipt endpoint validates warehouse scope.
- [ ] Intake links back to harvest and pickup request.
- [ ] Duplicate receipt is idempotent.
- [ ] Existing processing flow works from the created intake.

## Acceptance Criteria

1. Farm harvest creates pickup request and notification, not intake.
2. Warehouse dashboard can list pickup requests.
3. Warehouse Manager can dispatch pickup to logistics.
4. Driver return/warehouse receipt creates intake.
5. Existing processing/inventory flow works after intake creation.
6. `ADMIN` sees aggregate status but no default operator buttons.
7. `WAREHOUSE_MGR` is scoped to assigned warehouses.

## Test Checklist

- [ ] Unit: HarvestWorker creates pickup request.
- [ ] Unit: duplicate harvest event is idempotent.
- [ ] Unit: receipt creates intake once.
- [ ] Service: list pickup requests scoped by warehouse.
- [ ] Integration: harvest -> pickup request -> dispatch -> return -> intake.
- [ ] Regression: existing processing tests migrated to new intake setup.
