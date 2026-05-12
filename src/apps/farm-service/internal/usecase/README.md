# Layer: Usecase (Lớp điều phối nghiệp vụ)

## 🏢 Vai trò trong Clean Architecture
Lớp này chứa các logic cụ thể cho từng tính năng của ứng dụng (Application Business Rules). Nó đóng vai trò là "người điều phối" (Orchestrator).

## 🎯 Scope (Phạm vi)
- **Nên viết:**
    - Logic luồng công việc (ví dụ: "Khi tạo Harvest, phải kiểm tra Farm có tồn tại không, sau đó mới lưu").
    - Định nghĩa các **Interfaces** (Repository) để yêu cầu lớp Infrastructure thực thi.
    - Gọi các service bên ngoài thông qua interface.
- **Không nên viết:**
    - Code thực thi SQL cụ thể.
    - Logic phân tích Request gRPC/HTTP (đó là việc của lớp Delivery).

## 🚀 Quy tắc vàng
Lớp Usecase chỉ biết đến lớp **Domain**. Nó không được biết lớp Infrastructure hay Delivery đang dùng công nghệ gì (Postgres hay MySQL, gRPC hay REST).
