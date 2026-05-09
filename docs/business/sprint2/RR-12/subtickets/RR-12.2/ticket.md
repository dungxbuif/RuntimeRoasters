# [RR-12.2] Kiểm soát quyền truy cập cho giao thức gRPC

- **Mục tiêu:** Tự động hóa việc kiểm tra quyền cho tất cả các cuộc gọi dịch vụ nội bộ thực hiện qua giao thức gRPC.
- **Mô tả:** Triển khai cơ chế chặn và kiểm tra quyền tại tầng giao tiếp gRPC, đảm bảo mọi yêu cầu đều được xác thực quyền hạn trước khi xử lý.

## ✅ Tiêu chí chấp nhận (Acceptance Criteria)
1. Mọi yêu cầu gRPC đều đi qua bộ lọc kiểm tra quyền (Interceptor).
2. Hệ thống nhận diện được hành động và tài nguyên đang được yêu cầu dựa trên tên phương thức gRPC.
3. Trả về lỗi "Từ chối truy cập" (Permission Denied) nếu người dùng không đủ quyền hạn.
