# Technical Design: Sprint 6 — Inventory Consistency & Saga Lite (Reference Aligned)

Mục tiêu: Hoàn thiện luồng Saga Choreography (Retail <-> Warehouse) sử dụng các pattern từ dự án UrbanX: **Soft Reservations**, **Exactly-once Consumer**, và **Transactional Inbox**.

---

## 1. Cơ chế Tạm giữ hàng (Inventory Reservations)
*Tham chiếu: UrbanX Inventory Service logic*

Để chống bán lố (Over-selling) và đảm bảo tính nhất quán, Warehouse Service áp dụng mô hình 2 chỉ số:

### 1.1. Cấu trúc bảng `inventory`
- `quantity_available`: Số lượng thực tế trong kho.
- `quantity_reserved`: Số lượng đang bị tạm giữ cho các đơn hàng chưa thanh toán.

**Công thức khả dụng**: `AvailableForSale = quantity_available - quantity_reserved`.

### 1.2. Logic Reservation (Soft Commit)
Khi nhận `OrderCreated`:
1. Kiểm tra: `IF (quantity_available - quantity_reserved) >= req.quantity`.
2. Tạm giữ: `UPDATE inventory SET quantity_reserved = quantity_reserved + req.quantity WHERE id = ...`.
3. Ghi nhật ký: `INSERT INTO inventory_reservations (order_id, quantity, status='RESERVED')`.

### 1.3. Chốt chặn an toàn (Defensive Check)
Khi giải phóng hàng (Release/Rollback), áp dụng logic tương tự `Math.Max(0, ...)` để tránh số lượng bị âm:
```go
newReserved := inventory.QuantityReserved - reservation.Quantity
if newReserved < 0 {
    newReserved = 0 // "Phòng vệ dữ liệu" tránh lỗi làm lệch kho
}
```

---

## 2. Exactly-once Consumer (Inbox Pattern)
*Tham chiếu: UrbanX Transactional Inbox*

Đảm bảo mỗi sự kiện đặt hàng chỉ được trừ kho đúng một lần duy nhất, kể cả khi Kafka gửi lặp tin nhắn.

### 2.1. Transactional Inbox
Mọi thao tác thay đổi kho PHẢI nằm trong cùng 1 DB Transaction với việc ghi vào bảng `inbox_events`:
1. `BEGIN TRANSACTION`
2. `INSERT INTO inbox_events (id) VALUES (kafka_message_id)` -> Nếu trùng ID, DB sẽ báo lỗi Unique Constraint và rollback toàn bộ.
3. Thực hiện `Update Inventory`.
4. `COMMIT`

---

## 3. Saga Choreography & Compensation
*Tham chiếu: UrbanX Saga Flow*

- **Happy Path**: `OrderCreated` -> `StockReserved` -> `PaymentSucceeded` (Mock) -> `OrderConfirmed`.
- **Compensation Path**: `OrderCreated` -> `StockReserved` -> `PaymentFailed` (Mock) -> `OrderCancelledEvent` -> **Warehouse nhả kho (Release Reserved)**.

---

## 4. Anti-Corruption Layer (ACL) cho Mock Payment
*Tham chiếu: UrbanX Payment Gateways*

Ngay từ giai đoạn Mock, chúng ta sẽ thiết kế lớp ACL để bảo vệ core logic:
- Định nghĩa interface `IPaymentGateway` trong `payment-service`.
- Triển khai `MockGateway` thực thi interface này.
- Giúp việc chuyển sang **Stripe thật** ở Sprint 11 chỉ là việc thay thế lớp thực thi (Implementation), không sửa đổi UseCase.

---
*TechLead Signed-off: 2026-05-10*
