# [RR-16] Số hóa Tài sản Nông trại (Digitizing Farm Assets)

- **Tóm tắt (Summary):** Triển khai tính năng đăng ký mới và quản lý danh sách các nông trại trong hệ thống.
- **Độ ưu tiên (Priority):** `HIGH`
- **Loại (Type):** Feature

---

## 📖 Câu chuyện người dùng (User Story)
> Là một **Người nông dân (Farmer)**, tôi muốn có thể đăng ký các khu vực canh tác của mình lên hệ thống, kèm theo thông tin về vị trí và giống cà phê, để tôi có thể bắt đầu ghi nhật ký thu hoạch cho từng lô hàng sau này.

## 💰 Giá trị nghiệp vụ (Business Value)
Việc số hóa thông tin nông trại là tiền đề cho tính năng **Truy xuất nguồn gốc (Traceability)**. Nó giúp minh bạch hóa xuất xứ hạt cà phê cho người tiêu dùng cuối cùng.

---

## 🔍 Điều kiện nghiệm thu (Acceptance Criteria)

### Kịch bản 1: Đăng ký nông trại thành công
- **Giả sử:** Tôi là một Farmer đã đăng nhập.
- **Khi:** Tôi gửi thông tin nông trại (Tên, Địa chỉ, Diện tích, Giống cà phê).
- **Thì:** Hệ thống phải lưu trữ thông tin này và gán quyền sở hữu nông trại đó cho tài khoản của tôi.

### Kịch bản 2: Hiển thị danh sách nông trại cá nhân
- **Giả sử:** Tôi đã đăng ký 2 nông trại khác nhau.
- **Khi:** Tôi truy cập trang danh sách nông trại của mình.
- **Thì:** Tôi phải thấy chính xác 2 nông trại đó và **không thấy** bất kỳ nông trại nào của người khác.

### Kịch bản 3: Ràng buộc dữ liệu
- **Giả sử:** Tôi nhập diện tích nông trại bằng 0 hoặc số âm.
- **Khi:** Tôi nhấn "Lưu".
- **Thì:** Hệ thống phải từ chối và hiển thị thông báo lỗi "Diện tích phải lớn hơn 0".
