# Sprint 7: Real-time Logistics

**Epic Goal:** Manage transportation and real-time shipment tracking.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-24](./RR-24/ticket.md) | [BA] Transportation Dispatch Process | 🕒 To Do | Logistics Manager |
| [RR-25](./RR-25/ticket.md) | [Tech] Logistics Service & Driver Tracking | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Valkey GEO:** Store and query driver locations in real-time.
- **Service Integration:** Listen for events from Warehouse to trigger shipments.
- **State Machine:** Manage shipment statuses (PENDING -> ASSIGNED -> IN_TRANSIT -> DELIVERED).
- **Driver Simulator:** Script to simulate driver movement on a map.
