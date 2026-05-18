# RuntimeRoasters Roadmap

## Phases

### Phase 1: Foundation (Completed)
**Goal:** Setup basic infrastructure and Auth.

### Phase 2: Farm Core (Completed)
**Goal:** Implement basic Farm CRUD.

### Phase 3: Reliable Harvesting (Completed)
**Goal:** Implement Outbox for Farm Service.

### Phase 4: Warehouse Foundation (Completed)
**Goal:** Implement Warehouse service and Inventory.

### Phase 5: Processing & State Machine (Completed)
**Goal:** Coffee processing flow.

### Phase 6: Retail & Order Saga
**Goal:** Build the Retail Service and Saga Orchestrator for Order management.
**Requirements:** [RR-22, RR-23]
**Plans:** 4 plans
- [ ] 06-01-PLAN.md — Retail Service Core & Order Management
- [ ] 06-02-PLAN.md — Outbox Relay & Kafka Integration
- [ ] 06-03-PLAN.md — Saga Choreography & State Machine
- [ ] 06-04-PLAN.md — Delivery & Integration

### Phase 7: Real-time Logistics
**Goal:** Build the Logistics Service and Driver Tracking system.
**Requirements:** [RR-24, RR-25]
**Plans:** 4 plans
- [ ] 07-01-PLAN.md — Foundation, DB/Valkey Setup, and Shipment Domain
- [ ] 07-02-PLAN.md — Kafka Consumers (Warehouse events) and Shipment Assignment
- [ ] 07-03-PLAN.md — GPS Tracking (Valkey GEO) and Real-time Updates
- [ ] 07-04-PLAN.md — Shipment State Machine and Delivery Confirmation

### Phase 8: Financial Integrity
**Goal:** Payment Service and Stripe.
