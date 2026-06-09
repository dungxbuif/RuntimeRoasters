---
artifact_type: detail_design
id: RR-31-DETAIL-DESIGN
status: done
owner: ai
trace:
  ticket: docs/work/tickets/sprint11/RR-31/ticket.md
  tests: docs/work/tickets/sprint11/RR-31/test_verification.md
  docs_review: docs/work/tickets/sprint11/RR-31/docs_review.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  master_docs:
    - docs/architecture/SDD/MASTER_DATA.md
    - docs/architecture/API.md
---

# RR-31 Detail Design: Complete Admin Demo Seed

## Decision

Admin Initialize Data is the canonical demo seed entry point. The legacy
task seed:demo direct-DB path remains present only until live parity evidence
passes across PostgreSQL, Elasticsearch, Cassandra, Kafka, and service APIs.

## Implemented Slice

- auth-service seed orchestration fans out to Payment, Trace, and Audit when
  service URLs are configured, with local and docker fallback URLs.
- Payment exposes internal POST /v1/system/seed and GET /v1/system/status
  and writes deterministic succeeded, failed, and refunded demo payment rows
  through payment-owned repositories.
- Trace exposes internal seed/status endpoints and builds deterministic trace
  documents by feeding canonical CloudEvents through trace-service handlers.
- Audit exposes internal seed/status endpoints and builds deterministic audit
  rows by feeding canonical CloudEvents through audit-service handlers, keeping
  audit hash creation inside audit-service.
- src/apps/auth-service/.env.example documents the new downstream service URLs
  used by the seed orchestrator.

## Seed Data Semantics

- Demo prefix: rr31-demo.
- Reference time: 2026-06-08T08:00:00Z.
- Seed reruns are deterministic and service-owned.
- Normal seed does not truncate operational tables.
- force=true is accepted by the service endpoints for orchestration parity, but
  the current implemented slice is non-destructive. The destructive reset policy
  remains behind the RR-31 removal gate.

## Available Mechanism

The optimized availability design from Sprint 7 remains unchanged: retail menu
availability is pre-materialized into store-owned menu inventory state instead
of repeatedly summing lot movement floors at runtime. RR-31 uses that seed
contract as input and adds downstream finance/projection history, rather than
changing retail inventory availability tables.

## Legacy Seeder Gate

Do not remove task seed:demo or src/scripts/seed_demo/main.go until:

- Admin seed has live proof for Payment, Trace/Elasticsearch, and
  Audit/Cassandra.
- Admin seed is run twice on a live stack with stable counts.
- User-created runtime rows survive a normal re-seed.
- Force reset behavior is explicitly approved and tested.
- Docs and runbooks no longer refer to direct DB seeding as a supported path.

## Known Unknowns

- Live Elasticsearch projection parity is not yet proven in this slice.
- Live Cassandra audit parity is not yet proven in this slice.
- Auth UI result aggregation may need a follow-up if the frontend expects only
  the original four service statuses.
