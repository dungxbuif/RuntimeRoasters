# Design

## Documentation Architecture

`DEMO_FLOW_CONTRACT.md` is the canonical source of truth for demo behavior and simulation boundaries for this feedback item. Product docs should link to it when discussing demo execution, current-vs-planned behavior, or simulation labels.

Story history remains archive-only. It can explain why decisions were made, but it must not override the canonical contract.

## Truth Labels

Each demo flow uses one label:

| Label | Meaning |
| --- | --- |
| `REAL_RUNTIME` | Backed by current runtime services, persistence, auth, and events. |
| `BACKEND_SIMULATED` | Created or fast-forwarded by backend/demo seeding while preserving target data shape. |
| `BROWSER_SIMULATED` | Driven by client-side timers, route replay, or UI actions that post to real backend APIs. |
| `UI_MOCK_ONLY` | Visual placeholder, fallback, or static UX element not backed by current source-of-truth data. |
| `PLANNED` | Intended future behavior, not current demo behavior. |

## Flow Matrix

| Flow | Truth Label | Notes |
| --- | --- | --- |
| OIDC login | `REAL_RUNTIME` | Kratos/Hydra login and local token handling are real runtime behavior. |
| Admin bootstrap | `REAL_RUNTIME` plus `BACKEND_SIMULATED` | Bootstrap endpoints are real; historical demo data may be inserted directly for immediate dashboards. |
| Farm harvest | `REAL_RUNTIME` | Farm service persists harvests and emits events. |
| Warehouse pickup | `REAL_RUNTIME` | Warehouse and logistics coordinate pickup requests, assignments, and shipment milestones. |
| Driver route simulation | `BROWSER_SIMULATED` plus `REAL_RUNTIME` | Route points come from static browser data; accepted GPS updates are persisted and published by backend. |
| Retail paid order | `REAL_RUNTIME` | Retail order creation, payment intent, warehouse reservation, and downstream events are runtime flow. |
| Payment webhook/simulated payment | `BACKEND_SIMULATED` plus `REAL_RUNTIME` | Provider settlement is demo/simulated; webhook signature and business validation still run. |
| Warehouse stock reservation | `REAL_RUNTIME` | Warehouse service owns reservation logic and emits stock events. |
| Logistics delivery | `REAL_RUNTIME` plus `BROWSER_SIMULATED` | Assignment and milestones are runtime; GPS motion source is browser route replay. |
| Trace CQRS | `REAL_RUNTIME` | Trace consumes Kafka events and builds Postgres/Elasticsearch read models. |
| Finance dashboard | `REAL_RUNTIME` plus demo controls | Ledger reads payment data; Pass/Fail buttons send signed demo webhook payloads. |
| Public QR trace | `PLANNED` or route-specific | The public trace document shape is supported conceptually, but shortcut/list UX may still be missing. |

## Multi-Account Demo Standard

Use a separate browser profile or isolated browser context for each persona. Do not use multiple accounts in the same browser profile unless the test intentionally includes logout plus Kratos session reset.

Recommended contexts:

- Profile A: `ADMIN`
- Profile B: `STORE_MGR`
- Profile C: `WAREHOUSE_MGR`
- Profile D: `DRIVER`
- Public/incognito: root topology or public trace

## Cross-Link Strategy

- Product docs point readers to this feedback contract for demo truth labels.
- Specialized docs explain logistics realtime, trace CQRS, order fulfillment, multi-account sessions, and UI audit.
- Later implementation tickets can promote these docs into a permanent product docs location if the team wants them outside feedback history.
