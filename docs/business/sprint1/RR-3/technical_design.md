# Technical Plan - [RR-3] Core Framework - Base App & Error Handling

## 🎯 Chiến lược triển khai (Strategy)
- Xây dựng base app quản lý vòng đời ứng dụng.
- Chuẩn hóa format lỗi theo RFC 9457.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Implement `pkg/base/app.go`.
    - Quản lý signal và wait group cho graceful shutdown.
    - Setup cơ chế Health Check cơ bản.
- [ ] Implement `pkg/errs/problem.go` map các error Go sang struct Problem JSON.

## 🧪 Xác minh (Verification)
- [ ] Kiểm tra cơ chế Graceful Shutdown bằng lệnh kill.
- [ ] Mock lỗi và kiểm tra response body có đúng chuẩn RFC 9457.