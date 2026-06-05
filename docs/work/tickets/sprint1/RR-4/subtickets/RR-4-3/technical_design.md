# Technical Plan - [RR-4-3] Dependency Injection — Google Wire scaffold

## 🎯 Chiến lược triển khai (Strategy)
- Thiết lập cấu trúc Wire cho việc injection dependencies.
- Tạo ProviderSet chung cho các thành phần infra.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Cài đặt wire tool.
- [ ] Tạo file `wire.go` với các ProviderSet.
- [ ] Chạy `wire gen` để sinh code.

## 🧪 Xác minh (Verification)
- [ ] Kiểm tra `main.go` được rút gọn tối đa.
- [ ] Build project thành công sau khi inject.