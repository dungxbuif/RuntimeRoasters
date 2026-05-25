# Technical Plan - [RR-4-2] pkg/base — RegisterGateway & ServeSwagger

## 🎯 Chiến lược triển khai (Strategy)
- Mở rộng App struct trong pkg/base để lưu trữ và khởi chạy gateway mux.
- Tích hợp server file tĩnh cho Swagger JSON.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Thêm phương thức `RegisterGateway` vào App.
- [ ] Thêm phương thức `ServeSwagger` vào App.
- [ ] Cập nhật `app.Run` để khởi chạy HTTP server cùng với gRPC server.

## 🧪 Xác minh (Verification)
- [ ] Kiểm tra mount thành công gateway mux.
- [ ] Kiểm tra truy cập Swagger JSON qua HTTP.