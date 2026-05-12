# Technical Design - [RR-15] Farm Service Bootstrapping

Mục tiêu: Thiết lập hạ tầng code và nối dây (wiring) thủ công cho Farm Service, tuân thủ tiêu chuẩn [Canonical Template Standards](../../../architecture/reference/canonical-template-standards.md).

## 📂 1. Cấu trúc thư mục (Vertical Slice Preparation)
- `src/apps/farm-service/cmd/main.go`: Khởi tạo và chạy `application.Run()`.
- `src/apps/farm-service/internal/app/init.go`: **Manual DI Implementation**. Khởi tạo DB, Valkey, Casbin Engine, JWKS Provider và các Handlers.
- `api/runtime/farm/v1/farm.proto`: Định nghĩa contract API CRUD.

## ⚙️ 2. Manual DI & Interceptor Strategy
- **Security Gate 1 (AuthN)**: Nhúng `authgrpc.GRPCUnaryInterceptor` vào gRPC Server.
- **Security Gate 2 (AuthZ)**: Nhúng `casbingrpc.GRPCUnaryInterceptor` sử dụng `ResilientReader` kết nối tới `auth-service`.
- **Logging**: Khởi tạo Global Logger trong `init.go` và sử dụng `logger.FromContext(ctx)` trong các tầng xử lý lỗi.

## 📋 Sub-tasks & Checklist
- [ ] Thiết kế `farm.proto` với đầy đủ message CRUD.
- [ ] Implement hàm `InitializeApp()` trong `internal/app/init.go`.
- [ ] Đăng ký gRPC service `FarmService` vào bộ khung `base.App`.
- [ ] Cấu hình port `8083` (HTTP) và `50053` (gRPC) trong file `.env`.
- [ ] **Validation:** Chạy `go build ./apps/farm-service/...` thành công.
