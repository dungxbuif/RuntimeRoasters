# Runtime Roasters — Demo Identities & Personas

Tài liệu này liệt kê danh sách các tài khoản người dùng mẫu (Demo Users) được chuẩn hóa để phục vụ quá trình kiểm thử (UAT) và trình diễn hệ thống (Showcase).

---

## 📋 1. Danh sách User Demo (Standardized Roles)

Mọi tài khoản dưới đây đều được gán Role viết hoa (UPPERCASE) theo đúng chuẩn hệ thống mới.

| Nhóm chức năng | Tên người dùng | Email (Username) | Password mẫu | Role (Mã hệ thống) | `store_ids` |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Quản trị hệ thống** | System Administrator | `admin@runtimeroasters.com` | `Hello@123` | **`ADMIN`** | `[]` |
| **Quản trị nông nghiệp**| Agri Supply Admin | `agri.admin@runtimeroasters.com` | `Hello@123` | **`FARM_ADMIN`** | `[]` |
| **Quản lý Nông trại** | Manager Sơn La | `manager.sonla@runtimeroasters.com` | `Hello@123` | **`FARM_MANAGER`** | `[]` |
| **Quản lý Nông trại** | Manager Cầu Đất | `manager.caudat@runtimeroasters.com` | `Hello@123` | **`FARM_MANAGER`** | `[]` |
| **Nhà máy chế biến** | Roast Master | `processor@runtimeroasters.com` | `Hello@123` | **`PROCESSOR`** | `[]` |
| **Vận tải (Logistics)** | Driver Alpha | `driver@runtimeroasters.com` | `Hello@123` | **`DRIVER`** | `[]` |
| **Cửa hàng bán lẻ** | Mgr Hanoi Hoan Kiem | `mgr.hn.hoankiem@runtimeroasters.com` | `Hello@123` | **`STORE_MGR`** | `["11111111-1111-1111-1111-111111111101"]` |
| **Cửa hàng bán lẻ** | Mgr HCM Dist 1 | `mgr.hcm.q1@runtimeroasters.com` | `Hello@123` | **`STORE_MGR`** | `["11111111-1111-1111-1111-111111111103"]` |
| **Cửa hàng bán lẻ** | Mgr Danang Hai Chau | `mgr.dn.haichau@runtimeroasters.com` | `Hello@123` | **`STORE_MGR`** | `["11111111-1111-1111-1111-111111111105"]` |


---

## 🎭 2. Kịch bản Demo Gợi ý (The Storyline)

Để thể hiện được toàn bộ sức mạnh của kiến trúc Microservices và Phân quyền, PO có thể thực hiện theo luồng sau:

### Phase 1: Onboarding (Quyền ADMIN)
- **Hành động**: Đăng nhập bằng `admin@...`.
- **Thực hiện**: Vào `/admin/users`, tạo các tài khoản `FARM_ADMIN` và `FARM_MANAGER`.
- **Giá trị**: Chứng minh tính năng **Admin-only Creation** thông qua Auth Proxy.

### Phase 2: Tài nguyên & Gán quyền (Quyền FARM_ADMIN)
- **Hành động**: Đăng nhập bằng `agri.admin@...`.
- **Thực hiện**: Vào `/admin/farms`, tạo mới "Nông trại Cầu Đất" và gán Manager là `manager.caudat@...`.
- **Giá trị**: Chứng minh khả năng điều phối tài nguyên hệ thống.

### Phase 3: Vận hành cục bộ (Quyền FARM_MANAGER)
- **Hành động**: Đăng nhập bằng `manager.caudat@...`.
- **Thực hiện**: Vào **Farm Origin** dashboard.
- **Giá trị**: Chứng minh **Data Scoping (ABAC)** — Manager này tuyệt đối không thấy nông trại của Manager kia.

### Phase 4: Theo dõi mẻ hàng (Sắp tới - Sprint 4+)
- **Hành động**: Đăng nhập bằng các role vận hành (`FARM_MANAGER`, `PROCESSOR`, `DRIVER`).
- **Giá trị**: Chứng minh tính nhất quán dữ liệu xuyên suốt chuỗi cung ứng.

---

## 🛠️ 3. Lưu ý Kỹ thuật
- **OIDC Flow**: Hệ thống sử dụng Ory Hydra để cấp phát token. Sau khi Login thành công tại Kratos, hệ thống tự động gọi `Accept Login` tới Hydra.
- **JWT Claims**: Role, email, org và store scope của người dùng được nhúng trực tiếp vào JWT (`role`, `email`, `org_id`, `store_ids`) để các service verify in-memory.
- **Seeding**: Tài khoản `admin@runtimeroasters.com` được nạp tự động qua container `rr-seeder` khi khởi động infra.
