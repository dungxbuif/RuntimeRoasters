# [RR-8] UI Integration - Farm & Batch Management

- **Summary:** Kết nối UI Dashboard với Farm Service API.
- **Priority:** `HIGH`
- **Description:** Thực hiện luồng tương tác thực tế từ giao diện web xuống Backend để quản lý dữ liệu.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Hiển thị danh sách Farm
- **Given:** Có dữ liệu Farm trong Backend.
- **When:** Tôi truy cập trang Quản lý Nông trại.
- **Then:** Một bảng hiển thị danh sách tất cả các Farm phải hiện ra với đầy đủ thông tin tên và vị trí.

### Scenario 2: Tạo mẻ thu hoạch mới từ UI
- **Given:** Người dùng đã chọn một Nông trại.
- **When:** Người dùng điền form "Tạo mẻ" và nhấn "Lưu".
- **Then:** Một request POST phải được gửi sang Farm Service và dữ liệu sau đó phải được hiển thị mới lại trên UI.

### Scenario 3: Xử lý lỗi từ UI
- **Given:** Backend trả về lỗi (ví dụ: Tên nông trại quá ngắn).
- **When:** UI nhận được lỗi chuẩn RFC 7807.
- **Then:** UI phải hiển thị thông báo lỗi chi tiết cho người dùng một cách thân thiện.

---

## 🛠️ Technical Notes
- Sử dụng Custom Hook (`useFarms`, `useCreateBatch`).
- Sử dụng Axios interceptor để log trace ID.

## 📋 Sub-tasks
- [ ] Xây dựng trang Quản lý Nông trại.
- [ ] Xây dựng Form tạo Mẻ thu hoạch.
- [ ] Viết API client hooks sử dụng React Query.
