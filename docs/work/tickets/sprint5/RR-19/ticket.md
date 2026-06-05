# [RR-19] [BA] Quản lý Kho & Chuỗi giá trị (Intake to Batch)

**User Story:**
Dưới vai trò là **Quản lý Kho (Warehouse Manager)**, tôi muốn kiểm soát toàn bộ vòng đời của hàng hóa từ khi rời nông trại đến khi trở thành thành phẩm đóng gói, để tôi có thể đảm bảo số lượng tồn kho và khả năng truy xuất nguồn gốc.

**Business Context:**
Trong mô hình demo, chúng ta đơn giản hóa việc sản xuất nhưng vẫn giữ chặt chẽ tính nguyên tắc của kho bãi. Một **Lô hàng (Batch)** bán ra thị trường có thể được tạo thành từ nhiều **Mẻ rang (Roast Runs)** nhỏ để tối ưu công suất máy móc.

---

## 🔄 Luồng Nghiệp vụ (Workflow)

1.  **Nhập kho (Intake):**
    - Hệ thống nhận `Harvest_Created` event.
    - Tạo bản ghi `Intake_Batch` (Trạng thái: `RECEIVED`).
2.  **Chế biến (Processing - High Level):**
    - Thủ kho chuyển trạng thái Batch sang `PROCESSING`.
    - Thợ rang thực hiện N mẻ rang (Run). Với mỗi mẻ, chỉ ghi nhận: `Input_Weight` và `Output_Weight`.
3.  **Hoàn tất & Cấp mã (Finalize & Stocking):**
    - Khi đủ số lượng hoặc kết thúc ca, Quản lý bấm "Finalize Batch".
    - Hệ thống tổng hợp: `Total_Output = Sum(Run_Outputs)`.
    - Cấp mã **Production Batch ID** (theo chuẩn thương mại).
    - Cập nhật tồn kho thành phẩm (Roasted Coffee SKU).

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)

- **Mối quan hệ 1-N:** 1 Intake/Production Batch chứa 0..N Roast Runs.
- **Tính toán hao hụt:** Hệ thống tự động tính `% Loss` tổng thể của cả Batch dựa trên tổng Input và tổng Output.
- **Trạng thái hợp lệ:** Chỉ được phép Finalize Batch khi đã có ít nhất 1 Roast Run được ghi nhận.
- **Traceability:** Phải truy xuất được từ Production Batch ID ngược về mã Harvest ID ban đầu.

---

## ✅ Acceptance Criteria (AC)

### Scenario 1: Tiếp nhận hàng từ Farm
- **Given:** Farm Service vừa gửi event mẻ thu hoạch 500kg.
- **When:** Hệ thống xử lý event.
- **Then:** Một bản ghi "Lô chờ chế biến" xuất hiện trong kho với khối lượng 500kg thô.

### Scenario 2: Ghi nhận mẻ rang (Roast Run)
- **Given:** Một Batch đang ở trạng thái `PROCESSING`.
- **When:** Thợ rang nhập: Run #1 - In: 12kg, Out: 10kg.
- **Then:** Hệ thống lưu nhật ký mẻ rang và hiển thị tổng sản lượng hiện tại là 10kg.

### Scenario 3: Đóng lô thành phẩm
- **When:** Quản lý chọn "Finalize Batch".
- **Then:**
    - Trạng thái Batch đổi thành `STOCKED`.
    - Hệ thống sinh mã `PROD-ARABICA-20260514-001`.
    - Tồn kho SKU tương ứng tăng lên.
