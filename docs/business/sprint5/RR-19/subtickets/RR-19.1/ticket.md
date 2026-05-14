# [RR-19.1] Tiếp nhận mẻ thu hoạch & Khởi tạo Lô sản xuất (Batch)

**User Story:**
Dưới vai trò là **Quản lý Kho (Warehouse Manager)**, tôi muốn hệ thống tự động nhận thông tin mẻ thu hoạch từ Farm để khởi tạo một **Lô sản xuất (Production Batch)**, làm thực thể gốc để gom các mẻ rang sau này.

**Business Context:**
Đây là bước "mở sổ" cho một đợt sản xuất. Một Lô hàng thương mại (Batch) sẽ được liên kết trực tiếp với một đợt thu hoạch từ Farm để đảm bảo tính truy xuất.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Lắng nghe Thông báo:** Hệ thống tự động nhận event `Harvest_Created` từ Kafka.
2. **Khởi tạo Lô (Auto):** Hệ thống tự động tạo một bản ghi `Production Batch` ở trạng thái `RECEIVED`.
3. **Cân thực tế (Intake):** Nhân viên nhập trọng lượng thực nhận tại kho.
4. **Xác nhận (Confirm):**
    - Hệ thống so sánh với trọng lượng Farm khai báo.
    - Nếu lệch > 2%, yêu cầu nhập lý do.
5. **Cấp mã Batch:** Hệ thống sinh mã thương mại ban đầu cho Lô hàng.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Mối quan hệ:** 1 Production Batch tham chiếu tới 1 Harvest ID.
- **Trạng thái:** Mặc định ban đầu là `RECEIVED`.
- **Hao hụt vận chuyển:** Được ghi nhận ngay tại bước này nếu có sai lệch lớn.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Tự động tạo Lô khi có Harvest
- **Given:** Farm Service gửi event thu hoạch 500kg.
- **When:** Warehouse Service nhận event.
- **Then:** Một `Production Batch` mới được tạo với `Status = RECEIVED`.

### Scenario 2: Xác nhận cân nặng thực tế
- **When:** Nhân viên nhập cân nặng 495kg (Lệch < 2%).
- **Then:** Hệ thống cho phép "Confirm" và chuyển sang giai đoạn chờ chế biến.

---

## 📊 Yêu cầu Dữ liệu
- `Harvest_ID` (FK)
- `Intake_Weight` (Thực tế)
- `Intake_Note` (Nếu có hao hụt)
