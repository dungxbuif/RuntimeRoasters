---
artifact_type: test_verification
id: TEST-RR-URG-07A
status: complete
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

## PostgreSQL Service Seed

- Pass:
  `RETAIL_TEST_DATABASE_URL=postgres://user:password@rr-postgres:5432/retail_07a_test?sslmode=disable go test ./apps/retail-service/internal/usecase -run TestRetailSeedPostgresIntegration -count=1 -v`
- Runtime: Go 1.26.2 container on the same Docker network as PostgreSQL 16.
- Result: migration-backed seed ran twice idempotently; counts, trace codes,
  ledger balances, and 210 materialized availability rows passed.

## Startup Schema Policy

- Retail startup no longer runs GORM `AutoMigrate`.
- Schema must be created by the reviewed SQL migration before service startup.
- Admin seed creates/reconciles business data only; it does not create tables.

## UAT

Not required for schema/seed-only 07A. User-facing UAT belongs to 07D/07E.
