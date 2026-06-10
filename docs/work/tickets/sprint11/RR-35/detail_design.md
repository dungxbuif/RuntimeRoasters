---
artifact_type: detail_design
id: DESIGN-RR-35-valkey-geo
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
  backlog_item: BL-002
  requirement: docs/requirements/SPEC.md
  phase: sprint11
  ticket_or_bug: docs/work/tickets/sprint11/RR-35/ticket.md
  test_verification: TBD
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: TBD
  adrs: []
  master_docs_touched: [docs/architecture/ARCHITECTURE.md, docs/architecture/API.md]
---

# Detail Design: Valkey Geo Tracking (RR-35)

Current Phase: Detail Design | Next action: implementation of the plan

## Problem

- Problem statement: Current GPS tracking uses HTTP polling which is inefficient for real-time visualization.
- Why now: BL-002 is prioritized for the real-time logistics demo.
- Source ticket/bug: docs/work/tickets/sprint11/RR-35/ticket.md
- Success outcome: Real-time, event-driven GPS movement on maps with three visibility tiers.

## Context Loaded

- `docs/CONTEXT.md`
- `docs/work/BACKLOG.md`
- `docs/standards/QUALITY_BAR.md`
- `docs/standards/VALIDATION.md`
- `src/apps/logistics-service/internal/usecase/service.go`
- `src/apps/socket-service/internal/usecase/service.go`
- `src/pkg/events/contracts.go`

## Brownfield Scope

- Touched modules/files: 
    - `src/apps/logistics-service/internal/app/app.go`
    - `src/apps/logistics-service/internal/app/init.go`
    - `src/apps/logistics-service/internal/usecase/service.go`
    - `src/apps/socket-service/internal/usecase/service.go`
    - `src/apps/socket-service/internal/app/init.go`
    - `src/apps/client-app/src/app/(dashboard)/dashboard/logistics/page.tsx`
    - `src/apps/client-app/src/app/(dashboard)/dashboard/driver/page.tsx`
- Direct dependencies inspected: Kafka, Redis (Valkey), WebSocket.
- Contracts affected: WebSocket `gps_update` message, Kafka `logistics.gps.updated` event.
- Known unknowns: Exact memory overhead of Valkey Geo keys under high-frequency updates (to be monitored).
- Scope expansion reason: N/A.

## Small Task Exemption

- Small task exemption: no
- Reason: Changes API contracts (WS) and architecture boundaries (Event-driven flow).
- Impact checked: API=yes, DB=yes, Security=yes, Runtime=yes, Standards=no

## Proposed Approach

1. **Ingress (Socket Service):**
   - Update `readPump` to handle `{"type": "gps_update", "payload": {...}}`.
   - Verify `DRIVER` role via `client.scope.Role`.
   - Publish `events.LogisticsGPSUpdated` CloudEvent to Kafka.
2. **Persistence (Logistics Service):**
   - Consume `events.TopicLogisticsGPSUpdated`.
   - Call `GEOADD` on Valkey key `drivers:locations` with `driver_id` as member.
   - Update PostgreSQL `vehicles` current coordinates for historical/snapshot consistency.
3. **Egress (Socket Service):**
   - Consume Kafka and broadcast as `topology.event` to all relevant clients.
   - Apply visibility filtering: Public clients see anonymized payload (Name/Plate only).
4. **Frontend (Client App):**
   - Refactor `LogisticsMap` to accept `mode` prop.
   - Root page: `public` mode (anonymized).
   - Manager Dashboards: `pro` mode (full telemetry).
   - Admin Dashboard: `admin` mode (+ metrics).

## Alternatives Considered

| Alternative | Pros | Cons | Decision |
| --- | --- | --- | --- |
| HTTP Polling (Current) | Simple, stateless. | High overhead, not real-time, battery drain on mobile. | Rejected |
| SSE (Server-Sent Events) | Simpler than WS for egress. | Uni-directional, still needs HTTP for ingress updates. | Rejected |

## Impacted Areas

- Code/modules: `logistics-service`, `socket-service`, `client-app`.
- Product behavior: Real-time map updates.
- API/contracts: New WebSocket message types.
- Data/schema: Valkey Geo keys, PostgreSQL vehicle coordinates.
- Security/auth: Role-based WS ingestion and egress filtering.
- Deployment/runtime: Kafka and Valkey dependencies must be healthy.
- Docs: SPEC, ARCHITECTURE, API.

## Test Plan

- Unit: `logistics-service` usecase tests for `GEOADD` and Kafka consumption.
- Integration: `socket-service` test for WS -> Kafka flow.
- E2E: Playwright test for driver update appearing on admin map.
- UAT: Manual visual demo.
- Manual/platform: Check Valkey geo keys via `redis-cli`.
- Docs review: Reconcile all docs.

## Validation Matrix Impact

- Update required: yes
- Row(s): BL-002
- Reason if no update required: N/A

## Reconciliation Plan

- Requirements docs: update
- Architecture docs: update
- API docs: update
- ERD docs: update
- ADR: no change (covered by design patterns)
- Context: update

## Docs Review

- Code changed but docs unchanged reason: N/A
- Requirements updated or not needed reason: Updated to reflect real-time tracking.
- Architecture updated or not needed reason: Updated for event-driven flow.
- API updated or not needed reason: Updated for WS contracts.
- ERD/data updated or not needed reason: Updated for Valkey Geo storage.
- ADR created or not needed reason: Existing patterns used.

## Verification Results

- Command: TBD
- Result: TBD
- Notes: TBD
