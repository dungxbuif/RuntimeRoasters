# [RR-12] [Epic] Casbin RBAC Authorization Engine

**User Story:**
Dưới vai trò là **Security Lead**, tôi muốn triển khai một bộ máy phân quyền linh hoạt và tập trung để đảm bảo mỗi người dùng chỉ có thể thực hiện các hành động trong phạm vi quyền hạn được cấp.

**Business Context:**
Phân quyền là rào chắn cuối cùng để bảo vệ dữ liệu. Một sai sót trong phân quyền có thể dẫn đến việc rò rỉ thông tin nhạy cảm hoặc thao tác trái phép trên dữ liệu tài chính/kho hàng.

---

## 📋 Danh sách Sub-tickets
- **RR-12.1:** Định nghĩa bộ quy tắc (Policy) Casbin.
- **RR-12.2:** Triển khai Enforcer tại tầng Domain.
- **RR-12.3:** gRPC Interceptor để kiểm tra quyền tự động.
- **RR-12.4:** Tích hợp với Auth Service để đồng bộ quyền thời gian thực.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Role-Based:** Quyền được gán cho Vai trò (Role), không gán trực tiếp cho Người dùng (User).
- **Default Deny:** Mọi request mặc định là bị từ chối nếu không có rule cho phép rõ ràng.
- **Service Scoping:** Quyền được phân tách theo từng Service (Farm, Warehouse, Retail).

---

## ✅ Acceptance Criteria (AC)
1. **Successful AuthZ:** User có role `FARM_MANAGER` có thể tạo Harvest.
2. **Access Denied:** User có role `FARM_MANAGER` không thể truy cập API của Retail Service.
3. **Audit Trail:** Mọi hành động từ chối quyền (Access Denied) phải được log lại.
