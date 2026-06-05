# [RR-19.3] Đóng lô & Nhập kho Thành phẩm SKU

**User Story:**
Dưới vai trò là **Quản lý Kho (Warehouse Manager)**, tôi muốn thực hiện chốt sổ Lô sản xuất (Finalize Batch) để hệ thống tự động sinh mã Production Batch ID chính thức và cập nhật số lượng tồn kho thành phẩm.

**Business Context:**
Khi toàn bộ các mẻ rang nhỏ hoàn tất, hàng hóa được đóng bao và dán nhãn Batch thương mại để xuất bán cho Retail. Việc Finalize sẽ "khóa" dữ liệu sản xuất và chính thức đưa hàng vào kho.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Yêu cầu Đóng lô:** Quản lý chọn Batch đã hoàn tất rang.
2. **Xác nhận số lượng:** Hệ thống hiển thị tổng sản lượng (Sum of Runs). Quản lý bấm "Finalize".
3. **Sinh mã Batch ID:** Hệ thống cấp mã chính thức (VD: `PROD-CD-20260514-001`).
4. **Nhập kho (Auto Stock-in):**
    - Trạng thái Batch thành `STOCKED`.
    - Tồn kho SKU tương ứng (VD: `CD-ARABICA-ROASTED`) tăng thêm đúng bằng sản lượng Batch.
5. **Thông báo:** Bắn event `Warehouse_Stock_Updated` qua Kafka.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Tính bất biến:** Sau khi Finalize, không được phép thêm Roast Run vào Batch đó nữa.
- **Traceability:** Mã Production Batch ID phải ánh xạ được về danh sách các Run IDs đã tạo ra nó.
- **SKU Mapping:** Tự động xác định SKU dựa trên `Coffee_Type` và `Origin_Code`.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Đóng lô và sinh mã ID
- **Given:** Batch đang có 50kg thành phẩm từ 5 mẻ rang.
- **When:** Tôi nhấn "Finalize".
- **Then:** Hệ thống sinh mã `PROD-AR-20260514-001`.
- **And:** Batch chuyển sang `STOCKED`.

### Scenario 2: Cập nhật kho tự động
- **When:** Lô hàng 50kg được chốt.
- **Then:** Số lượng `Available Quantity` trong kho tăng thêm 50kg ngay lập tức.

---

## 📊 Yêu cầu Dữ liệu
- `Production_Batch_ID` (Generated)
- `Final_Yield` (Auto-sum)
- `Stock_Location` (Vị trí kho)
