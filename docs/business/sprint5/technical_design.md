# Technical Design: Sprint 5 — Distributed Order Orchestration (Finalized)

Mục tiêu: Đảm bảo việc tạo đơn hàng và phát hành sự kiện là một hành động nguyên tử (Atomic), chống trùng lặp tại cửa ngõ API.

---

## 1. Transactional Outbox Pattern (Relay Worker)

### 1.1. DB Schema (Retail Service)
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    total_amount DECIMAL(12,2),
    status VARCHAR(50) DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, PROCESSED, FAILED
    retry_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 1.2. UseCase: `CreateOrder`
Sử dụng GORM Transaction để đảm bảo tính Atomicity:
1. Mở Transaction.
2. Tạo bản ghi `orders`.
3. Tạo bản ghi `outbox_events` (Chứa CloudEvents metadata).
4. Commit.

### 1.3. Relay Worker Logic
- Tần suất quét: 500ms.
- Logic: Quét các event `PENDING`, publish lên Kafka topic `retail.order.events`.
- Sau khi Kafka xác nhận (Ack): Cập nhật trạng thái thành `PROCESSED`.
- Ưu điểm: Đảm bảo "At-least-once delivery" mà không cần hạ tầng phức tạp.

---

## 2. Idempotency Tầng 1 (API Level - Comprehensive)

### 2.1. Middleware Xử lý (Valkey/Redis)
- **Key Strategy**: `idemp:{user_id}:{idempotency_key}`.
- **Quy trình 3 giai đoạn**:
    1. **Locking**: `SET {key} "PROCESSING" NX EX 30`. Nếu thất bại -> Trả về `425 Too Early`.
    2. **Execution**: Chạy Business Logic.
    3. **Commit**: Cập nhật Redis với giá trị:
       ```json
       {
         "status": "COMPLETED",
         "http_code": 201,
         "body": "..." 
       }
       ```
- **TTL**: 24 giờ.

---

## 3. Event Schema (CloudEvents Standard)
Mọi sự kiện phát ra từ Retail Service phải tuân thủ:
- `id`: Unique Event ID.
- `source`: `/retail-service`.
- `type`: `order.created.v1`.
- `data`: Business payload.
- `idempotency_key`: Đính kèm để các service sau sử dụng nếu cần.

---
*TechLead Signed-off: 2026-05-10*
