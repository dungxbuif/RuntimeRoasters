# Kafka Engineering Conventions (HA & Scalability)

Goal: Ensure the system achieves High Availability (HA) and safe Scale-out capabilities in a multi-pod environment.

---

## 1. Triple-Shield Architecture

### Shield 1: Partition Keys (Stream Consistency)
- **Rule:** Every event belonging to the same entity (Aggregate) **MUST** have the same Partition Key.
- **Implementation:** 
    - Producer (Farm): `msg.Key = []byte(harvest_id)`.
    - Result: All events for a single batch will always land in the same Partition and be processed sequentially by the same Pod.

### Shield 2: Inbox Pattern (Idempotency - Anti-duplication)
- **Rule:** Every received message must be checked for uniqueness before processing business logic.
- **Implementation:**
    - Use an `inbox_events` table (UNIQUE `message_id`).
    - Store the `message_id` and execute Business Logic within the same Database Transaction.

### Shield 3: Distributed Locking (Shared Resource Protection)
- **Rule:** When updating shared resources (e.g., Total SKU Inventory), a distributed lock must be used.
- **Implementation:**
    - Tool: **Valkey Redlock** (`github.com/go-redsync/redsync`).
    - Key format: `lock:inventory:{sku}`.

---

## 2. Infrastructure Configuration (Production-ready)
- **Replication Factor:** 3.
- **Min In-sync Replicas:** 2.
- **Acks:** `all` (Ensure absolute data safety for the supply chain).
- **Consumer Group:** Each service uses a unique `group.id` (e.g., `warehouse-service-group`).
