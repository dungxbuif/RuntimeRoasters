# [RR-7] Control App — SigNoz (Observability)

- **Goal:** Sử dụng SigNoz + ClickHouse làm Control Plane để giám sát hệ thống.
- **Business Value:** Cung cấp công cụ mạnh mẽ để quan sát luồng request giữa các microservices, hỗ trợ tối ưu hóa hiệu năng và tìm lỗi.
- **Priority:** `HIGH`

## 🔍 Acceptance Criteria

### Scenario 1: SigNoz UI accessible
- **Given:** `docker compose up -d` hoàn tất.
- **When:** Tôi truy cập `localhost:3301`.
- **Then:** SigNoz UI hiển thị và sẵn sàng truy vấn traces.

### Scenario 2: Trace visibility
- **Given:** Demo Service và KrakenD đã cấu hình OTLP exporter.
- **When:** Một request được gửi qua hệ thống.
- **Then:** Trace tương ứng xuất hiện trong SigNoz với đầy đủ các span từ Gateway đến Service.
