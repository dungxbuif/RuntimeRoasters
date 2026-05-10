# Technical Design - [RR-18] Flow 3: Harvest Management

Mục tiêu: Quản lý Lô thu hoạch (Harvest) theo mô hình Aggregate Root.

## 🏗️ 1. Database Schema (PostgreSQL)
Bảng `harvests`:
```sql
CREATE TABLE harvests (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    harvest_date DATE NOT NULL,
    quantity DECIMAL(12,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'HARVESTED',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## 🧠 2. Aggregate Root Pattern (Security)
- **Constraint**: Không thể tạo `Harvest` nếu không kiểm tra quyền sở hữu `Farm`.
- **UseCase `RecordHarvest`**:
    1. Kiểm tra `FarmID` tồn tại và thuộc về User đang login.
    2. Nếu hợp lệ -> Ghi `harvests` record.

## 🛠️ 3. Transactional Outbox (Preparation)
- **Kỹ thuật**: Mọi hàm Repository ghi `harvests` phải chấp nhận `*gorm.DB` (để hỗ trợ transaction).
- **Mục tiêu**: Sau này khi tạo Harvest, sẽ đồng thời ghi một Event vào bảng `outbox` để bắn lên Kafka.

## 📋 Sub-tasks & Checklist
- [ ] Tạo bảng `harvests`.
- [ ] Implement `domain.Harvest` và `HarvestRepository`.
- [ ] Xây dựng UseCase `RecordHarvest` (kèm logic check Ownership của Farm).
- [ ] Viết unit test cho Aggregate Root logic (Chủ nông trại X không hái được cho nông trại Y).
