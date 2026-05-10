# Technical Design - [RR-17] Flow 2: View, Update, Delete Farm

Mục tiêu: Hoàn thiện vòng đời quản lý nông trại (CRUD) với kiểm tra quyền sở hữu dữ liệu.

## 🏗️ 1. Domain & Repository
- **Interface (Bổ sung):**
    ```go
    type FarmRepository interface {
        // ... (các hàm cũ)
        GetByID(ctx context.Context, id string) (*Farm, error)
        Update(ctx context.Context, farm *Farm) error
        Delete(ctx context.Context, id string) error
    }
    ```

## 🛠️ 2. Infrastructure (Postgres)
- **File:** `internal/infrastructure/postgres/farm_repo.go` (Thêm các file `get.go`, `update.go`, `delete.go`)
- **Technique:** Hàm `Update` sử dụng `db.Save(farm)` của GORM.

## 🧠 3. UseCase (Business Logic)
- **Logic `UpdateFarm` (Data Scope Enforcement):**
    1. Gọi Repo `GetByID(ctx, req.ID)`. Nếu không thấy -> `errs.ErrNotFound`.
    2. Lấy `CurrentUserID` từ Identity Context.
    3. **KIỂM TRA:** `farm.OwnerID == CurrentUserID`.
    4. Nếu KHÁC: Trả về `errs.ErrForbidden` (Chống truy cập chéo dữ liệu).
    5. Nếu TRÙNG: Cập nhật các trường (Name, Location, Area, Type) và gọi Repo `Update`.

- **Logic `DeleteFarm`:** Tương tự `UpdateFarm`, phải kiểm tra quyền sở hữu trước khi thực thi xóa.

## 📡 4. Delivery (gRPC)
- **Hàm `GetFarm`:** Gọi UseCase, map lỗi `ErrForbidden` sang gRPC `PermissionDenied`.
- **Hàm `UpdateFarm`:** Nhận `UpdateFarmRequest`, gọi UseCase, trả về Entity đã cập nhật.
- **Hàm `DeleteFarm`:** Trả về boolean `success`.

## 📋 Sub-tasks
- [ ] Implement `GetByID`, `Update`, `Delete` trong Repository.
- [ ] Viết UseCase `GetFarm`, `UpdateFarm`, `DeleteFarm` kèm logic check Ownership.
- [ ] Hoàn thiện gRPC Handler cho các phương thức tương ứng.
