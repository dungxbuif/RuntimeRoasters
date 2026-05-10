# RR-12.4: Resilient Reader Engine Design (pkg/base/casbin)

## 1. Mục tiêu (Goal)
Xây dựng một "siêu Enforcer" có khả năng tự đồng bộ và duy trì hoạt động trong môi trường phân tán nhiều biến động, đảm bảo tính sẵn sàng cao (High Availability).

## 2. Thiết kế chi tiết 3 trụ cột (3-Pillar Strategy)

### 2.1 Trụ cột 1: Bootstrapping (gRPC Snapshot)
- **Hành động:** Khi khởi tạo `NewResilientEngine`, engine sẽ thực hiện một cuộc gọi gRPC đồng bộ tới `AuthService.GetPolicies`.
- **Resilience:** Sử dụng thư viện `github.com/cenkalti/backoff/v4` để retry. Nếu sau 5 phút vẫn không thể kết nối tới Auth Service, Engine có thể nạp một bộ chính sách mặc định (fail-safe) để tránh làm chết service Reader.

### 2.2 Trụ cột 2: Real-time Update (Kafka Watcher)
- **Hành động:** Engine chạy một goroutine lắng nghe topic `auth.policy.changed`.
- **Logic:** Khi nhận được message "RELOAD", engine không parse dữ liệu từ Kafka (để tránh message quá lớn) mà đơn giản là kích hoạt hàm `syncFromGRPC()` để làm mới dữ liệu.

### 2.3 Trụ cột 3: Safety Net (Polling)
- **Hành động:** Sử dụng `time.Ticker` chạy định kỳ mỗi 5 phút.
- **Logic:** Gọi `syncFromGRPC()` để đảm bảo nếu Kafka bị mất tin nhắn hoặc có lỗi không quán, dữ liệu sẽ được sửa chữa (Self-healing).

## 3. Cấu trúc Source Code
Vị trí: `src/pkg/base/casbin/`
```text
.
├── engine.go           # Logic điều phối 3 trụ cột
├── adapter_memory.go   # Custom Casbin adapter để load/store policy trong RAM
├── watcher_kafka.go    # Tích hợp Kafka consumer
└── config.go           # Định nghĩa các tham số (AuthService URL, Kafka Brokers)
```

## 4. Implementation Snippet (Pseudo-code)
```go
func (e *ResilientEngine) syncFromGRPC() {
    resp, err := e.grpcClient.GetPolicies(ctx, &pb.GetPoliciesRequest{})
    if err != nil {
        log.Error("Failed to sync policies", err)
        return
    }
    
    // Clear old policies and load new ones into Memory Adapter
    e.adapter.Clear()
    for _, rule := range resp.Rules {
        e.adapter.AddRule(rule)
    }
    e.enforcer.LoadPolicy()
}
```

## 5. Xác minh (Verification)
- **Unit Test:** Mock gRPC server và Kafka để kiểm tra Engine có gọi đúng hàm sync khi có sự kiện hay không.
- **Integration Test:** Chạy Auth Service và một Reader, thay đổi quyền ở Auth Service và kiểm tra log của Reader xem có nhận được update trong < 1s hay không.
