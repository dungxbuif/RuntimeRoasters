# Sprint 5: Warehouse & Inventory Core (Traceability)

**Status:** 🚧 In Progress
**Goal:** Focus on Warehouse operations, managing the flow of goods from the Farm to the Warehouse, performing high-level Processing, and packaging into Production Batches to prepare for sales.

---

## 📋 Ticket Status (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-19](./RR-19/ticket.md) | [Epic] Warehouse & Value Chain Management (Intake to Batch) | 🚧 In Progress | BA/PO |
| [RR-20](./RR-20/technical_design.md) | [Tech] Warehouse Scaffolding & State Machine | 🕒 To Do | Tech Lead |
| [RR-21](./RR-21/technical_design.md) | [Tech] Inventory Reservation (Saga Participant) | 🕒 To Do | Tech Lead |

---

## 💡 Business Vision
- **Inventory Fidelity:** Ensure raw inventory (Green Beans) and finished inventory (Roasted Beans) are always accurate in real-time.
- **Batch Aggregation:** Apply the model **1 Production Batch = N Roast Runs**. Allows merging several small roast runs into a single large commercial lot.
- **Simplified Processing:** Do not dive deep into temperature sensors; focus on state transitions and output volume confirmation.

## 📊 Expected Results (Sprint Goals)
- Receive harvest notifications from the Farm and automatically create intake records.
- Transition batch statuses through stages: `RECEIVED` -> `PROCESSING` -> `STOCKED`.
- Close production lots, automatically calculate total volume from small roast runs, and issue **Production Batch IDs**.
