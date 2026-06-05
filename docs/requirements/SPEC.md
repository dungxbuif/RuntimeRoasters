---
artifact_type: requirement_spec
id: SPEC
status: active
owner: human
---

# Product Specification
> Current demo truth labels and simulation boundaries live in [**FB-20260528-08 Demo Flow Contract**](../feedbacks/items/FB-20260528-08-demo-flow-contract/DEMO_FLOW_CONTRACT.md).

This document defines the core business specifications of the **Runtime Roasters** platform, focusing on operational processes, control rules, and supply value within the coffee value chain. This serves as a guiding document for strategic project management and handover.

---

## 1. Vision & Strategic Objectives

**Runtime Roasters** digitalizes the entire journey of coffee beans from the farm to the consumer. The primary goal is to establish absolute trust through transparency of origin, strict quality control at every stage, and optimization of coordination among participants: Farms, Processors, Warehousing, and Retail points.

The system ensures that every served product has a complete digital profile, authenticating its history from cultivation to final service.

---

## 2. Key Business Entities

| Entity | Role in the Supply Chain |
| :--- | :--- |
| **Farm** | The originator of raw materials, responsible for seed quality, cultivation processes, and growing region authentication. |
| **Batch** | The identification unit for a volume of products at a specific stage. A batch may change in characteristics (fresh, green, finished) but retains its identification throughout. |
| **Warehouse** | The hub for storage management, preservation, and where value-added conversion operations (Processing, Roasting) occur. |
| **Order** | Represents customer demand, serving as the trigger for inventory reservation and supply coordination. |
| **Shipment** | The physical movement of goods between nodes in the value chain. |
| **Vehicle / Driver** | Fleet resources that are automatically seeded (auto-seeded) when a Warehouse is created (Keep it simple). Driver accounts execute route simulation and confirm shipment milestones. |
| **Notification** | A role-scoped operational prompt shown on dashboards when a user action is required. |
| **Trace / Correlation** | Technical and business identifiers used to connect domain events, telemetry, and UI journey visualization. |

---

## 3. Core Business Rules & Processes

### 3.1. Origin Authentication & Harvesting
The system establishes data discipline right from the start of the supply chain.
- **Growing Region Authentication**: Every harvest batch must be confirmed to match the geographical location of the registered farm and be under the management of authorized personnel.
- **Actual Yield Control**: The system performs automatic reconciliation between declared harvest yields and the farm's expected productivity. Any significant discrepancies are flagged as "Business Anomalies" to prevent the introduction of unverified raw materials.

### 3.2. Processing & Loss Management
The Roasting process is the most critical link in product value transformation.
- **Historical Linking**: When a green coffee batch is put into roasting, the system automatically inherits the entire origin profile from the source batch to the new finished batch.
- **Standard Loss Thresholds**: The system applies rules for allowable weight loss percentages (depending on the roast level). If the actual ratio falls outside the safety threshold, the system requires confirmation and justification from the person in charge to ensure volume transparency.
- **Inventory Rotation Rules**: Priority is given to shipping batches produced earlier (FIFO) to ensure maximum freshness and quality for the customer.

### 3.3. Supply Coordination & Transaction Processing
Ensures absolute consistency between sales and fulfillment capability.
- **Temporary Reservation Mechanism**: As soon as a customer places an order, the system reserves a corresponding amount of stock in the warehouse. This prevents overselling beyond the actual available quantity.
- **Order Status Coordination**: The system automatically coordinates the steps: Inventory Check ➔ Payment Confirmation ➔ Release Order. All these steps must occur synchronously; if any step fails, previous steps are automatically rolled back to restore the original inventory state.

### 3.4. Transportation & Journey Monitoring
- **Auto-seeding Resources**: Vehicles and drivers are automatically generated (auto-seeded) as demo resources when a Warehouse is created, bypassing the need for manual ADMIN creation of individual drivers.
- **Shipment Dispatching**: Warehouse Manager assigns available vehicles and drivers to inbound farm pickup or outbound retail delivery.
- **Driver Client Simulation**: In production-demo mode, the driver logs in with a `DRIVER` account, starts the route simulation from the UI, and the browser posts GPS/status updates to backend as if it were a driver device.
- **Progress Tracking**: The goods' journey is continuously monitored in real time on Warehouse, Logistics, Farm, Retail, Trace, and public architecture views.
- **Demo Confirmation Rule**: The main demo uses driver-only milestone confirmation for pickup, delivery, and return. Farm/store dual confirmation can be added later as a stricter business mode.
- **Mandatory Return**: After both farm pickup and retail delivery, the driver must complete a return-to-base leg. This is not a background status reset.
- **Backend Source of Truth**: The frontend may animate the route, but backend services validate assignment, persist accepted updates, emit events, broadcast realtime state, and drive trace/audit records.

### 3.5. Realtime Control Plane & Traceability
- **Business Events First**: UI flow diagrams should use backend domain events to explain the business story.
- **OTel Enrichment**: OpenTelemetry data may be used to show which services/components participated, latency, and distributed trace connectivity.
- **Public Showcase**: The Client App root page may show `ArchitectureTopology` without login using sanitized topology/demo data.
- **Private Streams**: Private dashboards and driver simulation streams must be authenticated and scoped by role/entity assignment.

### 3.6. Persistence Responsibilities
- **PostgreSQL** is the source of truth for operational state: harvests, pickup requests, shipments, orders, inventory, assignments, and workflow status.
- **PostgreSQL JSONB** is allowed for flexible metadata and event payload snapshots, but role scope, status, entity IDs, and workflow-critical fields must be typed/indexed columns.
- **Elasticsearch** is the CQRS read model for fast business traceability search and timeline views.
- **Cassandra** stores immutable append-only audit logs and may also store short-lived trace-service history if implemented with a separate schema/TTL.
- **Valkey** stores realtime/coordination data such as driver GEO locations, liveness TTLs, idempotency cache, and inventory locks.

---

## 4. Authorization Model & Business Responsibilities

Information access rights and operational authorities are strictly established based on actual job responsibilities:

1.  **System Administrator (`ADMIN`)**: 
    - **Entity Creation**: The sole authority for creating Accounts (with Roles), Farms, Warehouses, and Retail Stores.
    - **Assignment**: Responsible for binding Managers to their respective entities.
    - **Global Overview**: Views aggregate system health and chain integrity.
2.  **Farm Manager (`FARM_MANAGER`)**: 
    - **Operational Authority**: Exclusive right to execute Farm APIs and resources.
    - **Harvesting**: Declares harvests and manages cultivation profiles for assigned farms.
3.  **Warehouse Manager (`WAREHOUSE_MGR`)**: 
    - **Integrated Operations**: Manages warehouse inventory and processing (Roasting).
    - **Logistics Control**: Inherits all logistics management logic. Receives pickup notifications and manages the outbound queue.
    - **Fleet Assignment**: Authority to assign Vehicles and Drivers (created by Admin) to specific shipments.
    - **Simplified Fleet**: No CRUD UI exists for Drivers/Vehicles; they are managed through Seed Data.
4.  **Store Manager (`STORE_MGR`)**: 
    - **Retail focus**: Exclusive authority over retail operations, store demand (Orders), and incoming delivery verification for assigned stores.
5.  **Driver (`DRIVER`)**: 
    - **Execution**: Views only assigned shipments, starts route simulation, and confirms physical milestones (Pickup/Delivery/Return).

**Policy Management Rule:** To keep authorization simple, dynamic policy updates are avoided. To change a user's role or permissions, the user account must be deleted and recreated. Permissions are automatically assigned when entities are bound.

---

## 5. Security Logs & Traceability Reporting

All data in the value chain is aggregated into a "Traceability Report," including:
- **Cultivation Profile**: Seed variety, harvest time, and growing region characteristics.
- **Processing Profile**: Roasting process, loss metrics, and packaging time.
- **Transportation Profile**: Movement history, transit points, and personnel in charge.
- **Realtime Profile**: Socket/SSE milestone updates shown during the demo.
- **Service Profile**: OTel-derived service/component participation for architecture visualization.
- **Public QR Profile**: Sanitized trace documents for retail-sold cups/order items shown to unauthenticated visitors through QR trace URLs.

The system utilizes secure storage techniques to ensure that once these records are confirmed, they cannot be altered, creating an authentic and permanent proof of product quality.

---
*Runtime Roasters Business Specification — Version 1.0 (2026)*
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
# Business Logic: Coffee Batch Lifecycle & Traceability

This document defines the core business rules for managing batches (Batch) within the Runtime Roasters system.

## 1. Definition of a Batch (What is a Batch?)

A "Batch" represents a specific volume of coffee beans processed together in a single stage. Traceability must be preserved from the moment of harvest (Farm) through packaging (Warehouse) and delivery (Retail).

## 2. Identification Rules (Batch ID Generation)

Batch IDs are automatically generated to uniquely identify a batch at each stage.

**Format:** `RR-{ServiceCode}-{OriginCode}-{YYYYMMDD}-{Sequence}`

| Component | Meaning | Example |
| :--- | :--- | :--- |
| `RR` | Project Prefix | `RR` |
| `ServiceCode` | `H` (Harvest), `P` (Processing), `S` (Stocked) | `P` |
| `OriginCode` | Region Code (e.g., `CD`: Cau Dat, `BMT`: Buon Ma Thuot) | `CD` |
| `YYYYMMDD` | Creation Date | `20260515` |
| `Sequence` | Daily sequence number (4 digits) | `0001` |

## 3. Batch Lifecycle

### Stage 1: Harvesting (Farm Service)
- When the Farm Manager confirms a harvest, an `RR-H-...` ID is generated.
- Status: `HARVESTED`.
- Included Information: Bean type, raw weight, harvest date, farm.

### Stage 2: Pickup Request & Inbound Logistics
- Warehouse receives `farm.harvest.created` and creates a pickup request, not an intake.
- Warehouse dashboard shows a pickup notification.
- `WAREHOUSE_MGR` dispatches an available vehicle/driver.
- `DRIVER` starts the Driver Client simulation, posts GPS/status updates, confirms pickup/loading, and confirms return to warehouse.
- Statuses: `PICKUP_REQUESTED`, `PICKUP_ASSIGNED`, `IN_TRANSIT_TO_FARM`, `PICKED_UP`, `RETURNING_TO_WAREHOUSE`.

### Stage 3: Intake & Pre-processing (Warehouse Service)
- Intake is created only after the driver return/warehouse receipt milestone.
- Status: `RECEIVED`.
- Conversion Code: From `RR-H-...` to `RR-P-...` (still maintaining the link to the original Harvest ID).

### Stage 4: Processing (Roasting)
- Changes status to `HULLING` -> `ROASTING`.
- **Weight Loss Rule:** The roasting process reduces weight by 12% - 20%. The system automatically calculates the `Expected Yield` based on the `Intake Weight`.
- If the output weight deviates by more than 5% from the `Expected Yield`, the Operator is required to enter an Anomaly Note.

### Stage 5: Finished Goods Entry (Stocking)
- After processing is complete, the batch is packaged.
- Conversion Code: Changes to `RR-S-...` (Stocked).
- Status: `STOCKED`.
- Inventory data is updated for the corresponding SKU.

## 4. Inventory Rules

- **Available Quantity:** The actual weight available for sale (`Total Stocked` - `Reserved`).
- **Reserved Quantity:** Weight "held" by orders currently in processing (Saga).
- **FIFO (First-In, First-Out):** Priority is given to shipping batches (`RR-S-...`) with older production dates to ensure freshness.
- **Concurrency Guard:** Reservations must use a distributed lock such as Valkey/Redlock keyed by warehouse and SKU or inventory lot.
- **Idempotency:** Reservation and release operations must be idempotent so duplicate SAGA events do not oversell or double-release inventory.

## 5. Warehouse UI And Role Behavior

| Role | Warehouse UI Visibility | Allowed Actions |
| :--- | :--- | :--- |
| `ADMIN` | Aggregate counts, health, active shipment/inventory summary | Create/assign warehouses and managers. No default processing/finalize buttons. |
| `WAREHOUSE_MGR` | Assigned warehouse inbound queue, intakes, batches, inventory, outbound queue | Dispatch pickup/delivery, receive returned shipment, create batch, start/finalize processing. |
| `PROCESSOR` | Processing queue if role split is enabled | Start/update processing steps only. |
| `DRIVER` | Assigned shipment only through Driver Client | GPS/status updates and milestone confirmations. |

Warehouse dashboard sections:

- Pickup request queue.
- Active inbound shipments.
- Receipt/intake creation.
- Processing batches.
- Inventory and reservation status.
- Outbound retail dispatch queue.
- Realtime notifications from socket/SSE service.

## 6. Events

Inbound:

- `warehouse.pickup.requested`
- `logistics.pickup.assigned`
- `logistics.pickup.arrived_at_warehouse`
- `warehouse.pickup.received`
- `warehouse.intake.created`

Processing/inventory:

- `warehouse.batch.created`
- `warehouse.batch.processing_started`
- `warehouse.batch.finalized`
- `warehouse.stock.updated`
- `warehouse.stock.reserved`
- `warehouse.stock.reservation_failed`

Outbound:

- `warehouse.dispatch.requested`
- `logistics.delivery.assigned`

---
**BA Approval Required for any changes to these rules.**
# Business Logic: Retail Order Fulfillment

This document defines the downstream paid retail order flow from store demand to warehouse reservation and driver delivery.

## 1. Purpose

Paid retail orders trigger warehouse reservation and delivery. The production-demo flow may use simulated payment providers, but it should still model the business sequence as paid order -> reservation -> dispatch -> delivery -> driver return.

## 2. Canonical Flow

1. `STORE_MGR` logs in.
2. User opens assigned Retail dashboard.
3. User creates a paid retail order for an assigned store.
4. Retail service publishes order event.
5. Payment remains pending until the demo Stripe webhook Pass/Fail action.
6. Warehouse reserves finished inventory with concurrency guard.
7. Warehouse dashboard shows outbound dispatch request.
8. `WAREHOUSE_MGR` assigns vehicle/driver.
9. `DRIVER` opens assigned shipment in Driver Client.
10. Driver starts route simulation; browser posts GPS/status updates to backend.
11. Driver confirms store arrival and delivery, then returns to base before the
    order becomes completed.
12. Retail order reaches delivery-completed status after driver delivery confirmation.
13. Driver return-to-base is mandatory and must be recorded.

## 3. Roles

| Role | Retail Visibility | Allowed UI Actions |
| :--- | :--- | :--- |
| `ADMIN` | Aggregate store/order counts and health | Create stores, assign store managers, view overview. |
| `STORE_MGR` | Assigned stores through `store_ids` | Create paid retail order, view incoming delivery status. |
| `WAREHOUSE_MGR` | Outbound dispatch queue for assigned warehouse | Assign vehicle/driver after stock reservation. |
| `DRIVER` | Assigned delivery shipments only | Start route simulation, update GPS, confirm arrival/delivery/return. |

Store-side receipt confirmation is optional strict mode. It is not required for the main production-demo flow.

## 4. UI Behavior

Retail dashboard should show:

- Assigned store selector.
- Create paid order action.
- Order/payment/reservation status.
- Incoming delivery notification.
- Delivery progress from Driver Client updates.
- Completed state after driver confirmation.

Retail dashboard should not directly update driver GPS or shipment route state.

## 5. State Rules

Retail order statuses:

- `CREATED`
- `PAYMENT_PENDING` or `PAYMENT_SIMULATED`
- `RESERVED`
- `DISPATCH_REQUESTED`
- `IN_TRANSIT`
- `DELIVERED`
- `COMPLETED`
- `FAILED`

Completion rule for demo:

- Main demo: driver delivery confirmation completes delivery.
- Driver return-to-base is still required after delivery.
- Strict future mode: store-side receipt confirmation may be required before completion.

## 6. Paid Order Decision

Decision: final demo uses paid order.

| Term | Meaning | When To Use |
| :--- | :--- | :--- |
| Paid order | A demand created from a customer/store order after payment success or simulated payment success. | Required for final demo. |
| Replenishment request | A store restock request independent from a direct customer payment. | Out of scope for final demo unless reopened later. |

Implementation guidance:

- The downstream flow is: paid order -> payment success or simulated success -> warehouse reservation -> dispatch -> driver delivery -> driver return.
- The UI should label the action as paid order or checkout/order placement, not generic replenishment.
- Replenishment can reuse the downstream pipeline in the future, but it is not part of this ticket.

## 7. Events

Retail:

- `retail.order.created`
- `retail.delivery.received`

Warehouse/logistics:

- `warehouse.stock.reserved`
- `warehouse.stock.reservation_failed`
- `warehouse.dispatch.requested`
- `logistics.delivery.assigned`
- `logistics.delivery.departed`
- `logistics.delivery.arrived_at_store`
- `logistics.delivery.driver_confirmed`
- `logistics.delivery.completed`
- `logistics.driver.returned_to_base`

## 8. Authorization

`STORE_MGR` can only see and create demand for assigned stores. Missing `store_ids` should fail closed: no store data and no create action.
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
# Logistics Simulation Guide

This document explains how the real-time logistics simulation works in Runtime Roasters, from seeding locations to visualizing driver movement on actual road paths.

---

## 1. Seeding Strategy

The system uses a set of predefined locations in Vietnam to provide a realistic "Farm-to-Cup" experience. These locations are seeded into the `logistics-service` database.

### Core Locations (POIs)
- **Farmers:** Specific farms in Cầu Đất (Dalat), Buôn Ma Thuột, and Pleiku.
- **Roasteries:** Processing centers in major industrial zones (Sóng Thần, Hòa Lạc, Hòa Khánh).
- **Retailers:** Retail stores in central districts of HCM, Hanoi, and Da Nang.

**Runtime Seed File:** `src/apps/logistics-service/internal/seed/logistics.json`

`deployments/logistics-seed.sql` is legacy/reference only. Do not use it as the runtime source of truth unless the logistics seed strategy is explicitly reworked.

---

## 2. Realistic Route Generation (OSRM)

Unlike simple straight-line movement, our simulation uses actual road data. We use the **Open Source Routing Machine (OSRM)** to calculate paths between our seeded locations.

### The Route Generator Tool
Located at `src/scripts/generate_routes.go`, this tool:
1. Iterates through all logical pairs of locations (e.g., Farm -> Roastery).
2. Calls the OSRM Public API: `http://router.project-osrm.org/route/v1/driving/{lng,lat;lng,lat}`.
3. Extracts the coordinate list (polyline) representing the actual road path.
4. Exports the results to a static JSON file: `src/apps/logistics-service/testdata/routes.json`.

**Why static files?** To ensure the simulation remains stable, offline-capable for demos, and prevents hitting rate limits on public APIs during development.

---

## 3. Real-time Movement Simulation

The primary production-demo simulator is the authenticated **Driver Client**. The driver browser animates seeded route points and posts GPS/status updates to backend like a real driver device.

### How it works:
1. **Login:** User authenticates as `DRIVER`.
2. **Assigned Shipment:** Driver Client fetches only assigned shipments.
3. **Load Paths:** UI reads the seeded route from `routes.json` or an API backed by the same route data.
4. **Replay Path:** UI iterates through the coordinate list at a configured demo speed, usually 30-45 seconds per leg.
5. **Update Location:** Every few seconds, UI posts the current coordinate and shipment status to Logistics Service.
6. **Backend Validation:** Logistics Service validates that the driver is assigned to the shipment before accepting the update.
7. **Freshness Tracking:** Logistics Service stores current coordinates in **Valkey GEO** and sets a TTL key to track driver liveness.
8. **Persistence:** Shipment state and milestone confirmations remain persisted in Postgres.

`src/scripts/simulate_drivers.go` was researched/documented as a possible standalone simulator, but its implementation must be verified before use. If absent, keep Driver Client simulation as the main demo path and implement a backend/script fallback only for unattended demos.

---

## 4. UI Visualization Flow

The Frontend leverages this data to create a modern, animated tracking experience.

1. **Static Routes:** The UI fetches the path coordinates (via API or static asset) to draw lines on the map.
2. **Live Updates:** Dashboards connect to the socket/realtime service or trace-service SSE/WebSocket when private operational data is streamed.
3. **Event Pipeline:** 
   - Driver Client -> Logistics Service (`POST /v1/logistics/drivers/location` or equivalent)
   - Logistics Service -> Postgres shipment state + Valkey GEO/liveness
   - Logistics Service -> Kafka (`logistics.gps.updated`, `logistics.shipment.status_changed`)
   - Socket/Trace Service -> Consumer (Kafka) -> Client (SSE/WS)
4. **Rendering:** The UI updates the driver icon's position on the map based on the live coordinates, snapping them to the predefined road path for a smooth experience.

---

## 5. Running the Simulation

1. **Start logistics-service:**
   ```bash
   cd src/apps/logistics-service
   air
   ```
   On startup the service loads `src/apps/logistics-service/internal/seed/logistics.json` and idempotently inserts vehicles, drivers, and locations.
2. **Generate routes (if POIs changed):**
   ```bash
   go run src/scripts/generate_routes.go
   ```
3. **Run the demo simulation:**
   - Log in with `driver@runtimeroasters.com`.
   - Open the assigned Driver/Logistics shipment screen.
   - Click Start to replay the seeded route.
   - Watch Logistics/Warehouse/Trace/Architecture screens receive realtime updates.

4. **Optional unattended simulator:**
   ```bash
   go run src/scripts/simulate_drivers.go
   ```
   This command is optional and only valid if the script exists in the current implementation.
# Public QR Trace Demo

This document defines the public traceability demo shown to unauthenticated users.

## 1. Purpose

Runtime Roasters is also a learning/showcase project. The public Client App should let visitors inspect the journey of a coffee unit that was actually sold by a retail store, such as a cup label or receipt/order item, without logging in.

The public trace page must show real sold-item trace data, not a purely mocked visual and not a generic product catalog profile.

## 2. Public User Flow

1. Visitor opens the public Client App trace showcase page.
2. Visitor clicks a button such as "Show Demo QR Codes".
3. Client displays a list of QR codes for pre-seeded sold cups/order items.
4. Visitor scans or clicks a QR.
5. Public trace page opens a trace document.
6. Trace document shows the purchased item's journey:
   - farm and harvest
   - warehouse pickup
   - driver return to warehouse
   - intake
   - processing
   - inventory
   - paid retail order
   - retail delivery
   - driver return to base

## 3. Data Source

The public page should use deterministic retail sale seed data:

- Seeded farms, warehouses, stores, products, shipments, paid orders, and sold order items.
- Seeded sold items must link to reserved inventory lots and upstream batch/harvest lineage.
- Seeded or replayed domain events consumed by trace-service.
- Trace-service projects the business journey into Elasticsearch.
- Optional OTel/service participation can be attached as technical detail.

Public trace should query trace-service or a public-safe API backed by Elasticsearch. It should not read Cassandra or internal service databases directly.

QR codes must not be generated for generic unsold catalog products. A QR code is meaningful only when it identifies a public trace document for a sold retail unit/order item.

## 4. QR Content

QR codes should point to stable public trace URLs:

- `/trace/public/:trace_code`
- or `/qr/:trace_code`

The QR payload should not expose internal primary keys. Use a UI-issued `cup_id` as the public QR token; backend persists it on the sold sale item and trace-service exposes the same value as public `trace_code`.

## 5. Visibility Rules

Allowed on public trace page:

- farm name/region
- sold item display label
- product/batch display code
- high-level processing milestones
- logistics milestones and route summary
- store/city display
- timestamps suitable for demo
- service/component participation if sanitized

Not allowed:

- private user IDs
- raw JWT/session data
- payment provider secrets
- internal webhook payloads
- admin/audit-only details

## 6. Role Relationship

This page is public/no-auth. It is separate from operational dashboards:

- `ADMIN` prepares data and assignments.
- `FARM_MANAGER`, `WAREHOUSE_MGR`, `STORE_MGR`, and `DRIVER` create the real operational journey.
- Public visitor only reads sanitized trace output.
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
3. Client displays QR list from seeded sold cups/order items.
4. Visitor opens a QR trace URL.
5. Trace page loads a real trace document from trace-service/Elasticsearch.
6. Page shows sanitized Farm -> Warehouse -> Retail journey plus optional service/component detail.

## 4. Socket/Auth Rules

- Driver Client simulation requires authenticated `DRIVER`.
- Private dashboard realtime streams require JWT and role/entity scoping.
- Public topology demo uses trace-service pull APIs for config/history and
  socket-service WebSocket for push updates. Socket-service has no durable DB;
  Valkey stores only ephemeral session/presence/reconnect state.
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
# Runtime Roasters — Product Requirements Document (PRD)

| Field            | Value                                                    |
| :--------------- | :------------------------------------------------------- |
| **Project**      | Runtime Roasters                                         |
| **Author**       | PM (AI-assisted)                                         |
| **Status**       | Draft v1.2                                               |
| **Last Updated** | 2026-04-16                                               |
| **Type**         | Technical Showcase                           |
| **Developer**    | Solo                                                     |
| **Deployment**   | Docker on Cloud VPS (HA-ready design)                    |

---

## 1. Executive Summary

**Runtime Roasters** is a `Showcase` project simulating a coffee supply chain management platform from farm to retail cup (Farm-to-Cup). The primary objective is to demonstrate the ability to design and implement a production-grade distributed `Microservices` architecture, using 100% `Golang` for the `Backend` and `ReactJS` for the `Frontend`.

The project is **not** intended to solve a real-world commercial problem but instead focuses on:
- Showcasing complex `Design Patterns` (Saga, CQRS, Outbox, Event Sourcing).
- Proving the capability to handle `Distributed Systems`.
- Providing a visual demonstration of the architecture through a "Control Plane Visualization" `Dashboard`.

**Evaluation Audience:** Employers, Technical Reviewers, TAs (Technical Architects).

---

## 2. Goals & Non-Goals

### 2.1 Goals
| #  | Goal                                                                                   | Measurement                                         |
| :- | :------------------------------------------------------------------------------------- | :----------------------------------------------------- |
| G1 | Showcase `Microservices` architecture with 9+ independent services communicating via `Kafka`/`gRPC` | All services running and communicating on `Docker Compose` |
| G2 | Successfully implement the `Saga Pattern` (Choreography) with `Compensating Actions`    | Demo rollback scenarios when `Warehouse` is out of stock + Auto-Refund via `Stripe` |
| G3 | Implement `Transactional Outbox` + `Inbox` (Idempotency)                               | Demo disconnecting `Kafka`, ensuring data remains consistent after recovery |
| G4 | Implement `CQRS` with `PostgreSQL` (Write) + `Elasticsearch` (Read)                     | Coffee traceability retrieval in < 100ms               |
| G5 | Visual architecture dashboard (Control Plane Visualization)                            | Viewers see real-time data flow between services       |
| G6 | Design for `High Availability` (HA-ready)                                              | Architecture can scale horizontally without refactoring |
| G7 | Integrate `Payment Gateway` (Stripe) with full `Webhook Security` + `Idempotency`      | Demo B2B payments and auto-refunds on Saga failure     |

### 2.2 Non-Goals
- ❌ No native mobile app (responsive web only).
- ❌ No processing of GPS data from real hardware devices (using simulated data).
- ❌ No complex financial/accounting management system (basic payment flow only).
- ❌ No multi-tenant targeting (single-tenant demo only).
- ❌ No real money processing in the demo environment (using `Stripe Test Mode`).

---

## 3. User Roles & Personas

| Role         | System Code  | Short Description                                       | Primary Permissions                                          |
| :----------- | :----------- | :------------------------------------------------------- | :----------------------------------------------------------- |
| **Farm Manager** | `FARM_MANAGER` | Manages coffee farms, declares harvests                | CRUD farms, create harvest batches, view history             |
| **Processor** | `PROCESSOR`  | Manages processing factory, roasting                    | Receive raw beans, create roast batches, issue `Batch ID`, packaging |
| **Driver**   | `DRIVER`     | Drives transport vehicles, updates GPS                  | Accept trips, update delivery status, send GPS coordinates    |
| **Store Manager** | `STORE_MGR` | Manages retail stores                                  | View inventory, create supply requests, receive goods        |
| **Admin**    | `ADMIN`      | System Administrator                                   | View full `Dashboard`, `Audit log`, manage users             |
| **Customer** | `end_user` | End consumer (no login required)                       | Scan QR, view traceability                                   |

---

## 4. Business Flows

### 4.1 Main Flow: Farm-to-Cup Pipeline

```
[Farm Manager]  [Processor]       [Warehouse]    [Payment]     [Logistics]    [Retail]      [End User]
   │                │                  │              │              │             │              │
   ├─ Harvest ────►│                  │              │              │             │              │
   │ (Harvest       │                  │              │              │             │              │
   │  Created)      ├─ Roasting ─────►│              │              │             │              │
   │                │ (BatchProcessed) │              │              │             │              │
   │                │                  ├─ Intake ────►│             │             │              │
   │                │                  │(InventoryAdd)│              │             │              │
   │                │                  │              │              │             ├─ Order       │
   │                │                  │              │              │             │(OrderCreated)│
   │                │                  │              │◄─────────────┤─────────────┤              │
   │                │                  │              │ Payment      │             │              │
   │                │                  │              │(PaymentIntent│             │              │
   │                │                  │              │  Created)    │             │              │
   │                │                  │              │──► Stripe ──►│             │              │
   │                │                  │              │  (Webhook)   │             │              │
   │                │                  │◄─────────────┤──────────────┤             │              │
   │                │                  │ Reserve stock│              │             │              │
   │                │                  │(StockReserved│              │             │              │
   │                │                  │              │              │◄────────────┤              │
   │                │                  │              │              │ Dispatch    │              │
   │                │                  │              │              ├─ GPS ──────►│              │
   │                │                  │              │              ├─ Delivery ─►│              │
   │                │                  │              │              │             │              │
   │                │                  │              │              │             │  Scan QR ────┤
   │                │                  │              │              │             │  Traceability│
```

### 4.2 Saga Flow: Supply Order (with Payment)

This is the most complex flow, implementing the `Saga Pattern` (Choreography) combined with a `Payment Gateway`:

| Step | Service          | Action                                | Event Emitted                | Failure → Compensation                           |
| :--- | :--------------- | :------------------------------------ | :--------------------------- | :----------------------------------------------- |
| 1    | `Retail`         | Store creates supply order            | `SupplyOrderCreated`         | —                                                |
| 2    | `Payment`        | Create `PaymentIntent` on Stripe      | `PaymentIntentCreated`       | —                                                |
| 3    | Frontend         | Display Stripe payment form           | —                            | User cancels → `PaymentCancelled` → Cancel order |
| 4    | `Payment`        | Receive Stripe Webhook, confirm payment| `PaymentCompleted`           | `PaymentFailed` → Cancel order                   |
| 5    | `Warehouse`      | Check & reserve stock                 | `StockReserved`              | `StockReserveFailed` → Refund Stripe + Cancel order |
| 6    | `Logistics`      | Find vehicle & create trip            | `ShipmentAssigned`           | `ShipmentFailed` → Release stock + Refund Stripe  |
| 7    | `Logistics`      | Driver accepts and transports         | `ShipmentInTransit`          | —                                                |
| 8    | `Logistics`      | Successful delivery                   | `ShipmentDelivered`          | —                                                |
| 9    | `Warehouse`      | Confirm stock deduction               | `StockDeducted`              | —                                                |
| 10   | `Retail`         | Receive goods, update inventory       | `SupplyOrderCompleted`       | —                                                |

**Compensating Actions (Rollback):**
- If Step 5 fails (out of stock) → `Payment` receives `StockReserveFailed`, calls `Stripe Refund API`, emits `PaymentRefunded`. `Retail` receives it → order changes to `REJECTED`.
- If Step 6 fails (no vehicle available) → `Warehouse` receives `ShipmentFailed`, releases stock. `Payment` receives the event, calls `Stripe Refund API`.
- **Principle:** If money was deducted on Stripe, there must be a path to refund it. Never leave money "hanging."

### 4.3 Payment Flow Details

#### 4.3.1 Payment Sequence

```
  [Frontend/POS]         [Payment Service]           [Stripe]             [Kafka]
       │                        │                       │                    │
       │── POST /orders ───────►│                       │                    │
       │   (SupplyOrderCreated  │                       │                    │
       │    from Kafka)         │                       │                    │
       │                        ├── Create PaymentIntent►│                   │
       │                        │◄── client_secret ──────┤                   │
       │◄── client_secret ──────┤                       │                    │
       │                        │                       │                    │
       │── Stripe.js submit ───►│───────────────────────►│                   │
       │   (card form)          │                       │                    │
       │                        │                       │                    │
       │                        │◄── Webhook POST ──────┤                   │
       │                        │   (payment_intent.     │                   │
       │                        │    succeeded)          │                   │
       │                        │                       │                    │
       │                        ├── HMAC verify ────────►│ (validate sig)    │
       │                        ├── Check Inbox ────────►│ (idempotency)     │
       │                        ├── Save + Outbox ──────►│                   │
       │                        │                       │  ──► PaymentCompleted
       │                        │                       │                    │
```

#### 4.3.2 Payment State Machine

```
                    ┌──────────────────────────────────────────┐
                    │                                          │
  ┌─────────┐   PaymentIntent   ┌──────────────┐              │
  │ CREATED ├──── created ─────►│   PENDING    │              │
  └─────────┘                   └──────┬───────┘              │
                                       │                      │
                          ┌────────────┼────────────┐         │
                     Webhook OK    Webhook FAIL   User cancel  │
                          │            │            │         │
                   ┌──────▼──────┐ ┌───▼──────┐ ┌──▼──────┐  │
                   │  SUCCEEDED  │ │  FAILED  │ │CANCELLED│  │
                   └──────┬──────┘ └──────────┘ └─────────┘  │
                          │                                   │
                   Saga Failure                                │
                   (Stock/Ship fail)                           │
                          │                                   │
                   ┌──────▼──────┐                            │
                   │  REFUNDED   ├────────────────────────────┘
                   └─────────────┘
```

**`Payment` Statuses:**

| Status       | Description                                                 |
| :----------- | :---------------------------------------------------------- |
| `CREATED`    | Order just created, waiting for `PaymentIntent` creation    |
| `PENDING`    | `PaymentIntent` created on Stripe, waiting for user payment |
| `SUCCEEDED`  | Stripe confirms successful payment (via Webhook)            |
| `FAILED`     | Stripe reports payment failure (card declined, etc.)        |
| `CANCELLED`  | User cancels payment before completion                      |
| `REFUNDED`   | Refunded via Stripe Refund API (due to Saga compensation)   |

#### 4.3.3 Payment Service Database Schema

```sql
-- Main Table: Record payments
CREATE TABLE payments (
    id                       UUID PRIMARY KEY,
    order_id                 UUID NOT NULL UNIQUE,
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    amount                   DECIMAL(12,2) NOT NULL,
    currency                 VARCHAR(3) DEFAULT 'VND',
    status                   VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    stripe_refund_id         VARCHAR(255),
    created_at               TIMESTAMPTZ DEFAULT NOW(),
    updated_at               TIMESTAMPTZ DEFAULT NOW()
);

-- Inbox: Anti-duplication for Webhook processing (Idempotency)
CREATE TABLE inbox_stripe_events (
    event_id    VARCHAR(255) PRIMARY KEY,  -- Stripe event ID (evt_xxx)
    event_type  VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);

-- Outbox: Ensure events reach Kafka
CREATE TABLE outbox_events (
    id            UUID PRIMARY KEY,
    aggregate_id  UUID NOT NULL,
    event_type    VARCHAR(100) NOT NULL,
    payload       JSONB NOT NULL,
    published     BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);
```

#### 4.3.4 Webhook Security (HMAC Validation)

Every Webhook from Stripe must pass through `HMAC` middleware before processing:

1. The demo client sends authenticated `POST /v1/webhooks/stripe` with the
   `Stripe-Signature` header.
2. The `Payment Service` loads the active demo signing key from
   `payment_webhook_keys` and verifies the raw body using `crypto/hmac` (Go).
3. Compares the calculated hash with `Stripe-Signature` → if they match, process; otherwise, reject (`HTTP 403`).
4. **Principle:** Never trust any request to `/webhooks/*` without `HMAC verification`.

> **Future Expansion:** The `Payment Service` architecture is designed following the `Payment Gateway Abstraction`. Stripe is the first implementation. Other providers (VNPay, MoMo, PayOS) can be added by implementing the same interface without affecting the Saga flow.

### 4.4 QR Traceability

When an `end_user` scans the QR code on a retail-sold cup label or receipt/order item, the system returns a public trace document for that purchased item:

| Information           | Source Service  | Example                                    |
| :-------------------- | :-------------- | :---------------------------------------- |
| Farm Name             | `Farm`          | "Son La Farm - Anh Minh"                  |
| Cultivation Area      | `Farm`          | "Plot A3, Son La, altitude 1200m"         |
| Coffee Variety        | `Farm`          | "Arabica Catimor"                         |
| Harvest Date          | `Farm`          | "2026-01-15"                              |
| Processing Method     | `Process`       | "Washed Process"                          |
| Roasting Date         | `Process`       | "2026-02-01"                              |
| Roast Level           | `Process`       | "Medium Roast"                            |
| Batch ID              | `Process`       | `RR-SL-2026-0215-A3`                      |
| Sold Item             | `Retail`        | "Cup/order item sold at Runtime Cafe Hanoi" |
| Transport Route       | `Logistics`     | "Son La → Hanoi (320km, 6h)"              |
| Shop Intake Date      | `Retail`        | "2026-02-03"                              |

**Batch ID format:** `RR-{REGION_CODE}-{YEAR}-{MMDD}-{LOT}`

---

## 5. System Architecture Overview

### 5.1 Architectural Foundations

#### Codebase: Monorepo
The entire project resides in a single repository with a shared structure:
- **`api/`** — Shared `gRPC` `Proto` definitions (contracts between services).
- **`pkg/`** — Shared libraries: `Database Wrapper`, `Middleware` (`Auth`, `Tracing`, `Logging`), `Kafka Producer/Consumer` base.
- **`services/`** — Each Microservice is an independent module.

#### Service Structure: Clean Architecture (go-clean-arch v4)
Each service applies `Clean Architecture` with 3 distinct layers:

| Layer             | Content                                                                       |
| :---------------- | :---------------------------------------------------------------------------- |
| **Domain**        | Entities, Value Objects, Repository Interfaces. Completely framework-agnostic. |
| **UseCase**       | Business logic, Saga step handlers, Command/Query handlers.                   |
| **Infrastructure**| `Gin` HTTP handlers, `gRPC` servers, `PostgreSQL` repos, `Kafka` producers.   |

#### Infrastructure: Polyglot Persistence
Multi-modal storage, using the right tool for each data type:

| Database          | Role                                                          |
| :---------------- | :------------------------------------------------------------ |
| `PostgreSQL`      | Core `ACID` transactions — `Source of Truth` for each Microservice |
| `Apache Cassandra`| Perpetual raw Event storage (`Audit` / `Event Sourcing`)      |
| `Elasticsearch`   | High-speed lookup, `CQRS Read Model` for traceability         |
| `Valkey`          | `Cache`, `Distributed Lock`, real-time `GPS` coordinates       |

#### Deployment: Docker on Proxmox
Deployed 100% via `Docker Compose`. The host environment is a self-managed `Proxmox` server. All infrastructure (DBs, Kafka, Services) runs as containers, easily migratable to a Cloud VPS.

### 5.2 Service Inventory

| #  | Service                | Protocol                    | Database / Storage                    | Kafka Role        |
| :- | :--------------------- | :-------------------------- | :------------------------------------ | :---------------- |
| 1  | `Client App`           | HTTP `:3000`                | —                                     | —                 |
| 2  | `API Gateway`          | HTTP `:8081`                | —                                     | —                 |
| 3  | `Auth Service`         | HTTP `:8082` / gRPC `:50052`| PostgreSQL `auth_db`, Ory Kratos/Hydra | Producer/Consumer |
| 4  | `Farm Service`         | HTTP `:8083` / gRPC `:50053`| PostgreSQL `farm_db`                  | Producer          |
| 5  | `Retail Service`       | HTTP `:8084` / gRPC `:50054`| PostgreSQL `retail_db`                | Producer          |
| 6  | `Logistics Service`    | HTTP `:8085` / gRPC `:50055`| PostgreSQL `logistics_db`, Valkey      | Producer/Consumer |
| 7  | `Payment Service`      | HTTP `:8086` / gRPC `:50056`| PostgreSQL `payment_db`               | Producer/Consumer |
| 8  | `Trace Service`        | HTTP `:8087` / gRPC `:50057`| PostgreSQL `trace_db`, Elasticsearch   | Consumer          |
| 9  | `Audit Service`        | HTTP `:8088` / gRPC `:50058`| PostgreSQL `audit_db`, Cassandra       | Consumer          |
| 10 | `Warehouse Service`    | HTTP `:8089` / gRPC `:50059`| PostgreSQL `warehouse_db`             | Producer/Consumer |
| 11 | `Kratos Public`        | HTTP `:4433`                | Ory Kratos                            | —                 |
| 12 | `Hydra Public`         | HTTP `:4444`                | Ory Hydra                             | —                 |
| 13 | `Kafka UI`             | HTTP `:8090`                | Kafka                                 | —                 |
| 14 | `SigNoz UI`            | HTTP `:3301`                | ClickHouse                            | —                 |
| 15 | `Kibana`               | HTTP `:5601`                | Elasticsearch                         | —                 |

#### `Webhook Service` (Ingress Gateway) — Details

The `Webhook Service` is a specialized gateway, **completely separate** from the `API Gateway`, responsible for receiving all external data streams:

| Source           | Endpoint                         | Processing                                               |
| :--------------- | :------------------------------- | :------------------------------------------------------- |
| `Stripe`         | `POST /webhooks/stripe`          | `HMAC-SHA256` verify → `Inbox` dedup → publish to Kafka  |
| `VNPay`          | `POST /webhooks/vnpay`           | `HMAC-SHA512` verify → `Inbox` dedup → publish to Kafka  |
| `IoT Device`     | `POST /webhooks/iot/gps`         | Token auth → normalization → publish `logistics.gps.updated` |

**Principle:** The `Webhook Service` **does not contain business logic**. Its sole task is: signature verification → dedup using `Inbox` → payload normalization → publish to `Kafka`. Business logic is handled in the corresponding service.

### 5.3 Kafka Topic Design (Draft)

| Topic                              | Producer(s)      | Consumer(s)                         |
| :--------------------------------- | :--------------- | :---------------------------------- |
| `auth.user.events`                 | Auth             | Farm, Audit, Gateway (Cache)        |
| `farm.harvest.created`             | Farm             | Processing, Trace, Audit            |
| `process.batch.completed`          | Processing       | Warehouse, Trace, Audit             |
| `retail.order.created`             | Retail           | Payment, Trace, Audit               |
| `payment.intent.created`           | Payment          | Retail, Trace, Audit                |
| `payment.completed`                | Payment webhook  | Warehouse, Retail, Trace, Audit     |
| `payment.simulated_completed`      | Payment          | Legacy/demo compatibility           |
| `payment.failed`                   | Payment          | Retail, Trace, Audit                |
| `payment.refunded`                 | Payment          | Retail, Trace, Audit                |
| `warehouse.stock.reserved`         | Warehouse        | Logistics, Payment, Trace, Audit    |
| `warehouse.stock.reservation_failed` | Warehouse      | Payment, Retail, Trace, Audit       |
| `warehouse.stock.released`         | Warehouse        | Retail, Trace, Audit                |
| `logistics.delivery.assigned`      | Logistics        | Warehouse, Retail, Trace, Audit     |
| `logistics.shipment.failed`        | Logistics        | Warehouse, Payment, Retail, Trace, Audit |
| `logistics.delivery.completed`     | Logistics        | Warehouse, Retail, Trace, Audit     |
| `logistics.gps.updated`            | Logistics        | Trace                               |
| `retail.order.completed`           | Retail           | Trace, Audit                        |

### 5.4 HA-Ready Design Principles

Although the demo is deployed on a single-node `Docker Compose`, the architecture is pre-designed to scale:

| Component           | HA Strategy                                                      |
| :------------------ | :--------------------------------------------------------------- |
| `Go Services`       | Stateless — horizontal scale by adding container replicas         |
| `PostgreSQL`        | Each service has its own DB (database-per-service) → independent scaling |
| `Kafka`             | Multi-partition topics, consumer groups for parallel processing   |
| `Valkey`            | Supports Cluster mode (Sentinel/Cluster) for GPS data            |
| `Elasticsearch`     | Shard/Replica strategy for read-model                            |
| `API Gateway`       | Stateless, can be placed behind a Load Balancer                  |

### 5.5 Technical Mechanisms

#### A. Decentralized Authorization

The `API Gateway` only handles **Authentication**: validating the `JWT` and extracting `Roles` from the token, then passing them down to services via `HTTP Headers` (`X-User-ID`, `X-User-Roles`).

Each Microservice integrates `Casbin` into its internal `Middleware` layer. Based on its own `policy.csv` file, the service **self-authorizes access** to each API endpoint — independent of the Gateway, increasing autonomy and reducing central load.

```
[Client] ──► [API Gateway] ──► JWT validate + extract Roles ──► Header: X-Roles=STORE_MGR
                                                                       │
                                                               [Retail Service]
                                                               Casbin Middleware
                                                               policy.csv: STORE_MGR CAN POST /orders
                                                                       │
                                                               ✅ Allow / ❌ Deny
```

#### B. Dual Idempotency

The system protects against duplication at **two independent levels**:

| Level | Request Type | Mechanism | Storage | TTL |
| :--- | :----------- | :----- | :------ | :-- |
| **Level 1** (Synchronous) | HTTP REST API | Header `Idempotency-Key` stored in `Valkey` | `Valkey` | 24h |
| **Level 2** (Asynchronous) | Kafka Consumer + Webhook | `Inbox Pattern` writing `event_id` to `PostgreSQL` | `PostgreSQL` | Permanent |

- **Level 1:** Client sends `Idempotency-Key: <uuid>` in the header. `API Gateway` middleware checks `Valkey`. If the key exists → return cached response, no re-processing.
- **Level 2:** Consumer (Kafka/Webhook) extracts `event_id`, opens a `Transaction`: checks the `inbox_events` table → if present → rollback and skip → if absent → save and process.

#### C. Configuration Management

Uses a combination of `.env` files and the `viper` library (Go):
- Each service has its own configuration file (`config.yaml` + `.env` override).
- `viper` automatically reads environment variables and supports hot-reloading.
- `Docker Compose` passes `Secrets` and `DB Host` via the `environment` block.
- No hardcoded credentials in the source code.

#### D. Distributed Tracing

`OpenTelemetry` + `SigNoz/ClickHouse`. Workflow:

1. `API Gateway` assigns a unique `Trace-ID` to each incoming request.
2. The `Trace-ID` is injected into:
   - `gRPC Context` (metadata) when calling internal services.
   - `HTTP Headers` (`traceparent`) when calling external services.
   - `Kafka Message Headers` when publishing events.
3. Each service creates a child `Span`, linked to the root `Trace-ID`.
4. `SigNoz UI` renders a waterfall/timeline displaying the entire transaction across multiple services.

> **Showcase value:** A `Trace-ID` from the store's order placement to Saga completion (through Payment → Warehouse → Logistics) is visualized intuitively on SigNoz — very impressive for TA reviewers.

---

## 6. Phased Delivery Plan

### Phase 1: Foundation
> **Goal:** Build the Monorepo skeleton, CI pipeline, and the first 3 core services.

| Deliverable                                 | Pattern showcase                     |
| :------------------------------------------ | :----------------------------------- |
| Project scaffold (Monorepo + shared libs)   | `DDD` project structure              |
| `API Gateway` (Gin + JWT parse + routing)   | `API Gateway Pattern`, `Rate Limiting` |
| `Identity Service` (Ory Kratos integration) | `OAuth2`/`OIDC`                      |
| `Farm Service` (Farm CRUD + harvest)        | `Transactional Design (Outbox)`      |
| Kafka + PostgreSQL infra (Docker Compose)   | `Event-Driven` base                  |
| Casbin middleware (shared lib)              | `Decentralized Authorization`        |
| gRPC proto definitions (shared)            | `Protocol Buffers`                   |

### Phase 2: Supply Chain + Payment
> **Goal:** Complete the Farm → Process → Warehouse pipeline, integrate Stripe, and full Saga flow.

| Deliverable                                    | Pattern showcase                         |
| :--------------------------------------------- | :--------------------------------------- |
| `Processing Service` (roasting, Batch ID)      | `Event-Driven Consumer/Producer`         |
| `Warehouse Service` (inventory, reserve/release)| `Saga Participant`                       |
| `Retail Service` (supply order placement)      | `Saga Orchestrator`                      |
| `Payment Service` (Stripe PaymentIntent + Webhook) | `Webhook HMAC`, `Inbox Pattern`      |
| Payment Gateway Abstraction (interface-based)   | `Strategy Pattern` (Stripe, VNPay, etc.) |
| Stripe Webhook + HMAC validation middleware     | `Zero Trust` (external data)             |
| Full Saga flow + Auto-Refund Compensation      | `Saga Pattern` (Choreography)           |
| `Inbox Pattern` (Idempotency on Consumer + Webhook) | `Transactional Inbox`             |

### Phase 3: Logistics & Real-time
> **Goal:** GPS tracking, transportation, and real-time data.

| Deliverable                                   | Pattern showcase                   |
| :-------------------------------------------- | :--------------------------------- |
| `Logistics Service` (dispatch, GPS tracking)  | `Geo-spatial` (Valkey GEO commands)|
| GPS simulator (fake DRIVER coordinates)        | Real-time data pipeline            |
| Valkey integration (cache + distributed lock)  | `Distributed Lock`, `Cache-aside`  |
| mTLS for gRPC between services                | `Zero Trust Architecture`          |

### Phase 4: Observability & Traceability
> **Goal:** CQRS read-model, audit trail, and distributed tracing.

| Deliverable                                    | Pattern showcase                  |
| :--------------------------------------------- | :-------------------------------- |
| `Traceability Service` (Kafka → Elasticsearch) | `CQRS` (Read-model projection)    |
| QR code generation + traceability lookup       | `CQRS Query` endpoint             |
| `Audit Service` (Kafka → Cassandra, hash chain)| `Event Sourcing` (Lite)           |
| `Hash Chaining` for data integrity             | `Data Integrity Pattern`          |
| OpenTelemetry + SigNoz integration             | `Distributed Tracing`             |
| Prometheus + Grafana dashboards                | `Observability Stack`             |

### Phase 5: Control Plane Visualization Dashboard (Frontend)
> **Goal:** ReactJS Dashboard combining Operational UI + System Visualization.

| Deliverable                                     | Tech showcase                    |
| :---------------------------------------------- | :------------------------------- |
| `Monitor Service` (Kafka → SSE/WebSocket)       | Real-time event broadcasting     |
| ReactJS app (Vite + React Flow + Framer Motion) | Modern frontend stack            |
| Split-screen: App View + System View            | `Service Mesh Visualization`     |
| Isometric 3D service map with animated edges    | `React Flow` animated edges      |
| Visual metaphors (Green Bean, Roasted Bean, etc.)| Domain-specific UI language      |
| Chaos Control panel (kill Kafka, set stock = 0) | `Resiliency` demonstration       |
| CQRS Time Machine (time slider retrieval)       | CQRS visual demo                 |

---

## 7. Non-Functional Requirements

| Category         | Requirement                                                         |
| :--------------- | :------------------------------------------------------------------ |
| **Performance**  | CQRS read query (traceability lookup) < 100ms                      |
| **Performance**  | GPS update latency < 500ms (from Logistics → Monitor → Dashboard)  |
| **Availability** | HA-ready design: Stateless services, database-per-service           |
| **Scalability**  | Kafka multi-partition, consumer group ready                          |
| **Security**     | JWT Authentication (Ory Kratos), Casbin Authorization per-service   |
| **Security**     | mTLS for all gRPC internal communication                            |
| **Security**     | HMAC signature validation for Stripe Webhook (anti-spoofing)        |
| **Security**     | Stripe Test Mode only — no real money processing in the demo        |
| **Idempotency**  | `Inbox Pattern` for Stripe Webhook (at-least-once → exactly-once)   |
| **Integrity**    | Hash chaining on Audit log (Cassandra) anti-tampering               |
| **Observability**| Distributed tracing (OpenTelemetry + SigNoz) on every request       |
| **Observability**| Prometheus metrics + Grafana dashboard for each service             |
| **Deployment**   | Full Docker Compose for local dev                                   |
| **Deployment**   | Docker images ready for Cloud VPS deployment                        |

---

## 8. Technical Decisions

| Decision                     | Choice                      | Reason                                                             |
| :--------------------------- | :-------------------------- | :----------------------------------------------------------------- |
| Backend language             | `Go` 1.22+                  | Performance, concurrency, ecosystem for microservices              |
| Code organization            | `Monorepo`                  | Share `proto`, `pkg` libs; easy cross-service change management    |
| Service architecture         | `Clean Architecture` v4     | Separate Domain/UseCase/Infra — testable, replaceable adapters      |
| HTTP framework               | `Gin`                       | Lightweight, high-performance REST                                 |
| Internal RPC                 | `gRPC` + `Protobuf`         | Type-safe, high-speed, contracts defined in `api/`                 |
| Message broker               | `Apache Kafka`              | Industry standard for event-driven architecture                    |
| Relational DB                | `PostgreSQL` v15            | ACID, mature, database-per-service                                 |
| Document DB                  | `Apache Cassandra`          | Wide-column store, flexible schema for audit logs                  |
| Search engine                | `Elasticsearch`             | Full-text search + CQRS read-model                                 |
| Cache / Real-time            | `Valkey`                    | Valkey alternative, GEO commands for GPS, Idempotency-Key store    |
| Payment gateway              | `Stripe` (Test Mode)        | Industry standard, excellent API docs, Webhook support             |
| Payment abstraction          | `Strategy + Factory Pattern`| `PaymentProvider` interface + `ProviderFactory` → easy to add VNPay |
| Identity                     | `Ory Kratos`                | Open-source, self-hosted identity management                       |
| Authorization                | `Casbin` + `policy.csv`     | Embeddable RBAC, decentralized per-service, Gateway only authn    |
| Config management            | `viper` + `.env`            | Flexible, hot-reload, no hardcoded secrets, Docker-friendly        |
| Dependency Injection         | Manual DI                   | Avoid magic, easy to debug, Composition Root at main.go            |
| Reliability                  | Selective Outbox            | Applied to critical flows to ensure consistency                    |
| Frontend                     | `ReactJS` (Vite)            | Combined Operational UI + Control Plane Visualization Dashboard    |
| Visualization                | `React Flow`                | Node-based UI for service mesh visualization                       |
| Animation                    | `Framer Motion`             | Micro-animations, glow effects                                     |
| State management             | `Zustand`                   | Lightweight state for real-time WebSocket data                     |
| Tracing                      | `OpenTelemetry` + `SigNoz`  | Trace-ID from Gateway, propagate via gRPC/Kafka/HTTP headers       |
| Metrics                      | `Prometheus` + `Grafana`    | Industry standard monitoring stack                                 |
| Host infrastructure          | `Proxmox`                   | Self-hosted hypervisor, VM-based Docker environment                |
| Deployment                   | `Docker` + `Docker Compose` | Containerized, Cloud VPS ready                                     |

---

## 9. Risks & Mitigations

| Risk                                          | Impact | Mitigation                                              |
| :--------------------------------------------- | :----- | :------------------------------------------------------ |
| Solo developer → scope creep                   | High   | Strict phased delivery, MVP-first mindset                |
| Kafka learning curve                           | Medium | Start with single-partition, scale later                 |
| Too many databases to manage                   | Medium | Docker Compose manages entire infrastructure             |
| Visualization Dashboard too complex            | High   | Phase 5 — only build after backend is stable             |
| HA design cannot be verified on single-node    | Low    | Document HA strategy, verify via architecture review     |
| Stripe Webhook replay/spoofing                 | High   | HMAC validation + Inbox idempotency pattern              |
| Money "hanging" when Saga fails post-payment   | High   | Auto-refund compensation, strict Payment state machine    |
| Stripe API version changes                     | Low    | Abstract via interface, easy to swap provider            |

---

## 10. Success Criteria

The project is considered a **success** when:

- [ ] All 9+ services run synchronously on `Docker Compose` without errors.
- [ ] `Saga Rollback` scenario demonstrated (out of stock → auto-refund Stripe → order reversed).
- [ ] `Outbox Pattern` demonstrated (kill Kafka → recover → events flush successfully).
- [ ] `Inbox Pattern` demonstrated (duplicate Webhook sent → processed only once).
- [ ] Successful Stripe payment (Test Mode) and refund when Saga fails.
- [ ] QR scan returns full Farm-to-Cup history of a retail-sold cup/order item.
- [ ] Control Plane Visualization Dashboard displays real-time data flow between services.
- [ ] Reviewers/TAs can read the code and understand each `Pattern` applied.
- [ ] System can be deployed to a Cloud VPS using `docker-compose up -d`.

---

## Appendix A: Glossary

| Term                   | Definition                                                                          |
| :--------------------- | :---------------------------------------------------------------------------------- |
| `Batch ID`             | Unique identifier for a finished coffee batch: `RR-{REGION}-{YEAR}-{MMDD}-{LOT}`    |
| `Outbox Pattern`       | Writing events to DB in the same transaction as business data, worker pushes later  |
| `Inbox Pattern`        | Storing processed events to avoid duplicate processing (idempotency)               |
| `Saga`                 | Sequence of distributed transactions with compensating actions for rollback          |
| `CQRS`                 | Separation of write (Command) and read (Query) into two different systems            |
| `Clean Architecture`   | 3-layer architecture (Domain/UseCase/Infra) — framework-agnostic, testable           |
| `Polyglot Persistence` | Using multiple types of databases, each suited to specific data characteristics      |
| `Hash Chaining`        | Each audit record contains the hash of the previous record, creating an immutable chain |
| `Control Plane Visualization` | Real-time dashboard monitoring the entire system architecture               |
| `PaymentIntent`        | Stripe object representing a payment transaction waiting for processing              |
| `Webhook Service`      | Specialized Ingress Gateway receiving and validating external data (Stripe, IoT)     |
| `HMAC`                 | Hash-based Message Authentication Code — digital signature for integrity            |
| `Dual Idempotency`     | 2-level protection: Valkey for HTTP (sync) + Inbox Pattern for Kafka/Webhook (async) |
| `Idempotency-Key`      | UUID created by the client, sent in HTTP header to prevent duplicate execution      |
| `ProviderFactory`      | Factory creating Payment adapters (Stripe/VNPay) based on config                    |
| `Compensating Action`  | Rollback action (e.g., Refund) when a Saga step fails                               |
| `Proxmox`              | Open-source self-hosted hypervisor platform, running Docker VM                      |
| `Trace-ID`             | Unique ID generated at the Gateway, propagated throughout the system for debugging   |

---

## Appendix B: Reference Documents

| Document                                                                  | Location                                        |
| :------------------------------------------------------------------------ | :---------------------------------------------- |
| System Architecture                                                       | `docs/product/TECH.md`                           |
| Data Models & Persistence                                                 | `docs/product/TECH.md`                           |
| UI/UX Visual Ideas (Visualization Dashboard)                              | `docs/product/ui-ux/DESIGN.md`                   |
| REST API Specifications                                                   | `docs/product/TECH.md`                           |
| gRPC Contract Definitions                                                 | `docs/product/TECH.md`                           |
| Demo Setup Runbook                                                        | `docs/product/DEMO_SETUP_RUNBOOK.md`             |
