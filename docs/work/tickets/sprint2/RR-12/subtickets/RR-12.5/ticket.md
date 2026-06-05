# [RR-12.5] Integration: Demo Service Wiring & Middleware

- **Mục tiêu:** Áp dụng hệ thống phân quyền vào dịch vụ thực tế và kiểm thử khả năng bảo vệ API đầu cuối.
- **Mô tả:** Triển khai các thành phần thực thi quyền (Interceptors/Middleware) tại `pkg/base/casbin` và tích hợp chúng vào `demo-service` để hoàn tất luồng bảo mật.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. gRPC Interceptor được triển khai tại `pkg/base/casbin/interceptor_grpc.go`.
2. Gin Middleware được triển khai tại `pkg/base/casbin/middleware_gin.go`.
3. `demo-service` tích hợp thành công Resilient Engine và bảo vệ các endpoint bằng Interceptor/Middleware mới.
4. Kiểm thử thành công: User với Role hợp lệ (ví dụ: `farm_manager`) có quyền gọi API, User không hợp lệ bị từ chối với mã lỗi `PermissionDenied`.
5. Audit Log: Log đầy đủ thông tin về các yêu cầu bị từ chối truy cập.

## 🛠 Task list cho Developer
- [ ] Phát triển gRPC Unary Interceptor tự động ánh xạ RPC Method thành Casbin Object/Action.
- [ ] Phát triển Gin Middleware hỗ trợ các endpoint REST (nếu có).
- [ ] Cấu hình Dependency Injection (DI) trong `demo-service` để nạp Engine.
- [ ] Đăng ký Interceptor vào gRPC server của `demo-service`.
- [ ] Thực hiện End-to-End test với các kịch bản Role khác nhau.
