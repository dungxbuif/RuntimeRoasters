# Sprint 12: Hardening & Grand Finale

**Epic Goal:** Thắt chặt bảo mật toàn hệ thống, xử lý nợ kỹ thuật tồn đọng và trình diễn khả năng tự phục hồi (Resiliency).

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-12.1](./RR-12.1/ticket.md) | [BA] Trực quan hóa Control Plane & Chaos Demo | 🕒 To Do | Admin |
| [RR-12.2](./RR-12.2/ticket.md) | [Tech] Zero Trust with mTLS & Chaos Engineering | 🕒 To Do | Tech Lead |
| [RR-13](../sprint2/RR-13/ticket.md) | [Tech] Security: End-to-End Auth Integration | 🕒 To Do | Tech Lead |
| [RR-14](../sprint2/RR-14/ticket.md) | [Tech] Security: Token Revocation (Logout) | 🕒 To Do | Tech Lead |
| [RR-36](../sprint4/RR-36.md) | [Tech] PgBouncer Integration & Connection Pooling | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Mutual TLS (mTLS):** Bảo mật mọi kết nối gRPC nội bộ.
- **Chaos Engineering:** Mô phỏng chết Broker, chết Service để demo tính HA.
- **System Hardening:** Hoàn thiện cơ chế Revoke Token/Logout và tối ưu hóa DB Connection Pool qua PgBouncer.
