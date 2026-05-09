# Sprint 4: Distributed Saga & Messaging

**Goal:** Triển khai cơ chế giao dịch phân tán (Distributed Transactions) sử dụng mô hình Saga Choreography và Kafka để đảm bảo tính nhất quán dữ liệu giữa các microservices.

---

## 🎯 Sprint Goal

> **Hệ thống có thể xử lý luồng đặt hàng (Order) xuyên suốt từ Retail Service đến Warehouse Service, đảm bảo tính nguyên tử thông qua cơ chế bù trừ (Compensation) nếu có lỗi xảy ra.**

---

## 📋 Roadmap & Flows

Sprint này tập trung vào luồng xử lý đơn hàng phân tán:

### 1. Atomic Order Persistence (Outbox Pattern)
- Lưu đơn hàng và sự kiện vào DB trong cùng một transaction.
- Đảm bảo sự kiện luôn được gửi tới Kafka ít nhất một lần (At-least-once delivery).

### 2. Saga Participant: Stock Reservation
- Warehouse Service lắng nghe sự kiện `OrderCreated`.
- Thực hiện giữ hàng (reserve stock) và phản hồi kết quả.

### 3. Saga Rollback & Compensation
- Xử lý kịch bản thất bại (hết hàng).
- Hoàn tác trạng thái đơn hàng tại Retail Service.

### 4. Read-optimized View (CQRS)
- Tổng hợp dữ liệu từ nhiều sự kiện để phục vụ truy vấn nhanh.

---

## 🎫 Tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-19](./RR-19.md) | Retail Service & Outbox Pattern | 🕒 To Do |
| [RR-20](./RR-20.md) | Warehouse Service & Saga Participant | 🕒 To Do |
| [RR-21](./RR-21.md) | Saga Rollback & Compensations | 🕒 To Do |
| [RR-22](./RR-22.md) | Trace Service & CQRS | 🕒 To Do |
