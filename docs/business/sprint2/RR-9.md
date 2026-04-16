# [RR-9] Kafka & Transactional Outbox Setup

- **Summary:** Thiết lập Kafka Cluster (Redpanda) và triển khai Transactional Outbox Pattern tại Farm Service.
- **Priority:** `CRITICAL`
- **Depends on:** `RR-1` (Infrastructure), `RR-6` (Farm Business Logic)
- **Description:** Đây là nền tảng của Sprint 2. Mọi sự kiện nghiệp vụ (thu hoạch mẻ mới, cập nhật lô hàng) phải được đảm bảo đến Kafka thông qua cơ chế Outbox — không mất dữ liệu kể cả khi broker tạm thời sập.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Outbox ghi event cùng transaction nghiệp vụ
- **Given:** Farm Service nhận request tạo mẻ thu hoạch mới.
- **When:** Business logic lưu `HarvestBatch` vào Postgres.
- **Then:** Một bản ghi `OutboxEvent{topic: "farm.harvest.created", payload: ...}` phải được ghi vào bảng `outbox` **trong cùng một DB transaction**. Nếu transaction rollback, outbox record cũng phải mất theo.

### Scenario 2: Outbox Worker relay lên Kafka
- **Given:** Có bản ghi pending trong bảng `outbox`.
- **When:** Outbox Worker chạy (polling interval: 5s).
- **Then:** Worker publish event lên Kafka topic `farm.harvest.created` với `acks=all`, sau đó đánh dấu record là `published=true`. Nếu Kafka unavailable, Worker retry với Exponential Backoff — **không xóa record**.

### Scenario 3: Idempotency khi Worker chạy song song
- **Given:** Triển khai 2 replicas của Farm Service (2 Outbox Workers).
- **When:** Cả 2 Worker cùng poll bảng `outbox`.
- **Then:** Sử dụng `SELECT ... FOR UPDATE SKIP LOCKED` để đảm bảo mỗi event chỉ được xử lý bởi 1 Worker. Không có event nào bị publish 2 lần lên Kafka.

---

## 🛠️ Technical Notes

- **Outbox Table Schema:**
  ```sql
  CREATE TABLE outbox_events (
      id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      topic       VARCHAR(255) NOT NULL,
      payload     JSONB        NOT NULL,
      published   BOOLEAN      DEFAULT FALSE,
      created_at  TIMESTAMPTZ  DEFAULT NOW(),
      published_at TIMESTAMPTZ
  );
  CREATE INDEX idx_outbox_unpublished ON outbox_events (published, created_at) WHERE published = FALSE;
  ```
- **Kafka Config:** `acks=all`, `enable.idempotence=true`, `Topic Replication Factor=3` (theo HA Design).
- **Worker:** Goroutine riêng biệt trong Farm Service process. Không phải service mới.

## 📋 Sub-tasks
- [ ] Tạo bảng `outbox_events` trong Farm DB migration.
- [ ] Implement `OutboxRepository` interface tại `pkg/outbox`.
- [ ] Implement Outbox Worker (Goroutine) với `SELECT FOR UPDATE SKIP LOCKED`.
- [ ] Wrap `CreateHarvestBatch` use-case để ghi outbox trong cùng transaction.
- [ ] Integration test: Tắt Kafka → Tạo batch → Bật Kafka → Kiểm tra event đến.
