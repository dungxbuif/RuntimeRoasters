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
    - Derive `message_id` using the shared helper, not raw offsets:
        1. `topic:event_id` when the payload has `event_id`.
        2. `topic:key` when Kafka key is present.
        3. `topic-partition-offset` only as a legacy fallback.

**Important:** Kafka offsets are not stable across local broker recreation, topic resets, or demo-state resets. Consumers must not use only `topic-partition-offset` for business idempotency. Use `pkg/kafka.MessageID(msg)` unless there is a documented reason not to.

### Shield 2.5: Trace Context Across Async Boundaries
- **Rule:** Every Kafka message produced inside an active request/SAGA must carry W3C `traceparent`.
- **Implementation:**
    - Producers use `pkg/kafka.NewProducer`, which injects the current OTel context.
    - Consumers use `pkg/kafka.NewConsumer`, which extracts `traceparent` before invoking handlers.
    - Transactional outbox rows must persist `traceparent`/`tracestate` because the relay runs later and cannot rely on the original request context.

**Important:** A relay publishing from `context.Background()` without restoring outbox trace metadata breaks the distributed trace.

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
