# [RR-4-3] Dependency Injection — Google Wire scaffold

- **Summary:** Thiết lập Google Wire để tự động hóa việc wiring dependencies trong Farm Service.
- **Parent:** [RR-4](../ticket.md)
- **Priority:** `HIGH`

---

## 📖 User Story

> As a developer, I want dependencies to be wired automatically via Wire, so that `main.go` only calls one function and adding a new layer (e.g. repository) requires no changes to the composition root.

---

## 🔍 Acceptance Criteria

### Scenario 1: main.go sạch sau khi Wire
- **Given:** Wire scaffold đã được setup.
- **When:** Tôi mở `farm-service/cmd/main.go`.
- **Then:** File chỉ chứa: load config → gọi `InitializeApp(cfg)` → gọi `app.Run(...)`. Không có khởi tạo thủ công nào khác.

### Scenario 2: Thêm dependency mới không cần sửa main.go
- **Given:** Một `FarmRepository` mới được tạo ở RR-5.
- **When:** Provider của `FarmRepository` được thêm vào `ProviderSet`.
- **Then:** Chạy `wire gen` → `wire_gen.go` cập nhật tự động, `main.go` không thay đổi.

### Scenario 3: Build thành công sau wire gen
- **Given:** `wire gen` đã chạy xong.
- **When:** Tôi chạy `go build ./...`.
- **Then:** Toàn bộ module biên dịch không có lỗi.

---

## 🛠️ Technical Notes
- Tool: `github.com/google/wire`.
- Scope của ticket này: chỉ wire `Config → DB → Redis → App`. Provider cho Repository và UseCase là placeholder cho RR-5, RR-6.
- File `wire_gen.go` được checked-in vào repo.
