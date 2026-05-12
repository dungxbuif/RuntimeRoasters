# Ticket RR-6.2: [Tech] Saga Participant: Stock Reservation & Locking

**Mục tiêu:** Xây dựng dịch vụ Warehouse và cơ chế giữ chỗ hàng an toàn trong môi trường phân tán.

---

## 🛠️ Yêu cầu Kỹ thuật

### 1. Service Scaffolding
- Tạo folder `src/apps/warehouse-service`.
- Thiết lập kiến trúc Clean Architecture tương tự các service trước.
- Kết nối PostgreSQL và Valkey.

### 2. Stock Reservation API (Saga Protocol)
- Implement 3 phương thức gRPC:
    - `ReserveStock`: Giảm `available`, tăng `reserved`, tạo bản ghi `reservations`.
    - `ConfirmStock`: Giảm `reserved`, giảm `total`, chuyển status reservation sang `CONFIRMED`.
    - `ReleaseStock`: Tăng `available`, giảm `reserved`, chuyển status reservation sang `CANCELLED`.

### 3. Distributed Locking
- Tích hợp Valkey Lock để bảo vệ các thao tác Update tồn kho.
- Đảm bảo cơ chế Fail-safe: Nếu Valkey chết, phải có log cảnh báo và fallback/tạm dừng để bảo vệ dữ liệu.

### 4. Database Integrity
- Sử dụng SQL Constraint để ngăn chặn tồn kho âm (Check Constraint).
- Đảm bảo tính lũy đẳng (Idempotency) dựa trên `order_id`.

---

## 🧪 Kiểm thử (Verification)
- **Concurrency Test:** Chạy script gọi 100 request `ReserveStock` đồng thời cho cùng 1 mẻ hạt chỉ còn 10kg -> Chỉ 1 request thành công.
- **Saga Simulation:** Giả lập chuỗi gọi `Reserve` -> `Release` và kiểm tra tồn kho quay về trạng thái ban đầu chính xác 100%.
- **Idempotency Test:** Gọi `Reserve` 2 lần với cùng 1 `order_id` -> Hệ thống chỉ trừ kho 1 lần.
