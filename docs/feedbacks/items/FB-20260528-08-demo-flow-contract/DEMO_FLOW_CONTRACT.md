# Demo Flow Contract

This document is the source of truth for the current Runtime Roasters demo flow boundaries for `FB-20260528-08`.

## Glossary

| Term | Meaning |
| --- | --- |
| Real runtime | The current app executes through services, auth, persistence, and events. |
| Backend simulation | Backend/demo tooling creates or fast-forwards records/events for demo readiness. |
| Browser simulation | The browser drives a timed or scripted action, usually posting to real backend APIs. |
| UI mock | Static, fallback, or visual-only data that is not the current source of truth. |
| Planned | Desired future behavior that is not implemented as current demo runtime. |

## Canonical Demo Storyline

1. `ADMIN` logs in through OIDC and initializes deterministic demo data if needed.
2. A farm harvest or seeded historical harvest starts the supply-chain story.
3. Warehouse receives pickup/stock work and assigns resources.
4. `DRIVER` uses the Driver Client to confirm milestones and replay route GPS.
5. `STORE_MGR` creates a retail order.
6. Payment service creates a demo payment intent.
7. Finance user presses Pass or Fail to send a signed demo webhook.
8. Warehouse reserves stock only after successful payment.
9. Logistics fulfills delivery and return-to-base milestones.
10. Trace service consumes emitted events and builds queryable provenance documents.
11. Dashboards and topology views show operational state, traceability, and observability context.

## Current Flow Matrix

| Flow | Label | Current Contract |
| --- | --- | --- |
| OIDC login | `REAL_RUNTIME` | Kratos/Hydra issue sessions/tokens; app stores local token and calls protected APIs. |
| Admin bootstrap | `REAL_RUNTIME`, `BACKEND_SIMULATED` | Bootstrap endpoints are real; historical read-model data may be inserted directly for an immediate demo state. |
| Farm harvest | `REAL_RUNTIME` | Farm service owns harvest creation and emits events. |
| Warehouse pickup | `REAL_RUNTIME` | Warehouse and logistics services coordinate pickup request, assignment, and milestone events. |
| Driver route simulation | `BROWSER_SIMULATED`, `REAL_RUNTIME` | Browser replays static route points; backend validates, stores GPS, and publishes events. |
| Retail paid order | `REAL_RUNTIME` | Retail order creation triggers event-driven payment, warehouse, logistics, trace, and audit flow. |
| Payment webhook/simulated payment | `BACKEND_SIMULATED`, `REAL_RUNTIME` | Provider settlement is simulated/demo; signature and business validation remain meaningful. |
| Warehouse stock reservation | `REAL_RUNTIME` | Reservation is service-owned runtime logic; compensation/refund events handle failure paths. |
| Logistics delivery | `REAL_RUNTIME`, `BROWSER_SIMULATED` | Assignment/milestones are backend state; moving vehicle positions are browser-replayed GPS. |
| Trace CQRS | `REAL_RUNTIME` | Kafka events build Postgres and Elasticsearch read models. |
| Finance dashboard | `REAL_RUNTIME` plus demo controls | Reads payment ledger; Pass/Fail sends signed demo webhook payloads. |
| Public QR trace | `PLANNED` | Public trace shape is part of product direction; current UI still needs shortcuts/listing for easy lookup. |
| Saga visualizer/topology | `REAL_RUNTIME`, `UI_MOCK_ONLY` depending section | Root topology can read trace topology/history; any static visual fallback must be treated as illustrative. |

## Intentionally Omitted Flows

- Real customer checkout.
- Real payment settlement with external provider clearing.
- Human multi-step order approval.
- Real driver mobile device telematics.
- Public customer account lifecycle.
- Full accounting/finance system.
- Farm/store dual confirmation.
- Multi-warehouse optimization.
- Real route planning, geofencing, ETA, and dispatch optimization.

## Why These Are Omitted

The demo focuses on distributed-system architecture and role-scoped operations: OIDC login, service-local authorization, Kafka choreography, transactional eventing, read models, auditability, and operational dashboards. The omitted flows would add product scope and external dependency weight without improving the main showcase goal.

## Specialized Docs

- [LOGISTICS_REALTIME.md](LOGISTICS_REALTIME.md)
- [TRACE_CQRS.md](TRACE_CQRS.md)
- [ORDER_APPROVAL_AND_FULFILLMENT.md](ORDER_APPROVAL_AND_FULFILLMENT.md)
- [MULTI_ACCOUNT_DEMO.md](MULTI_ACCOUNT_DEMO.md)
- [DEMO_UI_AUDIT.md](DEMO_UI_AUDIT.md)
