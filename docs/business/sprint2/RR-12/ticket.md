# [RR-12] Fine-grained Authorization (Casbin)

- **Summary:** Tích hợp Casbin với mô hình Hybrid RBAC+ABAC: Casbin kiểm tra Role (Gate 2), Repository filter quyền sở hữu dữ liệu bằng SQL.
- **Priority:** MEDIUM
- **Type:** Infrastructure

---

## 📖 User Story
> As a security officer, I want fine-grained access control that enforces both role-based permissions at the service boundary and data ownership at the repository level.

## 🔍 Acceptance Criteria
1. Request phải qua được Casbin RBAC check dựa trên Role.
2. Dữ liệu trả về phải được filter theo quyền sở hữu (Ownership).
3. Hỗ trợ phân cấp Role (Admin kế thừa Farmer).

## 🎫 Subtickets (Phân tích kỹ thuật & Triển khai)
1. [RR-12.1: Auth Service: Bootstrap & Persistence](./subtickets/RR-12.1/ticket.md) - Khởi tạo Auth Service, cấu hình Postgres Gorm Adapter và Casbin model.
2. [RR-12.2: Auth Service: gRPC Snapshot API](./subtickets/RR-12.2/ticket.md) - Triển khai gRPC API cung cấp bộ policy snapshot để bootstrap các dịch vụ khác.
3. [RR-12.3: Auth Service: Event-Driven Publishing](./subtickets/RR-12.3/ticket.md) - Tích hợp Kafka để publish sự kiện mỗi khi policy thay đổi.
4. [RR-12.4: Base Package: Resilient Reader Engine](./subtickets/RR-12.4/ticket.md) - Phát triển bộ engine đồng bộ 3 trụ cột (gRPC boot, Kafka live, Polling fallback) tại `pkg/base/casbin`.
5. [RR-12.5: Integration: Demo Service Wiring](./subtickets/RR-12.5/ticket.md) - Tích hợp bộ reader engine vào demo-service và kiểm thử đầu cuối.
