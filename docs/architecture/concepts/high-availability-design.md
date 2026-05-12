# High Availability (HA) & Kafka Messaging Strategy

Tài liệu này định nghĩa cấu trúc hạ tầng HA và chiến lược phân nhóm sự kiện trên Kafka để đảm bảo hệ thống Runtime Roasters không có điểm yếu chí tử (Single Point of Failure - SPoF).

---

## 1. Chiến lược Kafka: Grouping & Scaling

Để đảm bảo hiệu năng và khả năng mở rộng, chúng ta không dùng một topic duy nhất.

### 1.1. Phân nhóm Topic (Topic-per-Domain/Epic)
Chúng ta gom nhóm các sự kiện theo Domain nghiệp vụ để tối ưu hóa việc quản lý và phân quyền:
- `auth.user.events`: Chứa `USER_CREATED`, `ROLE_ASSIGNED`.
- `retail.order.events`: Chứa `ORDER_CREATED`, `ORDER_CANCELLED`.
- `warehouse.stock.events`: Chứa `STOCK_RESERVED`, `STOCK_FAILED`.
- `logistics.tracking.events`: Chứa `GPS_UPDATED`, `SHIPMENT_MOVED`.

### 1.2. Kỹ thuật Consumer Group (Horizontal Scaling)
- Mỗi microservice khi lắng nghe sự kiện sẽ sử dụng một **Consumer Group ID** riêng (VD: `warehouse-service-group`).
- Khi chúng ta scale Warehouse Service lên 5 instances, Kafka sẽ tự động chia các Partitions cho 5 instances này. Đảm bảo mỗi message chỉ được xử lý đúng 1 lần bởi 1 instance (Load Balancing).

### 1.3. Partitioning & Ordering
- Sử dụng **Business Key** (VD: `order_id` hoặc `batch_id`) làm **Kafka Key**.
- Kafka đảm bảo tất cả message có cùng Key sẽ luôn vào cùng một Partition.
- Điều này cực kỳ quan trọng để đảm bảo **Thứ tự sự kiện** (VD: `ORDER_CREATED` phải được xử lý trước `ORDER_CANCELLED`).

---

## 2. Thiết kế High Availability (HA) Toàn diện

Hệ thống được thiết kế để chạy đa node ở mọi tầng.

### 2.1. Tầng Dữ liệu (Persistence HA)
- **PostgreSQL**: Triển khai mô hình **Master-Slave Replication**.
    - Node Master: Xử lý Write.
    - Node Slave (Replica): Xử lý Read.
    - Sử dụng **PgBouncer** để tự động chuyển hướng kết nối nếu Master gặp sự cố.
- **Valkey/Valkey**: Triển khai **Valkey Sentinel** hoặc **Cluster Mode** để tự động Failover.

### 2.2. Tầng Truyền thông (Kafka HA)
- **Broker Clustering**: Tối thiểu 3 Kafka Brokers chạy song song.
- **Replication Factor = 3**: Mỗi message được sao chép sang 3 node khác nhau. Nếu 1-2 node chết, dữ liệu vẫn an toàn.
- **Zookeeper Ensemble**: Chạy 3-5 nodes Zookeeper để bầu chọn Leader (Quorum).

### 2.3. Tầng Cổng kết nối (Gateway HA)
- **KrakenD**: Chạy stateless trong Docker Swarm hoặc Kubernetes với ít nhất 2 replicas.
- Sử dụng một **External Load Balancer** (như Nginx hoặc Cloud LB) phía trước để check health và chia tải cho các node KrakenD.

### 2.4. Tầng Dịch vụ (Microservices HA)
- Mọi service (Farm, Auth, Retail...) đều được thiết kế **Stateless**.
- Trạng thái phiên làm việc (Session) lưu tập trung trong Valkey/Kratos.
- Cho phép scale ngang (Horizontal Scaling) vô hạn mà không mất dữ liệu.

---

## 3. Tổng kết mô hình HA

| Thành phần | Cơ chế HA | Mục tiêu |
| :--- | :--- | :--- |
| **Identity** | Ory Kratos (2+ instances) | Luôn có thể Login. |
| **API Gateway** | KrakenD (2+ instances) | Không nghẽn cổ chai. |
| **Message Broker** | Kafka (3 Brokers, RF=3) | Không mất sự kiện. |
| **Database** | Postgres (Primary + Standby) | Data an toàn tuyệt đối. |
| **Connection** | PgBouncer (Stateless) | Pooling ổn định. |

---
*TechLead Signed-off: 2026-05-10*
