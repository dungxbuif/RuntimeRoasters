# Sprint 3: Farm Service Business Logic (Vertical Slice)

**Goal:** Hoàn thiện dịch vụ Quản lý Nông trại (Farm Service) với đầy đủ tính năng CRUD và logic nghiệp vụ, áp dụng phong cách Vertical Slice (đi theo luồng tính năng) kết hợp với Clean Architecture.

---

## 📅 Roadmap & Flows

Sprint này tập trung vào 3 luồng dữ liệu chính:

### 1. Bootstrapping & Foundation
- Dựng khung sườn service mới dựa trên `demo-service`.
- Nối dây (wiring) cơ sở hạ tầng (DB, Redis, Gateway).
- **Patterns:** Composition Root, Dependency Injection (Wire).

### 2. Flow 1: Quản lý Danh sách & Tạo mới Nông trại
- Đi từ UI Form/Table -> KrakenD -> Farm Service -> Postgres.
- **Patterns:** Rich Domain Model, Repository Pattern, DTO, BFF.

### 3. Flow 2: Chi tiết, Cập nhật & Xóa
- Hoàn thiện các thao tác quản lý vòng đời nông trại.
- **Patterns:** Optimistic Locking, Unit of Work (WithTx).

### 4. Flow 3: Quản lý Lô hàng (Batch Management)
- Thiết kế quan hệ Aggregate Root giữa Farm và Batch.
- Chuẩn bị sẵn cấu trúc cho Outbox Pattern.

---

## 🎫 Tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-15](./RR-15/ticket.md) | Farm Service Bootstrapping | 🕒 To Do |
| [RR-16](./RR-16/ticket.md) | Flow 1: Create & List Farm | 🕒 To Do |
| [RR-17](./RR-17/ticket.md) | Flow 2: Update & Delete Farm | 🕒 To Do |
| [RR-18](./RR-18/ticket.md) | Flow 3: Batch Management | 🕒 To Do |
