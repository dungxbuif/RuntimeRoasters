---
artifact_type: feedback_log
id: FEEDBACK_LOG
status: active
owner: shared
human_fields:
  - raw_feedback
  - source
  - triage_override
ai_fields:
  - type_recommendation
  - severity
  - converted_artifact
  - notes
shared_fields:
  - feedback_items
  - status
updated: 2026-06-05
---

# Feedback Log

This file replaces the deleted legacy `docs/feedbacks/**` tree as the Harness v1 intake funnel.

Legacy feedback was reviewed from `4404600^` and consolidated here. Detailed implementation notes were merged into current tickets, standards, or architecture docs when still relevant.

## Triage Rules

- Raw user feedback enters this file first.
- Bugs convert to `docs/work/bugs/`.
- Features/enhancements convert to `docs/work/tickets/` or `docs/requirements/USER_STORIES.md`.
- Closed feedback remains here as durable context unless it requires an ADR.

## Feedback Items

| ID | Date | Raw Feedback | Source | Type | Status | Converted Artifact | Notes |
| --- | --- | --- | --- | --- | --- | --- | --- |
| FB-20260528-01 | 2026-05-28 | Architecture diagram needs fixed layout and grouped frames. | Legacy `docs/feedbacks` | Enhancement | closed | `docs/architecture/ui-ux/DESIGN.md` | Implemented as fixed topology/grouped architecture view. |
| FB-20260528-02 | 2026-05-28 | Buttons need pointer cursor. | Legacy `docs/feedbacks` | Bug | closed | `docs/standards/CODE.md` | Global/component styling expectation retained. |
| FB-20260528-03 | 2026-05-28 | Improve deterministic seed data, roles, finance integrity, and traceability dashboard data. | Legacy `docs/feedbacks` | Feature | closed | `docs/requirements/MASTER_DATA.md` | Seed and demo baseline rules retained; RR-URG-07 adds sold-cup seed requirements. |
| FB-20260528-04 | 2026-05-28 | Logistics map should use light mode and hide default routes. | Legacy `docs/feedbacks` | Enhancement | closed | `docs/architecture/ui-ux/DESIGN.md` | Current AGENTS rules also require active-shipment route filtering. |
| FB-20260528-05 | 2026-05-28 | Admin/policy gaps: ADMIN should not run farm/retail domain actions; add admin resource setup and Warehouse Manager fleet assignment. | Legacy `docs/feedbacks` | Feature | closed | `docs/requirements/domain/README.md`, `docs/decisions/0006-two-gate-authz-casbin.md` | Role mandates retained in AGENTS and domain docs. |
| FB-20260528-06 | 2026-05-28 | Retail orders API/UI routing, CORS, and method issues. | Legacy `docs/feedbacks` | Bug | closed | `docs/architecture/API.md` | KrakenD endpoint concerns retained in API/architecture docs. |
| FB-20260528-07 | 2026-05-28 | Casbin login/session bug around Kratos session and missing OAuth challenge. | Legacy `docs/feedbacks` | Bug | closed | `docs/engineering/TROUBLESHOOTING.md` | Troubleshooting note retained. |
| FB-20260528-08 | 2026-05-28 | Demo flow contract: distinguish real runtime, backend simulation, browser simulation, and UI mock. | Legacy `docs/feedbacks` | Feature | active | `docs/requirements/REQUIREMENTS.md`, `docs/work/tickets/sprint10/` | RR-URG sprint is the active execution surface. |
| FB-20260528-09 | 2026-05-28 | Docs gap/streamlining and redundant explorer cleanup. | Legacy `docs/feedbacks` | Maintenance | active | `docs/work/BACKLOG.md`, `docs/work/TRACEABILITY.md` | Harness v1 migration continues under docs reconciliation. |

## Current Intake Notes

- RR-URG-07 changed from generic product QR to UI-issued `product_id` for sold cups/items.
- Payment remains out of scope for the public cup purchase flow.
- Backend derives origin from `product_id` plus `retail_inventory_lot` lineage.
