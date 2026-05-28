# Demo UI Audit

## API-Backed UI

| UI Surface | Current Classification |
| --- | --- |
| `/dashboard/logistics` shipments/locations | API-backed when logistics API returns data. |
| `/dashboard/driver` shipment actions | API-backed milestone calls. |
| `/dashboard/driver` GPS publish | Browser-simulated movement posted to real backend API. |
| `/dashboard/traceability` document lookup | API-backed by trace document lookup. |
| `/dashboard/finance` payment ledger | API-backed by payment service. |
| `/dashboard/finance` Pass/Fail | Demo control that sends signed webhook payloads to payment service. |
| `/dashboard/retail` order creation | API-backed where current retail service route is used. |
| `/dashboard/retail/orders` | Routing/view surface for retail orders. |
| `/dashboard/warehouse` queues | API-backed where APIs return data; static activity sections need labeling. |

## Browser-Simulated UI

- Driver route replay from `/data/routes.json`.
- GPS timer loop.
- Logistics dashboard simulation toggle that posts periodic GPS for active real shipments.

## UI Mock / Visualizer Sections

- Fallback shipments in `/dashboard/logistics` when logistics API returns empty.
- Static stats/cards/logs in some dashboard pages.
- Saga/topology visual sections when they are not connected to current trace/order data.
- Map route lines when displayed without active shipment context.

## Recommended Cleanup Backlog

- Label mock/fallback sections in UI when they remain.
- Remove logistics fallback shipments after seeding is rich enough.
- Add traceability shortcuts for latest seeded/live order, shipment, harvest, or batch IDs.
- Add light-mode logistics map polish.
- Render route lines only for active shipments.
- Add clear labels for browser-simulated GPS vs backend-persisted GPS.
- Replace static dashboard logs with API-backed event/feed data where feasible.

## Review Language

Use precise labels in demos:

- "API-backed"
- "browser-simulated"
- "backend-simulated"
- "mock fallback"
- "planned"

Avoid presenting fallback data as live production state.
