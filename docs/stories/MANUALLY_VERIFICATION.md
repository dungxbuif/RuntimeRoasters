# 🧪 Runtime Roasters: Master System Validation Manual (Updated)

Tài liệu này hướng dẫn quy trình xác thực hệ thống từ khi setup mới, đảm bảo tính đúng đắn của Phân quyền (RBAC) và Luồng nghiệp vụ (Saga), bao gồm cả các kịch bản Happy Path và kịch bản Rollback.

---

## 🛠️ Global Accounts (Dữ liệu Seed)

| Vai trò | Email Đăng nhập | Entity Quản lý | Quyền hạn Chính (UI Constraints) |
| :--- | :--- | :--- | :--- |
| **System Admin** | `admin@runtimeroasters.com` | Toàn hệ thống | Overview, Manage Users, Farm/Warehouse/Retail Registry. **KHÔNG** thấy Create Order, Harvest. |
| **Farm Manager** | `mgr.caudat@runtimeroasters.com` | Cầu Đất Farm | Chỉ thấy Cầu Đất Farm. Được tạo Harvest. |
| **Warehouse Mgr**| `mgr.songthan@runtimeroasters.com` | Kho Sóng Thần | Quản lý kho, điều phối xe (Assign Driver), tạo Batch (Roasting). |
| **Store Manager**| `mgr.hcm.q1@runtimeroasters.com` | Store Quận 1 | Chỉ thấy Store Q1. Được tạo Đơn hàng (Create Order). |
| **Driver** | `driver.songthan01@runtimeroasters.com` | Xe 51C-11111 | Chỉ thấy chuyến hàng được giao. Bấm Start Route, Confirm. |

*Mật khẩu chung:* `Hello@123`

---

## 🏁 Giai đoạn 1: Khởi tạo Hệ thống (System Bootstrap & Admin Verification)

### Bước 1: Reset Môi trường
1. Chạy `task env:reset`.
2. Đảm bảo Backend (`task be`) và Frontend (`task fe`) đã khởi động.

### Bước 2: Admin Đăng nhập & Xác thực UI Constraints
1. Truy cập `http://localhost:3000`, đăng nhập bằng `admin@runtimeroasters.com`.
2. **Xác thực UI Admin (RBAC Strict):**
   - Phải thấy menu: **Manage Users**, **Farm Registry**, **Logistics Map**, **Intelligence Hub**.
   - **Xác thực Không hiển thị:** KHÔNG được thấy menu **Harvest Declaration**, **Create Order**. Admin không có quyền vận hành chi tiết.
3. **Thao tác Khởi tạo:** Nếu hiện popup "System Bootstrap", bấm **Initialize Data** để nạp và đồng bộ master data.
   - `auth-service` sẽ điều phối gọi gRPC/REST nội bộ đến các service thông qua mã bảo mật `X-Internal-Secret` và tự động ánh xạ email sang các Kratos UUID động. Đảm bảo không bị lỗi 409 Conflict hoặc 401 Unauthorized.
4. **Xác nhận Kết quả Seed:**
   - Vào **Users**, thấy các Manager & Driver.
   - Vào **Farm Registry / Logistics Map**, thấy Farm/Store/Warehouse đã được Assigned đúng Manager.

---

## 🏁 Giai đoạn 2: Vận hành Thượng nguồn (Farm -> Warehouse) - Happy Path

### Bước 3: Farm Manager Khai báo Thu hoạch (Harvest)
1. Đăng xuất Admin, đăng nhập bằng `mgr.caudat@runtimeroasters.com`.
2. **Xác thực UI:** Phải thấy menu **Harvest Declaration**, CHỈ hiển thị thông tin Farm Cầu Đất. KHÔNG thấy Create Order.
3. **Thao tác:** Tạo một bản ghi thu hoạch mới (VD: 100kg Arabica).
4. **Hành vi hệ thống:** Service Farm sẽ tạo `farm.harvest.created`. Trạng thái harvest sẽ hiển thị `PICKUP_REQUESTED` sau khi kho nhận thông báo.

### Bước 4: Warehouse Manager Điều phối Inbound (Pickup)
1. Đăng xuất, đăng nhập bằng `mgr.songthan@runtimeroasters.com`.
2. **Xác thực UI:** Thấy Dashboard quản lý Kho Sóng Thần, Inbound Queue, Outbound Queue.
3. **Thao tác:** Tìm Pickup Request từ Farm Cầu Đất -> Bấm **Assign Driver** -> Chọn `driver.songthan01`.

### Bước 5: Driver Thực hiện Shipment (Inbound)
1. Đăng xuất, đăng nhập bằng `driver.songthan01@runtimeroasters.com`.
2. **Thao tác:** Vào **My Shipments**. Bấm **Start Route**, chờ xe chạy giả lập trên bản đồ, xác nhận các Milestone (Confirm Pickup, Confirm Warehouse Arrival).
3. **Hành vi:** Sau khi Driver hoàn tất, kho sẽ tự động tạo `Intake` và Intake được cộng vào Inventory.

---

## 🏁 Giai đoạn 3: Chế biến & Hạ nguồn (Warehouse -> Retail) - Happy Path

### Bước 6: Warehouse Tạo Roasting Batch
1. Đăng nhập lại `mgr.songthan@runtimeroasters.com`.
2. **Thao tác:** Từ hàng Intake (Coffee nhân), tạo **Roasting Batch** và Finalize để chuyển sang cà phê rang (Finished Goods).
3. **Xác thực:** Kiểm tra bảng **Inventory** (Kho Sóng Thần) đã có số lượng hàng sẵn sàng bán.

### Bước 7: Store Manager Đặt hàng (SAGA: Order -> Payment -> Warehouse)
1. Đăng xuất, đăng nhập bằng `mgr.hcm.q1@runtimeroasters.com`.
2. **Xác thực UI:** Chỉ thấy menu **Create Order** cho Store Q1. KHÔNG thấy Harvest hay Inventory của Kho.
3. **Thao tác:** Tạo đơn hàng (Retail Order) với số lượng nhỏ hơn số dư tồn kho. Thanh toán bằng Stripe Demo (Nhập thẻ test Pass).
4. **SAGA Happy Path Check:** 
   - Thanh toán thành công -> Phát sinh `payment.intent.created` & `payment.simulated_completed`.
   - Hệ thống tự động `Reserve Stock` trong Kho. Trạng thái đơn hàng đổi thành `RESERVED` / `DISPATCH_REQUESTED`.

### Bước 8: Giao hàng (Outbound Logistics)
1. Đăng nhập lại `mgr.songthan@runtimeroasters.com`.
2. Thấy đơn Dispatch cho Store Q1. Assign `driver.songthan01`.
3. Đăng nhập `driver.songthan01`, Start Route từ Kho -> Store, Confirm Delivery, **bắt buộc Confirm Return to Base**.
4. Đơn hàng chuyển sang `COMPLETED`.

---

## 🏁 Giai đoạn 4: Kiểm tra SAGA Rollback (Xác thực Logic Rủi ro)

### Kịch bản 4A: Rollback do Lỗi Thanh Toán (Payment Failure)
1. Đăng nhập `mgr.hcm.q1@runtimeroasters.com`.
2. Tạo Order, nhưng khi ở màn Stripe Checkout, sử dụng thẻ Test báo Fail (hoặc gọi webhook fail).
3. **SAGA Check:**
   - Hệ thống không gửi request `Reserve Stock` tới kho.
   - Đơn hàng tự động bị `CANCELLED`. Inventory không bị ảnh hưởng.

### Kịch bản 4B: Rollback do Thiếu Tồn Kho (Out of Stock / Race Condition)
1. Cần 2 tab đăng nhập `mgr.hcm.q1@runtimeroasters.com` (Store Q1) và `mgr.hcm.q2@runtimeroasters.com` (Store Q2).
2. Tồn kho hiện có: `100 kg`.
3. Cùng lúc tạo đơn hàng Q1 (`80 kg`) và Q2 (`50 kg`). Giả sử cả 2 đều Pass Payment.
4. **SAGA Check:**
   - Đơn hàng đầu tiên (Q1) sẽ Reserve Stock thành công (`100 - 80 = 20 kg`). Trạng thái -> `RESERVED`.
   - Đơn hàng thứ 2 (Q2) sẽ gặp lỗi Reserve Stock do chỉ còn 20kg. Phát sinh `warehouse.stock.reservation_failed`.
   - Hệ thống (Retail Service) nhận event Failed, tự động gọi API refund (simulated) payment.
   - Đơn hàng Q2 bị `CANCELLED`.

---

## 🏁 Giai đoạn 5: Provenance Proof (Audit Toàn trình)

1. Đăng nhập lại Admin hoặc để Public.
2. Truy cập **Provenance Trace** (`/dashboard/traceability` hoặc Public QR).
3. Nhập mã Order hoặc Batch vừa hoàn thành (ở Bước 8).
4. **Expectation:** Vẽ đúng Sơ đồ Hành trình (Trace):
   - Farm (Harvest) -> Driver (Pickup) -> Warehouse (Roast) -> Order (Reservation) -> Driver (Delivery).
   - Trace ID phải liền mạch (Correlation thông qua Kafka Outbox).
