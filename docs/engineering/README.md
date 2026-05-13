# RuntimeRoasters Engineering & Developer Documentation

Chào mừng bạn đến với khu vực tài liệu kỹ thuật chi tiết dành cho lập trình viên. Thư mục này chứa các hướng dẫn thực thi, tiêu chuẩn code, sơ đồ database và các tài liệu tham khảo kỹ thuật hàng ngày.

---

# 🗺️ Engineering Master Index

### 1. 📏 Tiêu chuẩn & Quy chuẩn (Standards)
- [Canonical Template Standards](./standards/canonical-template-standards.md) - Quy chuẩn cấu trúc project Go mẫu.
- [Modular Bootstrap Design](./standards/modular-bootstrap-design.md) - Thiết kế bộ khởi tạo modular cho microservices.

### 2. 📖 Tham khảo Kỹ thuật (Reference)
- [Core Framework Package (pkg)](./reference/core-framework-pkg.md) - Hướng dẫn sử dụng các thư viện dùng chung trong `src/pkg/`.
- [Infrastructure Port Map](./reference/infrastructure-port-map.md) - Danh sách Port và cấu hình hạ tầng.

### 3. 🗄️ Cơ sở dữ liệu (Database)
- [Data Models](./database/data-models.md) - Thiết kế các thực thể và quan hệ database tổng thể.
- [Farm Service DB Schema](./database/farm-service-db-schema.md) - Chi tiết schema cho dịch vụ Farm.

### 4. 🎨 Frontend & Client App
- [Frontend Auth Library Design](./frontend/fe_auth_library_design.md) - Thiết kế thư viện xác thực cho React/Next.js.

---

## 🏗️ Nguyên tắc phát triển Kỹ thuật
- **Dry (Don't Repeat Yourself):** Tận dụng tối đa các package trong `src/pkg/`.
- **Type-Safety:** Mọi giao tiếp internal bắt buộc dùng gRPC/Protobuf.
- **Explicit over Implicit:** Ưu tiên code tường minh, Manual DI, không dùng magic.

---
*Cập nhật lần cuối: 2026-05-13*
