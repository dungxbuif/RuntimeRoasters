# Ticket RR-5.2: [Tech] Roastery Service Initialization & State Machine

**Mục tiêu:** Xây dựng khung dịch vụ Roastery và bộ máy quản lý trạng thái sản xuất.

---

## 🛠️ Yêu cầu Kỹ thuật

### 1. Service Scaffolding
- Tạo folder `src/apps/roastery-process-service`.
- Thiết lập `main.go`, Manual DI, và các layer `domain`, `usecase`, `infrastructure` (Clean Architecture).
- Cấu hình gRPC Server và kết nối PostgreSQL.

### 2. State Machine Implementation
- Hiện thực hóa logic kiểm tra `validTransitions` trong Domain layer.
- Đảm bảo tính nhất quán (Consistency): Trạng thái chỉ được thay đổi thông qua UseCase hợp lệ.

### 3. Kafka Consumer & Inbox Pattern
- Viết Consumer lắng nghe topic `farm.harvest.events`.
- Implement `InboxRepository` để lưu vết `message_id`.
- Đảm bảo Idempotency: Không xử lý lại các sự kiện đã có trong Inbox.

### 4. Batch ID Generator
- Viết utility hoặc service để sinh mã theo định dạng đã thỏa thuận.
- Đảm bảo mã sinh ra là duy nhất (Unique) trong hệ thống.

---

## 🧪 Kiểm thử (Verification)
- **Unit Test State Machine:** Test tất cả các trường hợp chuyển trạng thái (Valid & Invalid).
- **Integration Test Consumer:** Gửi một message giả lập vào Kafka và kiểm tra DB Roastery có tự tạo bản ghi `processing_batches` không.
- **Concurrency Test:** Giả lập 2 request cập nhật trạng thái cùng lúc cho 1 batch để verify tính toàn vẹn.
