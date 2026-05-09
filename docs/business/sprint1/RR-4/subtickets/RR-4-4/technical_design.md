# Technical Plan - [RR-4-4] Infrastructure — KrakenD, Apps Shell & Core DBs

## 🎯 Chiến lược triển khai (Strategy)
- Sử dụng Docker Compose để khởi chạy toàn bộ stack hạ tầng.
- Cấu hình KrakenD routing cho các service.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Cập nhật `docker-compose.yaml` với KrakenD, Postgres, Redis và App Shells.
- [ ] Cấu hình `krakend.json`.
- [ ] Khởi tạo project Next.js cho Client và Control Apps.

## 🧪 Xác minh (Verification)
- [ ] Kiểm tra trạng thái container.
- [ ] Test routing từ KrakenD đến service.