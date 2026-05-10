# ADR 0002: Lựa chọn SigNoz thay cho Stack truyền thống (Prometheus/Jaeger/Loki)

## Trạng thái
**Accepted**

## Bối cảnh (Context)
Để đảm bảo khả năng quan sát (Observability) cho hệ thống phân tán, ban đầu đội ngũ định sử dụng stack truyền thống: Prometheus cho Metrics, Loki cho Logs, và Jaeger cho Traces.
Tuy nhiên, cấu hình và vận hành 3 hệ thống riêng lẻ này rất phức tạp. Việc tương quan (correlation) giữa Trace ID và Log cũng đòi hỏi cấu hình thủ công tốn kém thời gian.

## Quyết định (Decision)
Sử dụng **SigNoz** làm nền tảng tập trung (Centralized Hub) cho cả 3 trụ cột (Logs, Metrics, Traces).
- Cài đặt SigNoz qua Docker Compose.
- Sử dụng trực tiếp chuẩn **OpenTelemetry (OTLP)** từ các Go Microservices bắn thẳng về SigNoz mà không cần dựng thêm OTel Collector riêng lẻ tại mỗi node.

## Hậu quả (Consequences)
- **Tích cực:** Quản lý hạ tầng đơn giản hơn (chỉ 1 UI duy nhất). Tự động liên kết giữa Log và Trace (chỉ cần search TraceID là ra toàn bộ Log). Chuẩn OTLP tránh vendor lock-in.
- **Tiêu cực:** SigNoz sử dụng ClickHouse làm storage backend, đòi hỏi thêm resource (RAM/CPU) trên môi trường development.

## Nguồn tham khảo
- **Sprint:** Sprint 1.
- **Ticket/Thảo luận:** Tích hợp `pkg/telemetry` và Logging framework.
- Xem chi tiết tại: [Observability Strategy](../concepts/observability-strategy.md).
