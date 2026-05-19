# Business Logic: Farm Operations

This document defines how farm-side operations affect the end-to-end production-demo flow.

## 1. Purpose

Farm operations are the origin of the coffee journey. A harvest declaration must create a traceable business event and a warehouse pickup request, but it must not create warehouse intake immediately.

Canonical flow:

1. `FARM_MANAGER` logs in.
2. `FARM_MANAGER` opens the assigned Farm dashboard.
3. User creates a harvest for an assigned farm.
4. Farm service persists the harvest and publishes `farm.harvest.created`.
5. Warehouse creates an inbound pickup request and notification.
6. Farm dashboard shows pickup status once a driver is assigned.
7. Driver handles pickup, loading, and return milestones in the main demo flow.

## 2. Roles

| Role | Farm Visibility | Allowed UI Actions |
| :--- | :--- | :--- |
| `ADMIN` | Aggregate farm counts and health | Create farms, assign managers, view overview. No default harvest/confirmation operator buttons. |
| `FARM_MANAGER` | Assigned farms only | Create harvest, view harvest status, view pickup status. |
| `WAREHOUSE_MGR` | Warehouse pickup queue, not farm editing | Dispatch vehicle/driver for pickup request. |
| `DRIVER` | Assigned pickup shipments only | Start route simulation, update GPS, confirm pickup/loading and return milestones. |

## 3. UI Behavior

Farm dashboard should show:

- Assigned farms.
- Create harvest action.
- Harvest list with pickup status.
- Pickup notification/status after warehouse dispatch.
- Driver progress if the farm has an active pickup.

Farm dashboard should not be the primary simulator. Movement simulation belongs to the authenticated Driver Client.

## 4. State Rules

Harvest statuses:

- `CREATED`: harvest persisted by farm service.
- `PICKUP_REQUESTED`: warehouse has created pickup request.
- `PICKUP_ASSIGNED`: logistics has assigned vehicle/driver.
- `PICKED_UP`: driver confirmed cargo pickup/loading.
- `ARRIVED_WAREHOUSE`: driver returned to warehouse.
- `INTAKE_CREATED`: warehouse receipt created an intake.

The direct shortcut `HarvestCreated -> Intake` is not the canonical product behavior.

## 5. Events

Farm-originated:

- `farm.harvest.created`

Warehouse/logistics events visible to farm:

- `warehouse.pickup.requested`
- `logistics.pickup.assigned`
- `logistics.pickup.arrived_at_farm`
- `logistics.pickup.loading_confirmed`
- `logistics.pickup.return_started`
- `logistics.pickup.arrived_at_warehouse`
- `warehouse.intake.created`

## 6. Authorization

`FARM_MANAGER` can only create and view harvests for assigned farms. If no farm assignment exists, the UI should show no operational farm actions.

Driver confirmation is enough for the main demo. Farm-side dual confirmation may be added later as strict mode, but it is not required for the production-demo flow.
