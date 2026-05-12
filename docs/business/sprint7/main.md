# Sprint 7: Retail & Order Saga

**Epic Goal:** Xây dựng hệ thống đặt hàng và điều phối chuỗi cung ứng tự động qua Saga Pattern.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-7.1](./RR-7.1/ticket.md) | [BA] Hệ thống đặt hàng tại quầy (POS) | 🕒 To Do | Store Manager |
| [RR-7.2](./RR-7.2/ticket.md) | [Tech] Saga Orchestrator (Choreography) | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Saga Choreography:** Điều phối luồng qua Kafka Events (Order -> Warehouse -> Logistics).
- **Idempotent Consumer:** Đảm bảo mỗi sự kiện chỉ được xử lý đúng một lần.
- **Compensating Actions:** Cơ chế hoàn tác nếu một bước trong chuỗi thất bại.
