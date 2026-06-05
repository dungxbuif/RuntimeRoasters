---
artifact_type: requirements_index
id: REQUIREMENTS
status: active
owner: shared
human_fields:
  - priority
  - acceptance
  - requirement_source
ai_fields:
  - requirement_rows
  - status_updates
  - trace_links
shared_fields:
  - functional_requirements
  - non_functional_requirements
updated: 2026-06-05
---

# Requirements

This index summarizes current accepted requirements. Detailed behavior remains in `docs/requirements/SPEC.md` and `docs/requirements/domain/README.md`.

## Functional Requirements

| ID | Requirement | Priority | Source | Status |
| --- | --- | --- | --- | --- |
| REQ-AUTH-001 | Users authenticate through OIDC and protected APIs enforce JWT plus service-local Casbin. | high | Product/security docs | active |
| REQ-FARM-001 | Farm Manager declares harvests for assigned farms. | high | Domain docs | active |
| REQ-WH-001 | Warehouse Manager dispatches pickup/delivery work and manages warehouse inventory operations. | high | Domain docs, RR-URG sprint | active |
| REQ-ORDER-001 | Store Manager creates paid retail orders that trigger payment, warehouse reservation, logistics delivery, and trace/audit events. | high | RR-URG-05 | active |
| REQ-TRACE-001 | Trace-service projects business events into queryable Postgres/Elasticsearch read models. | high | RR-URG-01, RR-URG-07 | active |
| REQ-QR-001 | Public QR trace resolves a UI-issued `product_id` for a sold cup/item, not an unsold catalog product. | high | RR-URG-07 user clarification | ready |
| REQ-SEED-001 | Demo seed creates deterministic users, resources, inventory, and trace-ready data. | high | Bootstrap docs, RR-URG-08 | active |
| REQ-UI-001 | Client app exposes role-scoped dashboards and public no-auth topology/QR pages. | high | RR-URG-09 | active |

## Non-Functional Requirements

| ID | Requirement | Category | Status |
| --- | --- | --- | --- |
| NFR-SEC-001 | Protected operations must fail closed when auth, role, or record scope is missing. | security | active |
| NFR-OBS-001 | Cross-service demo flows should preserve W3C trace context through HTTP, gRPC, Kafka, and DB operations. | observability | active |
| NFR-DATA-001 | Demo seed and migrations must be idempotent for local reset/replay. | reliability | active |
| NFR-TEST-001 | High-risk tickets require unit, integration, E2E, and platform evidence or explicit skip reasons. | validation | active |
| NFR-DOCS-001 | Behavior, API, schema, security, and runtime changes require docs reconciliation before closure. | process | active |

## Demo Simulation Boundaries

Runtime Roasters demo shortcuts must simulate real-world edges, not replace the backend source of truth.

| Label | Meaning | Allowed Examples |
| --- | --- | --- |
| `REAL_RUNTIME` | Normal application path through APIs, DB state, Kafka, auth, trace, and audit. | OIDC, RBAC/Casbin, retail orders, stock reservation, dispatch, trace projection. |
| `BACKEND_SIMULATED` | Idempotent seed/bootstrap data used to make dashboards immediately demo-ready. | Historical read-model rows, seeded users/resources/routes, provider settlement emulator. |
| `BROWSER_SIMULATED` | Browser behaves like a real external actor or device while backend validates and persists state. | Driver GPS route replay, timed milestone progression. |
| `UI_MOCK_ONLY` | Visual fallback only; never the source of truth for an accepted demo flow. | Offline placeholder visuals while an API is unavailable. |

Simulation rules:

- Payment simulation uses signed Stripe-compatible webhook paths; retail must not post unsigned mock payment success directly.
- Driver/GPS simulation replays checked-in route data from the browser and posts each accepted tick to logistics APIs.
- Warehouse dispatch, retail sale, inventory movement, trace, and audit state must be real backend state.
- Time compression is allowed only when domain transitions and Kafka events are still preserved.
- Public QR demo data must represent sold cups/items with UI-issued `product_id`; chain-wide menu entries use `menu_item_id`.
