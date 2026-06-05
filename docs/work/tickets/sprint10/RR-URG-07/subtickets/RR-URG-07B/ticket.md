---
artifact_type: ticket
id: RR-URG-07B
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

# RR-URG-07B: Retail Demo Sale APIs

## Goal

Implement the retail-service APIs that let a user buy coffee without payment and let managers view sold cups/items.

This ticket depends on [RR-URG-07A](../RR-URG-07A/ticket.md).

## Scope

- Public product/menu read APIs.
- Public demo sale API.
- Manager sold-items APIs.
- Retail sale transaction logic.
- `retail.sale.completed` event contract and outbox publish.

## APIs

Public/no-auth:

- `GET /v1/public/retail/menu-items`
- `GET /v1/public/retail/stores/:store_id/menu-items`
- `POST /v1/public/retail/sales`

Authenticated manager/admin:

- `GET /v1/retail/sales`
- `GET /v1/retail/sales/:sale_id`
- `GET /v1/retail/stores/:store_id/sold-items`

Internal support for trace-service:

- `GET /v1/internal/retail/sold-items/:product_id` or equivalent gRPC method.

## Sale Request

`POST /v1/public/retail/sales` accepts:

- `store_id`
- `items[]`
  - `menu_item_id`
  - `quantity`
  - `product_id`

The UI creates `product_id`; backend must not generate it.

## Sale Transaction

In one DB transaction:

1. Validate store and menu items.
2. Validate each `product_id` is non-empty and unique.
3. Resolve `menu_item_id` to SKU and menu metadata.
4. Select one retail inventory lot for the store/SKU.
   - Demo selection may be random.
   - Production replacement is FIFO/FEFO.
5. Create `retail_sales`.
6. Create `retail_sale_items`.
7. Set `trace_code = product_id`.
8. Decrement `retail_inventory_lots.available_quantity`.
9. Insert `retail_stock_movements`.
10. Insert outbox event `retail.sale.completed`.

## Event Contract

Topic: `retail.sale.completed`

Payload:

- `event_id`
- `sale_id`
- `invoice_no`
- `store_id`
- `items[]`
  - `sale_item_id`
  - `menu_item_id`
  - `product_id`
  - `sku`
  - `product_name`
  - `quantity`
  - `retail_inventory_lot_id`
  - `trace_code`
  - `public_url`
  - `source_batch_id`
  - `source_harvest_id`
  - `source_warehouse_id`
- `total_amount`
- `sold_at`
- `occurred_at`

## Acceptance Criteria

- User can buy coffee without login and without payment.
- API returns invoice and QR URL per sold item.
- Duplicate `product_id` is rejected.
- Store inventory decreases after purchase.
- Stock movement history is persisted.
- Manager can list sold cups/items for assigned store.

## Verification

- Unit: sale requires store, items, `menu_item_id`, and sold-item `product_id`.
- Unit: duplicate `product_id` fails.
- Unit: inventory decrement and stock movement happen in one transaction.
- Integration: `POST /v1/public/retail/sales` creates invoice and sale items.
- Integration: outbox contains `retail.sale.completed`.
- Auth integration: manager sold-items endpoint is store-scoped.
