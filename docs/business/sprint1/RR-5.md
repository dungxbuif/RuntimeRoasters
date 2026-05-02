# [RR-5] Observability Stack — SigNoz & Full Infrastructure

- **Summary:** Hoàn thiện `docker-compose.yaml` với SigNoz (ClickHouse) và các thành phần nâng cao (Redpanda, Cassandra, Elasticsearch).
- **Priority:** `HIGH`

---

## Acceptance Criteria

### Scenario 1: SigNoz & ClickHouse
- **Then:** `rr-signoz` và `rr-clickhouse` chạy ổn định, nhận OTLP tại `4317/4318`.

### Scenario 2: Nâng cao (Storage & Broker)
- **Then:** Redpanda, Cassandra và Elasticsearch khởi động thành công, sẵn sàng cho các Sprint sau.
