# [RR-20] [Tech] Warehouse Service: Infrastructure & State Machine

**User Story:**
Dưới vai trò là **Tech Lead**, tôi muốn thiết lập dịch vụ Warehouse với kiến trúc Clean Architecture và bộ máy State Machine mạnh mẽ để hiện thực hóa các yêu cầu nghiệp vụ trong RR-19.

**Technical Value:**
Việc tách biệt Logic trạng thái giúp code dễ bảo trì và mở rộng khi quy trình sản xuất trở nên phức tạp hơn. Việc hợp nhất Process và Warehouse giúp giảm độ trễ và sự phức tạp của hệ thống.

---

## 🛠️ Yêu cầu Kỹ thuật (Technical Requirements)

### 1. Service Scaffolding
- Cấu trúc thư mục: `src/apps/warehouse-service`.
- Port: HTTP `8083`, gRPC `50053` (Theo `docs/engineering/port-plan.md`).
- Framework: Go Clean Architecture (Domain, Usecase, Infrastructure).

### 2. Database Schema (Postgres)
- Bảng `processing_batches`: Lưu vết vòng đời (Intake weight, Yield weight, Status).
- Bảng `inventory_items`: Lưu tổng tồn kho khả dụng cho từng SKU.
- Bảng `inventory_logs`: Lưu vết mọi biến động kho (Audit trail).

### 3. State Machine (Domain Layer)
- Định nghĩa các trạng thái: `RECEIVED`, `HULLING`, `ROASTING`, `STOCKED`.
- Hiện thực logic `CanTransitionTo` để chặn các bước nhảy trạng thái sai quy trình.

### 4. Kafka Integration
- Consumer lắng nghe Topic `farm.harvest.created` để tự động khởi tạo Batch ở trạng thái `RECEIVED`.
- Sử dụng **Inbox Pattern** để đảm bảo Idempotency (Không tạo 2 batch cho 1 harvest).

---

## 🧪 Chiến lược Kiểm thử (Verification)
- **Unit Test:** Kiểm tra State Machine với các bộ input hợp lệ và không hợp lệ.
- **Integration Test:** Gửi mock event vào Kafka và kiểm tra DB có tạo bản ghi `processing_batches` không.
- **API Test:** Sử dụng Postman/Curl để test flow cập nhật trọng lượng đầu ra và chuyển sang `STOCKED`.
