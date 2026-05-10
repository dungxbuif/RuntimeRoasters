# Resilient AuthZ Sync Architecture

Để đảm bảo hệ thống Microservices hoạt động ổn định, phân quyền tập trung nhưng thực thi phi tập trung, dự án áp dụng mô hình **Resilient Sync Architecture** thông qua Casbin.

## 1. Mục tiêu (Goals)
- **Tập trung hóa:** Admin chỉ cần cấu hình quyền tại một nơi duy nhất.
- **Độ trễ bằng 0 (Zero Latency):** Microservices thực thi kiểm tra quyền in-memory, không phải gọi qua network cho mỗi request.
- **Khả năng chịu lỗi (Resilience):** Microservice vẫn có thể tự khởi động và phân quyền ngay cả khi Auth Service hay Kafka sập.

## 2. Kiến trúc Tổng thể

Dựa trên mô hình **Centralized Management, Distributed Enforcement**:

1. **Centralized Auth Service (The Writer):** Một dịch vụ trung tâm quản lý toàn bộ các bản ghi `casbin_rule` trong Postgres. Đây là nơi duy nhất thực hiện các thao tác Ghi (Add/Remove policy).
2. **Kafka & Zookeeper (The Bus):** Sử dụng Kafka (phiên bản Zookeeper-based) làm kênh truyền tin thời gian thực. Khi có sự thay đổi quyền, Auth Service bắn Event lên Kafka.
3. **In-memory Enforcer (The Reader):** Các Microservices sử dụng **Casbin v3.x** để thực thi kiểm tra quyền trực tiếp trên RAM, đảm bảo tốc độ tối đa.

## 3. Chiến lược "3 Trụ Cột" Chống Lỗi

Hệ thống dựa trên 3 trụ cột để đảm bảo tính nhất quán (Consistency) và khả năng phục hồi (Resilience):

- **Bootstrapping (gRPC Snapshot):** Khi một service khởi động, hành động đầu tiên của nó là gọi gRPC tới Auth Service để tải về toàn bộ "Snapshot" quyền hiện tại. Việc này giúp service sẵn sàng nhận traffic. Nếu gọi gRPC thất bại, service sẽ áp dụng **Exponential Backoff** để thử lại thay vì crash.
- **Live Update (Kafka):** Trong suốt quá trình hoạt động, service lắng nghe các thay đổi thời gian thực qua Kafka watcher để cập nhật RAM ngay lập tức, đảm bảo độ trễ cập nhật quyền gần như tức thời.
- **Self-Healing (Polling/Re-sync):** Để phòng trường hợp service lỡ mất một Event từ Kafka (do network split, Kafka downtime), cứ mỗi 10-15 phút, service tự động chạy background job đồng bộ lại toàn bộ dữ liệu từ Auth Service qua gRPC. Cơ chế này đóng vai trò "chữa lành" những sai lệch dữ liệu.
