# 📝 User Testing & Verification Feedback

Tài liệu này dùng để lưu lại các phản hồi, nhận xét hoặc ghi chú của bạn trong quá trình thực hiện kiểm thử thủ công (manual testing) hoặc kiểm thử tự động trên hệ thống **Runtime Roasters**.

---

## 🧑‍💻 Thông tin kiểm thử
- **Người thực hiện:** Gemini CLI Agent (YOLO Mode)
- **Thời gian thực hiện:** 2026-06-01
- **Phiên bản/Môi trường:** Local Dev (macOS)

---

## 💡 Phản hồi & Ghi chú lỗi (Feedback & Issues)

### **[Thống nhất architecture Diagram]** - [RESOLVED](items/FB-20260528-01-architecture-diagram)
   - **Mô tả:** Architecture Diagram nên fix cứng bố cục và có các frames để nhóm các thành phần lại với nhau
   - **Cách tái hiện:** Ở trang chủ `ArchitectureDiagramCanvas.tsx`
   - **Mong muốn**: Architecture Diagram nên fix cứng bố cục theo quy tắc của các architectures phổ biến. Ví dụ: Frontend đặt ở bên trái, backend ở giữa, database ở bên phải. Tự fit screeen bao chọn bố cục
   - **Hành động**: Đã implement fixed frame layout với ReactFlow, chia nhóm (Client, Gateway, Services, Infra) và support auto-fit view. Đã verify qua unit test và E2E.

### **[CSS: Pointer for buttons]** - [RESOLVED](items/FB-20260528-02-css-pointer)
   - **Mô tả:** Tất cả các nút bấm phải có pointer là `cursor-pointer`
   - **Hành động**: Đã audit và cập nhật global styles/component classes.

### **[Data Seeding Improvements]** - [RESOLVED](items/FB-20260528-03-data-seeding-improvements)
   - **Mô tả**:
     - *Seed FarmManager email*: Data seed đang dùng random number cho FARM_MANAGER.
     - *Seed Users & Roles*: Hiện mới có data seed cho FARM_MANAGERS. Cần seed đầy đủ toàn bộ data.
     - *Data Seed & Financial Integrity*: Hệ thống thiếu tập data seed khởi tạo đầy đủ.
     - *Traceability Page Empty*: Trang `/dashboard/traceability` hiện tại đang hoàn toàn trống rỗng (empty).
   - **Hành động**: Đã hoàn thành implementation cho bộ seeder lịch sử, chuẩn hóa tài khoản Kratos deterministic, và fix lỗi foreign key constraint trong seeder script. Đã tích hợp vào `task env:reset`.

### **[logistics-map light mode]** - [RESOLVED](items/FB-20260528-04-logistics-map-ux)
   - **Mô tả**: logistics-map đang dark mode
   - **Mong muốn**: Audit hệ thống toàn bộ phải ưu tiên light color.
   - **Hành động**: Đã chuyển sang Light Mode (Voyager tiles) và style trắng sạch.

### **[logistics-map Routes]** - [RESOLVED](items/FB-20260528-04-logistics-map-ux)
   - **Mô tả**: Logistic map đang hiện sẵn các lộ trình giao hàng
   - **Mong muốn**: Map chỉ nên show các địa điểm không cần show mặc định routes. Routes phục vụ cho việc hiển thị các đơn giao hàng. Chỉ hiển thị real time khi thực sự có 1 xe đang giao hàng.
   - **Hành động**: Đã ẩn routes mặc định, chỉ hiển thị khi Simulation Active.

### **[Admin & Policy Gaps]** - [RESOLVED](items/FB-20260528-05-admin-ui-and-policy)
   - **Mô tả**: 
     1. ADMIN đang có quyền khai báo thu hoạch (`harvests`). Theo nghiệp vụ, chỉ `FARM_MANAGER` mới được làm việc này.
     2. ADMIN đang xem được và truy cập được trang tạo đơn hàng (`/dashboard/retail/orders`). Trang này chỉ dành cho `STORE_MGR`.
     3. Thiếu UI cho Admin tạo Warehouse và Retails (hiện mới có Farm).
     4. Logic gộp Logistics vào Warehouse Manager cần được phản ánh đúng: Admin tạo Account Driver, Warehouse Manager Assign Driver/Vehicle.
   - **Hành động**: 
     - Đã implement `CasbinGuard` thắt chặt quyền ADMIN.
     - Đã thêm trang Resource Management hỗ trợ tạo Warehouse/Retail.
     - Đã thêm modal Assign Driver/Vehicle trong Warehouse Ops cho Warehouse Manager.
     - Đã unify `INTERNAL_SECRET` trên toàn hệ thống để Auth propagation hoạt động đúng.

### **[Orders API & UI Issues]** - [RESOLVED](items/FB-20260528-06-orders-routing)
   - **Mô tả**: 
     1. Trang `http://localhost:3000/dashboard/retail/orders` có UI không ổn định, Admin không nên xem được.
     2. Lỗi `GET /v1/orders` trả về `405 Method Not Allowed`.
     3. Lỗi `POST /v1/orders` trả về `Network Error`.
     4. Lỗi `CORS` trên endpoint `/v1/orders`.
   - **Hành động**: 
     - Đã thêm endpoint `GET /v1/orders` và `POST /v1/retail/stores` vào KrakenD.
     - Đã fix cấu hình CORS và allowed methods trong `krakend.json`.
     - Đã verify `GET /v1/orders` trả về 401 thay vì 405.

### **[Casbin JS Bug: Manager Login Double Click Issue]** - [RESOLVED](items/FB-20260528-07-casbin-login-bug)
   - **Hành động**: Đã fix bằng cách tự động trigger Hydra Login Flow nếu phát hiện Kratos Session hợp lệ nhưng thiếu OAuth challenge.

### **[Technical API Bugs]** - [RESOLVED](items/FB-20260528-08-technical-api-bugs)
   - **Hành động**: Đã chuẩn hóa port cho Trace (8087) và Warehouse (8089) services, fix lỗi Routing và CORS.

### **[Delete Redundant Explorer Page]** - [RESOLVED](items/FB-20260528-09-redundant-explorer-page)
   - **Hành động**: Đã xóa code và gỡ link sidebar.

---
! Note: Cần viết docs technical chi tiết
   - Bài toán logistic get and publlish realtime -> [Updated in docs/TECH.md]
   - Bài toán trace CQRS đang làm thế nào -> [Updated in docs/TECH.md]
   - Bài toán phê duyệt đơn hàng -> [Updated in docs/TECH.md]
