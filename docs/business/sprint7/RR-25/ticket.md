# [RR-25] Logistics Service & Driver Tracking

- **Summary:** Triển khai hạ tầng kỹ thuật cho dịch vụ Logistics, tích hợp Valkey GEO và Kafka.
- **Priority:** `CRITICAL`
- **Role:** Tech Lead

---

## 🔍 Acceptance Criteria

### Scenario 1: Driver Location Update
- **Given:** Tài xế gửi tọa độ (Lat, Long).
- **When:** API nhận request.
- **Then:** Lưu tọa độ vào Valkey sử dụng `GEOADD` và publish sự kiện `logistics.gps.updated` lên Kafka.

### Scenario 2: Finding Nearest Driver
- **Given:** Một chuyến hàng cần vận chuyển.
- **When:** Gọi hàm tìm kiếm tài xế.
- **Then:** Sử dụng `GEORADIUS` hoặc `GEOSEARCH` trong Valkey để trả về danh sách tài xế gần nhất.

---

## 📋 Sub-tickets
- [ ] **RR-25.1**: Scaffolding Logistics Service (Clean Architecture).
- [ ] **RR-25.2**: Triển khai Valkey Repository cho GPS data.
- [ ] **RR-25.3**: Implement Shipment State Machine (PENDING -> DELIVERED).
- [ ] **RR-25.4**: Viết Script Simulator giả lập 5-10 tài xế di chuyển thực tế trên bản đồ.
- [ ] **RR-25.5**: [Compatibility] Implement Kafka Consumer xử lý chính xác payload `warehouse.stock.updated` từ Sprint 5.
