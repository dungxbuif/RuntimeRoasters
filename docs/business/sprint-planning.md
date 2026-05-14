# Lộ trình Phát triển RuntimeRoasters (Master Plan)

**Mục tiêu cốt lõi:** Xây dựng hệ thống Microservices mẫu mực — kiến trúc phân tán, Event-Driven, Saga Pattern, Observability (Ưu tiên Manual DI & Selective Outbox).

---

## 🚀 Tổng quan Roadmap (Sprints 4-12)

| Sprint | Goal | Key Deliverables (Business + Tech) |
| :--- | :--- | :--- |
| **S4** | **The Resilient Farm** | Transactional Outbox, Reliable Harvesting API, CloudEvents Blueprint. |
| **S5** | **Warehouse & Processing** | Unified Warehouse Service, Coffee Processing State Machine, Inventory & Saga Participant. |
| **S6** | **Retail & Order Saga** | Retail Service, Saga Orchestrator (Choreography), Idempotent Ordering. |
| **S7** | **Real-time Logistics** | Logistics Service, Driver Trips, Valkey GEO Real-time Tracking. |
| **S8** | **Financial Integrity** | Payment Service, Stripe Integration, Webhook HMAC, Saga Rollbacks. |
| **S9** | **Traceability (CQRS)** | Trace Service, Elasticsearch Read-Model, 360-degree Bean Journey View. |
| **S10** | **Observability & Audit** | Audit Service (Cassandra), System Health Dashboards. |
| **S11** | **Hardening & Finale** | mTLS Zero Trust, Token Revocation, PgBouncer, Chaos Engineering. |
| **S12** | **Project Polish** | Final technical debt cleanup, Load testing, and Final Demo. |

---

## 📋 Chi tiết các Sprint (Tiến độ)

### [x] Sprint 1-3: Foundation & Farm Core (Completed)
Hoàn tất hạ tầng cơ bản và các API CRUD cho Farm Service.

### [x] Sprint 4: The Resilient Farm (Completed)
- **Goal:** Đảm bảo thu hoạch không bao giờ mất dữ liệu và dọn dẹp nợ kỹ thuật Sprint 3.
- **Tickets:** [RR-4.0], [RR-4.1], [RR-4.2].
- Chi tiết: [`docs/business/sprint4/main.md`](./sprint4/main.md)

### [ ] Sprint 5: Warehouse & Processing
- **Goal:** Hợp nhất quy trình nhà máy và kho bãi. Xử lý hạt từ lúc nhập kho sau thu hoạch, qua rang xay đến lưu kho thành phẩm và giữ hàng cho Saga.
- Chi tiết: [`docs/business/sprint5/main.md`](./sprint5/main.md)

### [ ] Sprint 6: Retail & Order Saga
- **Goal:** Đặt hàng và điều phối chuỗi cung ứng tự động (Choreography Saga).
- Chi tiết: [`docs/business/sprint6/main.md`](./sprint6/main.md)

### [ ] Sprint 7: Real-time Logistics
- **Goal:** Theo dõi xe vận chuyển trên bản đồ.
- Chi tiết: [`docs/business/sprint7/main.md`](./sprint7/main.md)

### [ ] Sprint 8: Financial Integrity
- **Goal:** Thanh toán thực và hoàn tiền tự động.
- Chi tiết: [`docs/business/sprint8/main.md`](./sprint8/main.md)

### [ ] Sprint 9: Traceability (CQRS)
- **Goal:** Truy xuất nguồn gốc tốc độ cao.
- Chi tiết: [`docs/business/sprint9/main.md`](./sprint9/main.md)

### [ ] Sprint 10: Observability & Audit
- **Goal:** Giám sát toàn diện và Audit log bất biến.
- Chi tiết: [`docs/business/sprint10/main.md`](./sprint10/main.md)

### [ ] Sprint 11: Hardening & Grand Finale
- **Goal:** Bảo mật mTLS, xử lý nợ kỹ thuật tồn đọng (PgBouncer, Token Revocation).
- Chi tiết: [`docs/business/sprint11/main.md`](./sprint11/main.md)

### [ ] Sprint 12: Project Polish
- **Goal:** Load testing, dọn dẹp code và Demo tổng thể.
- Chi tiết: [`docs/business/sprint12/main.md`](./sprint12/main.md)
