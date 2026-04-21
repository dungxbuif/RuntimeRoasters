# [RR-4-2] `pkg/base` — RegisterGateway & ServeSwagger

- **Summary:** Mở rộng base framework để hỗ trợ mount grpc-gateway HTTP mux và serve file OpenAPI.
- **Parent:** [RR-4](../ticket.md)
- **Priority:** `HIGH`

---

## 📖 User Story

> As a developer, I want the base `App` to support HTTP-to-gRPC transcoding out of the box, so that any service only needs to register a gateway mux and automatically gets REST endpoints without running a second process.

---

## 🔍 Acceptance Criteria

### Scenario 1: Mount grpc-gateway mux
- **Given:** Farm Service gọi `app.RegisterGateway(mux)` trong quá trình khởi tạo.
- **When:** Một HTTP request đến `/v1/farms/demo`.
- **Then:** Request được forward đến đúng gRPC handler — không trả về 404.

### Scenario 2: Serve Swagger JSON
- **Given:** Farm Service gọi `app.ServeSwagger("path/to/farm.swagger.json")`.
- **When:** Tôi gọi `GET /swagger/farm.json`.
- **Then:** Server trả về nội dung JSON hợp lệ của OpenAPI spec.

### Scenario 3: Không ảnh hưởng các route hiện có
- **Given:** `RegisterGateway()` đã được mount.
- **When:** Tôi gọi `GET /health/live`.
- **Then:** Health check vẫn trả về 200 bình thường.

---

## 🛠️ Technical Notes
- Thay đổi chỉ trong `pkg/base/app.go`.
- grpc-gateway mux mount tại `/v1/*` — không conflict với Gin routes hiện có (`/health/*`, `/swagger/*`).
- Library: `github.com/grpc-ecosystem/grpc-gateway/v2`.
