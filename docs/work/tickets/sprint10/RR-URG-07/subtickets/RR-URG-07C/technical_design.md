---
artifact_type: detail_design
id: DESIGN-RR-URG-07C
status: done
owner: ai
approval: approved
trace:
  backlog_item: BL-001
  requirement: REQ-QR-001
  phase: PHASE-2
  ticket_or_bug: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07C/ticket.md
  test_verification: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07C/test_verification.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: docs/work/tickets/sprint10/RR-URG-07/subtickets/RR-URG-07C/docs_review.md
  adrs:
    - docs/decisions/0010-retail-sales-and-availability.md
---

# RR-URG-07C Detail Design

## Status

- Status: `done`
- Approval: approved by the human instruction to complete all RR-URG-07 work on 2026-06-06.

## Problem

Retail owns the sold cup and its lot lineage, while trace-service owns public
trace composition. The public QR token cannot be resolved by trace-service
today, and the existing client trace page is static.

## Context And Brownfield Scope

- Read Harness context, standards, Phase 2, RR-URG-07A/B/C, API/SDD, Retail
  internal sold-item contract, and trace-service PostgreSQL/Elasticsearch paths.
- Touch trace-service config/usecase/routes/tests, Retail internal read contract,
  KrakenD, and API/SDD docs.
- No new persistent table is required. Retail remains operational truth;
  trace-service composes a public response at read time.

## Approach

1. Extend Retail internal support with a bounded sold-item list and store display
   fields. Both internal endpoints keep `X-Internal-Secret` enforcement and are
   never exposed through KrakenD.
2. Add a timeout-bound Retail HTTP client owned by trace-service.
3. Add anonymous `GET /v1/public/trace-items` and
   `GET /v1/public/trace/:trace_code` routes.
4. Resolve `trace_code == product_id` through Retail, query trace events and
   documents by source batch/harvest IDs, and compose one public DTO.
5. Return only whitelisted business fields and sanitized milestones. Never
   expose raw event payloads, user IDs, auth data, secrets, or payment details.
6. PostgreSQL remains the fallback when Elasticsearch is unavailable.

## Public Contract

- Detail includes sold cup, invoice, product/store display, lot/batch/harvest/
  warehouse IDs, ordered business milestones, service participation, trace IDs,
  and `public_safe: true`.
- List contains only persisted sold cups and their trace URLs.
- Unknown product returns `404`; Retail timeout/unavailability returns `502`.
- List limit is bounded to 100; default is 25.

## Alternatives

| Alternative | Decision |
| --- | --- |
| Copy Retail sale rows into trace DB synchronously | Rejected: duplicates operational truth and misses seeded sales. |
| Derive origin from QR token | Rejected: token has no provenance semantics. |
| Return raw trace documents/events | Rejected: unsafe public surface. |

## Test And Reconciliation

- Unit: sanitizer/DTO composition, unknown item, Retail dependency failure,
  bounded list, and trace-code mapping.
- HTTP integration: no-auth detail/list and 404.
- Live: seeded Retail product resolves through running Retail + Trace services.
- Reconcile API, SDD, ADR if boundaries change, validation matrix, context,
  backlog, ticket evidence, and changelog.
