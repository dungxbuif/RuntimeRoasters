# Technical Design: Sprint 9 — Traceability & Search (Reference Aligned)

Mục tiêu: Xây dựng hệ thống truy xuất nguồn gốc (Traceability) tối ưu hiệu năng sử dụng mô hình **CQRS** và **Elasticsearch** từ dự án UrbanX.

---

## 1. Kiến trúc CQRS (Command Query Responsibility Segregation)
*Tham chiếu: UrbanX Catalog Service CQRS*

Chúng ta tách biệt hoàn toàn luồng dữ liệu nghiệp vụ (Write) và luồng dữ liệu tra cứu (Read):

### 1.1. Write Side (Command Side)
- **Source**: Các sự kiện từ Farm Service (`batch.created`), Warehouse (`stock.reserved`), Retail (`order.confirmed`).
- **Storage**: Các service này lưu trữ dữ liệu nghiệp vụ vào PostgreSQL của riêng chúng (Isolated DBs).

### 1.2. Read Side (Query Side - Trace Service)
- **Aggregator**: Trace Service lắng nghe tất cả các topic Kafka liên quan.
- **Elasticsearch Index**: Lưu trữ bản ghi đã được **Denormalize** (Phi chuẩn hóa).
    - Thay vì phải join nhiều bảng từ nhiều service, toàn bộ hành trình của hạt cà phê được gói gọn trong 1 document duy nhất trong Elasticsearch.
- **Tốc độ**: Đảm bảo phản hồi < 100ms cho các truy vấn tra cứu hành trình mẻ hàng (QR Code Scan).

---

## 2. Thiết kế Elasticsearch (Search Model)
*Tham chiếu: UrbanX Search Service logic*

### 2.1. Mapping Index: `coffee_traceability`
Chúng ta sử dụng các kiểu dữ liệu nâng cao của Elasticsearch để tối ưu tìm kiếm:
- `batch_id`: `keyword` (Primary lookup key).
- `farm_location`: `geo_point` (Để hiển thị bản đồ vùng trồng trên UI).
- `history_timeline`: `nested` object. Chứa mảng các sự kiện:
  ```json
  {
    "status": "HARVESTED",
    "timestamp": "2026-05-10T...",
    "location": "Sơn La",
    "description": "Thu hoạch Arabica thượng hạng"
  }
  ```

---

## 3. Đồng bộ dữ liệu (Event-Driven Sync)
*Tham chiếu: UrbanX Outbox & Consumer Flow*

### 3.1. Luồng đồng bộ
1. **Event Trigger**: Khi có thay đổi tại Farm/Warehouse... sự kiện được bắn lên Kafka qua Outbox Pattern.
2. **Trace Consumer**: Trace Service (`KafkaTraceEventConsumer`) nhận tin.
3. **Upsert Logic**: 
    - Nếu `batch_id` chưa tồn tại -> Tạo mới document.
    - Nếu `batch_id` đã có -> Cập nhật (Push) thêm 1 mốc sự kiện mới vào mảng `history_timeline`.
4. **Resiliency**: Nếu Elasticsearch bị treo, Trace Service sẽ tự động đọc lại các sự kiện từ Kafka offset gần nhất để cập nhật bù khi hệ thống online trở lại.

---

## 4. Search API Excellence
*Tham chiếu: UrbanX API Gateway Forwarding*

- **Endpoint**: `GET /v1/trace/{batch_id}`.
- **Gateway**: KrakenD định tuyến trực tiếp tới Trace Service.
- **Response**: Trả về Full Journey (Timeline) đã được tối ưu cho giao diện Mobile (QR Landing Page).

---
*TechLead Signed-off: 2026-05-10*
