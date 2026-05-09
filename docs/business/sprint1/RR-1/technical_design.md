# Technical Plan - [RR-1] Infrastructure Kick-off

## 🎯 Chiến lược triển khai (Strategy)
- Sử dụng Docker Compose để orchestrate các dịch vụ hạ tầng.
- Tách biệt cấu hình và dữ liệu thông qua volumes.
- Đảm bảo tính sẵn sàng cao cho các thành phần quan trọng.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Hoàn thiện `deployments/docker-compose.yaml`.
    - Sử dụng Image `postgres:16-alpine`.
    - Sử dụng `Redpanda` bản mới nhất để tương thích gRPC/Kafka API.
    - Cấu hình Jaeger với OTLP protocol hỗ trợ gRPC.
- [ ] Hoàn thiện `deployments/init-db.sql`.
- [ ] Thiết lập mạng nội bộ cho các container.

## 🧪 Xác minh (Verification)
- [ ] Chạy lệnh `docker-compose up -d` và kiểm tra trạng thái container.
- [ ] Thực thi và kiểm tra logs của từng service.
- [ ] Kiểm tra kết nối đến Postgres và các database đã tạo.
- [ ] Truy cập Redpanda Console (`localhost:8080`) và Jaeger UI (`localhost:16686`).
