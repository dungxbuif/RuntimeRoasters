# [RR-4-5] Farm Service — GetDemoFarm handler

- **Summary:** Implement handler tối giản để verify toàn bộ luồng kỹ thuật hoạt động end-to-end.
- **Parent:** [RR-4](../ticket.md)
- **Priority:** `HIGH`
- **Depends on:** RR-4-1, RR-4-2, RR-4-3, RR-4-4

---

## 📖 User Story

> As a developer, I want a working `GetDemoFarm` endpoint that returns hardcoded data, so that I can verify the full path from Gateway → grpc-gateway → gRPC handler is correctly wired before building real business logic.

---

## 🔍 Acceptance Criteria

### Scenario 1: gRPC handler trả về dữ liệu
- **Given:** Farm Service đang chạy.
- **When:** Tôi gọi gRPC method `FarmService.GetDemoFarm`.
- **Then:** Handler trả về một `Farm` object với dữ liệu hardcoded hợp lệ.

### Scenario 2: Gọi qua HTTP trả về JSON
- **Given:** Farm Service đang chạy với grpc-gateway mounted.
- **When:** Tôi gọi `GET localhost:8080/v1/farms/demo`.
- **Then:** Response là JSON hợp lệ map đúng với message `Farm` trong proto.

### Scenario 3: Full path qua Gateway
- **Given:** KrakenD và Farm Service đều đang chạy.
- **When:** Tôi gọi `GET localhost:8081/v1/farms/demo`.
- **Then:** Response JSON hợp lệ được trả về — xác nhận toàn bộ luồng RR-4-1 đến RR-4-4 hoạt động đúng.

---

