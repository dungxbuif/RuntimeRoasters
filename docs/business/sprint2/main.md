# Sprint 2: Identity & Access Control

**Trạng thái:** ✅ Hoàn thành (Completed)
**Mục tiêu:** Xây dựng hệ thống định danh (Identity) và phân quyền (Authorization) tập trung, bảo vệ toàn bộ API và cung cấp luồng Đăng nhập/Đăng xuất chuẩn OIDC.

---

## 📋 Trạng thái Ticket (Kanban)

| Ticket | Summary | Status | Role |
| :--- | :--- | :--- | :--- |
| [RR-9](./RR-9/ticket.md) | [Tech] Identity Infrastructure (Ory Kratos & Hydra) | ✅ Done | Security Eng |
| [RR-10](./RR-10/ticket.md) | [BA] Client Auth Flow: Login/Consent UI | ✅ Done | Product Owner |
| [RR-11](./RR-11/ticket.md) | [Tech] JWT Validation & Propagation | ✅ Done | Tech Lead |
| [RR-12](./RR-12/ticket.md) | [Epic] Casbin RBAC Authorization Engine | ✅ Done | Tech Lead |
| [RR-13](./RR-13/ticket.md) | [Tech] E2E Integration: Secure Service-to-Service | ✅ Done | Backend |
| [RR-14](./RR-14/ticket.md) | [Tech] Token Revocation & Valkey Blacklist | ✅ Done | Tech Lead |

---

## 💡 Tầm nhìn Nghiệp vụ (Business Vision)
- **Zero Trust Architecture:** Mọi request vào hệ thống đều phải được định danh và kiểm tra quyền một cách nghiêm ngặt.
- **Seamless Experience:** Người dùng chỉ cần đăng nhập một lần (SSO) để sử dụng tất cả các dịch vụ trong hệ sinh thái.
- **Fine-grained Control:** Phân quyền đến từng hành động cụ thể (VD: Farmer chỉ được tạo Harvest, không được duyệt Payment).

## 📊 Kết quả đạt được (Sprint Result)
- Tích hợp thành công bộ giải pháp Ory (Kratos/Hydra) cho quản lý Identity.
- Triển khai **Three-Gate Security** (Gateway -> JWT -> Casbin).
- Hệ thống hỗ trợ thu hồi Token tức thì qua Valkey.
