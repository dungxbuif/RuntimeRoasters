# Technical Design: Sprint 10 — Reliability & Observability (Reference Aligned)

Mục tiêu: Đạt được trạng thái "Hệ thống trong suốt" (Transparent System) thông qua Distributed Tracing từ dự án UrbanX và lưu trữ nhật ký bất biến trên Cassandra.

---

## 1. Toàn diện hóa OpenTelemetry (Full Instrumentation)
*Tham chiếu: UrbanX ServiceDefaults & Distributed Tracing*

Chúng ta áp dụng triết lý "Mọi request đều có dấu vân tay":

### 1.1. W3C Trace Context Propagation
- **Gateway (KrakenD)**: Khởi tạo Root Span cho mỗi request từ Browser.
- **Inter-service (gRPC)**: Sử dụng gRPC Unary/Stream Interceptors để chuyển tiếp Trace-ID qua metadata.
- **Async (Kafka)**: Tích hợp OTel Header Injector/Extractor vào Kafka Producer/Consumer. 
- **Kết quả**: Một Trace duy nhất sẽ nối liền từ Retail -> Warehouse -> Payment -> Kafka -> Trace Service.

### 1.2. Trạm trung chuyển OTel Collector
- Cấu hình một service trung tâm nhận dữ liệu OTLP từ tất cả microservices.
- Export Traces tới **Jaeger/SigNoz** và Metrics tới **Prometheus**.

---

## 2. Nhật ký bất biến (Audit Service & Cassandra)
*Tham chiếu: Triết lý "Không bao giờ mất dấu" của UrbanX*

Trong khi UrbanX lưu Audit vào Postgres/Elastic, Runtime Roasters nâng cấp lên **Cassandra** để đảm bảo khả năng ghi cực nhanh và dữ liệu không thể xóa sửa.

### 2.1. Cassandra Schema Design
- **Table**: `audit_logs`.
- **Partition Key**: `batch_id` (để truy vấn lịch sử mẻ hàng nhanh).
- **Clustering Key**: `occurred_at` (để sắp xếp theo thời gian).
- **Trường đặc biệt**: `previous_hash` & `current_hash`. Chúng ta áp dụng **Hash Chaining** để phát hiện ngay lập tức nếu có ai đó cố tình can thiệp vào DB.

---

## 3. Giám sát hệ thống (Monitoring - RED Pattern)
*Tham chiếu: UrbanX Aspire Dashboard Metrics*

Xây dựng bộ chỉ số sức khỏe cho từng microservice:
- **Rate**: Số lượng yêu cầu/giây.
- **Errors**: Tỷ lệ lỗi (HTTP 5xx, gRPC Codes != OK).
- **Duration**: Độ trễ (P50, P95, P99 Latency).
- **Kafka Lag**: Theo dõi độ trễ của các Consumer (Đặc biệt quan trọng cho Saga và CQRS Sync).

---

## 4. Tích hợp Health Checks
*Tham chiếu: UrbanX /health endpoints*

Mọi service phải cung cấp endpoint `/health`:
- **Liveness**: Báo hiệu service còn sống (Process is running).
- **Readiness**: Báo hiệu service đã sẵn sàng nhận traffic (DB connected, Kafka connected).
- **Dashboard**: Tích hợp trạng thái Xanh/Đỏ vào trang quản trị Admin.

---
*TechLead Signed-off: 2026-05-10*
