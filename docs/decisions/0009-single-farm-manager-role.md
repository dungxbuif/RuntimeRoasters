---
artifact_type: adr
id: ADR-0009
status: accepted
owner: shared
human_fields:
  - decision_approval
  - final_status
ai_fields:
  - context
  - alternatives_considered
  - consequences
  - linked_work
shared_fields:
  - decision
  - trace
trace:
  requirement: REQ-AUTH-001
  phase: PHASE-2
  tickets_or_bugs:
    - RR-URG-07A
  detail_design: docs/work/tickets/sprint10/RR-URG-07A-DETAIL_DESIGN.md
  master_docs:
    - docs/requirements/MASTER_DATA.md
    - docs/requirements/REQUIREMENTS.md
  release_notes: docs/releases/CHANGELOG.md
---

# ADR 0009: Use One Farm Operator Role

## Status

- Status: accepted
- Date: 2026-06-06
- Decision approval: explicit human direction in the active work session

## Context

The system exposed two overlapping farm roles. No demo identity used the additional administrative farm role, while its policies duplicated or overlapped `FARM_MANAGER`. Keeping both values caused drift across Kratos schema, Casbin policies, farm usecases, frontend types, seed documentation, and tests.

## Decision

Use `FARM_MANAGER` as the only farm-domain operator role.

- `ADMIN` remains a separate system/resource administration role with explicit policies.
- `FARM_MANAGER` operates farms and harvests subject to ownership scoping.
- No secondary farm-administration role is accepted by identity schemas, UI types, policies, usecases, or seed data.

## Alternatives Considered

- Keep two farm roles with different privileges: rejected because no accepted business distinction or seeded persona exists.
- Rename `FARM_MANAGER` to the administrative role: rejected because `FARM_MANAGER` is already the active identity, policy, UI, and seed vocabulary.

## Consequences

- Positive: one canonical farm role across identity, authorization, UI, seed, and docs.
- Positive: removes duplicate Casbin policy and special-case scoping behavior.
- Negative: any external identity carrying the removed value must be migrated to `FARM_MANAGER` before its next login/token issuance.
- Neutral: `ADMIN` permissions remain explicit and are not inherited through a farm role.

## Linked Work

- Authorization model: `docs/decisions/0006-two-gate-authz-casbin.md`
- Master roles: `docs/requirements/MASTER_DATA.md`
- Seed preparation: `docs/work/tickets/sprint10/RR-URG-07A-DETAIL_DESIGN.md`
