# Sprint 10: Emergency Final Demo Sprint

**Goal:** Connect the Farm -> Warehouse -> Retail demo path with real runtime services, role-scoped actions, logistics simulation, SAGA events, traceability, and audit/realtime evidence.

---

## 📋 Ticket Status (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-URG-01](./RR-URG-01/ticket.md) | Verify And Fix End-to-End TraceId Propagation | ✅ Done | Tech Lead |
| [RR-URG-02](./RR-URG-02/ticket.md) | Event Contracts, IDs, And Storage Boundaries | ✅ Done | Tech Lead |
| [RR-URG-03](./RR-URG-03/ticket.md) | Warehouse HTTP APIs And Harvest-To-Pickup Flow | ✅ Done | Warehouse Mgr |
| [RR-URG-04](./RR-URG-04/ticket.md) | Logistics Shipment Legs, Driver Client Updates, And Mandatory Return | ✅ Done | Logistics Mgr |
| [RR-URG-05](./RR-URG-05/ticket.md) | Paid Order Fulfillment To Warehouse Reservation And Delivery | ✅ Done | Tech Lead |
| [RR-URG-06](./RR-URG-06/ticket.md) | Realtime Notifications And Socket/SSE Broadcasts | ✅ Done | Backend |
| [RR-URG-07](./RR-URG-07/ticket.md) | Public Sold-Cup QR Trace Umbrella | 🕒 In Progress | Tech Lead |
| [RR-URG-08](./RR-URG-08/ticket.md) | Role Authorization, Demo Accounts, And Seeded Journey Data | 🕒 To Do | Security Eng |
| [RR-URG-09](./RR-URG-09/ticket.md) | Client Dashboards, Driver Simulation UI, And Public QR Trace Page | 🕒 To Do | Frontend |
| [RR-URG-10](./RR-URG-10/ticket.md) | Demo Runbook And End-To-End Verification | 🕒 To Do | Team |

---

## 🛠️ Technical Focus
- **Distributed Tracing:** Connect gateway, services, Kafka, and Postgres/Elasticsearch read models with W3C propagation.
- **Saga Orchestration:** Choreographed SAGA using Kafka event publishing with transactional outbox logic.
- **Real-time Awareness:** Live WebSocket socket service broadcasting notifications and tracking driver simulation.
- **Data Provenance:** Tracing coffee menu inventory items to source lots, batches, and farms via public QR codes.
