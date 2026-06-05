# Quy Hoạch Dữ Liệu Khởi Tạo Hệ Thống (Master Seed Data Specifications)

Tài liệu này đặc tả toàn bộ quy hoạch dữ liệu seed (Seed Data), cấu trúc phân quyền (RBAC Constraints), ràng buộc nghiệp vụ, kiến trúc cơ sở dữ liệu ERP phân tán, và chiến lược triển khai API-driven seeding cho hệ thống **Runtime Roasters**.

---

# Bản Đồ Kiến Trúc Dữ Liệu ERP & Quan Hệ Thực Thể (ERP Database & Entity Relationships Map)

Hệ thống Runtime Roasters được vận hành dựa trên kiến trúc **Database-per-Service** (Mỗi microservice sở hữu một cơ sở dữ liệu độc lập). Dữ liệu nghiệp vụ thực tế phát sinh liên tục thông qua các kịch bản chuỗi cung ứng (Farm-to-Cup). Dưới đây là phân bổ chi tiết các bảng dữ liệu trong ERP phân tán và mối liên kết logic giữa các dịch vụ.

## 1. Bản Đồ Dữ Liệu Các Dịch Vụ (Operational Database Schemas)

### 1.1. Farm Service (`farm_db`)
Chịu trách nhiệm quản lý nguồn cung cấp nguyên liệu thô đầu vào tại các nông trại đối tác.
* **`farms` (Danh mục Nông trại - Master Data)**:
  * `id` (BIGSERIAL - Primary Key)
  * `name` (Tên nông trại), `location` (Vùng địa lý: CAU_DAT, BUON_MA_THUOT, PLEIKU...)
  * `latitude`, `longitude` (Tọa độ địa lý)
  * `area` (Diện tích gieo trồng), `coffee_type` (Loại hạt: ARABICA, ROBUSTA...)
  * `owner_id` (UUID - Mã định danh Kratos của `FARM_MANAGER`)
* **`harvests` (Mùa vụ thu hoạch - Transactional Data)**:
  * `id` (BIGSERIAL - Primary Key)
  * `farm_id` (Liên kết với `farms.id` - Foreign Key)
  * `quantity` (Sản lượng thu hoạch tính bằng kg)
  * `harvest_date` (Ngày thu hoạch), `status` (Trạng thái mùa vụ: NEW, PROCESSING, COMPLETED)
  * `owner_id` (UUID - Mã định danh Kratos của manager tạo vụ)

### 1.2. Warehouse Service (`warehouse_db`)
Quản lý chuỗi cung ứng trung nguồn: thu gom cà phê tươi, sơ chế, rang xay thành phẩm, và quản lý kho hàng thành phẩm.
* **`pick_up_requests` (Lệnh thu gom cà phê - Transactional Data)**:
  * `id` (Primary Key - UUID/Text)
  * `harvest_id` (ID mùa vụ cần thu gom từ `farm_db`)
  * `farm_id` (Mã nông trại xuất phát), `warehouse_id` (Mã kho đích)
  * `quantity` (Khối lượng thu hoạch thực tế), `coffee_type` (Loại cà phê)
  * `status` (Trạng thái: REQUESTED, DISPATCHED, RECEIVED)
  * `shipment_id` (Mã chuyến vận chuyển liên kết sang `logistics_db`)
* **`production_batches` (Lô sản xuất rang xay - Transactional Data)**:
  * `id` (Primary Key), `batch_id` (Mã lô sản xuất duy nhất)
  * `warehouse_id` (Mã kho thực hiện), `status` (Trạng thái chế biến: DRAFT, PROCESSING, COMPLETED)
  * `total_input_weight`, `total_output_weight` (Khối lượng nguyên liệu trước/sau khi rang)
  * `weight_loss_percent` (Tỷ lệ hao hụt khối lượng khi rang - chỉ số tối ưu hóa chất lượng)
* **`intakes` (Phiếu nhập kho nguyên liệu - Transactional Data)**:
  * `id` (Primary Key)
  * `harvest_id` (ID mùa vụ liên kết), `pickup_id` (Mã lệnh thu gom liên kết)
  * `quantity` (Khối lượng thực nhận), `status` (Trạng thái: UNASSIGNED, PROCESSING, BATCHED)
  * `batch_id` (Liên kết với `production_batches.id` - Foreign Key)
* **`roast_runs` (Mẻ rang chi tiết - Transactional Data)**:
  * `id` (Primary Key)
  * `batch_id` (Liên kết với `production_batches.id` - Foreign Key)
  * `run_number` (Số lượt rang trong lô), `input_weight`, `output_weight` (Cân nặng hạt trước/sau rang)
* **`inventories` (Kho hàng tồn kho thành phẩm - Master/Transactional Data)**:
  * `id` (Primary Key)
  * `sku` (Mã SKU cà phê thành phẩm - Duy nhất)
  * `coffee_type`, `origin_code` (Loại cà phê và mã vùng xuất xứ)
  * `warehouse_id` (Mã kho chứa)
  * `available_quantity` (Số lượng tồn kho sẵn sàng xuất cho các cửa hàng bán lẻ)

### 1.3. Retail Service (`retail_db`)
Quản lý chuỗi cung ứng hạ nguồn: vận hành các cửa hàng bán lẻ và thu nhận đơn hàng từ người tiêu dùng.
* **`stores` (Danh mục Cửa hàng bán lẻ - Master Data)**:
  * `id` (UUID - Primary Key)
  * `name` (Tên cửa hàng), `city` (Thành phố), `address` (Địa chỉ)
  * `manager_id`, `manager_email` (Thông tin quản lý cửa hàng `STORE_MGR`)
* **`orders` (Đơn đặt hàng bán lẻ - Transactional Data)**:
  * `id` (UUID - Primary Key)
  * `store_id` (Mã cửa hàng nhận đơn - Foreign Key)
  * `items` (Danh sách sản phẩm mua dạng JSONB)
  * `total_amount` (Tổng giá trị đơn hàng), `status` (Trạng thái đơn: PENDING, PAID, SHIPPED, DELIVERED, CANCELLED)
  * `idempotency_key` (Khóa chống trùng lặp giao dịch)

### 1.4. Payment Service (`payment_db`)
Quản lý tích hợp thanh toán (Stripe, VNPay) mô phỏng trong chuỗi cung ứng.
* **`payments` (Giao dịch thanh toán - Transactional Data)**:
  * `id` (UUID - Primary Key)
  * `order_id` (Mã đơn hàng liên kết sang `retail_db` - Duy nhất)
  * `store_id` (Mã cửa hàng xảy ra giao dịch)
  * `provider` (Kênh thanh toán: STRIPE, VNPAY)
  * `provider_ref` (Mã tham chiếu giao dịch phía đối tác thanh toán)
  * `amount` (Số tiền giao dịch), `status` (Trạng thái: PENDING, COMPLETED, FAILED, REFUNDED)
* **`webhook_events` (Nhật ký Callbacks - Audit Data)**:
  * `id` (UUID - Primary Key)
  * `provider` (Kênh thanh toán), `event_id` (Mã sự kiện webhook), `status` (Kết quả xử lý)

### 1.5. Logistics Service (`logistics_db`)
Trọng tâm liên kết toàn bộ chuỗi cung ứng: quản lý đội xe trung chuyển hạt thô từ Farm về Warehouse, và giao sản phẩm từ Warehouse đến Stores.
* **`vehicles` (Danh mục Phương tiện - Master Data)**:
  * `id` (Primary Key - VD: `VEHICLE-DEMO-001`)
  * `plate_number` (Biển số xe), `type` (Loại xe: TRUCK, VAN)
  * `capacity_kg` (Tải trọng tối đa), `home_warehouse_id` (Mã kho chủ quản)
  * `status` (Trạng thái vận hành: IDLE, BUSY, MAINTENANCE)
* **`drivers` (Danh mục Tài xế - Master Data)**:
  * `id` (UUID - Primary Key)
  * `user_id` (Liên kết với mã tài khoản Kratos của Driver)
  * `name` (Họ tên), `phone` (Số điện thoại)
  * `vehicle_id` (Liên kết với `vehicles.id` - Foreign Key)
  * `status` (Trạng thái: IDLE, DRIVING, OFF)
* **`locations` (Bản đồ Vị trí Địa lý - Static Master Data)**:
  * `id` (Primary Key - Map với Farm ID, Warehouse ID, Store ID)
  * `name` (Tên địa điểm), `type` (Phân loại: FARM, WAREHOUSE, RETAILER)
  * `lat`, `lng` (Tọa độ GPS phục vụ tính toán lộ trình)
* **`shipments` (Chuyến hàng vận chuyển - Transactional Data)**:
  * `id` (UUID - Primary Key)
  * `type` (Phân loại vận chuyển: HARVEST_PICKUP, RETAIL_DELIVERY)
  * `order_id` / `harvest_id` (Mã đơn hàng hoặc mã mùa vụ được vận chuyển)
  * `driver_id` (Tài xế thực hiện chuyến hàng - Foreign Key)
  * `vehicle_id` (Phương tiện thực hiện chuyến hàng)
  * `origin_location_id` / `destination_location_id` (Địa điểm đi và đến)
  * `status` (Trạng thái hành trình: PENDING, ASSIGNED, DEPARTED, ARRIVED_ORIGIN, LOADED, DEPARTED_ORIGIN, ARRIVED_DESTINATION, DELIVERED)

---

## 2. Bản Đồ Mối Quan Hệ Giữa Các Thực Thể (Entity Relationship Diagram - ERD)

Dưới đây là sơ đồ Mermaid thể hiện mối quan hệ logic xuyên suốt toàn bộ chuỗi cung ứng **Runtime Roasters**. Do hệ thống sử dụng kiến trúc microservices phân tán, các đường liên kết giữa các bảng thuộc các database khác nhau (được ký hiệu bằng nét đứt) sẽ được thực thi thông qua **Event-Driven Saga (Kafka)** hoặc **API Gateway Aggregation**:

```mermaid
erDiagram
    %% FARM SERVICE DATABASE
    FARMS ||--o{ HARVESTS : "has"
    FARMS {
        bigserial id PK
        string name
        enum location
        decimal latitude
        decimal longitude
        decimal area
        enum coffee_type
        uuid owner_id
    }
    HARVESTS {
        bigserial id PK
        bigint farm_id FK
        uuid owner_id
        enum coffee_type
        decimal quantity
        timestamp harvest_date
        enum status
    }

    %% WAREHOUSE SERVICE DATABASE
    PRODUCTION_BATCHES ||--o{ INTAKES : "processes"
    PRODUCTION_BATCHES ||--o{ ROAST_RUNS : "contains"
    PRODUCTION_BATCHES {
        string id PK
        string batch_id UK
        string warehouse_id
        string status
        decimal total_input_weight
        decimal total_output_weight
        decimal weight_loss_percent
    }
    INTAKES {
        string id PK
        string harvest_id FK "Cross-DB link"
        string pickup_id FK "Cross-DB link"
        string warehouse_id
        enum coffee_type
        decimal quantity
        string status
        string batch_id FK
    }
    ROAST_RUNS {
        string id PK
        string batch_id FK
        int run_number
        decimal input_weight
        decimal output_weight
        string status
    }
    INVENTORIES {
        string id PK
        string sku UK
        enum coffee_type
        string origin_code
        string warehouse_id
        decimal available_quantity
    }
    PICK_UP_REQUESTS {
        string id PK
        string harvest_id UK "Cross-DB link"
        string farm_id
        string warehouse_id
        string status
        string shipment_id FK "Cross-DB link"
    }

    %% RETAIL SERVICE DATABASE
    STORES ||--o{ ORDERS : "receives"
    STORES {
        uuid id PK
        string name
        string city
        string address
        string manager_email
    }
    ORDERS {
        uuid id PK
        uuid store_id FK
        jsonb items
        decimal total_amount
        string status
        string idempotency_key
    }

    %% PAYMENT SERVICE DATABASE
    PAYMENTS {
        uuid id PK
        uuid order_id UK "Cross-DB link"
        uuid store_id
        string provider
        string provider_ref UK
        decimal amount
        string status
    }

    %% LOGISTICS SERVICE DATABASE
    VEHICLES ||--o| DRIVERS : "assigned to"
    DRIVERS ||--o{ SHIPMENTS : "delivers"
    VEHICLES {
        string id PK
        string plate_number UK
        string type
        decimal capacity_kg
        string home_warehouse_id
        string status
    }
    DRIVERS {
        uuid id PK
        string user_id UK "Kratos link"
        string name
        string phone
        string vehicle_id FK
        string status
        boolean is_available
    }
    SHIPMENTS {
        uuid id PK
        string type
        string order_id FK "Cross-DB link"
        string harvest_id FK "Cross-DB link"
        uuid driver_id FK
        string vehicle_id
        string origin_location_id FK
        string destination_location_id FK
        string status
    }
    LOCATIONS {
        string id PK
        string name
        string type
        decimal lat
        decimal lng
    }

    %% CROSS SERVICE LOGICAL RELATIONSHIPS (Kafka / REST)
    HARVESTS ..> PICK_UP_REQUESTS : "Triggers pickup request"
    PICK_UP_REQUESTS ..> SHIPMENTS : "Generates logistics shipment"
    ORDERS ..> PAYMENTS : "Requires payment authorization"
    PAYMENTS ..> SHIPMENTS : "Successful payment triggers retail delivery shipment"
    INVENTORIES ..> ORDERS : "Checks and reserves SKU stock"
```

### Sơ đồ Luồng & Khớp nối Hệ thống Dưới Dạng ASCII (ASCII Architecture & ER Flow)
```text
               +-------------------------------------------------------+
               |                    [ Ory Kratos ]                     |
               |             (Central User Identity DB)                |
               +---------------------------+---------------------------+
                                           | (authz UUIDs)
                                           v
+------------------+     (Kafka Event)     +------------------+     (REST API)     +------------------+
|   [ farm_db ]    | - - - - - - - - - - > | [ warehouse_db ] | < - - - - - - - -  |  [ retail_db ]   |
|                  |                       |                  |                    |                  |
|  +------------+  |  farm.harvest.declared|  +------------+  | inventory.reserved |  +------------+  |
|  |   farms    |  |                       |  | inventories|  |                    |  |   stores   |  |
|  +-----+------+  |                       |  +------------+  |                    |  +-----+------+  |
|        | (1)     |                       |        ^         |                    |        | (1)     |
|        |         |                       |        | updates |                    |        |         |
|        v (N)     |                       |  +-----+------+  |                    |        v (N)     |
|  +------------+  |                       |  | roast_runs |  |                    |  +------------+  |
|  |  harvests  |  |                       |  +-----+------+  |                    |  |   orders   |  |
|  +------------+  |                       |        | (N)     |                    |  +-----+------+  |
+--------|---------+                       |        v (1)     |                    +--------|---------+
         |                                 |  +------------+  |                             |
         | (Triggers Pickup)               |  | batches    |  |                             | (Requires
         v                                 |  +-----+------+  |                             v Payment)
+---------------------------------------+  |        ^ (1)     |                    +------------------+
|           [ logistics_db ]            |  |        |         |                    |  [ payment_db ]  |
|                                       |  |        v (N)     |                    |                  |
|  +------------+      +-------------+  |  |  +------------+  |                    |  +------------+  |
|  | pick_ups   |<---->|  shipments  |  |  |  |  intakes   |  |                    |  |  payments  |  |
|  +------------+      +------+------+  |  +------------------+                    |  +------------+  |
|                             |         |                                          +------------------+
|                             v (uses)  |
|  +------------+      +------+------+  |
|  |  drivers   |<---->|  vehicles   |  |
|  +------------+      +-------------+  |
+---------------------------------------+
```

---

## 3. Phân Tích Luồng Dữ Liệu Nghiệp Vụ ERP (Operational Workflows)

### 3.1. Luồng Cung ứng Thô (Farm-to-Warehouse Flow)
* **Khai báo thu hoạch (Harvesting)**: `FARM_MANAGER` thực hiện tạo bản ghi `harvests` (Trạng thái: `NEW`) liên kết với nông trại `farms` thuộc quyền sở hữu của mình. Sự kiện `farm.harvest.declared` phát ra Kafka.
* **Yêu cầu lấy hàng (Pickup Request)**: `warehouse-service` nghe sự kiện, tự động sinh bản ghi `pick_up_requests` (Trạng thái: `REQUESTED`) để liên kết `harvest_id` đến kho hàng `warehouse_id` gần nhất.
* **Điều phối vận chuyển (Logistics Seeding)**: `WAREHOUSE_MGR` phê duyệt lệnh thu gom, hệ thống gọi sang `logistics-service` để tự động gán `driver_id` và `vehicle_id` rảnh rỗi tương ứng với kho hàng để tạo chuyến xe `shipments` (Phân loại: `HARVEST_PICKUP`).
* **Tiếp nhận & Rang xay (Processing)**: Khi xe cập kho, `warehouse-service` cập nhật trạng thái `pick_up_requests` sang `RECEIVED`, đồng thời tạo phiếu nhập kho `intakes`. Quản lý kho gộp nhiều `intakes` thô vào một lô sản xuất `production_batches` để tiến hành thực hiện các mẻ rang `roast_runs`. 
* **Cập nhật tồn kho (Inventory)**: Khi lô sản xuất hoàn thành, cà phê thô biến thành hạt rang chín. Hệ thống tự động cập nhật số lượng tồn kho `inventories` tương ứng với mã `sku` cụ thể của dòng sản phẩm đó.

### 3.2. Luồng Bán Lẻ & Phân Phối (Retail-to-Delivery Flow)
* **Đặt hàng (Ordering)**: Khách hàng mua sắm tại cửa hàng `stores`. `retail-service` tạo đơn hàng `orders` ở trạng thái `PENDING`.
* **Thanh toán (Payment)**: Hệ thống sinh yêu cầu thanh toán `payments` liên kết với `order_id` sang `payment-service`. Khách hàng thanh toán qua cổng giả lập hoặc Stripe thành công. Sự kiện `payment.completed` phát ra Kafka.
* **Cắt giảm tồn kho (Stock Reservation)**: `warehouse-service` tiêu thụ sự kiện thanh toán, trừ trực tiếp số lượng tồn kho ở `inventories` tương ứng với mã SKU sản phẩm bán ra.
* **Giao hàng (Logistics Dispatch)**: `logistics-service` tiếp nhận thông tin đơn hàng đã thanh toán, `WAREHOUSE_MGR` thực hiện điều phối chuyến vận chuyển `shipments` (Phân loại: `RETAIL_DELIVERY`) giao hạt rang từ Kho trung chuyển (`warehouse_id`) tới địa chỉ cửa hàng bán lẻ (`destination_store_id`). Tài xế (`driver.user_id`) tiếp nhận lộ trình thông qua thiết bị di động và hoàn thành đơn hàng.

---

## 4. Danh sách Vai Trò & Phân Quyền (Role Definitions & Authorization)

Hệ thống Runtime Roasters áp dụng mô hình phân quyền **Centralized Management, Distributed Enforcement** (Quản lý tập trung qua Casbin, Thực thi phi tập trung tại từng microservice). Dưới đây là định nghĩa vai trò của 5 nhóm tài khoản:


* **`ADMIN` (System Administrator)**
  * **Phạm vi**: Quyền xem và tạo các thông tin danh mục gốc (Master Data). Đối với dữ liệu phát sinh từ nghiệp vụ domain do các Manager tạo ra (như Harvests, Orders, Processing Logs, Shipments), **ADMIN có quyền xem (Read-only) để phục vụ giám sát toàn hệ thống, nhưng tuyệt đối không có quyền can thiệp, chỉnh sửa hay khởi tạo (No Write/Edit access)**.
  * **Trách nhiệm**: Khởi tạo, quản lý tài khoản người dùng (Users - bao gồm Drivers) và thiết lập dữ liệu danh mục gốc (Farms, Warehouses, Stores) cùng việc gán quản lý tương ứng cho các thực thể đó.
  * **Cơ chế khởi tạo**: Tài khoản Admin duy nhất được seed trực tiếp thông qua cơ sở dữ liệu Kratos ([seed-admin.json](../deployments/kratos/seed-admin.json)), không dùng SQL seed.
* **`FARM_MANAGER` (Farm Manager)**
  * **Phạm vi**: Quản lý độc quyền các tài nguyên thuộc phạm vi Nông trại được gán.
  * **Trách nhiệm**: Thực hiện các nghiệp vụ nông trại như khai báo mùa vụ thu hoạch (Harvest Declarations).
* **`WAREHOUSE_MGR` (Warehouse & Logistics Manager)**
  * **Phạm vi**: Quản lý toàn bộ nghiệp vụ Kho hàng được gán, có toàn quyền quản lý Phương tiện (Vehicles) và Điều phối Vận tải (Logistics Fleet Assignment).
  * **Trách nhiệm**: Quản lý xuất/nhập kho, tồn kho, chế biến và toàn quyền quyết định về đội xe (Vehicles). Để tinh gọn hệ thống, logic Logistics được gộp trực tiếp vào vai trò này.
  * **Lưu ý UI/UX**: Hệ thống **không cung cấp giao diện CRUD** cho Phương tiện (Vehicles) hay Tài xế (Drivers). Các thực thể này được cung cấp hoàn toàn thông qua Seed Data tĩnh. `WAREHOUSE_MGR` nắm giữ toàn quyền điều động (Assign) tài xế và phương tiện cho các chuyến vận chuyển.
* **`STORE_MGR` (Store Manager)**
  * **Phạm vi**: Quản lý độc quyền Cửa hàng bán lẻ được gán.
  * **Trách nhiệm**: Tiếp nhận và quản lý đơn hàng bán lẻ (Retail Orders), theo dõi doanh thu cửa hàng.
* **`DRIVER` (Logistics Driver)**
  * **Phạm vi**: Tài xế vận chuyển.
  * **Trách nhiệm**: Nhận lệnh điều động giao hàng từ bộ phận Logistics, cập nhật trạng thái hành trình vận chuyển.

---

## 2. Quy Tắc Ràng Buộc Dữ Liệu (Data Integrity Constraints)

Để đảm bảo tính nhất quán của dữ liệu hệ thống, quá trình seeding và vận hành phải tuân thủ nghiêm ngặt các quy tắc ràng buộc sau:

1. **Ràng buộc Quản lý (Manager Assignment)**:
   * Mỗi **Farm** phải được gán duy nhất cho **1 Farm Manager** quản lý trực tiếp.
   * Mỗi **Store** phải được gán duy nhất cho **1 Store Manager** quản lý trực tiếp.
   * Mỗi **Warehouse** phải được gán duy nhất cho **1 Warehouse Manager** quản lý trực tiếp.
2. **Ràng buộc Logistics (Logistics Grid)**:
   * Mỗi **Warehouse** được cấu hình tối thiểu **3 phương tiện** và tối đa **5 phương tiện**.
   * Mỗi **Phương tiện (Vehicle)** phải được gán cố định cho duy nhất **1 Tài xế (Driver)**.
3. **Bản đồ Tuyến đường & Tọa độ Địa lý (Geographical Routing)**:
   * Hệ thống vận hành dựa trên tọa độ địa lý thực tế (Latitude & Longitude) của các Farms, Warehouses và Stores.
   * Để tối ưu hóa quá trình tính toán khoảng cách và mô phỏng lộ trình vận chuyển thực tế, script chuyên biệt tại [main.go](../src/scripts/generate_routes/main.go) được sử dụng để tự động tính toán và khởi tạo (generate) các tuyến đường (routes). Dữ liệu sau đó được kết xuất và lưu giữ tại tệp tĩnh [routes.json](../src/apps/logistics-service/testdata/routes.json) để phục vụ lập kế hoạch vận tải.

---

## 3. Dữ Liệu Thực Thể Cố Định (Master Data Grid)

Hệ thống thiết lập sẵn một mạng lưới thực thể cố định tại Việt Nam để chạy các kịch bản mô phỏng SAGA và Logistics:

### 3.1. Cửa Hàng Bán Lẻ (Stores)
* **Hoan Kiem Store** | Hà Nội | *2 Lý Thái Tổ* | Quản lý: `mgr.store.hoankiem@runtimeroasters.com`
* **Cau Giay Store** | Hà Nội | *102 Trần Thái Tông* | Quản lý: `mgr.store.caugiay@runtimeroasters.com`
* **District 1 Store** | TP. Hồ Chí Minh | *45 Lê Thánh Tôn* | Quản lý: `mgr.store.q1@runtimeroasters.com`
* **District 7 Store** | TP. Hồ Chí Minh | *Phú Mỹ Hưng* | Quản lý: `mgr.store.q7@runtimeroasters.com`
* **Hai Chau Store** | Đà Nẵng | *15 Bạch Đằng* | Quản lý: `mgr.store.haichau@runtimeroasters.com`

### 3.2. Nông Trại Cà Phê (Farms)
* **K'Ho Coffee Farm** | Vùng: `CAU_DAT` | Tọa độ: `11.9404, 108.4442` | Diện tích: `15.5 ha` | Loại hạt: `ARABICA` | Quản lý: `mgr.farm.kho@runtimeroasters.com`
* **Cau Dat Arabica** | Vùng: `CAU_DAT` | Tọa độ: `11.8950, 108.5380` | Diện tích: `45.0 ha` | Loại hạt: `ARABICA` | Quản lý: `mgr.farm.caudat@runtimeroasters.com`
* **Son Pacamara Farm** | Vùng: `CAU_DAT` | Tọa độ: `11.8500, 108.5000` | Diện tích: `12.0 ha` | Loại hạt: `ARABICA` | Quản lý: `mgr.farm.sonpacamara@runtimeroasters.com`
* **Aeroco Coffee** | Vùng: `BUON_MA_THUOT` | Tọa độ: `12.6660, 108.0380` | Diện tích: `20.0 ha` | Loại hạt: `ROBUSTA` | Quản lý: `mgr.farm.aeroco@runtimeroasters.com`
* **Trung Nguyen Village** | Vùng: `BUON_MA_THUOT` | Tọa độ: `12.7000, 108.0500` | Diện tích: `5.0 ha` | Loại hạt: `ROBUSTA` | Quản lý: `mgr.farm.trungnguyen@runtimeroasters.com`
* **Chu Se Estate** | Vùng: `PLEIKU` | Tọa độ: `14.0100, 108.0400` | Diện tích: `30.0 ha` | Loại hạt: `ROBUSTA` | Quản lý: `mgr.farm.chuse@runtimeroasters.com`

### 3.3. Kho Hàng Trung Chuyển (Warehouses)
* **Hoa Lac Warehouse** | Mã kho: `WAREHOUSE-HN-001` | Tọa độ: `21.010, 105.530` | Quản lý: `mgr.warehouse.hn@runtimeroasters.com`
* **Song Than Warehouse** | Mã kho: `WAREHOUSE-HCM-001` | Tọa độ: `10.880, 106.750` | Quản lý: `mgr.warehouse.hcm@runtimeroasters.com`
* **Hoa Khanh Warehouse** | Mã kho: `WAREHOUSE-DN-001` | Tọa độ: `16.080, 108.150` | Quản lý: `mgr.warehouse.dn@runtimeroasters.com`

### 3.4. Đội Xe Vận Tải & Tài Xế (Vehicles & Drivers)

#### Hoa Lac Warehouse (`WAREHOUSE-HN-001`) - Đội xe 4 Phương tiện
* **Vehicle HN-V01**: Biển kiểm soát: `29C-888.01` | Dòng xe: `Hyundai Mighty 110S` (Tải trọng: 7 tấn) | Tài xế: `driver.hn.01@runtimeroasters.com`
* **Vehicle HN-V02**: Biển kiểm soát: `29C-888.02` | Dòng xe: `Isuzu NPR400` (Tải trọng: 3.5 tấn) | Tài xế: `driver.hn.02@runtimeroasters.com`
* **Vehicle HN-V03**: Biển kiểm soát: `29C-888.03` | Dòng xe: `Ford Transit Van` (Tải trọng: 1.5 tấn) | Tài xế: `driver.hn.03@runtimeroasters.com`
* **Vehicle HN-V04**: Biển kiểm soát: `29C-888.04` | Dòng xe: `Suzuki Super Carry` (Tải trọng: 500kg) | Tài xế: `driver.hn.04@runtimeroasters.com`

#### Song Than Warehouse (`WAREHOUSE-HCM-001`) - Đội xe 4 Phương tiện
* **Vehicle HCM-V01**: Biển kiểm soát: `51D-999.01` | Dòng xe: `Hyundai Mighty 110S` (Tải trọng: 7 tấn) | Tài xế: `driver.hcm.01@runtimeroasters.com`
* **Vehicle HCM-V02**: Biển kiểm soát: `51D-999.02` | Dòng xe: `Isuzu NPR400` (Tải trọng: 3.5 tấn) | Tài xế: `driver.hcm.02@runtimeroasters.com`
* **Vehicle HCM-V03**: Biển kiểm soát: `51D-999.03` | Dòng xe: `Ford Transit Van` (Tải trọng: 1.5 tấn) | Tài xế: `driver.hcm.03@runtimeroasters.com`
* **Vehicle HCM-V04**: Biển kiểm soát: `51D-999.04` | Dòng xe: `Suzuki Super Carry` (Tải trọng: 500kg) | Tài xế: `driver.hcm.04@runtimeroasters.com`

#### Hoa Khanh Warehouse (`WAREHOUSE-DN-001`) - Đội xe 3 Phương tiện
* **Vehicle DN-V01**: Biển kiểm soát: `43C-777.01` | Dòng xe: `Isuzu NPR400` (Tải trọng: 3.5 tấn) | Tài xế: `driver.dn.01@runtimeroasters.com`
* **Vehicle DN-V02**: Biển kiểm soát: `43C-777.02` | Dòng xe: `Ford Transit Van` (Tải trọng: 1.5 tấn) | Tài xế: `driver.dn.02@runtimeroasters.com`
* **Vehicle DN-V03**: Biển kiểm soát: `43C-777.03` | Dòng xe: `Suzuki Super Carry` (Tải trọng: 500kg) | Tài xế: `driver.dn.03@runtimeroasters.com`

---

## 4. Tổng Hợp Danh Sách Tài Khoản Người Dùng (Consolidated Users Directory)

Dưới đây là danh sách toàn bộ **26 tài khoản người dùng** được định nghĩa chính thức trong kịch bản khởi tạo hệ thống:

| Email | Họ và Tên | Vai trò (`Role`) | Thực thể được gán (`Associated Entity`) |
|---|---|---|---|
| `admin@runtimeroasters.com` | System Admin | `ADMIN` | Quản lý tài khoản & Dữ liệu gốc (Master Data & Users) |
| `mgr.farm.kho@runtimeroasters.com` | K'Ho Farm Manager | `FARM_MANAGER` | **K'Ho Coffee Farm** |
| `mgr.farm.caudat@runtimeroasters.com` | Cau Dat Farm Manager | `FARM_MANAGER` | **Cau Dat Arabica Farm** |
| `mgr.farm.sonpacamara@runtimeroasters.com` | Son Pacamara Manager | `FARM_MANAGER` | **Son Pacamara Farm** |
| `mgr.farm.aeroco@runtimeroasters.com` | Aeroco Farm Manager | `FARM_MANAGER` | **Aeroco Coffee Farm** |
| `mgr.farm.trungnguyen@runtimeroasters.com` | Trung Nguyen Manager | `FARM_MANAGER` | **Trung Nguyen Village** |
| `mgr.farm.chuse@runtimeroasters.com` | Chu Se Farm Manager | `FARM_MANAGER` | **Chu Se Estate** |
| `mgr.warehouse.hn@runtimeroasters.com` | Hoa Lac Warehouse Manager | `WAREHOUSE_MGR` | **Hoa Lac Warehouse** (`WAREHOUSE-HN-001`) |
| `mgr.warehouse.hcm@runtimeroasters.com` | Song Than Warehouse Manager | `WAREHOUSE_MGR` | **Song Than Warehouse** (`WAREHOUSE-HCM-001`) |
| `mgr.warehouse.dn@runtimeroasters.com` | Hoa Khanh Warehouse Manager | `WAREHOUSE_MGR` | **Hoa Khanh Warehouse** (`WAREHOUSE-DN-001`) |
| `mgr.store.hoankiem@runtimeroasters.com` | Hoan Kiem Store Manager | `STORE_MGR` | **Hoan Kiem Store** |
| `mgr.store.caugiay@runtimeroasters.com` | Cau Giay Store Manager | `STORE_MGR` | **Cau Giay Store** |
| `mgr.store.q1@runtimeroasters.com` | District 1 Store Manager | `STORE_MGR` | **District 1 Store** |
| `mgr.store.q7@runtimeroasters.com` | District 7 Store Manager | `STORE_MGR` | **District 7 Store** |
| `mgr.store.haichau@runtimeroasters.com` | Hai Chau Store Manager | `STORE_MGR` | **Hai Chau Store** |
| `driver.hn.01@runtimeroasters.com` | HN Driver 01 | `DRIVER` | **Hoa Lac Warehouse** (Vehicle: `HN-V01`) |
| `driver.hn.02@runtimeroasters.com` | HN Driver 02 | `DRIVER` | **Hoa Lac Warehouse** (Vehicle: `HN-V02`) |
| `driver.hn.03@runtimeroasters.com` | HN Driver 03 | `DRIVER` | **Hoa Lac Warehouse** (Vehicle: `HN-V03`) |
| `driver.hn.04@runtimeroasters.com` | HN Driver 04 | `DRIVER` | **Hoa Lac Warehouse** (Vehicle: `HN-V04`) |
| `driver.hcm.01@runtimeroasters.com` | HCM Driver 01 | `DRIVER` | **Song Than Warehouse** (Vehicle: `HCM-V01`) |
| `driver.hcm.02@runtimeroasters.com` | HCM Driver 02 | `DRIVER` | **Song Than Warehouse** (Vehicle: `HCM-V02`) |
| `driver.hcm.03@runtimeroasters.com` | HCM Driver 03 | `DRIVER` | **Song Than Warehouse** (Vehicle: `HCM-V03`) |
| `driver.hcm.04@runtimeroasters.com` | HCM Driver 04 | `DRIVER` | **Song Than Warehouse** (Vehicle: `HCM-V04`) |
| `driver.dn.01@runtimeroasters.com` | DN Driver 01 | `DRIVER` | **Hoa Khanh Warehouse** (Vehicle: `DN-V01`) |
| `driver.dn.02@runtimeroasters.com` | DN Driver 02 | `DRIVER` | **Hoa Khanh Warehouse** (Vehicle: `DN-V02`) |
| `driver.dn.03@runtimeroasters.com` | DN Driver 03 | `DRIVER` | **Hoa Khanh Warehouse** (Vehicle: `DN-V03`) |

---

## 5. Kiến Trúc & Chiến Lược Triển Khai Seeding (Seeding Architecture)

### 5.1. Cơ chế Migrations phi tập trung
* Mỗi microservice chịu trách nhiệm hoàn toàn về việc di trú cấu trúc bảng (database migration) của chính nó. 
* Các kịch bản khởi động hoặc tệp cấu hình triển khai phải đảm bảo chạy các tiến trình migration tại từng service một cách độc lập để hỗ trợ tối đa việc scale-out và deploy hạ tầng sau này.
* **Kịch bản Khởi tạo & Dọn dẹp môi trường (Dev Deploy Script)**: Để hỗ trợ các kỹ sư phát triển triển khai nhanh và khởi tạo/seed toàn bộ hệ thống chỉ bằng 1 câu lệnh, dự án đã xây dựng sẵn script chuyên dụng tại [reset-env.sh](../scripts/reset-env.sh). Khi chạy script này:
  1. Dọn dẹp hoàn toàn tài nguyên cũ (`docker compose down -v`).
  2. Dựng các dịch vụ hạ tầng (`docker compose up -d`).
  3. Kích hoạt toàn bộ migrations của các microservice (`task migrate:all`).
  4. Seed dữ liệu hạ tầng tĩnh gồm Admin User và OAuth2 Clients (`task seed:infra`).

### 5.2. Chuẩn hóa Schema
Tất cả các bảng lưu trữ thực thể nghiệp vụ quan trọng bắt buộc phải cấu hình cấu trúc audit chuyên nghiệp để đồng bộ dữ liệu trace, bao gồm các trường:
* `created_at` (Thời điểm tạo)
* `updated_at` (Thời điểm cập nhật mới nhất)
* `created_by` (Định danh người thực hiện - liên kết với ID tài khoản Kratos)

### 5.3. Quy trình & Chiến lược API-driven Seeding (Bootstrapping Flow)

Quy trình seeding dữ liệu mẫu được thiết kế khép kín, bảo mật và trực quan thông qua giao diện Quản trị (Admin UX) theo sơ đồ hoạt động sau:

#### 1. Trải nghiệm Lần đầu tiên (First-run Experience - FRX)
* Khi khởi động hệ thống lần đầu, cơ sở dữ liệu hoàn toàn trống (chỉ có duy nhất tài khoản `ADMIN` được định nghĩa tĩnh trong [seed-admin.json](../deployments/kratos/seed-admin.json) thông qua script [seed.sh](../deployments/seed.sh)).
* Khi Admin đăng nhập lần đầu, giao diện Frontend (`client-app`) thông qua thành phần `SystemBootstrapModal.tsx` sẽ tự động gửi yêu cầu kiểm tra trạng thái seeding qua cổng `/system/status` trên toàn bộ hệ thống.
* Nếu hệ thống ghi nhận trạng thái chưa khởi tạo (unseeded), giao diện sẽ hiển thị một Modal cảnh báo nổi bật kèm nút **"Initialize Data"** (Khởi tạo Dữ liệu).

#### 2. Luồng Điều phối Khởi tạo Dữ liệu (Bootstrap Flow Code)
Khi Admin click vào nút **"Initialize Data"**:
1. **Frontend Call**: Client-app gửi request bảo mật đến API khởi tạo của `auth-service` (Endpoint: `/v1/auth/seed`).
2. **Kratos Provisioning**: 
   * `auth-service` tiếp nhận yêu cầu, tiến hành giao tiếp với Ory Kratos Admin API để tạo động toàn bộ **25 tài khoản người dùng** còn lại (bao gồm Farm Managers, Warehouse Managers, Store Managers, và Drivers) đã định nghĩa ở mục 4.
   * `auth-service` thu thập toàn bộ ID duy nhất (Kratos UUID) do Kratos trả về và đóng gói thành một bản đồ định danh tài khoản (`users_map` ánh xạ `email -> Kratos UUID`).
3. **Downstream Dispatch**:
   * Dùng `users_map` vừa thu được, `auth-service` thực hiện các cuộc gọi **gRPC/HTTP song song** đến các API seed chuyên biệt của các dịch vụ hạ nguồn bao gồm:
     * **Farm Service**: `/v1/system/seed` (hoặc gRPC)
     * **Warehouse Service**: `/v1/system/seed` (hoặc gRPC)
     * **Retail Service**: `/v1/system/seed` (hoặc gRPC)
     * **Logistics Service**: `/v1/system/seed` (hoặc gRPC)
   * Các dịch vụ hạ nguồn tự giải nén `users_map` để ánh xạ chính xác thông tin chủ sở hữu (Kratos UUID) vào dữ liệu Nông trại, Kho hàng, Cửa hàng và Vận tải của mình rồi tiến hành lưu xuống DB của service đó.

#### 3. Cơ chế Bảo mật API-KEY (Security Protection)
Để ngăn chặn các cuộc tấn công chiếm quyền hoặc cố tình khởi tạo lại dữ liệu trái phép:
* Toàn bộ các API endpoint `/v1/system/seed` tại `auth-service` và các service hạ nguồn bắt buộc phải được bảo vệ bằng cơ chế xác thực **API-KEY**.
* **Cấu hình**: API Key được định cấu hình thống nhất qua biến môi trường **`SYSTEM_SEED_API_KEY`** tại tất cả các service:
  ```env
  SYSTEM_SEED_API_KEY=your_highly_secure_seed_api_key_here
  ```
* **Thực thi**: Khi `auth-service` gọi sang các service khác, nó phải đính kèm key này trong header (ví dụ: `X-System-Seed-Key: <key>`). Các service hạ nguồn sẽ so khớp Header này với biến môi trường `SYSTEM_SEED_API_KEY` của chính nó trước khi chấp thuận chạy script seeding dữ liệu.

---

### 3.4. Mô Phỏng Dữ Liệu Nghiệp Vụ Mẫu Hệ Thống (Logic 3: Rich Demo Data Seeder Simulator)

Để kiểm thử nhanh toàn bộ dashboard mà không cần thao tác thủ công hàng trăm đơn hàng, hệ thống tích hợp công cụ simulator chạy độc lập chèn dữ liệu trực tiếp vào cơ sở dữ liệu (Bypass Kafka).

#### 1. Thiết Kế & Quy Tắc
* **Truy vấn Động (No Hardcoded IDs)**: Tuyệt đối không được hardcode bất kỳ ID nào của Warehouse, Farm hay Store trong script. Chương trình bắt buộc phải truy vấn tuần tự dữ liệu Master đã seed trước đó từ các database tương ứng (`farm_db`, `warehouse_db`, `retail_db`) để lấy ID thực tế tại runtime và liên kết chính xác cho các bước tiếp theo.
* **Quy Trình Giả Lập**:
  * Mô phỏng **22 chu kỳ cung ứng nông sản**: Harvests -> Warehouse Pickups -> Roastery Intakes -> Roasting Production Batches (hao hụt `12% - 22%`) -> Cập nhật tồn kho.
  * Mô phỏng **46 đơn hàng bán lẻ SAGA**: Retail Orders -> Payments (85% thành công, 10% chờ, 5% lỗi thẻ) -> Warehouse Stock Reserved -> Logistics Shipment Assignments (Driver & Vehicle) -> Delivered.
  * **Trị số Thời Gian**: Tất cả các bản ghi (`created_at`, `occurred_at`) được backdate phân bổ ngẫu nhiên trong vòng 30 ngày qua để đảm bảo đồ thị phân tích trên UI được vẽ chính xác và sinh động.

#### 2. Kịch Bản Xác Minh Tính Đúng Đắn Của Dữ Liệu Mẫu (Data Verification Strategy)

Sau khi khởi chạy `task seed:demo` (thực thi file `src/scripts/seed_demo/main.go`), ta tiến hành verify tính đúng đắn của dữ liệu mẫu theo các tiêu chí sau:

* **Tính Nhất Quán Giữa Các DB (Cross-DB Consistency)**:
  * Kiểm tra `harvest_id` từ `farm_db.harvests` phải khớp với `harvest_id` được liên kết trong `warehouse_db.pick_up_requests`, `warehouse_db.intakes`, và `logistics_db.shipments`.
  * Kiểm tra `order_id` từ `retail_db.orders` phải khớp tuyệt đối với `payment_db.payments` và `logistics_db.shipments.order_id`.
* **Xác Minh Elasticsearch (Trace Documents Search)**:
  * Chạy truy vấn tìm kiếm dữ liệu đã tổng hợp:
    ```bash
    curl -s http://localhost:9200/coffee_traceability/_search?size=1 | jq
    ```
  * **Tiêu chuẩn đạt**: Phải trả về tài liệu dạng JSON chứa đầy đủ cấu trúc của một chu kỳ nghiệp vụ bao gồm: timeline các bước, các mốc thời gian, chi tiết order và shipment tương ứng.
* **Xác Minh Cassandra (Immutable Cryptographic Audit Logs)**:
  * Chạy truy vấn kiểm tra lịch sử audit log bất biến:
    ```bash
    docker exec -it rr-cassandra cqlsh -e "SELECT * FROM runtime_roasters_audit.audit_logs LIMIT 5;"
    ```
  * **Tiêu chuẩn đạt**: Các bản ghi audit phải được định vị chính xác theo `partition_key` (ví dụ: `order`, `harvest`), lưu đúng `occurred_at`, payload khớp dữ liệu gốc và có `current_hash` mô phỏng SHA-256 chuỗi blockchain.
* **Xác Minh Giao Diện (Frontend Dashboard UI)**:
  * Truy cập `http://localhost:3000/dashboard` và `/dashboard/topology-mesh`.
  * **Tiêu chuẩn đạt**: Bản đồ mesh và biểu đồ doanh thu bán lẻ, biểu đồ hao hụt mẻ rang phải hiển thị dữ liệu lịch sử phong phú, trực quan, không có trạng thái trống (empty state).

