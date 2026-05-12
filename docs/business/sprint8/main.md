# Sprint 8: Real-time Logistics

**Epic Goal:** Theo dõi hành trình vận chuyển và cập nhật vị trí tài xế theo thời gian thực.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-8.1](./RR-8.1/ticket.md) | [BA] Điều phối và Theo dõi chuyến xe | 🕒 To Do | Logistics Manager |
| [RR-8.2](./RR-8.2/ticket.md) | [Tech] Real-time GPS Tracking with Valkey GEO | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Valkey GEO:** Lưu trữ và tính toán khoảng cách tọa độ GPS tốc độ cao.
- **Webhook Ingress:** Tiếp nhận dữ liệu GPS từ Simulator/App qua Webhook Service.
- **SSE/Websocket:** Đẩy vị trí trực tiếp lên dashboard quản trị.
