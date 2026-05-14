# Business Logic: Coffee Batch Lifecycle & Traceability

Tài liệu này định nghĩa các quy tắc nghiệp vụ cốt lõi cho việc quản lý mẻ hàng (Batch) trong hệ thống Runtime Roasters.

## 1. Định nghĩa Mẻ hàng (What is a Batch?)

Một "Mẻ hàng" đại diện cho một khối lượng hạt cà phê cụ thể được xử lý cùng nhau trong một công đoạn. Tính truy xuất nguồn gốc (Traceability) phải được bảo toàn từ lúc thu hoạch (Farm) đến khi đóng gói (Warehouse) và giao hàng (Retail).

## 2. Quy tắc Định danh (Batch ID Generation)

Mã Batch ID được sinh ra tự động để định danh duy nhất mẻ hàng tại từng giai đoạn.

**Định dạng:** `RR-{ServiceCode}-{OriginCode}-{YYYYMMDD}-{Sequence}`

| Thành phần | Ý nghĩa | Ví dụ |
| :--- | :--- | :--- |
| `RR` | Tiền tố dự án | `RR` |
| `ServiceCode` | `H` (Harvest), `P` (Processing), `S` (Stocked) | `P` |
| `OriginCode` | Mã vùng (e.g., `CD`: Cầu Đất, `BMT`: Buôn Ma Thuột) | `CD` |
| `YYYYMMDD` | Ngày tạo mã | `20260515` |
| `Sequence` | Số thứ tự trong ngày (4 chữ số) | `0001` |

## 3. Vòng đời của Mẻ hàng (Batch Lifecycle)

### Giai đoạn 1: Thu hoạch (Farm Service)
- Khi Farmer xác nhận thu hoạch, một mã `RR-H-...` được sinh ra.
- Trạng thái: `HARVESTED`.
- Thông tin đi kèm: Loại hạt, khối lượng thô, ngày thu hoạch, nông trại.

### Giai đoạn 2: Tiếp nhận & Sơ chế (Warehouse Service)
- Warehouse nhận sự kiện từ Farm, tạo bản ghi mẻ chế biến.
- Trạng thái: `RECEIVED`.
- Mã chuyển đổi: Từ `RR-H-...` sang `RR-P-...` (vẫn giữ link tới mã Harvest gốc).

### Giai đoạn 3: Chế biến (Roasting)
- Chuyển trạng thái sang `HULLING` (Xát vỏ) -> `ROASTING` (Rang).
- **Quy tắc Hao hụt (Weight Loss):** Quá trình rang làm giảm 12% - 20% trọng lượng. Hệ thống tự động tính toán `Expected Yield` dựa trên `Intake Weight`.
- Nếu trọng lượng đầu ra lệch quá 5% so với `Expected Yield`, yêu cầu Operator nhập lý do (Anomaly Note).

### Giai đoạn 4: Nhập kho thành phẩm (Stocking)
- Sau khi chế biến xong, mẻ hàng được đóng gói.
- Mã chuyển đổi: Sang `RR-S-...` (Stocked).
- Trạng thái: `STOCKED`.
- Dữ liệu tồn kho (Inventory) được cập nhật tăng cho SKU tương ứng.

## 4. Quy tắc Tồn kho (Inventory Rules)

- **Available Quantity:** Trọng lượng thực tế có thể bán (`Total Stocked` - `Reserved`).
- **Reserved Quantity:** Trọng lượng bị "treo" bởi các đơn hàng đang chờ xử lý (Saga).
- **FIFO (First-In, First-Out):** Ưu tiên xuất kho các mẻ hàng (`RR-S-...`) có ngày sản xuất cũ hơn để đảm bảo độ tươi mới.

---
**BA Approval Required for any changes to these rules.**
