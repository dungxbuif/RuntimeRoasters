# Sprint 7: Real-time Logistics

**Epic Goal:** Quản lý vận chuyển và theo dõi chuyến hàng thời gian thực.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-24](./RR-24.md) | [BA] Quy trình điều phối vận chuyển | 🕒 To Do | Logistics Manager |
| [RR-25](./RR-25.md) | [Tech] Logistics Service & Driver Tracking | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Valkey GEO:** Lưu trữ và truy vấn vị trí tài xế theo thời gian thực.
- **Service Integration:** Lắng nghe sự kiện từ Warehouse để kích hoạt chuyến hàng.
- **State Machine:** Quản lý trạng thái chuyến hàng (Pending -> In Transit -> Delivered).
