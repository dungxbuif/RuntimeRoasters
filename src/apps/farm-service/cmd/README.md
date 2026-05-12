# Thư mục: Cmd (Điểm khởi đầu - Entry Point)

## 🏢 Vai trò trong Clean Architecture
Đây là nơi chứa hàm `main()`. Nhiệm vụ duy nhất của nó là gọi hàm khởi tạo ứng dụng và kích hoạt lệnh chạy Server.

## 🎯 Scope (Phạm vi)
- **Nên viết:**
    - Gọi hàm `InitializeApp()` từ package `internal/app`.
    - Xử lý các tín hiệu hệ thống (OS Signals) để tắt ứng dụng an toàn (Graceful Shutdown).
    - In ra các thông tin khởi chạy ban đầu.
- **Không nên viết:**
    - Tuyệt đối không viết logic cấu hình hay nghiệp vụ ở đây.

## 🚀 Quy tắc vàng
File `main.go` nên cực kỳ ngắn gọn (thường dưới 50 dòng code).
