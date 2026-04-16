# Lộ trình Phát triển RuntimeRoasters (Portfolio Edition)

**Mục tiêu cốt lõi:** Xây dựng một project cực mạnh về Technical (Showcase cho vị trí Senior Backend), thể hiện khả năng thiết kế hệ thống phân tán, xử lý dữ liệu lớn và kiến trúc Microservices phức tạp bằng Golang.

**Chiến lược (The PM & TA View):**
*   **Nguồn lực:** 1 Solo Developer (Full-stack/Backend focus).
*   **Độ dài Sprint:** 1 Tuần / Sprint (Giữ nhịp độ nhanh, tập trung dứt điểm từng pattern).
*   **Phương pháp:** Vertical Slice (Cắt dọc). Không làm dàn trải tất cả các API, mà chọn những luồng (flow) khó nhất để làm từ A-Z (từ DB lên tới Kafka rồi sang Service khác).

---

## 🚀 Phase 1: The Foundation & "Hello World" of Events
*Mục tiêu: Đặt nền móng vững chắc và chứng minh khả năng giao tiếp bất đồng bộ.*

### Sprint 1: Kỷ nguyên Hạ tầng & Core Framework
*   **Task 1 (DevOps):** Viết file `compose.yaml` dựng toàn bộ Local Infra: PostgreSQL, Kafka (hoặc Redpanda cho nhẹ), Valkey, Elasticsearch, Jaeger. Đảm bảo tất cả giao tiếp được với nhau.
*   **Task 2 (TA):** Xây dựng "Internal Framework" tại `pkg/`:
    *   `pkg/logger`: Tích hợp Zap, tự động nhận Trace-ID.
    *   `pkg/response`: Implement chuẩn lỗi **RFC 7807** (Problem Details).
    *   `pkg/database`: Helper kết nối Postgres (GORM hoặc sqlx).
*   **Task 3 (TA):** Setup API Gateway cơ bản (Kong hoặc KrakenD) để làm chốt chặn.

### Sprint 2: The First Vertical Slice (Outbox & CQRS)
*Chứng minh bạn biết cách giải quyết bài toán Dual-Write và tách biệt luồng Đọc/Ghi.*
*   **Task 1 (Farm Service):** Xây dựng API `POST /farms/batches` (Tạo mẻ cà phê).
    *   *Hardcore element:* Implement **Transactional Outbox**. Khi tạo mẻ, lưu vào DB và lưu 1 record vào bảng `outbox` trong cùng 1 Transaction.
*   **Task 2 (Farm Worker):** Viết 1 background goroutine đọc bảng `outbox` và publish event `BatchCreated` lên Kafka một cách an toàn (At-least-once delivery).
*   **Task 3 (Trace Service - CQRS):** Xây dựng consumer lắng nghe Kafka topic. Nhận event `BatchCreated` và thực hiện Upsert vào **Elasticsearch**.
*   **Task 4 (Trace Service):** Xây dựng API `GET /traces/{batch_id}` đọc siêu tốc từ Elasticsearch.

*=> Hết Sprint 2: Bạn đã có thể chém gió trong CV về Event-Driven, Outbox Pattern và CQRS.*

---

## ⚡ Phase 2: Distributed Transactions (Saga Pattern)
*Mục tiêu: Xử lý bài toán khó nhất của Microservices - Giao dịch phân tán.*

### Sprint 3: The Orchestrator & The Participant
*   **Task 1 (Retail Service):** Xây dựng API `POST /orders` (Tạo đơn nhập hàng cho tiệm).
    *   Trạng thái ban đầu: `PENDING`.
    *   Publish event `OrderCreated` lên Kafka (vẫn dùng Outbox Pattern).
*   **Task 2 (Warehouse Service):** Lắng nghe `OrderCreated`.
    *   Thực hiện logic: Kiểm tra tồn kho -> Reserve (giữ chỗ) số lượng.
    *   *Thành công:* Publish `InventoryReserved`.
    *   *Thất bại (Hết hàng):* Publish `InventoryFailed`.

### Sprint 4: External Webhooks & The Rollback (Compensating Action)
*Chứng minh bạn biết cách handle 3rd party an toàn và xử lý lỗi dây chuyền.*
*   **Task 1 (Webhook Service):** Xây dựng API nhận Webhook từ Stripe (giả lập thanh toán).
    *   *Hardcore element:* Verify chữ ký **HMAC**. Áp dụng **Inbox Pattern** (lưu `Message_ID` vào DB) để chống Duplicate Webhook (Idempotency). Sau đó publish `PaymentCompleted` lên Kafka.
*   **Task 2 (Retail Service - Saga Coordinator):** Lắng nghe các event từ Warehouse và Payment.
    *   Nếu nhận `InventoryFailed` hoặc timeout Payment: Chuyển order sang `CANCELLED`.
    *   *Saga Rollback:* Publish event `OrderCancelled`.
*   **Task 3 (Warehouse & Payment Service):** Lắng nghe `OrderCancelled` để thực hiện bù trừ (Release tồn kho đã giữ, gọi API Stripe để Refund).

*=> Hết Phase 2: Resume của bạn có thể ghi "Designed and implemented Saga Choreography with automatic compensating actions & dual idempotency".*

---

## 🛡️ Phase 3: "God Mode" & Security
*Mục tiêu: Đạt chuẩn Production-grade, biến project thành một hệ thống hoàn hảo.*

### Sprint 5: Traceability & Real-time
*   **Task 1 (Logistics Service):** Viết API cập nhật tọa độ GPS liên tục. Lưu trạng thái vào **Valkey** (GEO caching) thay vì DB để tối ưu performance.
*   **Task 2 (Observability):** Chăm chút lại OpenTelemetry. Đảm bảo 1 request từ lúc Retail tạo Order, chui qua Kafka đến Warehouse, trả về KQ... đều nối chung 1 `Trace-ID` và vẽ lên biểu đồ Gantt của **Jaeger** thật đẹp (Dùng để cap màn hình gắn vào README/CV).

### Sprint 6: Zero Trust & Auth
*   **Task 1 (Security):** Tích hợp Ory Kratos hoặc tự viết 1 service Auth cấp phát JWT.
    *   Cấu hình Gateway offload việc verify JWT.
*   **Task 2 (Authorization):** Tích hợp **Casbin** vào middleware của Go để phân quyền chi tiết (RBAC/ABAC).
*   **Task 3 (Internal Sec):** Setup **mTLS** (chứng chỉ số) cho các kết nối gRPC nội bộ giữa các service (nếu bạn dùng gRPC cho luồng sync).

---

