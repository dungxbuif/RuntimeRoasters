# [RR-10] [BA] Client Auth Flow: Login/Consent UI

**User Story:**
Dưới vai trò là **Người dùng hệ thống**, tôi muốn có một giao diện đăng nhập hiện đại và an toàn để tôi có thể truy cập vào các tính năng quản lý của RuntimeRoasters.

**Business Context:**
Giao diện đăng nhập là điểm tiếp xúc đầu tiên của người dùng với hệ thống. Nó cần đảm bảo tính chuyên nghiệp, an toàn và hỗ trợ đầy đủ các tính năng như Quên mật khẩu, Đăng ký và Quản lý phiên làm việc.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Truy cập:** Người dùng truy cập vào Client App.
2. **Chuyển hướng:** Hệ thống nhận thấy chưa đăng nhập, chuyển hướng sang Login Page (Ory Kratos).
3. **Xác thực:** Người dùng nhập Email/Password.
4. **Cấp quyền (Consent):** Nếu ứng dụng yêu cầu quyền mới, hiển thị màn hình Consent (Ory Hydra).
5. **Hoàn tất:** Chuyển hướng về Client App với Access Token hợp lệ.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Password Strength:** Mật khẩu tối thiểu 8 ký tự, bao gồm chữ hoa, chữ thường và số.
- **Session Timeout:** Phiên làm việc kéo dài 24h. Sau 24h yêu cầu đăng nhập lại.
- **Multi-device:** Hỗ trợ đăng nhập trên nhiều thiết bị đồng thời.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Đăng nhập thành công
- **Given:** Tôi có tài khoản hợp lệ `admin@runtimeroasters.com`.
- **When:** Tôi nhập đúng thông tin tại màn hình Login.
- **Then:** Hệ thống hiển thị thông báo "Chào mừng quay trở lại" và chuyển tôi vào Dashboard.

### Scenario 2: Đăng nhập thất bại (Sai pass)
- **When:** Tôi nhập sai mật khẩu 3 lần.
- **Then:** Hệ thống hiển thị cảnh báo và yêu cầu chờ 30s trước khi thử lại (Security policy).

---

## 📊 Yêu cầu Dữ liệu
- `Email` (Required)
- `Password` (Required)
- `Remember Me` (Boolean)
