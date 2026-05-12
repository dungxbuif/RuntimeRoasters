# Sprint 6: Warehouse Inventory

**Epic Goal:** Quản lý kho trung tâm và cung cấp cơ chế giữ chỗ hàng (Reservation) an toàn cho các đơn hàng.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-6.1](./RR-6.1/ticket.md) | [BA] Quản lý tồn kho thành phẩm | 🕒 To Do | Warehouse Keeper |
| [RR-6.2](./RR-6.2/ticket.md) | [Tech] Saga Participant: Stock Reservation & Locking | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Distributed Locking:** Sử dụng Valkey để tránh bán quá số lượng (Overselling).
- **Saga Participant:** Hỗ trợ các API `Reserve`, `Confirm`, `Cancel` cho luồng Saga.
- **Inventory Snapshot:** Tối ưu hóa việc đọc tồn kho thực tế.
