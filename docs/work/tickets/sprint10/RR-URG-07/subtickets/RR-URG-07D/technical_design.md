---
artifact_type: detail_design
id: DESIGN-RR-URG-07D
status: done
owner: ai
approval: approved
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  ticket_or_bug: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07D/ticket.md
  test_verification: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07D/test_verification.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07D/docs_review.md
---

# RR-URG-07D Detail Design

## Status

- Status: `done`
- Approval: approved by the human instruction to complete all RR-URG-07 work on 2026-06-06.

## Problem

There is no usable public sale UI, invoice QR, or real public trace page. Store
operators also cannot see sold cups in the existing dashboard.

## Design Direction

- Extend the existing white/slate/cobalt Runtime Roasters visual system.
- Public kiosk uses a two-column layout: store/menu selection on the left and a
  persistent order summary on the right.
- Invoice is a focused success modal with one product ID, trace link, and real
  generated QR per physical cup.
- Public trace prioritizes business milestones; service/trace metadata is
  secondary and collapsible/quiet.
- Store dashboard adds a sold-cups table using existing DataGrid conventions.
- Responsive layouts collapse to one column without hiding primary actions.

## Behavior

1. Load public stores/menu availability.
2. Generate one `RR-CUP-*` product ID per selected quantity before submission.
3. Send one API item per cup and a generated `Idempotency-Key`.
4. Render invoice and QR links from the returned sale.
5. Fetch public trace by route code with loading, 404, failure, and success
   states; remove all static demo journey data.
6. Fetch authenticated sold items by selected store in the store dashboard.

## Test And Reconciliation

- Unit/component: product ID generation, request expansion, invoice QR/link,
  trace loading/404/success, sold-item table.
- E2E: mocked fail-fast buy -> invoice -> trace and dashboard sold-item flow.
- Platform: real gateway/services for one seeded/new product in RR-URG-07E.
- Reconcile API-facing client types, SDD UI flow, context, matrix, tickets,
  evidence, and changelog. No ERD change.
