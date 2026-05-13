# 🌊 Detailed Data Flow & Integration Steps

Tài liệu này mô tả chi tiết các bước thực thi kỹ thuật cho từng luồng nghiệp vụ quan trọng.

---

## 1. Luồng Thu hoạch & Truy xuất (Farm to Trace Flow)

Đây là luồng minh chứng cho **CQRS** và **Transactional Outbox**.

### Bước 1: Client gửi yêu cầu thu hoạch
*   **Action:** UI gửi `POST /v1/harvests` tới API Gateway.
*   **Auth:** Gateway kiểm tra JWT, thêm `X-User-Role: Farmer` vào header.

### Bước 2: Xử lý tại Farm Service (Write Side)
1.  **Delivery layer:** Nhận JSON, validate cấu hình (coffee_type, quantity).
2.  **UseCase layer:** Gọi `CreateHarvest`.
3.  **Repository layer:** Mở một Database Transaction (Kế hoạch RR-4.2):
    *   `INSERT INTO harvests (...)`
    *   `INSERT INTO outbox_events (event_type, payload) VALUES ('farm.harvest.created.v1', '...')`
4.  **Commit Transaction:** Dữ liệu được lưu vĩnh viễn vào Postgres.

### Bước 3: Đẩy sự kiện lên Kafka (Asynchronous)
1.  **Relay Worker** (Goroutine ngầm): Quét bảng `outbox_events` nơi status là `PENDING`.
2.  **Publisher:** Đẩy message theo chuẩn CloudEvents lên Kafka topic `farm.harvest.events`.
3.  **Action:** Cập nhật record outbox thành `COMPLETED`.

### Bước 4: Cập nhật Read-Model tại Trace Service (Read Side)
1.  **Consumer:** Trace Service lắng nghe topic `farm.harvest.events`.
2.  **Idempotency Check:** Kiểm tra `Message_ID` đã xử lý chưa (Inbox Pattern).
3.  **Indexing:** Thực hiện **Upsert** vào Elasticsearch index `coffee_traces`.
4.  **Status:** Dữ liệu sẵn sàng để tìm kiếm full-text.

---

## 2. Luồng Thanh toán & Saga (Payment & Saga Flow)

### Bước 1: Retail Service khởi tạo Saga
*   Tạo Order với trạng thái `PENDING`.
*   Ghi Outbox sự kiện `OrderCreated`.

### Bước 2: Warehouse Service giữ hàng
*   Nghe `OrderCreated`.
*   Dùng **Valkey Distributed Lock** để khóa mã hàng.
*   Trừ tồn kho tạm thời (Reserve).
*   Ghi Outbox `InventoryReserved`.

### Bước 3: Webhook Service nhận tiền
*   Stripe gửi Webhook `payment_intent.succeeded`.
*   Webhook Service verify **HMAC signature**.
*   Ghi Inbox (tránh duplicate) -> Publish `PaymentCompleted`.

### Bước 4: Saga Kết thúc
*   Retail Service nghe đủ 2 sự kiện thành công -> Chuyển Order thành `SUCCESS`.
*   Nếu 1 trong 2 thất bại -> Kích hoạt **Compensating Action** (Hoàn tiền / Nhả kho).
