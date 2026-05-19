# Role And UI Matrix

This document maps user roles to dashboard actions and explains how each UI action advances the end-to-end business flow.

## 1. Role Summary

| Role | Primary Responsibility | Default Scope |
| :--- | :--- | :--- |
| `ADMIN` | Setup, assignment, aggregate overview | Global aggregate only |
| `FARM_MANAGER` | Harvest creation and farm pickup visibility | Assigned farms |
| `WAREHOUSE_MGR` | Pickup dispatch, receipt, processing, inventory, retail dispatch | Assigned warehouses |
| `STORE_MGR` | Store demand and incoming delivery visibility | Assigned `store_ids` |
| `DRIVER` | Assigned shipment simulation and milestone confirmation | Assigned shipments |
| `PROCESSOR` | Optional production execution role | Assigned warehouse/batch |

`ADMIN` is not an all-powerful daily operator by default. Support override, if implemented, must be explicit and audited.

## 2. UI Action Matrix

| UI Area | Action | Allowed Role | Flow Effect |
| :--- | :--- | :--- | :--- |
| Admin | Create manager account | `ADMIN` | Enables scoped operational user. |
| Admin | Create farm/warehouse/store/vehicle/driver | `ADMIN` | Creates business nodes and fleet resources. |
| Admin | Assign manager/fleet | `ADMIN` | Defines record-level scope. |
| Farm Dashboard | Create harvest | `FARM_MANAGER` | Publishes `farm.harvest.created`; creates warehouse pickup request. |
| Warehouse Dashboard | Dispatch pickup | `WAREHOUSE_MGR` | Assigns vehicle/driver; creates logistics pickup shipment. |
| Driver Client | Start route | `DRIVER` | Starts UI route animation and posts GPS/status updates to backend. |
| Driver Client | Confirm pickup/loading | `DRIVER` | Marks cargo picked up; broadcasts pickup milestone. |
| Driver Client | Confirm warehouse return | `DRIVER` | Triggers warehouse receipt/intake flow. |
| Warehouse Dashboard | Create processing batch | `WAREHOUSE_MGR` or `PROCESSOR` | Converts received intake into processing batch. |
| Warehouse Dashboard | Finalize batch | `WAREHOUSE_MGR` | Creates finished inventory and publishes stock update. |
| Retail Dashboard | Create paid order | `STORE_MGR` | Starts payment/reservation/dispatch SAGA. |
| Warehouse Dashboard | Dispatch retail delivery | `WAREHOUSE_MGR` | Assigns vehicle/driver for store delivery. |
| Driver Client | Confirm store delivery | `DRIVER` | Completes delivery in main demo flow. |
| Public Root | View `ArchitectureTopology` | Public/no auth | Shows sanitized architecture/topology/demo data. |
| Trace/Public QR | View business journey | Authenticated roles, or public sanitized view | Shows domain events plus OTel/service participation. |

## 3. Dashboard Effects By Flow

### Farm -> Warehouse

1. `FARM_MANAGER` creates harvest.
2. Warehouse dashboard receives pickup notification.
3. `WAREHOUSE_MGR` dispatches driver.
4. Driver Client starts route and posts GPS.
5. Logistics map updates in realtime.
6. Driver confirms pickup and return.
7. Warehouse receives return event and creates intake.

### Warehouse -> Retail

1. `STORE_MGR` creates paid order.
2. Payment simulation and warehouse reservation run.
3. Warehouse dashboard receives outbound dispatch request.
4. `WAREHOUSE_MGR` dispatches driver.
5. Driver Client simulates delivery route.
6. Retail dashboard watches incoming delivery status.
7. Driver confirms delivery.
8. Retail order reaches completed state after delivery and required return flow.

### Public QR Trace

1. Visitor opens public trace showcase page.
2. Visitor clicks generate/show demo QR codes.
3. Client displays QR list from prepared seed products.
4. Visitor opens a QR trace URL.
5. Trace page loads a real trace document from trace-service/Elasticsearch.
6. Page shows sanitized Farm -> Warehouse -> Retail journey plus optional service/component detail.

## 4. Socket/Auth Rules

- Driver Client simulation requires authenticated `DRIVER`.
- Private dashboard realtime streams require JWT and role/entity scoping.
- Public root architecture stream can be unauthenticated only if sanitized.
- Socket events are display transport only; backend persisted state remains source of truth.

## 5. Data Source Rules For UI

| UI Need | Data Source |
| :--- | :--- |
| Current operational state and allowed actions | Service APIs backed by PostgreSQL source-of-truth rows |
| Live location/status animation | Socket/SSE events backed by accepted backend GPS/status updates |
| Business journey search and traceability timeline | Elasticsearch CQRS read model |
| Immutable audit trail | Cassandra audit logs |
| Architecture/service participation visualization | OTel spans plus selected domain events |
| Driver liveness and current GEO | Valkey GEO/TTL, with persisted shipment state in PostgreSQL |

UI should not query Cassandra directly for normal operational screens. UI should not depend on JSONB fields for authorization or workflow transitions unless those fields are explicitly promoted into API contracts.

## 6. Fail-Closed Rules

- `STORE_MGR` without `store_ids` sees no store records and cannot create demand.
- `FARM_MANAGER` without farm assignment sees no farm operations.
- `DRIVER` without assigned shipment sees no simulation controls.
- `WAREHOUSE_MGR` without warehouse assignment sees no dispatch/processing controls.
