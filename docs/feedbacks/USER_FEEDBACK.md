# 📝 User Testing & Verification Feedback

Tài liệu này dùng để lưu lại các phản hồi, nhận xét hoặc ghi chú của bạn trong quá trình thực hiện kiểm thử thủ công (manual testing) hoặc kiểm thử tự động trên hệ thống **Runtime Roasters**.

---

## 🧑‍💻 Thông tin kiểm thử
- **Người thực hiện:** 
- **Thời gian thực hiện:** 
- **Phiên bản/Môi trường:** Local Dev (macOS)

---

## 📊 Kịch bản & Tiến trình kiểm thử thủ công

| STT | Kịch bản kiểm thử (Test Scenario) | Kết quả mong đợi | Trạng thái (OK / Lỗi / Ghi chú) |
|---|---|---|---|
| 1 | **Reset dữ liệu & docker compose** <br>`./deployments/reset-demo-state.sh` | Hệ thống tự khởi động container và dọn dẹp sạch db/topics. | |
| 2 | **Khởi chạy hệ thống** <br>`task all` | Toàn bộ các dịch vụ và web frontend khởi chạy thành công. | |
| 3 | **Đăng nhập OIDC (Ory Kratos)** | Đăng nhập tài khoản `admin@runtimeroasters.com` thành công. | |
| 4 | **Đặt hàng & Luồng SAGA** | Tạo đơn hàng mới ở Retail dashboard kích hoạt luồng Kafka, cập nhật trạng thái kho và giao hàng hoàn tất. | |
| 5 | **Đồng bộ hóa Casbin & Phân quyền** | Các quyền hạn thực hiện đúng vai trò được gán, không bị lỗi bypass. | |

---

## 💡 Phản hồi & Ghi chú lỗi (Feedback & Issues)

### **[Thống nhất architecture Diagram]**
   - **Mô tả:** Architecture Diagram nên fix cứng bố cục và có các frames để nhóm các thành phần lại với nhau
   - **Cách tái hiện:**:L Ở trang chủ `ArchitectureDiagramCanvas.tsx`
   - **Mong muốn**: Architecture Diagram nên fix cứng bố cục theo quy tắc của các architectures phổ biến. Ví dụ: Frontend đặt ở bên trái, backend ở giữa, database ở bên phải. Tự fit screeen bao chọn bố cục

### **[CSS: Pointer for buttons]**
   - **Mô tả:** Tất cả các nút bấm phải có pointer là `cursor-pointer`

### **[Seed Farmanager mail to create user for FarmManager]**
   - **Mô tả**: Data seed dand dùng random number cho FARM_MANAGER `mgr.1779846686075@runtimeroasters.com`, `mgr.1779846871243@runtimeroasters.com`, `mgr.1779847137516@runtimeroasters.com`, `mgr.1779847207810@runtimeroasters.com`, `mgr.1779847323031@runtimeroasters.com`
   - **Cách tái hiện**: Đang nhập hệ thống bằng Account Admin lần đầu. Initializae Data. View: `http://localhost:3000/dashboard/users`
   - **Mong muốn**: Dùng chính tên rút gọn của farm mà managers được assign để làm email

### **[Seed Users]**
   - **Mô tả**: Hiện mới có data seed cho FARM_MANAGERS
   - **Cách tái hiện**: Đang nhập hệ thống bằng Account Admin lần đầu. Initializae Data. View: `http://localhost:3000/dashboard/users`
   - **Mong muốn**: Seed đầy đủ toàn bộ data trong hệ thống => Sẵn sàng để demo toàn hệ thống theo 1 cách hợp lý. Cần Plan với PO và người review trước khi làm

### **[logistics-map light mode]**
   - **Mô tả**: logistics-map đang dark mode
   - **Mong muốn**: Audit hệ thống toàn bộ phải ưu tiên light color. Chưa cần apply các mode

### **[logistics-map Routes]**
   - **Mô tả**: Logistic map đang hiện sẵn các lộ trình giao hàng
   - **Mong muốn**: Map chỉ nên show các địa điểm không cần show mặc định routes. Routes phục vụ cho việc hiển thị các đơn giao hàng. Chỉ hiển thị real time khi thực sự có 1 xe đang giao hàng. => Hiện tại logic đang sai. Cần plan với Reviewer

### **[Data Seed]**
   - **Mô tả**: Khi intitialize. Tôi cần hệ thống có sẵn 1 tập data seed
   - **Mong muốn**:Nhiều data chưa có ví dụ như user các role (Đã đề cập)
     + Ngoài ra Financial Integrity. Cũng chưa có gì

### **[MISSING ADMIN UI]**
   - **Mô tả**: Mới chỉ show số lượng kho bãi, farms. Nhưng chưa thấy list chi tiết, 
   - **Mong muốn**: Admin ngoài thấy các số lượng tổng quan. Còn cần xem được chi tiết từng item và quản lý được các mục đó. Ví dụ: xem chi tiết từng kho bãi, từng lô, từng user, từng đơn hàng..., Danh sách. Có thể có page riêng hoặc tận dụng Map. Cần plan
   
### **[User & Policy Management]**
   - **Mô tả**: Để giữ cho cơ chế phân quyền đơn giản (Keep it simple), khi muốn thay đổi quyền của user thì chỉ cần xóa user đó đi tạo lại. Quyền hạn trong hệ thống đã được định nghĩa sẵn khi gán (assign) hoặc khi user thực hiện hành động sẽ tự động thêm quyền cho các entities tương ứng.
   - **Mong muốn**: Thống nhất luồng nghiệp vụ quản trị phân quyền theo hướng tinh gọn này để tránh phức tạp hóa hệ thống quản lý chính sách động, giúp việc kiểm thử và vận hành trở nên tường minh.

### **[Orders Page Not Found]**
   - **Mô tả**: Khi truy cập đường dẫn `http://localhost:3000/dashboard/orders` báo lỗi trang không tìm thấy (Not Found).
   - **Mong muốn**: Cần sửa lại định tuyến (routing) hoặc tạo trang quản lý đơn hàng cho Admin/Users để hiển thị danh sách đơn hàng.

### **[Traceability Page Empty]**
   - **Mô tả**: Trang `/dashboard/traceability` hiện tại đang hoàn toàn trống rỗng (empty).
   - **Mong muốn**: Cần cung cấp tập dữ liệu mẫu (data seed) hoàn chỉnh cho phần lịch sử nguồn gốc hạt cà phê để Admin có thể xem được biểu đồ truy vết mẫu ngay khi khởi tạo hệ thống.

### **[Business Gap Analysis: Admin Resource Creation & Assignment]**
   - **Mô tả**: Có sự không đồng nhất (Gap) lớn giữa tài liệu thiết kế nghiệp vụ (`docs/product/SPEC.md`) và thực tế triển khai ở giao diện Frontend:
     1. **Theo Business Docs**: `ADMIN` chịu trách nhiệm tạo thủ công tất cả: Nông trại (Farms), Kho bãi (Warehouses), Cửa hàng (Stores/Retail), Đội xe (Vehicles) và Tài xế (Drivers) rồi gán thủ công từng người quản trị.
     2. **Thực tế triển khai**:
        - Giao diện Admin chỉ có UI tạo Nông trại (`FARM`) và gán Manager. Hoàn toàn **thiếu UI** để tạo và gán quản lý cho **Kho bãi (Warehouses)** và **Cửa hàng (Stores/Retail)**.
        - Hệ thống đang dùng cơ chế **Auto-seed đơn giản hóa (Keep it simple)**: Khi một kho bãi/nông trại được tạo, đội tài xế/xe liên quan sẽ tự động được sinh mẫu ngầm để sẵn sàng chạy thử nghiệm thay vì bắt ADMIN tạo thủ công từng dòng dữ liệu xe/tài xế.
   - **Mong muốn**: Cập nhật lại tài liệu `docs/product/SPEC.md` để ghi nhận cơ chế "Tự động seed tài xế/xe khi tạo kho bãi" (Keep it simple), đồng thời lên kế hoạch phát triển các trang UI quản lý và gán Quản lý cho Kho bãi và Cửa hàng.

### **[Missing Users Sidebar Menu Item]**
   - **Mô tả**: Không có tùy chọn (item) trên thanh Menu Sidebar để đi tới trang quản lý người dùng `http://localhost:3000/dashboard/users`. Người dùng hiện tại phải tự gõ thủ công URL trên thanh địa chỉ của trình duyệt để truy cập.
   - **Mong muốn**: Bổ sung liên kết (Link) "Quản lý người dùng" hoặc "Users" vào thanh Sidebar Menu dành cho tài khoản Admin để thuận tiện cho việc truy cập nhanh.

### **[Casbin JS Bug: Manager Login Double Click Issue]**
   - **Mô tả**: Khi đăng nhập bằng tài khoản Manager (ví dụ: `mgr.1779846686075@runtimeroasters.com`), giao diện gặp lỗi Casbin khiến người dùng phải bấm click đăng nhập lần thứ hai mới vào được Dashboard.
   - **Cách tái hiện**: Đăng nhập bằng tài khoản Manager, kiểm tra F12 Console trình duyệt sẽ thấy thông báo lỗi:
     ```text
     [browser] [Casbin] Failed to initialize Casbin policies: TypeError: (0 , __TURBOPACK__imported__module__...newEnforcer) is not a function
         at CasbinProvider.useEffect.initCasbin (src/lib/auth/casbin.tsx:46:38)
     ```
   - **Nguyên nhân kỹ thuật**: Thư viện `casbin.js` trong môi trường Turbopack (Next.js) không export hàm `newEnforcer` theo cách thông thường (`import { newEnforcer } from 'casbin.js'`), dẫn đến hàm bị undefined và gây crash luồng khởi tạo chính sách phân quyền ở lần tải đầu tiên.
   - **Mong muốn**: Cần sửa lại cách import thư viện `casbin.js` tại [casbin.tsx](file:///Users/dungxbuif/workspace/RuntimeRoasters/src/apps/client-app/src/lib/auth/casbin.tsx) (ví dụ: sử dụng default import `import * as casbin` hoặc import tương thích với Next.js/Turbopack) để luồng phân quyền chạy mượt mà ngay từ lần đăng nhập đầu tiên.





   