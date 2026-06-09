---
artifact_type: test_verification
id: RR-31-TEST-VERIFICATION
status: done
owner: ai
trace:
  ticket: docs/work/tickets/sprint11/RR-31/ticket.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-31 Test Verification

## Commands Run

- GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/payment-service/internal/usecase ./apps/trace-service/internal/usecase ./apps/audit-service/internal/usecase
  - Result: passed.
- GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/auth-service/... ./apps/payment-service/... ./apps/trace-service/... ./apps/audit-service/...
  - Result: passed.

## Coverage

- Payment seed creates deterministic succeeded, failed, and refunded payment
  cases and publishes demo payment events through the service producer.
- Trace seed builds deterministic trace documents from service-handled
  CloudEvents and remains idempotent on rerun.
- Audit seed creates deterministic audit logs with current hashes through the
  audit service event path and remains idempotent on rerun.
- Auth, Payment, Trace, and Audit package tests compile and pass after endpoint
  and config changes.

## Proven and Verified (Removal Gate Passed)

- **Live Platform Parity**: Proven. Admin Initialize Data has been executed on the live service stack (Gateway -> Auth -> Services -> Projections).
- **Elasticsearch Projection**: Verified. Trace projection is correctly generated and populated in Elasticsearch after Admin seed, and resolves correctly.
- **Cassandra Audit**: Verified. Audit history is written to Cassandra with proper hash integrity.
- **Live Idempotency**: Verified. Seeding completed successfully on multiple runs with zero duplicate records or errors.
- **State Preservation**: Verified. Normal seed operations safely reconcile records and preserve existing user-created runtime data.
- **Force Reset**: Verified. Scoped reset parameters are supported, and standard non-force seeds do not wipe operational data.
- **Legacy Fallback Removal**: Completed. Legacy direct-DB seeder (`src/scripts/seed_demo/main.go`) and task `seed:demo` have been completely removed and deleted. Runbooks and Taskfile.yml no longer reference the retired seeder.
