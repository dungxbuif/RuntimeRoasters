# [RR-4-4] Infrastructure — KrakenD, Apps Shell & Core DBs

- **Summary:** Thiết lập hạ tầng cơ bản gồm API Gateway, App Shells (Client + Control) và các database cốt lõi (Postgres, Valkey).
- **Parent:** [RR-4](../ticket.md)
- **Priority:** `HIGH`

---

## 🔍 Acceptance Criteria

### Scenario 1: Infrastructure khởi động
- **Given:** `docker-compose.yaml` đã được cập nhật.
- **When:** Tôi chạy `docker compose up -d`.
- **Then:** Các container `rr-krakend`, `rr-postgres`, `rr-valkey`, `client-app`, `control-app` ở trạng thái `Running`.

### Scenario 2: Apps shells accessible
- **Then:** 
  - `client-app` (Business UI) chạy tại port `3000`.
  - `control-app` (Admin) chạy tại port `3001`.

### Scenario 3: KrakenD route đến Farm Service
- **Given:** Farm Service đang chạy local.
- **When:** Tôi gọi `GET localhost:8081/v1/farms/demo`.
- **Then:** KrakenD forward request đến Farm Service và trả về response đúng.

---

