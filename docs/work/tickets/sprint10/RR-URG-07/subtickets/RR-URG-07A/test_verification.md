---
artifact_type: test_verification
id: TEST-RR-URG-07A
status: in_review
owner: ai
updated: 2026-06-06
trace:
  ticket: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/ticket.md
  detail_design: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07A/technical_design.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07A Test Verification

## Automated

- Pass:
  `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/retail-service/...`
- Covered:
  - 42 menu items match required size/price cells.
  - deterministic 25-cup demo dataset.
  - seed idempotency and stable record counts.
  - ledger-to-lot balance reconciliation.
  - 210 materialized store/menu availability rows.
  - re-seed preserves additional runtime stock movements.

## PostgreSQL Migration

- Pass: applied `000001_init.up.sql` to fresh `retail_07a_test` with
  `psql -v ON_ERROR_STOP=1` on PostgreSQL 16.
- Result: all tables and indexes were created without SQL errors.

## Remaining Verification

- Skipped: opt-in PostgreSQL service seed integration test.
- Command:
  `RETAIL_TEST_DATABASE_URL=... go test ./apps/retail-service/internal/usecase -run TestRetailSeedPostgresIntegration -count=1 -v`
- Reason: localhost network access was denied by the execution sandbox and the
  escalation session did not complete.
- Residual risk: seed upserts passed on SQLite and schema passed on PostgreSQL
  separately, but their combined live PostgreSQL path remains to be run.

## UAT

Not required for schema/seed-only 07A. User-facing UAT belongs to 07D/07E.
