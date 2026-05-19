# Sprint 4: The Resilient Farm (Transactional Outbox)

**Status:** ✅ Completed
**Goal:** Ensure every harvest batch is recorded with 100% data integrity even when the distributed system encounters failures. Implement the Transactional Outbox Blueprint as a project-wide pattern.

---

## 📋 Ticket Status (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-4.0](./RR-4.0/ticket.md) | [Tech] Refactor Farm Service: Migration & DI cleanup | ✅ Done | Tech Lead |
| [RR-4.1](./RR-4.1/ticket.md) | [BA] Harvest Batch Declaration (Harvesting Management) | ✅ Done | Farm Manager |
| [RR-4.2](./RR-4.2/ticket.md) | [Tech] Transactional Outbox: Event Reliability | ✅ Done | Tech Lead |
| [RR-22](./RR-22.md) | [Tech] Trace Service & CQRS Bootstrap | ✅ Done | Backend |

---

## 💡 Business Vision
- **Data Integrity:** "A single bean drops, the system knows." Absolutely no harvest data should be lost when transitioning to the factory.
- **Real-time Awareness:** Downstream departments (Warehouse) receive notifications as soon as goods leave the farm.
- **Professionalism:** Adopt the CloudEvents 1.0 standard for professional and extensible inter-service communication.

## 📊 Sprint Result
- Successfully implemented the `outbox_events` table and Relay Worker in the Farm Service.
- Completed the Harvesting API integrated with the GORM Transaction mechanism (Atomic Write).
- Harvest events are successfully published to Kafka and ready for consumption by the Warehouse.
