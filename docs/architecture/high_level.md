# RuntimeRoasters — High-Level Overview

**RuntimeRoasters** là nền tảng quản lý chuỗi cung ứng và logistics mô phỏng vòng đời của hạt cà phê từ nông trại đến tách cà phê bán lẻ (Farm-to-Cup).

Dự án được xây dựng như bản showcase kỹ thuật, 100% Golang, áp dụng DDD và các mẫu thiết kế phân tán hạng nặng.

---

## Câu chuyện Doanh nghiệp

Trong văn hóa Việt, cà phê không chỉ là thức uống mà là "sợi dây" kết nối xã hội. **RuntimeRoasters** số hóa toàn bộ "mạch sống" này — mỗi hạt cà phê đều có danh tính số, từ nông trường qua nhà máy, logistics, đến tách cà phê trên tay khách hàng.

### Các điểm chạm chính

1. **Thượng nguồn (Farm):** Nông dân cập nhật diện tích và khai báo mẻ cà phê vừa thu hoạch.
2. **Trung nguồn (Processing & Warehouse):** Nhà máy tiếp nhận, bóc vỏ, phơi, rang, đóng gói thành Batch ID.
3. **Vận tải (Logistics):** Tài xế nhận điều phối, vận chuyển, cập nhật GPS thời gian thực.
4. **Hạ nguồn (Retail):** Quản lý cửa hàng theo dõi tồn kho, nhập hàng.
5. **Truy xuất & Quản trị:** Quét QR → xem hành trình đầy đủ. Admin giám sát toàn cảnh qua Control Plane.

---

## Thành phần Hệ thống

| Service | Chức năng | Patterns |
| :--- | :--- | :--- |
| **client-app** | Business UI (user) | React/Next.js |
| **control-app** | Control Plane (admin) | BFF, Transparent Proxy |
| **Gateway (KrakenD)** | Authentication offload, Rate Limiting, Routing | Gateway Pattern |
| **demo-service** | Sprint 1 canonical template | Clean Arch, OTel, DI |
| **Farm** | Quản lý nông hộ, vườn cây, thu hoạch | Outbox Pattern |
| **Process** | Chế biến mẻ rang, cấp Batch ID | Event-Driven |
| **Logistics** | Điều phối xe, tracking GPS | Real-time Geo (Redis) |
| **Warehouse** | Reserve, xuất/nhập kho | Saga (Participant) |
| **Retail** | Cửa hàng đặt hàng, tiêu thụ | Saga (Orchestrator) |
| **Trace** | Tổng hợp hành trình vào Elasticsearch | CQRS |
| **Audit** | Lưu Kafka history vào Cassandra | Event Sourcing (Lite) |

---

## Ngăn xếp Công nghệ

- **Backend:** Go 1.22+ với Go Workspaces; Google Wire (DI); Gin (HTTP); gRPC (internal RPC).
- **API Gateway:** KrakenD.
- **Message Broker:** Redpanda (Kafka-compatible).
- **Databases:** PostgreSQL (source of truth), Elasticsearch (CQRS read), Cassandra (audit), Redis (cache/geo).
- **Observability:** SigNoz (traces + metrics + logs, ClickHouse-backed). OTel SDK trong mọi service. Xem [telemetry.md](./telemetry.md).
- **Frontend:**
  - **client-app:** Next.js 15 App Router (Business UI).
  - **control-app:** Next.js 15 (Admin Dashboard + SigNoz Proxy). swagger-ui-react cho API Explorer.
- **Security:** Ory Kratos (Identity Provider), Ory Hydra (OAuth2/OIDC Provider), Casbin (RBAC/ABAC).
  - *Note:* Để xem chi tiết luồng đăng nhập SSO phân tán (Zero-Consent) và cấu hình Identity, tham khảo [Identity SSO Diagram](./identity-flow.md) và [Identity SSO Implementation Details](../engineering/identity-sso-implementation.md).
---

## Data Flow Patterns

### Control Plane (Admin — Dashboard)

```
Browser → control-app :3001 → internal services
```

### Observability (SigNoz UI)

```
Browser → SigNoz :3301 (Direct Access)
```

### Business UI (User — Direct)

```
Browser → client-app /app/* → KrakenD :8081 → Microservices (gRPC)
```

---

## Distributed Tracing

SigNoz + OTel SDK tạo waterfall trace xuyên suốt toàn hệ thống. Mỗi request có `trace_id` duy nhất lan truyền qua HTTP Headers, gRPC Metadata, và Kafka Headers (W3C Trace Context).

Admin xem trace trực tiếp tại SigNoz UI (:3301).
