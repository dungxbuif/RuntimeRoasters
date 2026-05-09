# Technical Plan - [RR-4-5] Farm Service — GetDemoFarm handler

## 🎯 Chiến lược triển khai (Strategy)
- Implement handler gRPC đơn giản với dữ liệu giả.
- Đảm bảo handler được đăng ký với cả gRPC server và gateway mux.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Tạo handler `GetDemoFarm` trong Farm Service.
- [ ] Đăng ký handler trong `InitializeApp`.

## 🧪 Xác minh (Verification)
- [ ] Gọi gRPC method trực tiếp.
- [ ] Gọi qua HTTP gateway.
- [ ] Gọi qua API Gateway (KrakenD).