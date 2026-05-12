# Sprint 10: Traceability (CQRS)

**Epic Goal:** Cung cấp khả năng truy xuất nguồn gốc 360 độ từ hạt cà phê đến tách cà phê với tốc độ tối ưu.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-10.1](./RR-10.1/ticket.md) | [BA] Dashboard truy xuất nguồn gốc (Bean Journey) | 🕒 To Do | End User |
| [RR-10.2](./RR-10.2/ticket.md) | [Tech] CQRS Read-Model with Elasticsearch | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **CQRS Pattern:** Tách biệt luồng ghi (Postgres) và luồng đọc (Elasticsearch).
- **Denormalization:** Tổng hợp dữ liệu từ nhiều service vào một document duy nhất.
- **Trace Propagation:** Sử dụng OpenTelemetry để gắn kết các bước trong chuỗi.
