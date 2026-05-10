# Plan: RR-13 End-to-End Secure Integration (Verification)

## 1. Overview
Mục tiêu là kiểm chứng mô hình bảo mật **Two-Gate Security** đã triển khai:
- **Gate 1 (Gateway):** KrakenD kiểm tra `scope` của JWT.
- **Gate 2 (Service):** Microservice kiểm tra `role` của người dùng qua Casbin.

## 2. Kịch bản kiểm chứng (Test Scenarios)

### S1: Unauthorized (No Token)
- **Hành động:** Gọi `GET http://localhost:8081/v1/demo/ping` không kèm Authorization header.
- **Kết quả mong đợi:** KrakenD trả về `401 Unauthorized`.

### S2: Forbidden - Gate 1 (Invalid Scope)
- **Hành động:** Gửi JWT hợp lệ nhưng không có scope `demo:read`.
- **Kết quả mong đợi:** KrakenD trả về `403 Forbidden`.

### S3: Forbidden - Gate 2 (Invalid Role)
- **Hành động:** Gửi JWT có scope `demo:read` (qua được Gate 1) nhưng user có role `guest` (bị Casbin tại service chặn).
- **Kết quả mong đợi:** Service trả về `403 Forbidden`.

### S4: Success (Full Access)
- **Hành động:** Gửi JWT có scope `demo:read` và role `admin`.
- **Kết quả mong đợi:** Trả về `200 OK`.

## 3. Implementation Plan
1. **KrakenD Config:** Bổ sung `auth/validator` vào `krakend.json` cho endpoint `/v1/demo/ping`. (Đã làm)
2. **E2E Test Suite:** Viết bộ test Go trong `src/apps/demo-service/__tests__/e2e/security_test.go`.
3. **Trace Verification:** Kiểm tra SigNoz để xác nhận `user.id` được gắn vào span.

## 4. Sub-tasks
- [x] Configure KrakenD with JWT Validator (Gate 1).
- [ ] Implement Go E2E Test Suite skeleton.
- [ ] Implement Token Generator helper for testing.
- [ ] Run verification against running stack.
