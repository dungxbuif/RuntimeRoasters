# Sprint 11: Observability & Audit

**Epic Goal:** Giám sát "sức khỏe" hệ thống và lưu trữ bằng chứng giao dịch không thể sửa đổi.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-11.1](./RR-11.1/ticket.md) | [BA] Hệ thống giám sát và Audit Trail | 🕒 To Do | Auditor |
| [RR-11.2](./RR-11.2/ticket.md) | [Tech] Audit Service with Apache Cassandra | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Apache Cassandra:** Lưu trữ Event Log khổng lồ với khả năng write-heavy.
- **Hash-chained Events:** Đảm bảo tính toàn vẹn của Audit Log (chống sửa đổi).
- **Unified Dashboard:** Tích hợp Prometheus, Grafana và Signoz.
