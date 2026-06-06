---
artifact_type: detail_design
id: DESIGN-RR-URG-07B
status: done
owner: ai
approval: approved
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
  ticket_or_bug: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07B/ticket.md
  test_verification: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07B/test_verification.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07B/docs_review.md
  adrs:
    - docs/decisions/0010-retail-sales-and-availability.md
  master_docs_touched:
    - docs/architecture/API.md
    - docs/architecture/SDD/README.md
---

# RR-URG-07B Detail Design

## Status

- ID: `DESIGN-RR-URG-07B`
- Status: `done`
- Ticket: `RR-URG-07B`
- Approval: `approved` in chat on 2026-06-06
- Author: AI
- Updated: 2026-06-06

## Problem

RR-URG-07A provides the retail schema, deterministic master/demo data, inventory
ledger, sold-cup identity, and materialized availability. Retail still lacks the
HTTP contracts and transactional use case needed to browse a menu, create a
no-payment demo sale, inspect sold cups, and emit the event consumed by later
trace work.

Success means a fresh seeded environment is immediately usable through public
sale APIs while preserving inventory correctness under concurrent PostgreSQL
requests and enforcing store scope on manager reads.

## Context Loaded

- `docs/CONTEXT.md`
- `docs/work/BACKLOG.md`
- `docs/standards/README.md`
- `docs/standards/QUALITY_BAR.md`
- `docs/standards/VALIDATION.md`
- `docs/work/phases/PHASE-2.md`
- `docs/work/tickets/sprint10/RR-URG-07/ticket.md`
- `docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/technical_design.md`
- `docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07B/ticket.md`
- `docs/requirements/REQUIREMENTS.md`
- `docs/architecture/API.md`
- `docs/architecture/SDD/README.md`
- `docs/decisions/0010-retail-sales-and-availability.md`
- `src/apps/retail-service/internal/app/`
- `src/apps/retail-service/internal/domain/`
- `src/apps/retail-service/internal/usecase/`
- `src/pkg/events/`
- `deployments/krakend/krakend.json`

## Brownfield Scope

- Touched modules: retail HTTP routes, sale/menu use cases, event contracts,
  gateway routes, tests, and API/SDD work docs.
- Direct dependencies: 07A tables and seed, identity store scope, shared HTTP
  guards, CloudEvents, transactional outbox, and Kafka producer worker.
- Contracts affected: public and authenticated HTTP APIs plus
  `retail.sale.completed`.
- Known unknowns: trace-service projection and public trace response belong to
  RR-URG-07C and are not implemented here.
- Scope expansion: none. No payment, QR image generation, trace projection, or
  frontend work.

## Small Task Exemption

- Small task exemption: no
- Reason: adds public API, authenticated API, internal API, transactional data
  behavior, and an integration event contract.
- Impact checked: API=yes, DB=no, Security=yes, Runtime=yes, Standards=no

## API Contract

### Public Menu

- `GET /v1/public/retail/menu-items` returns active chain-wide menus and their
  active menu items. It does not claim store availability.
- `GET /v1/public/retail/stores/:store_id/menu-items` returns active menu items
  joined to `store_menu_inventories`, including `available` and
  `available_units`.
- Both endpoints are deliberately outside authentication middleware.
- The store endpoint returns `404` for an unknown/inactive store.

### Public Sale

`POST /v1/public/retail/sales` accepts:

```json
{
  "store_id": "STORE-HN-HOANKIEM",
  "items": [
    {
      "menu_item_id": "MI-ESPRESSO-SINGLE",
      "product_id": "RR-CUP-20260606-0001"
    }
  ]
}
```

- One request item represents one physical sold cup. There is no `quantity`
  aggregation because every cup requires its own unique `product_id`, lineage,
  trace code, and QR URL.
- `product_id` is supplied by the client and must be non-empty and unique both
  within the request and across persisted sale items.
- `Idempotency-Key` is a required header and is persisted on `sales`. Reusing a
  key returns the existing completed sale; using a new key with an existing
  `product_id` returns `409 Conflict`.
- A successful `201 Created` response includes invoice data, totals, and each
  sold item's `product_id`, `trace_code`, lineage identifiers, and `public_url`.
- Validation failures return `400`; missing store/menu item returns `404`;
  duplicate identity or exhausted inventory returns `409`.

### Manager And Internal Reads

- Authenticated routes:
  - `GET /v1/retail/sales`
  - `GET /v1/retail/sales/:sale_id`
  - `GET /v1/retail/stores/:store_id/sold-items`
- Store managers are restricted by `identity.StoreScopeFromContext`; broader
  roles may use their existing unrestricted scope. A single resource outside
  scope is returned as `404` to avoid disclosing its existence.
- `GET /v1/internal/retail/sold-items/:product_id` returns the operational
  sold-item and lineage record required by trace-service.
- The internal route requires an exact `X-Internal-Secret` match and is not
  exposed as a public KrakenD route.

## Sale Transaction

All writes happen in one PostgreSQL transaction:

1. Validate request shape, idempotency key, unique product IDs, active store,
   and active menu items.
2. Resolve each `menu_item_id` to its stock SKU, serving quantity, display
   metadata, and price.
3. Process required SKUs in stable sorted order to reduce deadlock risk.
4. For each cup, select the first eligible inventory lot ordered by
   `received_at, id`, then lock it with `SELECT ... FOR UPDATE`.
5. Eligibility requires matching store/SKU, `AVAILABLE` status, non-expired
   stock, and `available_quantity >= serving_quantity`.
6. Insert one `sale`, one `sale_item` per cup, and one negative `SOLD`
   `stock_movement` per cup. Set `trace_code = product_id`.
7. Decrement the locked lot using a guarded update. Mark it `DEPLETED` only
   when its balance reaches zero; menu availability is computed separately.
8. Rebuild affected `store_menu_inventories` rows once after all movements.
   Request reads never execute runtime `SUM/FLOOR`; they read this materialized
   projection maintained by write paths.
9. Insert one `retail.sale.completed` CloudEvent into the outbox in the same
   transaction.

Any failure rolls back the sale, sale items, lot balances, movements,
availability projection, and outbox record together.

## Event Contract

- Topic/type: `retail.sale.completed`
- Source: `/services/retail-service`
- Subject: `sales/{sale_id}`
- CloudEvent metadata includes `event_id`, `correlationid`, `storeid`, event
  time, and trace context when available.
- Payload includes sale/invoice/store totals and one item payload per sold cup:
  `sale_item_id`, `menu_item_id`, `product_id`, `sku`, `product_name`, price,
  inventory lot, trace code, public URL, batch, harvest, and warehouse lineage.
- The topic is added to `events.TraceableTopics`; the existing outbox worker
  publishes it without a second database write path.

## Availability And Performance Decision

- `store_menu_inventories` remains the read-optimized projection for menu APIs.
- Availability is rebuilt only for affected store/SKU combinations after a
  committed inventory mutation. The implementation may reuse the 07A helper
  initially, but must expose an affected-SKU path so sale cost scales with the
  changed SKUs rather than all menu items.
- Required indexes from 07A support lot lookup by store/SKU/status/expiry and
  sold-item lookup by store/product/time.
- The authoritative audit source remains `stock_movements`; a reconciliation
  command/test compares projection values against ledger-derived balances.

## Failure And Concurrency Handling

- Row locking plus guarded decrement prevents overselling under concurrent
  PostgreSQL requests.
- Stable SKU lock ordering reduces cross-SKU deadlocks; a database deadlock is
  surfaced as a retryable conflict rather than partially committing.
- Expired, depleted, and insufficient lots are skipped. No eligible lot yields
  `409` with the unavailable menu item identified.
- Idempotent replay returns the original completed response without another
  movement or event.
- Event publish failure marks the row failed; the existing worker retries both
  pending and failed rows. It does not roll back an already committed sale.

## Alternatives Considered

| Alternative | Pros | Cons | Decision |
| --- | --- | --- | --- |
| Runtime `SUM/FLOOR` on every menu request | No projection writes | Cost grows with movement history and duplicates business logic | Rejected |
| Random lot selection | Simple demo | Nondeterministic lineage and poor stock rotation | Rejected |
| Aggregate `quantity` per request item | Smaller request | Cannot assign one unique identity and lineage per physical cup | Rejected |
| Publish Kafka directly in the sale transaction | Immediate event | Database/event dual-write failure window | Rejected |
| Expose internal sold-item route through public gateway | Easy trace access | Leaks operational lineage surface | Rejected |

## Impacted Areas

- Code: retail app/usecase/domain DTOs, event contracts, gateway configuration,
  and focused tests.
- Product behavior: public menu browsing and no-payment cup purchase become
  immediately usable from seeded data.
- API/contracts: new public, manager, and internal endpoints.
- Data/schema: no migration change expected; uses RR-URG-07A schema.
- Security: public sale/menu routes are anonymous by design; manager routes use
  guards and store scope; internal route uses shared secret.
- Runtime: outbox starts producing `retail.sale.completed`.
- Docs: API master, SDD implementation/decision notes, ticket evidence, matrix,
  backlog, and context.

## Test Plan

- Unit:
  - request validation and one-cup item contract.
  - duplicate/empty `product_id` rejection.
  - menu/store validation and error mapping.
  - FIFO lot selection, serving decrement, movement, totals, and projection.
  - rollback and idempotent replay.
  - manager store-scope filtering and internal-secret rejection.
- PostgreSQL integration:
  - apply migration and seed to a clean live database.
  - create sale and verify invoice, item lineage, lot balance, movement,
    availability, and outbox atomically.
  - run concurrent requests against the final available units and prove no
    negative balance/oversell.
  - replay idempotency key and prove stable row counts.
- HTTP integration:
  - public endpoints work without a token.
  - protected endpoints reject anonymous callers and enforce manager scope.
  - internal endpoint rejects missing/wrong secret and accepts configured one.
- Contract:
  - `retail.sale.completed` round-trips as a valid CloudEvent.
- E2E/UAT: deferred to RR-URG-07D/07E for browser QR flow; API-level UAT is
  required here.

## Validation Matrix Impact

- Update required: yes
- Row: `Public Sold-Cup QR Trace`
- RR-URG-07B supplies unit, PostgreSQL integration, HTTP/auth, and event evidence;
  browser and public trace evidence remain planned for later subtickets.

## Reconciliation Plan

- Requirements: preserve accepted one-cup identity semantics; no new business
  requirement.
- API: document exact requests, responses, status codes, and auth boundaries.
- Architecture/SDD: record projection maintenance, FIFO locking, idempotency,
  outbox, and the decision not to calculate availability on read.
- ERD: no change unless implementation reveals a missing constraint.
- ADR: ADR-0010 remains authoritative; amend only if implementation changes a
  durable decision.
- Work docs: create verification and docs-review evidence, update matrix,
  context, backlog, and ticket status through the lifecycle.

## Approval Gate

Approved by the human in chat on 2026-06-06.
