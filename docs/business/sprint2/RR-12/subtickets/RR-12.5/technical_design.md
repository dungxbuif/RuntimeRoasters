# RR-12.5: System Integration & Resilience Testing

## 1. Mục tiêu (Goal)
Đảm bảo tính toàn vẹn và khả năng chịu lỗi của hệ thống phân quyền trong môi trường chạy thực tế (Docker Compose / K8s).

## 2. Kế hoạch tích hợp

### 2.1 Deployment Configuration
- Cập nhật `docker-compose.dev.yaml`:
    - Thêm `auth-service` với các biến môi trường cấu hình Postgres và Kafka.
    - Cấu hình `demo-service` trỏ endpoint `AUTH_SERVICE_GRPC` về `auth-service:50051`.

### 2.2 Wire-up tại Demo Service
- Trong `main.go` của `demo-service`:
    1. Khởi tạo `casbin.ResilientEngine`.
    2. Đăng ký `casbin.UnaryServerInterceptor(engine)` vào gRPC Server.

## 3. Resilience Test Cases

| Case | Hành động | Kết quả mong muốn |
|---|---|---|
| **Happy Path** | Sửa quyền trong DB -> Publish Kafka. | Reader cập nhật trong < 1s. |
| **Kafka Failure** | Tắt Kafka -> Sửa quyền trong DB. | Reader cập nhật sau tối đa 5 phút (Polling). |
| **Auth Service Down** | Tắt Auth Service -> Khởi động Reader. | Reader retry liên tục và khởi động thành công ngay khi Auth Service bật lại. |
| **Network Partition** | Cắt kết nối gRPC giữa Reader và Writer. | Reader vẫn hoạt động với bản snapshot cũ trong memory. |

## 4. Xác minh (Verification)
- Thực hiện chuỗi test case trên trong môi trường dev.
- Kiểm tra log của cả Writer và Reader để đảm bảo các tín hiệu đồng bộ được gửi/nhận đúng.
