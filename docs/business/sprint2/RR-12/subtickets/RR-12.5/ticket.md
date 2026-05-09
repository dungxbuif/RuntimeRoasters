# [RR-12.5] Tích hợp phân quyền vào Demo Service

- **Mục tiêu:** Áp dụng thực tế cơ chế phân quyền vào một dịch vụ cụ thể để kiểm chứng hiệu quả.
- **Mô tả:** Cấu hình và tích hợp bộ kiểm soát quyền vào Demo Service, đảm bảo dịch vụ này vận hành đúng theo các quy tắc đã định nghĩa.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Demo Service từ chối các yêu cầu từ người dùng không có quyền (ví dụ: Guest không được tạo Farm).
2. Demo Service cho phép các yêu cầu hợp lệ từ người dùng có quyền (ví dụ: Farmer tạo Farm của chính mình).
3. Hệ thống ghi log rõ ràng về các quyết định cho phép hoặc từ chối truy cập.
