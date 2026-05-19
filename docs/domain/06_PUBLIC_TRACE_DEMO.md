# Public QR Trace Demo

This document defines the public traceability demo shown to unauthenticated users.

## 1. Purpose

Runtime Roasters is also a learning/showcase project. The public Client App should let visitors inspect prepared coffee product journeys without logging in.

The public trace page must show real prepared trace data, not a purely mocked visual.

## 2. Public User Flow

1. Visitor opens the public Client App trace showcase page.
2. Visitor clicks a button such as "Generate Demo QR Codes".
3. Client displays a list of QR codes for pre-seeded products.
4. Visitor scans or clicks a QR.
5. Public trace page opens a trace document.
6. Trace document shows the prepared product journey:
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

The public page should use prepared seed data:

- Seeded farms, warehouses, stores, products, shipments, and orders.
- Seeded or replayed domain events consumed by trace-service.
- Trace-service projects the business journey into Elasticsearch.
- Optional OTel/service participation can be attached as technical detail.

Public trace should query trace-service or a public-safe API backed by Elasticsearch. It should not read Cassandra or internal service databases directly.

## 4. QR Content

QR codes should point to stable public trace URLs:

- `/trace/public/:trace_code`
- or `/qr/:trace_code`

The QR payload should not expose internal primary keys if avoidable. Use a public `trace_code` or product code mapped by trace-service.

## 5. Visibility Rules

Allowed on public trace page:

- farm name/region
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
