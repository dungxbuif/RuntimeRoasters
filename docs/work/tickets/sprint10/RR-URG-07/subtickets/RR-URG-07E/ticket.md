---
artifact_type: ticket
id: RR-URG-07E
status: ready
owner: shared
priority: high
lane: normal
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  detail_design: not_required
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-URG-07E: Evidence And Platform Verification

## Goal

Close RR-URG-07 with reliable evidence across unit, integration, E2E, and platform checks.

This ticket depends on [RR-URG-07A](../RR-URG-07A/ticket.md) through [RR-URG-07D](../RR-URG-07D/ticket.md).

## Evidence Scope

Verify the complete simplified flow:

1. Seed store products and retail inventory lots.
2. UI or API creates a demo sale without payment.
3. UI-issued `product_id` is persisted on sale item.
4. Retail inventory decreases.
5. Stock movement is recorded.
6. Trace-service resolves public trace by `product_id`.
7. Public UI shows real sold-cup trace.

## Required Test Types

Unit:

- retail sale validation.
- duplicate `product_id` rejection.
- inventory decrement transaction.
- public trace sanitizer.
- trace `product_id` lookup logic.

Integration:

- retail migrations and seed.
- retail sale API writes invoice, items, inventory, movement, outbox.
- trace public API resolves seeded and newly created `product_id`.
- store manager sold-items API is scoped.

E2E:

- public POS buy flow.
- invoice and QR display.
- sold-items by store.
- public trace by QR.

Platform:

- local infra plus required services.
- one instance per Kafka consumer.
- Postgres evidence for retail sale and stock movement.
- trace-service Postgres/Elasticsearch evidence for public trace.
- public HTTP response for same `product_id`.

## Playwright Guardrails

- Use production runtime when dev server is unstable.
- Set explicit action/navigation/test timeouts.
- Prefer mocked auth for manager-only UI checks.
- Avoid waiting on open-ended realtime streams.
- Do not leave browser/server processes running after evidence unless user explicitly asks.

## Acceptance Criteria

- All required evidence is captured in `docs/work/tickets/sprint10/evidence/`.
- Harness matrix marks RR-URG-07 split tickets with unit/integration/E2E/platform evidence as applicable.
- Any skipped live/platform step has a concrete reason and next command.
- Final report maps every acceptance criterion to evidence.
