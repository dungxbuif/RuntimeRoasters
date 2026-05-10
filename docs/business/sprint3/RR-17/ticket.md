# [RR-17] Flow 2: View, Update, Delete Farm

- **Summary:** Triển khai các tính năng quản lý chi tiết, cập nhật và xóa Nông trại.
- **Priority:** `MEDIUM`
- **Type:** Feature

---

## 🔍 Acceptance Criteria

### Scenario 1: Cập nhật thông tin Nông trại
- **Given:** Tôi đang ở trang chi tiết nông trại.
- **When:** Tôi thay đổi thông tin và nhấn "Cập nhật".
- **Then:** Thông tin mới được lưu lại và phiên bản (version) của bản ghi được tăng lên.

### Scenario 2: Xử lý xung đột khi cập nhật (Optimistic Locking)
- **Given:** Hai người dùng cùng mở trang chỉnh sửa một nông trại.
- **When:** Người thứ nhất lưu thành công, sau đó người thứ hai nhấn lưu.
- **Then:** Người thứ hai phải nhận được thông báo lỗi "Dữ liệu đã bị thay đổi bởi người khác" (ErrOptimisticLock).

### Scenario 3: Xóa Nông trại
- **Given:** Tôi muốn ngừng quản lý một nông trại.
- **When:** Tôi nhấn "Xóa" và xác nhận.
- **Then:** Nông trại bị xóa khỏi danh sách (hoặc đánh dấu Soft Delete).

---

### 1. Repository Technique
- **Update Function:**
```go
// postgres/farm_repo.go
func (r *farmRepo) Update(ctx, farm *domain.Farm) error {
    return r.db.WithContext(ctx).Save(farm).Error
}
```


### 2. UseCase Logic (Data Scope Enforcement)
- **`UpdateFarm` Logic:**
  1. Gọi Repo `GetByID(ctx, id)`.
  2. KIỂM TRA: `farm.OwnerID == identity.FromContext(ctx).Subject`.
  3. NẾU SAI: Return `errs.ErrForbidden`. (Cấm sửa hàng của người khác).
  4. NẾU ĐÚNG: Cập nhật các trường và gọi Repo `Update`.

### 3. Error Handling
- Map `errs.ErrConflict` sang gRPC status `codes.Aborted`.
- Map `errs.ErrForbidden` sang gRPC status `codes.PermissionDenied`.
- Sử dụng helper tại `pkg/errs` để chuẩn hóa lỗi RFC 9457.
