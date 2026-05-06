# Runtime Roasters — Product Requirements Document (PRD)

| Field            | Value                                                    |
| :--------------- | :------------------------------------------------------- |
| **Project**      | Runtime Roasters                                         |
| **Author**       | PM (AI-assisted)                                         |
| **Status**       | Draft v1.2                                               |
| **Last Updated** | 2026-04-16                                               |
| **Type**         | Technical Portfolio / Showcase                           |
| **Developer**    | Solo                                                     |
| **Deployment**   | Docker on Cloud VPS (HA-ready design)                    |

---

## 1. Executive Summary

**Runtime Roasters** là một dự án `Portfolio/Showcase` mô phỏng nền tảng quản lý chuỗi cung ứng cà phê từ nông trại đến tách bán lẻ (Farm-to-Cup). Mục tiêu chính là trình diễn năng lực thiết kế và triển khai kiến trúc `Microservices` phân tán ở production-grade, sử dụng 100% `Golang` cho `Backend` và `ReactJS` cho `Frontend`.

Dự án **không** hướng tới việc giải quyết bài toán thương mại thực tế, mà tập trung vào:
- Showcase các `Design Patterns` phức tạp (Saga, CQRS, Outbox, Event Sourcing).
- Chứng minh khả năng xử lý hệ thống phân tán (`Distributed Systems`).
- Trình diễn trực quan kiến trúc qua "Control Plane Visualization" `Dashboard`.

**Đối tượng đánh giá:** Nhà tuyển dụng, Technical Reviewer, TA (Technical Architect).

---

## 2. Goals & Non-Goals

### 2.1 Goals (Mục tiêu)
| #  | Goal                                                                                   | Đo lường bằng                                         |
| :- | :------------------------------------------------------------------------------------- | :----------------------------------------------------- |
| G1 | Showcase kiến trúc `Microservices` với 9+ services độc lập giao tiếp qua `Kafka`/`gRPC` | Tất cả services chạy và giao tiếp được trên `Docker Compose` |
| G2 | Triển khai thành công `Saga Pattern` (Choreography) với `Compensating Actions`          | Demo kịch bản Rollback khi `Warehouse` hết hàng + Auto-Refund via `Stripe` |
| G3 | Triển khai `Transactional Outbox` + `Inbox` (Idempotency)                               | Demo ngắt `Kafka`, data vẫn nhất quán sau khi khôi phục |
| G4 | Triển khai `CQRS` với `PostgreSQL` (Write) + `Elasticsearch` (Read)                     | Truy xuất nguồn gốc cà phê < 100ms                     |
| G5 | Dashboard trực quan hóa kiến trúc (Control Plane Visualization)                                            | Người xem thấy data flow real-time giữa các services   |
| G6 | Thiết kế sẵn cho `High Availability` (HA-ready)                                        | Kiến trúc có thể scale horizontal mà không refactor     |
| G7 | Tích hợp `Payment Gateway` (Stripe) với đầy đủ `Webhook Security` + `Idempotency`      | Demo thanh toán B2B, auto-refund khi Saga fail          |

### 2.2 Non-Goals (Không làm)
- ❌ Không build mobile app native (chỉ responsive web).
- ❌ Không xử lý dữ liệu GPS từ thiết bị phần cứng thật (dùng data giả lập).
- ❌ Không làm hệ thống quản lý tài chính/kế toán phức tạp (chỉ payment flow cơ bản).
- ❌ Không target multi-tenant (chỉ single-tenant demo).
- ❌ Không xử lý tiền thật ở môi trường demo (sử dụng `Stripe Test Mode`).

---

## 3. User Roles & Personas

| Role         | Mã hệ thống  | Mô tả ngắn                                              | Quyền chính                                                  |
| :----------- | :----------- | :------------------------------------------------------- | :----------------------------------------------------------- |
| **Nông dân** | `farmer`     | Chủ nông trại cà phê, khai báo thu hoạch                 | CRUD nông trại, tạo lô thu hoạch, xem lịch sử               |
| **Quản đốc** | `processor`  | Quản lý nhà máy chế biến, rang xay                       | Tiếp nhận hạt thô, tạo mẻ rang, cấp `Batch ID`, đóng gói    |
| **Tài xế**   | `driver`     | Lái xe vận chuyển, cập nhật GPS                          | Nhận cuốc, cập nhật trạng thái vận chuyển, gửi tọa độ GPS    |
| **QL Cửa hàng** | `store_mgr` | Quản lý cửa hàng bán lẻ                                | Xem tồn kho, tạo yêu cầu cung ứng, tiếp nhận hàng           |
| **Admin**    | `admin`      | Quản trị viên hệ thống                                   | Xem toàn bộ `Dashboard`, `Audit log`, quản lý user           |
| **Khách hàng** | `end_user` | Người tiêu dùng cuối (không cần đăng nhập)              | Quét QR, xem truy xuất nguồn gốc                             |

---

## 4. Business Flows (Luồng nghiệp vụ chính)

### 4.1 Luồng chính: Farm-to-Cup Pipeline

```
[Farmer]        [Processor]       [Warehouse]    [Payment]     [Logistics]    [Retail]      [End User]
   │                │                  │              │              │             │              │
   ├─ Thu hoạch ──►│                  │              │              │             │              │
   │ (Harvest       │                  │              │              │             │              │
   │  Created)      ├─ Rang xay ─────►│              │              │             │              │
   │                │ (BatchProcessed) │              │              │             │              │
   │                │                  ├─ Nhập kho ──►│             │             │              │
   │                │                  │(InventoryAdd)│              │             │              │
   │                │                  │              │              │             ├─ Đặt hàng    │
   │                │                  │              │              │             │(OrderCreated)│
   │                │                  │              │◄─────────────┤─────────────┤              │
   │                │                  │              │ Thanh toán   │             │              │
   │                │                  │              │(PaymentIntent│             │              │
   │                │                  │              │  Created)    │             │              │
   │                │                  │              │──► Stripe ──►│             │              │
   │                │                  │              │  (Webhook)   │             │              │
   │                │                  │◄─────────────┤──────────────┤             │              │
   │                │                  │ Reserve stock│              │             │              │
   │                │                  │(StockReserved│              │             │              │
   │                │                  │              │              │◄────────────┤              │
   │                │                  │              │              │ Điều xe     │              │
   │                │                  │              │              ├─ GPS ──────►│              │
   │                │                  │              │              ├─ Giao hàng ►│              │
   │                │                  │              │              │             │              │
   │                │                  │              │              │             │  Quét QR ────┤
   │                │                  │              │              │             │  Nguồn gốc   │
```

### 4.2 Saga Flow: Đặt hàng cung ứng (Supply Order) — Có thanh toán

Đây là luồng phức tạp nhất, triển khai `Saga Pattern` (Choreography) kết hợp `Payment Gateway`:

| Bước | Service          | Action                                | Event phát ra                | Failure → Compensation                           |
| :--- | :--------------- | :------------------------------------ | :--------------------------- | :----------------------------------------------- |
| 1    | `Retail`         | Cửa hàng tạo đơn cung ứng            | `SupplyOrderCreated`         | —                                                |
| 2    | `Payment`        | Tạo `PaymentIntent` trên Stripe       | `PaymentIntentCreated`       | —                                                |
| 3    | Frontend         | Hiển thị form thanh toán Stripe        | —                            | User huỷ → `PaymentCancelled` → Hủy đơn          |
| 4    | `Payment`        | Nhận Stripe Webhook, xác nhận thanh toán | `PaymentCompleted`         | `PaymentFailed` → Hủy đơn                        |
| 5    | `Warehouse`      | Kiểm tra & giữ chỗ hàng              | `StockReserved`              | `StockReserveFailed` → Refund Stripe + Hủy đơn   |
| 6    | `Logistics`      | Tìm xe & tạo chuyến                  | `ShipmentAssigned`           | `ShipmentFailed` → Release stock + Refund Stripe  |
| 7    | `Logistics`      | Tài xế nhận và vận chuyển             | `ShipmentInTransit`          | —                                                |
| 8    | `Logistics`      | Giao hàng thành công                  | `ShipmentDelivered`          | —                                                |
| 9    | `Warehouse`      | Xác nhận xuất kho                     | `StockDeducted`              | —                                                |
| 10   | `Retail`         | Tiếp nhận hàng, cập nhật kho          | `SupplyOrderCompleted`       | —                                                |

**Compensating Actions (Hoàn tác):**
- Nếu Bước 5 fail (hết hàng) → `Payment` nhận `StockReserveFailed`, gọi `Stripe Refund API`, phát `PaymentRefunded`. `Retail` nhận được → đơn chuyển `REJECTED`.
- Nếu Bước 6 fail (không có xe) → `Warehouse` nhận `ShipmentFailed`, release stock. `Payment` nhận event, gọi `Stripe Refund API`.
- **Nguyên tắc:** Tiền đã trừ ở Stripe thì phải có đường hoàn lại. Không bao giờ để tiền "treo".

### 4.3 Payment Flow Chi tiết

#### 4.3.1 Luồng thanh toán (Payment Sequence)

```
  [Frontend/POS]         [Payment Service]           [Stripe]             [Kafka]
       │                        │                       │                    │
       │── POST /orders ───────►│                       │                    │
       │   (SupplyOrderCreated  │                       │                    │
       │    from Kafka)         │                       │                    │
       │                        ├── Create PaymentIntent►│                   │
       │                        │◄── client_secret ──────┤                   │
       │◄── client_secret ──────┤                       │                    │
       │                        │                       │                    │
       │── Stripe.js submit ───►│───────────────────────►│                   │
       │   (card form)          │                       │                    │
       │                        │                       │                    │
       │                        │◄── Webhook POST ──────┤                   │
       │                        │   (payment_intent.     │                   │
       │                        │    succeeded)          │                   │
       │                        │                       │                    │
       │                        ├── HMAC verify ────────►│ (validate sig)    │
       │                        ├── Check Inbox ────────►│ (idempotency)     │
       │                        ├── Save + Outbox ──────►│                   │
       │                        │                       │  ──► PaymentCompleted
       │                        │                       │                    │
```

#### 4.3.2 Payment State Machine

```
                    ┌──────────────────────────────────────────┐
                    │                                          │
  ┌─────────┐   PaymentIntent   ┌──────────────┐              │
  │ CREATED ├──── created ─────►│   PENDING    │              │
  └─────────┘                   └──────┬───────┘              │
                                       │                      │
                          ┌────────────┼────────────┐         │
                     Webhook OK    Webhook FAIL   User cancel  │
                          │            │            │         │
                   ┌──────▼──────┐ ┌───▼──────┐ ┌──▼──────┐  │
                   │  SUCCEEDED  │ │  FAILED  │ │CANCELLED│  │
                   └──────┬──────┘ └──────────┘ └─────────┘  │
                          │                                   │
                   Saga Failure                                │
                   (Stock/Ship fail)                           │
                          │                                   │
                   ┌──────▼──────┐                            │
                   │  REFUNDED   ├────────────────────────────┘
                   └─────────────┘
```

**Trạng thái `Payment`:**

| Status       | Mô tả                                                      |
| :----------- | :---------------------------------------------------------- |
| `CREATED`    | Đơn hàng vừa tạo, chờ tạo `PaymentIntent`                  |
| `PENDING`    | `PaymentIntent` đã tạo trên Stripe, chờ user thanh toán    |
| `SUCCEEDED`  | Stripe xác nhận thanh toán thành công (via Webhook)         |
| `FAILED`     | Stripe báo thanh toán thất bại (thẻ bị từ chối, hết tiền)  |
| `CANCELLED`  | User huỷ bỏ thanh toán trước khi hoàn tất                   |
| `REFUNDED`   | Hoàn tiền qua Stripe Refund API (do Saga compensation)      |

#### 4.3.3 Payment Service Database Schema

```sql
-- Bảng chính: Ghi nhận thanh toán
CREATE TABLE payments (
    id                       UUID PRIMARY KEY,
    order_id                 UUID NOT NULL UNIQUE,
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    amount                   DECIMAL(12,2) NOT NULL,
    currency                 VARCHAR(3) DEFAULT 'VND',
    status                   VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    stripe_refund_id         VARCHAR(255),
    created_at               TIMESTAMPTZ DEFAULT NOW(),
    updated_at               TIMESTAMPTZ DEFAULT NOW()
);

-- Inbox: Chống xử lý lặp Webhook (Idempotency)
CREATE TABLE inbox_stripe_events (
    event_id    VARCHAR(255) PRIMARY KEY,  -- Stripe event ID (evt_xxx)
    event_type  VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);

-- Outbox: Đảm bảo event lên Kafka
CREATE TABLE outbox_events (
    id            UUID PRIMARY KEY,
    aggregate_id  UUID NOT NULL,
    event_type    VARCHAR(100) NOT NULL,
    payload       JSONB NOT NULL,
    published     BOOLEAN DEFAULT FALSE,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);
```

#### 4.3.4 Webhook Security (HMAC Validation)

Mọi Webhook từ Stripe đều phải qua middleware `HMAC` trước khi xử lý:

1. Stripe gửi `POST /api/v1/webhooks/stripe` kèm header `Stripe-Signature`.
2. `Payment Service` dùng `Webhook Secret` để băm payload nhận được bằng `crypto/hmac` (Go).
3. So sánh hash tính được với `Stripe-Signature` → khớp thì xử lý, không khớp thì reject (`HTTP 403`).
4. **Nguyên tắc:** Không tin bất kỳ request nào đến `/webhooks/*` mà không qua `HMAC verification`.

> **Mở rộng tương lai:** Kiến trúc `Payment Service` được thiết kế theo `Payment Gateway Abstraction`. Stripe là implementation đầu tiên. Các provider khác (VNPay, MoMo, PayOS) có thể được thêm vào bằng cách implement cùng interface mà không ảnh hưởng Saga flow.

### 4.4 Truy xuất nguồn gốc (QR Traceability)

Khi `end_user` quét mã QR trên bao bì cà phê, hệ thống trả về:

| Thông tin             | Nguồn service    | Ví dụ                                    |
| :-------------------- | :-------------- | :---------------------------------------- |
| Tên nông trại         | `Farm`          | "Nông trại Sơn La - Anh Minh"             |
| Khu vực canh tác      | `Farm`          | "Lô A3, Sơn La, độ cao 1200m"             |
| Giống cà phê          | `Farm`          | "Arabica Catimor"                          |
| Ngày thu hoạch        | `Farm`          | "2026-01-15"                               |
| Phương pháp chế biến  | `Process`       | "Washed Process"                           |
| Ngày rang             | `Process`       | "2026-02-01"                               |
| Độ rang               | `Process`       | "Medium Roast"                             |
| Batch ID              | `Process`       | `RR-SL-2026-0215-A3`                      |
| Lộ trình vận chuyển   | `Logistics`     | "Sơn La → Hà Nội (320km, 6h)"             |
| Ngày nhập cửa hàng    | `Retail`        | "2026-02-03"                               |

**Batch ID format:** `RR-{MÃ_VÙNG}-{NĂM}-{MMDD}-{LÔ}`

---

## 5. System Architecture (Tổng quan)

### 5.1 Architectural Foundations

#### Codebase: Monorepo
Toàn bộ project nằm trong một repository duy nhất với cấu trúc chia sẻ:
- **`api/`** — `gRPC` `Proto` definitions dùng chung (contracts giữa các services).
- **`pkg/`** — Shared libraries: `Database Wrapper`, `Middleware` (`Auth`, `Tracing`, `Logging`), `Kafka Producer/Consumer` base.
- **`services/`** — Mỗi Microservice là một module độc lập.

#### Service Structure: Clean Architecture (go-clean-arch v4)
Mỗi service áp dụng `Clean Architecture` với 3 lớp tách biệt rạch ròi:

| Lớp              | Nội dung                                                                      |
| :---------------- | :---------------------------------------------------------------------------- |
| **Domain**        | Entities, Value Objects, Repository Interfaces. Hoàn toàn không phụ thuộc framework. |
| **UseCase**       | Business logic, Saga step handlers, Command/Query handlers.                   |
| **Infrastructure**| `Gin` HTTP handlers, `gRPC` servers, `PostgreSQL` repos, `Kafka` producers.  |

#### Infrastructure: Polyglot Persistence
Lưu trữ đa phương thức, mỗi loại dữ liệu dùng đúng công cụ:

| Database          | Vai trò                                                       |
| :---------------- | :------------------------------------------------------------ |
| `PostgreSQL`      | Giao dịch lõi `ACID` — `Source of Truth` cho mỗi Microservice |
| `Apache Cassandra`   | Lưu trữ Event thô vĩnh cửu (`Audit` / `Event Sourcing`)       |
| `Elasticsearch`   | Tra cứu tốc độ cao, `CQRS Read Model` cho truy xuất nguồn gốc |
| `Redis`          | `Cache`, `Distributed Lock`, tọa độ `GPS` thời gian thực      |

#### Deployment: Docker on Proxmox
Triển khai 100% qua `Docker Compose`. Môi trường host là máy chủ `Proxmox` tự quản. Toàn bộ infra (DBs, Kafka, Services) chạy dưới dạng container, dễ dàng migrate lên Cloud VPS.

### 5.2 Service Inventory

| #  | Service                | Protocol                    | Database              | Kafka Role        |
| :- | :--------------------- | :-------------------------- | :-------------------- | :---------------- |
| 1  | `API Gateway`          | HTTP `:8000`                | —                     | —                 |
| 2  | `Identity Service`     | HTTP `:4433`                | Ory Kratos            | —                 |
| 3  | `Webhook Service`      | HTTP `:8088`                | PostgreSQL (Inbox)    | Producer only     |
| 4  | `Farm Service`         | HTTP `:8081` / gRPC `:9081` | PostgreSQL            | Producer          |
| 5  | `Processing Service`   | HTTP `:8082` / gRPC `:9082` | PostgreSQL            | Producer/Consumer |
| 6  | `Warehouse Service`    | HTTP `:8084` / gRPC `:9084` | PostgreSQL            | Producer/Consumer |
| 7  | `Logistics Service`    | HTTP `:8083` / gRPC `:9083` | PostgreSQL + Redis   | Producer/Consumer |
| 8  | `Retail Service`       | HTTP `:8085` / gRPC `:9085` | PostgreSQL            | Producer/Consumer |
| 9  | `Payment Service`      | HTTP `:8087` / gRPC `:9087` | PostgreSQL            | Producer/Consumer |
| 10 | `Traceability Service` | HTTP `:8086`                | Elasticsearch         | Consumer          |
| 11 | `Audit Service`        | Background worker           | Apache Cassandra      | Consumer          |
| 12 | `Monitor Service`      | HTTP `:8090` (SSE/WS)       | —                     | Consumer          |

#### `Webhook Service` (Ingress Gateway) — Chi tiết

`Webhook Service` là cửa ngõ chuyên biệt, **tách biệt hoàn toàn** khỏi `API Gateway`, chịu trách nhiệm tiếp nhận mọi luồng dữ liệu từ bên ngoài:

| Nguồn            | Endpoint                         | Xử lý                                                    |
| :--------------- | :------------------------------- | :------------------------------------------------------- |
| `Stripe`         | `POST /webhooks/stripe`          | `HMAC-SHA256` verify → `Inbox` dedup → publish to Kafka  |
| `VNPay`          | `POST /webhooks/vnpay`           | `HMAC-SHA512` verify → `Inbox` dedup → publish to Kafka  |
| `IoT Device`     | `POST /webhooks/iot/gps`         | Token auth → chuẩn hóa → publish `logistics.gps.updated` |

**Nguyên tắc:** `Webhook Service` **không chứa business logic**. Nhiệm vụ duy nhất: xác thực chữ ký → dedup bằng `Inbox` → chuẩn hóa payload → publish lên `Kafka`. Business logic xử lý ở service tương ứng.

### 5.3 Kafka Topic Design (Draft)

| Topic                              | Producer(s)      | Consumer(s)                         |
| :--------------------------------- | :--------------- | :---------------------------------- |
| `farm.harvest.created`             | Farm             | Processing, Trace, Audit            |
| `process.batch.completed`          | Processing       | Warehouse, Trace, Audit             |
| `retail.order.created`             | Retail           | Payment, Trace, Audit               |
| `payment.intent.created`           | Payment          | Retail, Trace, Audit                |
| `payment.completed`                | Payment          | Warehouse, Retail, Trace, Audit     |
| `payment.failed`                   | Payment          | Retail, Trace, Audit                |
| `payment.refunded`                 | Payment          | Retail, Trace, Audit                |
| `warehouse.stock.reserved`         | Warehouse        | Logistics, Payment, Trace, Audit    |
| `warehouse.stock.reserve-failed`   | Warehouse        | Payment, Retail, Trace, Audit       |
| `warehouse.stock.released`         | Warehouse        | Retail, Trace, Audit                |
| `logistics.shipment.assigned`      | Logistics        | Warehouse, Retail, Trace, Audit     |
| `logistics.shipment.failed`        | Logistics        | Warehouse, Payment, Retail, Trace, Audit |
| `logistics.shipment.delivered`     | Logistics        | Warehouse, Retail, Trace, Audit     |
| `logistics.gps.updated`            | Logistics        | Monitor, Trace                      |
| `retail.order.completed`           | Retail           | Trace, Audit                        |

### 5.4 HA-Ready Design Principles

Mặc dù triển khai demo trên single-node `Docker Compose`, kiến trúc được thiết kế sẵn để scale:

| Component           | HA Strategy                                                      |
| :------------------ | :--------------------------------------------------------------- |
| `Go Services`       | Stateless — scale horizontal bằng cách thêm container replicas   |
| `PostgreSQL`        | Mỗi service có DB riêng (database-per-service) → scale độc lập   |
| `Kafka`             | Multi-partition topics, consumer groups cho parallel processing   |
| `Redis`            | Hỗ trợ Cluster mode (Sentinel/Cluster) cho GPS data              |
| `Elasticsearch`     | Shard/Replica strategy cho read-model                            |
| `API Gateway`       | Stateless, có thể đặt sau Load Balancer                          |

### 5.5 Technical Mechanisms (Cơ chế Kỹ thuật)

#### A. Phân Quyền Phân Tán (Decentralized Authorization)

`API Gateway` chỉ đảm nhiệm **xác thực (Authentication)**: validate `JWT` và trích xuất `Roles` từ token, sau đó truyền xuống các service qua `HTTP Header` (`X-User-ID`, `X-User-Roles`).

Mỗi Microservice tích hợp `Casbin` vào lớp `Middleware` nội bộ. Dựa trên file `policy.csv` riêng của mình, service **tự phán quyết quyền truy cập (Authorization)** vào từng API endpoint — không phụ thuộc Gateway, tăng tính tự chủ và giảm tải tập trung.

```
[Client] ──► [API Gateway] ──► JWT validate + extract Roles ──► Header: X-Roles=store_mgr
                                                                       │
                                                               [Retail Service]
                                                               Casbin Middleware
                                                               policy.csv: store_mgr CAN POST /orders
                                                                       │
                                                               ✅ Allow / ❌ Deny
```

#### B. Tính Lũy Đẳng Kép (Dual Idempotency)

Hệ thống bảo vệ chống trùng lặp ở **hai tầng độc lập**:

| Tầng | Loại request | Cơ chế | Storage | TTL |
| :--- | :----------- | :----- | :------ | :-- |
| **Tầng 1** (Synchronous) | HTTP REST API | Header `Idempotency-Key` lưu vào `Redis` | `Redis` | 24h |
| **Tầng 2** (Asynchronous) | Kafka Consumer + Webhook | `Inbox Pattern` ghi `event_id` vào `PostgreSQL` | `PostgreSQL` | Vĩnh viễn |

- **Tầng 1:** Client gửi `Idempotency-Key: <uuid>` trong header. `API Gateway` middleware kiểm tra `Redis`. Nếu key đã tồn tại → trả response cached, không xử lý lại.
- **Tầng 2:** Consumer (Kafka/Webhook) trích xuất `event_id`, mở `Transaction`: kiểm tra `inbox_events` table → nếu đã có → rollback và bỏ qua → nếu chưa có → lưu và xử lý.

#### C. Quản Lý Cấu Hình (Configuration Management)

Sử dụng kết hợp `.env` files và thư viện `viper` (Go):
- Mỗi service có file cấu hình riêng (`config.yaml` + `.env` override).
- `viper` tự động đọc biến môi trường, hỗ trợ hot-reload.
- `Docker Compose` truyền `Secrets` và `DB Host` qua `environment` block.
- Không hardcode bất kỳ credential nào trong source code.

#### D. Dấu Vết Phân Tán (Distributed Tracing)

`OpenTelemetry` + `Jaeger`. Luồng hoạt động:

1. `API Gateway` cấp phát `Trace-ID` duy nhất cho mỗi request đến.
2. `Trace-ID` được tiêm vào:
   - `gRPC Context` (metadata) khi gọi internal service.
   - `HTTP Headers` (`traceparent`) khi gọi external service.
   - `Kafka Message Headers` khi publish event.
3. Mỗi service tạo `Span` con, liên kết với `Trace-ID` gốc.
4. `Jaeger UI` render `Gantt Chart` hiển thị toàn bộ timeline của giao dịch qua nhiều services.

> **Showcase value:** Một `Trace-ID` từ lúc cửa hàng đặt đơn đến khi Saga hoàn tất (qua Payment → Warehouse → Logistics) hiển thị trực quan trên Jaeger — rất ấn tượng với TA reviewer.

---

## 6. Phased Delivery Plan

### Phase 1: Foundation (Nền tảng)
> **Mục tiêu:** Dựng khung xương Monorepo, CI pipeline, và 3 core services đầu tiên.

| Deliverable                                 | Pattern showcase                     |
| :------------------------------------------ | :----------------------------------- |
| Project scaffold (Monorepo + shared libs)   | `DDD` project structure              |
| `API Gateway` (Gin + JWT parse + routing)   | `API Gateway Pattern`, `Rate Limiting` |
| `Identity Service` (Ory Kratos integration) | `OAuth2`/`OIDC`                      |
| `Farm Service` (CRUD nông trại + thu hoạch) | `Transactional Outbox`               |
| Kafka + PostgreSQL infra (Docker Compose)   | `Event-Driven` base                  |
| Casbin middleware (shared lib)              | `Decentralized Authorization`        |
| gRPC proto definitions (shared)            | `Protocol Buffers`                   |

### Phase 2: Supply Chain + Payment (Chuỗi cung ứng & Thanh toán)
> **Mục tiêu:** Hoàn thành pipeline Farm → Process → Warehouse, tích hợp Stripe, và Saga flow đầy đủ.

| Deliverable                                    | Pattern showcase                         |
| :--------------------------------------------- | :--------------------------------------- |
| `Processing Service` (rang xay, cấp Batch ID)  | `Event-Driven Consumer/Producer`         |
| `Warehouse Service` (tồn kho, reserve/release) | `Saga Participant`                       |
| `Retail Service` (đặt đơn cung ứng)            | `Saga Orchestrator`                      |
| `Payment Service` (Stripe PaymentIntent + Webhook) | `Webhook HMAC`, `Inbox Pattern`      |
| Payment Gateway Abstraction (interface-based)   | `Strategy Pattern` (Stripe, VNPay, etc.) |
| Stripe Webhook + HMAC validation middleware     | `Zero Trust` (external data)             |
| Saga flow đầy đủ + Auto-Refund Compensation     | `Saga Pattern` (Choreography)           |
| `Inbox Pattern` (Idempotency trên Consumer + Webhook) | `Transactional Inbox`             |

### Phase 3: Logistics & Real-time
> **Mục tiêu:** GPS tracking, vận chuyển, và data real-time.

| Deliverable                                   | Pattern showcase                   |
| :-------------------------------------------- | :--------------------------------- |
| `Logistics Service` (điều xe, tracking GPS)    | `Geo-spatial` (Redis GEO commands)|
| GPS simulator (fake driver coordinates)        | Real-time data pipeline            |
| Redis integration (cache + distributed lock)  | `Distributed Lock`, `Cache-aside`  |
| mTLS cho gRPC giữa các services               | `Zero Trust Architecture`          |

### Phase 4: Observability & Traceability
> **Mục tiêu:** CQRS read-model, audit trail, và distributed tracing.

| Deliverable                                    | Pattern showcase                  |
| :--------------------------------------------- | :-------------------------------- |
| `Traceability Service` (Kafka → Elasticsearch) | `CQRS` (Read-model projection)   |
| QR code generation + truy xuất nguồn gốc       | `CQRS Query` endpoint            |
| `Audit Service` (Kafka → Apache Cassandra, hash chain)  | `Event Sourcing` (Lite)           |
| `Hash Chaining` cho data integrity              | `Data Integrity Pattern`          |
| OpenTelemetry + Jaeger integration              | `Distributed Tracing`             |
| Prometheus + Grafana dashboards                 | `Observability Stack`             |

### Phase 5: Control Plane Visualization Dashboard (Frontend)
> **Mục tiêu:** ReactJS Dashboard kết hợp Operational UI + System Visualization.

| Deliverable                                     | Tech showcase                    |
| :---------------------------------------------- | :------------------------------- |
| `Monitor Service` (Kafka → SSE/WebSocket)        | Real-time event broadcasting     |
| ReactJS app (Vite + React Flow + Framer Motion)  | Modern frontend stack            |
| Split-screen: App View + System View             | `Service Mesh Visualization`     |
| Isometric 3D service map với animated edges       | `React Flow` animated edges      |
| Visual metaphors (Green Bean, Roasted Bean, etc.) | Domain-specific UI language      |
| Chaos Control panel (ngắt Kafka, set tồn kho = 0) | `Resiliency` demonstration      |
| CQRS Time Machine (time slider truy xuất)         | CQRS visual demo                 |

---

## 7. Non-Functional Requirements

| Category         | Requirement                                                         |
| :--------------- | :------------------------------------------------------------------ |
| **Performance**  | CQRS read query (truy xuất nguồn gốc) < 100ms                      |
| **Performance**  | GPS update latency < 500ms (từ Logistics → Monitor → Dashboard)    |
| **Availability** | HA-ready design: Stateless services, database-per-service           |
| **Scalability**  | Kafka multi-partition, consumer group ready                          |
| **Security**     | JWT Authentication (Ory Kratos), Casbin Authorization per-service   |
| **Security**     | mTLS cho tất cả gRPC internal communication                         |
| **Security**     | HMAC signature validation cho Stripe Webhook (chống giả mạo)        |
| **Security**     | Stripe Test Mode only — không xử lý tiền thật trong demo            |
| **Idempotency**  | `Inbox Pattern` cho Stripe Webhook (at-least-once → exactly-once)   |
| **Integrity**    | Hash chaining trên Audit log (Apache Cassandra) chống tamper                |
| **Observability**| Distributed tracing (OpenTelemetry + Jaeger) trên mọi request      |
| **Observability**| Prometheus metrics + Grafana dashboard cho mỗi service              |
| **Deployment**   | Full Docker Compose cho local dev                                    |
| **Deployment**   | Docker images sẵn sàng deploy lên Cloud VPS                         |

---

## 8. Technical Decisions (Đã chốt)

| Decision                     | Lựa chọn                    | Lý do                                                              |
| :--------------------------- | :-------------------------- | :----------------------------------------------------------------- |
| Backend language             | `Go` 1.22+                  | Performance, concurrency, ecosystem cho microservices              |
| Code organization            | `Monorepo`                  | Chia sẻ `proto`, `pkg` libs; dễ quản lý cross-service changes     |
| Service architecture         | `Clean Architecture` v4     | Tách Domain/UseCase/Infra — testable, replaceable adapters         |
| HTTP framework               | `Gin`                       | Lightweight, high-performance REST                                 |
| Internal RPC                 | `gRPC` + `Protobuf`         | Type-safe, high-speed, contracts định nghĩa tại `api/`            |
| Message broker               | `Apache Kafka`              | Industry standard cho event-driven architecture                    |
| Relational DB                | `PostgreSQL` v15            | ACID, mature, database-per-service                                 |
| Document DB                  | `Apache Cassandra`          | Wide-column store, flexible schema cho audit log, raw event storage |
| Search engine                | `Elasticsearch`             | Full-text search + CQRS read-model                                 |
| Cache / Real-time            | `Redis`                    | Redis alternative, GEO commands cho GPS, Idempotency-Key store     |
| Payment gateway              | `Stripe` (Test Mode)        | Industry standard, excellent API docs, Webhook support             |
| Payment abstraction          | `Strategy + Factory Pattern` | `PaymentProvider` interface + `ProviderFactory` → dễ thêm VNPay  |
| Identity                     | `Ory Kratos`                | Open-source, self-hosted identity management                       |
| Authorization                | `Casbin` + `policy.csv`     | Embeddable RBAC, decentralized per-service, Gateway chỉ authn     |
| Config management            | `viper` + `.env`            | Flexible, hot-reload, không hardcode secrets, Docker-friendly      |
| Frontend                     | `ReactJS` (Vite)            | Kết hợp Operational UI + Control Plane Visualization Dashboard                        |
| Visualization                | `React Flow`                | Node-based UI cho service mesh visualization                       |
| Animation                    | `Framer Motion`             | Micro-animations, glow effects                                     |
| State management             | `Zustand`                   | Lightweight state cho real-time WebSocket data                     |
| Tracing                      | `OpenTelemetry` + `Jaeger`  | Trace-ID từ Gateway, propagate qua gRPC/Kafka/HTTP headers         |
| Metrics                      | `Prometheus` + `Grafana`    | Industry standard monitoring stack                                 |
| Host infrastructure          | `Proxmox`                   | Self-hosted hypervisor, VM-based Docker environment                |
| Deployment                   | `Docker` + `Docker Compose` | Containerized, Cloud VPS ready                                     |

---

## 9. Risks & Mitigations

| Risk                                          | Impact | Mitigation                                              |
| :--------------------------------------------- | :----- | :------------------------------------------------------ |
| Solo developer → scope creep                   | High   | Strict phased delivery, MVP-first mindset                |
| Kafka learning curve                           | Medium | Start with single-partition, scale later                 |
| Too many databases to manage                   | Medium | Docker Compose quản lý toàn bộ infra                     |
| Control Plane Visualization Dashboard quá phức tạp                | High   | Phase 5 — chỉ làm sau khi backend ổn định               |
| HA design không thể verify trên single-node    | Low    | Ghi rõ HA strategy trong docs, verify bằng architecture review |
| Stripe Webhook bị replay/giả mạo              | High   | HMAC validation + Inbox idempotency pattern              |
| Tiền "treo" khi Saga fail sau thanh toán        | High   | Auto-refund compensation, Payment state machine nghiêm ngặt |
| Stripe API thay đổi version                    | Low    | Abstract qua interface, dễ swap provider                  |

---

## 10. Success Criteria

Dự án được coi là **thành công** khi:

- [ ] Toàn bộ 9+ services chạy đồng bộ trên `Docker Compose` mà không có lỗi.
- [ ] Demo được `Saga Rollback` scenario (hết hàng → auto-refund Stripe → đơn hàng tự hoàn tác).
- [ ] Demo được `Outbox Pattern` (ngắt Kafka → khôi phục → events flush thành công).
- [ ] Demo được `Inbox Pattern` (gửi Webhook trùng → hệ thống chỉ xử lý 1 lần).
- [ ] Stripe thanh toán thành công (Test Mode) và tiền được refund khi Saga fail.
- [ ] Quét QR trả về đầy đủ lịch sử Farm-to-Cup của một `Batch ID`.
- [ ] Control Plane Visualization Dashboard hiển thị data flow real-time giữa các services (bao gồm Payment flow).
- [ ] Reviewer/TA có thể đọc code và hiểu rõ từng `Pattern` được áp dụng.
- [ ] Hệ thống có thể deploy lên Cloud VPS bằng `docker-compose up -d`.

---

## Appendix A: Glossary

| Thuật ngữ              | Định nghĩa                                                                          |
| :--------------------- | :---------------------------------------------------------------------------------- |
| `Batch ID`             | Mã định danh duy nhất cho một mẻ cà phê thành phẩm: `RR-{VÙNG}-{NĂM}-{MMDD}-{LÔ}` |
| `Outbox Pattern`       | Ghi event vào DB cùng transaction với data nghiệp vụ, worker đẩy lên Kafka sau      |
| `Inbox Pattern`        | Lưu event đã xử lý để tránh xử lý lặp (idempotency)                                |
| `Saga`                 | Chuỗi giao dịch phân tán với cơ chế hoàn tác (compensation) khi có lỗi              |
| `CQRS`                 | Tách luồng ghi (Command) và đọc (Query) thành hai hệ thống riêng biệt               |
| `Clean Architecture`   | Kiến trúc 3 lớp (Domain/UseCase/Infra) — framework-agnostic, testable               |
| `Polyglot Persistence` | Dùng nhiều loại database khác nhau, mỗi loại phù hợp với đặc thù dữ liệu           |
| `Hash Chaining`        | Mỗi record audit chứa hash của record trước đó, tạo chuỗi bất biến (trên Cassandra) |
| `Control Plane Visualization`             | Dashboard giám sát toàn bộ kiến trúc hệ thống theo thời gian thực                   |
| `PaymentIntent`        | Object Stripe đại diện cho một giao dịch thanh toán đang chờ xử lý                  |
| `Webhook Service`      | Ingress Gateway chuyên biệt tiếp nhận và xác thực dữ liệu từ bên ngoài (Stripe, IoT) |
| `HMAC`                 | Hash-based Message Authentication Code — chữ ký số xác thực tính toàn vẹn           |
| `Dual Idempotency`     | Bảo vệ 2 tầng: Redis cho HTTP (sync) + Inbox Pattern cho Kafka/Webhook (async)     |
| `Idempotency-Key`      | UUID do client tạo, gửi trong HTTP header để chống thực thi lặp khi retry           |
| `ProviderFactory`      | Factory tạo ra Payment adapter (Stripe/VNPay) dựa trên config, dễ mở rộng           |
| `Compensating Action`  | Hành động hoàn tác (VD: Refund) khi một bước trong Saga bị lỗi                      |
| `Proxmox`              | Nền tảng hypervisor mã nguồn mở tự host, chạy Docker VM cho toàn bộ hệ thống        |
| `Trace-ID`             | ID duy nhất phát sinh tại Gateway, lan truyền xuyên suốt hệ thống để debug Saga     |

---

## Appendix B: Reference Documents

| Tài liệu                                                                 | Vị trí                                         |
| :------------------------------------------------------------------------ | :---------------------------------------------- |
| System Architecture                                                       | `docs/architecture/system-architecture.md`       |
| Database Schema Design                                                    | `docs/architecture/database-schema.md`           |
| UI/UX Visual Ideas (Control Plane Visualization Dashboard)                                   | `docs/ui-ux/visual-ideas.md`                     |
| REST API Specifications                                                   | `docs/api/rest-api.md`                           |
| gRPC Contract Definitions                                                 | `docs/api/grpc-contracts.md`                     |
| Deployment & Setup Guide                                                  | `docs/deployment/setup-guide.md`                 |
