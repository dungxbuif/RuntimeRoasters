# Layer: Delivery (Lớp giao tiếp/Vận chuyển)

## 🏢 Vai trò trong Clean Architecture
Lớp này là cửa ngõ duy nhất để người dùng hoặc các hệ thống khác tương tác với ứng dụng. Nó xử lý giao thức (gRPC, HTTP, CLI...).

## 🎯 Scope (Phạm vi)
- **Nên viết:**
    - Parse (phân tích) request đầu vào từ Protobuf hoặc JSON.
    - Validate định dạng dữ liệu thô (ví dụ: chuỗi có phải email không).
    - Gọi Usecase để xử lý nghiệp vụ.
    - Format response trả về (Status code, Error message).
- **Không nên viết:**
    - Logic nghiệp vụ phức tạp.
    - Truy vấn database trực tiếp.

## 🚀 Quy tắc vàng
Delivery chỉ là lớp "phiên dịch" giữa thế giới bên ngoài và lớp Usecase bên trong.
