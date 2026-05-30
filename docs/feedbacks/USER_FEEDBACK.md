# 📝 User Testing & Verification Feedback

Tài liệu này dùng để lưu lại các phản hồi, nhận xét hoặc ghi chú của bạn trong quá trình thực hiện kiểm thử thủ công (manual testing) hoặc kiểm thử tự động trên hệ thống **Runtime Roasters**.

---

## 🧑‍💻 Thông tin kiểm thử
- **Người thực hiện:** 
- **Thời gian thực hiện:** 
- **Phiên bản/Môi trường:** Local Dev (macOS)

---

## 💡 Phản hồi & Ghi chú lỗi (Feedback & Issues)

### **[Thống nhất architecture Diagram]** - [IN_PROGRESS](items/FB-20260528-01-architecture-diagram)
   - **Mô tả:** Architecture Diagram nên fix cứng bố cục và có các frames để nhóm các thành phần lại với nhau
   - **Cách tái hiện:**:L Ở trang chủ `ArchitectureDiagramCanvas.tsx`
   - **Mong muốn**: Architecture Diagram nên fix cứng bố cục theo quy tắc của các architectures phổ biến. Ví dụ: Frontend đặt ở bên trái, backend ở giữa, database ở bên phải. Tự fit screeen bao chọn bố cục

### **[CSS: Pointer for buttons]** - [RESOLVED](items/FB-20260528-02-css-pointer)
   - **Mô tả:** Tất cả các nút bấm phải có pointer là `cursor-pointer`

### **[Data Seeding Improvements]** - [IN_PROGRESS](items/FB-20260528-03-data-seeding-improvements)
   - **Mô tả**:
     - *Seed FarmManager email*: Data seed đang dùng random number cho FARM_MANAGER.
     - *Seed Users & Roles*: Hiện mới có data seed cho FARM_MANAGERS. Cần seed đầy đủ toàn bộ data.
     - *Data Seed & Financial Integrity*: Hệ thống thiếu tập data seed khởi tạo đầy đủ.
     - *Traceability Page Empty*: Trang `/dashboard/traceability` hiện tại đang hoàn toàn trống rỗng (empty).
   - **Mong muốn chung**: Gộp chung thành 1 kế hoạch Data Seeding toàn diện. Xây dựng bộ data mẫu có tính lịch sử (Saga flow) để trang Traceability và Finance có sẵn số liệu demo. Cập nhật lại tài liệu kịch bản demo (GUIDE.md).
     - Seed users email nên lấy tên của các entity mà account đó quản lý hoặc tên rút gọn của roles.
   - **Comment**: Đã hoàn thành implementation cho bộ seeder lịch sử (7 ngày), chuẩn hóa tài khoản Kratos deterministic, và tích hợp vào `task env:reset`. Sẵn sàng để Verify theo `docs/stories/MANUALLY_VERIFICATION.md`.

### **[logistics-map light mode]** - [TODO](items/FB-20260528-04-logistics-map-ux)
   - **Mô tả**: logistics-map đang dark mode
   - **Mong muốn**: Audit hệ thống toàn bộ phải ưu tiên light color. Chưa cần apply các mode

### **[logistics-map Routes]** - [TODO](items/FB-20260528-04-logistics-map-ux)
   - **Mô tả**: Logistic map đang hiện sẵn các lộ trình giao hàng
   - **Mong muốn**: Map chỉ nên show các địa điểm không cần show mặc định routes. Routes phục vụ cho việc hiển thị các đơn giao hàng. Chỉ hiển thị real time khi thực sự có 1 xe đang giao hàng. => Hiện tại logic đang sai. Cần plan với Reviewer

### **[Admin & Policy Gaps]** - [IN_PROGRESS](items/FB-20260528-05-admin-ui-and-policy)
   - **Mô tả**: 
     1. ADMIN đang có quyền khai báo thu hoạch (`harvests`). Theo nghiệp vụ, chỉ `FARM_MANAGER` mới được làm việc này.
     2. ADMIN đang xem được và truy cập được trang tạo đơn hàng (`/dashboard/retail/orders`). Trang này chỉ dành cho `STORE_MGR`.
     3. Thiếu UI cho Admin tạo Warehouse và Retails (hiện mới có Farm).
     4. Logic gộp Logistics vào Warehouse Manager cần được phản ánh đúng: Admin tạo Account Driver, Warehouse Manager Assign Driver/Vehicle (không cần UI CRUD cho Fleet, dùng seed data).
   - **Mong muốn**: Thắt chặt chính sách Casbin và PermissionGuard. Hoàn thiện UI Admin cho các thực thể còn thiếu. Đảm bảo đúng vai trò tác nghiệp.

### **[Orders API & UI Issues]** - [TODO](items/FB-20260528-06-orders-routing)
   - **Mô tả**: 
     1. Trang `http://localhost:3000/dashboard/retail/orders` có UI không ổn định, Admin không nên xem được.
     2. Lỗi `GET /v1/orders` trả về `405 Method Not Allowed`.
     3. Lỗi `POST /v1/orders` trả về `Network Error`.
     4. Lỗi `CORS` trên endpoint `/v1/orders`.
   - **Mong muốn**: Sửa lỗi routing KrakenD (đã fix sơ bộ), kiểm tra kết nối Retail Service, và cải thiện UI trang Order. Chỉ `STORE_MGR` của cửa hàng mới được tạo order.

### **[Casbin JS Bug: Manager Login Double Click Issue]** - [IN_PROGRESS](items/FB-20260528-07-casbin-login-bug)
   - **Mô tả**: Khi đăng nhập bằng tài khoản Manager, giao diện gặp lỗi Casbin khiến người dùng phải bấm click đăng nhập lần thứ hai mới vào được Dashboard.
   - **Nguyên nhân kỹ thuật**: Thư viện `casbin.js` trong môi trường Turbopack (Next.js) không export hàm `newEnforcer` theo cách thông thường.

### **[Technical API Bugs]** - [REOPENED](items/FB-20260528-08-technical-api-bugs)
   - **Mô tả**: Ghi nhận các lỗi 405, Network Error, CORS trên hệ thống Orders API.
   - **Mong muốn**: Đảm bảo Gateway và Microservices thông suốt cho luồng SAGA.

### **[Delete Redundant Explorer Page]** - [RESOLVED](items/FB-20260528-09-redundant-explorer-page)
   - **Mô tả**: Trang `System Explorer` (`/dashboard/explorer`) không còn cần thiết.
   - **Hành động**: Đã xóa code và gỡ link sidebar.

---
! Note: cẦn viết docs technical chi tiết
   - Bài toán logistic get and publlish realtime
   - Bài toán trace CQRS đang làm thế nào 
   - Bài toán phê duyệt đơn hàng
