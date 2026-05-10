# ADR 0003: Phân quyền bằng mô hình Two-Gate Hybrid (Casbin + Data Scoping)

## Trạng thái
**Accepted**

## Bối cảnh (Context)
Trong hệ thống Multi-tenant và Microservices, việc phân quyền không chỉ dừng lại ở việc "User này có quyền gọi API này hay không?" (RBAC), mà còn phải kiểm tra "User này có sở hữu dữ liệu này hay không?" (ABAC). Nếu dồn toàn bộ logic này vào code, các UseCase sẽ bị phình to bởi các câu lệnh `if-else`.

## Quyết định (Decision)
Áp dụng mô hình **Two-Gate Hybrid AuthZ**:
1. **Gate 1 (Biên giới/Middleware):** Sử dụng **Casbin** để kiểm tra RBAC (Role-Based Access Control). Chặn các request không hợp lệ ngay tại Middleware của Gin hoặc Interceptor của gRPC. Sử dụng gRPC Method Name làm tên định danh Resource trong file policy để đồng nhất.
2. **Gate 2 (Database Layer):** Áp dụng **Data Scoping** trong tầng Repository. Mọi truy vấn SQL đều bị ép thêm mệnh đề `WHERE owner_id = $1` để bảo vệ dữ liệu ở mức dòng (Row-level security).

Ngoài ra, Casbin policy sẽ được đồng bộ theo kiến trúc **Resilient Sync** (Snapshot qua gRPC + Live Update qua Kafka).

## Hậu quả (Consequences)
- **Tích cực:** Tầng Business Logic (UseCase) hoàn toàn sạch sẽ, không chứa logic phân quyền. Bảo mật nhiều lớp (Defense in depth).
- **Tiêu cực:** Đòi hỏi Developer phải luôn nhớ truyền `ownerID` vào các hàm Repository. Cấu trúc đồng bộ quyền (Casbin Sync) làm tăng độ phức tạp của hạ tầng.

## Nguồn tham khảo
- **Sprint:** Sprint 2.
- **Ticket:** RR-12 (Fine-grained Authorization).
- Xem chi tiết tại: [Resilient AuthZ Sync](../concepts/resilient-authz-sync.md).
