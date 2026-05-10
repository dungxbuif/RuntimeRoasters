# RuntimeRoasters — System Architecture & Development Blueprint

Tài liệu này định nghĩa kiến trúc tổng thể, các tiêu chuẩn kỹ thuật và lộ trình thực thi cho dự án **RuntimeRoasters**.

---

## 1. Tầm nhìn Kiến trúc (Architectural Vision)

Mục tiêu là xây dựng một hệ thống Microservices **"Mạnh mẽ - Tin cậy - Quan sát được"**. Hệ thống không chỉ giải quyết bài toán nghiệp vụ Chuỗi cung ứng cà phê mà còn là một bản showcase về các mẫu thiết kế (Design Patterns) hiện đại nhất trong hệ sinh thái Go.

### Nguyên tắc cốt lõi (Core Principles):
- **Abstraction First:** Không phụ thuộc vào framework. Mọi thành phần hạ tầng (DB, Queue, Cache) đều được trừu tượng hóa qua Interface.
- **No Hardcoding:** Tuyệt đối không hardcode. Sử dụng cấu hình đa môi trường qua Viper & Environment Variables.
- **Database per Service:** Đảm bảo tính độc lập và khả năng mở rộng riêng biệt cho từng service.
- **Standardization:** Thống nhất format Log, Error (RFC 9457), và cách khởi tạo Service thông qua bộ khung `pkg/` dùng chung.
- **Sprint Discipline:** Sprint 1 = infrastructure only. Sprint 2+ = business logic only. Không setup infra trong sprint nghiệp vụ.

---

## 2. App Architecture

#### `apps/client-app/` (Next.js 15 — Unified UI)
- **Pattern:** Direct Pattern.
- **Role:** Dành cho tất cả người dùng (Nông dân, Nhà máy, Retailer) và Admin.
- **Flow:** Browser → KrakenD :8081 → Microservices.
- **Main Modules:**
    - `/app/farms`: Quản lý nông hộ (Sprint 2).
    - `/app/batches`: Theo dõi mẻ hàng (Sprint 2).
    - `/app/logistics`: Theo dõi vận chuyển (Sprint 3).
    - `/admin/*`: Các tính năng quản trị hệ thống.

---

## 3. Sơ đồ Hạ tầng Tổng thể

```mermaid
graph TB
    subgraph "External World"
        Browser([Browser / Mobile])
    end

    subgraph "Client App"
        CA[client-app :3000]
    end

    subgraph "API Gateway"
        GW[KrakenD :8081]
    end

    subgraph "Microservices"
        DS[demo-service :8080/:50051]
        AS[auth-service :8082/:50052]
        WS[warehouse-service Sprint 3+]
    end

    subgraph "Data Persistence"
        PG[(PostgreSQL :54321)]
        RD[(Redis :6379)]
    end

    subgraph "Message Broker"
        KF[Kafka :9092/:9094]
        ZK[Zookeeper :2181]
        KUI[Kafka UI :8082]
        KUI --- KF
        KF --- ZK
    end

    Browser -->|:3000| CA
    CA -->|REST :8081| GW
    GW -->|gRPC| DS
    GW -->|gRPC| AS
    DS --> PG
    DS --> RD
    AS --> PG
    DS -.->|Kafka| KF
    AS -.->|Kafka| KF
```

---

## 4. Port Map (Final — no conflicts)

| Service | Port | Note |
| :--- | :--- | :--- |
| client-app | 3000 | Unified UI (Business + Admin) |
| KrakenD | 8081 | API Gateway |
| demo-service HTTP | 8080 | |
| demo-service gRPC | 50051 | |
| auth-service HTTP | 8082 | (RR-12) |
| auth-service gRPC | 50052 | |
| PostgreSQL | 54321 | (Mapped from 5432) |
| Redis | 6379 | |
| Kafka | 9094 | (External) |
| Kafka UI | 8082 | |
| Zookeeper | 2181 | |

---

## 5. Thiết kế Framework nội bộ (`pkg/`)

| Package | Vai trò |
| :--- | :--- |
| `pkg/config` | Viper: load từ `.env`, YAML, env vars. `BaseConfig` embed vào mọi service config. |
| `pkg/logger` | Zap: JSON cho prod, console cho dev. `FromContext(ctx)` tự inject trace_id. |
| `pkg/base` | Application lifecycle: gRPC/HTTP server, health checks, graceful shutdown, OTel init. |
| `pkg/errs` | RFC 9457 Problem Details. `GinErrorHandler()` middleware. `SubProblem[]` support. |
| `pkg/database` | GORM wrapper: connection pool, otelsql auto-instrumentation, migration scaffold. |
| `pkg/redis` | go-redis wrapper: redisotel tracing + metrics. |
| `pkg/telemetry` | OTel TracerProvider + MeterProvider + W3C propagator. |

---

## 6. Các Mẫu Thiết kế Kỹ thuật Nâng cao

1. **Saga Pattern & Distributed Transactions:**
   - Phase 1 (Choreography): Outbox Pattern + Kafka. Compensating Actions tự động.
   - Phase 2 (Orchestration): Temporal.io.

2. **Transactional Outbox & Inbox:**
   - Outbox: Atomicity giữa DB write và event publish.
   - Inbox: Idempotency bằng `Message_ID`.

3. **CQRS & Real-time Traceability:**
   - Write: PostgreSQL (ACID transactions).
   - Read: Elasticsearch (full-text search).
   - Sync: Trace Service consume Kafka → upsert Elasticsearch.

4. **Standardized API Response (RFC 9457):**
   - `type`, `title`, `status`, `detail`, `instance`, `trace_id`, `errors[]`.
   - `application/problem+json` Content-Type.

5. **Resilient Centralized Authorization:** Một `auth-service` quản lý tập trung (Writer), các service khác thực thi in-memory (Reader) qua **Casbin v3** với cơ chế đồng bộ 3 trụ cột (gRPC Snapshot + Kafka Live + Polling).

6. **Hash Chaining (Cassandra):** Audit log chống gian lận nội bộ.

---

## 8. Bảo mật & Xác thực

- **Identity:** Ory Kratos (JWT cấp phát + JWKS).
- **Gateway:** KrakenD forward JWT và đính kèm Authorization header.
- **In-service auth:** Services tự validate JWT offline bằng JWKS cache.
- **Fine-grained AuthZ:** Sử dụng **Centralized Auth Service** và **Casbin v3**. Các service kéo quyền về RAM qua gRPC lúc khởi động và cập nhật thời gian thực qua Kafka.

---

## 9. Cấu trúc Monorepo

```text
RuntimeRoasters/
├── api/                    # Proto definitions + buf toolchain (Shared)
│   └── runtime/
│       └── farm/v1/
├── apps/
│   ├── demo-service/       # Primary service & boilerplate template
│   └── auth-service/       # Policy Manager (RR-12)
├── src/
│   └── pkg/                # Go Workspaces shared packages
│       ├── base/
│       ├── config/
│       ├── database/
│       ├── errs/
│       ├── logger/
│       ├── redis/
│       └── telemetry/
├── deployments/
│   ├── docker-compose.yaml
│   ├── krakend/
│   │   └── krakend.json
│   └── init-db.sql
└── docs/
    ├── architecture/
    └── business/
```

---

## 10. Sprint Roadmap

| Sprint | Goal | Key Deliverable (RR-x) |
| :--- | :--- | :--- |
| **Sprint 1** | Foundation & Identity | Identity Infra, Client Auth Flow (RR-1 to RR-10) |
| **Sprint 2** | Security Core | JWT Validation, Centralized AuthZ (RR-11 to RR-14) |
| **Sprint 3** | Farm Service Logic | Vertical Slice CRUD, DDD, Repository (RR-15 to RR-18) |
| **Sprint 4** | Distributed Systems | Saga Pattern, Retail, Warehouse (RR-19 to RR-22) |

**Nguyên tắc:** Sprint 1 hoàn tất toàn bộ infra. Sprint 2+ chỉ viết business logic — không setup thêm bất kỳ infrastructure nào.
