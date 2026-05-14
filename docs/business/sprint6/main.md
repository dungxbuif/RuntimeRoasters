# Sprint 6: Retail & Order Saga

**Epic Goal:** Xây dựng hệ thống đặt hàng và điều phối chuỗi cung ứng tự động qua Saga Pattern.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-22](./RR-22.md) | [BA] Hệ thống đặt hàng tại quầy (POS) | 🕒 To Do | Store Manager |
| [RR-23](./RR-23.md) | [Tech] Retail Service & Saga Orchestrator | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Saga Choreography:** Điều phối luồng qua Kafka Events (Order -> Warehouse -> Logistics).
- **Transactional Outbox:** Đảm bảo lưu Order và bắn Event là một Transaction.
- **Idempotency:** Sử dụng `Idempotency-Key` cho API và `Message_ID` cho Consumer.
