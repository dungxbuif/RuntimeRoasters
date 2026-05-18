 Hệ thống (Infrastructure, Database, Services) đã được dọn dẹp, khởi động lại và seed dữ liệu đầy đủ. Các Service backend
  và frontend hiện đang chạy ngầm ổn định.

  Dưới đây là Kịch bản test thủ công (Manual Test Script) và Data Test cho hệ thống từ Sprint 1 đến Sprint 5 theo đúng các
  luồng nghiệp vụ hiện tại:

  🗄️ Dữ liệu Test (Đã được Seed sẵn)
   1. System Admin (Quyền tối cao)
      - Email: admin@runtimeroasters.com
      - Password: Hello@123
   2. Farm Manager (Quản lý Nông trại)
      - Email: manager.caudat@runtimeroasters.com
      - Password: password123
  ---

  📝 Kịch bản Test Thủ Công

  🟢 Sprint 1 & 2: Public Showcase, Identity & Access Control
  Trọng tâm: Kiểm tra UI cơ bản, chặn truy cập trái phép và luồng Đăng nhập (SSO qua Ory).

   - Test 1.1: Trải nghiệm người dùng vãng lai (Public)
     - Hành động: Mở trình duyệt ẩn danh, truy cập http://localhost:3000/.
     - Kỳ vọng: Xem được trang chủ, Showcase công khai nhưng không thấy sidebar quản trị.
   - Test 1.2: Chặn truy cập trái phép (Security Boundary)
     - Hành động: Gõ trực tiếp URL http://localhost:3000/dashboard/users khi chưa đăng nhập.
     - Kỳ vọng: Hệ thống tự động đẩy (redirect) về trang Login.
   - Test 1.3: Đăng nhập Admin
     - Hành động: Đăng nhập với tài khoản Admin (admin@runtimeroasters.com).
     - Kỳ vọng: Đăng nhập thành công, chuyển hướng vào Dashboard, thấy được toàn bộ các mục trên Sidebar (Quản lý User,
       Farm, v.v.).

  🟢 Sprint 3: Farm CRUD & Phân quyền dữ liệu (Data Isolation)
  Trọng tâm: Test quyền hạn (RBAC/ABAC) giữa Admin và Manager.

   - Test 3.1: Admin cấp phát Farm (Direct Assignment)
     - Hành động: (Đang đăng nhập Admin) Truy cập Farm Registry, bấm tạo Farm mới. Đặt tên "Highland Cau Dat", diện tích bất
       kỳ. Tại dropdown "Owner", chọn gán cho manager.caudat@runtimeroasters.com. Bấm lưu.
     - Kỳ vọng: Tạo Farm thành công. Admin thấy Farm này có trên lưới dữ liệu và Admin có thể click "Xóa" hoặc "Sửa" Farm
       này.
   - Test 3.2: Đăng xuất Admin
     - Hành động: Bấm Đăng xuất (Logout) góc màn hình.
     - Kỳ vọng: Quay lại màn hình Login.
   - Test 3.3: Quản lý nhìn thấy Farm của mình (Self-Ownership & Isolation)
     - Hành động: Đăng nhập bằng tài khoản Manager (manager.caudat@runtimeroasters.com). Truy cập Farm Registry.
     - Kỳ vọng: 
       - CHỈ thấy Farm "Highland Cau Dat" (do Admin vừa giao) hoặc các Farm do chính Manager tạo. Không thấy các Farm của
         người khác.
       - Giao diện KHÔNG hiển thị nút "Delete" đối với Farm (vì Manager không có quyền xóa, chỉ Admin có quyền).
   - Test 3.4: Quản lý tự tạo Farm
     - Hành động: Bấm tạo Farm mới. Đặt tên "Manager's Private Farm".
     - Kỳ vọng: Farm tạo thành công và hệ thống tự động gán owner_id là Manager hiện tại (không cần chọn Owner).

  🟢 Sprint 4: Khai báo mẻ thu hoạch (Smart Harvest) & Gửi Event
  Trọng tâm: Test tính toàn vẹn dữ liệu, chặn số lượng 0, đảm bảo Event gửi đi.

   - Test 4.1: Chặn khai báo lỗi
     - Hành động: Truy cập mục Harvests (Khai báo Thu hoạch). Bấm "Declare New Harvest", chọn một Farm, nhập khối lượng = 0.
     - Kỳ vọng: Hệ thống báo lỗi "Số lượng phải lớn hơn 0" (hoặc nút Record bị vô hiệu hóa), không cho phép lưu.
   - Test 4.2: Ghi nhận mẻ thu hoạch chuẩn
     - Hành động: Nhập khối lượng 150.5 kg, chọn loại hạt, sau đó Record/Save.
     - Kỳ vọng: Mẻ thu hoạch được lưu thành công và xuất hiện trong sổ cái (Ledger Table) trên UI.

  🟢 Sprint 5: Warehouse Intake (Luồng tiếp nhận Kafka - Backend/Monitor)
  Trọng tâm: Đảm bảo Warehouse tự động nhận mẻ thu hoạch (Event-Driven) từ Sprint 4 qua Transactional Outbox & Kafka.

   - Test 5.1: Kiểm tra đồng bộ Kafka
     - Hành động: Sau khi hoàn thành Test 4.2 ở trên, bạn hãy mở terminal và kiểm tra log của warehouse-service (hoặc mở
       Kafka UI tại http://localhost:8090).
     - Kỳ vọng: Log của warehouse-service hiển thị đã nhận được/consume sự kiện farm.harvest.events và tiến hành lưu/chuyển
       trạng thái mẻ hàng thành RECEIVED dưới DB. Nếu UI của Warehouse đã hoàn thiện, bạn có thể kiểm tra danh sách nhập kho
       để thấy mẻ hàng vừa thu hoạch hiển thị chờ xử lý.

  > 💡 Mẹo: Bạn có thể mở UI của Frontend tại: http://localhost:3000/ để bắt đầu quy trình test. Các luồng Auth hiện tại đã
  đồng bộ với hệ thống Ory Kratos/Hydra nội bộ.
