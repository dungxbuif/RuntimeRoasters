# Sprint 6: Retail & Order Saga

**Epic Goal:** Build an automated ordering and supply chain coordination system using the Saga Pattern.

---

## 📋 Tickets

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-21.5](./RR-21.5.md) | [Tech] Retail Store Seeding & Management | 🕒 To Do | Tech Lead |
| [RR-22](./RR-22.md) | [BA] Point of Sale (POS) Ordering System | 🕒 To Do | Store Manager |
| [RR-23](./RR-23.md) | [Tech] Retail Service & Saga Orchestrator | 🕒 To Do | Tech Lead |

---

## 🛠️ Technical Focus
- **Saga Choreography:** Coordinate flows via Kafka Events (Order -> Warehouse -> Logistics).
- **Transactional Outbox:** Ensure saving the Order and emitting the Event is a single Transaction.
- **Idempotency:** Use `Idempotency-Key` for APIs and `Message_ID` for Consumers.
