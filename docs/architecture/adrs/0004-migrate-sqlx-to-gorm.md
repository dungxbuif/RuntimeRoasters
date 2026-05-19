# ADR 0004: Chuyển đổi từ sqlx sang GORM

## Trạng thái
**Accepted**

## Bối cảnh (Context)
Trong Sprint 1 và đầu Sprint 2, dự án sử dụng `sqlx` để tương tác với cơ sở dữ liệu PostgreSQL. `sqlx` mang lại tốc độ cao và cho phép viết Raw SQL linh hoạt.
Tuy nhiên, khi tiến vào Sprint 3 (Farm Service - RR-16) và các tính năng phức tạp hơn, việc quản lý quan hệ (Relationships), tự động hóa các thao tác CRUD cơ bản, và xử lý Transaction bằng raw SQL trở nên lặp đi lặp lại và dễ xảy ra lỗi con người.

## Quyết định (Decision)
Thay thế hoàn toàn `sqlx` bằng **GORM**.
- Cập nhật wrapper trong `pkg/database/postgres.go` để khởi tạo kết nối thông qua `gorm.DB`.
- Sửa lại các Repositories (ví dụ: FarmRepository) để tận dụng cú pháp Chain của GORM (`db.Where().First()`, `db.Create()`).
- Bọc logic quản lý Database Transaction bằng hàm `WithTx` sử dụng `gorm.Transaction()`.

## Hậu quả (Consequences)
- **Tích cực:** Tăng tốc độ phát triển (Developer Velocity) do không phải viết boilerplate SQL cho các thao tác cơ bản. Dễ dàng xử lý các quan hệ dữ liệu phức tạp (Has-Many, Belongs-To) mà không cần JOIN thủ công.
- **Tiêu cực:** Có độ trễ nhất định (overhead) so với Raw SQL do GORM sử dụng reflection. Đội ngũ cần hiểu rõ cơ chế Preload của GORM để tránh lỗi N+1 Query.

## Nguồn tham khảo
- **Sprint:** Cuối Sprint 2 / Đầu Sprint 3.
- **Ticket:** RR-16 (Technical Design: Farm Repository).
- **Thảo luận:** Yêu cầu chuyển đổi ORM từ phía Developer ("Chuyển thành gorm tôi sẽ dùng gorm").
