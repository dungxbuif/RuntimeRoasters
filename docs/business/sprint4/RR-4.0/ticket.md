# Ticket RR-4.0: [Tech] Refactor Farm Service (Sprint 3 Fixes)

**Mục tiêu:** Khắc phục triệt để các vi phạm kiến trúc và bảo mật được phát hiện trong đợt Code Review cuối Sprint 3, đảm bảo nền tảng "Bulletproof" trước khi triển khai tính năng mới.

---

## ✅ Yêu cầu chi tiết

### 1. Refactor Clean Architecture (Dependency Rule)
- Di chuyển interface `FarmRepository` sang tầng `usecase`.
- Đảm bảo `farm_usecase.go` KHÔNG import bất kỳ thứ gì từ tầng `infrastructure`.
- Thực hiện Manual DI tại `internal/app/init.go`.

### 2. Casbin Level 2 Integration (Query-level AuthZ)
- Inject `casbin.Engine` vào Repository.
- Xây dựng GORM Scope `applyAuthZScope(ctx)` để tự động hóa việc kiểm tra quyền sở hữu dữ liệu (`owner_id`) dựa trên chính sách Casbin.
- Loại bỏ hoàn toàn các câu lệnh `if-else` kiểm tra Role thủ công trong Repository và UseCase.

### 3. Chuẩn bị Hạ tầng Outbox
- Tạo migration bổ sung bảng `outbox_events` (chuẩn bị cho ticket RR-4.2).

---

## 🧪 Kiểm thử (Verification)
- **Unit Test:** Đảm bảo UseCase hoạt động với Mock Repository (đã chuyển sang tầng UseCase).
- **Integration Test:** Kiểm tra API Get/List/Update/Delete trả về kết quả chính xác theo quyền của User mà không cần logic role-check thủ công.

---
## 🔗 Tài liệu tham khảo
- [Sprint 3 Code Review](../sprint3/CODE_REVIEW.md)
- [Technical Design (Sprint 4)](../technical_design.md)
