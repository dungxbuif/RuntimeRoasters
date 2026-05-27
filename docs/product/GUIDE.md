# 📘 Runtime Roasters Developer Guide

This guide provides a comprehensive technical breakdown of the **Runtime Roasters** platform. It serves as the primary technical reference for developers, covering architectural patterns, coding standards, and operational details.

---

## 🏗️ 1. Architectural Philosophy

Runtime Roasters is built as a highly available, event-driven microservices ecosystem. The system is designed so that scaling from 1 to N instances is purely an infrastructure operation with zero logic changes.

### Core Principles
- **Clean Architecture:** Strict separation of concerns (Domain ➔ UseCase ➔ Infrastructure ➔ Delivery).
- **Abstraction First:** All infrastructure (DB, Queue, Cache) is abstracted via Interfaces.
- **Calculated Consistency:** **Transactional Outbox** ensures critical events are never lost; **Inbox Pattern** ensures idempotency.
- **Dual Idempotency:** Two layers of protection: `Idempotency-Key` (HTTP) and `Transactional Inbox` (Kafka).
- **Observability by Design:** Distributed tracing (W3C) is propagated across Gateway, gRPC, and Kafka.
- **Zero-Trust Security:** Internal model with decentralized authorization (Casbin).

---

## 🌊 2. System Data Flows

### 🔐 Identity & SSO
1.  **Identity:** Managed by **Ory Kratos**. Handles user profiles, registration, and login.
2.  **Authentication:** Managed by **Ory Hydra**. Issues OIDC/OAuth2 tokens.
3.  **Gatekeeper:** KrakenD API Gateway validates JWT signatures and `scope` claims.
4.  **Authorization:** Each service runs a **Casbin Enforcer**. Policies are synced from `auth-service` via Kafka (`auth.policy.changed`) and gRPC snapshots.

### 📦 Order SAGA (Choreography)
1.  **Retail Service:** Creates Order (Pending) ➔ Publishes `OrderCreated` event.
2.  **Warehouse Service:** Receives `OrderCreated` ➔ Reserves Stock ➔ Publishes `StockReserved` (or `StockReservationFailed`).
3.  **Payment Service:** Receives `StockReserved` ➔ Processes Payment ➔ Publishes `PaymentCompleted`.
4.  **Retail Service:** Receives `PaymentCompleted` ➔ Updates Order to `Confirmed`.
5.  *Compensation:* If any step fails, services publish compensation events to roll back previous actions.

---

## 🛠️ 3. Development Workflow

### Project Structure (src/apps/)
Each service follows this layout:
- `cmd/`: Process entry point (Main).
- `internal/`:
    - `domain/`: Business entities and interfaces.
    - `usecase/`: Business logic.
    - `infrastructure/`: Implementations (DB, Kafka, Valkey).
    - `delivery/`: Transport handlers (gRPC, Gin).
- `config/`: Configuration mapping.

### 📦 3.1. Internal Framework (`src/pkg/`)
All services leverage a shared core framework to ensure consistency and reduce boilerplate:
- **`base`**: The service orchestrator. Handles gRPC/HTTP server lifecycle, graceful shutdown, and cross-cutting auth/security middlewares.
- **`database`**: GORM-based abstraction for PostgreSQL with built-in support for Transactional Outbox.
- **`kafka`**: High-level wrapper for Segmentio/Kafka-Go, implementing the Outbox/Inbox patterns.
- **`valkey`**: Client for Valkey (Redis alternative) used for distributed locking and real-time tracking.
- **`errs`**: Centralized error handling using RFC 9457 (Problem Details).
- **`telemetry`**: OpenTelemetry integration for tracing and metrics.
- **`logger`**: Structured JSON logging via Uber-Zap with Trace-ID propagation.

### 🧪 3.2. Testing Strategy
We employ a multi-layered testing strategy to ensure system reliability:
- **Unit Tests**: Focus on `domain` and `usecase` logic. Mocking is used for all external dependencies (DB, Kafka).
- **Integration Tests**: Verify the interaction between `infrastructure` implementations and real (Docker-based) dependencies.
- **E2E Tests**: (Future) Simulated supply chain runs across the full stack.
- **Contract Testing**: Protobuf definitions serve as the strict contract between services.

### Key Tools
- **Go Workspaces (`go.work`):** We use a Go workspace to manage multiple internal modules (e.g., `src` containing apps and pkg). This allows for seamless cross-module development and dependency resolution.
- **Taskfile:** Automated tasks for infra, proto generation, and seeding.
- **Air:** Hot-reload for Go development.
- **Buf:** Modern Protobuf management and code generation.
- **KrakenD:** Powerful API gateway for orchestration and security.

### 🔄 3.3. Environment Reset & Seeding
To maintain a clean development state or prepare for a fresh demo:
1.  **Full Reset:** Run `task env:reset`. This stops all containers, wipes all volumes (DB/Kafka/Valkey), restarts infra, and runs schema migrations.
2.  **Infrastructure Seed:** Automatically handled by `env:reset`, this creates the default `ADMIN` user and OAuth2 clients.
3.  **Business Data Seeding:** Unlike infrastructure, business data (Farms, Stores, Drivers) is **not** auto-seeded. After logging in as `ADMIN`, a bootstrap modal will appear on the Dashboard to trigger manual seeding via protected microservice endpoints.

### Coding Standards
- **DRY:** Use shared logic in `src/pkg/` (logger, database, telemetry).
... Applied fuzzy match at line 118-132.
- **Type-Safety:** All internal communication must use gRPC/Protobuf.
- **Explicit over Implicit:** Manual dependency injection is preferred over magic containers.

---

## 🧭 4. Operational Reference

For the current install/start/seed/playground runbook, use:

- [Demo Setup Runbook](./DEMO_SETUP_RUNBOOK.md)
- [Master Technical Specification](./TECH.md)

### Business/BA Docs

For BA review, product ownership, role responsibilities, and UI-facing business flows, use the **[Master Business Doc](./domain/README.md)**.

### Local Service URLs
- **Client App:** `http://localhost:3000`
- **KrakenD (Gateway):** `http://localhost:8081`
- **Auth Service:** HTTP `http://localhost:8082`, gRPC `localhost:50052`
- **Farm Service:** HTTP `http://localhost:8083`, gRPC `localhost:50053`
- **Retail Service:** HTTP `http://localhost:8084`, gRPC `localhost:50054`
- **Logistics Service:** HTTP `http://localhost:8085`, gRPC `localhost:50055`
- **Payment Service:** HTTP `http://localhost:8086`, gRPC `localhost:50056`
- **Trace Service:** HTTP `http://localhost:8087`, gRPC `localhost:50057`
- **Audit Service:** HTTP `http://localhost:8088`, gRPC `localhost:50058`
- **Warehouse Service:** HTTP `http://localhost:8089`, gRPC `localhost:50059`
- **Socket Service:** HTTP `http://localhost:8091`, gRPC `localhost:50060`
- **Kafka UI:** `http://localhost:8090`
- **Kibana (ES):** `http://localhost:5601`
- **Elasticsearch:** `http://localhost:9200`
- **SigNoz UI:** `http://localhost:3301`
- **OTLP:** gRPC `localhost:4317`, HTTP `localhost:4318`
- **Postgres:** `localhost:54321`
- **Cassandra:** `localhost:9042`
- **Kratos Public:** `http://localhost:4433`
- **Hydra Public:** `http://localhost:4444`

### Demo Stripe Webhook Flow

- `GET /v1/payments/demo/stripe-webhook-key` returns the seeded demo Stripe
  signing key for authenticated roles.
- `POST /v1/webhooks/stripe` is the canonical authenticated webhook endpoint.
  The request must include `Stripe-Signature: t=<unix>,v1=<hmac>`.
- Finance UI Pass emits `payment.completed`; Fail emits `payment.failed`.
- Warehouse reservation starts after `payment.completed`, not immediately after
  order creation.
- Retail order completion waits for `logistics.driver.returned_to_base`.

### Socket + Topology Flow

- Public root topology pulls config/history from trace-service:
  - `GET /v1/traces/public/topology/config`
  - `GET /v1/traces/public/topology/history`
- Public realtime connects to socket-service WebSocket:
  - `GET /v1/realtime/public/topology/ws?flow_id=...`
- Private dashboard realtime uses:
  - `GET /v1/realtime/stream?scope=dashboard&flow_id=...`
- `socket-service` is DB-free. It uses Kafka for events and Valkey TTL keys for
  socket session, presence, and reconnect smoothing.
- Generate internal service keys with `scripts/socket-keygen.sh <service-name>`.

### Kafka Topics
- `auth.policy.changed`: Broadcasts Casbin policy updates.
- `auth.user.events`: User lifecycle events.
- `farm.harvest.events`: Harvest records from plantations.
- `order.saga.events`: Main topic for order choreography.

---

## 🗺️ 5. Detailed Guides & Decision Records

For a deeper dive into specific system behaviors and the reasoning behind them:
- 🏛️ **[Master Technical Specification](./TECH.md)**: Architecture, service catalog, contracts, storage, and infrastructure.
- 🎨 **[UI/UX Design Notes](./ui-ux/DESIGN.md)**: Next.js patterns and dashboard UI/UX standards.
- 📜 **[Master Business Specification](./domain/README.md)**: The single source of truth for business flows, roles, and domain rules.
- 📋 **[Product/System Specification](./SPEC.md)**: The broader technical feature specs and history.
- 📅 **[Sprint Roadmap](../stories/ROADMAP.md)**: Development phases and task breakdown.
- 📜 **[Architecture Decisions (ADRs)](../decisions/)**: History of critical technical choices.

---
*This document is the "Living Constitution" of Runtime Roasters. Keep it updated as the architecture evolves.*
# 🛠️ Runtime Roasters Engineering Log

This document tracks technical challenges, bug fixes, and significant implementation milestones. It serves as a historical record of the system's evolution.

---

## 🐞 Bugs Found & Fixed

### Auth & Security
- **Login Accept 404:** Fixed an issue where `POST /v1/auth/login/accept` returned 404 because Gin's `NoRoute` handler interfered with the gRPC gateway.
- **OIDC Callback Loop:** Fixed the application landing page to properly capture the access token, strip it from the URL, and redirect to the dashboard.
- **Casbin Method Mapping:** Resolved a 403 error in `farm-service` by correctly mapping gRPC methods to Casbin actions (Get/List ➔ `read`, Delete ➔ `delete`, others ➔ `write`).
- **JWKS Protection:** Implemented an Nginx proxy (`rr-identity`) to protect Hydra's JWKS endpoint using an internal secret, ensuring only authorized services can fetch public keys.

### Service Interactions
- **Auth Polling Fallback:** Added a background bootstrap mechanism in `farm-service` to resiliently fetch Casbin snapshots if `auth-service` is not ready at startup.
- **Transactional Consistency:** Fixed a race condition in the Outbox worker where events could be published twice; added database-level locking during the "Mark as Processed" phase.
- **Elasticsearch Mapping:** Updated the `trace-service` read model to use correct keyword mappings for UUID fields, enabling accurate filtering by `order_id`.

---

## 🚀 Key Milestones

### Sprint 6-10: The SAGA Flow
Successfully connected the end-to-end coffee order flow:
1.  **Retail Service:** Initiates Order.
2.  **Payment Service:** Creates a pending Stripe/VNPay intent.
3.  **Finance UI:** Posts the signed demo Stripe webhook Pass/Fail result.
4.  **Warehouse Service:** Reserves coffee bean stock after payment success.
5.  **Logistics Service:** Assigns a driver and calculates real-time route.
6.  **Trace Service:** Fills the Elasticsearch read model for the "Track My Order" feature.
7.  **Audit Service:** Writes immutable logs to Cassandra for compliance.

### Documentation & Knowledge Base (Current)
Executed a comprehensive overhaul of the project's documentation to move from "Service-centric" to "Domain-centric" knowledge:
- Created the **Master System Architecture Specification** as the definitive technical source.
- Compiled the **Business & Domain Specification** (handover-ready) to define core supply chain rules (Roasting loss, FIFO, Batching).
- Consolidated all architectural flows into a single visual reference.
- Synchronized the `.planning/codebase/` map with the current implementation state.

### Multi-Region Design
Implemented W3C Trace-ID propagation across all transport layers (HTTP ➔ gRPC ➔ Kafka), allowing a single request to be tracked as it hops through the entire ecosystem.

---
*Last Updated: May 2026*
# master_demo_data.md — Runtime Roasters Master Demo Dataset

Tài liệu này tổng hợp toàn bộ dữ liệu mẫu (Master Data) của hệ thống Runtime Roasters, phục vụ cho mục đích phát triển, kiểm thử (UAT) và trình diễn hệ thống.

---

## 📋 1. Danh sách Người dùng (Identities)

Tất cả mật khẩu mặc định: `Hello@123`

| Tên người dùng | Email (Username) | Role | Mục đích Demo |
| :--- | :--- | :--- | :--- |
| **System Admin** | `admin@runtimeroasters.com` | `ADMIN` | Quản trị toàn hệ thống, tạo user. |
| **Agri Admin** | `agri.admin@runtimeroasters.com` | `FARM_ADMIN` | Quản lý danh mục vùng trồng, gán sở hữu. |
| **Manager Cầu Đất** | `manager.caudat@runtimeroasters.com` | `FARM_MANAGER` | Quản lý vận hành tại vùng Đà Lạt. |
| **Manager Buôn Ma Thuột**| `manager.bmt@runtimeroasters.com` | `FARM_MANAGER` | Quản lý vận hành tại Đắk Lắk. |
| **Farm Manager K'Ho** | `manager.kho@runtimeroasters.com` | `FARM_MANAGER` | Ghi chép nhật ký thu hoạch thực địa. |
| **Roast Master** | `processor@runtimeroasters.com` | `PROCESSOR` | Tiếp nhận cà phê nhân và chế biến. |
| **Driver Alpha** | `driver@runtimeroasters.com` | `DRIVER` | Vận chuyển hàng hóa giữa các điểm. |

---

## ☕ 2. Danh sách Nông trại (Farms - Famous in Vietnam)

Dữ liệu này được thiết kế để khớp với các vùng trồng nổi tiếng thực tế tại Việt Nam.

| Tên Nông trại | Địa điểm (Location) | Diện tích | Giống Cà phê | Chủ sở hữu (Owner) |
| :--- | :--- | :--- | :--- | :--- |
| **K'Ho Coffee Farm** | Lạc Dương, Đà Lạt | 15.5 ha | Arabica Heirloom | `manager.caudat@...` |
| **Cau Dat Arabica** | Cầu Đất, Đà Lạt | 45.0 ha | Arabica Typica | `manager.caudat@...` |
| **Son Pacamara Farm** | Trạm Hành, Đà Lạt | 12.0 ha | Arabica Pacamara | `manager.caudat@...` |
| **Aeroco Coffee** | Ea Kao, Buôn Ma Thuột | 20.0 ha | Specialty Robusta | `manager.bmt@...` |
| **Trung Nguyen Village**| Tân Lợi, Buôn Ma Thuột | 5.0 ha | Robusta | `manager.bmt@...` |
| **Chư Sê Estate** | Chư Sê, Gia Lai | 30.0 ha | Robusta | `agri.admin@...` |

---

## 📦 3. Nhật ký Thu hoạch Mẫu (Harvests)

| Farm | Ngày thu hoạch | Sản lượng (kg) | Trạng thái |
| :--- | :--- | :--- | :--- |
| K'Ho Coffee Farm | 2026-05-10 | 500.00 | `NEW` |
| K'Ho Coffee Farm | 2026-05-11 | 350.00 | `PROCESSING` |
| Aeroco Coffee | 2026-05-09 | 1200.00 | `SHIPPED` |
| Cau Dat Arabica | 2026-05-12 | 800.00 | `NEW` |

---

## 🛠️ 4. Technical JSON (Dành cho Seeding/API Test)

### 4.1 JSON Tạo Identity (Kratos)
```json
{
  "traits": {
    "email": "manager.caudat@runtimeroasters.com",
    "name": "Manager Cầu Đất",
    "role": "FARM_MANAGER"
  },
  "credentials": { "password": { "config": { "password": "Hello@123" } } }
}
```

### 4.2 JSON Tạo Nông trại (Farm Service API)
```json
{
  "name": "K'Ho Coffee Farm",
  "location": "Lạc Dương, Đà Lạt",
  "area": 15.5,
  "coffee_type": "Arabica Heirloom",
  "owner_id": "[ID_CỦA_MANAGER_CAU_DAT]"
}
```

---

## 📖 5. Ghi chú Bảo mật & Vận hành
1. **Ownership**: Nông trại chỉ hiển thị cho người dùng có `id` khớp với `owner_id` (trừ `ADMIN` và `FARM_ADMIN` có quyền view-all).
2. **Role Mapping**: Hệ thống tự động đồng bộ Role từ Kratos sang Casbin qua cơ chế **Bootstrapping Sync** của Auth Service.
3. **Validation**: Khi tạo Farm mới, `area` phải > 0 và `location` nên chọn trong danh sách các tỉnh Tây Nguyên.
# Runtime Roasters — Demo Identities & Personas

Tài liệu này liệt kê danh sách các tài khoản người dùng mẫu (Demo Users) được chuẩn hóa để phục vụ quá trình kiểm thử (UAT) và trình diễn hệ thống (Showcase).

---

## 📋 1. Danh sách User Demo (Standardized Roles)

Mọi tài khoản dưới đây đều được gán Role viết hoa (UPPERCASE) theo đúng chuẩn hệ thống mới.

| Nhóm chức năng | Tên người dùng | Email (Username) | Password mẫu | Role (Mã hệ thống) | `store_ids` | `warehouse_ids` |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Quản trị hệ thống** | System Administrator | `admin@runtimeroasters.com` | `Hello@123` | **`ADMIN`** | `[]` | `[]` |
| **Quản trị nông nghiệp**| Agri Supply Admin | `agri.admin@runtimeroasters.com` | `Hello@123` | **`FARM_ADMIN`** | `[]` | `[]` |
| **Quản lý Nông trại** | Manager Sơn La | `manager.sonla@runtimeroasters.com` | `Hello@123` | **`FARM_MANAGER`** | `[]` | `[]` |
| **Quản lý Nông trại** | Manager Cầu Đất | `manager.caudat@runtimeroasters.com` | `Hello@123` | **`FARM_MANAGER`** | `[]` | `[]` |
| **Warehouse** | Warehouse Manager Hanoi | `warehouse.hn@runtimeroasters.com` | `Hello@123` | **`WAREHOUSE_MGR`** | `[]` | `["WAREHOUSE-HN-001"]` |
| **Nhà máy chế biến** | Roast Master | `processor@runtimeroasters.com` | `Hello@123` | **`PROCESSOR`** | `[]` | `["WAREHOUSE-HN-001"]` |
| **Vận tải (Logistics)** | Driver Alpha | `driver@runtimeroasters.com` | `Hello@123` | **`DRIVER`** | `[]` | `[]` |
| **Cửa hàng bán lẻ** | Mgr Hanoi Hoan Kiem | `mgr.hn.hoankiem@runtimeroasters.com` | `Hello@123` | **`STORE_MGR`** | `["11111111-1111-1111-1111-111111111101"]` | `[]` |
| **Cửa hàng bán lẻ** | Mgr HCM Dist 1 | `mgr.hcm.q1@runtimeroasters.com` | `Hello@123` | **`STORE_MGR`** | `["11111111-1111-1111-1111-111111111103"]` | `[]` |
| **Cửa hàng bán lẻ** | Mgr Danang Hai Chau | `mgr.dn.haichau@runtimeroasters.com` | `Hello@123` | **`STORE_MGR`** | `["11111111-1111-1111-1111-111111111105"]` | `[]` |


---

## 🎭 2. Kịch bản Demo Gợi ý (The Storyline)

Để thể hiện được toàn bộ sức mạnh của kiến trúc Microservices và Phân quyền, PO có thể thực hiện theo luồng sau:

### Phase 1: Onboarding (Quyền ADMIN)
- **Hành động**: Đăng nhập bằng `admin@...`.
- **Thực hiện**: Vào `/admin/users`, tạo các tài khoản `FARM_ADMIN` và `FARM_MANAGER`.
- **Giá trị**: Chứng minh tính năng **Admin-only Creation** thông qua Auth Proxy.

### Phase 2: Tài nguyên & Gán quyền (Quyền FARM_ADMIN)
- **Hành động**: Đăng nhập bằng `agri.admin@...`.
- **Thực hiện**: Vào `/admin/farms`, tạo mới "Nông trại Cầu Đất" và gán Manager là `manager.caudat@...`.
- **Giá trị**: Chứng minh khả năng điều phối tài nguyên hệ thống.

### Phase 3: Vận hành cục bộ (Quyền FARM_MANAGER)
- **Hành động**: Đăng nhập bằng `manager.caudat@...`.
- **Thực hiện**: Vào **Farm Origin** dashboard.
- **Giá trị**: Chứng minh **Data Scoping (ABAC)** — Manager này tuyệt đối không thấy nông trại của Manager kia.

### Phase 4: Farm -> Warehouse Pickup Demo
- **Màn A**: Đăng nhập `manager.caudat@runtimeroasters.com`, tạo harvest từ Farm dashboard.
- **Màn B**: Đăng nhập Warehouse Manager khi account seed có sẵn, mở Warehouse dashboard để xem pickup request và dispatch driver.
- **Màn C**: Đăng nhập `driver@runtimeroasters.com`, mở Driver/Logistics shipment screen, bấm Start để chạy route simulation, sau đó confirm pickup và return.
- **Màn quan sát**: Mở Logistics map, Trace page, hoặc public `ArchitectureTopology` để xem realtime update.
- **Giá trị**: Chứng minh physical logistics backbone, role-scoped UI, socket/SSE update, và traceability.

### Phase 5: Warehouse -> Retail Delivery Demo
- **Màn A**: Đăng nhập `mgr.hn.hoankiem@runtimeroasters.com` hoặc store manager tương ứng, tạo paid order.
- **Màn B**: Warehouse Manager dispatch retail delivery sau khi stock reserved.
- **Màn C**: `driver@runtimeroasters.com` bấm Start để chạy delivery route và confirm delivery.
- **Màn quan sát**: Retail dashboard xem incoming delivery status; Logistics map và Trace page xem realtime journey.
- **Giá trị**: Chứng minh order SAGA, inventory reservation, logistics delivery, và scoped store visibility.

### Phase 6: Public QR Trace Demo
- **Màn public**: Không đăng nhập, mở public trace showcase page.
- **Thực hiện**: Bấm nút generate/show demo QR codes.
- **Kết quả**: Client hiển thị danh sách QR cho các sản phẩm seed sẵn.
- **Thực hiện**: Click hoặc scan một QR.
- **Kết quả**: Public trace page hiển thị trace thật từ dữ liệu đã chuẩn bị: farm, pickup, warehouse intake, processing, paid order, delivery, driver return.
- **Giá trị**: Chứng minh traceability public mà không expose dữ liệu nhạy cảm.

---

## 🛠️ 3. Lưu ý Kỹ thuật
- **OIDC Flow**: Hệ thống sử dụng Ory Hydra để cấp phát token. Sau khi Login thành công tại Kratos, hệ thống tự động gọi `Accept Login` tới Hydra.
- **JWT Claims**: Role, email, org, store scope và warehouse scope của người dùng được nhúng trực tiếp vào JWT (`role`, `email`, `org_id`, `store_ids`, `warehouse_ids`) để các service verify in-memory.
- **Seeding**: Tài khoản `admin@runtimeroasters.com` được nạp tự động qua container `rr-seeder` khi khởi động infra.
- **Driver Simulation**: Main demo dùng Driver Client sau khi login bằng `DRIVER`; frontend replay route seed và bắn GPS/status về backend.
- **Realtime**: Private dashboard streams cần auth/role scope. Public root `ArchitectureTopology` có thể no-auth nếu chỉ hiển thị dữ liệu sanitized.
