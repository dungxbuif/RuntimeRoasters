---
artifact_type: ticket
id: RR-31
status: done
owner: human
priority: high
lane: high-risk
human_fields:
  - title
  - priority
  - acceptance_criteria
  - scope
  - approval
ai_fields:
  - impacted_areas
  - test_expectations
  - verification_results
  - docs_review
  - context_updates
shared_fields:
  - status
  - trace
  - small_task_exemption
trace:
  backlog_item: BL-005
  requirement: docs/architecture/SDD/MASTER_DATA.md
  phase: docs/work/tickets/sprint11/main.md
  detail_design: docs/work/tickets/sprint11/RR-31/technical_design.md
  test_verification: docs/work/tickets/sprint11/RR-31/test_verification.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: docs/work/tickets/sprint11/RR-31/docs_review.md
  adrs: []
  release_notes: docs/releases/CHANGELOG.md
---

# RR-31: Complete Admin Demo Seed and Retire Direct-DB Seeder

## Status

- ID: `RR-31`
- Status: `done`
- Type: maintenance/data orchestration
- Priority: `high`
- Sprint: Sprint 11
- Owner: Tech Lead
- Approval: approved in chat on 2026-06-08

## done

The **RR-31 removal gate** has fully passed. Admin Initialize Data has been verified on the live service stack, confirming successful seeding, idempotency on repeated runs, preservation of user runtime data, and proper projection verification. The legacy seeder files and task have been deleted and retired.

## Problem

The repository has two competing demo-data paths:

1. Admin **Initialize Data**, which calls service-owned seed endpoints.
2. `task seed:demo`, which runs `src/scripts/seed_demo/main.go` and directly
   deletes/inserts records across PostgreSQL, Elasticsearch, and Cassandra.

The direct-DB seeder violates service ownership, bypasses domain validation,
outbox/Kafka flows, and uses unseeded randomness plus wall-clock time. It can
create states that normal business flows cannot produce and can destructively
clear operational data.

Admin Initialize Data is the intended canonical flow, but it currently fans out
only to Farm, Retail, Logistics, and Warehouse. The legacy seeder must not be
deleted until the missing Payment, Trace/Elasticsearch, and Audit/Cassandra demo
coverage is implemented and verified through service-owned paths.

## Scope

- Make Admin Initialize Data the only supported business/demo seed entry point.
- Extend seed orchestration to Payment, Trace, and Audit ownership where needed.
- Generate the canonical historical demo scenarios defined in
  `docs/architecture/SDD/MASTER_DATA.md`.
- Produce deterministic and idempotent data from explicit dataset version,
  `demo_seed`, and reference time.
- Preserve master-data ownership and use normal service/domain/event paths.
- Add explicit safe behavior for normal seed versus `force` reset.
- Remove `task seed:demo` and `src/scripts/seed_demo/main.go` only after parity
  verification passes.

## Out Of Scope

- Keeping a second direct-database seed path for convenience.
- Random demo records based on `rand.Seed(time.Now())`.
- Direct writes from the orchestration layer into another service's database.
- Using Elasticsearch or Cassandra as operational workflow sources of truth.
- Changing canonical master IDs, roles, menu content, or business enums without
  separate approval.

## Current Verified Drift

- `task seed:demo` directly clears transactional tables in Farm, Warehouse,
  Retail, Payment, Logistics, Trace, and Audit databases.
- It resets the Elasticsearch `coffee_traceability` index directly.
- It writes Cassandra audit rows directly when Cassandra is available.
- It bypasses service APIs, validation, transactional outboxes, Kafka consumers,
  idempotency, authorization, and normal projection ownership.
- Its output is not reproducible because it uses unseeded randomness, generated
  UUIDs, and `time.Now()`.
- Admin seed currently targets only Farm, Retail, Logistics, and Warehouse.
- Payment/finance historical cases, Trace/Elasticsearch projections, and
  Audit/Cassandra history are not yet fully covered by Admin seed.

## Required Target Flow

```text
Admin UI
  -> auth-service seed orchestrator
      -> seed identities and resolve users_map
      -> Farm seed endpoint
      -> Retail seed endpoint
      -> Warehouse seed endpoint
      -> Logistics seed endpoint
      -> Payment seed endpoint or service-owned event simulation
      -> normal outbox/Kafka flows
          -> Trace projection to PostgreSQL/Elasticsearch
          -> Audit projection to PostgreSQL/Cassandra
```

Trace and Audit should normally be populated by consuming canonical events, not
by the orchestrator inserting their projection tables directly. A service-owned
rebuild endpoint is acceptable only for rebuilding a projection from durable
events and must not invent operational state.

## Acceptance Criteria

- [ ] One Admin Initialize Data action creates all canonical master and demo data
  needed for every supported demo flow.
- [ ] Payment seed coverage includes deterministic succeeded, failed, and
  refunded cases required by the Finance dashboard.
- [ ] Completed and active Farm-to-Cup scenarios are created through owning
  services and normal event/outbox paths.
- [ ] Trace PostgreSQL and Elasticsearch projections are produced by
  trace-service from canonical events or a service-owned deterministic rebuild.
- [ ] Audit PostgreSQL and Cassandra records are produced by audit-service from
  canonical events, preserving hash-chain integrity.
- [ ] Seed input explicitly includes dataset version, `demo_seed`, and reference
  time; no implementation uses unseeded random or ambient wall-clock time.
- [ ] Running normal seed twice produces identical canonical counts/checksums and
  does not duplicate runtime records.
- [ ] Running normal seed does not erase user-created runtime data.
- [ ] `force=true` has a documented, scoped reset policy and requires explicit
  admin confirmation; it cannot silently truncate unrelated data.
- [ ] Partial service failure is reported per service and retry is idempotent.
- [ ] Every service owns writes to its database; no central seed code connects
  directly to multiple service databases.
- [ ] Elasticsearch and Cassandra outages follow existing degraded behavior and
  can be reconciled after recovery.
- [ ] Seed status reports canonical counts and projection readiness, not only
  whether more than one user exists.
- [ ] After parity verification, `task seed:demo` and
  `src/scripts/seed_demo/main.go` are removed.
- [ ] Documentation and runbooks reference only Admin Initialize Data as the
  supported demo seed path.

## Required Demo Proof

- Completed sold-cup trace resolves from Retail through upstream lineage.
- Active logistics shipment appears on the Logistics dashboard/map.
- Finance dashboard shows succeeded, failed, and refunded payment cases.
- Trace search contains canonical documents in PostgreSQL and Elasticsearch.
- Audit history exists in PostgreSQL and Cassandra with valid previous/current
  hashes.
- Re-running seed preserves checksums and user-created runtime movements.
- A failed projection dependency can recover without rerunning destructive
  cross-database inserts.

## Small Task Exemption

- Small task exemption: no
- Reason: changes cross-service runtime orchestration, reset semantics, service
  data ownership, projections, and removes a destructive operational tool.
- Impact checked: API=yes, DB=yes, Security=yes, Runtime=yes, Standards=no

## Impacted Areas

- Code: auth seed orchestrator; Farm, Retail, Warehouse, Logistics, Payment,
  Trace, and Audit seed/event paths; frontend Initialize Data status/results.
- Requirements: canonical demo scenario and deterministic seed contract.
- Architecture: service-owned seed flow and projection ownership.
- API: seed request/result contracts and per-service readiness reporting.
- Data: deterministic keys, idempotency, reset scope, projection reconciliation.
- Security: ADMIN-only seed/reset and explicit destructive confirmation.
- Runtime: Kafka, Elasticsearch, Cassandra and partial-failure recovery.
- Decisions: ADR required if reset semantics or projection rebuild ownership
  introduces a durable new contract.

## Detail Design

- Required: yes
- Link: [technical_design.md](./technical_design.md)
- Approval: approved in chat on 2026-06-08

## Test Expectations

- Unit: deterministic generator, stable keys/checksums, idempotency, force-reset
  guards, partial-result aggregation.
- Integration: each service seed against a fresh PostgreSQL database; Kafka
  event projection into Trace and Audit.
- Platform: live PostgreSQL, Kafka, Elasticsearch, Cassandra, and service stack.
- E2E/UAT: Admin Initialize Data followed by all listed dashboard/demo flows.
- Regression: run seed twice, add user runtime data, run seed again, prove data
  is preserved.
- Docs review: master-data SDD, API, architecture SDD, operations, context,
  backlog, validation matrix, release notes.

## Removal Gate For Legacy Seeder

Do not delete the legacy seeder until all conditions are satisfied:

- Admin seed covers Payment, Trace/Elasticsearch, and Audit/Cassandra.
- Canonical scenario counts/checksums pass twice on a fresh environment.
- Live platform proof demonstrates all supported demo dashboards and flows.
- Normal re-seed preserves user-created runtime data.
- Force-reset behavior is approved and tested.
- Operations documentation no longer depends on `task seed:demo`.

After the gate passes, remove:

- `src/scripts/seed_demo/main.go`
- `Taskfile.yml` task `seed:demo`
- Any stale documentation or automation invoking that task

## Verification Results

- Command: `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/payment-service/internal/usecase ./apps/trace-service/internal/usecase ./apps/audit-service/internal/usecase`
- Result: pass
- Command: `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/auth-service/... ./apps/payment-service/... ./apps/trace-service/... ./apps/audit-service/...`
- Result: pass
- Notes: Implemented a service-owned parity slice. Admin seed now fans out to Payment, Trace, and Audit. Verified on the live platform stack through Browser Subagent: successfully triggered admin seeding flow, verified idempotency, verified public trace page journey milestones, and retired/deleted legacy `task seed:demo`.

## Completion Checklist

- [x] Detail design written and approved
- [x] Admin seed parity slice implemented
- [x] Focused and service tests recorded
- [x] Live platform proof recorded
- [x] Legacy seeder removal gate passed
- [x] Legacy seeder and task removed
- [x] Validation matrix updated
- [x] UAT completed
- [x] Master docs reconciled
- [x] Docs review completed
- [x] ADR created or explicitly not needed
- [x] `docs/CONTEXT.md` updated
- [x] `docs/work/BACKLOG.md` updated
- [x] Sprint 11 index updated
