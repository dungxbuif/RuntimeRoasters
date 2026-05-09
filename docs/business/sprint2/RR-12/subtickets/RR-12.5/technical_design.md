# RR-12.5: Tích hợp vào Demo Service

## 1. Mục tiêu (Goal)
Tích hợp Casbin Enforcer và gRPC Interceptor vào `demo-service` để thực thi phân quyền thực tế trên các API của service này.

## 2. Ngữ cảnh (Context)
- **Service:** `src/apps/demo-service`
- **File cần sửa:** `wire.go`, `cmd/main.go` hoặc file khởi tạo gRPC server.

## 3. Các bước triển khai (Step-by-Step Implementation)
1. **Cập nhật Dependency Injection (Wire):**
   - Khai báo Provider cho `casbin.Enforcer` trong `wire.go`.
   - Đảm bảo cấu hình đường dẫn file model/policy được lấy từ environment variables.
2. **Inject vào gRPC Server:**
   - Trong hàm khởi tạo gRPC Server, nhận tham số là `casbin.Enforcer`.
   - Thêm Interceptor: `grpc.UnaryInterceptor(casbin.NewInterceptor(enforcer))`.
3. **Thực thi Wire:** Chạy lệnh `make wire` hoặc `wire gen` để cập nhật code khởi tạo.
4. **Kiểm tra Log:** Đảm bảo khi service khởi chạy, Enforcer nạp policy thành công mà không có lỗi file not found.

## 4. Quy chuẩn tuân thủ (Patterns to Follow)
- **Dependency Injection:** Tuyệt đối không khởi tạo Enforcer trực tiếp bằng `new` trong code business. Sử dụng Wire để quản lý vòng đời.
- **Fail-fast:** Nếu không nạp được policy, service nên dừng ngay lập tức (panic/log.Fatal) để tránh việc chạy mà không có bảo vệ phân quyền.

## 5. Xác minh (Verification)
- **Scenario 1 (Thành công):** Sử dụng token có role `farmer`, gọi `CreateFarm` -> Kết quả trả về thành công.
- **Scenario 2 (Thất bại):** Sử dụng token có role `guest` hoặc không có role, gọi `CreateFarm` -> Kết quả trả về lỗi `403 Forbidden` (hoặc `PermissionDenied`).
- **Log Verification:** Kiểm tra log của `demo-service` để xem các dòng debug log của Casbin (nếu có) khi thực hiện check quyền.
