# [RR-18] Ghi nhật ký Thu hoạch (Harvest Logging)

- **Tóm tắt (Summary):** Triển khai tính năng khai báo sản lượng thu hoạch thực tế cho từng nông trại, tạo nguồn dữ liệu gốc cho chuỗi cung ứng.
- **Độ ưu tiên (Priority):** `HIGH`
- **Loại (Type):** Feature

---

## 📖 Câu chuyện người dùng (User Story)
> Là một **Quản lý nông trại (Farm Manager)**, tôi muốn có thể ghi lại khối lượng cà phê vừa thu hoạch được tại một nông trại cụ thể, để tôi có bằng chứng về sản lượng và sẵn sàng gửi hàng tới nhà máy chế biến.

## 💰 Giá trị nghiệp vụ (Business Value)
Đây là hành động **kích hoạt** chuỗi cung ứng. Dữ liệu thu hoạch là cơ sở để tính toán hiệu suất nông trại và là thông tin quan trọng nhất mà người tiêu dùng muốn xem khi quét mã QR.

---

## 🔍 Điều kiện nghiệm thu (Acceptance Criteria)

### Kịch bản 1: Khai báo thu hoạch hợp lệ
- **Giả sử:** Tôi là chủ sở hữu của "Nông trại Sơn La".
- **Khi:** Tôi gửi thông tin thu hoạch (Ngày hái, Khối lượng: 500kg).
- **Thì:** Hệ thống phải lưu bản ghi này và liên kết chính xác với "Nông trại Sơn La".

### Kịch bản 2: Ngăn chặn khai báo cho nông trại không thuộc sở hữu
- **Giả sử:** Tôi cố tình gửi lệnh thu hoạch cho nông trại của người khác.
- **Khi:** Hệ thống xử lý.
- **Thì:** Hệ thống phải từ chối và thông báo lỗi.

### Kịch bản 3: Thông tin bắt buộc
- **Giả sử:** Tôi để trống ngày thu hoạch hoặc khối lượng.
- **Khi:** Tôi nhấn "Lưu".
- **Thì:** Hệ thống phải báo lỗi "Thông tin thu hoạch không đầy đủ".
