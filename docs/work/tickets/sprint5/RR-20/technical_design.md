# Technical Design: [RR-20] Warehouse Service Architecture & Batch Logic

**Role:** Tech Lead
**Context:** Thiết kế hệ thống quản lý Kho tập trung vào luồng Batching (1 Batch - N Runs).

---

## 1. Data Model Strategy (1-N Relationship)

Hệ thống tuân thủ mô hình **1 Production Batch = N Roast Runs** để đảm bảo tính ổn định (Consistency) và khả năng scale trong sản xuất.

### 1.1. Bảng `production_batches` (Thực thể quản lý)
- `id`: UUID.
- `harvest_id`: Liên kết tới Farm Service (Nguồn gốc hạt thô).
- `status`: `RECEIVED` (Mới nhập thô), `PROCESSING` (Đang rang), `STOCKED` (Đã hoàn thành lô).
- `production_batch_id`: Mã thương mại (VD: `PROD-AR-20260514-001`).
- `coffee_type`: Loại hạt (Arabica/Robusta).
- `origin_code`: Vùng nguyên liệu.
- `total_input_weight`: Tổng khối lượng thô đã dùng (kg).
- `total_output_weight`: Tổng khối lượng chín thu được (kg).

### 1.2. Bảng `roast_runs` (Thực thể vật lý - "Mẻ")
- `id`: UUID.
- `batch_id`: FK trỏ tới `production_batches`.
- `input_weight`: Khối lượng thô của mẻ rang này.
- `output_weight`: Khối lượng chín của mẻ rang này.
- `created_at`: Thời điểm rang.

---

## 2. Kafka HA & Scalability Strategy

Để đảm bảo hệ thống scale được nhiều Pods mà không lỗi dữ liệu:

### 2.1. Partitioning & Key
- **Topic:** `farm.harvest.events` (Partitions: 3).
- **Partition Key:** Sử dụng `harvest_id`. Đảm bảo các event liên quan đến cùng 1 vụ thu hoạch luôn được xử lý bởi cùng 1 Pod, giữ đúng thứ tự.

### 2.2. Idempotency (Inbox Pattern)
- Mọi event nhận từ Kafka được lưu vào bảng `inbox_events` cùng với Business Logic trong 1 Transaction.
- Nếu nhận trùng `message_id`, hệ thống tự động bỏ qua (Ignore).

### 2.3. Distributed Locking (Valkey)
- Khi thực hiện `FinalizeBatch` và cập nhật `inventory_items`, service phải lấy lock theo SKU: `lock:inventory:{sku}`.
- Tránh tranh chấp dữ liệu khi nhiều Pods cùng nhập kho các Batch khác nhau cho cùng một loại sản phẩm.

---

## 3. Inventory Integration

Khi Batch chuyển sang `STOCKED`:
1. Tính `Total_Output = Sum(Roast_Runs.output_weight)`.
2. Update `inventory_items` (UPSERT).
3. Emit event `warehouse.stock.updated`.
