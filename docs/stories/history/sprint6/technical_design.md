# Sprint 6 Technical Design: Retail & Order Saga

## 1. Overview
Sprint 6 focuses on building the **Retail Service**, which serves as the starting point for the **Order Saga (Choreography)** flow. This service manages ordering from retail stores (Stores) and tracks order status throughout the supply chain.

## 2. Domain Models

### 2.1 Retail Store
Represents the retail locations of Runtime Roasters.

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `name` | String | Store name |
| `city` | String | City (Hanoi, HCM, Da Nang) |
| `address` | String | Detailed address |
| `manager_id` | UUID | ID of the Store Manager |
| `manager_email`| String | Email of the Store Manager (used for seeding/audit) |
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

**Business Rule:** Each store has one Store Manager. Only the Store Manager has the authority to place orders (Create Order) for their respective store.

### 2.2 Order
Manages order status and customer/staff information.

| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `store_id` | UUID | Foreign Key -> Stores |
| `total_amount` | Decimal | Total order value |
| `status` | Enum | PENDING, PREPARING, SHIPPING, COMPLETED, REJECTED |
| `idempotency_key` | String | Duplication prevention (Unique Index) |

## 3. Saga Flow (Choreography)

1. **Retail:** Create Order (`PENDING`) + Save `retail.order.created` to Outbox.
2. **Warehouse:** Receive `retail.order.created` -> Check inventory -> Reserve Stock -> Publish `warehouse.stock.reserved` (or `failed`).
3. **Logistics:** Receive `warehouse.stock.reserved` -> Find vehicle -> Publish `logistics.shipment.assigned`.
4. **Retail:** Update status based on subsequent events:
    - `warehouse.stock.reserved` -> `PREPARING`
    - `logistics.shipment.assigned` -> `SHIPPING`
    - `logistics.shipment.delivered` -> `COMPLETED`
    - Any `failed` event -> `REJECTED` (Rollback/Compensate).

## 4. Technical Mechanisms

### 4.1 Transactional Outbox
Uses the same pattern as `farm-service`:
- `db.Transaction(func(tx *gorm.DB) error { ... })`
- Atomically save `orders` and `outbox_events`.

### 4.2 Idempotency
- **API Level:** `Idempotency-Key` header stored in Redis/Valkey (24h TTL).
- **DB Level:** Unique Index on `orders.idempotency_key`.

### 4.3 Authorization
Uses `Casbin` middleware to ensure that the `STORE_MGR` role is only permitted to create orders for the `store_id` they manage.
