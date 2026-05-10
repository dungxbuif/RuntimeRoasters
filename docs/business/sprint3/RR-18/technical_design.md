# Technical Design - [RR-18] Flow 3: Batch Management

Mục tiêu: Quản lý Lô hàng (Batch) theo mô hình Aggregate Root với Nông trại.

## 🏗️ 1. Domain & Repository
- **Entity `Batch`:**
    ```go
    type Batch struct {
        ID          string    `gorm:"primaryKey;type:uuid"`
        FarmID      string    `gorm:"not null;index"`
        HarvestDate time.Time `gorm:"not null"`
        Quantity    float64   `gorm:"not null"`
        Status      string    `gorm:"default:'PENDING'"`
        CreatedAt   time.Time
        UpdatedAt   time.Time
    }
    ```
- **Aggregate Root Pattern:** Mọi thao tác với `Batch` phải thông qua kiểm tra tính hợp lệ của `Farm`.

## 🛠️ 2. Transactional Outbox (Preparation)
- **Logic:** Khi tạo `Batch`, cần chuẩn bị sẵn hạ tầng để ghi sự kiện vào bảng `outbox`.
- **Technique:** Repository nên hỗ trợ nhận vào một `gorm.DB` instance để thực thi nhiều câu lệnh trong cùng 1 transaction.

## 🧠 3. UseCase (Business Logic)
- **Logic `CreateBatch`:**
    1. Kiểm tra sự tồn tại của `Farm` (qua `FarmRepo.GetByID`).
    2. Kiểm tra quyền sở hữu `Farm` của User hiện tại.
    3. Thực hiện lưu `Batch` vào Database.
    4. (Dự phòng) Chèn bản ghi sự kiện `BatchCreated` vào bảng `outbox`.

## 📋 Sub-tasks
- [ ] Thiết kế bảng `batches` với Foreign Key tới `farms`.
- [ ] Triển khai `BatchRepository`.
- [ ] Viết UseCase `CreateBatch` có check Aggregate Root.
- [ ] Thiết lập helper Transaction cho GORM (VD: hàm `WithTx`).
