# RR-URG-04: Logistics Shipment Legs, Driver Client Updates, And Mandatory Return

## Priority

P1. Depends on RR-URG-02 and pairs with RR-URG-03.

## Problem

Current logistics is destination-centric. The final demo needs shipment legs, vehicles, assigned drivers, browser-driven Driver Client simulation, validated GPS updates, milestone confirmation, and mandatory return-to-base.

## Scope

- Add vehicle/driver/shipment leg model.
- Add `/v1/logistics/*` HTTP APIs matching frontend/gateway.
- Support Driver Client simulation with authenticated `DRIVER`.
- Enforce mandatory return leg for farm pickup and retail delivery.

## Implementation Details

### 1. Domain Model

Shipment fields:

- `id`
- `type`: `FARM_PICKUP`, `RETAIL_DELIVERY`, `RETURN_TO_BASE`
- `status`
- `current_leg`
- `route_id`
- `origin_location_id`
- `destination_location_id`
- `farm_id`
- `harvest_id`
- `warehouse_id`
- `store_id`
- `order_id`
- `driver_id`
- `vehicle_id`
- timestamps:
  - `assigned_at`
  - `departed_at`
  - `arrived_at_origin_at`
  - `loaded_at`
  - `departed_origin_at`
  - `arrived_destination_at`
  - `delivered_at`
  - `return_started_at`
  - `returned_at`

Vehicle fields:

- `id`
- `plate_number`
- `capacity_kg`
- `home_warehouse_id`
- `status`
- `current_latitude`
- `current_longitude`

Driver fields:

- `id`
- `user_id`
- `vehicle_id`
- `status`
- `current_shipment_id`

Checklist:

- [x] Migrations added through AutoMigrate-compatible domain fields.
- [x] Repositories added through GORM-backed usecase queries.
- [x] Existing shipment data migration handled or seeded reset documented.

### 2. Logistics APIs

Add/align:

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

Checklist:

- [x] KrakenD paths use `/v1/logistics/*`.
- [x] Driver location update requires `DRIVER`.
- [x] Driver can update only assigned shipment.
- [x] Warehouse Manager can dispatch only assigned warehouse shipments.

### 3. Driver Client Simulation Contract

Driver Client sends:

- `driver_id` from identity/session, not trusted from request body alone.
- `shipment_id`
- `lat`
- `lng`
- `heading`
- `speed`
- `route_index`
- `status`
- `occurred_at`

Backend validates:

- user role is `DRIVER`.
- driver record maps to user.
- shipment assigned to driver.
- status transition is valid.
- route index belongs to shipment route.

Checklist:

- [x] Invalid driver cannot post GPS.
- [x] Driver cannot update another shipment.
- [x] GPS accepted updates publish `logistics.gps.updated`.
- [x] Status changes publish canonical pickup/delivery status events.

### 4. Mandatory Return

Farm pickup:

- warehouse/base -> farm
- pickup/loading confirmed by driver
- farm -> warehouse/base return
- returned status required before intake

Retail delivery:

- warehouse/base -> retail store
- delivery confirmed by driver
- retail store -> warehouse/base return
- returned status required before driver/vehicle available

Checklist:

- [x] Shipment cannot become terminal `COMPLETED` until return leg recorded.
- [x] Vehicle cannot become `IDLE` until return leg recorded.
- [x] Driver cannot become `AVAILABLE` until return leg recorded.
- [x] Return events emitted.

### 5. Driver Confirmation API Matrix

These are the canonical driver actions for the main demo. Farm/store dual confirmation is out of scope for this sprint.

| Flow | UI Button | API | Actor | Required State Before | State/Event After |
| :--- | :--- | :--- | :--- | :--- | :--- |
| Farm pickup | Start pickup route | `POST /v1/logistics/shipments/:id/depart` | `DRIVER` | `ASSIGNED` | `IN_TRANSIT_TO_FARM`, `logistics.pickup.departed` |
| Farm pickup | Arrived at farm | `POST /v1/logistics/shipments/:id/arrive` | `DRIVER` | `IN_TRANSIT_TO_FARM` | `ARRIVED_AT_FARM`, `logistics.pickup.arrived_at_farm` |
| Farm pickup | Confirm pickup/loading | `POST /v1/logistics/shipments/:id/confirm-load` | `DRIVER` | `ARRIVED_AT_FARM` | `PICKED_UP`, `logistics.pickup.loading_confirmed` |
| Farm pickup | Start return | `POST /v1/logistics/shipments/:id/return` | `DRIVER` | `PICKED_UP` | `RETURNING_TO_WAREHOUSE`, `logistics.pickup.return_started` |
| Farm pickup | Returned to warehouse | `POST /v1/logistics/shipments/:id/arrive` | `DRIVER` | `RETURNING_TO_WAREHOUSE` | `ARRIVED_WAREHOUSE`, `logistics.pickup.arrived_at_warehouse` |
| Retail delivery | Start delivery route | `POST /v1/logistics/shipments/:id/depart` | `DRIVER` | `ASSIGNED` | `IN_TRANSIT_TO_STORE`, `logistics.delivery.departed` |
| Retail delivery | Arrived at store | `POST /v1/logistics/shipments/:id/arrive` | `DRIVER` | `IN_TRANSIT_TO_STORE` | `ARRIVED_AT_STORE`, `logistics.delivery.arrived_at_store` |
| Retail delivery | Confirm delivery | `POST /v1/logistics/shipments/:id/confirm-delivery` | `DRIVER` | `ARRIVED_AT_STORE` | `DELIVERED`, `logistics.delivery.driver_confirmed` |
| Retail delivery | Start return | `POST /v1/logistics/shipments/:id/return` | `DRIVER` | `DELIVERED` | `RETURNING_TO_BASE`, `logistics.delivery.return_started` |
| Retail delivery | Returned to base | `POST /v1/logistics/shipments/:id/arrive` | `DRIVER` | `RETURNING_TO_BASE` | `RETURNED_TO_BASE`, `logistics.driver.returned_to_base` |

### 6. Warehouse Handoff Points

Warehouse is not the driver simulator. Warehouse owns dispatch and receipt/intake.

| Handoff | Trigger | Warehouse Behavior | Owning Ticket |
| :--- | :--- | :--- | :--- |
| Pickup request created | `farm.harvest.created` | Create `PickupRequest`, notify dashboard | RR-URG-03 |
| Pickup dispatched | Warehouse Manager dispatch button | Select vehicle/driver, ask logistics to create shipment | RR-URG-03 + RR-URG-04 |
| Pickup returned | `logistics.pickup.arrived_at_warehouse` | Mark pickup arrived and enable receipt/intake | RR-URG-03 |
| Intake created | Warehouse receipt action or configured auto-receipt | Create `Intake`, publish `warehouse.intake.created` | RR-URG-03 |
| Retail delivery dispatched | `warehouse.stock.reserved` and dispatch button | Create dispatch request and logistics delivery shipment | RR-URG-05 + RR-URG-04 |
| Driver returned after retail delivery | `logistics.driver.returned_to_base` | Vehicle/driver become available; order can become `COMPLETED` | RR-URG-04 + RR-URG-05 |

### 7. State Machine Summary

Farm pickup shipment:

`ASSIGNED -> IN_TRANSIT_TO_FARM -> ARRIVED_AT_FARM -> PICKED_UP -> RETURNING_TO_WAREHOUSE -> ARRIVED_WAREHOUSE -> COMPLETED`

Retail delivery shipment:

`ASSIGNED -> IN_TRANSIT_TO_STORE -> ARRIVED_AT_STORE -> DELIVERED -> RETURNING_TO_BASE -> RETURNED_TO_BASE -> COMPLETED`

The terminal `COMPLETED` state is invalid until the required return leg is recorded.

## Acceptance Criteria

1. Logistics creates pickup and retail delivery shipments with explicit legs.
2. Driver Client simulation posts GPS updates accepted by backend.
3. Driver milestone buttons update backend state.
4. Driver return-to-base is mandatory.
5. Realtime events emit for GPS, status, pickup, delivery, and return.
6. Store/Farm/Warehouse/Logistics dashboards can display active shipment state.

## Test Checklist

- [x] Unit: shipment state machine valid transitions.
- [x] Unit: invalid transition denied.
- [x] Unit: return-to-base required before completion/availability.
- [x] Unit: driver assignment authorization.
- [x] Service: GPS update scoped to assigned driver.
- [x] Integration: pickup route outbound + return covered at service boundary.
- [x] Integration: retail route outbound + return covered at service boundary.

## Manual Evidence

Use [manual-test-guide.md](./manual-test-guide.md) for live verification commands and expected database evidence.
