# Sprint 4: The Resilient Farm

**Epic Goal:** Đảm bảo mọi mẻ thu hoạch được ghi nhận 100% không mất dữ liệu ngay cả khi hệ thống phân tán gặp sự cố.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-4.0](./RR-4.0/ticket.md) | [Tech] Refactor Farm Service (Sprint 3 Fixes) | 🕒 To Do | Tech Lead |
| [RR-4.1](./RR-4.1/ticket.md) | [BA] Khai báo mẻ thu hoạch thông minh | 🕒 To Do | Farmer |
| [RR-4.2](./RR-4.2/ticket.md) | [Tech] Transactional Outbox Blueprint | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Transactional Outbox:** Đảm bảo Atomicity giữa DB update và Event publishing.
- **Relay Worker:** Cơ chế polling/streaming dữ liệu từ bảng Outbox sang Kafka.
- **CloudEvents Standard:** Áp dụng chuẩn CloudEvents 1.0 cho mọi message payload.
