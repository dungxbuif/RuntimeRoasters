# Ticket RR-5.1: [BA] Quy trình chế biến mẻ hạt (Processing Flow)

**User Story:**
Dưới vai trò là một **Quản lý nhà máy (Roastery Manager)**, tôi muốn quản lý quy trình chuyển đổi mẻ hạt từ lúc nhận từ trang trại cho đến khi ra thành phẩm hạt rang, để đảm bảo chất lượng và khả năng truy vết.

**Business Value:**
Giai đoạn chế biến là nơi giá trị của hạt cà phê được nâng cao. Việc minh bạch hóa quy trình này giúp tăng niềm tin cho khách hàng và quản lý hiệu quả năng suất nhà máy.

---

## ✅ Acceptance Criteria (AC)

### Scenario 1: Tiếp nhận mẻ hạt tự động
- **Given:** Farm Service vừa phát hành một sự kiện thu hoạch mới.
- **Then:** Hệ thống Roastery tự động tạo một mẻ chế biến ở trạng thái `RECEIVED` kèm theo thông tin Farm gốc.

### Scenario 2: Cập nhật công đoạn sản xuất
- **Given:** Mẻ hạt đang ở trạng thái `RECEIVED`.
- **When:** Worker bắt đầu bóc vỏ (Hulling).
- **Then:** Trạng thái chuyển sang `HULLING`.
- **And:** Ghi lại thời gian và người thực hiện vào nhật ký (Logs).

### Scenario 3: Hoàn tất và cấp mã Batch ID
- **Given:** Mẻ hạt đang ở công đoạn Rang (Roasting).
- **When:** Quá trình rang kết thúc.
- **Then:** Trạng thái chuyển sang `COMPLETED`.
- **And:** Hệ thống tự động sinh mã `Batch ID` (VD: RR-A-20260515-X9K2).
- **And:** Mẻ hạt sẵn sàng để nhập Kho (Warehouse).

### Scenario 4: Vi phạm quy trình
- **Given:** Mẻ hạt đang ở trạng thái `RECEIVED`.
- **When:** Worker cố tình cập nhật thẳng lên `COMPLETED`.
- **Then:** Hệ thống từ chối và yêu cầu tuân thủ đúng trình tự (Hulling -> Drying -> Roasting).

---

## 📊 Dữ liệu yêu cầu
- `batch_id`: String (Generated)
- `worker_id`: UUID
- `processing_steps`: List of (Step, Timestamp, Note)

## 🔗 Liên kết kỹ thuật
- Xem Technical Design tại: [technical_design.md](../technical_design.md)
