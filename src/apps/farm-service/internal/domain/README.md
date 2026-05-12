# Layer: Domain (Lõi của hệ thống)

## 🏢 Vai trò trong Clean Architecture
Đây là lớp trong cùng nhất của củ hành kiến trúc. Nó chứa các **Entities** (Thực thể) và **Business Logic** (Quy tắc nghiệp vụ) cốt lõi.

## 🎯 Scope (Phạm vi)
- **Nên viết:**
    - Các Struct định nghĩa dữ liệu (ví dụ: `Farm`, `Harvest`).
    - Các phương thức Validation thuần túy (ví dụ: `Validate()`).
    - Các Hằng số (Constants) hoặc Enums liên quan đến nghiệp vụ.
    - Không được phụ thuộc vào bất kỳ thư viện bên ngoài nào (ngoại trừ thư viện Go chuẩn hoặc các package helper của dự án).
- **Không nên viết:**
    - Code liên quan đến Database (SQL, GORM tags - trừ trường hợp đặc biệt).
    - Code liên quan đến giao tiếp (gRPC, JSON).
    - Code gọi đến các service khác.

## 🚀 Quy tắc vàng
Lớp Domain **không được phép import** bất kỳ package nào từ các lớp `usecase`, `infrastructure`, hay `delivery`. Nó là lớp độc lập hoàn toàn.
