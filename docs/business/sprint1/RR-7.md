# [RR-7] Control App — Jaeger (Observability)

- **Summary:** Sử dụng Jaeger All-in-one làm Control Plane để giám sát hệ thống (Distributed Tracing).
- **Priority:** `HIGH`
- **Status:** `TO_DO`
- **Depends on:** RR-5

---

## User Story

> As an admin, I want to use Jaeger as the central tracing platform, so that I can monitor service interactions and performance without complex setup.

---

## Acceptance Criteria

### Scenario 1: Jaeger UI accessible
- **Given:** `docker compose up -d` hoàn tất.
- **When:** Tôi truy cập `localhost:16686`.
- **Then:** Jaeger UI hiển thị và sẵn sàng truy vấn traces.

### Scenario 2: Trace visibility
- **Given:** Demo Service và KrakenD đã cấu hình OTLP exporter.
- **When:** Một request được gửi qua hệ thống.
- **Then:** Trace tương ứng xuất hiện trong Jaeger với đầy đủ các span từ Gateway đến Service.

---

## 🛠️ Technical Notes
- Sử dụng Jaeger All-in-one (In-memory storage cho môi trường dev).
- Cổng OTLP: `4317` (gRPC), `4318` (HTTP).
- Giao diện UI: `16686`.
