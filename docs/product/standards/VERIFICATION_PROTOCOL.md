# 🧪 Runtime Roasters: Master System Validation Manual

Tài liệu này là **Cẩm nang Xác thực Toàn diện** (Master Validation Manual) cho hệ thống Runtime Roasters. Tài liệu được thiết kế theo cấu trúc mô-đun (Modular), tuân thủ nghiêm ngặt [Testing & Verification Guidelines](./TESTING_GUIDELINES.md).

---

## 🛠️ Global Configuration & Resources
*Sử dụng các thông tin này làm hằng số (Constants) xuyên suốt mọi giai đoạn xác thực.*

| Resource | Value / Endpoint | Purpose |
| :--- | :--- | :--- |
| **Admin Email** | `admin@runtimeroasters.com` | Tài khoản quản trị tối cao |
| **Default Password**| `Hello@123` | Mật khẩu cho mọi tài khoản Seed/Demo |
| **Frontend URL** | `http://localhost:3000` | Giao diện người dùng |
| **API Gateway** | `http://localhost:8081` | KrakenD Entrypoint |
| **Kafka UI** | `http://localhost:8090` | Giám sát luồng sự kiện (Events) |
| **SigNoz (OTel)** | `http://localhost:3301` | Tracing & Metrics |
| **Kibana (Search)** | `http://localhost:5601` | Kiểm tra CQRS Read Model |
| **Socket Service** | `ws://localhost:8091` | WebSocket realtime topology/dashboard fanout |
| **Trace Topology Public Config** | `http://localhost:8081/v1/traces/public/topology/config` | Public topology pull source |
| **Trace Topology Public History** | `http://localhost:8081/v1/traces/public/topology/history` | Public topology replay/history |

---

## 🏁 Flow 1.0: Fresh System Setup & ADMIN Bootstrap
*Đây là Flow nền tảng, đảm bảo "Luật chơi" và "Thế giới dữ liệu" được khởi tạo chính xác.*

### Bước 1: Khởi động & Dọn sạch môi trường
*   **Hành động:** Chạy lệnh `task env:reset`
*   **Xác thực hạ tầng:**
    *   Truy cập **System Explorer** (`/dashboard/explorer`): Phải thấy Swagger UI và KrakenD Gateway phản hồi 200.
    *   Truy cập **Kafka UI** (`:8090`): Phải thấy các topic `auth.*` và `logistics.*`.
*   **Trạng thái DB:**
    *   `identity_db`: Phải có 1 identity duy nhất (`admin@runtimeroasters.com`).
    *   `farm_db`, `retail_db`: Các bảng nghiệp vụ phải trống.

### Bước 2: Đăng nhập & Kích hoạt Seeding (ADMIN Only)
*   **Hành động:**
    1.  Mở `localhost:3000`, đăng nhập bằng `admin@runtimeroasters.com` / `Hello@123`.
    2.  **UI Requirement:** Sau khi đăng nhập, hệ thống điều hướng về **Manage Users**.
    3.  **Bootstrap Trigger:** Một **System Bootstrap Modal** phải tự động hiện lên nếu dữ liệu chưa được nạp.
    4.  *Ghi chú:* Nếu modal không hiện, có thể truy cập `http://localhost:3000/dashboard?bootstrap=true` để cưỡng bức hiển thị.
*   **Expectation:** Bấm **"Initialize Data"**. Modal báo `Success` và tự đóng sau 2 giây.

### Bước 3: Xác thực Dữ liệu Master (Dữ liệu nền tảng)
Kiểm tra tính hiện diện của dữ liệu tại các trang quản trị:

#### A. Manage Users (Identity List)
*   **UI Requirement:** Phải có bộ lọc **"Filter Role"** và ô **"Search"**.
*   **Xác thực:** Lọc theo Role `Logistics Driver` ➔ Phải thấy 3 tài xế: `driver@...`, `driver.beta@...`, `driver.hcm@...`.
*   **Status:** Tất cả user seed phải có trạng thái **"Active"**.

#### B. Manage Farms (Farm Registry)
*   **Xác thực:** Danh sách phải đủ **6 Farms**.
*   **Data Check:**
    *   `FARM-CAUDAT-001` (K'Ho Coffee Farm)
    *   `FARM-CAUDAT-002` (Cau Dat Arabica)
    *   `FARM-BMT-001` (Aeroco Coffee)
    *   ... (và 3 farm khác).

#### C. Intelligence Hub (Global KPI)
*   **Hành động:** Truy cập menu **Intelligence Hub** (biểu đồ analytics).
*   **UI Requirement:** Phải thấy 4 thẻ KPI chính:
    1.  **Pilot Farms:** Hiển thị số 6.
    2.  **Active Personnel:** Hiển thị tổng số User đã seed (khoảng 10-15).
    3.  **Retail Stores:** Hiển thị số 5.
    4.  **Fleet Status:** Hiển thị "3 Active".
*   **System Integrity:** Các dịch vụ `Identity`, `Kafka`, `Traceability` phải báo trạng thái **"Healthy"** (màu xanh).

---

## 🖥️ UI Visibility Rules (Dành cho ADMIN)
*Giữ giao diện tối giản để tập trung vào quản trị.*

| Màn hình | Item nên hiện | Item nên ẩn |
| :--- | :--- | :--- |
| **Sidebar** | Users, Farms, Account, Intelligence Hub, Trace, Explorer | Các dashboard vận hành lẻ (Harvest, Logistics Monitor) |
| **Intelligence Hub** | KPI Cards, Chain Integrity Status, Anomaly Watch | Các biểu đồ chi tiết của từng lô hàng |
| **Manage Users** | Search, Role Filter, Identity Table, Create Button | Plaintext Passwords |

---

## 🛰️ Monitoring Endpoints (Hậu kiểm kỹ thuật)
*Xác nhận luồng dữ liệu ngầm.*

1.  **Kafka Flow:** Topic `logistics.locations.updated` ➔ Phải có message cho 14 địa điểm.
2.  **Identity Admin:** `http://localhost:4434/admin/identities` ➔ Trả về JSON danh sách toàn bộ User.
3.  **Search Read Model:** `http://localhost:9200/trace_events/_search` ➔ Phải bắt đầu có dữ liệu index.
4.  **Socket Health:** `http://localhost:8091/health/live` và `/health/ready` ➔ Trả về 200 khi Valkey sẵn sàng.
5.  **Topology Config:** `GET /v1/traces/public/topology/config` qua KrakenD ➔ Trả về nodes `service.socket`, `service.trace`, `infra.kafka`, `infra.valkey`.

---

## 🏁 Flow 2.0: Luồng Thượng nguồn (Farm ➔ Warehouse)
*Mục tiêu: Xác thực tính minh bạch và kỷ luật vật lý từ lúc thu hoạch đến khi nhập kho.*

### Bước 1: Khai báo Thu hoạch (Farm Manager)
1.  **Hành động:** Đăng nhập một Manager (ví dụ `mgr.hn.hoankiem@...`).
2.  **Thao tác:** Khai báo 5,000kg cà phê Arabica.
3.  **Xác thực Nghiệp vụ:**
    *   **ID Format:** Mã sinh ra phải có dạng `RR-H-{ORIGIN}-[DATE]-XXXX`.
    *   **GPS Guard:** Hệ thống kiểm tra tọa độ tại farm để chống khai báo sai lệch.
    *   **Outbox:** Kiểm tra `farm_db` ➔ Bảng `outbox_events` phải có record `farm.harvest.created`.

### Bước 2: Điều phối & Vận hành Vật lý (Warehouse Mgr & Driver)
1.  **Dispatch:** Đăng nhập `WAREHOUSE_MGR` ➔ Thấy Pickup Request trạng thái `WAITING` ➔ Gán Driver và Xe trống.
2.  **Simulation:** Đăng nhập `DRIVER` ➔ Bấm **"Start Route"**.
3.  **Xác thực Realtime:**
    *   **GPS Heartbeat:** Theo dõi **Logistics Monitor** ➔ Biểu tượng xe phải di chuyển và cập nhật vị trí mỗi 5-10 giây.
    *   **Offline Buffer:** Thử nghiệm kịch bản rớt mạng ➔ Xác nhận tọa độ được upload bù (Batch upload) khi có mạng lại.
4.  **Confirm Pickup:** Tài xế xác nhận "Đã bốc hàng" tại Farm.

### Bước 3: Chặng về & Nhập kho (Intake)
1.  **Mandatory Return:** Tài xế phải về tới Kho và bấm **"Return to Base"**.
2.  **Intake:** Warehouse Mgr kiểm tra danh sách xe đã về ➔ Bấm **"Create Intake"**.
3.  **Expectation:** Lô hàng chính thức nằm trong kho. Trạng thái Pickup Request chuyển sang `COMPLETED`.

---

## 🏁 Flow 3.0: Luồng Chế biến & Hạ nguồn (Warehouse ➔ Retail)
*Mục tiêu: Xác thực biến đổi giá trị, rào chắn bảo vệ kho và luồng SAGA Đơn hàng.*

### Bước 1: Rang xay & Kỷ luật Hao hụt
1.  **Thao tác:** Warehouse Mgr tạo Production Batch từ lô thô ➔ Thực hiện Rang.
2.  **Test Case (Bất thường):** Nhập khối lượng đầu ra hụt > 5% so với tiêu chuẩn (12-20%).
3.  **Expectation:** Hệ thống **cưỡng bức** yêu cầu nhập **"Anomaly Note"**. 
4.  **Hậu kiểm:** Quét mã lô thành phẩm (`RR-S`) ➔ Phải thấy lý do hao hụt hiển thị minh bạch.

### Bước 2: Đặt hàng & SAGA Fulfillment
1.  **Order:** Store Manager tạo đơn hàng thành phẩm.
2.  **Payment Pending Gate:**
    *   Ngay sau khi tạo đơn, Payment phải ở trạng thái `PENDING`.
    *   Kafka phải có `payment.intent.created`.
    *   Kafka chưa được có `warehouse.stock.reserved` cho đơn này trước khi webhook thành công.
3.  **Stripe Demo Webhook Pass:**
    *   Truy cập Finance dashboard bằng bất kỳ role đã đăng nhập.
    *   Bấm **Pass** trên payment đang `PENDING`.
    *   Client phải gọi `GET /v1/payments/demo/stripe-webhook-key`, ký payload Stripe-compatible bằng `HMAC-SHA256`, rồi gửi `POST /v1/webhooks/stripe`.
    *   **Expectation:** Payment chuyển `SUCCEEDED`, Kafka có `payment.completed`, Warehouse mới bắt đầu reserve stock.
4.  **Stripe Demo Webhook Fail:**
    *   Tạo đơn mới, bấm **Fail** trên payment đang `PENDING`.
    *   **Expectation:** Payment chuyển `FAILED`, Kafka có `payment.failed`, order chuyển `REJECTED/FAILED`, Warehouse **không** tạo stock reservation.
5.  **SAGA Rollback (Edge Case):**
    *   Giả lập kho hết hàng sau khi `payment.completed`.
    *   **Expectation:** Hệ thống phát sự kiện `warehouse.stock.reservation_failed` ➔ Payment refund/compensation chạy ➔ Đơn hàng chuyển sang `FAILED/REJECTED`.

### Bước 3: Giao hàng & Chốt vòng đời (Provenance Proof)
1.  **Delivery:** Driver thực hiện chặng giao tới Cửa hàng ➔ Cập nhật GPS realtime.
2.  **Completion:** Đơn hàng chỉ hoàn tất sau khi Driver xác nhận đã quay về Kho (Mandatory Return).
3.  **The Proof:** Lấy mã Trace ID của gói cà phê tại Store ➔ Sử dụng **Provenance Trace**.
4.  **Expectation:** Thấy bản đồ nối liền mạch toàn bộ dòng đời hạt cà phê.

---

## 🏁 Flow 4.0: Observability & Trace Integrity
*Mục tiêu: Đảm bảo "Sợi chỉ đỏ" dữ liệu (W3C Trace ID) xuyên suốt 100% hệ thống phân tán.*

### Bước 1: Kích hoạt SAGA & Bắt vết
1.  **Hành động:** Thực hiện 1 luồng Order SAGA (Flow 3.0).
2.  **Xác thực kỹ thuật:**
    *   **Trace Context:** Truy cập SigNoz (`:3301`) ➔ Tìm kiếm bằng Trace ID từ log của KrakenD.
    *   **Expectation:** Phải thấy một cây đồ thị (Span tree) nối liền: `Gateway ➔ Retail ➔ Kafka ➔ Payment ➔ Warehouse ➔ Logistics`.
    *   **Storage Check:** Vào Elasticsearch (`:9200`) ➔ Kiểm tra record trong index `trace_events` ➔ Phải chứa metadata `trace_id` khớp hoàn toàn với SigNoz.

---

## 🏁 Flow 5.0: Real-time Awareness (WebSocket Only)
*Mục tiêu: Xác thực tính tức thời và bảo mật của luồng thông báo đẩy. Sprint 6/7 chỉ dùng WebSocket; không dùng SSE.*

### Bước 1: Public Topology Push/Pull
1.  **Pull Config:** Mở `http://localhost:3000/` ở tab ẩn danh.
    *   UI phải render `ArchitectureTopology` từ `GET /v1/traces/public/topology/config`.
    *   Không yêu cầu login.
    *   Flow selector phải có các flow canonical như `Paid Order Fulfillment`, `Delivery Return To Base`, `Socket Push/Pull`.
2.  **Pull History:** Sau khi chạy một order/payment/delivery flow, gọi:
    *   `GET http://localhost:8081/v1/traces/public/topology/history?flow_id=flow.retail.paid-order-fulfillment`
    *   **Expectation:** Response chỉ chứa payload đã sanitize: topic/status/flow/node/edge/source/trace/time/demo-safe short IDs.
3.  **Push Realtime:** Root topology phải mở WebSocket:
    *   `ws://localhost:8091/v1/realtime/public/topology/ws?flow_id=...`
    *   Khi Kafka có event mới, node/edge liên quan phải highlight mà không cần F5.
4.  **Fallback:** Tắt socket-service tạm thời ➔ UI không crash, trạng thái chuyển sang polling và vẫn refresh từ trace history.

### Bước 2: Private Dashboard Stream Scope
1.  **Authenticated Stream:** Mở dashboard khi đã login.
    *   UI/API private dùng `GET /v1/realtime/stream?scope=dashboard&flow_id=...`.
    *   KrakenD route phải yêu cầu JWT/Casbin.
2.  **Role Scope:**
    *   `ADMIN`: nhận được toàn bộ private events.
    *   `STORE_MGR`: chỉ nhận event cùng `store_id`.
    *   `WAREHOUSE_MGR`: chỉ nhận event cùng `warehouse_id`.
    *   `DRIVER`: chỉ nhận event cùng `driver_id` hoặc shipment được assign.
3.  **Security Check:** Mở 2 trình duyệt cho 2 Manager khác nhau ➔ Manager A thực hiện hành động scoped theo store/warehouse ➔ Manager B **không được phép** nhận private realtime event của Manager A.

### Bước 3: Internal Socket Publish API
1.  **Keygen:** Chạy `scripts/socket-keygen.sh warehouse-service`.
    *   Plaintext `SOCKET_API_KEY` chỉ set ở caller service.
    *   Hash/scopes set ở socket-service env `SOCKET_INTERNAL_API_KEYS`.
2.  **Fail Closed:**
    *   Không gửi `X-RR-API-Key` ➔ `401`.
    *   Gửi key sai ➔ `401`.
    *   Gửi key đúng nhưng thiếu `events:publish` ➔ `403`.
3.  **Publish:** Gửi `POST /internal/v1/socket/events` với key hợp lệ.
    *   **Expectation:** socket-service publish Kafka topic `socket.broadcast.requested`.
    *   trace-service nhận event và public/private topology history có entry tương ứng nếu visibility cho phép.

---

## 🏁 Flow 6.0: Public Transparency Show (The QR Trace)
*Mục tiêu: Xác thực trải nghiệm khách hàng cuối và tính an toàn của dữ liệu công khai.*

### Bước 1: Trình diễn Truy xuất (No-Auth)
1.  **Hành động:** Mở trang `localhost:3000/api-docs` (hoặc trang Public Trace) trong tab Ẩn danh (Guest).
2.  **Thao tác:** Nhập một Trace ID hợp lệ (ví dụ mã `RR-S-...`).
3.  **Xác thực "Sạch" (Data Sanitization):**
    *   **Hiển thị:** Phải thấy bản đồ di chuyển, tên nông trại, và mốc thời gian rang xay.
    *   **Ẩn:** Tuyệt đối **không** được thấy Doanh thu, ID tài khoản Kratos, hoặc Tên riêng của Tài xế (Ẩn danh hóa).
4.  **Expectation:** Tốc độ phản hồi trang truy xuất phải < 3 giây.

### Bước 2: Public Topology Không Lộ Dữ liệu Nhạy cảm
1.  **Hành động:** Mở root `/` ở tab ẩn danh trong lúc flow đang chạy.
2.  **Xác thực:** Public topology/history không được hiển thị:
    *   JWT/session/token.
    *   Kratos identity ID.
    *   Raw webhook payload.
    *   Tên riêng tài xế hoặc thông tin tài chính nhạy cảm.
3.  **Expectation:** Chỉ thấy trạng thái hệ thống, topic, node/edge, trace id và các short IDs demo-safe.

---

## 🚀 Final Acceptance Criteria (The "Demo Day" Gate)
Hệ thống chỉ được coi là "Ready for Demo" khi:
1.  Toàn bộ **Flow 1.0 đến 6.0** đạt trạng thái 🟢 PASS.
2.  Mọi mẻ rang sai lệch > 5% đều có bằng chứng **Anomaly Note** trên Public Trace.
3.  Không có hiện tượng **"Nhập kho ảo"**: Mọi lô hàng Intake đều có GPS chặng Pickup và Return đi kèm.
4.  Paid order chỉ reserve stock sau `payment.completed` từ Stripe-compatible webhook Pass.
5.  Webhook Fail phát `payment.failed` và không tạo reservation.
6.  Root `ArchitectureTopology` hoạt động no-auth bằng pull history + WebSocket push, có polling fallback.
7.  Private realtime stream được JWT/Casbin bảo vệ và lọc theo role/entity scope.

---
*Manual Version: 1.5 | Enterprise Handover Ready.*
