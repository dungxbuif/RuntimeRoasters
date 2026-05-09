# [RR-6] Demo Service — Clean Architecture Boilerplate + Full Stack Vertical Slice

- **Goal:** Xây dựng apps/demo-service/ với Clean Architecture đầy đủ.
- **Business Value:** Tạo ra một mẫu chuẩn (canonical template) giúp tăng tốc độ phát triển các service tiếp theo và đảm bảo tính nhất quán trong kiến trúc.
- **Priority:** `HIGH`

## 🔍 Acceptance Criteria

### Scenario 1: Clean Architecture đầy đủ layers
- **Given:** `apps/demo-service/` đã được implement.
- **When:** Tôi mở cấu trúc thư mục.
- **Then:** Có đầy đủ `domain/`, `usecase/`, `infrastructure/`, `delivery/`. Dependency flow: `delivery → usecase (interface) → domain`. Infrastructure implements usecase interfaces. `go build ./...` passes.

### Scenario 2: Wire DI hoạt động
- **Given:** Demo Service sử dụng Google Wire.
- **When:** Tôi mở `cmd/main.go`.
- **Then:** File chỉ có: load config → `InitializeApp(cfg)` → `app.Run(...)`. Mọi dependency được wire tự động.

### Scenario 3: Kết nối infrastructure
- **Given:** Demo Service đang chạy với docker-compose up.
- **When:** Tôi gọi `GET /health/ready`.
- **Then:** Response JSON trả về status của PostgreSQL, Redis, và Kafka — tất cả `healthy`.

### Scenario 4: Demo endpoint qua Gateway
- **Given:** KrakenD và Demo Service đều đang chạy.
- **When:** Tôi gọi `GET localhost:8081/v1/demo/ping`.
- **Then:** Response JSON hợp lệ được trả về với `message` và `timestamp`.

### Scenario 5: Trace visible trong SigNoz
- **Given:** Demo Service đang chạy, SigNoz đang chạy.
- **When:** Tôi gọi `GET /v1/demo/ping`.
- **Then:** SigNoz UI hiển thị waterfall trace với các spans: HTTP, gRPC, DB query, Redis, Kafka produce — tất cả trong cùng 1 trace_id.

### Scenario 6: Kafka message có trace context
- **Given:** Demo Service produce Kafka message khi gọi ping endpoint.
- **When:** Consumer nhận message.
- **Then:** Kafka message headers có `traceparent` header theo W3C format — trace không bị đứt tại biên giới Kafka.

### Scenario 7: Log correlation
- **Given:** Demo Service đang xử lý request.
- **When:** Tôi đọc stdout log.
- **Then:** Mọi dòng log trong cùng request có cùng `trace_id` và `span_id` field trong JSON.