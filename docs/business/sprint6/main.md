# Sprint 6: Observability & Traceability

**Goal:** Triển khai hệ thống truy xuất nguồn gốc (Traceability) sử dụng CQRS và nâng cao khả năng quan sát hệ thống (Observability) qua OTel và Cassandra.

---

## 🎯 Sprint Goal

> **Hệ thống cho phép truy xuất nguồn gốc mẻ cà phê từ QR code và giám sát mọi giao dịch thông qua Distributed Tracing và Audit Log.**

---

## 📋 Roadmap & Flows

Sprint này tập trung vào tính minh bạch của dữ liệu và hệ thống:

### 1. Traceability Service (CQRS)
- Tổng hợp dữ liệu từ các service Farm, Processing, Warehouse, Retail.
- Xây dựng Read Model tối ưu trên Elasticsearch để truy vấn nguồn gốc cực nhanh.

### 2. Audit Service & Cassandra
- Lưu trữ mọi sự kiện thô (Raw Events) từ Kafka vào Apache Cassandra.
- Triển khai Hash Chaining để đảm bảo tính bất biến và chống giả mạo của Audit Log.

### 3. Full Instrumentation (OTel)
- Tích hợp OpenTelemetry SDK vào toàn bộ các microservices.
- Export Traces tới Jaeger/SigNoz để phân tích hiệu năng và debug Saga flow.

### 4. System Monitoring
- Triển khai Prometheus và Grafana.
- Xây dựng Dashboards theo dõi sức khỏe hệ thống (CPU, RAM, Kafka Lag, Error Rates).

---

## 🎫 Tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-27](./RR-27.md) | Traceability Service & CQRS Read Model | 🕒 To Do |
| [RR-35](./RR-35.md) | Elasticsearch Integration & Index Design | 🕒 To Do |
| [RR-28](./RR-28.md) | Audit Service & Cassandra Storage | 🕒 To Do |
| [RR-29](./RR-29.md) | Full OpenTelemetry & Jaeger Setup | 🕒 To Do |
| [RR-30](./RR-30.md) | Prometheus & Grafana Monitoring | 🕒 To Do |
