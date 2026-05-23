# Observability

Runtime Roasters uses OpenTelemetry for runtime tracing and SigNoz/ClickHouse for local demo observability.

## Local Stack

- SigNoz UI: `http://localhost:3301`
- OTLP gRPC: `localhost:4317`
- OTLP HTTP: `localhost:4318`
- Collector config: `deployments/otel-collector-config.yaml`
- Compose services:
  - `signoz`
  - `signoz-clickhouse`
  - `otel-collector`
  - `signoz-telemetrystore-migrator`
  - `signoz-zookeeper-1`

`signoz-zookeeper-1` is only for ClickHouse/SigNoz coordination. Kafka remains the official Apache Kafka broker and does not use ZooKeeper.

## Setup Notes For Agents

Use this exact order when bringing up observability:

```bash
docker compose -f deployments/docker-compose.dev.yaml up -d signoz-zookeeper-1 signoz-init-clickhouse signoz-clickhouse signoz-telemetrystore-migrator otel-collector signoz
```

The service dependencies are intentional:

- `signoz-zookeeper-1`: ClickHouse/SigNoz coordination only. Never connect Kafka to it.
- `signoz-init-clickhouse`: downloads and installs the `histogramQuantile` executable function used by SigNoz queries.
- `signoz-clickhouse`: telemetry store for traces, metrics, logs, metadata, and meter data.
- `signoz-telemetrystore-migrator`: runs bootstrap, sync, and async ClickHouse migrations before the collector starts.
- `otel-collector`: SigNoz OTel collector listening on `4317/4318`.
- `signoz`: UI/query service exposed at `localhost:3301`.

Files to update together:

- `deployments/docker-compose.dev.yaml`
- `deployments/otel-collector-config.yaml`
- `deployments/signoz/clickhouse/cluster.xml`
- `deployments/signoz/clickhouse/custom-function.xml`

Common failure points:

- ClickHouse starts before `signoz-init-clickhouse` completes: the custom `histogramQuantile` function may be missing.
- Collector starts before `signoz-telemetrystore-migrator` exits successfully: trace tables may be missing.
- Collector starts with OpAMP manager mode before SigNoz org/agent enrollment exists: the runtime config can be replaced by `nop` receivers/exporters. In local dev, run collector directly with `--config=/etc/otel-collector-config.yaml`.
- Services log `connection refused localhost:4317`: `otel-collector` is not running or host port `4317` is occupied.
- SigNoz UI is healthy but traces are missing: verify service env uses `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317` for host-local services, then check `docker logs rr-otel-collector`.
- Kafka references any ZooKeeper/coordinator service: this is wrong. Kafka remains independent Apache Kafka.

## Trace Evidence Rule

RR-URG-01 is only fully closed when one live demo flow shows the same `trace_id` in:

- service logs,
- Kafka propagated `traceparent`,
- trace-service Postgres/Elasticsearch documents,
- SigNoz traces.

## Service Endpoint

Host-local Go services should export OTLP to:

```text
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
```

If a service later runs inside Docker, use:

```text
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
```
