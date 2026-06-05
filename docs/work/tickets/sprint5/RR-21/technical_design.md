# Technical Design: [RR-21] Saga Participant & Inventory Reservation

**Role:** Tech Lead
**Context:** Triển khai cơ chế tham gia Saga cho Warehouse Service, xử lý giữ hàng và rollback.

---

## 1. Technical Strategy

### 1.1. Distributed Locking
- **Tool:** Valkey (Redis) với thư viện `github.com/go-redsync/redsync`.
- **Key Format:** `lock:inventory:{sku}`.
- **TTL:** 5-10 giây (Đủ để hoàn tất DB transaction).
- **Rationale:** Đảm bảo khi có hàng ngàn request `ReserveStock`, việc kiểm tra `available_quantity` và trừ hàng được thực hiện tuần tự trên từng SKU, tránh overselling.

### 1.2. Idempotency (Inbox Pattern)
- **Key:** `order_id` + `action` (e.g., `RESERVE` or `CANCEL`).
- **Storage:** Bảng `inbox_events` trong Warehouse DB.
- **Flow:** Nhận Kafka Event -> Kiểm tra Inbox -> Thực thi logic -> Lưu Inbox + Update Inventory (Atomic Transaction).

---

## 2. Database Schema Additions

```sql
-- Quản lý giữ hàng (Saga state)
CREATE TABLE stock_reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    sku VARCHAR(50) NOT NULL,
    quantity DECIMAL(15,2) NOT NULL,
    status VARCHAR(20) NOT NULL, -- PENDING, COMMITTED, CANCELLED
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index để tìm nhanh theo order_id khi rollback
CREATE INDEX idx_reservations_order_id ON stock_reservations(order_id);
```

---

## 3. Implementation Steps

1.  **Distributed Lock Setup:** Tích hợp `redsync` vào Infrastructure layer.
2.  **ReserveStock UseCase:**
    - Acquire Lock.
    - Start DB Transaction.
    - Check `inventory_items.available_quantity`.
    - If enough: `UPDATE inventory_items SET available_quantity = available_quantity - X, reserved_quantity = reserved_quantity + X`.
    - `INSERT INTO stock_reservations`.
    - Commit Transaction & Release Lock.
    - Emit `warehouse.stock.reserved` event.
3.  **Compensating UseCase (CancelReservation):**
    - Ngược lại với Reserve: Tăng `available`, giảm `reserved`.
    - Update trạng thái reservation sang `CANCELLED`.

---

## 4. Verification Plan
- **Concurrency Test:** Chạy 50 goroutines cùng lúc gọi `ReserveStock` cho cùng 1 SKU có 100 đơn vị tồn kho. Verify kết quả cuối cùng phải chính xác 0 và không có transaction nào bị lỗi data race.
- **Saga Rollback Test:** Giả lập `OrderCreated` -> `StockReserved` -> `OrderCancelled`. Verify `available_quantity` quay về giá trị ban đầu.
