# [RR-5] Docker Complete — SigNoz + Full Compose Overhaul

- **Summary:** Hoàn thiện `docker-compose.yaml` với SigNoz observability stack và fix toàn bộ port conflicts.
- **Priority:** `HIGH`

---

## User Story

> As a developer, I want a single `docker compose up -d` to bring up the complete infrastructure including observability, so that I can develop without manually configuring any monitoring tool.

---

## Acceptance Criteria

### Scenario 1: SigNoz khởi động và nhận OTLP data
- **Given:** `docker-compose.yaml` đã được cập nhật.
- **When:** Tôi chạy `docker compose up -d`.
- **Then:** Container `rr-signoz` và `rr-clickhouse` ở trạng thái `Running`. SigNoz nhận OTLP data tại port `4317` (gRPC) và `4318` (HTTP).

### Scenario 2: SigNoz UI không expose ra host
- **Given:** `docker compose up -d` đã hoàn tất.
- **When:** Tôi thử truy cập `localhost:3301` trực tiếp trên browser.
- **Then:** Connection bị từ chối — port 3301 không được bind ra host.

### Scenario 3: Redpanda Console không conflict
- **Given:** `docker compose up -d` đã hoàn tất.
- **When:** Tôi kiểm tra danh sách port đang lắng nghe.
- **Then:** Redpanda Console chạy tại `8090`, không conflict với bất kỳ service nào khác.

### Scenario 4: Toàn bộ containers start không lỗi
- **Given:** `docker-compose.yaml` đã đúng.
- **When:** Tôi chạy `docker compose up -d && docker compose ps`.
- **Then:** Tất cả containers có status `running` hoặc `healthy`.
