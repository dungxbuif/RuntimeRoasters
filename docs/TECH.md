# 🛠️ Runtime Roasters: Technical Details

Tài liệu này giải thích chi tiết các bài toán kỹ thuật trọng tâm trong hệ thống Runtime Roasters.

---

## 🚚 1. Bài toán Logistics: Real-time Get & Publish

### Luồng xử lý
Hệ thống Logistics sử dụng mô hình Event-Driven kết hợp với WebSocket để cập nhật trạng thái và vị trí tài xế thời gian thực.

- **Publish (Đẩy vị trí)**: 
  - Mobile App (hoặc Simulator) gửi tọa độ GPS về `logistics-service` qua REST API `POST /v1/logistics/gps`.
  - `logistics-service` cập nhật vị trí vào **Valkey** (Redis clone) để truy xuất cực nhanh và đồng thời đẩy một Event vào Kafka topic `logistics.gps.updated`.
  - `socket-service` lắng nghe topic này và broadcast qua **WebSocket** tới tất cả Clients đang xem bản đồ.

- **Get (Lấy thông tin)**:
  - Khi người dùng vào Dashboard, Client gọi `GET /v1/logistics/shipments` để lấy danh sách chuyến hàng và vị trí hiện tại của các tài xế (lấy từ cache Valkey).
  - Sau đó, kết nối WebSocket sẽ đảm nhận việc cập nhật "sống" các vị trí này mà không cần reload trang.

---

## 🔍 2. Bài toán Traceability: CQRS & Event Sourcing

Hệ thống truy xuất nguồn gốc (Traceability) được thiết kế theo nguyên lý CQRS (Command Query Responsibility Segregation).

- **Command Side (Write)**: Các dịch vụ nghiệp vụ (Farm, Warehouse, Retail) thực hiện các tác vụ thay đổi trạng thái (thu hoạch, nhập kho, bán hàng). Sau khi thay đổi DB thành công, họ đẩy một CloudEvent vào Kafka.
- **Event Storage**: Một `trace-service` chuyên biệt lắng nghe tất cả các topics quan trọng và lưu vết vào **Elasticsearch** (cho tìm kiếm nhanh) và **Cassandra** (cho tính bền vững và audit lâu dài).
- **Query Side (Read)**: 
  - Khi cần xem cây truy xuất của một ly cà phê, Client gọi `trace-service`. 
  - `trace-service` thực hiện truy vấn ngược (Recursive Search) trong Elasticsearch để dựng lại toàn bộ lịch sử từ Cup -> Retail -> Warehouse -> Farm.
  - Kết quả được trả về dưới dạng JSON đồ thị (Graph) để hiển thị trên UI.

---

## 📦 3. Bài toán Phê duyệt đơn hàng (Order Approval)

Quy trình xử lý đơn hàng từ Retailer về Warehouse tuân thủ mô hình SAGA (Orchestration-based).

1. **Khởi tạo**: `STORE_MGR` tạo đơn hàng trên Dashboard. `retail-service` lưu đơn hàng với trạng thái `PENDING` và phát event `retail.order.created`.
2. **Kiểm tra kho**: `warehouse-service` nhận event, thực hiện giữ chỗ (Reserve) hàng trong kho:
   - Nếu đủ hàng: Phát event `warehouse.stock.reserved`.
   - Nếu thiếu hàng: Phát event `warehouse.stock.reservation_failed`.
3. **Thanh toán**: `payment-service` nhận event `reserved`, thực hiện trừ tiền (hoặc giả lập thanh toán):
   - Thành công: Phát `payment.completed`.
   - Thất bại: Phát `payment.failed`.
4. **Phê duyệt & Điều phối**:
   - `retail-service` nhận `payment.completed` và chuyển trạng thái đơn hàng sang `APPROVED`.
   - `warehouse-service` nhận `payment.completed` và tạo một `DispatchRequest` trong hàng đợi chờ Warehouse Manager điều xe và tài xế.
5. **Giao hàng**: Warehouse Manager chọn Tài xế & Xe trong giao diện `Warehouse Ops`, gọi lệnh `Dispatch`. Một chuyến hàng (Shipment) được tạo ra trong `logistics-service`.
