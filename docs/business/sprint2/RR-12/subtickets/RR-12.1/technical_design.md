# RR-12.1: Auth Service Bootstrap & Persistence Design

## 1. Mục tiêu (Goal)
Thiết lập dịch vụ trung tâm quản lý chính sách phân quyền (`auth-service`) với khả năng lưu trữ bền vững vào PostgreSQL thông qua GORM Adapter.

## 2. Cấu trúc thư mục (Folder Structure)
Dịch vụ sẽ được đặt tại `src/apps/auth-service/`:
```text
src/apps/auth-service/
├── cmd/
│   └── main.go             # Điểm khởi đầu của ứng dụng
├── config/
│   └── config.go           # Load env và cấu hình (DB, Kafka, gRPC)
└── internal/
    ├── domain/             # Business logic & Entities
    ├── repository/         # Gorm Adapter & DB persistence
    └── service/            # Orchestration layer
```

## 3. Cấu hình Casbin Persistence

### 3.1 Thư viện sử dụng
- Core: `github.com/casbin/casbin/v3`
- Adapter: `github.com/casbin/gorm-adapter/v3`

### 3.2 Model Definition (`internal/repository/casbin/model.conf`)
Sử dụng mô hình RBAC mạnh mẽ:
```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
# Cho phép Admin full quyền OR Kiểm tra Role Hierarchy + Object Match (Glob) + Action Match
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))
```

### 3.3 Gorm Adapter Setup
Trong `internal/repository/casbin/repo.go`:
1. Sử dụng kết nối DB từ `pkg/database/postgres.go`.
2. Khởi tạo adapter: `adapter, _ := gormadapter.NewAdapterByDB(db)`.
3. Tạo Enforcer: `enforcer, _ := casbin.NewEnforcer(modelPath, adapter)`.

## 4. Quy trình khởi tạo (Logic Flow)
1. **Load Config:** Đọc các biến môi trường (DB Host, Port, Credentials).
2. **Connect DB:** Khởi tạo kết nối Postgres toàn cục.
3. **Auto-Migrate:** Gorm Adapter sẽ tự động tạo bảng `casbin_rule` nếu chưa tồn tại.
4. **Load Policy:** `enforcer.LoadPolicy()` để nạp toàn bộ policy từ DB vào memory của Auth Service (dùng cho các API write sau này).

## 5. Lưu ý cho Developer
- Sử dụng `go:embed` để đóng gói `model.conf` vào file binary, tránh lỗi thiếu file khi deploy.
- Đảm bảo `CasbinRule` table được tạo đúng schema trong database `runtime_roasters`.
