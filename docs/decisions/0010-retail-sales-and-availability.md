---
artifact_type: adr
id: ADR-0010
status: accepted
owner: shared
decision_approval: approved_in_chat_2026-06-06
trace:
  requirement: REQ-QR-001
  phase: PHASE-2
  tickets_or_bugs:
    - RR-URG-07A
  detail_design: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/technical_design.md
  master_docs:
    - docs/architecture/SDD/README.md
    - docs/architecture/SDD/MASTER_DATA.md
    - docs/architecture/ERD.md
  release_notes: docs/releases/CHANGELOG.md
---

# ADR 0010: Retail Sales And Materialized Availability

## Status

Accepted on 2026-06-06.

## Context

RR-URG-07 needs a chain-wide menu, store inventory lineage, one QR identity per
sold cup, and fast store-menu availability. Retail already owns its database,
so table names do not need a `retail_` prefix. Running `SUM(FLOOR(...))`
over all lots on every public menu request would add repeated read cost.

The project is in development and permits rewriting initial migrations after
explicit approval.

## Decision

- Use `menus`, `menu_items`, `inventory_lots`, `sales`,
  `sale_items`, `stock_movements`, and `store_menu_inventories`.
- Keep `orders` for warehouse supply orders; POS purchases use `sales`.
- Treat each priced drink-size in SAMPLE_MENU.md as one sellable menu item.
- Treat each sale item as one sold cup, `product_id`, and QR.
- Keep stock movements as ledger authority and lot available quantity as cache.
- Materialize store availability and rebuild it on stock mutation/reconciliation.
- Apply floor per lot so one cup cannot combine provenance across lots.
- Rewrite Retail `000001_init.up.sql` during development.

## Alternatives Considered

- Aggregate every GET: simpler writes, repeated public-query cost.
- Valkey cache: fast reads, but invalidation and dual-write failure modes.
- Store availability on menu items: cannot represent individual stores.
- Reuse orders for POS: conflates supply fulfillment and sold-cup sales.
- Prefix tables with `retail_`: redundant inside the Retail database.

## Consequences

- Positive: bounded indexed reads and transactional consistency.
- Positive: explicit catalog, supply, sale, lineage, and sold-cup ownership.
- Positive: availability can be rebuilt from the ledger.
- Negative: stock mutations perform additional recomputation writes.
- Negative: recipe or shared-SKU changes require affected rows to reconcile.
- Neutral: production will need forward-only migrations before deployed data
  must be preserved.
