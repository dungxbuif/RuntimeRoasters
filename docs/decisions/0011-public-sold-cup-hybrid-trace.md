---
artifact_type: adr
id: ADR-0011
status: accepted
owner: shared
decision_approval: approved_in_chat_2026-06-06
trace:
  requirement: REQ-QR-001
  phase: PHASE-2
  tickets_or_bugs:
    - RR-URG-07C
    - RR-URG-07D
  master_docs:
    - docs/architecture/API.md
    - docs/architecture/SDD/README.md
---

# ADR 0011: Public Sold-Cup Hybrid Trace

## Context

Retail owns sold cups and inventory lineage. Trace-service owns public journey
composition. Seeded sold cups exist before Kafka projection, and raw trace
documents contain operational fields unsuitable for anonymous responses.

## Decision

- Treat `trace_code` as the UI-issued Retail `product_id`.
- Trace-service calls a secret-protected Retail internal API for the sold cup
  and lineage, then combines it with trace PostgreSQL/Elasticsearch data.
- Return a strict public DTO rather than raw events/documents.
- Keep seeded and newly sold cups immediately resolvable even when a Kafka
  projection or Elasticsearch is unavailable.
- Keep the Retail internal API off KrakenD; only trace-service public list and
  detail routes are exposed anonymously.

## Alternatives

- Duplicate all sold-cup rows into trace storage: rejected due to stale copies
  and seeded-data projection gaps.
- Publicly expose Retail lineage endpoints: rejected due to operational data
  leakage and ownership confusion.
- Encode provenance into QR tokens: rejected because IDs are identity, not
  trustworthy lineage documents.

## Consequences

- Positive: immediate demo readiness and explicit ownership.
- Positive: a small whitelist defines the public security boundary.
- Negative: public trace availability depends on a synchronous Retail call.
- Negative: the shared internal secret remains demo-grade and should become
  service identity/mTLS in production.

