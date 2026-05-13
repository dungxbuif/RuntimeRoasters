# Runtime Roasters - Farm-to-Cup Coffee Supply Chain Platform

[English](#english) | [Tiếng Việt](#tiếng-việt)

---

<a name="english"></a>

## English Version

**Runtime Roasters** is a farm-to-cup coffee supply chain and logistics platform built to demonstrate a production-grade **Microservices** architecture. It simulates the entire lifecycle of a coffee bean—from the farm to the retail cup—ensuring transparency, traceability, and reliability.

The system is developed 100% in `Golang`, adhering to `Domain-Driven Design` (`DDD`) principles and implementing sophisticated distributed patterns to ensure scalability and data integrity.

> **Design Philosophy — HA from Design, not HA from Patch:** High Availability is not a feature added after the system is built — it is a set of architectural constraints enforced from the first line of code. Every service is designed so that scaling from 1 to N instances is purely an infrastructure operation: add a container, done. No logic changes, no coordination code, no sticky sessions to worry about. See [`docs/architecture/README.md`](docs/architecture/README.md) for the full HA checklist.

---

### 🎓 What This Project Teaches

This project serves as a comprehensive showcase of industry-standard patterns and technologies in the `Go` ecosystem:

- **Microservices Architecture:** 9+ `Independent Services` in a `Monorepo`, each following `Clean Architecture` (Domain / UseCase / Infrastructure).
- **Event-Driven Communication:** Asynchronous coordination using `Apache Kafka`.
- **Saga Pattern** (`Choreography`): Managing multi-step flows (`Order` ➔ `Payment` ➔ `Warehouse` ➔ `Logistics`) with automatic compensation (`Rollbacks` + `Auto-Refund`).
- **CQRS** (`Command Query Responsibility Segregation`): Separating transactional writes (`PostgreSQL`) from high-speed traceability reads (`Elasticsearch`).
- **Transactional Outbox** & **Inbox:** Applied selectively to critical business flows to guarantee message delivery and prevent duplicate processing (`Idempotency`).
- **Dual Idempotency:** `Idempotency-Key` in `Valkey` (HTTP, sync) + `Inbox Pattern` in `PostgreSQL` (async) — two independent layers.
- **Payment Gateway Integration:** `Stripe` `PaymentIntents`, `Webhook` handling with `HMAC` signature validation, auto-refund `Compensating Actions`.
- **Webhook Ingress Gateway:** Dedicated `Webhook Service` validates all external data (`Stripe`, `VNPay`, `IoT`) with `HMAC` before publishing to `Kafka`.
- **Zero Trust Architecture:** `mTLS` for all internal `gRPC`, `Decentralized Authorization` via `Casbin` per-service.
- **Distributed Tracing:** `Trace-ID` injected at `Gateway`, propagated through `gRPC Context` and `Kafka Headers`, visualized in `Jaeger` Gantt charts.
- **Standardized API Response (RFC 7807):** Unified error handling using the `Problem Details` standard for consistent client-side processing.
- **Real-time Geo-Tracking:** Live GPS management using `Valkey` (high-performance `Valkey` alternative).

---

### 📖 The Business Story

In Vietnamese culture, coffee is a "thread" that connects society. However, there is often a gap in information between the farmer and the consumer. **Runtime Roasters** digitizes this entire "lifeline" (Mạch sống).

Each coffee bean is given a digital identity. From the smallest data point at the plantation, through processing stages at the factory, to nationwide logistics, and finally into the customer's cup, every step is recorded and verifiable.

**The main touchpoints are:**
1. **Upstream (Farm):** Farmers update cultivation area and declare newly harvested coffee batches.
2. **Midstream (Processing & Warehouse):** Factory managers receive raw beans, proceed to hulling, drying, roasting, and packaging into finished batches (Batch ID). Goods are stored at the central warehouse.
3. **Transportation (Logistics):** When a dispatch order is received, the driver receives the trip, transports goods from the factory to retail locations, and continuously updates GPS coordinates in real-time.
4. **Downstream (Retail):** Store managers track inventory at the point of sale, send replenishment requests, and receive goods from drivers.
5. **Traceability & Audit:** End users can scan QR codes to view the entire journey of the coffee cup. Administrators have tools for panoramic monitoring and audit logs to prevent fraud.

---

### 🧩 System Overview

| Service     | Responsibility                                                                                        | Key Patterns            |
| :---------- | :---------------------------------------------------------------------------------------------------- | :---------------------- |
| `Gateway`   | Entry point, Authentication, Rate Limiting                                                            | `API Gateway Pattern`   |
| `Webhook`   | Ingress gateway for external data (`Stripe`, `VNPay`, `IoT GPS`) — `HMAC` validate + publish to Kafka | `Inbox Pattern`, `HMAC` |
| `Farm`      | Farmer management, plantation tracking, and harvest coffee batches                                    | `Transactional Design`  |
| `Process`   | Mill/Roastery operations, processing raw beans into `Batch ID`s                                       | `Event-Driven`          |
| `Warehouse` | Inventory management, stock reservation (`Saga Participant`)                                          | `Transactional DB`      |
| `Logistics` | Route coordination, real-time `GPS` tracking via `Valkey`                                              | `Geo-spatial Tracking`  |
| `Retail`    | Store ordering, consumption tracking (`Saga Orchestrator`)                                            | `State Machine`         |
| `Payment`   | `Stripe`/`VNPay` integration via `Strategy Pattern`, `Refund` management                              | `Inbox Pattern`, `HMAC` |
| `Trace`     | Unified traceability engine aggregate data into `Elasticsearch`                                       | `CQRS` (`Read model`)   |
| `Audit`     | Immutable event logging stored in `Apache Cassandra`                                                  | `Event Sourcing` (Lite) |

---

### 💻 Technology Stack

#### Backend

- **Language:** `Go` (`Golang`) 1.25+ with **Go Workspaces**
- **Dependency Injection:** **Manual DI** (Composition Root tại `main.go`)
- **Codebase:** `Monorepo` — shared `api/` (`gRPC` Proto), `pkg/` (Middleware, DB Wrapper)
- **Service structure:** `Clean Architecture` (Domain / UseCase / Infrastructure)
- **Communication:** `gRPC` (Internal), `Gin`/`REST` (External)
- **Message Broker:** `Apache Kafka`
- **Payment:** `Stripe` (`PaymentIntents`, `Webhooks`, Test Mode) via `Strategy + Factory Pattern`
- **Config:** `viper` + `.env` — no hardcoded secrets
- **Observability:** `OpenTelemetry`, `Jaeger`, `Prometheus` & `Grafana`

#### Data Persistence

- **PostgreSQL:** `Source of Truth` (`ACID` transactions)
- **Elasticsearch:** Search engine for traceability and history
- **Apache Cassandra:** `Audit logs` and hash-chained event storage
- **Valkey:** `Distributed locking`, `caching`, and real-time `GPS` tracking

#### Security & IAM

- **Ory Kratos & Hydra:** `Identity` and `OAuth2/OIDC` management
- **Casbin:** Decentralized `policy-based authorization`
- **mTLS:** `Mutual TLS` for secure `gRPC` service-to-service communication

---

### 🏗️ Technical Deep Dive

#### Clean Architecture & Shared Framework

Every service follows `go-clean-arch v4` with 3 strict layers: **Domain**, **UseCase**, and **Infrastructure**. Common logic (Observability, Auth headers, RFC 7807 Errors, Health Checks) is consolidated in a shared `pkg/` library, reducing boilerplate across all 9+ services.

#### Saga (Choreography)

Coordinating a supply chain request involves multiple services. If a store requests a batch but the `Warehouse` fails to reserve stock, a `Compensating Event` is triggered — `Payment Service` auto-calls `Stripe Refund API`, ensuring no money gets stuck.

#### Zero Trust & High-Performance Auth (Centralized Identity & Distributed Validation)

The system implements **Centralized Identity** but **Distributed Validation**. **Ory Kratos** manages user identities, while **Ory Hydra** acts as the OIDC Provider to issue standard JWTs. To achieve maximum performance, the `API Gateway` acts purely as a proxy, forwarding requests with tokens directly to downstream services. Upon startup, internal services (e.g., Farm, Catalog) fetch Public Keys (JWKS) from the **Hydra JWKS endpoint**. They perform **in-memory JWT validation** locally for every request without making external network calls to the identity stack. After validating the token, each service executes its own decentralized authorization (via `Casbin`), ensuring a `Zero Trust` network with zero-latency authentication.

#### Dual Idempotency

Two independent layers protect against duplicate processing: (1) `Idempotency-Key` header stored in `Valkey` (24h TTL) for synchronous HTTP retries, (2) `Inbox Pattern` in `PostgreSQL` for async `Kafka` consumers and `Webhooks`.

#### Transactional Outbox

To prevent `dual-write` problems in critical flows, services write events to an `outbox` table _in the same transaction_ as business logic. A dedicated worker relays these to `Kafka`, guaranteeing at-least-once delivery even if the broker is temporarily down.

#### Webhook Ingress Gateway

A dedicated `Webhook Service` (separate from `API Gateway`) handles all inbound external data. It validates `HMAC signatures` (Stripe, VNPay), deduplicates via `Inbox Pattern`, then publishes normalized events to `Kafka`. Zero business logic lives here.

#### Distributed Tracing (OpenTelemetry)

Every request gets a unique `Trace-ID` at the `Gateway`. This ID is injected into `gRPC` metadata, `HTTP traceparent` headers, and `Kafka` message headers — propagating across the entire Saga. `Jaeger` renders a `Gantt chart` of the full transaction lifecycle.

---

### 🚀 Getting Started

#### Prerequisites

- `Go` 1.25+
- `Docker` & `Docker Compose`
- `Apache Kafka` / `Redpanda`

#### Installation

1. Clone the repository
2. Spin up infrastructure: `docker-compose up -d`
3. Run services: `go run ./cmd/...`

---

<a name="tiếng-việt"></a>

## Bản tiếng Việt

**Runtime Roasters** là một nền tảng quản lý chuỗi cung ứng và logistics mô phỏng vòng đời của hạt cà phê từ nông trại đến tách cà phê bán lẻ (Farm-to-Cup). Dự án được xây dựng như một bản showcase kỹ thuật cho kiến trúc `Microservices` hiện đại.

Hệ thống được phát triển 100% bằng `Golang`, tuân thủ nghiêm ngặt các nguyên tắc `Domain-Driven Design` (`DDD`) và áp dụng các mẫu thiết kế phân tán phức tạp để đảm bảo tính mở rộng, độ tin cậy và khả năng truy xuất nguồn gốc minh bạch.

> **Triết lý thiết kế — HA from Design, không phải HA from Patch:** Độ sẵn sàng cao (HA) không phải tính năng được bổ sung sau khi hệ thống đã chạy — đây là tập hợp các ràng buộc kiến trúc được áp đặt từ dòng code đầu tiên. Mỗi service được thiết kế để scale từ 1 lên N instances chỉ là thao tác hạ tầng thuần túy: thêm container, xong. Không thay đổi logic, không cần code điều phối, không lo sticky session. Xem [`docs/architecture/README.md`](docs/architecture/README.md) để biết đầy đủ HA checklist.

---

### 🎓 Dự án này mang lại kiến thức gì?

Đây là một ví dụ thực tế về việc áp dụng các `pattern` và công nghệ tiêu chuẩn công nghiệp trong hệ sinh thái `Go`:

- **Kiến trúc Microservices:** Hơn 9 `Independent Services` trong `Monorepo`, mỗi service theo `Clean Architecture` (Domain / UseCase / Infrastructure).
- **Giao tiếp hướng sự kiện (Event-Driven):** Phối hợp bất đồng bộ sử dụng `Apache Kafka`.
- **Saga Pattern** (`Choreography`): Quản lý luồng đa bước (`Order` ➔ `Payment` ➔ `Warehouse` ➔ `Logistics`) với hoàn tác tự động (`Rollbacks` + `Auto-Refund`).
- **CQRS:** Tách biệt luồng ghi (`PostgreSQL`) và đọc traoốc độ cao (`Elasticsearch`).
- **Transactional Outbox** & **Inbox:** Áp dụng có chọn lọc cho các luồng nghiệp vụ quan trọng để đảm bảo gửi tin nhắn và tránh xử lý lặp lại (`Idempotency`).
- **Tính Lũy Đẳng Kép (Dual Idempotency):** `Idempotency-Key` trong `Valkey` (HTTP sync) + `Inbox Pattern` trong `PostgreSQL` (async) — hai lớp bảo vệ độc lập.
- **Tích hợp Payment Gateway:** `Stripe` `PaymentIntents`, xử lý `Webhook` với `HMAC`, hoàn tiền tự động qua `Strategy + Factory Pattern`.
- **Webhook Ingress Gateway:** `Webhook Service` riêng biệt xác thực `HMAC` từ `Stripe`, `VNPay`, thiết bị `IoT` trước khi publish lên `Kafka`.
- **Kiến trúc Zero Trust:** `mTLS` cho `gRPC` nội bộ + `Casbin` phân quyền phi tập trung tại từng service.
- **Distributed Tracing:** `Trace-ID` phát sinh tại `Gateway`, lan truyền qua `gRPC Context` và `Kafka Headers`, hiển thị `Gantt Chart` trên `Jaeger`.
- **Phản hồi lỗi chuẩn hóa (RFC 7807):** Sử dụng chuẩn `Problem Details` để đồng nhất cách xử lý lỗi giữa Backend và các loại Client.
- **Theo dõi GPS thời gian thực:** Sử dụng `Valkey` (GEO commands, high-performance `Valkey` alternative).

---

### 📖 Câu chuyện Doanh nghiệp

Trong văn hóa Việt, cà phê không chỉ là thức uống mà là "sợi dây" kết nối xã hội. Tuy nhiên, sự đứt gãy thông tin giữa người nông dân và tách cà phê trên tay khách hàng vẫn còn hiện hữu. **Runtime Roasters** ra đời để số hóa toàn bộ "mạch sống" này.

Dự án mô phỏng một hệ sinh thái nơi mỗi hạt cà phê đều có một danh tính số. Từ những dữ liệu nhỏ nhất ở nông trường, đi qua các công đoạn xử lý tại nhà máy, đến các mạng lưới logistics xuyên suốt đất nước để cuối cùng hội tụ tại những tách cà phê thơm ngon.

---

### 🧩 Tổng quan Hệ thống

| Dịch vụ     | Trách nhiệm chính                                                                                            | Các Pattern áp dụng     |
| :---------- | :----------------------------------------------------------------------------------------------------------- | :---------------------- |
| `Gateway`   | Cửa ngõ, Xác thực, Giới hạn lưu lượng (`Rate Limiting`)                                                      | `API Gateway Pattern`   |
| `Webhook`   | Ingress gateway tiếp nhận dữ liệu bên ngoài (`Stripe`, `VNPay`, `IoT GPS`) — xác thực `HMAC` + publish Kafka | `Inbox Pattern`, `HMAC` |
| `Farm`      | Quản lý nông hộ, vườn cây và thu hoạch theo mẻ                                                               | `Transactional Design`  |
| `Process`   | Tiếp nhận hạt thô, bóc vỏ, rang xay và cấp `Batch ID`                                                        | `Event-Driven`          |
| `Warehouse` | Quản lý kho, giữ chỗ hàng (`Saga Participant`)                                                               | `Transactional DB`      |
| `Logistics` | Điều phối vận tải, tracking `GPS` qua `Valkey`                                                                | `Geo-spatial Tracking`  |
| `Retail`    | Cửa hàng đặt hàng, quản lý tiêu thụ (`Saga Orchestrator`)                                                    | `State Machine`         |
| `Payment`   | Tích hợp `Stripe`/`VNPay` qua `Strategy Pattern`, quản lý `Refund`                                           | `Inbox Pattern`, `HMAC` |
| `Trace`     | Tổng hợp hành trình hạt cà phê vào `Elasticsearch`                                                           | `CQRS` (`Read model`)   |
| `Audit`     | Lưu lịch sử bất biến và `Audit log` vào `Apache Cassandra`                                                   | `Event Sourcing` (Lite) |

---

### 💻 Ngăn xếp Công nghệ

#### Backend

- **Ngôn ngữ:** `Go` (`Golang`) 1.25+ với **Go Workspaces**
- **Dependency Injection:** **Manual DI** (Composition Root tại `main.go`)
- **Codebase:** `Monorepo` — chia sẻ `api/` (`gRPC` Proto), `pkg/` (Middleware, DB Wrapper)
- **Kiến trúc service:** `Clean Architecture` (Domain / UseCase / Infrastructure)
- **Giao tiếp:** `gRPC` (Nội bộ), `Gin`/`REST` (Bên ngoài)
- **Message Broker:** `Apache Kafka`
- **Thanh toán:** `Stripe` qua `Strategy + Factory Pattern` — dễ thêm `VNPay`, `MoMo`
- **Cấu hình:** `viper` + `.env` — không hardcode secrets
- **Giám sát:** `OpenTelemetry`, `Jaeger`, `Prometheus` & `Grafana`

#### Cơ sở dữ liệu

- **PostgreSQL:** `Source of Truth` (Giao dịch `ACID`)
- **Elasticsearch:** Công cụ tìm kiếm truy xuất nguồn gốc
- **Apache Cassandra:** Lưu trữ `Audit log` và `Hash-chained event` storage
- **Valkey:** `Distributed locking`, `Cache` và quản lý `GPS` thời gian thực

#### Bảo mật & IAM

- **Ory Kratos:** Quản lý danh tính (`Identity Management`)
- **Casbin:** Phân quyền dựa trên chính sách phi tập trung (`Policy-based Authorization`)
- **mTLS:** `Mutual TLS` cho giao tiếp `gRPC` an toàn giữa các dịch vụ

---

### 🏗️ Đi sâu vào Kỹ thuật

#### Kiến trúc sạch & Shared Framework

Mỗi service theo `go-clean-arch v4` với 3 lớp: **Domain**, **UseCase**, **Infrastructure**. Các logic dùng chung (Observability, Auth header, RFC 7807, Health Checks) được gom vào thư viện `pkg/` dùng chung, giúp giảm thiểu code lặp cho tất cả 9+ services.

#### Saga (Choreography)

Điều phối yêu cầu cung ứng liên quan đến nhiều dịch vụ. Nếu `Warehouse` hết hàng sau khi đã thanh toán, `Payment Service` tự động gọi `Stripe Refund API` — tiền không bao giờ bị "treo".

#### Zero Trust & High-Performance Auth (Bảo mật tập trung & Xác thực phân tán)

Hệ thống áp dụng mô hình **Bảo mật tập trung (Centralized Identity)** nhưng **Xác thực phân tán (Distributed Validation)**. Identity Provider (`Ory Kratos`) cấp phát JWT và quản lý Public Keys (JWKS). API Gateway đóng vai trò là "Người dẫn đường" (Proxy), chuyển tiếp (Forward) trực tiếp các yêu cầu kèm Token từ bên ngoài vào. Khi khởi động, các Microservices (như Catalog, Farm) tải về JWKS từ IDT. Khi nhận request, chúng tự động giải mã và xác thực Token ngay tại bộ nhớ (in-memory validation) bằng Public Keys đã cache. Nhờ vậy, hệ thống loại bỏ hoàn toàn độ trễ gọi mạng (network call) sang IDT trên mỗi request, giúp tăng hiệu năng tối đa trong khi vẫn duy trì kiến trúc `Zero Trust`.

#### Tính Lũy Đẳng Kép (Dual Idempotency)

Hai lớp bảo vệ độc lập: (1) `Idempotency-Key` trong `Valkey` (TTL 24h) cho HTTP retry, (2) `Inbox Pattern` trong `PostgreSQL` cho `Kafka` consumer và `Webhook`.

#### Transactional Outbox

Để tránh lỗi `dual-write` trong các luồng quan trọng, service ghi event vào bảng `outbox` _cùng một Transaction_ với logic nghiệp vụ. Worker riêng đẩy lên `Kafka`, đảm bảo event không mất dù `Broker` tạm thời chết.

#### Webhook Ingress Gateway

`Webhook Service` riêng biệt (khác `API Gateway`) tiếp nhận toàn bộ dữ liệu đầu vào từ `Stripe`, `VNPay`, `IoT GPS`. Xác thực `HMAC`, dedup bằng `Inbox Pattern`, rồi publish lên `Kafka`. Không chứa business logic.

#### Distributed Tracing (OpenTelemetry)

Mỗi yêu cầu được gán `Trace-ID` tại `Gateway`. ID này lan truyền qua `gRPC` metadata, `HTTP traceparent` header và `Kafka` message headers, cho phép xem toàn bộ vòng đời giao dịch Saga trên biểu đồ `Gantt` của `Jaeger`.

---

### 🚀 Bắt đầu

#### Yêu cầu hệ thống

- `Go` 1.25+
- `Docker` & `Docker Compose`
- `Apache Kafka` / `Redpanda`

#### Cài đặt

1. Clone repository
2. Chạy hạ tầng: `docker-compose up -d`
3. Chạy các services: `go run ./cmd/...`

---

## 🛡️ Giấy phép

MIT - Được tạo bởi cộng đồng **RuntimeRoasters**.
