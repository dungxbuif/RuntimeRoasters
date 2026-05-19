# Sprint 6 Technical Design: Retail & Order Saga

## 1. Overview
Sprint 6 tập trung vào xây dựng **Retail Service**, đóng vai trò là điểm khởi đầu cho luồng **Order Saga (Choreography)**. Service này quản lý việc đặt hàng từ các cửa hàng bán lẻ (Stores) và theo dõi trạng thái đơn hàng xuyên suốt chuỗi cung ứng.

## 2. Domain Models

### 2.1 Retail Store (Cửa hàng)
Đại diện cho các điểm bán lẻ của Runtime Roasters.

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `name` | String | Tên cửa hàng |
| `city` | String | Thành phố (Hanoi, HCM, Da Nang) |
| `address` | String | Địa chỉ chi tiết |
| `manager_id` | UUID | ID của trưởng cửa hàng (Store Manager) |
| `manager_email`| String | Email của trưởng cửa hàng (dùng cho seeding/audit) |
| `status` | Enum | ACTIVE, INACTIVE |

**Seed Data Plan:**
- **Hanoi:**
    - **Hoan Kiem Store** (2 Ly Thai To). Manager: `mgr.hn.hoankiem@runtimeroasters.com`
    - **Cau Giay Store** (102 Tran Thai Tong). Manager: `mgr.hn.caugiay@runtimeroasters.com` (TBD in Identities)
- **Ho Chi Minh City:**
    - **District 1 Store** (45 Le Thanh Ton). Manager: `mgr.hcm.q1@runtimeroasters.com`
    - **District 7 Store** (Phu My Hung). Manager: `mgr.hcm.q7@runtimeroasters.com` (TBD in Identities)
- **Da Nang:**
    - **Hai Chau Store** (15 Bach Dang). Manager: `mgr.dn.haichau@runtimeroasters.com`

**Business Rule:** Mỗi cửa hàng có 1 Trưởng cửa hàng (Store Manager). Chỉ Trưởng cửa hàng mới có quyền thực hiện đặt hàng (Create Order) cho cửa hàng của mình.

### 2.2 Order (Đơn hàng)
Quản lý trạng thái đơn hàng và thông tin khách hàng/nhân viên.

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `store_id` | UUID | Foreign Key -> Stores |
| `total_amount` | Decimal | Tổng giá trị đơn hàng |
| `status` | Enum | PENDING, PREPARING, SHIPPING, COMPLETED, REJECTED |
| `idempotency_key` | String | Chống trùng lặp (Unique Index) |

## 3. Saga Flow (Choreography)

1. **Retail:** Tạo Order (`PENDING`) + Lưu `retail.order.created` vào Outbox.
2. **Warehouse:** Nhận `retail.order.created` -> Kiểm tra kho -> Reserve Stock -> Bắn `warehouse.stock.reserved` (hoặc `failed`).
3. **Logistics:** Nhận `warehouse.stock.reserved` -> Tìm xe -> Bắn `logistics.shipment.assigned`.
4. **Retail:** Cập nhật status dựa trên các event tiếp theo:
    - `warehouse.stock.reserved` -> `PREPARING`
    - `logistics.shipment.assigned` -> `SHIPPING`
    - `logistics.shipment.delivered` -> `COMPLETED`
    - Bất kỳ `failed` event nào -> `REJECTED` (Rollback/Compensate).

## 4. Technical Mechanisms

### 4.1 Transactional Outbox
Sử dụng chung pattern với `farm-service`:
- `db.Transaction(func(tx *gorm.DB) error { ... })`
- Lưu `orders` và `outbox_events` đồng thời.

### 4.2 Idempotency
- **API Level:** `Idempotency-Key` header lưu vào Redis/Valkey (TTL 24h).
- **DB Level:** Unique Index trên `orders.idempotency_key`.

### 4.3 Authorization
Sử dụng `Casbin` middleware để check role `STORE_MGR` chỉ được phép tạo đơn cho `store_id` mà họ quản lý.
