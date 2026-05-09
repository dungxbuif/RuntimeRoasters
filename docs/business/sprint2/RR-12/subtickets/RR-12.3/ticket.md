# [RR-12.3] Kiểm soát quyền truy cập cho giao thức REST (Gin)

- **Mục tiêu:** Bảo vệ các điểm cuối API công cộng (REST) bằng cơ chế phân quyền tương đương với gRPC.
- **Mô tả:** Triển khai Middleware cho framework Gin để kiểm soát quyền truy cập dựa trên phương thức HTTP và đường dẫn API.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Các yêu cầu HTTP (REST) được kiểm tra quyền dựa trên phương thức (GET, POST, PUT, DELETE, v.v.).
2. Ánh xạ chính xác giữa đường dẫn API và các quy tắc phân quyền trong hệ thống.
3. Trả về mã lỗi 403 Forbidden nếu yêu cầu không hợp lệ về mặt quyền hạn.
