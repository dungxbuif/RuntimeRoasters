# [RR-22] Trace Service & CQRS

- **Summary:** Thiết lập nền tảng cho Trace Service sử dụng mô hình CQRS để theo dõi vòng đời sản phẩm.
- **Priority:** `HIGH`
- **Description:** Xây dựng service chuyên biệt để lắng nghe và tổng hợp tất cả sự kiện từ Farm đến Retail, chuẩn bị dữ liệu cho Read Model.

---

## 🔍 Acceptance Criteria

### Scenario 1: Tổng hợp sự kiện theo Batch ID
- **Given:** Các sự kiện từ nhiều service khác nhau có cùng `batch_id`.
- **When:** Trace Service consume các sự kiện này.
- **Then:** Hệ thống lưu trữ chúng theo trình tự thời gian, cho phép truy vấn trạng thái hiện tại của mẻ hàng.

---

## 📋 Sub-tasks
- [ ] **Sub-task 22.1**: Thiết lập DB cho Trace Service (chuẩn bị dữ liệu để map sang Elasticsearch sau này).
- [ ] **Sub-task 22.2**: Lắng nghe toàn bộ event từ Farm, Processing, Retail, Warehouse để lưu vết trạng thái của từng `Batch_ID`.
