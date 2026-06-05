# Technical Design - [RR-17] Flow 2: View, Update, Delete Farm

Mục tiêu: Hoàn thiện vòng đời CRUD kèm kiểm soát Data Scope.

## 🏗️ 1. Repository Enhancements
- **Methods:**
    - `Update(ctx, *domain.Farm) error`
    - `Delete(ctx, id string) error`
    - `GetByID(ctx, id string) (*domain.Farm, error)`

## 🧠 2. UseCase Security (Double Check)
- **Gate 2 Logic (UseCase Layer)**: Ngay cả khi Casbin cho phép vào API Update, UseCase vẫn phải kiểm tra Identity.
- **Workflow:**
    1. Lấy `CurrentUserID` từ Context.
    2. Gọi Repo `GetByID`.
    3. **Nếu `farm.OwnerID != CurrentUserID`** -> Trả về `errs.ErrForbidden`.

## 📋 Sub-tasks & Checklist
- [ ] Implement `GetByID`, `Update`, `Delete` tại Repository.
- [ ] Viết UseCase `UpdateFarm` và `DeleteFarm` kèm logic check Ownership.
- [ ] Map lỗi `ErrForbidden` sang gRPC status `PermissionDenied` trong Handler.
- [ ] Viết unit test cho kịch bản: Cập nhật thành công vs. Cập nhật của người khác (Fail).
