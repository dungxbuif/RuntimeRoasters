# [RR-17] Farm UseCase & API (gRPC/REST)

- **Summary:** Triển khai logic nghiệp vụ và cung cấp giao diện lập trình ứng dụng (API) cho các dịch vụ khác.
- **Priority:** `HIGH`
- **Type:** Feature

---

## 📖 User Story
> As a system integrator, I want to access farm data via standard APIs so that I can build mobile apps, web dashboards, or integrate other services with the farm management system.

## 💰 Business Value
Mở rộng khả năng tương tác của hệ thống qua các giao thức chuẩn. Đảm bảo logic nghiệp vụ (như kiểm tra giới hạn diện tích, quyền sở hữu) được thực thi nhất quán.

## 🔍 Acceptance Criteria

### Scenario 1: gRPC Interface
- **Given:** A technical client (another service).
- **When:** Calling `CreateFarm` or `GetFarm` via gRPC.
- **Then:** The system responds with Protobuf-encoded data.

### Scenario 2: REST Interface (via Gateway)
- **Given:** A web or mobile client.
- **When:** Calling `POST /v1/farms` or `GET /v1/farms/{id}`.
- **Then:** The system responds with JSON-encoded data.

### Scenario 3: Validation Rules
- **Given:** A farm creation request with negative area.
- **When:** The UseCase processes the request.
- **Then:** It returns a validation error "Area must be positive".

### Scenario 4: Authorization Enforcement
- **Given:** A user without proper roles.
- **When:** Accessing farm management APIs.
- **Then:** The system returns `403 Forbidden` (Gate 2 check).
