# [RR-18] Flow 3: Batch Management

- **Summary:** Quản lý các lô hàng (Batch) gắn liền với từng Nông trại.
- **Priority:** `MEDIUM`
- **Type:** Feature

---

## 🔍 Acceptance Criteria

### Scenario 1: Tạo Lô hàng mới cho Nông trại
- **Given:** Tôi đang xem chi tiết một Nông trại.
- **When:** Tôi tạo một Lô hàng (Batch) mới cho nông trại đó.
- **Then:** Lô hàng được lưu trữ với tham chiếu `farm_id` chính xác.

### Scenario 2: Kiểm tra ràng buộc Nông trại tồn tại
- **Given:** Tôi cố gắng tạo lô hàng cho một `farm_id` không tồn tại.
- **When:** Hệ thống thực hiện lưu trữ.
- **Then:** Hệ thống trả về lỗi "Nông trại không tồn tại".

---

## 👨‍💻 Developer Implementation Guide

### 1. Aggregate Root Pattern
- **Technique:** `Farm` là Aggregate Root cho `Batch`.
- **Constraint:** `Batch` không thể tồn tại độc lập. `farm_id` phải là `NOT NULL` và có Foreign Key tới bảng `farms`.

### 2. Transactional Outbox (Preparation)
- **Technique:** Sử dụng `database.WithTx` (Unit of Work) nếu có sẵn hoặc GORM `db.Transaction`.
- **Logic:** Khi lưu `Batch`, hãy thực hiện chèn record vào bảng `batches` VÀ `outbox` trong cùng 1 Transaction để đảm bảo tính nguyên tử.

### 3. Repository (`internal/infrastructure/postgres/batch_repo.go`)
- **Functions:**
  - `Create(ctx, *Batch) error`
  - `ListByFarm(ctx, farmID string) ([]*Batch, error)`

### 4. UseCase Logic
- **`CreateBatch`:**
  1. Kiểm tra sự tồn tại của `Farm` qua `FarmRepository`.
  2. Kiểm tra `Farm.OwnerID` để đảm bảo user có quyền tạo batch cho farm này.
  3. Lưu `Batch`.
