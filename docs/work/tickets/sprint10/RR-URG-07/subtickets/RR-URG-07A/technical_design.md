---
artifact_type: detail_design
id: DESIGN-RR-URG-07A
status: ready
owner: ai
approval: pending
human_fields:
  - approval
  - constraints
  - scope_decisions
ai_fields:
  - problem
  - context_loaded
  - brownfield_scope
  - proposed_approach
  - alternatives_considered
  - impacted_areas
  - test_plan
  - reconciliation_plan
shared_fields:
  - status
  - trace
  - small_task_exemption
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  ticket_or_bug: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/ticket.md
  test_verification: pending
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: pending
  adrs:
    - pending
  master_docs_touched:
    - docs/architecture/ERD.md
    - docs/requirements/MASTER_DATA.md
---

# RR-URG-07A Detail Design

## Status

- ID: `DESIGN-RR-URG-07A`
- Status: `ready`
- Ticket: `RR-URG-07A`
- Approval: `pending`
- Author: AI
- Updated: 2026-06-06

## Problem

Retail service currently persists only stores, supply orders, outbox, and inbox records. It cannot represent a chain-wide coffee menu, store inventory lineage, a completed consumer sale, a sold cup identity, or inventory movement history.

## RR-URG-07A must provide the minimum durable schema and deterministic demo seed required by later sale API, public trace, UI, and platform-evidence subtickets.

Success means an empty or existing retail database can receive the additive migration and idempotent seed, producing traceable sold cups without exposing an API in this ticket.

## Context Loaded

- `docs/CONTEXT.md`
- `docs/work/BACKLOG.md`
- `docs/standards/README.md`
- `docs/standards/QUALITY_BAR.md`
- `docs/standards/VALIDATION.md`
- `docs/work/phases/PHASE-2.md`
- `docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/ticket.md`
- `docs/requirements/REQUIREMENTS.md`
- `docs/requirements/MASTER_DATA.md`
- `docs/architecture/ERD.md`
- `src/apps/retail-service/internal/domain/models.go`
- `src/apps/retail-service/internal/app/init.go`
- `src/apps/retail-service/internal/usecase/service.go`
- `src/apps/retail-service/internal/seed/`
- `src/apps/retail-service/migrations/000001_init.up.sql`

## Brownfield Scope

- Touched modules/files:
  - retail domain models.
  - retail migration `000002`.
  - retail seed loader and checked-in seed JSON.
  - retail seed usecase and tests.
  - retail startup migration model list.
- Direct dependencies inspected:
  - existing store IDs and manager mapping.
  - current `SeedData` orchestration from auth-service.
  - GORM and Postgres migration conventions.
- Contracts affected: retail database schema and internal seed result/counts.
- Known unknowns:
  - Exact upstream batch/harvest/warehouse IDs available in the current warehouse seed must be verified before final seed values are committed.
  - Existing demo databases may contain schema drift; migration must be run against the current local `retail_db`, not only a fresh DB.
- Scope expansion: none. No HTTP API, trace projection, QR rendering, or UI work.

## Small Task Exemption

- Small task exemption: no
- Reason: changes persistent schema and deterministic seed behavior.
- Impact checked: API=no, DB=yes, Security=no, Runtime=yes, Standards=no

## Identifier Contract

- `retail_menu_items.id` is the stable chain-wide menu item identifier and is exposed as `menu_item_id` in application/API contracts.
- `retail_sale_items.menu_item_id` is a foreign key to `retail_menu_items.id`.
- `retail_sale_items.product_id` is the unique identity of one sold cup/item.
- `retail_sale_items.trace_code` equals `product_id` for public URL compatibility.
- A QR token alone does not encode provenance. Backend resolves provenance through sale item -> inventory lot -> upstream lineage.

## Proposed Approach

1. Add domain models and enums:
   - `RetailMenuItem`
   - `RetailInventoryLot`
   - `RetailSale`
   - `RetailSaleItem`
   - `RetailStockMovement`
   - sale status, movement type, and reference type constants.
2. Add additive migration `000002_retail_sales.up.sql`:
   - create the five tables without altering existing `stores` or `orders`.
   - add foreign keys inside `retail_db`.
   - keep upstream lineage IDs as indexed string/UUID-compatible values without cross-database foreign keys.
   - add unique indexes for menu item `product_code`, `sku`, `invoice_no`, `idempotency_key`, sold-item `product_id`, and `trace_code`.
   - add store/SKU/lot/sold-time query indexes needed by later APIs.
3. Include the five models in startup `AutoMigrate` only as a development compatibility layer. The SQL migration remains the reviewed schema contract.
4. Add checked-in seed files for menu items, inventory lots, sales, sale items, and stock movements.
5. Extend the seed loader to validate:
   - non-empty IDs and required fields.
   - unique menu item IDs/SKUs.
   - sale item menu item and lot references.
   - unique sold-item `product_id` and `trace_code == product_id`.
   - movement references and balance arithmetic.
6. Extend `SeedData` to run all retail seed writes in one GORM transaction:
   - upsert stores/menu items/lots.
   - upsert sales and sale items by deterministic IDs.
   - upsert movements by deterministic IDs.
   - set inventory-lot `available_quantity` to the deterministic expected balance.
7. Preserve current force behavior without deleting user-created records. Re-running seed updates deterministic rows only and does not duplicate stock movements.

## Data Constraints

- Money uses `DECIMAL(12,2)`, represented by the existing project numeric convention for this slice.
- Quantities use `DECIMAL(12,3)` to support cups/items now and possible weighted stock later.
- Sale status is constrained to `COMPLETED` for seeded history; later APIs may introduce additional lifecycle states.
- `quantity_delta` is positive for `RECEIVED`, negative for `SOLD`, and signed for `ADJUSTMENT`.
- `available_quantity >= 0`.
- Seed invariant:

```text
received_quantity = sum(RECEIVED quantity_delta)
available_quantity = sum(all quantity_delta)
```

- Stock movement history is authoritative. The initial `RECEIVED` movement is included once, and `received_quantity` is a denormalized intake snapshot used for reconciliation.

## Failure Handling

- Invalid embedded seed data fails before any database write.
- Any database failure rolls back the full retail seed transaction.
- Missing upstream lineage identifiers fail validation for public QR seed lots.
- Repeated seed calls produce the same row counts and balances.
- Existing unrelated stores/orders are untouched.

## Alternatives Considered

| Alternative | Pros | Cons | Decision |
| --- | --- | --- | --- |
| Store sold cups in existing `orders.items` JSONB | Minimal schema | Cannot enforce unique cup IDs, lot lineage, stock history, or efficient store queries | Rejected |
| Generate QR/provenance directly from product SKU | Simple QR | Incorrect business meaning; an unsold menu item is not a sold cup | Rejected |
| Use only GORM `AutoMigrate` | Less SQL | Unsafe against existing database drift and weak as a production schema contract | Rejected |
| Create sale and trace rows in trace-service only | Easy public lookup | Makes trace read model the operational source of truth | Rejected |

## Impacted Areas

- Code/modules: `apps/retail-service` domain, app initialization, seed, usecase, tests.
- Product behavior: deterministic demo state gains real sold cups and store inventory.
- API/contracts: no public API in 07A.
- Data/schema: five new retail tables and indexes.
- Security/auth: no change; seed remains internal/admin-orchestrated.
- Deployment/runtime: migration order gains `000002`; startup auto-migration includes additive models.
- Docs: ERD, master data, ticket, validation matrix, context, changelog.

## Test Plan

- Unit:
  - valid seed loads.
  - empty/missing identifiers rejected.
  - duplicate SKU, invoice, sold-item `product_id`, or movement ID rejected.
  - unknown product/lot references rejected.
  - `trace_code != product_id` rejected.
  - negative resulting inventory rejected.
- Integration:
  - apply `000001` then `000002` to an empty retail DB.
  - apply `000002` against current demo retail DB.
  - run seed twice and assert stable counts/balances.
  - force seed updates deterministic rows without duplicates.
  - transaction rolls back on an injected invalid reference/database error.
- E2E: not required for schema/seed-only 07A; covered by 07D/07E.
- UAT: not required because no user-facing flow is added.
- Platform/manual:
  - inspect table definitions and indexes with `psql`.
  - query Hoan Kiem and District 1 inventory and sold cups.
- Docs review: reconcile schema and seed facts with ERD/master data.

## Validation Matrix Impact

- Update required: yes
- Row: `Public Sold-Cup QR Trace`
- 07A may mark unit/integration seed/schema evidence only; E2E/platform remain planned until later subtickets.

## Reconciliation Plan

- Requirements docs: no behavior change beyond already accepted RR-URG-07 semantics.
- Architecture docs: no architecture-boundary change.
- API docs: no change.
- ERD docs: update with new retail entities and relationships.
- ADR: create a schema/data-ownership ADR because this introduces durable operational tables and identifier semantics.
- Context: update after implementation and verification.

## Approval Gate

Implementation is blocked until the human changes `approval: pending` to `approval: approved` or explicitly approves this design in chat.
