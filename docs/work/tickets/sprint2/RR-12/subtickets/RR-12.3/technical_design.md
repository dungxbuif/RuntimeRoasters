# RR-12.3: Auth Service Event-Driven Publishing Design

## 1. Mục tiêu (Goal)
Triển khai cơ chế "Single Writer" phát tín hiệu thay đổi chính sách thông qua Kafka để đồng bộ hóa trạng thái với các "Distributed Readers".

## 2. Kafka Configuration
- **Topic Name:** `auth.policy.changed`
- **Partitions:** 1 (để đảm bảo thứ tự các lệnh thay đổi nếu cần mở rộng sau này).
- **Format:** JSON.
- **Payload mẫu:** `{"action": "RELOAD", "timestamp": 123456789}`.

## 3. Implementation Logic

### 3.1 Kafka Writer Wrapper
Tại `internal/repository/kafka/publisher.go`:
```go
type PolicyPublisher struct {
    writer *kafka.Writer
}

func (p *PolicyPublisher) PublishReloadSignal(ctx context.Context) error {
    msg := kafka.Message{
        Value: []byte(`{"action": "RELOAD"}`),
    }
    return p.writer.WriteMessages(ctx, msg)
}
```

### 3.2 Intercepting Write Operations
Vì Auth Service là nơi duy nhất thực hiện thay đổi, chúng ta cần đảm bảo mọi thay đổi đều đi kèm với một tín hiệu Kafka.

Cách thực hiện tốt nhất: Tạo một `AuthService` struct bao bọc `casbin.Enforcer`:
```go
type AuthServiceImpl struct {
    enforcer  *casbin.Enforcer
    publisher *kafka.PolicyPublisher
}

func (s *AuthServiceImpl) AddPolicy(sub, obj, act string) error {
    // 1. Lưu vào DB (qua Enforcer)
    ok, err := s.enforcer.AddPolicy(sub, obj, act)
    if err != nil || !ok {
        return err
    }
    
    // 2. Phát tín hiệu qua Kafka
    return s.publisher.PublishReloadSignal(context.Background())
}
```

## 4. Resilience & Error Handling
- **Tách biệt lỗi:** Nếu `AddPolicy` vào DB thành công nhưng Kafka thất bại, hàm vẫn nên trả về thành công hoặc log warning. Reader sẽ dựa vào "Trụ cột 3" (Polling) để tự chữa lành dữ liệu sau vài phút.
- **Async Publishing:** Để không làm chậm API write, có thể gửi message Kafka vào một channel nội bộ để background worker xử lý.

## 5. Xác minh (Verification)
1. Sử dụng Kafka CLI tool (như `kcat` hoặc `kafka-console-consumer`) để lắng nghe topic `auth.policy.changed`.
2. Thực hiện gọi API thêm policy tại Auth Service.
3. Kiểm tra xem message "RELOAD" có xuất hiện trong Kafka hay không.
