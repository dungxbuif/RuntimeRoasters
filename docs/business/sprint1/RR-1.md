# [RR-1] Infrastructure Kick-off

- **Summary:** Thiết lập hạ tầng môi trường chạy đa dịch vụ (Multi-services).
- **Priority:** `CRITICAL`
- **Description:** Xây dựng file cấu hình Docker Compose bao gồm tất cả các thành phần data persistence, messaging và observability cho toàn hệ thống.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Khởi tạo hạ tầng cơ bản
- **Given:** Tôi đã cấu hình file `docker-compose.yaml` với Postgres, Redpanda, Elasticsearch, Valkey, Cassandra, và Jaeger.
- **When:** Tôi thực hiện lệnh `docker-compose up -d` tại thư mục `deployments/`.
- **Then:** Toàn bộ 7 containers phải ở trạng thái "Running" và không có lỗi khởi động.

### Scenario 2: Kiểm tra tính sẵn sàng của Database
- **Given:** Container Postgres đã khởi động thành công.
- **When:** Tôi kết nối vào Postgres instance.
- **Then:** Tôi phải thấy các cơ sở dữ liệu `farm_db`, `warehouse_db`, `retail_db`,... đã được tạo sẵn từ script khởi tạo.

### Scenario 3: Kiểm tra giao diện quản trị hạ tầng
- **Given:** Toàn bộ hệ thống hạ tầng đang chạy.
- **When:** Tôi truy cập `localhost:8080` (Redpanda) và `localhost:16686` (Jaeger).
- **Then:** Giao diện quản trị phải hiển thị và sẵn sàng để giám sát dữ liệu.

---

## 🛠️ Technical Notes
- Sử dụng Image `postgres:16-alpine`.
- Sử dụng `Redpanda` bản mới nhất để tương thích gRPC/Kafka API.
- Cấu hình Jaeger với OTLP protocol hỗ trợ gRPC.

## 📋 Sub-tasks
- [ ] Hoàn thiện `deployments/docker-compose.yaml`.
- [ ] Hoàn thiện `deployments/init-db.sql`.
- [ ] Thực thi và kiểm tra logs.
