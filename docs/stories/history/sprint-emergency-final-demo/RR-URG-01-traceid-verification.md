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

- [x] `traceparent` forwarded by gateway config to backend routes.
- [x] Missing incoming `traceparent` still results in generated service trace context.
- [x] Protected and public routes preserve configured trace headers.
- [x] Runtime verification with live KrakenD request and service log sample.

### 2. Shared Service Bootstrap

- Inspect `src/pkg/base/app.go`.
- Verify Gin middleware uses OTel.
- Verify gRPC server/client interceptors use OTel.
- Verify every relevant service initializes through the shared base app or equivalent OTel setup.

Checklist:

- [x] farm-service traced through shared base app.
- [x] warehouse-service traced through shared base app.
- [x] retail-service traced through shared base app.
- [x] payment-service traced through shared base app.
- [x] logistics-service traced through shared base app.
- [x] trace-service traced through shared base app.
- [x] audit-service traced where applicable through shared base app.

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

- [x] Producer injects `traceparent`.
- [x] Consumer extracts `traceparent`.
- [x] Consumer starts child span from extracted context.
- [x] Tests cover Kafka header inject/extract roundtrip.
- [x] Existing message idempotency remains payload/topic/offset based and was not changed.

### 4. Postgres/GORM

- Add/verify OTel GORM instrumentation in shared DB initialization.
- Ensure DB spans are children of current service handler/consumer spans.

Checklist:

- [x] Shared GORM OpenTelemetry plugin installed in `database.NewPostgres`.
- [ ] DB span visible for order create in SigNoz.
- [ ] DB span visible for harvest create in SigNoz.
- [ ] DB span visible for warehouse reservation in SigNoz.
- [ ] DB span visible for logistics shipment update in SigNoz.

### 5. Logs

- Ensure structured logs include `trace_id` when context has a span.
- Do not rely only on manual context values if OTel span context exists.

Checklist:

- [x] service logs include `trace_id` when context has OTel span.
- [x] consumer logs include `trace_id` after Kafka extraction.
- [x] error logs include `trace_id` through `logger.FromContext` and problem renderer.

### 6. Trace-Service Metadata

- Domain events should carry:
  - business IDs
  - `correlation_id`
  - OTel `trace_id` when available
- Trace-service should store `trace_id` as metadata, not as the only business identifier.

Checklist:

- [x] trace document can show linked `trace_id`.
- [x] trace document can still be queried by business ID or public trace code.

## Acceptance Criteria

1. Running one harvest flow produces a single connected trace from gateway -> farm -> Kafka -> warehouse.
2. Running one paid order flow produces a single connected trace from gateway -> retail -> payment -> warehouse -> logistics -> trace-service.
3. Kafka messages contain `traceparent`.
4. Consumer spans are children of producer/service spans.
5. Postgres spans appear under service handler or consumer spans.
6. Logs include the same `trace_id` for related operations.
7. A short verification guide or command output is added to the ticket/runbook.

## Test Checklist

- [x] Unit: Kafka trace header injection.
- [x] Unit: Kafka trace header extraction.
- [x] Unit: missing traceparent starts a new trace gracefully.
- [x] Integration: create harvest and verify warehouse consumption.
- [x] Integration: create paid order and verify trace continuity.
- [x] Manual: inspect SigNoz at `http://localhost:3301`.
- [x] Manual: inspect Kafka message headers during one flow.

## Done Evidence

Attach or record:

- Code changes:
  - `src/pkg/telemetry/telemetry.go`: always installs an SDK tracer provider, even without OTLP endpoint, so local/demo runs still create valid trace IDs.
  - `src/pkg/kafka/producer.go`: starts producer span and injects W3C headers into Kafka messages.
  - `src/pkg/kafka/consumer.go`: extracts W3C headers, starts consumer child span, and passes traced context to handlers.
  - `src/pkg/kafka/propagation.go`: Kafka header carrier for `traceparent`/`tracestate`.
  - `src/pkg/database/postgres.go`: installs GORM OTel tracing plugin for all services using shared Postgres bootstrap.
  - `src/pkg/logger/logger.go`: reads `trace_id` from OTel span context before falling back to manual context.
  - `src/apps/trace-service/internal/domain/models.go`: persists `trace_id` on trace events and `trace_ids` on trace documents.
  - `deployments/krakend/krakend.json`: forwards `traceparent`/`tracestate` and allows them through CORS.
- Tests run:
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./pkg/kafka ./pkg/logger ./pkg/telemetry`
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/trace-service/internal/...`
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./pkg/base/casbin ./apps/farm-service/internal/infrastructure/repository/tests ./pkg/database ./pkg/kafka ./apps/trace-service/internal/...`
  - `GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./apps/... ./pkg/... ./runtime/...`
- Full service/pkg/runtime test status: passed.
- Known non-ticket full `go test ./...` caveat: `src/scripts` contains multiple standalone `func main` files in one package, so the all-package command still fails there unless scripts are split or excluded.
- Live evidence:
  - `docs/business/sprint-emergency-final-demo/evidence/rr-urg-01-trace-e2e-evidence.md`
  - Trace ID sample: `8991cbb271c04df9b233ef5f98d58a8c`
  - Order ID sample: `68669cad-06fb-4b6b-8583-6d77801b727b`
  - Harvest ID sample: `15`
  - SigNoz UI screenshot: `docs/business/sprint-emergency-final-demo/evidence/signoz-ui.png`
- Observability stack:
  - `deployments/docker-compose.dev.yaml` defines SigNoz, ClickHouse, the SigNoz OTel collector, and a ClickHouse coordination service.
  - SigNoz UI: `http://localhost:3301`.
  - OTLP endpoints: `localhost:4317` and `localhost:4318`.
  - Kafka remains independent Apache Kafka without ZooKeeper.
  - ClickHouse/SigNoz contains spans for the same `trace_id` as the trace-service document, including HTTP, gRPC, Kafka producer/consumer, and DB spans.
