# [RR-19.3] Nhập kho thành phẩm & Cập nhật tồn kho tự động

**User Story:**
Dưới vai trò là **Thủ kho (Warehouse Keeper)**, tôi muốn hệ thống tự động cập nhật số lượng tồn kho ngay khi mẻ rang hoàn tất, để bộ phận bán hàng biết chính xác lượng hàng sẵn có.

**Business Context:**
Đây là bước cuối cùng trong chuỗi cung ứng nội bộ. Dữ liệu ở đây sẽ được phơi bày trực tiếp cho khách hàng (Retail), vì vậy tính chính xác là tuyệt đối.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Hoàn tất Đóng gói:** Sau khi rang và để nguội, mẻ hàng được đóng thành các gói/bao theo chuẩn.
2. **Xác nhận Nhập kho:**
    - Thủ kho kiểm tra lần cuối số lượng bao và mã Batch.
    - Bấm "Move to Stock".
3. **Cập nhật Tồn kho:**
    - Hệ thống chuyển mã từ `RR-P-...` (Processing) sang `RR-S-...` (Stocked).
    - Tăng số lượng `Available Quantity` trong bảng Inventory cho loại cà phê/vùng tương ứng.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **FIFO:** Hệ thống phải gắn nhãn ngày nhập kho để sau này gợi ý xuất kho theo thứ tự cũ trước mới sau.
- **Stock Limit:** Cảnh báo nếu tổng tồn kho vượt quá dung tích kho vật lý (Cấu hình trong setting).

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Nhập kho thành công
- **Given:** Mẻ hàng `RR-P-CD-001` rang xong đạt 85kg.
- **When:** Tôi xác nhận nhập kho.
- **Then:** Trạng thái mẻ hàng thành `STOCKED`.
- **And:** Tồn kho cà phê Cầu Đất (Arabica) tăng thêm 85kg.

---

## 📊 Yêu cầu Dữ liệu
- `Batch_ID`
- `Warehouse_Slot` (Vị trí kệ kho)
- `Stock_In_Date` (Mặc định là thời điểm xác nhận)
