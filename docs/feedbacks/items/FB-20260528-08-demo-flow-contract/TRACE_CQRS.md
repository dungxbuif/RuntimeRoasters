# Trace CQRS

## Current Implementation State

- Trace service consumes Kafka topics that represent traceable business events.
- It parses CloudEvents and extracts business IDs from extensions and event data.
- It stores event rows in Postgres `trace_events`.
- It rebuilds per-entity documents in Postgres `trace_documents`.
- It upserts Elasticsearch documents into `coffee_traceability`.
- Query path uses Elasticsearch first and Postgres fallback.
- Role scoping can filter non-admin access by store/warehouse/driver context.
- OpenTelemetry trace IDs enrich technical observability, but the business trace works from domain events and business IDs.

## Projection Path

```text
Kafka CloudEvent
  -> trace-service consumer
  -> extract batch/order/shipment/harvest/store/driver IDs
  -> insert trace_events
  -> rebuild trace_documents
  -> upsert Elasticsearch coffee_traceability
  -> UI lookup by known entity ID
```

## Simulated Pieces

- Bootstrap may directly insert historical read-model rows for immediate demo dashboards.
- OTel traces are observability enrichment, not the only source of business provenance.

## Omitted Pieces

- Traceability list/shortcut UI for discovering known entity IDs.
- Full public QR flow with unauthenticated customer journey polish.
- Rich search UI over Elasticsearch fields.
- Repair/replay tooling for large event history rebuilds.

## Traceability Page Caveat

The current traceability page requires a known Batch ID, Shipment ID, Order ID, or Harvest ID. If no seed or live flow has produced a known entity ID, the page can look empty. This is a demo data/discovery issue, not proof that trace CQRS is absent.

## Recommended Demo Script

1. Reset and bootstrap demo data, or run a fresh order/payment/logistics flow.
2. Capture a known order, shipment, harvest, or batch ID.
3. Open `/dashboard/traceability`.
4. Search for the known ID.
5. Confirm the document timeline shows events such as order creation, payment, warehouse reservation, logistics assignment, GPS/status, and delivery.
6. If Elasticsearch is down, verify Postgres fallback behavior through the same API path.

## Review Language

Use: "Trace CQRS is real runtime projection from Kafka events; demo seed may fast-forward read-model state."

Avoid: "Traceability is only a mock because the page starts empty."
