# [RR-1] Infrastructure Kick-off

**User Story:**
Dưới vai trò là **DevOps Engineer**, tôi muốn thiết lập một bộ cấu hình hạ tầng tập trung để toàn bộ đội ngũ phát triển có thể khởi động mọi thành phần cần thiết chỉ bằng một câu lệnh duy nhất.

**Business Context:**
Tính nhất quán của môi trường phát triển (Local Consistency) là yếu tố sống còn để giảm thiểu lỗi "Works on my machine". Hạ tầng vững chắc giúp đẩy nhanh tiến độ ở các Sprint sau.

---

## 🔍 Phạm vi Thực hiện (Scope)
Xây dựng hệ thống `docker-compose.yaml` tại thư mục `deployments/`, bao gồm:
- **Persistence:** PostgreSQL (với script khởi tạo database tự động), Valkey (Redis), Cassandra.
- **Messaging:** Kafka Cluster & Kafka UI.
- **Observability:** Jaeger, Elasticsearch.
- **API Gateway:** KrakenD.

---

## 🛠️ Quy tắc Kỹ thuật (Tech Rules)
- Sử dụng **Docker Networks** để cô lập traffic giữa các tầng.
- **Volume Persistence:** Phải mount volume ra host để không mất dữ liệu khi restart container.
- **Healthchecks:** Mọi container phải có healthcheck script để đảm bảo thứ tự khởi động.

---

## ✅ Acceptance Criteria (AC)
1. **End-to-End Startup:** Chạy `docker compose up -d` không báo lỗi, mọi service đều `healthy`.
2. **Database Auto-init:** Postgres tự động tạo `farm_db`, `auth_db`, `warehouse_db` khi khởi động lần đầu.
3. **UI Accessibility:** Truy cập được Kafka UI (:8080) và Jaeger UI (:16686) từ trình duyệt.
