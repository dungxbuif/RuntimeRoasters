# Sprint 3: Farm Service & Vertical Slice

**Trạng thái:** ✅ Hoàn thành (Completed)
**Mục tiêu:** Hoàn thiện dịch vụ Quản lý Nông trại (Farm Service) với đầy đủ tính năng nghiệp vụ, áp dụng phong cách Vertical Slice kết hợp với Clean Architecture để đảm bảo tính sẵn sàng cho quy trình thu hoạch.

---

## 📋 Trạng thái Ticket (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-15](./RR-15/ticket.md) | [Tech] Farm Service: Scaffolding & Composition Root | ✅ Done | Tech Lead |
| [RR-16](./RR-16/ticket.md) | [BA] Quản lý Danh mục Nông trại (Create & List) | ✅ Done | Product Owner |
| [RR-17](./RR-17/ticket.md) | [Tech] Farm CRUD: Update & Delete Logic | ✅ Done | Backend |
| [RR-18](./RR-18/ticket.md) | [BA] Quản lý Lô đất & Quy hoạch Vùng trồng | ✅ Done | Farmer |

---

## 💡 Tầm nhìn Nghiệp vụ (Business Vision)
- **Digital Farm Twin:** Mỗi nông trại thực tế phải được phản ánh chính xác trên hệ thống với các thông số về diện tích, độ cao và loại cà phê chủ đạo.
- **Resource Management:** Giúp người quản lý nắm bắt được năng lực sản xuất của từng vùng nguyên liệu (Cầu Đất, Buôn Ma Thuột,...).
- **Foundation for Traceability:** Thông tin nông trại là "gốc" của toàn bộ chuỗi truy xuất nguồn gốc phía sau.

## 📊 Kết quả đạt được (Sprint Result)
- Khởi tạo thành công Farm Service với đầy đủ các tầng Clean Architecture.
- Hoàn thiện bộ API quản lý Nông trại (CRUD).
- Tích hợp phân quyền Casbin: Chỉ Farmer/FarmAdmin mới được quản lý dữ liệu nông trại của họ.
