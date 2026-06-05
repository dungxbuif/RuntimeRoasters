# Technical Plan - [RR-7] Control App — SigNoz (Observability)

## 🎯 Chiến lược triển khai (Strategy)
- Sử dụng SigNoz + ClickHouse cho môi trường development/demo.
- Cấu hình OTLP gRPC endpoint cho các service.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Cấu hình SigNoz, ClickHouse, và OTel Collector trong `docker-compose.yaml`.
- [ ] Thiết lập cổng OTLP: `4317` (gRPC), `4318` (HTTP).

## 🧪 Xác minh (Verification)
- [ ] Truy cập SigNoz UI tại `localhost:3301`.
- [ ] Kiểm tra việc nhận trace từ Demo Service.
