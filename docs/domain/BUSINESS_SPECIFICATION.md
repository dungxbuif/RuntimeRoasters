# ☕ Runtime Roasters — Overall System Business Specification

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
| **Vehicle / Driver** | Demo fleet resources assigned by Warehouse/Logistics. Driver accounts execute route simulation and confirm shipment milestones. |
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

Information access rights are strictly established based on actual job responsibilities:

1.  **System Administrator (`ADMIN`)**: Creates manager accounts, farms, warehouses, stores, vehicles, and driver records; assigns managers; views aggregate health. It is not the normal operator for harvest, dispatch, receipt, GPS, or processing actions.
2.  **Farm Manager (`FARM_MANAGER`)**: Creates harvests and watches pickup status for assigned farms.
3.  **Warehouse Manager (`WAREHOUSE_MGR`)**: Receives pickup notifications, dispatches vehicles/drivers, receives returned shipments, creates intakes, manages processing and inventory.
4.  **Store Manager (`STORE_MGR`)**: Creates paid retail orders for assigned stores and watches incoming delivery status.
5.  **Driver (`DRIVER`)**: Views only assigned shipments, starts route simulation, posts GPS/location updates, and confirms pickup/delivery/return milestones.
6.  **Processor (`PROCESSOR`, optional split)**: Runs or assists processing steps if the final role split separates production from warehouse management.

---

## 5. Security Logs & Traceability Reporting

All data in the value chain is aggregated into a "Traceability Report," including:
- **Cultivation Profile**: Seed variety, harvest time, and growing region characteristics.
- **Processing Profile**: Roasting process, loss metrics, and packaging time.
- **Transportation Profile**: Movement history, transit points, and personnel in charge.
- **Realtime Profile**: Socket/SSE milestone updates shown during the demo.
- **Service Profile**: OTel-derived service/component participation for architecture visualization.
- **Public QR Profile**: Prepared product trace documents shown to unauthenticated visitors through QR trace URLs.

The system utilizes secure storage techniques to ensure that once these records are confirmed, they cannot be altered, creating an authentic and permanent proof of product quality.

---
*Runtime Roasters Business Specification — Version 1.0 (2026)*
