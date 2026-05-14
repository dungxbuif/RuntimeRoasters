# [RR-19.1] Tiếp nhận mẻ thu hoạch & Khởi tạo mã Batch

**User Story:**
Dưới vai trò là **Nhân viên Tiếp nhận (Intake Officer)**, tôi muốn hệ thống tự động nhận thông tin mẻ thu hoạch từ Farm để tôi có thể kiểm tra thực tế và cấp mã quản lý nhà máy mà không cần nhập liệu thủ công.

**Business Context:**
Đây là "cửa ngõ" đầu tiên của nhà máy. Sai sót ở bước này sẽ làm hỏng toàn bộ dữ liệu truy xuất nguồn gốc phía sau.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Lắng nghe Thông báo:** Hệ thống tự động hiển thị danh sách các mẻ thu hoạch (`Harvested`) đang được vận chuyển từ Farm đến.
2. **Kiểm tra Thực tế:** Khi xe hàng đến, nhân viên cân lại trọng lượng thực tế.
3. **Xác nhận Tiếp nhận:**
    - Nhân viên nhập trọng lượng thực tế nhận được.
    - Hệ thống so sánh với trọng lượng Farm khai báo (Cảnh báo nếu lệch > 2%).
4. **Cấp mã Batch ID:** Hệ thống tự động sinh mã `RR-P-{Origin}-{Date}-{Seq}` và in tem dán vào bao hàng.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Mã Batch:** Phải tham chiếu tuyệt đối tới `Harvest_ID` gốc.
- **Trạng thái:** Sau khi xác nhận, mẻ hàng chuyển sang trạng thái `RECEIVED`.
- **Hành động bắt buộc:** Phải chụp ảnh phiếu cân (Minh chứng dữ liệu).

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Tiếp nhận khớp dữ liệu
- **Given:** Có mẻ thu hoạch `H-001` nặng 100kg từ Farm.
- **When:** Tôi nhập cân nặng thực tế là 100kg và bấm "Confirm".
- **Then:** Hệ thống tạo Batch `RR-P-...-0001` trạng thái `RECEIVED`.

### Scenario 2: Cảnh báo lệch cân nặng
- **Given:** Mẻ thu hoạch khai báo 100kg.
- **When:** Tôi nhập cân nặng thực tế là 95kg.
- **Then:** Hệ thống hiển thị cảnh báo đỏ "Lệch quá 5% - Yêu cầu nhập lý do".
- **And:** Không cho phép Confirm nếu chưa nhập lý do.

---

## 📊 Yêu cầu Dữ liệu
- `Harvest_ID` (Mã tham chiếu)
- `Intake_Weight` (Cân nặng thực nhận)
- `Attachment_URL` (Link ảnh phiếu cân)
- `Note` (Ghi chú bất thường)
