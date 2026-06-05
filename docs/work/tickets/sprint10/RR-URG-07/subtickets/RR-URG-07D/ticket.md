---
artifact_type: ticket
id: RR-URG-07D
status: ready
owner: shared
priority: high
lane: normal
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  detail_design: required
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07D: Client Sale And QR UI

## Goal

Build the usable frontend for buying demo coffee, displaying invoices/QRs, listing sold cups by store, and opening public trace pages.

This ticket depends on [RR-URG-07B](../RR-URG-07B/ticket.md) and [RR-URG-07C](../RR-URG-07C/ticket.md).

## Scope

- Public POS/Kiosk sale UI.
- UI-generated `product_id` per sold item.
- Invoice modal/page after purchase.
- QR rendering from `product_id` public URL.
- Sold cups/items view for a selected store.
- Public trace page wired to real API.

## Public POS UI

Required controls:

- Store selector.
- Product grid from chain-wide menu.
- Quantity control.
- Buy button.
- Invoice result view.
- QR per sold cup/item.

UI must generate `product_id` before submitting the sale request.

Suggested format:

- `RR-CUP-{STORE_SHORT}-{YYYYMMDD}-{short_random}`

The exact format is UI-owned, but it must be unique enough for demo and backend must still enforce uniqueness.

## Sold Items UI

Required:

- Store selector.
- Table/list of sold cups/items for the selected store.
- Columns: sold time, invoice number, product, `product_id`, QR/open trace action.
- Empty state for stores with no sales.

## Public Trace Page

Replace static/mock trace data with:

- fetch `GET /v1/public/trace/:trace_code`.
- loading, not found, error, and success states.
- business milestones first.
- OTel/service participation only as secondary enrichment.

## Acceptance Criteria

- User can buy coffee from UI without login.
- UI sends `product_id` per sold item.
- Invoice view shows QR for each sold item.
- Sold-items UI shows cups sold by selected store.
- QR opens real public trace page.
- Public trace page no longer uses static mock journey.

## Verification

- Unit/component: `product_id` is created before sale request.
- Unit/component: invoice renders QR per item.
- E2E: buy coffee -> invoice -> QR -> public trace.
- E2E: sold-items view includes the newly sold cup.
- E2E must use fail-fast timeouts and must not leave Playwright hanging.
