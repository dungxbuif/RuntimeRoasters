# [RR-19.2] Quản lý quy trình chế biến (Rang/Sấy) & Tính hao hụt

**User Story:**
Dưới vai trò là **Thợ rang (Roaster Operator)**, tôi muốn ghi lại từng bước chế biến và trọng lượng sau khi rang để hệ thống tự động tính toán hiệu suất và chất lượng của mẻ hàng.

**Business Context:**
Chế biến là giai đoạn hạt cà phê thay đổi giá trị nhiều nhất. Việc theo dõi hao hụt giúp quản lý biết được máy rang có đang hoạt động ổn định hay không và mẻ hàng có đạt chuẩn hay không.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Bắt đầu Chế biến:** Thợ rang chọn một Batch đang ở trạng thái `RECEIVED`.
2. **Cập nhật Công đoạn:**
    - Chuyển sang `HULLING` (Xát vỏ).
    - Chuyển sang `ROASTING` (Đang rang).
3. **Kết thúc & Cân đầu ra:**
    - Sau khi rang xong, thợ rang nhập `Yield Weight` (Trọng lượng đầu ra).
    - Hệ thống tính toán `% Hao hụt` = `(1 - Yield/Intake) * 100`.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Hao hụt Tiêu chuẩn:** Hệ thống mặc định dải hao hụt là 12% - 20%.
- **Cảnh báo Chất lượng:** Nếu hao hụt < 10% (Chưa chín/ẩm) hoặc > 25% (Cháy/Thất thoát), hệ thống đánh dấu mẻ hàng là `QUALITY_WARNING`.
- **Thứ tự Trạng thái:** Không cho phép nhảy thẳng từ `RECEIVED` lên `STOCKED` mà không qua bước `ROASTING`.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Rang đạt chuẩn
- **Given:** Mẻ hàng có `Intake Weight` = 100kg.
- **When:** Tôi nhập `Yield Weight` = 85kg (Hao hụt 15%).
- **Then:** Hệ thống cho phép hoàn tất và chuyển trạng thái sang `PROCESSED`.

### Scenario 2: Rang cháy (Hao hụt quá cao)
- **Given:** `Intake Weight` = 100kg.
- **When:** Tôi nhập `Yield Weight` = 70kg (Hao hụt 30%).
- **Then:** Hệ thống yêu cầu nhập "Anomaly Note" và gắn flag `QUALITY_WARNING` cho mẻ hàng này.

---

## 📊 Yêu cầu Dữ liệu
- `Batch_ID` (Mã tham chiếu)
- `Processing_Type` (Rang/Sấy)
- `Yield_Weight` (Trọng lượng đầu ra)
- `Anomaly_Note` (Ghi chú nếu hao hụt bất thường)
