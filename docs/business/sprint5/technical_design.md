# Technical Design: Sprint 5 — Roastery Process Service

**Tác giả:** Tech Lead
**Trạng thái:** Approved
**Mục tiêu:** Khởi tạo dịch vụ chế biến và hiện thực hóa quy trình chuyển đổi từ hạt cà phê thô (Raw Beans) sang thành phẩm (Roasted Beans) với State Machine.

---

## 1. Kiến trúc Dịch vụ (Service Architecture)

`process-service` sẽ là một Microservice độc lập, tuân thủ nghiêm ngặt Clean Architecture và sử dụng **Manual Dependency Injection**.

### Sơ đồ luồng (Flow Diagram):
1. **Kafka Consumer** nhận `HarvestedEvent`.
2. **Inbox Pattern**: Kiểm tra bảng `inbox_events` để đảm bảo Idempotency.
3. **UseCase** tạo bản ghi `ProcessingBatch` mới với trạng thái `RECEIVED`.
4. **User** cập nhật trạng thái qua từng công đoạn: Hulling -> Drying -> Roasting.
5. **Final Step (Roasting)**: Sinh mã `Batch ID` (RR-TYPE-YYYYMMDD-RAND) và phát hành sự kiện `BatchRoastedEvent`.

---

## 2. Domain State Machine (The Workflow Engine)

Chúng ta thực hiện logic kiểm tra trạng thái ngay trong Domain Layer.

### Trạng thái & Chuyển đổi Hợp lệ:
- `RECEIVED` -> `HULLING`
- `HULLING` -> `DRYING`
- `DRYING` -> `ROASTING`
- `ROASTING` -> `COMPLETED`

**Quy tắc:** Không được nhảy cóc. Logic kiểm tra được bọc trong hàm `CanTransitionTo(next Status) bool`.

---

## 3. Transactional Inbox (Guaranteed Processing)

Để xử lý sự kiện từ Farm Service một cách tin cậy, Process Service phải sử dụng Inbox Pattern.

### Luồng xử lý Consumer:
1. Nhận message từ Kafka.
2. Mở Transaction:
    - Kiểm tra bảng `inbox_events` xem `message_id` đã tồn tại chưa.
    - Nếu có: Bỏ qua (Duplicate).
    - Nếu chưa: Lưu `inbox_events` và Tạo `ProcessingBatch`.
3. Commit Transaction.

---

## 4. Thiết kế Cơ sở dữ liệu (Database Schema)

### 4.1. Bảng `processing_batches`
| Cột | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Định danh nội bộ. |
| `harvest_id` | UUID (FK) | **Traceability Link** tới Farm gốc. |
| `batch_id` | VARCHAR(50) | Mã thành phẩm sau khi Rang. |
| `status` | VARCHAR(20) | RECEIVED, HULLING, DRYING, ROASTING, COMPLETED. |
| `logs` | JSONB | Nhật ký chi tiết từng bước (Audit Trail). |

### 4.2. Bảng `inbox_events`
| Cột | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `id` | UUID (PK) | |
| `message_id` | VARCHAR(255) | Unique ID từ Kafka Header. |
| `processed_at` | TIMESTAMPTZ | Thời điểm xử lý thành công. |

---

## 5. Failure Scenarios

- **Kafka gửi tin trùng**: Inbox Pattern chặn đứng dựa trên `message_id`.
- **Hệ thống sập khi đang Rang**: Mẻ hạt bị treo status `ROASTING`. Cần Job quét các batch quá hạn để cảnh báo.

---
**Tech Lead Signature**
