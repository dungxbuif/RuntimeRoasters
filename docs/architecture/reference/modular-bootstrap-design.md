# Technical Design: Modular Bootstrap Pattern (Reusability)

Mục tiêu: Loại bỏ sự trùng lặp (Boilerplate) tại `apps/*/internal/app/app.go` và `cmd/main.go` bằng cách chuyển logic điều phối (Orchestration) vào gói `pkg/base`.

## 1. Vấn đề hiện tại
Hiện tại, mỗi khi tạo một service mới (như Farm Service), chúng ta phải copy-paste khoảng 80% code từ `demo-service/internal/app/app.go`. Các phần lặp lại bao gồm:
- Khởi tạo Casbin background sync.
- Đăng ký Auth Middleware cho HTTP.
- Thiết lập Readiness check (ping DB/Valkey).
- Cấu hình Swagger/Gateway.
- Logic Graceful Shutdown.

## 2. Giải pháp: Orchestrator Pattern (Facade)
Chúng ta sẽ nâng cấp `pkg/base` để nó không chỉ cung cấp các hàm "đăng ký" mà còn đóng vai trò là một **Bộ điều phối (Orchestrator)** hoàn chỉnh.

### Mẫu thiết kế lựa chọn: **Template Method (via Composition) + Functional Hooks**.

### Cấu trúc dự kiến trong `pkg/base`:

#### A. Định nghĩa `Dependencies` (Common Components)
Chuyển các thành phần mà service nào cũng dùng vào một struct tập trung trong `base`:
```go
type Dependencies struct {
    DB           *database.DB
    RDB          *valkeyclient.Client
    CasbinEngine casbin.Engine
    KeyProvider  provider.KeyProvider
}
```

#### B. Định nghĩa `ServiceRegistrar` (Interface)
Mỗi service chỉ cần cung cấp các thông tin định danh và handlers của riêng nó:
```go
type ServiceRegistrar interface {
    RegisterGRPC(server *grpc.Server)
    RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
    RegisterHTTP(router *gin.RouterGroup) // Tuỳ chọn
}
```

#### C. Hàm Bootstrap tập trung (`Launcher`)
Hàm này sẽ thực thi toàn bộ luồng "chuẩn" mà chúng ta đang làm thủ công trong `app.go`:
```go
func (a *App) Launch(deps Dependencies, registrar ServiceRegistrar) {
    // 1. Tự động chạy Casbin Sync nếu có Engine
    if deps.CasbinEngine != nil {
        if reader, ok := deps.CasbinEngine.(*casbin.ResilientReader); ok {
            reader.StartBackgroundSync(context.Background())
        }
    }

    // 2. Gọi logic đăng ký riêng của Service
    registrar.RegisterGRPC(a.grpcServer)
    a.RegisterGateway(registrar.RegisterGateway, a.config.GRPCPort)

    // 3. Tự động nhúng Auth Middleware vào mọi request qua Gateway (nếu cần)
    // Hoặc cung cấp một group đã được bọc Auth sẵn cho Service
    
    // 4. Tự động đăng ký Readiness check cho DB & Valkey
    a.RegisterReadiness(func() error {
        if deps.DB != nil { /* ping */ }
        if deps.RDB != nil { /* ping */ }
        return nil
    })

    // 5. Chạy server
    a.Run(...)
}
```

## 3. Lợi ích sau khi áp dụng

### Code tại `farm-service/internal/app/app.go` sẽ chỉ còn:
```go
func (a *App) Run() {
    a.Base.Launch(a.Deps, a.Handler) // Chỉ 1 dòng duy nhất!
}
```

### Code tại `cmd/main.go` (Sử dụng Manual DI):
Hầu như không thay đổi, nhưng logic khởi tạo sẽ sạch hơn vì các tham số truyền vào đã được đóng gói trong `Dependencies`. Mọi dependency được khởi tạo và truyền vào thủ công tại Composition Root.

## 4. Lộ trình thực hiện (Sau Sprint 3)
1. **Refactor `pkg/base`**: Thêm struct `Dependencies` và hàm `Launch`.
2. **Standardize `Config`**: Đảm bảo mọi service config đều tương thích để `base` có thể đọc được các tham số chung.
3. **Migration**: Chuyển `demo-service` sang dùng `base.Launch` để kiểm chứng. Sau đó áp dụng cho `farm-service`.

## 💡 Triết lý
"Base là xương sống, Service là thịt". Xương sống lo các chức năng sinh tồn (Security, Health, Sync), Thịt lo các chức năng nghiệp vụ (Business Logic).
