# [RR-17] Duy trì Tính Chính xác của Dữ liệu Nông trại (Maintaining Farm Accuracy)

- **Tóm tắt (Summary):** Triển khai các tính năng cập nhật thông tin và xóa bỏ các nông trại không còn hoạt động.
- **Độ ưu tiên (Priority):** `MEDIUM`
- **Loại (Type):** Feature

---

## 📖 Câu chuyện người dùng (User Story)
> Là một **Người nông dân (Farmer)**, tôi muốn có thể thay đổi thông tin nông trại (như đổi tên hoặc cập nhật lại diện tích thực tế) hoặc xóa bỏ nông trại nếu tôi không còn canh tác ở đó nữa, để đảm bảo hồ sơ của tôi trên hệ thống luôn đúng với thực tế.

## 💰 Giá trị nghiệp vụ (Business Value)
Đảm bảo tính chính xác cho các báo cáo sản lượng. Việc người dùng có thể tự quản lý dữ liệu giúp giảm tải cho bộ phận hỗ trợ kỹ thuật.

---

## 🔍 Điều kiện nghiệm thu (Acceptance Criteria)

### Kịch bản 1: Cập nhật thông tin thành công
- **Giả sử:** Tôi đang ở trang chỉnh sửa nông trại "Đồi Chè A".
- **Khi:** Tôi đổi tên thành "Đồi Cà Phê A" và nhấn "Cập nhật".
- **Thì:** Tên mới phải được lưu lại và hiển thị chính xác ở mọi nơi.

### Kịch bản 2: Bảo vệ quyền sở hữu khi cập nhật/xóa
- **Giả sử:** Có một nông trại ID là `X` thuộc về Farmer A.
- **Khi:** Farmer B cố tình gửi lệnh cập nhật hoặc xóa cho nông trại `X`.
- **Thì:** Hệ thống phải từ chối và thông báo "Bạn không có quyền thực hiện hành động này".

### Kịch bản 3: Xóa nông trại
- **Giả sử:** Tôi muốn xóa nông trại của mình.
- **Khi:** Tôi xác nhận xóa.
- **Thì:** Nông trại đó không còn xuất hiện trong danh sách của tôi nữa.
