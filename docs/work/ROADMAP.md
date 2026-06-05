---
artifact_type: roadmap
id: ROADMAP
status: active
owner: human
human_fields:
  - milestones
  - priority
  - phase_order
ai_fields:
  - phase_links
  - status_summaries
shared_fields:
  - milestone_status
updated: 2026-06-05
---

# Roadmap

## Field Ownership

- Human owns milestones, priority, and phase order.
- AI maintains phase links and status summaries.

Use this file to group work into milestones or major phases.

## Roadmap Rules

- A roadmap item becomes executable only after it has a phase file in `docs/work/phases/`.
- Each phase should contain one or more tickets or bugs.
- Completed phases should link to release notes or changelog entries when relevant.

## Milestones

| Milestone | Goal | Status | Phase Files |
| --- | --- | --- | --- |
| M1 | Phase 1: System Bootstrap & Seeding | done | TBD |
| M2 | Phase 2: Logistics Real-time Tracking & Saga Simulator | done | TBD |
| M3 | Phase 3: Reliability & Traceability | active | [PHASE-3-RELIABILITY.md](phases/PHASE-3-RELIABILITY.md) |

## Legacy Roadmap Notes

Pre-Harness sprint discussion mapped the post-farm-management roadmap as:

- security and core infrastructure hardening.
- distributed order orchestration.
- inventory consistency and SAGA behavior.
- logistics and realtime delivery.
- geographic intelligence.
- traceability/search.
- reliability and observability.
- commerce/payment integration.
- control-plane and resiliency demo.

Current execution is driven by backlog rank and phase files, not by legacy sprint discussion files.
