# [RR-5] Farm Service - Persistence Layer (Repository)

- **Summary:** Cài đặt lớp lưu trữ dữ liệu (PostgreSQL) cho Farm Service.
- **Priority:** `MEDIUM`
- **Description:** Implement cơ chế truy xuất dữ liệu nông trại và mẻ cà phê theo chuẩn Clean Architecture.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Thiết kế Database Schema
- **Given:** Các thực thể `Farm` và `HarvestBatch`.
- **When:** Tôi thực hiện migration SQL.
- **Then:** Các bảng dữ liệu với quan hệ 1-N phải có cấu trúc đúng thiết kế (UUID, Foreign Keys, Index).

### Scenario 2: Lưu trữ Farm mới
- **Given:** Dữ liệu farm hợp lệ.
- **When:** Tôi gọi hàm `repository.CreateFarm(ctx, farm)`.
- **Then:** Bản ghi phải được lưu vĩnh viễn vào bảng `farms` và trả về thông tin đã lưu.

### Scenario 3: Quản lý Transaction cho mẻ thu hoạch
- **Given:** Yêu cầu tạo mới một Batch.
- **When:** Tôi thực hiện lệnh SQL.
- **Then:** Phải sử dụng Database Transaction để đảm bảo tính ACID khi ghi thêm các thông tin liên quan (nếu có).

---

## 🛠️ Technical Notes
- Sử dụng thư viện `jmoiron/sqlx` hoặc GORM.
- Sử dụng UUID v4 cho Primary Keys.

## 📋 Sub-tasks
- [ ] Thiết kế `sql/migrations/`.
- [ ] Implement `internal/repository/postgres/`.
- [ ] Viết Unit Test cơ bản cho Repository.
