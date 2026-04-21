# [RR-4-4] Infrastructure — KrakenD & Swagger UI

- **Summary:** Thêm API Gateway (KrakenD) và Swagger UI vào môi trường local Docker Compose.
- **Parent:** [RR-4](../ticket.md)
- **Priority:** `HIGH`

---

## 📖 User Story

> As a developer, I want an API Gateway and Swagger UI available via `docker compose up`, so that the local environment mirrors production topology and I can test APIs through the gateway from day one.

---

## 🔍 Acceptance Criteria

### Scenario 1: KrakenD khởi động
- **Given:** `docker-compose.yaml` đã được cập nhật.
- **When:** Tôi chạy `docker compose up -d`.
- **Then:** Container `rr-krakend` ở trạng thái `Running` tại port `8081`.

### Scenario 2: KrakenD route đến Farm Service
- **Given:** Farm Service đang chạy local.
- **When:** Tôi gọi `GET localhost:8081/v1/farms/demo`.
- **Then:** KrakenD forward request đến Farm Service và trả về response đúng — không trả về lỗi gateway.

### Scenario 3: Swagger UI khởi động và load contract
- **Given:** Container `rr-swagger-ui` đang chạy tại port `8082`.
- **When:** Tôi mở `localhost:8082` trên browser.
- **Then:** Swagger UI hiển thị đúng nội dung từ `farm.swagger.json` của Farm Service.

---

## 🛠️ Technical Notes
- API Gateway: `devopsfaith/krakend:2`.
- Swagger UI: `swaggerapi/swagger-ui`.
- KrakenD config: `deployments/krakend/krakend.json`.
- Trong môi trường dev, KrakenD trỏ đến `host.docker.internal` để gọi service chạy local.
