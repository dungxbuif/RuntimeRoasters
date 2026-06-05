# Technical Design: Sprint 7 - Real-time Logistics

## 1. Overview
Dịch vụ Logistics chịu trách nhiệm điều phối việc vận chuyển cà phê thành phẩm từ Kho (Warehouse) đến các Cửa hàng (Retail). Điểm nhấn kỹ thuật là việc xử lý dữ liệu vị trí thời gian thực (Real-time GPS) và cơ chế tìm kiếm tài xế thông minh.

## 2. Architecture Details

### 2.1 Technology Stack
- **Database**: PostgreSQL (Lưu thông tin Shipment, Driver, Vehicle).
- **Cache/Geo**: Valkey (Redis) - Chuyên biệt cho dữ liệu không gian (Geospatial).
- **Messaging**: Kafka (Sự kiện trạng thái và tọa độ).

### 2.2 Domain Models (Postgres)
- **Shipment**: `ID`, `OrderID`, `OriginWarehouseID`, `DestinationStoreID`, `DriverID`, `Status` (`PENDING`, `ASSIGNED`, `PICKED_UP`, `IN_TRANSIT`, `DELIVERED`).
- **Driver**: `ID`, `Name`, `Phone`, `Status` (`IDLE`, `BUSY`, `OFFLINE`).

### 2.3 Real-time GPS Flow (Valkey)
Tọa độ tài xế sẽ **không** lưu vào Postgres để tránh quá tải.
1. **Lưu trữ**: Sử dụng `GEOADD drivers:locations <long> <lat> <driver_id>`.
2. **Truy vấn**: Sử dụng `GEOSEARCH drivers:locations FROMLONLAT <long> <lat> BYRADIUS 5 km WITHDIST`.
3. **Expiration**: Sử dụng `SETEX driver:last_seen:<id>` để biết tài xế còn online hay không (vì GEO set không hỗ trợ TTL từng phần tử).

## 3. Workflow & Events

### 3.1 Luồng tạo chuyến hàng (Incoming)
- **Consume Event**: `warehouse.stock.updated` (Payload từ Sprint 5).
- **Integration Logic**: 
    - Logistics Service phải xử lý chính xác cấu trúc payload: `{batch_id, sku, quantity, timestamp}`.
    - Chuyến hàng (`Shipment`) sẽ lưu `BatchID` này làm tham chiếu ngược (Back-reference) để đảm bảo tính truy xuất nguồn gốc (Traceability).
- **Tương thích ngược**: 
    - Đảm bảo định dạng `BatchID` (ví dụ: `BATCH-SL-0518-123`) được kế thừa chính xác.
    - Trong giai đoạn demo, nếu event thiếu `DestinationStoreID`, Logistics sẽ mặc định gán một Cửa hàng demo (ví dụ: `STORE-001`) để hoàn tất luồng.
- **Trigger**: Tìm tài xế gần `OriginWarehouseID` ngay khi nhận được event.

### 3.2 Luồng di chuyển (Tracking)
- Driver App (Simulator) -> `POST /v1/drivers/location`.
- Backend:
    - `GEOADD` vào Valkey.
    - Publish Kafka: `logistics.gps.updated` (Payload: `{driver_id, lat, long, shipment_id}`).
- **Dashboard**: Lắng nghe Kafka qua WebSocket để vẽ tài xế di chuyển trên bản đồ mà không cần load lại trang.

### 3.3 Luồng hoàn tất (Saga Participant)
- Khi `Shipment` chuyển sang `DELIVERED`.
- Publish Kafka: `logistics.shipment.delivered`.
- **Retail Service** sẽ nhận event này để cộng kho cửa hàng và hoàn tất Order.

## 4. Driver Simulator Design
Một script Go độc lập sẽ chạy để demo:
1. Init 5 tài xế tại các tọa độ xung quanh kho hàng.
2. Mỗi 2 giây, tính toán tọa độ mới (tiến về phía đích) và gọi API update location.
3. Giả lập các tình huống: Tài xế từ chối chuyến, tài xế đi lạc đường.

## 5. Security & Idempotency
- **Idempotency**: Sử dụng `msgID` từ Kafka cho các sự kiện chuyển trạng thái Shipment.
- **Auth**: Chỉ tài xế được gán cho chuyến hàng mới được phép cập nhật tọa độ cho `shipment_id` đó.
