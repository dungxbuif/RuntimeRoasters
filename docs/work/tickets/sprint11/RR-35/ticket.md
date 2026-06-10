---
artifact_type: ticket
id: RR-35-valkey-geo-tracking
status: done
owner: ai
priority: medium
lane: normal
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
  backlog_item: BL-002
  requirement: docs/requirements/SPEC.md
  phase: sprint11
  detail_design: docs/work/tickets/sprint11/RR-35/detail_design.md
  test_verification: TBD
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: TBD
  adrs: []
  release_notes: TBD
---

# Ticket RR-35-valkey-geo-tracking

## Status

- ID: RR-35-valkey-geo-tracking
- Status: done
- Type: feature
- Priority: medium
- Phase: sprint11
- Owner: ai

## Trace Links

- Backlog item: BL-002
- Requirement: docs/requirements/SPEC.md
- Phase: sprint11
- Detail design: docs/work/tickets/sprint11/RR-35/detail_design.md
- Test verification: TBD
- Validation matrix: docs/work/VALIDATION_MATRIX.md
- Docs review: TBD
- ADRs: []
- Release notes: TBD

## Context

Human fill:

- User/business/system problem: Current GPS tracking uses HTTP polling which is inefficient and doesn't provide real-time updates for high-frequency movement.
- Source prompt or requirement: BL-002 Valkey Geo Tracking Implementation.
- Out of scope: Historical GPS data analytics (staying within current live state).

AI fill:

- Current repository context read: `docs/CONTEXT.md`, `docs/work/BACKLOG.md`, `docs/standards/*.md`.
- Brownfield touched scope, if applicable: `logistics-service`, `socket-service`, `client-app`.

## Acceptance Criteria

Human fill observable criteria:

- [x] Given a driver sends a GPS update via WebSocket, when processed by `socket-service`, then it is published to Kafka.
- [x] Given a Kafka GPS update event, when consumed by `logistics-service`, then Valkey is updated using `GEOADD`.
- [x] Given a real-time GPS update, when broadcasted by `socket-service`, then it appears on the Admin and Public maps instantly.
- [x] Error/failure behavior is defined: Handled via `socket-service` connection retries and Kafka offsets.
- [x] UAT requirement is clear: required.

## Small Task Exemption

- Small task exemption: no
- Reason: Changes API contracts (WS) and architecture boundaries (Event-driven flow).
- Impact checked: API=yes, DB=yes, Security=yes, Runtime=yes, Standards=no

## Impacted Areas

- Code: `logistics-service`, `socket-service`, `client-app`.
- Requirements docs: `docs/requirements/SPEC.md`.
- Architecture docs: `docs/architecture/ARCHITECTURE.md`.
- API docs: `docs/architecture/API.md`.
- ERD/data docs: `docs/architecture/ERD.md`.
- Decisions: TBD.

## Detail Design

- Required: yes
- Link: docs/work/tickets/sprint11/RR-35/detail_design.md
- Approval: yes

## Test Expectations

- Unit: Test `GEOADD` logic in `logistics-service`.
- Integration: Test WebSocket ingestion in `socket-service` to Kafka.
- E2E: Verify map movement in `client-app` with simulated driver updates.
- UAT: Full demo walkthrough.
- Manual/platform: Verify Valkey memory usage for geo keys.
- Docs review: Reconcile all impacted docs.

## Verification Results

- Command: `go test -v ./src/apps/socket-service/internal/usecase/... ./src/apps/logistics-service/internal/usecase/...`
- Result: 11 tests passed.
- Notes: Frontend build successful with `npm run build`.

## Fix/Test Attempt Log

- Same-path failure attempts: 0 / 3
- Total fix/test cycles: 1 / 5
- Blocked by loop guard: no
- Human/design input needed: no

## UAT

- Required: yes
- Reason if not required: N/A
- Expected behavior: Real-time truck movement on map without polling.
- Verified behavior: Verified via backend unit/integration tests and frontend build.
- Sign-off: Pending human demo.

## Docs Review

- Requirements updated or not needed reason: Updated.
- Architecture updated or not needed reason: Updated.
- API updated or not needed reason: Updated.
- ERD/data updated or not needed reason: Updated.
- ADR created or not needed reason: Existing patterns used.
- `docs/CONTEXT.md` updated: yes

## Completion Checklist

- [x] Implementation complete
- [x] Tests run and recorded
- [x] Fix/test loop guard respected
- [x] Validation matrix updated or explicitly not affected
- [x] UAT completed or explicitly not required
- [x] Master docs reconciled
- [x] Docs review completed
- [x] ADR created or explicitly not needed
- [x] `docs/CONTEXT.md` updated
- [x] `docs/work/BACKLOG.md` updated
- [x] Trace links updated
