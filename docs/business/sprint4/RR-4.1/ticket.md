# Ticket RR-4.1: [BA] Khai báo mẻ thu hoạch thông minh (Digital Harvesting)

**User Story:**
Dưới vai trò là một **Nông dân (Farmer)**, tôi muốn khai báo mẻ thu hoạch cà phê mới trực tiếp trên hệ thống, để sản lượng của tôi được ghi nhận và bắt đầu hành trình truy xuất nguồn gốc minh bạch.

**Business Value:**
Đây là "điểm khởi đầu" của toàn bộ chuỗi cung ứng. Nếu không có dữ liệu này, chúng ta không thể chứng minh được hạt cà phê đến từ đâu, loại gì và chất lượng ra sao.

---

## ✅ Acceptance Criteria (AC)

### Scenario 1: Khai báo thành công
- **Given:** Tôi đã đăng nhập và có quyền sở hữu Nông trại "Cầu Đất 01".
- **When:** Tôi nhập `Loại cà phê: Arabica`, `Khối lượng: 500kg`, `Ngày hái: Today`.
- **Then:** Hệ thống lưu trữ thành công và trả về mã số mẻ hạt (Harvest ID).
- **And:** Trạng thái mẻ hạt được đánh dấu là `NEW`.

### Scenario 2: Khai báo sai khối lượng
- **Given:** Tôi đang ở màn hình khai báo.
- **When:** Tôi nhập khối lượng là `-10` hoặc `0`.
- **Then:** Hệ thống phải từ chối và hiển thị lỗi: "Khối lượng phải lớn hơn 0".

### Scenario 3: Khai báo cho nông trại không thuộc sở hữu
- **Given:** Tôi là Farmer A.
- **When:** Tôi cố tình gửi yêu cầu tạo harvest cho `farm_id` của Farmer B.
- **Then:** Hệ thống trả về lỗi `403 Forbidden` hoặc `Not Found` để đảm bảo an ninh dữ liệu.

---

## 📊 Dữ liệu yêu cầu
- `farm_id`: uint64 (Bắt buộc)
- `coffee_type`: String Enum (ARABICA, ROBUSTA, CHERRY, CULI)
- `quantity`: Decimal (kg)
- `harvest_date`: ISO8601 Date / Unix Timestamp
- `status`: String Enum (NEW, PROCESSING, COMPLETED) - Hệ thống tự gán `NEW` khi khởi tạo.

## 🔗 Liên kết kỹ thuật
- Xem Technical Design tại: [technical_design.md](../technical_design.md)
