# [RR-16] Technical Design: Farm Repository & Database

**Status:** `DRAFT`
**Author:** Tech Lead

---

## 1. Context & Goal
Thiết kế Schema và triển khai `FarmRepository` sử dụng `GORM` để tương tác với PostgreSQL.

---

## 2. Database Schema

```sql
CREATE TABLE farms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    area DECIMAL(10, 2) NOT NULL,
    type VARCHAR(50) NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_farms_owner_id ON farms(owner_id);
CREATE INDEX idx_farms_location ON farms(location);
```

---

## 3. Implementation Details

### 3.1 Repository Interface
- **Path:** `internal/domain/repository.go`
- **Methods:**
    - `Create(ctx context.Context, farm *Farm) error`
    - `GetByID(ctx context.Context, id uuid.UUID) (*Farm, error)`
    - `List(ctx context.Context, filter FarmFilter) ([]*Farm, error)`
    - `Update(ctx context.Context, farm *Farm) error`
    - `Delete(ctx context.Context, id uuid.UUID) error`

### 3.2 GORM Implementation
- Sử dụng `db.Create()` cho Create.
- Sử dụng `db.Where().Find()` với dynamic filter cho List.
- Sử dụng `db.Save()` hoặc `db.Updates()` cho Update.

### 3.3 Data Ownership (ABAC)
- Mọi câu lệnh `SELECT`, `UPDATE`, `DELETE` phải luôn kèm theo `WHERE owner_id = $1` để đảm bảo an toàn dữ liệu.

---

## 4. Migration Plan
- Tạo file migration trong `deployments/init-db.sql` hoặc sử dụng một công cụ migration (nếu có).
