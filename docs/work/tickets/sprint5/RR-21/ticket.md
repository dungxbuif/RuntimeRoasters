# [RR-21] [Tech] Saga Participant & Inventory Reservation

**User Story:**
Dưới vai trò là **Architect**, tôi muốn Warehouse Service có khả năng tham gia vào các giao dịch phân tán (Saga) để thực hiện giữ chỗ hàng (Reserve) một cách an toàn và nhất quán.

**Technical Value:**
Đây là mảnh ghép quan trọng để hệ thống có thể xử lý các đơn hàng phức tạp mà không bị tình trạng "bán quá số lượng" (Overselling).

---

## 🛠️ Yêu cầu Kỹ thuật (Technical Requirements)

### 1. Saga Participant Logic
- Lắng nghe sự kiện `retail.order.created` từ Kafka.
- Thực hiện logic `ReserveStock`:
    - Kiểm tra `available_quantity` trong DB.
    - Nếu đủ: Giảm `available_quantity`, tăng `reserved_quantity`, lưu bản ghi `stock_reservations`. Phát event `warehouse.stock.reserved`.
    - Nếu không đủ: Phát event `warehouse.stock.failed`.

### 2. Compensating Action (Rollback)
- Lắng nghe sự kiện `retail.order.cancelled` hoặc `payment.failed`.
- Thực hiện logic `CancelReservation`:
    - Hoàn trả `reserved_quantity` về `available_quantity`.
    - Chuyển trạng thái reservation thành `CANCELLED`.

### 3. Distributed Locking (Valkey)
- Sử dụng Valkey để thực hiện Distributed Lock trên `SKU_ID` trong quá trình Reserve, tránh Race Condition khi có hàng ngàn đơn hàng cùng lúc.

---

## ✅ Acceptance Criteria (AC)
1. **At-least-once Processing:** Sử dụng Inbox Pattern để đảm bảo không xử lý lặp một event reservation.
2. **Data Consistency:** Sau một chuỗi Saga hoàn chỉnh (Thành công hoặc Thất bại), tổng số lượng `Available + Reserved` phải không đổi so với ban đầu.
3. **Performance:** Thời gian xử lý logic Reserve (bao gồm cả lock) phải < 100ms.
