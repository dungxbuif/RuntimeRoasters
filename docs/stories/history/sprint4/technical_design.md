# Technical Design: Sprint 4 — The Resilient Farm

**Tác giả:** Tech Lead
**Trạng thái:** Approved
**Cốt lõi:** Áp dụng **Transactional Outbox Pattern** để đảm bảo tính nguyên tử (Atomicity) giữa nghiệp vụ và sự kiện phát tán.

---

## 1. Kiến trúc Tổng thể (Architecture Overview)

Trong Microservices, việc cập nhật Database và gửi Message (Kafka) là hai hành động không thể bọc trong một Distributed Transaction (2PC). 
**Giải pháp:** Sử dụng bảng `outbox_events` làm hàng đợi tạm thời ngay trong chính Database của `farm-service`.

### Luồng dữ liệu (Data Flow):
1. **User** gọi `POST /v1/harvests`.
2. **UseCase** mở một Transaction:
    - Lưu bản ghi `harvests`.
    - Lưu bản ghi `outbox_events` (chứa thông tin mẻ thu hoạch).
3. **Transaction Commit** (Thành công 100% cả hai hoặc thất bại 100%).
4. **Relay Worker** (Go Routine) chạy ngầm:
    - Quét bảng `outbox_events` nơi `status = 'PENDING'`.
    - Publish message lên Kafka.
    - Cập nhật `processed_at = NOW()` và `status = 'COMPLETED'` sau khi Kafka xác nhận (Ack).

---

## 2. Thiết kế Cơ sở dữ liệu (Database Schema)

### 2.1. Bảng `harvests` (Mở rộng từ CRUD cơ bản)
```sql
CREATE TABLE harvests (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL,
    coffee_type VARCHAR(50) NOT NULL, -- Arabica, Robusta...
    quantity DECIMAL(12,2) NOT NULL,
    harvest_date TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) DEFAULT 'NEW', -- NEW, PROCESSING, COMPLETED
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 2.2. Bảng `outbox_events` (Mẫu chuẩn Blueprint)
```sql
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL, -- e.g., 'farm.harvest.created.v1'
    payload JSONB NOT NULL,            -- Data theo chuẩn CloudEvents
    metadata JSONB,                    -- TraceID, UserID, v.v.
    retry_count INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, PROCESSING, COMPLETED, FAILED
    processed_at TIMESTAMPTZ,          -- NULL nếu chưa gửi
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_outbox_unprocessed ON outbox_events(created_at) WHERE status = 'PENDING';
```

---

## 3. Cơ chế triển khai chi tiết (Deep-Dive Implementation)

### 3.1. Low-Level Relay Worker Logic (SELECT FOR UPDATE SKIP LOCKED)
Để hỗ trợ nhiều instance của Relay Worker chạy song song mà không gửi trùng tin, chúng ta sử dụng cơ chế khóa mức dòng của Postgres:

```sql
UPDATE outbox_events 
SET status = 'PROCESSING'
WHERE id IN (
    SELECT id FROM outbox_events 
    WHERE status = 'PENDING' 
    AND processed_at IS NULL
    ORDER BY created_at ASC 
    LIMIT 20 
    FOR UPDATE SKIP LOCKED
)
RETURNING *;
```
*Giải thích:* `SKIP LOCKED` giúp các instance khác không bị block, chúng sẽ bỏ qua các row đang bị instance này xử lý và lấy các row tiếp theo.

### 3.2. Distributed Tracing Propagation
Để OTel hiển thị đúng biểu đồ từ API -> DB -> Kafka -> Consumer:
1. **At API Layer:** Trích xuất `SpanContext` từ context.
2. **At Outbox Save:** Serialize `SpanContext` vào cột `metadata`.
3. **At Relay Worker:** Deserialize `metadata`, tạo một `Span` mới liên kết với `SpanContext` cũ (FollowsFrom), và inject `traceparent` vào Kafka Headers.

---

## 4. Failure Mode & Effects Analysis (FMEA)

| Tình huống | Kết quả | Cơ chế xử lý |
| :--- | :--- | :--- |
| **Kafka sập** | Message vẫn nằm trong DB | Relay Worker retry theo Exponential Backoff. |
| **Relay Worker sập** | Dữ liệu không bị mất | Khi Worker restart, nó tiếp tục quét các mẻ chưa xử lý. |
| **DB sập giữa chừng** | Không có dữ liệu nào được lưu | DB Transaction đảm bảo tính nguyên tử (Rollback toàn bộ). |
| **Duplicate Publish** | Kafka nhận 2 tin trùng | **Consumer side Idempotency** (Sử dụng `event.id` làm khóa trong bảng Inbox). |

---

## 5. Repository Pattern Implementation (The "Unit of Work")

Chúng ta sẽ implement một `TxManager` để bọc UseCase, đảm bảo Manual DI minh bạch:

```go
func (uc *HarvestUseCase) Create(ctx context.Context, input CreateInput) error {
    return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // 1. Business Logic & Save Harvest
        harvest := domain.NewHarvest(input)
        if err := uc.harvestRepo.Save(txCtx, harvest); err != nil {
            return err
        }

        // 2. Create Outbox Event in same transaction
        event := domain.NewOutboxEvent("harvest.created", harvest)
        return uc.outboxRepo.Save(txCtx, event)
    })
}
```

---

## 6. Gap Analysis & Steps

1. **Step 1:** Migration — Tạo bảng `harvests` và `outbox_events`.
2. **Step 2:** `pkg/database` — Viết `TxManager` hỗ trợ Manual DI.
3. **Step 3:** Repository — Viết logic `CreateHarvestWithOutbox(ctx, harvest, event)`.
4. **Step 4:** Relay Worker — Viết Background service chạy trong `main.go`.

---
**Tech Lead Signature**
