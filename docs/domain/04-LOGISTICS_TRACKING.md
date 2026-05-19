# Business Logic: Logistics Tracking

This document defines logistics movement, vehicle/driver state, route simulation, and realtime dashboard effects.

## 1. Purpose

Logistics is the physical backbone of the production-demo story. It connects:

- Farm harvest -> Warehouse pickup -> Intake.
- Warehouse stock reservation -> Retail delivery -> Completed demand.

The demo intentionally uses seeded routes and authenticated Driver Client simulation. This keeps the flow visible and controllable while still persisting every accepted location/status update in backend services.

## 2. Shipment Types

| Type | Origin | Destination | Trigger |
| :--- | :--- | :--- | :--- |
| `FARM_PICKUP` | Warehouse/base | Farm, then return to warehouse | Warehouse pickup request after harvest |
| `RETAIL_DELIVERY` | Warehouse | Retail store, then required return to warehouse/base | Warehouse dispatch after stock reservation |
| `RETURN_TO_BASE` | Current destination | Warehouse/base | After pickup or delivery completion |

## 3. Driver Client Simulation

Primary demo mode:

1. Driver logs in with a `DRIVER` account.
2. Driver opens assigned shipment.
3. UI loads seeded route points from backend/static route data.
4. Driver clicks Start.
5. Browser animates the vehicle along route points for configured duration, e.g. 30-45 seconds.
6. Browser posts GPS/status updates to backend at each tick.
7. Backend validates that the driver is assigned to the shipment.
8. Backend stores accepted updates, publishes events, and broadcasts realtime updates.
9. Driver confirms milestone buttons: departure, arrival, pickup/delivery, return.
10. Return-to-base is mandatory for both pickup and retail delivery flows.

Backend/script simulation may exist as fallback for unattended demos, but frontend movement must never bypass backend persistence/events.

## 4. Vehicle And Driver State

Vehicle statuses:

- `IDLE`
- `ASSIGNED`
- `EN_ROUTE_PICKUP`
- `LOADING`
- `EN_ROUTE_DROPOFF`
- `UNLOADING`
- `RETURNING_TO_BASE`
- `MAINTENANCE`

Driver statuses:

- `AVAILABLE`
- `ASSIGNED`
- `DRIVING`
- `WAITING_CONFIRMATION`
- `OFFLINE`

## 5. UI Role Behavior

| Role | Logistics UI Visibility | Allowed Actions |
| :--- | :--- | :--- |
| `ADMIN` | Aggregate fleet/shipment health | Create/assign fleet resources, no normal GPS/milestone actions. |
| `WAREHOUSE_MGR` | Assigned warehouse dispatch queues and active shipments | Dispatch vehicle/driver, monitor inbound/outbound progress. |
| `DRIVER` | Assigned shipment only | Start simulation, update GPS, confirm milestones. |
| `FARM_MANAGER` | Pickup status for assigned farms | Watch status only in main demo. |
| `STORE_MGR` | Incoming delivery status for assigned stores | Watch status only in main demo. |

## 6. Realtime Broadcasts

Socket/SSE service should broadcast:

- `logistics.gps.updated`
- `logistics.shipment.status_changed`
- `logistics.pickup.assigned`
- `logistics.pickup.arrived_at_farm`
- `logistics.pickup.loading_confirmed`
- `logistics.pickup.arrived_at_warehouse`
- `logistics.delivery.assigned`
- `logistics.delivery.arrived_at_store`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

Private realtime streams must be authenticated and role/entity-scoped. Public/no-auth realtime is only allowed for sanitized root architecture/topology showcase data.

## 7. Demo Pacing

The route animation and socket broadcast should be slow enough for a reviewer to see the flow:

- Default outbound duration: 30 seconds.
- Default return duration: 45 seconds.
- Tick interval: 2-3 seconds.
- Optional display delay: 500-1000 ms.

## 8. Trace And Audit

Every accepted GPS/status/milestone update should carry:

- `shipment_id`
- `driver_id`
- `vehicle_id`
- `route_id`
- business IDs such as `harvest_id`, `order_id`, `store_id`, `warehouse_id`
- `correlation_id`
- OTel `trace_id` where available

Audit must record dispatch, milestone confirmation, fast-forward, and override actions.
