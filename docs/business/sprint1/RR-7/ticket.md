# [RR-7] Control App — Jaeger (Observability)

- **Goal:** Sử dụng Jaeger All-in-one làm Control Plane để giám sát hệ thống.
- **Business Value:** Cung cấp công cụ mạnh mẽ để quan sát luồng request giữa các microservices, hỗ trợ tối ưu hóa hiệu năng và tìm lỗi.
- **Priority:** `HIGH`

## 🔍 Acceptance Criteria

### Scenario 1: Jaeger UI accessible
- **Given:** `docker compose up -d` hoàn tất.
- **When:** Tôi truy cập `localhost:16686`.
- **Then:** Jaeger UI hiển thị và sẵn sàng truy vấn traces.

### Scenario 2: Trace visibility
- **Given:** Demo Service và KrakenD đã cấu hình OTLP exporter.
- **When:** Một request được gửi qua hệ thống.
- **Then:** Trace tương ứng xuất hiện trong Jaeger với đầy đủ các span từ Gateway đến Service.