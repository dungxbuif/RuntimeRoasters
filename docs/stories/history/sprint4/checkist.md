Luồng **Tạo lô thu hoạch (Create Harvest Batch)** chính là "Điểm khởi thủy" (Origin) của toàn bộ chuỗi cung ứng Runtime Roasters. Hạt cà phê bắt đầu vòng đời từ đây, nên dữ liệu của nó phải đạt độ chuẩn xác tuyệt đối.

Để hoàn thiện luồng này ở mức Enterprise, bạn cần "lắp ráp" tất cả những khái niệm chúng ta đã bàn luận vào một API duy nhất. Dưới đây là Checklist 5 lớp kiến trúc bạn cần apply cho luồng `CreateHarvestBatch`:

---

### Lớp 1: Cổng bảo vệ (Security & Zero-Trust)

Trước khi request chạm vào code nghiệp vụ, nó phải bị chặn lại ở tầng HTTP/gRPC Middleware.

* **Xác thực (Authentication):** Middleware bóc JWT từ Header, giải mã ra `user_id` và `role`.
* **Phân quyền (Casbin Authorization):** Middleware ném thông tin vào RAM của Casbin `enforcer.Enforce(role, "/api/v1/farm/harvests", "POST")`. Chỉ người có quyền `FARM_MANAGER` mới được đi tiếp.

### Lớp 2: Kiểm duyệt đầu vào (Edge & Domain Validation)

* **Edge (DTO):** Request JSON đi vào `CreateHarvestRequest`. Dùng `go-playground/validator` chặn ngay nếu truyền thiếu `FarmID` hoặc `Quantity = 0`.
* **Domain:** Map DTO sang `domain.Harvest`. Gọi hàm `h.Validate()` bằng `ozzo-validation` để đảm bảo `Quantity` hợp lý (ví dụ không thể thu hoạch > 100 tấn/ngày cho 1 nông trại nhỏ) và gán trạng thái mặc định là `StatusNew`.
* **Traceability (Truy xuất nguồn gốc):** Đây là hệ thống chuỗi cung ứng. Bạn nên apply logic **tự động sinh Batch Code** (Mã lô hàng) duy nhất ngay tại đây. Ví dụ: `HRV-20260513-UUID`.

### Lớp 3: Giao dịch Cơ sở dữ liệu (The Core)

Đây là nơi bạn apply **Transactional Outbox**. Bạn không được phép dùng `db.Create()` đơn thuần. Bạn phải mở một Transaction.

```go
// Tầng Repository/Infrastructure
func (r *harvestRepo) CreateHarvestWithOutbox(ctx context.Context, h *domain.Harvest) error {
    // 1. Chuyển Domain -> Model
    harvestModel := toHarvestModel(h)
    
    // 2. Mở Transaction
    return r.db.Transaction(func(tx *gorm.DB) error {
        // 3a. Lưu vào bảng harvests
        if err := tx.Create(harvestModel).Error; err != nil {
            return err
        }
        
        // 3b. Tạo payload cho Event
        eventPayload := map[string]interface{}{
            "harvest_id": harvestModel.ID,
            "farm_id":    harvestModel.FarmID,
            "quantity":   harvestModel.Quantity,
            "status":     harvestModel.Status,
        }
        payloadBytes, _ := json.Marshal(eventPayload)

        // 3c. Lưu vào bảng outbox_events
        outboxEvent := OutboxEventModel{
            AggregateType: "Harvest",
            AggregateID:   fmt.Sprintf("%d", harvestModel.ID),
            EventType:     "HarvestBatchCreated",
            Payload:       string(payloadBytes),
            Status:        "PENDING",
        }
        
        if err := tx.Create(&outboxEvent).Error; err != nil {
            return err // Nếu lỗi, cả Harvest và Outbox đều bị Rollback
        }
        
        return nil // Commit thành công!
    })
}

```

### Lớp 4: Phát sóng Sự kiện (Event Publishing)

Sau khi Transaction thành công, request của người dùng (HTTP 200 OK) sẽ được trả về ngay lập tức để UI không bị treo. Nhiệm vụ còn lại thuộc về Background Worker.

* **Outbox Relay:** Bạn apply một Goroutine chạy ngầm trong Farm Service. Cứ 1 giây nó quét bảng `outbox_events` tìm các dòng `PENDING`.
* Nó bắn event `HarvestBatchCreated` lên **Apache Kafka** (Topic: `farm.harvest.events`).
* Nhận được ACK từ Kafka -> Update dòng outbox thành `PUBLISHED`.

### Lớp 5: Hậu kiểm & Chuyển giao (Downstream Reaction)

Luồng này không nằm ở Farm Service, mà nằm ở **Warehouse / Roastery Service**.

* Warehouse có một Kafka Consumer lắng nghe topic `farm.harvest.events`.
* Khi thấy `HarvestBatchCreated` với khối lượng 500kg, Warehouse Service sẽ:
1. Ghi nhận một dự báo (Forecasting) hoặc một "Phiếu nhập kho dự kiến" (Inbound Order).
2. Gửi Notification (Socket/Email) cho Xưởng Trưởng: *"Nông trại vừa tạo lô thu hoạch 500kg, chuẩn bị dọn chỗ trong kho chứa hạt thô nhé!"*



---

### Tóm tắt lại "To-Do List" cho luồng này:

1. [x] Thiết kế file `.proto` cho `CreateHarvest` (Đã có).
2. [x] Viết Struct ở tầng `domain` với `ozzo-validation` (Đã có).
3. [ ] Tạo thư mục `infrastructure/db/models` và viết `HarvestModel`, `OutboxEventModel`.
4. [ ] Cấu hình GORM kết nối Postgres.
5. [ ] Viết file SQL Migration tạo 2 bảng trên (Nếu dùng `golang-migrate`).
6. [ ] Viết `harvest_repository.go` chứa hàm Transactional Outbox như trên.
7. [ ] (Nâng cao) Viết `outbox_worker.go` quét DB và in log ra màn hình (chưa cần nối Kafka vội, in ra terminal để test logic chạy ngầm trước).

Bạn muốn bắt tay vào code tầng nào tiếp theo? Tôi có thể cung cấp đoạn code hoàn chỉnh cho cái `outbox_worker.go` (Goroutine chạy ngầm) vì phần này xử lý Concurrency (Đồng thời) trong Go khá thú vị và dễ bị lỗi nếu không cẩn thận.