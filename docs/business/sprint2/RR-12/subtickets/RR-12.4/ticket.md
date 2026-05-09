# [RR-12.4] Định nghĩa và quản lý Chính sách phân quyền

- **Mục tiêu:** Thiết lập các quy tắc phân quyền cụ thể cho các vai trò trong hệ thống (Admin, Farmer, v.v.).
- **Mô tả:** Xây dựng tệp cấu hình chứa các quy tắc phân quyền, xác định rõ "Ai được phép làm gì trên tài nguyên nào".

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Có bộ quy tắc rõ ràng cho các vai trò chủ chốt: admin, farmer, processor, warehouse_mgr, v.v.
2. Hỗ trợ phân cấp vai trò: ví dụ Admin tự động kế thừa quyền của tất cả các vai trò khác.
3. Cấu hình dễ dàng mở rộng và bảo trì mà không cần thay đổi mã nguồn cốt lõi.
