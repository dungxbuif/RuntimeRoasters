# Technical Design: [RR-20] Warehouse Service Infrastructure & State Machine

**Role:** Tech Lead
**Context:** Triển khai nền tảng kỹ thuật cho Warehouse Service (Unified), quản lý luồng từ Receipt -> Stocked.

---

## 1. Technical Strategy

### 1.1. Core Framework
- **Language:** Go 1.22+.
- **Architecture:** Clean Architecture (Domain-driven).
- **Manual DI:** Sử dụng `cmd/main.go` để khởi tạo repository, usecase và handler. Không dùng wire nếu project chưa scale quá lớn (để dễ debug).

### 1.2. State Machine Implementation
- **Library:** `github.com/looplab/fsm` (hoặc custom logic nếu đơn giản).
- **Transitions:**
    - `RECEIVED` -> `HULLING` (Event: `StartHulling`)
    - `HULLING` -> `ROASTING` (Event: `StartRoasting`)
    - `ROASTING` -> `STOCKED` (Event: `CompleteRoasting`)
- **Validation:** Mọi thay đổi trạng thái phải đi qua phương thức `domain.Batch.TransitionTo(nextStatus)`.

---

## 2. Database Schema (PostgreSQL)

```sql
-- Quản lý vòng đời mẻ hàng
CREATE TABLE processing_batches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    harvest_id UUID NOT NULL, -- Tham chiếu từ Farm Service
    batch_id VARCHAR(50) UNIQUE NOT NULL, -- Mã sinh theo BATCH_LOGIC.md
    status VARCHAR(20) NOT NULL, -- RECEIVED, HULLING, ROASTING, STOCKED
    coffee_type VARCHAR(20) NOT NULL,
    origin_code VARCHAR(10) NOT NULL,
    intake_weight DECIMAL(10,2) NOT NULL,
    yield_weight DECIMAL(10,2) DEFAULT 0,
    weight_loss_percent DECIMAL(5,2) DEFAULT 0,
    quality_flag VARCHAR(20) DEFAULT 'NORMAL', -- NORMAL, QUALITY_WARNING
    operator_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Quản lý tồn kho SKU
CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) UNIQUE NOT NULL, -- e.g., ARABICA-CD-ROASTED
    coffee_type VARCHAR(20) NOT NULL,
    origin_code VARCHAR(10) NOT NULL,
    total_quantity DECIMAL(15,2) DEFAULT 0,
    available_quantity DECIMAL(15,2) DEFAULT 0,
    reserved_quantity DECIMAL(15,2) DEFAULT 0,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Inbox Pattern
CREATE TABLE inbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id VARCHAR(255) UNIQUE NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 3. Implementation Steps

1.  **Scaffolding:** Tạo thư mục `apps/warehouse-service` và copy `pkg/` boilerplate.
2.  **Domain Layer:** Định nghĩa `Batch` entity và interface cho `Repository`.
3.  **UseCase Layer:**
    - `ProcessHarvestEvent`: Nhận event từ Kafka, check Inbox, tạo Batch.
    - `UpdateBatchStatus`: Xử lý chuyển trạng thái và tính toán `yield_weight`.
4.  **Infrastructure Layer:**
    - Implement GORM repository.
    - Setup Kafka Consumer cho topic `farm.harvest.events`.

---

## 4. Verification Plan
- **Unit Test:** `TestBatchStateMachine` để verify logic `CanTransitionTo`.
- **SQL Test:** Kiểm tra transaction cập nhật đồng thời `processing_batches` và `inventory_items` khi chuyển sang `STOCKED`.
- **Manual Check:** Dùng `grpcurl` gọi API `UpdateBatchStatus` và kiểm tra log.
