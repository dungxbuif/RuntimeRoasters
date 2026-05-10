# ADR 0002: Lựa chọn SigNoz thay cho Stack truyền thống (Prometheus/Jaeger/Loki)

**Trạng thái:** 🔴 REVERTED (Đã gỡ bỏ - 2026-05-10)

## Bối cảnh (Context)
Dự án ban đầu lựa chọn SigNoz làm nền tảng Observability tập trung để thay thế cho việc dựng lẻ tẻ các container Jaeger, Prometheus, v.v.

## Quyết định (Decision)
Dừng sử dụng SigNoz và ClickHouse trong giai đoạn hiện tại để giảm bớt tài nguyên máy local cho Developer. Hệ thống sẽ quay lại sử dụng các giải pháp nhẹ hơn (ví dụ: Jaeger chạy lẻ) hoặc tập trung hoàn thiện logic nghiệp vụ trước.

## Hậu quả (Consequences)
- Phải cập nhật lại `docker-compose.dev.yaml` và cấu hình telemetry trong code.
- Giảm tải cho RAM/CPU máy dev (ClickHouse chiếm khá nhiều tài nguyên).
