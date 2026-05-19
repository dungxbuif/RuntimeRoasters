# [RR-5] Observability Stack — SigNoz & Full Infrastructure

- **Goal:** Complete `docker-compose.yaml` with SigNoz and advanced components.
- **Business Value:** Provides comprehensive system monitoring, helping to detect performance issues and errors early.
- **Priority:** `HIGH`

## 🔍 Acceptance Criteria

### Scenario 1: SigNoz & ClickHouse
- **Then:** `rr-signoz` and `rr-clickhouse` run stably, receiving OTLP at `4317/4318`.

### Scenario 2: Advanced (Storage & Broker)
- **Then:** Redpanda, Cassandra, and Elasticsearch start successfully, ready for subsequent Sprints.
