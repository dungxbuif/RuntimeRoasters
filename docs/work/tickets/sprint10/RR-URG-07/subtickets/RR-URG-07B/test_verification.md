---
artifact_type: test_verification
id: TEST-RR-URG-07B
status: complete
owner: ai
updated: 2026-06-06
trace:
  ticket: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07B/ticket.md
  detail_design: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07B/technical_design.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07B Test Verification

## Automated Unit And HTTP Integration

- Pass:
  `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/retail-service/... ./pkg/events/...`
- Covered:
  - sale validation and duplicate product rejection.
  - atomic sale item, FIFO lot decrement, movement, availability, and outbox.
  - idempotent replay without duplicate movement or event.
  - manager store-scope filtering.
  - anonymous public menu/sale routes.
  - internal sold-item route rejects a missing secret and returns lineage with
    the configured secret.
  - `retail.sale.completed` participates in strict CloudEvent round-trip tests.

## Live PostgreSQL

- Fresh database: `rr_urg_07b_test` on PostgreSQL 16.
- Pass: applied `000001_init.up.sql` with `psql -v ON_ERROR_STOP=1`.
- Pass:
  `RETAIL_TEST_DATABASE_URL=postgres://user:password@localhost:54321/rr_urg_07b_test?sslmode=disable GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/retail-service/internal/usecase -run TestRetailSeedPostgresIntegration -count=1 -v`
- Result:
  - deterministic seed ran twice with stable counts and reconciled balances.
  - two simultaneous sales competed for one available cup.
  - exactly one sale succeeded and one returned inventory conflict.
  - final lot balance was zero, proving no negative inventory or oversell.

## Configuration And Hygiene

- `jq empty deployments/krakend/krakend.json`: pass.
- `git diff --check`: pass.

## UAT

API-level UAT passed through Gin route integration tests. Browser sale and QR
journey verification is not part of 07B and remains assigned to RR-URG-07D/07E.
