# Ticket RR-6.1: [BA] Quản lý tồn kho thành phẩm

**User Story:**
Dưới vai trò là một **Thủ kho (Warehouse Keeper)**, tôi muốn quản lý chính xác lượng tồn kho của từng mẻ cà phê rang thành phẩm, để đảm bảo hệ thống bán lẻ luôn biết được lượng hàng sẵn có để phục vụ khách hàng.

**Business Value:**
Tồn kho là tài sản. Việc quản lý chính xác giúp tối ưu hóa dòng vốn, tránh tình trạng cháy hàng ảo hoặc tồn kho quá mức. Đây cũng là chốt chặn quan trọng nhất để thực hiện các cam kết bán hàng.

---

## ✅ Acceptance Criteria (AC)

### Scenario 1: Nhập kho từ nhà máy
- **Given:** Roastery Service vừa hoàn tất mẻ rang `RR-A-20260515-X9K2` với khối lượng 100kg.
- **When:** Hệ thống Warehouse nhận được thông báo.
- **Then:** Tồn kho của Batch này tăng thêm 100kg.
- **And:** `available_qty = 100`, `reserved_qty = 0`.

### Scenario 2: Kiểm tra tồn kho khả dụng
- **Given:** Batch A có `available: 50kg`, `reserved: 10kg`.
- **When:** Hệ thống Retail hỏi "Có thể bán 45kg không?".
- **Then:** Hệ thống trả lời "Có".
- **When:** Hệ thống Retail hỏi "Có thể bán 55kg không?".
- **Then:** Hệ thống trả lời "Không đủ hàng".

### Scenario 3: Theo dõi vị trí hàng
- **Given:** Tôi là thủ kho.
- **When:** Tôi xem chi tiết Batch A.
- **Then:** Hệ thống hiển thị vị trí kệ (ví dụ: Khu A - Kệ 04).

---

## 📊 Dữ liệu yêu cầu
- `batch_id`: Mã từ Roastery.
- `quantity`: Số lượng (kg/gói).
- `location`: Vị trí vật lý.

## 🔗 Liên kết kỹ thuật
- Xem Technical Design tại: [technical_design.md](../technical_design.md)
