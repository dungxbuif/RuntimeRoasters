# 🕵️ Tech Lead Code Review: Farm Service (End of Sprint 3)

**Tác giả:** Tech Lead
**Trạng thái:** ❌ BLOCKING VIOLATIONS
**Đối tượng:** `src/apps/farm-service`

---

## 1. Kết luận Tổng quan (Verdict)

Mã nguồn hiện tại của **Farm Service** đang vi phạm nghiêm trọng các tiêu chuẩn kiến trúc (Clean Architecture) và tiêu chuẩn an ninh (Two-Gate AuthZ) đã đề ra. Hệ thống không ở trạng thái "Bulletproof" và không đủ điều kiện để bước vào Sprint 4.

**Điểm đánh giá kiến trúc:** 3/10

---

## 2. Các vi phạm Nghiêm trọng (Critical Findings)

### 🚨 Vi phạm 1: Dependency Rule (Nguyên tắc Phụ thuộc)
- **Vị trí:** `src/apps/farm-service/internal/usecase/farm_usecase.go`
- **Chi tiết:** Lớp `usecase` đang thực hiện `import` trực tiếp gói `repository` từ lớp `infrastructure`.
- **Hậu quả:** Phá vỡ tính đóng gói của kiến trúc Clean. UseCase bị phụ thuộc vào tầng thực thi hạ tầng. Đây là lỗi **BLOCKING** không thể chấp nhận.

### 🚨 Vi phạm 2: Interface Ownership (Quyền sở hữu Interface)
- **Vị trí:** `src/apps/farm-service/internal/infrastructure/repository/farm_repository.go`
- **Chi tiết:** Interface `FarmRepository` đang được định nghĩa tại tầng hạ tầng (Infrastructure).
- **Quy chuẩn:** Theo **ADR 0001**, Interface phải thuộc về người tiêu thụ nó (Consumer-owned). `FarmRepository` **phải** được định nghĩa tại tầng `usecase`.

### 🚨 Vi phạm 3: Security Gate 2 Violation (An ninh Tầng 2)
- **Vị trí:** `src/apps/farm-service/internal/infrastructure/repository/farm_repository.go`
- **Chi tiết:** Các hàm `GetByID`, `List`, `Update`, `Delete` đang sử dụng logic thủ công:
  ```go
  if userId.Role != "ADMIN" && userId.Role != "FARM_MANAGER" {
      query = query.Where("owner_id = ?", userId.Subject)
  }
  ```
- **Lỗi:**
    1. **Lọt concerns:** Tầng Repository không được phép biết về "Role" (đó là việc của Gate 1 - Casbin).
    2. **Cơ chế thủ công:** Việc kiểm tra ownership đang bị rải rác và lặp lại.
    3. **Casbin Level 2:** Hoàn toàn thiếu việc áp dụng Casbin ở tầng Query (Query-level AuthZ) như đã yêu cầu trong **ADR 0003**.

### 🚨 Vi phạm 4: Tech Debt cho Sprint 4 (Transactional Outbox)
- **Chi tiết:** Mặc dù đã kết thúc Sprint 3, Farm Service chưa hề có sự chuẩn bị nào cho **Transactional Outbox**. Không có bảng `outbox_events` trong migration, không có boilerplate cho Relay Worker.
- **Hệ quả:** Dẫn đến việc Sprint 4 sẽ bị quá tải hoặc tech debt bị đẩy về cuối dự án.

---

## 3. Đề xuất Khắc phục (Remediation)

1.  **Refactor Interface:** Di chuyển `FarmRepository` interface sang lớp `usecase`. Đảm bảo `farm_usecase.go` KHÔNG import bất cứ thứ gì từ lớp `infrastructure`.
2.  **Centralized Scoping:** Xây dựng một cơ chế Scoping dùng chung (GORM Scopes) để ép buộc kiểm tra `owner_id` một cách tự động, thay vì dùng `if-else` thủ công.
3.  **Casbin Level 2 Integration:** Implement cơ chế trích xuất Filter từ Casbin Engine và đẩy vào GORM Query.
4.  **Usecase Logic Purity:** Loại bỏ các check Role cứng (ví dụ: `userId.Role == "FARM_MANAGER"`) trong UseCase. Thay thế bằng các Domain logic hoặc Casbin policy.
5.  **Preparation:** Khởi tạo bảng `outbox_events` và cấu hình `Unit of Work (TxManager)` ngay lập tức để phục vụ Sprint 4.

---
**Ký tên**
*Tech Lead (RuntimeRoasters)*
