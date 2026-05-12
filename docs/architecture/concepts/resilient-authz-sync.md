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

## 3. Chiến lược "3 Trụ Cột" Chống Lỗi (Phased Implementation)

Để đảm bảo tiến độ phát triển, cơ chế đồng bộ được triển khai theo 2 giai đoạn:

### Giai đoạn 1: Resilience via Polling & Snapshots (Hiện tại)
Hiện tại, hệ thống tập trung vào tính ổn định và khả năng tự phục hồi mà không phụ thuộc vào Kafka:
- **Bootstrapping (gRPC Snapshot):** Khi một service khởi động, nó tải toàn bộ Snapshot quyền từ Auth Service.
- **Self-Healing (Polling):** Cứ mỗi 10-15 phút, service tự động đồng bộ lại toàn bộ dữ liệu.
- **Ưu điểm:** Giảm độ phức tạp vận hành, dễ debug trong giai đoạn MVP.

### Giai đoạn 2: Real-time via Kafka (Tech Debt - Sprint 12+)
Nâng cấp khả năng cập nhật thời gian thực:
- **Live Update (Kafka):** Lắng nghe thay đổi qua Kafka watcher để cập nhật RAM ngay lập tức.
- **Đồng bộ hóa:** Kết hợp cả Kafka (Real-time) và Polling (Chữa lành) để đạt được độ tin cậy tối đa.
- **Mục tiêu:** Giảm độ trễ cập nhật quyền từ vài phút xuống còn vài mili giây.

## 4. Advanced: Query-Level Authorization (ABAC to SQL)

Đối với các bài toán phân quyền dữ liệu phức tạp (vd: Warehouse, Retail), hệ thống hỗ trợ cơ chế chuyển đổi Casbin Policies thành câu lệnh SQL `WHERE`.

- **Cơ chế:** Sử dụng `e.GetAllowedObjectConditions` để trích xuất logic từ ABAC policies.
- **Pattern ứng dụng:** **GORM Scopes**. Tích hợp bộ lọc Casbin trực tiếp vào câu lệnh truy vấn của Repository để đảm bảo hiệu năng và tính bảo mật ở mức bản ghi.
- **Chiến lược áp dụng:**
    - **Simple Ownership (KISS):** Dùng Native GORM `WHERE` (vd: Farm Service).
    - **Dynamic Filtering:** Dùng Casbin Query Level (vd: Warehouse Service - Regional Isolation).

## 5. Identity to Authorization Sync (Kratos -> Casbin)

Để giải quyết vấn đề bất đồng bộ giữa kho Định danh (Kratos) và kho Phân quyền (Casbin), đặc biệt là với dữ liệu nạp sẵn (Seeded Data), hệ thống áp dụng cơ chế đồng bộ hóa chủ động:

- **Source of Truth:** Ory Kratos là nguồn sự thật duy nhất cho thông tin người dùng và các thuộc tính (traits) của họ, bao gồm cả `role`.
- **Bootstrapping Sync:** Khi `auth-service` khởi động, nó thực hiện một tiến trình quét toàn bộ Identity từ Kratos. Nếu phát hiện Identity nào chưa có bản ghi `g` (grouping policy) tương ứng trong Casbin, `auth-service` sẽ tự động tạo bản ghi đó dựa trên trait `role` của Identity.
- **Runtime Integrity:** Khi tạo người dùng qua API của `auth-service`, service đóng vai trò Orchestrator đảm bảo ghi dữ liệu vào cả Kratos và Casbin một cách nguyên tử (atomic-like).
- **Resilience:** Cơ chế này giúp hệ thống tự phục hồi quyền hạn ngay cả khi Database phân quyền bị xóa sạch hoặc khi có các tài khoản được nạp trực tiếp vào Kratos qua các kênh Admin/CLI.

> **⚠️ Lưu ý kỹ thuật (Technical Debt):** Giải pháp Bootstrapping Sync hiện tại là giải pháp tạm thời (Phase 1) để đảm bảo hệ thống hoạt động ổn định trong giai đoạn phát triển và test dữ liệu Seed. Trong tương lai (Phase 2), chúng ta cần nghiên cứu triển khai **Kratos Webhooks** hoặc cơ chế **Event-driven Sync** để đồng bộ hóa thời gian thực ngay khi có thay đổi Identity, tránh việc phải quét toàn bộ Database mỗi khi service khởi động.
