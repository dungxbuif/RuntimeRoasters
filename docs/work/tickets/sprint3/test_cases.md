# Sprint 3: Farm Service CRUD — Test Cases (Persona-based)

Tài liệu này liệt kê các bước kiểm thử cụ thể cho Sprint 3, tập trung vào tính năng Quản lý Nông trại và Phân quyền (Data Scoping).

---

## 👥 Persona 1: System Admin (`admin@runtimeroasters.com`)

### Test Case 1.1: Tạo Nông trại cho Manager (Direct Assignment)
- **Hành động**: Đăng nhập bằng Admin, truy cập `/manage/farms`, chọn "Provision New Farm".
- **Dữ liệu**: Tên: "Highland Arabica Node", Manager: "Manager Cầu Đất", Location: "Cầu Đất, Đà Lạt".
- **Mong đợi**: 
    - [ ] Farm được tạo thành công.
    - [ ] Admin thấy Farm này trong danh sách tổng.
    - [ ] `owner_id` của Farm khớp với ID của Manager Cầu Đất.

### Test Case 1.2: Xem danh sách tổng (Global Visibility)
- **Hành động**: Truy cập `/manage/farms`.
- **Mong đợi**:
    - [ ] Admin thấy toàn bộ các Farm hiện có trong hệ thống (kể cả farm của các manager khác nhau).

### Test Case 1.3: Xóa/Cập nhật Farm bất kỳ (Overriding Control)
- **Hành động**: Chọn 1 Farm thuộc sở hữu của Manager khác và thực hiện Edit/Delete.
- **Mong đợi**:
    - [ ] Thao tác thành công (Admin có quyền tối cao).

---

## 👥 Persona 2: Farm Manager (`manager.caudat@runtimeroasters.com`)

### Test Case 2.1: Tự tạo Nông trại (Self-Ownership)
- **Hành động**: Đăng nhập bằng Manager Cầu Đất, truy cập `/manage/farms`.
- **Mong đợi**:
    - [ ] Farm được tạo thành công.
    - [ ] Hệ thống tự động gán `owner_id` là chính Manager đang đăng nhập.

### Test Case 2.2: Phân quyền dữ liệu (Data Isolation - ABAC)
- **Hành động**: Manager Cầu Đất xem danh sách Farm.
- **Mong đợi**:
    - [ ] **CHỈ** thấy các Farm do mình quản lý.
    - [ ] Tuyệt đối không thấy Farm của Manager Buôn Ma Thuột hay Admin.

### Test Case 2.3: Chặn truy cập trái phép (Security Boundary)
- **Hành động**: Giả lập request (Postman/cURL) dùng Token của Manager Cầu Đất để `DELETE` một Farm ID thuộc về Manager khác.
- **Mong đợi**:
    - [ ] Phản hồi lỗi `404 Not Found` hoặc `403 Forbidden`.

---

## 👥 Persona 3: User vãng lai (Public / Unauthorized)

### Test Case 3.1: Xem Showcase công khai
- **Hành động**: Không đăng nhập, truy cập `/`.
- **Mong đợi**: 
    - [ ] Xem được Dashboard, Topology, và các biểu đồ Visualize.

### Test Case 3.2: Chặn truy cập trang Quản lý
- **Hành động**: Truy cập thẳng `/manage/users` khi chưa đăng nhập.
- **Mong đợi**:
    - [ ] Hệ thống hiển thị trang Login hoặc báo lỗi Unauthorized.

---

## 🛠 Hướng dẫn thực hiện Test thủ công
1. Chạy lệnh `task seed` để đảm bảo có đủ user mẫu.
2. Mở trình duyệt ẩn danh để test các Persona khác nhau.
3. Kiểm tra log của `farm-service` và `auth-service` để verify Trace ID.
