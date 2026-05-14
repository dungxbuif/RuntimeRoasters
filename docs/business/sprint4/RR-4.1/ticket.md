# [RR-4.1] [BA] Khai báo mẻ thu hoạch (Harvesting Management)

**User Story:**
Dưới vai trò là **Nông dân (Farmer)**, tôi muốn khai báo khối lượng mẻ hạt vừa thu hoạch được để tôi có thể chuyển giao hàng cho nhà máy và nhận được thanh toán sau này.

**Business Context:**
Thu hoạch là thời điểm "hái ra tiền". Đây là dữ liệu đầu vào quan trọng nhất cho toàn bộ chuỗi cung ứng. Việc ghi nhận chính xác loại hạt, khối lượng và thời gian thu hoạch giúp đảm bảo tính tươi mới và chất lượng sản phẩm.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Lựa chọn Nông trại:** Farmer chọn nông trại mình đang quản lý.
2. **Nhập dữ liệu Thu hoạch:**
    - Khối lượng (kg).
    - Loại hạt (Arabica/Robusta/...).
    - Ngày giờ thu hoạch.
3. **Xác nhận:** Bấm "Submit Harvest".
4. **Thông báo:** Hệ thống cấp mã Harvest ID (ví dụ: `RR-H-CD-20260514-001`) và thông báo cho Warehouse.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Hạn mức (Quota):** Cảnh báo nếu khối lượng thu hoạch vượt quá năng lực dự kiến của nông trại (> 500kg/hecta/mẻ).
- **Trình tự thời gian:** Không cho phép khai báo thu hoạch cho các ngày trong tương lai.
- **Traceability ID:** Hệ thống tự động sinh mã Harvest ID theo chuẩn `docs/business/BATCH_LOGIC.md`.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Khai báo thành công
- **Given:** Tôi chọn "Nông trại Cầu Đất", loại "Arabica", nặng 100kg.
- **When:** Tôi nhấn gửi.
- **Then:** Hệ thống lưu mẻ thu hoạch và gửi thông báo "Harvest Created" tới các bộ phận liên quan.

### Scenario 2: Sai lệch dữ liệu (Số âm)
- **When:** Tôi nhập khối lượng -10kg.
- **Then:** Hệ thống hiển thị lỗi "Khối lượng phải là số dương".

---

## 📊 Yêu cầu Dữ liệu
- `Farm_ID` (Required)
- `Quantity` (Decimal, Required)
- `Harvest_Date` (Required)
- `Coffee_Type` (Enum, Required)
