# ADR 0002: Use SigNoz For Demo Observability

**Status:** Accepted - 2026-05-23

## Context
The final demo needs one observability surface for OpenTelemetry traces, metrics, and logs. RR-URG-01 also needs waterfall evidence for the same `trace_id` already stored in trace-service business documents.

## Decision
Use SigNoz + ClickHouse in `deployments/docker-compose.dev.yaml`.

- SigNoz UI is exposed at `http://localhost:3301`.
- OTLP gRPC/HTTP are exposed at `localhost:4317` and `localhost:4318`.
- The SigNoz ClickHouse coordinator is observability infrastructure only.
- Kafka remains Apache Kafka without ZooKeeper.

## Consequences
- Local infrastructure is heavier, but the production-demo environment can show real OTel traces.
- Trace-service remains the business traceability/read-model boundary; SigNoz is not the business trace source of truth.
