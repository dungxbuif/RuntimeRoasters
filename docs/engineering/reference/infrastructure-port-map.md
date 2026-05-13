# Runtime Roasters — Unified Port Mapping (Final)

Tài liệu này là **Nguồn sự thật duy nhất (Single Source of Truth)** về quy hoạch cổng (ports) cho toàn bộ hệ thống. Mọi thay đổi về hạ tầng hoặc service mới phải cập nhật tại đây trước.

---

## 1. Hạ tầng & Công cụ (Infrastructure & Tools)

| Thành phần | Port (Host) | Ghi chú |
| :--- | :--- | :--- |
| **PostgreSQL** | 54321 | Direct access. |
| **PgBouncer** | 6432 | Connection Pooler (Sprint 12). |
| **Valkey** | 6379 | Caching, Idempotency & Blacklist. |
| **Kafka (Broker)** | 9094 | External access. |
| **Kafka UI** | 8090 | Quản lý Kafka (Tránh trùng với 808x). |
| **Zookeeper** | 2181 | |

---

## 2. Bảo mật & Gateway (Security & Ingress)

| Thành phần | Port (Host) | Ghi chú |
| :--- | :--- | :--- |
| **KrakenD Gateway** | 8081 | Cửa ngõ API duy nhất. |
| **Kratos Public** | 4433 | Login, Register, Profile. |
| **Hydra Public** | 4444 | OAuth2 / OIDC endpoints. |
| **Hydra Admin** | 4445 | Client management. |
| **Identity Proxy** | 4434 | Nginx bọc Hydra Public (Bảo vệ JWKS). |

---

## 3. Microservices (Go Services)

Quy tắc: 
- **HTTP**: 808x (8080, 8082, 8083...)
- **gRPC**: 5005x (50051, 50052, 50053...)

| Service | HTTP | gRPC | Sprint | Ghi chú |
| :--- | :--- | :--- | :--- | :--- |
| **client-app** | 3000 | — | S1 | Frontend (React/Next). |
| **demo-service** | 8080 | 50051 | S1 | Canonical Template. |
| **auth-service** | 8082 | 50052 | S4 | Policy & Identity Proxy. |
| **farm-service** | 8083 | 50053 | S4 | Farm Management. |
| **process-service** | 8084 | 50054 | S5 | Roastery Processing. |
| **warehouse-service**| 8085 | 50055 | S6 | Inventory & Saga Participant. |
| **retail-service** | 8086 | 50056 | S7 | Store POS & Saga Orchestrator. |
| **logistics-service**| 8087 | 50057 | S8 | Real-time GPS & Trips. |
| **payment-service** | 8088 | 50058 | S9 | Stripe/VNPay Integration. |
| **trace-service** | 8089 | 50059 | S10 | CQRS Traceability Engine. |
| **audit-service** | 8091 | 50061 | S11 | Cassandra Audit Logs. |
| **webhook-service** | 8092 | 50062 | S9 | HMAC Validation Ingress. |

---

## 4. Công cụ Quan sát (Observability)

| Công cụ | Port |
| :--- | :--- |
| **SigNoz Dash** | 3301 |
| **Prometheus** | 9090 |
| **Grafana** | 3001 |
| **Jaeger UI** | 16686 |

---
*Cập nhật lần cuối: 2026-05-13 bởi TechLead*
