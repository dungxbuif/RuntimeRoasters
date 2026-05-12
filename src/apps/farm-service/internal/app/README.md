# Thư mục: App (Nơi lắp ráp Dependency Injection)

## 🏢 Vai trò trong Clean Architecture
Thư mục này không thuộc về một Layer nghiệp vụ cụ thể nào. Nó đóng vai trò là **Main Component** - nơi khởi tạo và kết nối (wire) tất cả các lớp lại với nhau.

## 🎯 Scope (Phạm vi)
- **Nên viết:**
    - Khởi tạo kết nối Database, Redis, Kafka.
    - Khởi tạo các Repository, Usecase, và Handler.
    - Thực hiện "Dependency Injection" thủ công (truyền Repository vào Usecase, truyền Usecase vào Handler).
    - Cấu hình Middleware, gRPC Server.
- **Không nên viết:**
    - Bất kỳ logic nghiệp vụ nào.

## 🚀 Quy tắc vàng
Nếu bạn muốn thêm một Repository mới hoặc một Usecase mới, đây là nơi bạn phải vào để "khai báo" và kết nối chúng.
