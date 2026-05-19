# RR-URG-01: Verify And Fix End-to-End TraceId Propagation

## Priority

P0. This ticket must be completed before other emergency sprint implementation tickets.

## Problem

The final demo depends on showing a real journey through services. If `trace_id` is fragmented across HTTP, gRPC, Kafka, DB, and logs, the trace UI and public QR story will be unreliable.

## Scope

Verify and fix W3C trace context propagation across:

- KrakenD -> service HTTP/gRPC handler.
- service handler -> Postgres/GORM.
- service handler -> Kafka producer.
- Kafka message headers -> Kafka consumer.
- consumer handler -> downstream DB/event publishing.
- service logs.
- trace-service projection metadata.

## Implementation Details

### 1. KrakenD

- Inspect `deployments/krakend/krakend.json`.
- Ensure OTel is enabled for the gateway.
- Ensure `traceparent` is forwarded to backend services.
- If KrakenD creates a trace for inbound public requests, verify downstream services see the same root trace.

Checklist:

- [ ] `traceparent` present in request received by backend.
- [ ] Missing incoming `traceparent` still results in a generated trace.
- [ ] Protected and public routes both preserve trace context.

### 2. Shared Service Bootstrap

- Inspect `src/pkg/base/app.go`.
- Verify Gin middleware uses OTel.
- Verify gRPC server/client interceptors use OTel.
- Verify every relevant service initializes through the shared base app or equivalent OTel setup.

Checklist:

- [ ] farm-service traced.
- [ ] warehouse-service traced.
- [ ] retail-service traced.
- [ ] payment-service traced.
- [ ] logistics-service traced.
- [ ] trace-service traced.
- [ ] audit-service traced where applicable.

### 3. Kafka Producer/Consumer

- Inspect shared Kafka producer/consumer packages.
- Inject current OTel context into Kafka message headers using W3C propagator.
- Extract context from Kafka message headers before invoking business handlers.
- Preserve existing idempotency/event ID headers.

Required headers:

- `traceparent`
- `tracestate` if present
- existing `message_id` / event ID headers

Checklist:

- [ ] Producer injects `traceparent`.
- [ ] Consumer extracts `traceparent`.
- [ ] Consumer starts child span from extracted context.
- [ ] Tests cover publish/consume header roundtrip.
- [ ] Existing message idempotency still works.

### 4. Postgres/GORM

- Add/verify OTel GORM instrumentation in shared DB initialization.
- Ensure DB spans are children of current service handler/consumer spans.

Checklist:

- [ ] DB span visible for order create.
- [ ] DB span visible for harvest create.
- [ ] DB span visible for warehouse reservation.
- [ ] DB span visible for logistics shipment update.

### 5. Logs

- Ensure structured logs include `trace_id` when context has a span.
- Do not rely only on manual context values if OTel span context exists.

Checklist:

- [ ] service logs include `trace_id`.
- [ ] consumer logs include `trace_id`.
- [ ] error logs include `trace_id`.

### 6. Trace-Service Metadata

- Domain events should carry:
  - business IDs
  - `correlation_id`
  - OTel `trace_id` when available
- Trace-service should store `trace_id` as metadata, not as the only business identifier.

Checklist:

- [ ] trace document can show linked `trace_id`.
- [ ] trace document can still be queried by business ID or public trace code.

## Acceptance Criteria

1. Running one harvest flow produces a single connected trace from gateway -> farm -> Kafka -> warehouse.
2. Running one paid order flow produces a single connected trace from gateway -> retail -> payment -> warehouse -> logistics -> trace-service.
3. Kafka messages contain `traceparent`.
4. Consumer spans are children of producer/service spans.
5. Postgres spans appear under service handler or consumer spans.
6. Logs include the same `trace_id` for related operations.
7. A short verification guide or command output is added to the ticket/runbook.

## Test Checklist

- [ ] Unit: Kafka trace header injection.
- [ ] Unit: Kafka trace header extraction.
- [ ] Unit: missing traceparent starts a new trace gracefully.
- [ ] Integration: create harvest and verify trace continuity.
- [ ] Integration: create paid order and verify trace continuity.
- [ ] Manual: inspect SigNoz/Jaeger or configured OTel backend.
- [ ] Manual: inspect Kafka message headers during one flow.

## Done Evidence

Attach or record:

- Trace ID sample.
- Services/spans seen in trace.
- Kafka headers sample.
- Log lines showing same trace ID.
- Any remaining trace gaps.
