---
artifact_type: traceability_matrix
id: TRACEABILITY
status: active
owner: shared
human_fields:
  - requirement_source
  - release_scope
ai_fields:
  - trace_links
  - evidence_links
  - adr_links
shared_fields:
  - matrix_rows
updated: 2026-06-05
---

# Traceability

This file maps high-level Runtime Roasters requirements to tickets, evidence, decisions, and release notes.

## Trace Matrix

| Requirement | Phase | Ticket/Bug | Detail Design | Test Verification | Docs Review | ADR/Docs | Release |
| --- | --- | --- | --- | --- | --- | --- | --- |
| OIDC and two-gate authorization | Foundation | Sprint auth/RBAC work | Existing auth docs | Unit/E2E evidence in ticket history | Active docs reconciled | `docs/decisions/0006-two-gate-authz-casbin.md` | `docs/releases/CHANGELOG.md` |
| Farm-to-warehouse pickup flow | Emergency sprint | RR-URG-03, RR-URG-04 | Ticket docs | Pending final sprint evidence | Active | `docs/work/tickets/sprint10/` | TBD |
| Paid order fulfillment and delivery return | Emergency sprint | RR-URG-05 | Ticket docs | Passed targeted evidence before RR-URG-07 | Active | `docs/work/tickets/sprint10/RR-URG-05/ticket.md` | TBD |
| Realtime notifications | Emergency sprint | RR-URG-06 | Ticket docs | Unit, integration, E2E, platform evidence recorded | Active | `docs/work/tickets/sprint10/RR-URG-06/ticket.md` | TBD |
| Public sold-cup QR trace | Emergency sprint | RR-URG-07A to RR-URG-07E | Required per split ticket | Pending | Active | `docs/work/tickets/sprint10/RR-URG-07/ticket.md` | TBD |
| Deterministic demo seeding | Bootstrap/demo | RR-URG-08 | Ticket docs | Pending final seed verification | Active | `docs/requirements/MASTER_DATA.md` | TBD |

## Rules

- Every RR-URG-07 split ticket must update this matrix when implementation evidence exists.
- Do not mark public QR trace verified until unit, integration, E2E, and platform evidence are attached or explicitly scoped out.
- ADRs are required for durable changes to API/schema/auth/runtime boundaries.
