# [RR-3] Core Framework - Base App & Error Handling

- **Summary:** Xây dựng hệ thống Bootstrap cho mọi Microservice và chuẩn hóa phản hồi lỗi.
- **Priority:** `HIGH`
- **Description:** Implement thư viện khung khởi chạy (App Lifecycle) và chuyển đổi lỗi nghiệp vụ sang chuẩn RFC 9457 (Problem Details).

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Khởi tạo Service chuẩn hóa
- **Given:** Toàn bộ config đã sẵn sàng.
- **When:** Tôi thực hiện khởi tạo service qua hàm `base.NewApp(cfg)`.
- **Then:** Hệ thống phải tự động setup Health Checks và Graceful Shutdown.

### Scenario 2: Thực hiện tắt Service an toàn (Graceful Shutdown)
- **Given:** Service đang hoạt động và xử lý requests.
- **When:** Tôi nhấn `Ctrl+C` gửi tín hiệu `SIGINT`.
- **Then:** Service phải dừng nhận request mới và hoàn tất xử lý in-flight requests trước khi thoát.

### Scenario 3: Trả về lỗi theo chuẩn RFC 9457
- **Given:** Xảy ra lỗi nghiệp vụ (ví dụ: Không tìm thấy ID).
- **When:** Service trả về lỗi cho Client.
- **Then:** Payload trả về phải chứa đầy đủ các trường: `type`, `title`, `status`, `detail`, `instance`.

---

## 🛠️ Technical Notes
- Thư viện `pkg/base` quản lý signal và wait group.
- Thư viện `pkg/errs` map các error Go sang struct Problem JSON.

## 📋 Sub-tasks
- [ ] Implement `pkg/base/app.go`.
- [ ] Implement `pkg/errs/problem.go`.
