# Technical Design: Sprint 6 — Warehouse Service

**Tác giả:** Tech Lead
**Trạng thái:** Approved
**Mục tiêu:** Quản lý tồn kho vật lý và cung cấp cơ chế giữ chỗ hàng (Reservation) an toàn, chuẩn bị cho luồng Saga điều phối đơn hàng.

---

## 1. Kiến trúc Dịch vụ (Service Architecture)

`warehouse-service` quản lý số lượng tồn kho thực tế của các `Batch ID` thành phẩm. Service này đóng vai trò là **Saga Participant**.

### Luồng nghiệp vụ chính (Saga Steps):
1. **Reserve Stock**: Giữ chỗ hàng khi có đơn hàng mới. Giảm `available`, tăng `reserved`.
2. **Confirm Stock**: Trừ kho thật khi thanh toán thành công. Giảm `reserved`, giảm `total`.
3. **Release Stock**: Giải phóng hàng giữ chỗ nếu đơn hàng bị hủy. Tăng `available`, giảm `reserved`.

---

## 2. Concurrency Control (Distributed Locking)

Để tránh tình trạng "Bán quá số lượng" (Overselling) trong môi trường phân tán, chúng ta sử dụng **Valkey Distributed Lock** (Redlock).

### Thuật toán Reserve:
1. `Lock(batch_id)` trong Valkey với TTL (ví dụ: 5s).
2. Kiểm tra `available_qty >= requested_qty`.
3. Nếu OK: Thực hiện cập nhật DB và lưu vào bảng `reservations` status `PENDING`.
4. `Unlock(batch_id)`.

---

## 3. Defensive Programming: Chốt chặn số âm

Mọi câu lệnh Update kho phải có điều kiện bảo vệ ở tầng SQL để tránh lỗi logic ứng dụng:
```sql
UPDATE inventories 
SET available_qty = available_qty - ?, 
    reserved_qty = reserved_qty + ?
WHERE batch_id = ? 
AND available_qty >= ?; -- Chốt chặn cuối cùng
```

---

## 4. Thiết kế Cơ sở dữ liệu (Database Schema)

### 4.1. Bảng `inventories`
| Cột | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `batch_id` | VARCHAR(50) (PK) | |
| `available_qty` | DECIMAL | Số lượng có sẵn để bán. |
| `reserved_qty` | DECIMAL | Số lượng đang chờ thanh toán. |
| `total_qty` | DECIMAL | Tổng thực tế (`available + reserved`). |

### 4.2. Bảng `reservations`
| Cột | Kiểu dữ liệu | Mô tả |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Thường là `order_id`. |
| `batch_id` | VARCHAR(50) | |
| `status` | VARCHAR(20) | PENDING, CONFIRMED, CANCELLED. |

---

## 5. Failure Scenarios & Self-Healing

- **Orchestrator "quên" confirm/release**: Job quét bảng `reservations` quá 15 phút chưa confirm -> Tự động Release (TTL logic).
- **Valkey Lock sập**: Fallback về DB Row Locking (`SELECT FOR UPDATE`).

---
**Tech Lead Signature**
