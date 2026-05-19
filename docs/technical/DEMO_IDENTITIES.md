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

### Phase 4: Farm -> Warehouse Pickup Demo
- **Màn A**: Đăng nhập `manager.caudat@runtimeroasters.com`, tạo harvest từ Farm dashboard.
- **Màn B**: Đăng nhập Warehouse Manager khi account seed có sẵn, mở Warehouse dashboard để xem pickup request và dispatch driver.
- **Màn C**: Đăng nhập `driver@runtimeroasters.com`, mở Driver/Logistics shipment screen, bấm Start để chạy route simulation, sau đó confirm pickup và return.
- **Màn quan sát**: Mở Logistics map, Trace page, hoặc public `ArchitectureTopology` để xem realtime update.
- **Giá trị**: Chứng minh physical logistics backbone, role-scoped UI, socket/SSE update, và traceability.

### Phase 5: Warehouse -> Retail Delivery Demo
- **Màn A**: Đăng nhập `mgr.hn.hoankiem@runtimeroasters.com` hoặc store manager tương ứng, tạo paid order.
- **Màn B**: Warehouse Manager dispatch retail delivery sau khi stock reserved.
- **Màn C**: `driver@runtimeroasters.com` bấm Start để chạy delivery route và confirm delivery.
- **Màn quan sát**: Retail dashboard xem incoming delivery status; Logistics map và Trace page xem realtime journey.
- **Giá trị**: Chứng minh order SAGA, inventory reservation, logistics delivery, và scoped store visibility.

### Phase 6: Public QR Trace Demo
- **Màn public**: Không đăng nhập, mở public trace showcase page.
- **Thực hiện**: Bấm nút generate/show demo QR codes.
- **Kết quả**: Client hiển thị danh sách QR cho các sản phẩm seed sẵn.
- **Thực hiện**: Click hoặc scan một QR.
- **Kết quả**: Public trace page hiển thị trace thật từ dữ liệu đã chuẩn bị: farm, pickup, warehouse intake, processing, paid order, delivery, driver return.
- **Giá trị**: Chứng minh traceability public mà không expose dữ liệu nhạy cảm.

---

## 🛠️ 3. Lưu ý Kỹ thuật
- **OIDC Flow**: Hệ thống sử dụng Ory Hydra để cấp phát token. Sau khi Login thành công tại Kratos, hệ thống tự động gọi `Accept Login` tới Hydra.
- **JWT Claims**: Role, email, org và store scope của người dùng được nhúng trực tiếp vào JWT (`role`, `email`, `org_id`, `store_ids`) để các service verify in-memory.
- **Seeding**: Tài khoản `admin@runtimeroasters.com` được nạp tự động qua container `rr-seeder` khi khởi động infra.
- **Driver Simulation**: Main demo dùng Driver Client sau khi login bằng `DRIVER`; frontend replay route seed và bắn GPS/status về backend.
- **Realtime**: Private dashboard streams cần auth/role scope. Public root `ArchitectureTopology` có thể no-auth nếu chỉ hiển thị dữ liệu sanitized.
