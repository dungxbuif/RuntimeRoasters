# [RR-12.1] Auth Service: Bootstrap & Persistence

- **Mục tiêu:** Khởi tạo dịch vụ `auth-service` và thiết lập hệ thống lưu trữ chính sách (Policy) sử dụng Postgres và Casbin Gorm Adapter.
- **Mô tả:** Đây là bước nền tảng để biến `auth-service` thành "Single Writer" duy nhất, quản lý toàn bộ Database phân quyền của hệ thống.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Thư mục `src/apps/auth-service` được khởi tạo với cấu trúc Clean Architecture chuẩn của dự án.
2. Kết nối Postgres thành công, sử dụng `pkg/database/postgres.go`.
3. Tích hợp `github.com/casbin/gorm-adapter/v3` để quản lý bảng `casbin_rule`.
4. Định nghĩa file `model.conf` hỗ trợ RBAC Hierarchy và Wildcard matching.
5. `auth-service` khởi chạy không lỗi và tự động thực hiện migration cho các bảng cần thiết.

## 🛠 Task list cho Developer
- [ ] Khởi tạo project `auth-service` (main.go, config, internal folders).
- [ ] Cấu hình Gorm Adapter v3 kết nối tới Postgres.
- [ ] Nhúng file `model.conf` vào binary bằng `go:embed`.
- [ ] Triển khai hàm khởi tạo Enforcer tại Auth Service để sẵn sàng cho các thao tác Write.
