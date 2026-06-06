---
artifact_type: ticket
id: RR-URG-07
status: ready
owner: shared
priority: high
lane: high-risk
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  detail_design: split_across_RR-URG-07A_to_RR-URG-07E
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07: Public Sold-Cup QR Trace Umbrella

## Goal

Deliver a simple public trace demo:

1. User opens a no-auth coffee sale UI.
2. User selects a retail store and product from the chain-wide menu.
3. UI creates one `product_id` per sold cup/item and submits the sale.
4. `retail-service` creates an invoice, persists sale items, decrements store inventory, and stores the UI-issued `product_id`.
5. `trace-service` resolves `product_id` to the sold item and derives origin from `product_id` plus the retail inventory lot lineage.
6. QR opens a public trace page for that sold cup/item.

No payment, payment simulation, or delivery workflow is required inside the cup purchase flow.

## Split Tickets

| Ticket | Scope | Depends On |
| :--- | :--- | :--- |
| [RR-URG-07A](./subtickets/RR-URG-07A/ticket.md) | Menus, menu items, store inventory lots, sales, sold cups, movements, availability, seed data | RR-URG-05 |
| [RR-URG-07B](./subtickets/RR-URG-07B/ticket.md) | Public demo sale APIs, manager sold-items APIs, event contract | RR-URG-07A |
| [RR-URG-07C](./subtickets/RR-URG-07C/ticket.md) | Trace-service hybrid public lookup by `product_id`/`trace_code` | RR-URG-07A, RR-URG-07B |
| [RR-URG-07D](./subtickets/RR-URG-07D/ticket.md) | Public POS UI, invoice modal, QR rendering, sold-items UI | RR-URG-07B, RR-URG-07C |
| [RR-URG-07E](./subtickets/RR-URG-07E/ticket.md) | Unit/integration/E2E/platform evidence for full flow | RR-URG-07A to RR-URG-07D |

## Non-Negotiable Decisions

- QR token is UI-issued `product_id`.
- Backend validates and persists `product_id`; backend does not invent origin from the QR token.
- Backend derives source lineage from `product_id` plus selected `inventory_lot`.
- `menus` and `menu_items` own the chain-wide catalog; contracts use `menu_item_id`.
- Store inventory is separate from warehouse inventory.
- Each sold cup/item maps to exactly one `product_id`.
- Public trace data is sanitized and no-auth.

## Acceptance Criteria

- A visitor can buy coffee from a selected store without login or payment.
- The sale creates an invoice and one sale item per sold cup/product.
- Each sale item has a UI-issued `product_id` and QR URL.
- Store inventory decreases and stock movement history is recorded.
- Manager UI can show sold cups/items for a selected store.
- Public QR page resolves a `product_id` to a real sold-cup trace.
- Seed data includes several already-sold cups/items for immediate demo scanning.
