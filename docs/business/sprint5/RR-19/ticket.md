# [RR-19] [Epic] Chuỗi cung ứng hạt: Từ Nông trại đến Kho thành phẩm

**User Story:**
Dưới vai trò là **Giám đốc Vận hành (COO)**, tôi muốn có một hệ thống quản lý mẻ hàng xuyên suốt từ khi nhận nguyên liệu thô đến khi ra thành phẩm đóng gói, để tôi có thể kiểm soát chất lượng, tỷ lệ hao hụt và đảm bảo luôn có đủ hàng để bán.

**Business Context:**
Hiện tại, việc chuyển giao giữa nông trại và nhà máy đang được ghi chép thủ công, dẫn đến thất thoát dữ liệu và khó khăn khi truy vết nếu mẻ hàng gặp lỗi. Việc số hóa quy trình này là điều kiện tiên quyết để mở rộng quy mô sản xuất.

---

## 📋 Danh sách Sub-tickets

| Ticket | Summary | Priority |
| :--- | :--- | :--- |
| [RR-19.1](./subtickets/RR-19.1/ticket.md) | Tiếp nhận mẻ thu hoạch & Khởi tạo mã Batch | `CRITICAL` |
| [RR-19.2](./subtickets/RR-19.2/ticket.md) | Quản lý quy trình chế biến (Rang/Sấy) & Tính hao hụt | `HIGH` |
| [RR-19.3](./subtickets/RR-19.3/ticket.md) | Nhập kho thành phẩm & Cập nhật tồn kho tự động | `CRITICAL` |

---

## 🛠️ Quy tắc Nghiệp vụ Chung (Business Rules)
- Mọi mẻ hàng khi chuyển trạng thái phải ghi lại **Timestamp** và **Operator ID**.
- Mã Batch ID phải tuân thủ đúng định dạng tại `docs/business/BATCH_LOGIC.md`.
- Hệ thống không cho phép xóa Batch đã có dữ liệu chế biến (Chỉ cho phép Cancel với lý do cụ thể).

---

## ✅ Acceptance Criteria (Epic Level)
1. **End-to-End Traceability:** Từ một gói cà phê bất kỳ, hệ thống phải truy ngược được về mẻ rang tương ứng và mẻ thu hoạch gốc tại nông trại nào.
2. **Weight Integrity:** Tổng trọng lượng thành phẩm cộng với hao hụt phải khớp với trọng lượng đầu vào (sai số cho phép < 1%).
3. **Inventory Sync:** Số lượng tồn kho trên Dashboard phải cập nhật ngay lập tức khi một mẻ hàng chuyển sang trạng thái `STOCKED`.
