# Lộ trình Phát triển RuntimeRoasters (Portfolio Edition)

**Mục tiêu cốt lõi:** Showcase Senior Backend — kiến trúc phân tán, OTel, Clean Architecture, DI, API Gateway.

**Chiến lược:**
- 1 Solo Developer (Full-stack/Backend focus)
- Sprint 1 = **Infrastructure only** — mọi boilerplate, base code, infra setup phải hoàn tất
- Sprint 2+ = **Business logic only** — chỉ viết domain/usecase/repo, không setup thêm gì
- Ticket style: BDD AC only. Technical design do Dev tự tạo khi implement.

---

## Sprint 1 — Complete Infrastructure

**Goal:** Sau Sprint 1, mọi thứ đều chạy được end-to-end. Sprint 2 chỉ viết business logic.

**Demo Service** (`apps/demo-service/`) = canonical template cho tất cả services sau. Farm Service bắt đầu sạch từ Sprint 2 bằng cách copy từ đây.

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-1](./sprint1/RR-1/ticket.md) | Infra Kick-off — Monorepo, Go Workspaces | ✅ Done |
| [RR-2](./sprint1/RR-2/ticket.md) | Config & Logger — Viper, Zap, BaseConfig | ✅ Done |
| [RR-3](./sprint1/RR-3/ticket.md) | Base & Errs — RFC 9457, app lifecycle | ✅ Done |
| [RR-4](./sprint1/RR-4/ticket.md) | Proto + Wire + KrakenD + Client Shell | 🕒 To Do |
| [RR-5](./sprint1/RR-5/ticket.md) | Docker Complete — SigNoz + Full Compose | 🕒 To Do |
| [RR-6](./sprint1/RR-6/ticket.md) | Demo Service — Clean Arch + Full Stack Slice | 🕒 To Do |
| [RR-7](./sprint1/RR-7/ticket.md) | Client App — Control Plane | 🕒 To Do |
| [RR-8](./sprint1/RR-8/ticket.md) | Client App — Business UI Shell | 🕒 To Do |

Chi tiết: [`docs/business/sprint1/main.md`](./sprint1/main.md)

---

## Sprint 2 — Security & Access Control

**Goal:** Thiết lập nền tảng AuthN/AuthZ toàn diện. Mọi API call phải được bảo mật trước khi viết business logic.

| Ticket | Summary |
| :--- | :--- |
| [RR-9](./sprint2/RR-9/ticket.md) | Identity Server Infrastructure (Ory Kratos) |
| [RR-10](./sprint2/RR-10/ticket.md) | Client-Side Auth & Login Flow |
| [RR-11](./sprint2/RR-11/ticket.md) | Backend Security Core - JWT Validation |
| [RR-12](./sprint2/RR-12/ticket.md) | Fine-grained Authorization (Casbin) |
| [RR-13](./sprint2/RR-13/ticket.md) | End-to-End Secure Integration |
| [RR-14](./sprint2/RR-14/ticket.md) | Token Revocation & Global Logout |

Chi tiết: [`docs/business/sprint2/main.md`](./sprint2/main.md)

---

## Sprint 3 — Farm Service Business Logic

**Goal:** Farm Service đầy đủ CRUD. Copy boilerplate từ `apps/demo-service/`, chỉ viết business logic.

| Ticket | Summary |
| :--- | :--- |
| [RR-15](./sprint3/RR-15/ticket.md) | Farm Service Bootstrap & Domain |
| [RR-16](./sprint3/RR-16/ticket.md) | Flow 1: Create & List Farm |
| [RR-17](./sprint3/RR-17/ticket.md) | Flow 2: View, Update, Delete Farm |
| [RR-18](./sprint3/RR-18/ticket.md) | Flow 3: Batch Management |


Chi tiết: [`docs/business/sprint3/main.md`](./sprint3/main.md)

---

## Sprint 4 — Security & Core Infra Hardening

**Epic Goal:** Nền tảng hạ tầng và bảo mật đạt chuẩn "Production-ready".

| Ticket | Summary |
| :--- | :--- |
| [RR-13](../sprint2/RR-13/ticket.md) | Security: End-to-End Auth Integration |
| [RR-14](../sprint2/RR-14/ticket.md) | Security: Token Revocation (Logout) |
| [RR-36](./sprint4/RR-36.md) | PgBouncer Integration & Connection Pooling |

Chi tiết: [`docs/business/sprint4/main.md`](./sprint4/main.md)

---

## Sprint 5 — Distributed Order Orchestration

**Epic Goal:** Ghi đơn hàng và phát tán sự kiện một cách nguyên tử.

| Ticket | Summary |
| :--- | :--- |
| [RR-19](./sprint5/RR-19.md) | Retail Service & Outbox Pattern |

---

## Sprint 6 — Inventory Consistency & Saga Lite

**Epic Goal:** Hoàn thành luồng Saga 2 bước (Retail-Warehouse) đảm bảo tính lũy đẳng.

| Ticket | Summary |
| :--- | :--- |
| [RR-20](./sprint6/RR-20.md) | Warehouse Service & Saga Participant |
| [RR-21](./sprint6/RR-21.md) | Saga Rollback & Compensations |

---

## Sprint 7 — Logistics & Real-time Delivery

**Epic Goal:** Tự động điều phối vận chuyển khi kho đã giữ hàng.

| Ticket | Summary |
| :--- | :--- |
| [RR-23](./sprint7/RR-23.md) | Logistics Service Core (Shipments & Drivers) |

---

## Sprint 8 — Geographic Intelligence

**Epic Goal:** Theo dõi vị trí tài xế trên bản đồ theo thời gian thực.

| Ticket | Summary |
| :--- | :--- |
| [RR-24](./sprint8/RR-24.md) | Redis Geo Tracking Integration |
| [RR-25](./sprint8/RR-25.md) | GPS Simulator & Event Publishing |

---

## Sprint 9 — System-wide Traceability & Search

**Epic Goal:** Tìm kiếm và truy xuất toàn bộ lịch sử hạt cà phê trong < 100ms.

| Ticket | Summary |
| :--- | :--- |
| [RR-27](./sprint9/RR-27.md) | Traceability Service & CQRS Read Model |
| [RR-35](./sprint9/RR-35.md) | Elasticsearch Integration & Index Design |

---

## Sprint 10 — Reliability & Observability

**Epic Goal:** Hệ thống minh bạch (Transparent) và có Audit log không thể sửa đổi.

| Ticket | Summary |
| :--- | :--- |
| [RR-28](./sprint10/RR-28.md) | Audit Service & Cassandra Storage |
| [RR-29](./sprint10/RR-29.md) | Full OpenTelemetry & Jaeger Setup |
| [RR-30](./sprint10/RR-30.md) | Prometheus & Grafana Monitoring |

---

## Sprint 11 — Real-world Commerce Integration

**Epic Goal:** Xử lý thanh toán thực tế và cơ chế hoàn tiền (Refund) tự động.

| Ticket | Summary |
| :--- | :--- |
| [RR-37](./sprint11/RR-37.md) | Payment Service & Stripe API Integration |

---

## Sprint 12 — Control Plane & The Grand Finale

**Epic Goal:** Trực quan hóa toàn bộ hệ thống và chứng minh khả năng tự phục hồi.

| Ticket | Summary |
| :--- | :--- |
| [RR-31](./sprint12/RR-31.md) | Monitor Service & SSE Broadcasting |
| [RR-32](./sprint12/RR-32.md) | Frontend: React Flow System Topology |
| [RR-33](./sprint12/RR-33.md) | Chaos Control Panel (Resiliency UI) |
| [RR-34](./sprint12/RR-34.md) | Consumer QR Traceability Flow (Final UI) |
| [RR-26](./sprint12/RR-26.md) | Security: Internal gRPC mTLS |
