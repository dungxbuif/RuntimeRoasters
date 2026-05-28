# Impact Analysis

## Affected Documentation

- `docs/product/GUIDE.md`
- `docs/product/TECH.md`
- `docs/product/SPEC.md`
- `docs/product/standards/BOOTSTRAP.md`
- `docs/feedbacks/USER_FEEDBACK.md`
- `docs/feedbacks/items/**`
- Historical story docs under `docs/stories/**`, which should remain archive context rather than canonical truth.

## Affected UI Context

- `/dashboard/logistics`
- `/dashboard/driver`
- `/dashboard/traceability`
- `/dashboard/finance`
- `/dashboard/retail`
- `/dashboard/retail/orders`
- `/dashboard/warehouse`
- root `/` architecture topology
- `/dashboard/topology-mesh`

## Backend Concepts

- Kafka event flow for retail, payment, warehouse, logistics, trace, audit, and topology.
- Valkey driver GPS storage and liveness.
- Trace Postgres plus Elasticsearch CQRS projection.
- Payment simulation and webhook validation.
- Warehouse stock reservation and dispatch events.
- Driver shipment milestone events.

## Demo Risks

- Reviewers may think a mock/fallback UI section is production behavior.
- Agents may implement against stale sprint history rather than current product truth.
- Multi-account manual testing can fail when multiple personas share one browser profile.
- Traceability can look empty or broken if the demo does not provide a known entity ID.
- Logistics route lines can be misread as active delivery routes even when there is no active shipment.

## Expected Impact Of This Docs Packet

- Clearer demo script.
- Fewer ambiguous implementation requests.
- Safer review language for real, simulated, mock, and planned behavior.
- A backlog basis for later code cleanup without changing runtime behavior in this ticket.
