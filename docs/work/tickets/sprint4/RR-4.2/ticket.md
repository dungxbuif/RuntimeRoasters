# Ticket RR-4.2: [Tech] Transactional Outbox Blueprint

**Mục tiêu:** Xây dựng cơ chế "Chốt chặn dữ liệu" (Data Safety Net) để đảm bảo tính nhất quán giữa Database và Message Broker (Kafka).

---

## 🛠️ Yêu cầu Kỹ thuật

### 1. Database Level
- Triển khai Migration tạo bảng `outbox_events`.
- Đảm bảo bảng này nằm trong cùng Database Schema với `harvests`.

### 2. Repository Pattern
- Implement phương thức lưu trữ nguyên tử (Atomic):
  ```go
  func (r *repository) CreateHarvest(ctx context.Context, harvest *domain.Harvest, event *domain.OutboxEvent) error {
      return r.db.Transaction(func(tx *gorm.DB) error {
          // 1. Lưu Harvest
          // 2. Lưu Outbox Event
      })
  }
  ```

### 3. Background Relay Service
- Một Go routine chạy độc lập trong `farm-service`.
- **Nhiệm vụ:**
    - Polling bảng `outbox_events` theo batch (ví dụ: 20 records/lần).
    - Publish lên Kafka Topic `farm.harvest.events`.
    - Đánh dấu `processed_at` ngay sau khi nhận được Ack từ Kafka.
- **Yêu cầu:** Không được để rò rỉ bộ nhớ (Memory Leak) và phải handle được tín hiệu `Graceful Shutdown`.

### 4. CloudEvents Wrapping
- Sử dụng library `github.com/cloudevents/sdk-go` hoặc cấu trúc tương đương.
- Header Kafka phải chứa `traceparent` để hỗ trợ Distributed Tracing.

---

## 🧪 Kiểm thử (Verification)
- **Unit Test:** Mock DB Transaction để verify cả 2 bản ghi đều được gọi.
- **Integration Test:**
    1. Start `farm-service` nhưng KHÔNG start Kafka.
    2. Gọi API tạo Harvest -> Kiểm tra DB thấy `processed_at` là NULL.
    3. Start Kafka -> Kiểm tra sau vài giây `processed_at` được cập nhật và Kafka nhận được message.
