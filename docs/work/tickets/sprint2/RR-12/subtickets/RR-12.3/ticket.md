# [RR-12.3] Auth Service: Event-Driven Publishing

- **Mục tiêu:** Thông báo tức thời cho các dịch vụ khác (Readers) khi có thay đổi về chính sách phân quyền thông qua Kafka.
- **Mô tả:** Triển khai cơ chế phát tín hiệu (Change Signals) tại Auth Service. Đây là "Trụ cột 2" trong chiến lược Resilience, đảm bảo độ trễ cập nhật quyền trên toàn hệ thống ở mức thấp nhất (Near real-time).

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Kafka Producer được tích hợp vào `auth-service` sử dụng `github.com/segmentio/kafka-go`.
2. Mỗi khi có thao tác Write thành công (Add/Remove policy), một message được gửi vào topic `auth.policy.changed`.
3. Message nội dung đơn giản (ví dụ: `{"action": "RELOAD"}`) để kích hoạt các Reader làm mới bộ nhớ đệm.
4. Đảm bảo tính tin cậy: Tín hiệu thay đổi chỉ được gửi sau khi Database đã commit thành công.

## 🛠 Task list cho Developer
- [ ] Cấu hình Kafka Producer trong `auth-service` (broker list, topic name).
- [ ] Tạo Wrapper cho Casbin Enforcer tại Auth Service để intercept các lệnh Write.
- [ ] Triển khai hàm `PublishChangeSignal` gửi message vào Kafka.
- [ ] Xử lý lỗi: Log lỗi nếu Kafka không gửi được message nhưng không làm rollback transaction của DB (Eventual Consistency sẽ xử lý qua Polling).
