# Technical Plan - [RR-4-1] Proto Toolchain — buf setup & farm.proto

## 🎯 Chiến lược triển khai (Strategy)
- Sử dụng buf làm công cụ quản lý proto chính.
- Định nghĩa contract gRPC và REST gateway trong cùng một file proto.

## 🛠️ Các bước thực hiện (Implementation Steps)
- [ ] Cấu hình `buf.yaml` và `buf.gen.yaml`.
- [ ] Định nghĩa `farm.proto`.
- [ ] Viết Makefile để chạy các lệnh buf.

## 🧪 Xác minh (Verification)
- [ ] Chạy `make proto` và kiểm tra các file sinh ra.
- [ ] Chạy `make proto-lint`.
- [ ] Thử nghiệm breaking change và chạy `make proto-breaking`.