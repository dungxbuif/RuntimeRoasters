# [RR-4] API Contract - Farm Service Protobuf & Swagger

- **Summary:** Định nghĩa hợp đồng API (API Contract) cho Farm Service.
- **Priority:** `MEDIUM`
- **Description:** Xây dựng định nghĩa gRPC và REST (via Gateway) cho các thực thể Nông trại (Farm) và Mẻ thu hoạch (Batch).

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Định nghĩa Farm Entity
- **Given:** Protobuf đã được khai báo.
- **When:** Tôi thực hiện biên dịch gRPC.
- **Then:** Các struct `Farm` và các method `CreateFarm`, `GetFarm` phải có sẵn cho Go service.

### Scenario 2: Tự động hóa OpenAPI (Swagger)
- **Given:** File `.proto` đã có chú thích chuẩn REST.
- **When:** Tôi chạy script biên dịch.
- **Then:** Hệ thống phải tạo ra file `farm.swagger.json` cho Dashboard UI.

### Scenario 3: Tính đồng bộ giữa UI và Backend
- **Given:** Frontend dùng file Swagger vừa tạo.
- **When:** Frontend gọi API theo định nghĩa.
- **Then:** Backend phải trả về dữ liệu đúng kiểu và cấu trúc như đã cam kết.

---

## 🛠️ Technical Notes
- Sử dụng Protobuf v3.
- Tích hợp `protoc-gen-openapiv2` để sinh tài liệu API.

## 📋 Sub-tasks
- [ ] Thiết kế `api/proto/farm.proto`.
- [ ] Viết Makefile hỗ trợ biên dịch proto tự động.
- [ ] Cài đặt các công cụ sinh code.
