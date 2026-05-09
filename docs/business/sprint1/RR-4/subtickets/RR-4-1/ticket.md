# [RR-4-1] Proto Toolchain — buf setup & farm.proto

- **Summary:** Thiết lập bộ công cụ biên dịch Protobuf và định nghĩa contract tối giản cho Farm Service.
- **Parent:** [RR-4](../ticket.md)
- **Priority:** `HIGH`

---

## 📖 User Story

> As a developer, I want a single `make proto` command to compile `.proto` files into Go stubs and a `swagger.json`, so that the API contract is the single source of truth for both backend and frontend.

---

## 🔍 Acceptance Criteria

### Scenario 1: Biên dịch thành công
- **Given:** File `api/proto/farm/v1/farm.proto` tồn tại với định nghĩa `FarmService`.
- **When:** Tôi chạy `make proto`.
- **Then:** Các file `farm.pb.go`, `farm_grpc.pb.go`, `farm.pb.gw.go`, và `farm.swagger.json` được sinh ra trong `src/api/gen/go/farm/v1/` không có lỗi.

### Scenario 2: Lint proto
- **Given:** File `.proto` đã được viết.
- **When:** Tôi chạy `make proto-lint`.
- **Then:** Không có warning hoặc error nào được in ra.

### Scenario 3: Phát hiện breaking change
- **Given:** Một field number trong `.proto` bị thay đổi.
- **When:** Tôi chạy `make proto-breaking`.
- **Then:** Lệnh exit với lỗi mô tả rõ breaking change.

---

