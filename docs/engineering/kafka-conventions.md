# Kafka Engineering Conventions (HA & Scalability)

Mục tiêu: Đảm bảo hệ thống đạt trạng thái High Availability (HA) và có khả năng Scale-out an toàn trên môi trường nhiều Pods.

---

## 1. Kiến trúc 3 Lớp Bảo vệ (Triple-Shield)

### Lớp 1: Partition Keys (Tính nhất quán luồng)
- **Quy tắc:** Mọi event thuộc cùng một thực thể (Aggregate) **BẮT BUỘC** phải có cùng một Partition Key.
- **Thực thi:** 
    - Producer (Farm): `msg.Key = []byte(harvest_id)`.
    - Kết quả: Tất cả các sự kiện của một mẻ hàng sẽ luôn rơi vào cùng một Partition và được xử lý tuần tự bởi cùng một Pod.

### Lớp 2: Inbox Pattern (Chống trùng lặp - Idempotency)
- **Quy tắc:** Mỗi message nhận được phải được kiểm tra tính duy nhất trước khi xử lý logic nghiệp vụ.
- **Thực thi:**
    - Sử dụng bảng `inbox_events` (message_id UNIQUE).
    - Lưu `message_id` và thực hiện Business Logic trong cùng một Database Transaction.

### Lớp 3: Distributed Locking (Bảo vệ Resource chung)
- **Quy tắc:** Khi cập nhật các Resource dùng chung (ví dụ: Tổng tồn kho SKU), phải sử dụng khóa phân tán.
- **Thực thi:**
    - Tool: **Valkey Redlock** (`github.com/go-redsync/redsync`).
    - Key format: `lock:inventory:{sku}`.

---

## 2. Cấu hình Infrastructure (Production-ready)
- **Replication Factor:** 3.
- **Min In-sync Replicas:** 2.
- **Acks:** `all` (Đảm bảo an toàn dữ liệu tuyệt đối cho chuỗi cung ứng).
- **Consumer Group:** Mỗi service sử dụng một `group.id` riêng biệt (ví dụ: `warehouse-service-group`).
