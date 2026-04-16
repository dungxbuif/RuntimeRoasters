# [RR-6] Farm Service - Business Logic & Delivery

- **Summary:** Hoàn thiện logic nghiệp vụ và các cổng giao tiếp (HTTP/gRPC).
- **Priority:** `MEDIUM`
- **Description:** Implement lớp UseCase để xử lý quy tắc nghiệp vụ và lớp Delivery để tiếp nhận yêu cầu từ UI/Gateway.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Kiểm tra nghiệp vụ khi tạo mẻ thu hoạch
- **Given:** Một yêu cầu tạo mẻ thu hoạch cho Nông trại X.
- **When:** Nông trại X không tồn tại trong hệ thống.
- **Then:** Hệ thống phải trả về lỗi "Farm Not Found" với mã lỗi 404 chuẩn RFC 7807.

### Scenario 2: Xử lý request REST API thành công
- **Given:** Request `POST /api/v1/farms` hợp lệ từ Client.
- **When:** Hệ thống xử lý xong.
- **Then:** Client phải nhận được HTTP 201 Created kèm theo payload JSON đầy đủ thông tin farm vừa tạo.

### Scenario 3: Tiêm phụ thuộc (Dependency Injection)
- **Given:** Toàn bộ code các lớp Repository và UseCase đã hoàn thiện.
- **When:** Khởi tạo service trong `main.go`.
- **Then:** Phải sử dụng Constructor pattern để inject Repository vào UseCase và UseCase vào Handler.

---

## 🛠️ Technical Notes
- Sử dụng Framework `Gin` cho HTTP layer.
- Tuân thủ nghiêm ngặt Dependency Inversion (Interface-driven).

## 📋 Sub-tasks
- [ ] Implement `internal/usecase/`.
- [ ] Implement `internal/delivery/http/`.
- [ ] Viết `main.go` cho Farm Service sử dụng `pkg/base`.
