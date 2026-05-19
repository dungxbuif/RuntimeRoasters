# Phase 6: Retail & Order Saga - Research

**Researched:** 2024-05-24
**Domain:** Distributed Sagas, Retail POS, Transactional Outbox
**Confidence:** HIGH

## Summary
Phase 6 focuses on building the **Retail Service** and implementing the **Order Saga** using choreography. The Retail Service acts as the entry point for customer orders (POS), managing the `orders` state machine and ensuring consistency via the Transactional Outbox pattern. The Saga involves coordination between Retail (Ordering), Warehouse (Inventory Reservation), and Payment services.

**Primary recommendation:** Reuse the Outbox implementation from `farm-service` and the Kafka consumer patterns from `warehouse-service` to maintain architectural consistency. Use GORM for DB transactions to ensure atomic Order + Outbox writes.

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| Order Management | Retail Service | — | Authoritative source for order state. |
| Inventory Reservation| Warehouse Service| — | Owns the stock levels and locking logic. |
| Payment Processing | Payment Service | — | Integrates with Stripe; owns financial state. |
| Saga Orchestration | Retail Service | — | In a choreography model, Retail initiates and often monitors the final state. |

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| GORM | v1.25+ | ORM & Transactions | Existing standard in project (ADR-0004). |
| Kafka-go | v0.4+ | Event Streaming | Used in `pkg/kafka` for messaging. |
| Postgres | 15+ | Relational DB | Primary persistent store for services. |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|--------------|
| Valkey | 7.2+ | Distributed Locking | Used for SKU-level locking in Warehouse. |

## Architecture Patterns

### Recommended Project Structure
```
src/apps/retail-service/
├── cmd/
├── internal/
│   ├── app/         # Initialization
│   ├── domain/      # Order models, Outbox models
│   ├── usecase/     # Order creation, Saga response handling
│   ├── repo/        # GORM repositories
│   └── worker/      # Outbox Relay, Saga Consumers
```

### Pattern 1: Transactional Outbox
**What:** Writing an order and an event to the same DB in one transaction.
**When to use:** Every time an external event must be published based on a DB change.

### Pattern 2: Saga Choreography
**Flow:**
1. Retail: `OrderCreated` -> Published to Kafka.
2. Warehouse: Listens to `OrderCreated` -> Reserves Stock -> Publishes `InventoryReserved` (or `ReservationFailed`).
3. Payment: Listens to `InventoryReserved` -> Charges Card -> Publishes `PaymentCompleted`.
4. Retail: Listens to `PaymentCompleted` -> Status = `SUCCESS`.

## Common Pitfalls

### Pitfall 1: Dual Writes
**What goes wrong:** Updating the DB but failing to send the Kafka message.
**How to avoid:** Use the Transactional Outbox pattern (already implemented in `farm-service`).

### Pitfall 2: Lost Compensating Actions
**What goes wrong:** A failure occurs but the "undo" event isn't processed.
**How to avoid:** Ensure all consumers use the **Inbox Pattern** (idempotency) to safely retry failures.

## Assumptions Log
| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | Warehouse will use Valkey for locking | Architecture | Performance impact if using DB locks instead. |
| A2 | Choreography is preferred over Orchestration | Summary | May need a central orchestrator if logic gets too complex. |

## Open Questions
1. **Payment Webhook:** Will the Payment service be built in this phase or mocked? (Current roadmap suggests Phase 8, but Saga needs it). *Recommendation: Mock Payment success/fail events for Phase 6.*
