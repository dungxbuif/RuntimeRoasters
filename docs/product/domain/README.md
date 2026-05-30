# 🏛️ CHIẾN LƯỢC VẬN HÀNH & HÀNH TRÌNH GIÁ TRỊ HẠT CÀ PHÊ

**Runtime Roasters — Hệ sinh thái Minh bạch hóa Nông sản Kỹ thuật số**

Tài liệu này không chỉ là bản đặc tả, mà là **Bản Hiến Pháp Vận Hành** của Runtime Roasters. Nó kể lại câu chuyện về cách chúng ta bảo vệ giá trị của từng hạt cà phê, từ lúc còn ở trên cành cây cho đến khi nằm trong tay khách hàng cuối cùng. Mọi hành vi của hệ thống, mọi quyết định của người dùng và mọi rào chắn bảo mật đều phải tuân thủ nghiêm ngặt câu chuyện này.

---

## CHƯƠNG 1: KHỞI NGUỒN & TRIẾT LÝ MINH BẠCH

### 1.1. Nỗi đau của sự đứt gãy
Ngành cà phê đặc sản (Specialty Coffee) thường bị mất giá trị do một thực tế phũ phàng: Sự đứt gãy thông tin. Khi hạt cà phê rời khỏi nông trại, nó thường trở thành một con số vô danh trong chuỗi cung ứng. Nông dân bị ép giá vì không chứng minh được công sức canh tác, xưởng rang khó kiểm soát chất lượng đầu vào, và khách hàng cuối cùng phải trả tiền cho một "lời hứa" về nguồn gốc thay vì một "thực chứng" dữ liệu.

### 1.2. Sứ mệnh "Thực chứng Dữ liệu"
Runtime Roasters sinh ra để thay thế sự "tin tưởng cảm tính" bằng **"Minh bạch hóa 100% bằng dữ liệu thực"**. Chúng ta không chỉ ghi lại sổ sách; chúng ta số hóa sự di chuyển vật lý của hàng hóa. Một bao cà phê không thể tự mình "bay" vào kho; nó phải được vận chuyển, giám sát GPS và xác nhận bởi những con người thực thụ tại các điểm chạm.

---

## CHƯƠNG 2: HỆ SINH THÁI VẬN HÀNH ĐA LỚP

Hệ thống được cấu thành bởi một mạng lưới các nút thắt giá trị, mỗi nút có một ranh giới và trách nhiệm riêng biệt.

### 2.1. Các Nút thắt Giá trị (Nodes)
- **Nông trại (Farm):** Điểm khởi đầu của dòng chảy dữ liệu. Mỗi farm được định danh bằng tọa độ GPS và mã vùng (Origin Code).
- **Kho bãi & Xưởng rang (Warehouse/Roastery):** Trung tâm điều phối và biến đổi giá trị. Đây là nơi tiếp nhận hàng thô, thực hiện rang xay và đóng gói thành phẩm.
- **Cửa hàng bán lẻ (Retail Store):** Điểm chạm cuối cùng, nơi nhu cầu hàng hóa được kích hoạt và truyền ngược lại chuỗi cung ứng.

### 2.2. Đội xe Hậu cần (The Fleet)
Để đảm bảo tính vật lý, hệ thống quản lý một đội ngũ **Tài xế (Drivers)** và **Xe tải (Vehicles)** chuyên dụng.
- **Trạng thái Xe:** Một chiếc xe có vòng đời vận hành rõ ràng: *Sẵn sàng (Idle) ➔ Đang đi lấy hàng (En Route to Pickup) ➔ Đang bốc hàng (Loading) ➔ Đang đi giao hàng (En Route to Dropoff) ➔ Đang quay về căn cứ (Returning to Base).*
- **Kỷ luật Vận tải:** Một xe chỉ được gán cho một nhiệm vụ duy nhất tại một thời điểm để đảm bảo tính chính xác của dữ liệu GPS và bằng chứng bàn giao.

### 2.3. Rào chắn Scoping (Fail-closed)
Mọi nhân sự quản lý (Manager) đều vận hành trong một "không gian làm việc" (Scope) hữu hạn. Nếu bạn là Quản lý Farm A, bạn sẽ không thấy dữ liệu của Farm B. Nếu bạn chưa được ADMIN gán vào bất kỳ cơ sở nào, hệ thống sẽ ở trạng thái **Fail-closed** (Bị khóa hoàn toàn) để bảo vệ sự riêng tư và bảo mật dữ liệu.

### 2.4. Phân định Vai trò và Quyền hạn (Strict Role Mandates)
Hệ thống tuân thủ một bộ quy tắc quản trị phân quyền khắt khe để đảm bảo đúng người, đúng việc:
- **System Administrator (ADMIN)**: Đóng vai trò là người thiết lập hệ thống (Global Orchestrator). 
  - *Được phép*: Tạo mới các thực thể nền tảng (Accounts, Roles, Warehouses, Farms, Retail Stores) và thực hiện hành động "Gán" (Assign) các Manager tương ứng vào các thực thể này. Tạo tài khoản Driver.
  - *Không được phép*: Tham gia vào bất kỳ nghiệp vụ tác nghiệp nào (không tạo đơn hàng, không khai báo thu hoạch).
- **Farm Manager (FARM_MANAGER)**:
  - *Được phép*: Có toàn quyền thực hiện các APIs, truy cập resources và thực thi các nghiệp vụ liên quan đến nông trại được gán (ví dụ: Khai báo thu hoạch).
- **Warehouse Manager (WAREHOUSE_MGR)**:
  - *Được phép*: Quản lý toàn bộ tài nguyên và nghiệp vụ của kho bãi. Do yêu cầu rút gọn hệ thống, logic quản lý Logistics được gộp chung cho Warehouse Manager. Mặc dù lý thuyết cần có CRUD cho Xe và Tài xế, nhưng ở đây UI CRUD này được lược bỏ (thay bằng Seed Data). Warehouse Manager chỉ có quyền **điều phối/gán** Tài xế và Xe tải vào các chuyến hàng.
- **Store Manager (STORE_MGR)**:
  - *Được phép*: Chỉ tập trung quản lý các nghiệp vụ bán lẻ của cửa hàng được gán (ví dụ: Tạo đơn hàng).

---

## CHƯƠNG 3: HÀNH TRÌNH GIỮ GÌN GIÁ TRỊ (PRODUCT JOURNEY)

Toàn bộ vận hành của Runtime Roasters tuân thủ 3 chặng hành trình cốt lõi, không có lối tắt, không có ngoại lệ.

### 3.1. Chặng 1: Từ Cây cà phê đến cổng Kho (Inbound Journey)
Hệ thống thiết lập một kỷ luật sắt đá: **Cấm nhập kho ảo.**
1. **Khai báo thu hoạch:** Khi cà phê hái xong, Quản lý Farm tạo một Lô thô (`RR-H`). Hệ thống yêu cầu xác thực vị trí GPS tại farm để chống khai báo khống sản lượng.
2. **Triệu hồi Vận tải:** Ngay lập tức, một yêu cầu lấy hàng (Pickup Request) xuất hiện trên bảng điều khiển của Quản lý Kho. Quản lý Kho gán một Tài xế và Xe trống.
3. **Mô phỏng Di chuyển:** Tài xế bắt đầu chạy. Trong suốt hành trình, ứng dụng của tài xế phát tín hiệu GPS liên tục (mỗi 5-10 giây). Bất kỳ ai có quyền giám sát đều có thể theo dõi chiếc xe di chuyển trên bản đồ thời gian thực.
4. **Xác nhận bàn giao:** Tài xế xác nhận "Đã bốc hàng" tại Farm và "Đã về tới cổng Kho". Hệ thống ghi nhận các mốc thời gian này vào lịch sử vĩnh viễn.
5. **Chặng về bắt buộc (Return to Base):** Tài xế **phải** xác nhận đã quay về căn cứ an toàn mới được phép kết thúc nhiệm vụ. Lúc này, Quản lý Kho mới chính thức tạo được Phiếu nhập kho (Intake).

### 3.2. Chặng 2: Biến đổi tại Xưởng rang (Processing Journey)
Đây là giai đoạn "Rang xay" - nơi hạt cà phê tươi trở thành cà phê đặc sản.
1. **Kỷ luật Hao hụt:** Mỗi mẻ rang đều có tỷ lệ hao hụt khối lượng tiêu chuẩn (12% - 20%). 
2. **Bằng chứng Minh bạch:** Nếu kết quả rang xay sai lệch quá **5%** so với dự kiến, hệ thống **cưỡng bức** người vận hành phải nhập "Ghi chú bất thường" (Anomaly Note). Ghi chú này sẽ đi kèm theo bao cà phê đến tận tay khách hàng cuối cùng để giải thích lý do (ví dụ: máy lỗi, hạt cháy).
3. **Bảo toàn "Gen dữ liệu":** Mọi mã lô thành phẩm (`RR-S`) đều được "đúc" từ mã lô thô (`RR-H`), đảm bảo dòng máu dữ liệu chảy xuyên suốt không bị ngắt quãng.

### 3.3. Chặng 3: Từ Xưởng đến Tách cà phê (Outbound Journey)
Hệ thống sử dụng kịch bản phối hợp đa dịch vụ để đảm bảo "Chỉ bán thứ chúng ta thực sự có".
1. **Thanh toán & Khóa kho:** Khi Cửa hàng đặt hàng, hệ thống thực hiện thanh toán (mô phỏng). Ngay sau đó, một cơ chế **Khóa kỹ thuật số** sẽ tạm khóa số lượng hàng trong kho để tránh tình trạng hai cửa hàng cùng tranh nhau một món hàng cuối cùng (Overselling).
2. **Hoàn tiền tự động:** Nếu vì một sự cố mạng mà kho hết hàng ngay khi thanh toán xong, hệ thống tự động kích hoạt lệnh **Hoàn tiền (Refund)** và báo lỗi cho người dùng.
3. **Giao hàng Realtime:** Tài xế thực hiện giao hàng tới Cửa hàng với quy trình giám sát GPS và xác nhận bàn giao tương tự chặng 1.

---

## CHƯƠNG 4: KỶ LUẬT DỮ LIỆU & QUY CHUẨN TRUY XUẤT

Để duy trì vị thế là "Nguồn Chân Lý", Runtime Roasters áp dụng những quy chuẩn dữ liệu cực kỳ khắt khe.

### 4.1. Định danh Toàn cầu (Universal IDs)
Mọi đối tượng đều mang một mã định danh mang tính kể chuyện:
`RR-{Loại}-{Vùng}-{Ngày}-{Số Thứ Tự}`
(Ví dụ: `RR-H-CD-20260525-0001` - Lô thu hoạch tại Cầu Đất, ngày 25/05/2026, thứ tự 0001).

### 4.2. Chống trùng lặp & Bảo vệ giao dịch (Idempotency)
- **Cửa sổ 30 giây:** Hệ thống thiết lập một "khoảng lặng" 30 giây cho mọi thao tác quan trọng (Harvest/Order). Nếu người dùng ấn nút nhiều lần hoặc mạng chập chờn gửi trùng yêu cầu, hệ thống sẽ từ chối để bảo vệ tính nhất quán của số liệu.
- **Mã định danh giao dịch:** Mỗi lệnh từ UI đều có một "chìa khóa" riêng, đảm bảo backend chỉ thực hiện hành động đó đúng một lần duy nhất.

### 4.3. Phân định Trách nhiệm Lưu trữ
- **Dữ liệu Vận hành:** Lưu trạng thái "sống" của hệ thống ngay lúc này (Ai đang ở đâu, hàng còn bao nhiêu).
- **Nhật ký bất biến (Audit Trail):** Mọi nút bấm, mọi mẻ rang hỏng đều được ghi vào một kho lưu trữ đặc biệt, không thể sửa xóa hay giả mạo. Đây là bằng chứng thép khi có tranh chấp.
- **Truy xuất Siêu tốc:** Hệ thống tìm kiếm chuyên dụng đảm bảo khách hàng quét QR có thể thấy lại hành trình 5 năm của bao cà phê trong chưa đầy **3 giây**.

### 4.4. Quy chuẩn Báo cáo Công khai (Data Sanitization)
Vì niềm tin của khách hàng nhưng vẫn bảo vệ bí mật kinh doanh, dữ liệu trả về khi quét mã QR sẽ được lọc sạch:
- **Công khai:** Bản đồ di chuyển, thời gian thu hoạch/rang xay, tên nông trại, các chỉ số chất lượng.
- **Bảo mật:** Doanh thu, số điện thoại cá nhân tài xế, dữ liệu thanh toán nội bộ.

---
*Runtime Roasters — Vận hành bằng Kỷ luật, Chinh phục bằng Minh bạch.*
