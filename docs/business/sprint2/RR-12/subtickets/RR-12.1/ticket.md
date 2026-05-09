# [RR-12.1] Thiết lập nền tảng phân quyền Casbin (Foundation Setup)

- **Mục tiêu:** Xây dựng bộ khung kỹ thuật để hỗ trợ việc kiểm soát quyền truy cập tập trung cho toàn bộ hệ thống.
- **Mô tả:** Thiết lập Casbin Enforcer và các công cụ hỗ trợ để hệ thống có thể hiểu và thực thi các quy tắc phân quyền.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Hệ thống có khả năng nạp cấu hình phân quyền từ file (Model và Policy).
2. Các hành động cơ bản (xem, ghi, xóa) được ánh xạ chính xác từ các yêu cầu kỹ thuật.
3. Đảm bảo nền tảng hoạt động ổn định và sẵn sàng cho việc tích hợp vào các giao thức giao tiếp (gRPC, REST).
