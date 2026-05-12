# Runtime Roasters — Master Checklist

Tài liệu này là "Single Source of Truth" để rà soát toàn bộ hệ thống (Whole App) từ hạ tầng, bảo mật, code đến quy trình triển khai. Checklist này sẽ được cập nhật liên tục qua từng task.

---

## 🏗 PHASE 1: Hạ tầng & Cấu hình (Infrastructure & Config)
Đảm bảo nền tảng vững chắc trước khi chạy Application logic.

### 1.1. Databases & Storage
- [ ] **Postgres**: Đã chạy migration đầy đủ (`identity_db`, `demo_db`, ...).
- [ ] **Redis**: Đã cấu hình làm Session Store và Token Blacklist.
- [ ] **Kafka**: (Nếu sử dụng) Đã tạo sẵn các topic cần thiết (`token-revocation`, ...).
- [ ] **Healthchecks**: Mọi service hạ tầng phải có cấu hình `healthcheck` trong Docker/K8s.

### 1.2. Identity & Auth (Core Security)
- [ ] `EXPECTED_ISSUER`: Khớp với `iss` claim (vd: `http://localhost:4444/`).
- [ ] `JWKS_URL`: Trỏ đúng endpoint nội bộ (`http://identity:4434/...`).
- [ ] `INTERNAL_SECRET`: Đã đổi từ giá trị mặc định và đồng bộ giữa Nginx Proxy & Backend.
- [ ] **OIDC Client Secret**: Đã cấu hình chuỗi ngẫu nhiên bảo mật cao.

### 1.3. Environment Variables
- [ ] **Centralized Config**: Mọi biến môi trường nhạy cảm (`*_SECRET`, `*_PASSWORD`) phải được quản lý qua Secret Manager hoặc `.env.local`.
- [ ] **Frontend Env**: `NEXT_PUBLIC_*` đã được kiểm tra tính đúng đắn khi build production.

---

## 🔒 PHASE 2: Bảo mật & Resilience (Security & Resilience)
Đảm bảo hệ thống an toàn và có khả năng tự phục hồi.

### 2.1. Code Audit
- [ ] **Fail-Fast & Retry**: Code fetch dữ liệu khởi động (như JWKS) đã có Exponential Backoff Retry.
- [ ] **Error Handling**: Đã sử dụng chuẩn RFC 9457 (Problem Details), không leak thông tin lỗi nội bộ.
- [ ] **Input Validation**: Mọi API endpoint đã có validation cho request body/params.

### 2.2. Network Security
- [ ] **CORS**: Chỉ whitelist các domain chính thống.
- [ ] **API Gateway Mapping**: Mọi endpoint resource phải sử dụng danh từ số nhiều (Plural: `/v1/users`, `/v1/farms`).
- [ ] **Gateway Config Sync**: Đã chạy `force-recreate` hoặc `reload` gateway để đảm bảo bản đồ routing mới nhất được nạp.
- [ ] **TLS/SSL**: Đảm bảo HTTPS được cấu hình cho mọi traffic public.

---

## 🧪 PHASE 3: Kiểm chứng luồng & Tính năng (Feature Verification)
Xác nhận nghiệp vụ chạy đúng thực tế.

### 3.1. Auth Flow (OIDC)
- [ ] **Login/Logout**: Hoạt động trơn tru, xóa sạch session/cookie khi logout.
- [ ] **Seamless Consent**: Người dùng không bị hỏi lại quyền nếu đã được tin tưởng.
- [ ] **Identity Context**: Identity của user được truyền chính xác vào UseCase layer.

### 3.2. Observability (Giám sát)
- [ ] **Tracing**: Mỗi request đều sinh ra Trace ID, User ID xuất hiện trong Span Attributes.
- [ ] **Logging**: Log theo format JSON (Structured Logging) để dễ dàng query.
- [ ] **Metrics**: Các metrics cơ bản (Request count, Latency, Error rate) đã được export.

---

## 🚀 PHASE 4: Triển khai & Vận hành (Deployment & Ops)
Các bước cuối cùng trước khi "Go Live".

### 4.1. Orchestration
- [ ] **Dependency Order**: Microservices có `depends_on` kèm `condition: service_healthy`.
- [ ] **Init Containers**: (K8s) Có container chờ DB/Identity sẵn sàng trước khi app chạy.
- [ ] **Resource Limits**: Đã cấu hình CPU/Memory Requests & Limits cho từng container.

### 4.2. CI/CD
- [ ] **Unit Tests**: Đã pass 100% trước khi merge.
- [ ] **Linter**: Không còn lỗi linting nghiêm trọng.
- [ ] **Image Security**: Docker images đã được quét lỗ hổng (Scan vulnerabilities).

---
*Cập nhật lần cuối: 2026-05-07 bởi Antigravity*
- [ ] **Unit Tests**: Đã pass 100% trước khi merge.
- [ ] **Linter**: Không còn lỗi linting nghiêm trọng.
- [ ] **Image Security**: Docker images đã được quét lỗ hổng (Scan vulnerabilities).

---
*Cập nhật lần cuối: 2026-05-07 bởi Antigravity*
