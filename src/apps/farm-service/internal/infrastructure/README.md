# Layer: Infrastructure (Lớp thực thi kỹ thuật)

## 🏢 Vai trò trong Clean Architecture
Lớp này chứa các chi tiết kỹ thuật cụ thể. Đây là nơi các "hợp đồng" (Interface) từ lớp Usecase được thực thi bằng mã code thực tế.

## 🎯 Scope (Phạm vi)
- **Nên viết:**
    - Hiện thực hóa các Repository (truy vấn Database bằng GORM/SQL).
    - Các client gọi API bên ngoài.
    - Cấu hình Cache, Mailer, Logger...
    - Định nghĩa các "Models" dành riêng cho DB (ví dụ: `FarmModel`).
- **Không nên viết:**
    - Logic nghiệp vụ (Business Logic). Infrastructure chỉ nên làm nhiệm vụ lưu trữ và truy xuất dữ liệu theo yêu cầu.

## 🚀 Quy tắc vàng
Dữ liệu từ Database (Model) nên được chuyển đổi sang dữ liệu Domain trước khi trả về cho lớp Usecase để đảm bảo tính độc lập.
