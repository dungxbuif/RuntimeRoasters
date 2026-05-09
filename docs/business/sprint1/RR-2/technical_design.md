# Technical Plan - [RR-2] Core Framework - Config & Logger

## 🎯 Chiến lược triển khai (Strategy)
- Sử dụng Viper để quản lý cấu hình từ nhiều nguồn (env, file).
- Sử dụng Zap làm logger hiệu năng cao, hỗ trợ log có cấu trúc.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Implement `pkg/config/config.go` sử dụng `spf13/viper`.
- [ ] Implement `pkg/logger/logger.go` sử dụng `uber-go/zap`.
    - Hỗ trợ Log Level qua biến môi trường `LOG_LEVEL`.
    - Hỗ trợ định dạng Console cho Dev và JSON cho Prod.
- [ ] Tạo file `.env.example` làm mẫu.

## 🧪 Xác minh (Verification)
- [ ] Unit test cho module config.
- [ ] Kiểm tra định dạng log trên terminal trong môi trường Dev/Prod.