# [RR-12.2] Auth Service: gRPC Snapshot API

- **Mục tiêu:** Cung cấp cơ chế cho phép các dịch vụ khác (Readers) tải toàn bộ chính sách phân quyền khi khởi động (Bootstrapping).
- **Mô tả:** Triển khai gRPC endpoint `GetPolicies` tại Auth Service. Đây là "Trụ cột 1" trong chiến lược Resilience, giúp Reader có đủ dữ liệu để hoạt động ngay cả khi Kafka chưa sẵn sàng.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. File proto được định nghĩa tại `api/runtime/auth/v1/auth.proto` với phương thức `GetPolicies`.
2. Response trả về danh sách đầy đủ các bản ghi chính sách (`p`) và phân cấp role (`g`).
3. Logic xử lý tại `auth-service` truy vấn trực tiếp từ cơ sở dữ liệu (thông qua adapter) để đảm bảo snapshot mới nhất.
4. Tích hợp gRPC server vào `auth-service` và lắng nghe trên port cấu hình (mặc định 50051).

## 🛠 Task list cho Developer
- [ ] Định nghĩa `AuthService.GetPolicies` trong file proto.
- [ ] Generate code gRPC (sử dụng Buf hoặc protoc).
- [ ] Triển khai Handler `GetPolicies` trong `internal/service/auth.go`.
- [ ] Viết script test nhanh bằng `grpcurl` để kiểm tra kết quả trả về.
