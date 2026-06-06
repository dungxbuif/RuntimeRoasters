---
artifact_type: ticket
id: RR-URG-07C
status: ready
owner: shared
priority: high
lane: high-risk
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  detail_design: required
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07C: Trace-Service Public Sold-Cup Query

## Goal

Implement public trace lookup for a sold cup/item by UI-issued `product_id`.

This ticket depends on [RR-URG-07A](../RR-URG-07A/ticket.md) and [RR-URG-07B](../RR-URG-07B/ticket.md).

## Scope

- Add public lookup by `product_id`/`trace_code`.
- Build hybrid trace response by combining retail sold-cup data with upstream trace documents.
- Keep public response sanitized.
- Add public list of traceable sold items for demo.

## APIs

Public/no-auth:

- `GET /v1/public/trace-items`
- `GET /v1/public/trace/:trace_code`

`trace_code` is the public API name. It equals the UI-issued `product_id`.

Internal dependency:

- trace-service calls retail-service internal API/gRPC to resolve `product_id` into:
  - sale item
  - store display
  - product display
  - sold timestamp
  - source batch/harvest/warehouse IDs

## Hybrid Query Flow

1. Public request calls trace-service with `trace_code`.
2. Trace-service treats `trace_code` as `product_id`.
3. Trace-service fetches sold-cup details from retail-service.
4. Trace-service queries Elasticsearch/Postgres trace documents by source batch/harvest IDs.
5. Trace-service combines retail sale details with upstream Farm -> Warehouse journey.
6. Trace-service returns a sanitized public document.

## Public Document Shape

- `product_id`
- `trace_code`
- `sale_item`
- `invoice_no`
- `product`
- `store`
- `farm`
- `warehouse`
- `origin_batch_id`
- `harvest_id`
- `inventory_lot_id`
- `milestones`
- `service_participation`
- `trace_ids`
- `public_safe`

## Sanitization

Do not expose:

- internal user IDs
- JWT/session data
- payment provider data
- raw event payloads
- audit-only details
- internal service secrets

## Acceptance Criteria

- `GET /v1/public/trace/:product_id` resolves a sold cup without login.
- Missing or unknown `product_id` returns 404.
- Trace response includes retail sale and upstream origin.
- Trace response is safe for public display.
- Public list returns only sold cups/items, not unsold product catalog rows.

## Verification

- Unit: sanitizer removes private fields.
- Unit: `trace_code` maps to `product_id`.
- Integration: seeded `product_id` resolves to public trace document.
- Integration: unknown `product_id` returns 404.
- Integration: public list excludes unsold products.
