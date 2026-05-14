# Sprint 5: Warehouse & Inventory Core (Traceability)

**Trạng thái:** 🚧 Đang thực hiện (In Progress)
**Mục tiêu:** Tập trung vào nghiệp vụ Kho bãi (Warehouse), quản lý luồng luân chuyển hàng hóa từ Nông trại về Kho, thực hiện Chế biến (High-level) và đóng gói thành Lô sản phẩm (Production Batch) để sẵn sàng bán hàng.

---

## 📋 Trạng thái Ticket (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-19](./RR-19/ticket.md) | [Epic] Quản lý Kho & Chuỗi giá trị (Intake to Batch) | 🚧 In Progress | BA/PO |
| [RR-20](./RR-20/technical_design.md) | [Tech] Warehouse Scaffolding & State Machine | 🕒 To Do | Tech Lead |
| [RR-21](./RR-21/technical_design.md) | [Tech] Inventory Reservation (Saga Participant) | 🕒 To Do | Tech Lead |

---

## 💡 Tầm nhìn Nghiệp vụ (Business Vision)
- **Inventory Fidelity:** Đảm bảo tồn kho thô (Green Beans) và tồn kho thành phẩm (Roasted Beans) luôn chính xác theo thời gian thực.
- **Batch Aggregation:** Áp dụng mô hình **1 Production Batch = N Roast Runs (Mẻ)**. Cho phép gom nhiều mẻ rang nhỏ thành một lô hàng thương mại lớn.
- **Simplified Processing:** Không đi sâu vào cảm biến nhiệt độ, chỉ tập trung vào việc chuyển đổi trạng thái và xác nhận khối lượng đầu ra.

## 📊 Kết quả mong đợi (Sprint Goals)
- Nhận thông báo thu hoạch từ Farm và tự động tạo phiếu nhập kho.
- Chuyển trạng thái mẻ hàng qua các giai đoạn: `RECEIVED` -> `PROCESSING` -> `STOCKED`.
- Đóng lô sản xuất, tự động tính tổng khối lượng từ các mẻ rang nhỏ và cấp mã **Production Batch ID**.
