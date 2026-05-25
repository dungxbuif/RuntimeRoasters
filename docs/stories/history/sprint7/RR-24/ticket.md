# [RR-24] Quy trình điều phối vận chuyển (Shipment Orchestration)

- **Summary:** Định nghĩa quy trình nghiệp vụ từ khi hàng sẵn sàng tại kho đến khi giao tới cửa hàng bán lẻ.
- **Priority:** `HIGH`
- **Role:** Logistics Manager

---

## 🔍 Acceptance Criteria

### Scenario 1: Tự động tạo chuyến hàng
- **Given:** Warehouse service phát sự kiện `warehouse.stock.updated` cho một Production Batch.
- **When:** Hệ thống nhận được sự kiện.
- **Then:** Tạo một bản ghi `Shipment` ở trạng thái `PENDING`, gán thông tin từ kho gửi đến cửa hàng nhận.

### Scenario 2: Điều phối tài xế (Assignment)
- **Given:** Có chuyến hàng `PENDING`.
- **When:** Hệ thống tìm thấy tài xế rảnh trong bán kính 5km (qua Valkey GEO).
- **Then:** Chuyển trạng thái sang `ASSIGNED` và thông báo cho tài xế.

### Scenario 3: Hoàn tất giao hàng
- **Given:** Chuyến hàng đang `IN_TRANSIT`.
- **When:** Tài xế xác nhận đã giao hàng tại điểm đích.
- **Then:** Trạng thái chuyển thành `DELIVERED`, bắn sự kiện thông báo cho Retail Service.

---

## 📋 Sub-tickets
- [ ] **RR-24.1**: Thiết kế UI màn hình Dashboard theo dõi chuyến hàng (Map view).
- [ ] **RR-24.2**: Định nghĩa danh sách các trạng thái chuyến hàng (Shipment Statuses).
- [ ] **RR-24.3**: Thiết kế luồng xử lý khi không tìm thấy tài xế (Retry/Manual Assign).
- [ ] **RR-24.4**: [Integration] Kiểm thử luồng liên thông Warehouse -> Logistics, đảm bảo BatchID được kế thừa đúng.
