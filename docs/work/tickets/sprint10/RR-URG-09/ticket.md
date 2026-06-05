# RR-URG-09: Client Dashboards, Driver Simulation UI, And Public QR Trace Page

## Priority

P3. Depends on backend APIs enough to integrate; UI can be scaffolded earlier with strict API contracts.

## Problem

The app currently has visual dashboards but missing connected operational flows. Final demo needs role-specific UI actions, Driver Client simulation, realtime updates, public ArchitectureTopology, and public QR trace page.

## Scope

- Farm dashboard pickup status.
- Warehouse pickup/outbound dispatch UI.
- Driver Client simulation UI.
- Logistics live map from backend state.
- Retail paid order UI.
- Trace/public QR UI.
- Realtime stream integration.

## Implementation Details

### 1. Farm Dashboard

Required UI:

- assigned farm list.
- create harvest action.
- harvest/pickup status list.
- active pickup progress.

Checklist:

- [ ] Hide create harvest when no farm assignment.
- [ ] Show pickup request after harvest.
- [ ] Show driver assigned/status updates.

### 2. Warehouse Dashboard

Required UI:

- inbound pickup notification queue.
- pickup dispatch panel.
- active inbound shipments.
- receipt/intake panel.
- processing and inventory panels.
- outbound paid-order dispatch queue.

Checklist:

- [ ] Dispatch button visible only to `WAREHOUSE_MGR`.
- [ ] Receipt action scoped to assigned warehouse.
- [ ] Processing flow starts after intake only.
- [ ] Outbound queue appears after stock reservation.

### 3. Driver Client

Required UI:

- assigned shipment list/detail.
- route map.
- Start simulation button.
- status timeline.
- milestone buttons:
  - depart
  - arrived at farm/store
  - confirm pickup/loading or delivery
  - start return
  - confirm returned to base

Simulation behavior:

- Load seeded route points.
- Animate 30-45 seconds per leg.
- Post GPS/status updates to backend at each tick.
- Disable controls when shipment not assigned to logged-in driver.

Checklist:

- [ ] Requires `DRIVER` login.
- [ ] No shipment means no Start button.
- [ ] Posts accepted GPS updates.
- [ ] Mandatory return leg visible and actionable.

### 4. Logistics Dashboard

Required UI:

- real shipments from backend.
- no generic production mock fallback.
- vehicle/driver roster.
- live GPS markers from socket/SSE/backend polling.
- shipment detail drawer.

Checklist:

- [ ] Shows farm pickup and retail delivery.
- [ ] Shows return-to-base leg.
- [ ] Shows driver/vehicle status.

### 5. Retail Dashboard

Required UI:

- assigned store selector.
- create paid order/checkout action.
- payment/simulated payment status.
- reservation status.
- incoming delivery status.
- completed state.

Checklist:

- [ ] `STORE_MGR` scoped by `store_ids`.
- [ ] Paid order wording used.
- [ ] Replenishment wording not used in final demo.

### 6. Public ArchitectureTopology

Requirements:

- no auth required.
- sanitized service/topology/demo data only.
- can show live public-safe events.

Checklist:

- [ ] No private operational action data.
- [ ] No payment/user secrets.

### 7. Public QR Trace Page

Required UI:

- public button to show/generate demo QR codes.
- QR list for sold cups/items.
- click/scan opens public trace URL.
- trace page displays real trace-service/Elasticsearch document.
- optional service/component participation from OTel.

Checklist:

- [ ] Public QR list loads without login.
- [ ] QR opens real trace.
- [ ] Trace journey includes driver return.
- [ ] Missing OTel data does not break page.

## Acceptance Criteria

1. A reviewer can run farm pickup flow across Farm, Warehouse, Driver, Logistics, and Trace UI.
2. A reviewer can run paid order delivery flow across Retail, Warehouse, Driver, Logistics, and Trace UI.
3. Driver Client is the intentional simulator and posts to backend.
4. Public ArchitectureTopology works without login.
5. Public QR trace works without login and shows real sold-cup data.

## Test Checklist

- [ ] Component tests for role-gated actions where practical.
- [ ] E2E: Farm Manager create harvest.
- [ ] E2E: Warehouse Manager dispatch pickup.
- [ ] E2E: Driver route simulation and return.
- [ ] E2E: Store Manager create paid order.
- [ ] E2E: Public QR trace page.
- [ ] Visual check: two-screen realtime demo pacing.
