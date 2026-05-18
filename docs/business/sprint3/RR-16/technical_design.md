# Technical Design - [RR-16] Flow 1: Create & List Farm

Mục tiêu: Triển khai logic nghiệp vụ và lưu trữ cho thực thể Farm.

## 🏗️ 1. Database Schema (PostgreSQL)
Bảng `farms`:
```sql
CREATE TABLE farms (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location TEXT,
    area DECIMAL(10,2) NOT NULL,
    coffee_type VARCHAR(100),
    owner_id UUID NOT NULL, -- Farm Manager Subject ID from JWT
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_farms_owner ON farms(owner_id);
```

## 🛠️ 2. Repository Layer (Data Scoping)
- **File:** `internal/infrastructure/postgres/farm_repo.go`
- **Logic quan trọng:** Mọi hàm truy vấn danh sách PHẢI lọc theo `owner_id`.
- **Snippet:** `db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&farms)`

## 🧠 3. UseCase Layer (Validation & Ownership)
- **Logic `CreateFarm`:**
    1. Gọi `identity.FromContext(ctx)` để lấy `Subject` làm `owner_id`.
    2. Thực hiện Domain Validation: `Area > 0`. Trả về `errs.ErrValidation` nếu sai.
    3. Gọi Repository để lưu.

## 📋 Sub-tasks & Checklist
- [ ] Tạo Migration SQL cho bảng `farms`.
- [ ] Implement `internal/domain/farm.go` (Entity + Repository Interface).
- [ ] Implement `internal/infrastructure/postgres/farm_repo.go`.
- [ ] Implement `internal/usecase/create_farm.go` và `list_farms.go`.
- [ ] Viết unit test cho UseCase (Sử dụng Mock Repository).
