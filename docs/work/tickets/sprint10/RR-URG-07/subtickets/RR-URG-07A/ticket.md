---
artifact_type: ticket
id: RR-URG-07A
status: ready
owner: shared
priority: high
lane: high-risk
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  detail_design: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/technical_design.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07A: Retail Sale Schema And Seed Data

## Goal

Add the minimum retail data model needed for a simple coffee sale demo and public sold-cup QR trace.

This ticket does not implement UI or retail-service public lookup.

## Scope

- Add chain-wide retail menu/products.
- Add store-level retail inventory lots with upstream lineage.
- Add invoice/sale and sale item tables.
- Add stock movement history.
- Seed demo products, store inventory, and a few sold cups/items.

## Data Model

`retail_menu_items`:

- `id` (`menu_item_id` in application/API contracts)
- `product_code`
- `sku`
- `name`
- `price`
- `active`
- `created_at`
- `updated_at`

`retail_inventory_lots`:

- `id`
- `store_id`
- `sku`
- `source_batch_id`
- `source_harvest_id`
- `source_warehouse_id`
- `source_order_id` optional
- `received_quantity`
- `available_quantity`
- `received_at`
- `created_at`
- `updated_at`

`retail_sales`:

- `id`
- `store_id`
- `invoice_no`
- `status`
- `total_amount`
- `sold_at`
- `idempotency_key`
- `created_at`

`retail_sale_items`:

- `id`
- `sale_id`
- `store_id`
- `menu_item_id`
- `sku`
- `product_name`
- `retail_inventory_lot_id`
- `quantity`
- `unit_price`
- `product_id`
- `trace_code`
- `public_url`
- `sold_at`

`retail_stock_movements`:

- `id`
- `store_id`
- `sku`
- `retail_inventory_lot_id`
- `movement_type`: `RECEIVED`, `SOLD`, `ADJUSTMENT`
- `quantity_delta`
- `reference_type`: `SEED`, `SALE_ITEM`
- `reference_id`
- `occurred_at`

## Seed Requirements

- Seed at least 2 active chain-wide coffee products.
- Seed store inventory lots for Hoan Kiem and District 1 stores.
- Each inventory lot must have source lineage: batch, harvest, warehouse where known.
- Seed at least 3 completed sales/invoices.
- Each seeded sold item must have a deterministic UI-style `product_id`, for example `RR-CUP-HK-0001`.
- For seeded sold items, stock movements must already decrement inventory.

## Invariants

- `menu_item_id` references one chain-wide menu item.
- `product_id` is unique per sold cup/item.
- `trace_code` equals `product_id` for public trace compatibility.
- Each sale item references exactly one retail inventory lot.
- Inventory balance equals the sum of all stock movement deltas.
- Inventory-lot `received_quantity` equals the sum of its `RECEIVED` movements.
- Products are chain-wide; inventory availability is store-specific.

## Acceptance Criteria

- Retail DB migrations create all required tables and indexes.
- Seed can run idempotently.
- Seeded stores have products available for sale.
- Seeded sold cups/items exist and have QR-ready `product_id` values.
- No public QR item exists without source inventory lineage.

## Verification

- Unit: seed loader validates products and inventory lots.
- Integration: migration creates schema on empty DB.
- Integration: seed creates products, inventory lots, sales, sale items, and stock movements.
- SQL evidence: selected store has positive inventory and several sold items with unique `product_id`.
