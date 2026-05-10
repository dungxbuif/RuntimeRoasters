# Technical Design - [RR-15] Farm Service Bootstrapping

Mục tiêu: Thiết lập hạ tầng code và nối dây (wiring) cho Farm Service.

## 📂 1. Cấu trúc thư mục (Folder Structure)
- `src/apps/farm-service/cmd/main.go`
- `src/apps/farm-service/config/config.go`
- `src/apps/farm-service/internal/app/app.go`
- `src/apps/farm-service/internal/app/wire.go`
- `src/apps/farm-service/internal/delivery/grpc/`
- `src/apps/farm-service/internal/usecase/`
- `src/apps/farm-service/internal/infrastructure/postgres/`

## ⚙️ 2. Cấu hình (Config)
- **File:** `config/config.go`
- **Struct:** `Config` nhúng `config.BaseConfig`.
- **Fields cần thêm:** 
    - `AuthServiceAddr string` (Cổng gRPC của auth-service)
    - `JWKSURL string` (URL tới identity proxy)
    - `ExpectedIssuer string`

## 🔗 3. Dependency Injection (Wire)
- **File:** `internal/app/wire.go`
- **Providers yêu cầu:**
    - `provideKeyProvider`: Khởi tạo `provider.NewJWKSCache`.
    - `provideCasbinClient`: Khởi tạo `casbingrpc.NewAuthSnapshotClient`.
    - `provideCasbinEngine`: Khởi tạo `casbin.NewResilientReader` (ModelText sử dụng RBAC đơn giản).
    - `provideGRPCServerOptions`: Inject 2 interceptors:
        1. `authgrpc.GRPCUnaryInterceptor(keyProvider, cfg.ExpectedIssuer)`
        2. `casbingrpc.GRPCUnaryInterceptor(casbinEngine)`

## 🚀 4. App Lifecycle
- **File:** `internal/app/app.go`
- **Hàm `Run()`:**
    - Phải gọi `reader.StartBackgroundSync(ctx)` để đồng bộ quyền từ auth-service.
    - Đăng ký gRPC service và Gateway mux.

## 📋 Sub-tasks
- [ ] Copy boilerplate từ `demo-service`.
- [ ] Đổi toàn bộ chuỗi "demo-service" và "Demo" thành "farm-service" và "Farm" (bao gồm import paths).
- [ ] Chạy `wire` tại thư mục `internal/app`.
- [ ] Đăng ký port `8083` (HTTP) và `50053` (gRPC) trong `.env`.
