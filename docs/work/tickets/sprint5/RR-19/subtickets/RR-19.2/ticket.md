# [RR-19.2] Nhật ký Mẻ rang (Roast Runs) & Quản lý Gom lô

**User Story:**
Dưới vai trò là **Thợ rang (Roaster)**, tôi muốn ghi nhận kết quả của từng mẻ rang (Run) vào trong một **Lô sản xuất (Batch)** để hệ thống tự động cộng dồn sản lượng và tính toán hao hụt tổng thể.

**Business Context:**
Trong Specialty Coffee, máy rang thường có công suất nhỏ hơn so với nhu cầu một đơn hàng. Việc cho phép một Lô hàng (Batch) chứa nhiều mẻ rang (Run) giúp linh hoạt trong sản xuất nhưng vẫn đảm bảo tính đồng nhất chất lượng của cả lô.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Chọn Lô:** Thợ rang chọn một Batch đang ở trạng thái `PROCESSING`.
2. **Thực hiện Mẻ rang:** Chạy máy rang cho 1 lượng hạt thô (VD: 12kg).
3. **Ghi nhật ký (Run Log):** Nhập `Input_Weight` (thô) và `Output_Weight` (chín) của mẻ đó.
4. **Cộng dồn:** Hệ thống tự động cập nhật `Total_Output` hiện tại của cả Batch.
5. **Tiếp tục:** Thợ rang làm tiếp các mẻ khác cho đến khi đủ số lượng yêu cầu của Batch.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Mối quan hệ 1-N:** 1 Production Batch có thể chứa nhiều Roast Runs.
- **Tính toán tự động:** `% Loss (Batch) = (1 - Sum(Run_Outputs) / Sum(Run_Inputs)) * 100`.
- **Cảnh báo:** Nếu mẻ rang nào có hao hụt > 25%, hệ thống cảnh báo "Quality Deviation" cho quản lý.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Ghi nhận mẻ rang đầu tiên
- **Given:** Batch `A` mới chuyển sang `PROCESSING`.
- **When:** Tôi thêm Run #1 (In: 12kg, Out: 10kg).
- **Then:** Hệ thống lưu Run #1 và hiển thị Batch sản lượng là 10kg.

### Scenario 2: Gom nhiều mẻ vào 1 lô
- **Given:** Batch `A` đã có Run #1 (10kg chín).
- **When:** Tôi thêm Run #2 (In: 12kg, Out: 10.2kg).
- **Then:** Tổng sản lượng Batch `A` hiển thị là 20.2kg.

---

## 📊 Yêu cầu Dữ liệu
- `Run_ID` (Auto)
- `Batch_ID` (FK)
- `Input_Weight`
- `Output_Weight`
- `Roaster_Name`
