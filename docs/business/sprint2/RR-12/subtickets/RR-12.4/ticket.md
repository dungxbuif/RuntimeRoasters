# [RR-12.4] Base Package: Resilient Reader Engine

- **Mục tiêu:** Xây dựng một thư viện dùng chung (Engine) giúp các Reader service tự quản lý bộ nhớ đệm chính sách phân quyền một cách tin cậy.
- **Mô tả:** Triển khai "Resilience Engine" tại `pkg/base/casbin`. Đây là trái tim của hệ thống phân quyền phân tán, tích hợp cả 3 trụ cột đồng bộ: gRPC Snapshot, Kafka Live Update, và Polling Fallback.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Engine được triển khai tại `src/pkg/base/casbin/engine.go`.
2. **Trụ cột 1:** Tự động gọi gRPC tới Auth Service để lấy snapshot khi khởi động (có cơ chế Retry với Exponential Backoff).
3. **Trụ cột 2:** Tích hợp Kafka Consumer lắng nghe topic `auth.policy.changed` để gọi `LoadPolicy()` tức thì.
4. **Trụ cột 3:** Background Ticker tự động đồng bộ qua gRPC mỗi 5-10 phút (Safety net).
5. Hiệu năng: Phương thức `Enforce(sub, obj, act)` phải chạy hoàn toàn trong bộ nhớ (Local Memory), không gọi network.

## 🛠 Task list cho Developer
- [ ] Thiết lập struct `AuthEngine` bao bọc `casbin.Enforcer` và `MemoryAdapter`.
- [ ] Triển khai hàm `syncFromGRPC()` thực hiện call snapshot và cập nhật enforcer.
- [ ] Triển khai `watchKafka()` chạy trong goroutine để nhận tín hiệu reload.
- [ ] Triển khai `startPolling()` chạy trong goroutine để đồng bộ định kỳ.
- [ ] Cung cấp hàm khởi tạo `NewResilientEngine(cfg Config)` trả về interface Enforcer.
