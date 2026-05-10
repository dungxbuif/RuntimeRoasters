# System-Wide Architectural Design: Runtime Roasters

Tài liệu này định nghĩa các tiêu chuẩn kỹ thuật áp dụng xuyên suốt cho toàn bộ hệ thống microservices, đảm bảo tính đồng nhất, khả năng mở rộng và vận hành tin cậy.

---

## 1. Tiêu chuẩn Truyền thông & Sự kiện (Messaging Standard)

Để tránh tình trạng "Event Spaghetti", chúng ta áp dụng các tiêu chuẩn sau:

- **Protocol**: Sử dụng Kafka làm xương sống.
- **Event Schema**: Áp dụng chuẩn **CloudEvents** (hoặc tương đương) để mọi event đều có Metadata thống nhất:
    - `id`: Unique event ID.
    - `source`: Tên service phát hành.
    - `type`: Loại sự kiện (VD: `retail.order.created`).
    - `time`: RFC3339 timestamp.
    - `data`: Payload nghiệp vụ.
- **Reliability**: Mọi hành động phát hành sự kiện PHẢI qua **Transactional Outbox**.
- **Idempotency**: Mọi Consumer PHẢI triển khai **Inbox Pattern** để chống xử lý lặp.

## 2. Tiêu chuẩn Trạng thái & Nhất quán (Distributed Consistency)

- **Saga Pattern**: Ưu tiên **Choreography** (phối hợp phi tập trung) cho các luồng đơn giản (Order-Warehouse). Cân nhắc **Orchestration** (Temporal.io) nếu luồng nghiệp vụ vượt quá 5 bước.
- **Error Handling (API)**: Tuyệt đối tuân thủ **RFC 9457 (Problem Details)**. Mọi service phải dùng chung `pkg/errs` để trả về lỗi có Trace-ID.

## 3. Tiêu chuẩn Quan sát & Giám sát (Observability Standard)

- **Tracing**: Sử dụng chuẩn **W3C Trace Context**. Trace-ID phải được truyền từ Gateway -> gRPC -> Kafka -> Downstream services.
- **Metrics**: Áp dụng mô hình **RED** (Request Rate, Error Rate, Duration) cho tất cả các endpoint.
- **Logging**: Sử dụng Structured Logging (JSON). Mọi Log record phải đính kèm `trace_id` và `span_id` nếu có context.

## 4. Tiêu chuẩn Hạ tầng & DI (Standardized Bootstrap)

- **Internal Library (`pkg/`)**: 
    - Mọi service mới đều phải sử dụng bộ khung khởi tạo từ `pkg/base`.
    - Việc "nối dây" (Dependency Injection) thực hiện thủ công trong `internal/app/init.go` (ưu tiên sự minh bạch, tránh magic của các framework DI quá nặng).
- **Database Access**: Sử dụng **GORM** làm chuẩn ORM, nhưng cấm dùng `AutoMigrate` trên môi trường dùng chung. Phải có file SQL Migration riêng.

---

## 5. Danh sách các "Global Infrastructure Components"

1.  **PgBouncer**: Tập trung pool kết nối Postgres.
2.  **Ory Kratos/Hydra**: Quản lý danh tính và cấp phát Token.
3.  **Auth Service (Casbin Manager)**: Quản trị chính sách tập trung.
4.  **OTel Collector**: Trạm trung chuyển dữ liệu quan sát.
5.  **Redis/Valkey**: Shared cache và Idempotency store.

---
*TechLead Signed-off: 2026-05-10*
