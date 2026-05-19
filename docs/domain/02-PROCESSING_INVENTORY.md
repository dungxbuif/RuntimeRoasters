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
