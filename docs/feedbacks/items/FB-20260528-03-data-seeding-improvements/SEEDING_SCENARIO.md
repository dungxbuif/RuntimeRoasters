# E2E Saga Seeding Scenario: Demo-Ready Farm-to-Cup Dataset

Tài liệu này định nghĩa bộ dữ liệu seed mục tiêu cho ticket `FB-20260528-03-data-seeding-improvements`.

Mục tiêu không phải thay seed cũ bằng một bộ dữ liệu khác, mà là **kế thừa seed hiện tại**, loại bỏ dữ liệu random, và bổ sung dữ liệu lịch sử đủ giàu để reviewer mở hệ thống lên là có thể demo ngay:

- `/dashboard/users` có đủ identity theo role.
- `/dashboard/traceability` có sẵn trace document và timeline.
- `/dashboard/finance` có nhiều trạng thái thanh toán.
- Logistics có cả shipment đã hoàn tất và shipment đang chạy.
- Retail, Warehouse, Farm vẫn giữ các seed hiện có để không phá demo cũ.

## Quy ước chung

- `T-0` là thời điểm Admin bấm **Initialize Data**.
- Tất cả timestamp dùng relative time theo `T-0`, không hard-code ngày cụ thể.
- Không dùng random number, `Date.now()`, hoặc timestamp để tạo email, business key, provider ref, idempotency key, trace key.
- ID kỹ thuật có thể được gen động bằng code nếu không cần người dùng nhập lại hoặc không phải reference cố định trong tài liệu demo.
- Các giá trị cần ổn định cho demo phải deterministic: email, store/warehouse/location ID hiện có, SKU, batch code, trace key, provider ref, idempotency key.
- Seed phải idempotent: chạy lại không tạo bản ghi trùng.
- Default password cho toàn bộ account demo: `Hello@123`.

## 0. Seed Levels Và Lựa Chọn Hợp Lý

Hệ thống nên duy trì 2 level seed rõ ràng thay vì gom tất cả vào một chỗ.

### 0.1 SQL Seed: Static Master Data

SQL seed dùng cho dữ liệu nền ít thay đổi, cần tồn tại ngay sau migration hoặc reset DB.

Giữ ở SQL seed:

- Farms master records.
- Retail stores master records.
- Logistics locations.
- Vehicles.
- Drivers shell records.
- Warehouse inventory baseline.
- Payment webhook demo key.

Nguyên tắc:

- SQL seed được phép dùng stable ID hiện tại cho master data, vì các ID này đang được frontend, integration tests và docs tham chiếu.
- Không đưa dữ liệu lịch sử Saga phức tạp vào migration SQL nếu dữ liệu đó cần liên kết Kratos subject, event timeline, trace projection hoặc Elasticsearch.
- SQL seed chỉ nên upsert static rows, không thực hiện orchestration cross-service.

### 0.2 Admin Seed: Demo Scenario Data

Admin seed là luồng chạy khi ADMIN bấm **Initialize Data** trong UI. Đây là nơi hợp lý để tạo dữ liệu động có logic nghiệp vụ.

Admin seed xử lý:

- Kratos identities và Casbin grouping.
- Mapping farm owner theo Kratos subject thật.
- Historical orders, payments, shipments, warehouse processing.
- Trace events, trace documents, Elasticsearch index.
- Audit records.

Nguyên tắc:

- Được phép gen dynamic UUID bằng code cho primary key kỹ thuật như `orders.id`, `payments.id`, `shipments.id`, `trace_events.id`, `audit_logs.id`.
- Nhưng phải lưu mapping theo deterministic business key để idempotent, ví dụ `seed_key`, `idempotency_key`, `provider_ref`, `message_id`, `trace_id`, `batch_id`.
- Nếu table chưa có `seed_key`, dùng unique field sẵn có: `idempotency_key`, `provider_ref`, `message_id`, `sku`, `batch_id`.
- Dữ liệu trong docs nên ghi **Business Key** là nguồn tham chiếu chính; UUID cố định chỉ dùng khi table hiện tại đã có stable ID từ seed cũ.

## 1. Static Master Data Cần Giữ Lại

Bộ data hiện tại đã ổn ở phần master data. Không xóa hoặc thu nhỏ các nhóm này.

### 1.1 Farms

Giữ 6 farms hiện có trong `farm_db.farms`:

| Seed Key | Farm Name | Region | Coffee Type | Manager Email |
| --- | --- | --- | --- | --- |
| `FARM-CAUDAT-001` | K'Ho Coffee Farm | Cầu Đất | ARABICA | `mgr.caudat@runtimeroasters.com` |
| `FARM-CAUDAT-002` | Cau Dat Arabica | Cầu Đất | ARABICA | `mgr.caudat@runtimeroasters.com` |
| `FARM-CAUDAT-003` | Son Pacamara Farm | Cầu Đất | ARABICA | `mgr.caudat@runtimeroasters.com` |
| `FARM-BMT-001` | Aeroco Coffee | Buôn Ma Thuột | ROBUSTA | `mgr.bmt@runtimeroasters.com` |
| `FARM-BMT-002` | Trung Nguyen Village | Buôn Ma Thuột | ROBUSTA | `mgr.bmt@runtimeroasters.com` |
| `FARM-PLEIKU-001` | Chu Se Estate | Pleiku | ROBUSTA | `mgr.pleiku@runtimeroasters.com` |

Ghi chú triển khai sau này:

- `farms.owner_id` phải map về Kratos subject thật của manager tương ứng.
- Nếu farm table vẫn dùng numeric `id`, seed scenario dùng `Seed Key` làm external reference trong trace/logistics.

### 1.2 Retail Stores

Giữ 5 stores hiện có trong `retail_db.stores`:

| Store ID | Store Name | City | Manager Email |
| --- | --- | --- | --- |
| `11111111-1111-1111-1111-111111111101` | Hoan Kiem Store | Hanoi | `mgr.hn.hoankiem@runtimeroasters.com` |
| `11111111-1111-1111-1111-111111111102` | Cau Giay Store | Hanoi | `mgr.hn.caugiay@runtimeroasters.com` |
| `11111111-1111-1111-1111-111111111103` | District 1 Store | Ho Chi Minh City | `mgr.hcm.q1@runtimeroasters.com` |
| `11111111-1111-1111-1111-111111111104` | District 7 Store | Ho Chi Minh City | `mgr.hcm.q7@runtimeroasters.com` |
| `11111111-1111-1111-1111-111111111105` | Hai Chau Store | Da Nang | `mgr.dn.haichau@runtimeroasters.com` |

### 1.3 Warehouses / Roasteries

Giữ 3 warehouse locations hiện có trong logistics seed:

| Warehouse ID | Name | Region | Manager Email |
| --- | --- | --- | --- |
| `WAREHOUSE-HN-001` | Hanoi Roastery | Hanoi | `mgr.hoalac@runtimeroasters.com` |
| `WAREHOUSE-HCM-001` | Song Than Roastery | Binh Duong / HCM | `mgr.songthan@runtimeroasters.com` |
| `WAREHOUSE-DN-001` | Hoa Khanh Roastery | Da Nang | `mgr.hoakhanh@runtimeroasters.com` |

### 1.4 Logistics Fleet

Giữ 3 vehicles và 3 drivers hiện có, nhưng chuẩn hóa `drivers.user_id` thành deterministic email:

| Driver ID | Email | Name | Vehicle ID | Home Warehouse |
| --- | --- | --- | --- | --- |
| `22222222-2222-2222-2222-222222222201` | `driver.hoalac01@runtimeroasters.com` | Driver Alpha | `VEHICLE-DEMO-001` | `WAREHOUSE-HN-001` |
| `22222222-2222-2222-2222-222222222202` | `driver.hoalac02@runtimeroasters.com` | Driver Beta | `VEHICLE-DEMO-002` | `WAREHOUSE-HN-001` |
| `22222222-2222-2222-2222-222222222203` | `driver.songthan01@runtimeroasters.com` | Driver Gamma | `VEHICLE-DEMO-003` | `WAREHOUSE-HCM-001` |

## 2. Identity Seed

Seed đầy đủ các role cần cho demo. Email phải lấy theo entity mà account quản lý hoặc tên role rút gọn, không dùng số random.

| Email | Role | Name | Scope |
| --- | --- | --- | --- |
| `admin@runtimeroasters.com` | `ADMIN` | System Administrator | Global |
| `mgr.caudat@runtimeroasters.com` | `FARM_MANAGER` | Cau Dat Farm Manager | Farms Cầu Đất |
| `mgr.bmt@runtimeroasters.com` | `FARM_MANAGER` | Buon Ma Thuot Farm Manager | Farms Buôn Ma Thuột |
| `mgr.pleiku@runtimeroasters.com` | `FARM_MANAGER` | Pleiku Farm Manager | Farms Pleiku |
| `mgr.hoalac@runtimeroasters.com` | `WAREHOUSE_MGR` | Hoa Lac Warehouse Manager | `WAREHOUSE-HN-001` |
| `mgr.songthan@runtimeroasters.com` | `WAREHOUSE_MGR` | Song Than Warehouse Manager | `WAREHOUSE-HCM-001` |
| `mgr.hoakhanh@runtimeroasters.com` | `WAREHOUSE_MGR` | Hoa Khanh Warehouse Manager | `WAREHOUSE-DN-001` |
| `mgr.hn.hoankiem@runtimeroasters.com` | `STORE_MGR` | Hoan Kiem Store Manager | `11111111-1111-1111-1111-111111111101` |
| `mgr.hn.caugiay@runtimeroasters.com` | `STORE_MGR` | Cau Giay Store Manager | `11111111-1111-1111-1111-111111111102` |
| `mgr.hcm.q1@runtimeroasters.com` | `STORE_MGR` | HCM District 1 Store Manager | `11111111-1111-1111-1111-111111111103` |
| `mgr.hcm.q7@runtimeroasters.com` | `STORE_MGR` | HCM District 7 Store Manager | `11111111-1111-1111-1111-111111111104` |
| `mgr.dn.haichau@runtimeroasters.com` | `STORE_MGR` | Hai Chau Store Manager | `11111111-1111-1111-1111-111111111105` |
| `driver.hoalac01@runtimeroasters.com` | `DRIVER` | Driver Alpha | `22222222-2222-2222-2222-222222222201` |
| `driver.hoalac02@runtimeroasters.com` | `DRIVER` | Driver Beta | `22222222-2222-2222-2222-222222222202` |
| `driver.songthan01@runtimeroasters.com` | `DRIVER` | Driver Gamma | `22222222-2222-2222-2222-222222222203` |

Kratos traits cần có:

- `role`: role trong bảng trên.
- `store_ids`: chỉ set cho `STORE_MGR`.
- `warehouse_ids`: chỉ set cho `WAREHOUSE_MGR`.
- `org_id`: `runtime-roasters-demo`.

Casbin:

- Sync role từ Kratos traits sang Casbin grouping policy.
- Không seed policy riêng cho từng user nếu role policy đã đủ.

## 3. Historical Scenario A: Completed Farm-to-Cup Trace

Scenario chính để demo `/dashboard/traceability`. Đây là luồng hoàn tất từ Farm tới Store.

### 3.1 Business Keys And Generated IDs

Các key dưới đây là reference demo ổn định. Primary key kỹ thuật có thể gen động nếu có unique business key để lookup lại.

| Entity | Stable Demo Reference | ID Policy |
| --- | --- | --- |
| Trace ID | `trc-seed-caudat-songthan-hcm-q1` | Must stay stable |
| Harvest key | `seed-harvest-caudat-arabica-001` | Can map to generated `harvests.id` |
| Pickup request key | `seed-pickup-caudat-001` | Can map to generated/explicit pickup ID |
| Inbound shipment key | `seed-shipment-in-caudat-songthan-001` | Can map to generated `shipments.id` |
| Intake key | `seed-intake-caudat-001` | Can map to generated/explicit intake ID |
| Production batch code | `BATCH-CAUDAT-AR-MED-001` | Must stay stable for Traceability search |
| SKU | `SKU-AR-VN-LD-001` | Must stay stable |
| Order idempotency key | `seed-order-caudat-hcm-q1` | Can map to generated `orders.id` |
| Payment provider ref | `pi_seed_caudat_hcm_q1` | Can map to generated `payments.id` |
| Outbound shipment key | `seed-shipment-out-caudat-hcm-q1-001` | Can map to generated `shipments.id` |
| Farm location key | `FARM-CAUDAT-002` | Existing static seed ID |
| Warehouse ID | `WAREHOUSE-HCM-001` | Existing static seed ID |
| Store ID | `11111111-1111-1111-1111-111111111103` | Existing static seed ID |
| Driver ID | `22222222-2222-2222-2222-222222222203` | Existing static seed ID |
| Vehicle ID | `VEHICLE-DEMO-003` | Existing static seed ID |

### 3.2 Farm Records

| Table | Record |
| --- | --- |
| `harvests` | `seed_key = seed-harvest-caudat-arabica-001` if supported, `farm_id = Cau Dat Arabica`, `owner_id = mgr.caudat subject`, `coffee_type = ARABICA`, `quantity = 1000.00`, `status = COMPLETED`, `harvest_date = T-7d 08:00`, `notes = "Seeded specialty Arabica lot for flagship trace demo"` |

Expected event:

| Time | Topic | Main Payload |
| --- | --- | --- |
| `T-7d 08:00` | `farm.harvest.created` | harvest, farm, coffee type, quantity, trace ID |

### 3.3 Inbound Logistics Records

| Table | Record |
| --- | --- |
| `shipments` | `seed_key = seed-shipment-in-caudat-songthan-001` if supported, `type = FARM_PICKUP`, `farm_id = FARM-CAUDAT-002`, `harvest_id = resolved seed-harvest-caudat-arabica-001`, `warehouse_id = WAREHOUSE-HCM-001`, `driver_id = Driver Gamma`, `vehicle_id = VEHICLE-DEMO-003`, `status = DELIVERED`, `current_leg = COMPLETED` |

Expected events:

| Time | Topic | Status |
| --- | --- | --- |
| `T-6d 08:00` | `logistics.pickup.assigned` | ASSIGNED |
| `T-6d 09:00` | `logistics.pickup.departed` | DEPARTED |
| `T-6d 13:00` | `logistics.pickup.arrived_at_farm` | ARRIVED_AT_FARM |
| `T-6d 13:30` | `logistics.pickup.loading_confirmed` | LOADING_CONFIRMED |
| `T-6d 14:00` | `logistics.pickup.return_started` | RETURN_STARTED |
| `T-6d 20:00` | `logistics.pickup.arrived_at_warehouse` | ARRIVED_AT_WAREHOUSE |
| `T-6d 20:15` | `logistics.pickup.completed` | COMPLETED |

### 3.4 Warehouse Processing Records

| Table | Record |
| --- | --- |
| `pick_up_requests` | `seed_key = seed-pickup-caudat-001` if supported, `harvest_id = resolved seed-harvest-caudat-arabica-001`, `warehouse_id = WAREHOUSE-HCM-001`, `quantity = 1000.00`, `coffee_type = ARABICA`, `origin_code = VN-LD`, `status = RECEIVED` |
| `intakes` | `seed_key = seed-intake-caudat-001` if supported, `harvest_id = resolved seed-harvest-caudat-arabica-001`, `warehouse_id = WAREHOUSE-HCM-001`, `coffee_type = ARABICA`, `origin_code = VN-LD`, `quantity = 1000.00`, `status = ASSIGNED`, `batch_id = BATCH-CAUDAT-AR-MED-001` |
| `production_batches` | `id = BATCH-CAUDAT-AR-MED-001`, `batch_id = BATCH-CAUDAT-AR-MED-001`, `warehouse_id = WAREHOUSE-HCM-001`, `status = FINALIZED`, `total_input_weight = 1000.00`, `total_output_weight = 850.00`, `weight_loss_percent = 15.00` |
| `roast_runs` | `input_weight = 1000.00`, `output_weight = 850.00`, `status = COMPLETED` |
| `inventories` | `sku = SKU-AR-VN-LD-001`, `available_quantity = 840.00` after order reservation |

Expected events:

| Time | Topic | Status |
| --- | --- | --- |
| `T-5d 09:00` | `warehouse.pickup.received` | RECEIVED |
| `T-5d 10:00` | `warehouse.intake.created` | INTAKE_CREATED |
| `T-5d 16:00` | `warehouse.inventory.updated` | STOCKED |

### 3.5 Retail Order And Payment Records

| Table | Record |
| --- | --- |
| `orders` | `id = generated`, `store_id = HCM District 1`, `items = [{"sku":"SKU-AR-VN-LD-001","quantity":10}]`, `total_amount = 1200000.00`, `status = COMPLETED`, `idempotency_key = seed-order-caudat-hcm-q1` |
| `payments` | `id = generated`, `order_id = resolved by idempotency_key seed-order-caudat-hcm-q1`, `store_id = HCM District 1`, `provider = STRIPE`, `provider_ref = pi_seed_caudat_hcm_q1`, `amount = 1200000.00`, `currency = VND`, `status = SUCCEEDED`, `simulated = true` |

Expected events:

| Time | Topic | Status |
| --- | --- | --- |
| `T-2d 09:00` | `retail.order.created` | ORDER_CREATED |
| `T-2d 09:01` | `payment.intent.created` | PENDING |
| `T-2d 09:03` | `payment.simulated_completed` | SUCCEEDED |
| `T-2d 09:04` | `warehouse.stock.reserved` | RESERVED |
| `T-2d 09:05` | `warehouse.dispatch.requested` | DISPATCH_REQUESTED |

### 3.6 Outbound Delivery Records

| Table | Record |
| --- | --- |
| `shipments` | `seed_key = seed-shipment-out-caudat-hcm-q1-001` if supported, `type = RETAIL_DELIVERY`, `order_id = resolved by idempotency_key seed-order-caudat-hcm-q1`, `origin_warehouse_id = WAREHOUSE-HCM-001`, `destination_store_id = HCM District 1`, `driver_id = Driver Gamma`, `vehicle_id = VEHICLE-DEMO-003`, `status = RETURNED_TO_BASE`, `current_leg = COMPLETED` |

Expected events:

| Time | Topic | Status |
| --- | --- | --- |
| `T-1d 08:00` | `logistics.delivery.assigned` | ASSIGNED |
| `T-1d 08:30` | `logistics.delivery.departed` | DEPARTED |
| `T-1d 10:00` | `logistics.delivery.arrived_at_store` | ARRIVED_AT_STORE |
| `T-1d 10:10` | `logistics.delivery.driver_confirmed` | DRIVER_CONFIRMED |
| `T-1d 10:15` | `logistics.delivery.completed` | DELIVERED |
| `T-1d 11:00` | `logistics.driver.return_started` | RETURN_STARTED |
| `T-1d 12:30` | `logistics.driver.returned_to_base` | RETURNED_TO_BASE |

## 4. Historical Scenario B: Active Logistics Demo

Scenario này giúp Logistics Map có dữ liệu đang chạy, nhưng không hiển thị route mặc định khi chưa có shipment.

### 4.1 Business Keys And Generated IDs

| Entity | Stable Demo Reference | ID Policy |
| --- | --- | --- |
| Trace ID | `trc-seed-bmt-songthan-hcm-q7-active` | Must stay stable |
| Order idempotency key | `seed-order-bmt-hcm-q7-active` | Can map to generated `orders.id` |
| Payment provider ref | `vnpay_seed_bmt_hcm_q7_active` | Can map to generated `payments.id` |
| Shipment key | `seed-shipment-active-bmt-hcm-q7` | Can map to generated `shipments.id` |
| Store | `11111111-1111-1111-1111-111111111104` | Existing static seed ID |
| Warehouse | `WAREHOUSE-HCM-001` | Existing static seed ID |
| Driver | `22222222-2222-2222-2222-222222222203` | Existing static seed ID |
| Vehicle | `VEHICLE-DEMO-003` | Existing static seed ID |
| SKU | `SKU-RB-VN-DL-001` | Existing/static SKU |

### 4.2 Records

| Table | Record |
| --- | --- |
| `orders` | `id = generated`, `idempotency_key = seed-order-bmt-hcm-q7-active`, `status = SHIPPING`, `total_amount = 760000.00`, `items = [{"sku":"SKU-RB-VN-DL-001","quantity":8}]` |
| `payments` | `status = SUCCEEDED`, `provider = VNPAY`, `provider_ref = vnpay_seed_bmt_hcm_q7_active`, `amount = 760000.00`, `currency = VND` |
| `shipments` | `id = generated`, `seed_key = seed-shipment-active-bmt-hcm-q7` if supported, `status = IN_TRANSIT`, `current_leg = TO_STORE`, `departed_at = T-2h`, `arrived_destination_at = null`, `delivered_at = null` |
| `drivers` | Driver Gamma `status = BUSY`, `is_available = false`, `current_shipment_id = resolved seed-shipment-active-bmt-hcm-q7` |
| `vehicles` | `VEHICLE-DEMO-003 status = BUSY`, current coordinates between Song Than and HCM Q7 |

Expected events:

| Time | Topic | Status |
| --- | --- | --- |
| `T-3h` | `retail.order.created` | ORDER_CREATED |
| `T-2h 58m` | `payment.intent.created` | PENDING |
| `T-2h 55m` | `payment.simulated_completed` | SUCCEEDED |
| `T-2h 50m` | `warehouse.stock.reserved` | RESERVED |
| `T-2h 45m` | `logistics.delivery.assigned` | ASSIGNED |
| `T-2h` | `logistics.delivery.departed` | DEPARTED |
| `T-90m` | `logistics.gps.updated` | IN_TRANSIT |
| `T-60m` | `logistics.gps.updated` | IN_TRANSIT |
| `T-30m` | `logistics.gps.updated` | IN_TRANSIT |

UI expectation:

- Logistics map chỉ show route/vehicle movement cho shipment này.
- Các locations khác vẫn hiển thị là point-of-interest.
- Không auto-show route cho các shipment đã completed.

## 5. Historical Scenario C: Finance Integrity Cases

Scenario này giúp Finance không chỉ có một dòng success.

### 5.1 Failed Payment

| Entity | ID / Value |
| --- | --- |
| Trace ID | `trc-seed-payment-failed-hn-caugiay` |
| Order idempotency key | `seed-order-failed-hn-caugiay` |
| Payment provider ref | `pi_seed_failed_hn_caugiay` |
| Store | `11111111-1111-1111-1111-111111111102` |
| Amount | `430000.00 VND` |
| Provider | `STRIPE` |
| Payment Status | `FAILED` |
| Order Status | `REJECTED` |

Expected events:

- `retail.order.created`
- `payment.intent.created`
- `payment.failed`

### 5.2 Refunded Payment

| Entity | ID / Value |
| --- | --- |
| Trace ID | `trc-seed-payment-refunded-hn-hoankiem` |
| Order idempotency key | `seed-order-refunded-hn-hoankiem` |
| Payment provider ref | `pi_seed_refunded_hn_hoankiem` |
| Store | `11111111-1111-1111-1111-111111111101` |
| Amount | `980000.00 VND` |
| Provider | `STRIPE` |
| Refund Ref | `re_seed_hn_hoankiem_quality_issue` |
| Payment Status | `REFUNDED` |
| Order Status | `REJECTED` or `REFUNDED` depending current enum support |

Expected events:

- `retail.order.created`
- `payment.intent.created`
- `payment.simulated_completed`
- `payment.refunded`

## 6. Trace Read Models

Trace seed phải tạo đủ dữ liệu để UI không trống.

### 6.1 `trace_events`

Mỗi event trong các scenario trên cần insert vào `trace_db.trace_events` với:

- `message_id`: deterministic, ví dụ `seed:<topic>:<business-id>`.
- `topic`: đúng topic trong `pkg/events`.
- `trace_id`: trace ID của scenario.
- `flow_id`, `node_id`, `edge_id`, `pattern`, `source_service`, `visibility`: lấy từ topology projection hiện có.
- `payload`: CloudEvent JSON hợp lệ, có `correlationid`.
- `display_payload`: payload đã sanitize để topology/history render được.
- `occurred_at`: timestamp relative theo timeline.

### 6.2 `trace_documents`

Tối thiểu phải có documents:

| Entity ID | Entity Type | Store ID | Purpose |
| --- | --- | --- | --- |
| `BATCH-CAUDAT-AR-MED-001` | `batch_id` | `11111111-1111-1111-1111-111111111103` | Default Traceability demo |
| Resolved order from `seed-order-caudat-hcm-q1` | `order_id` | `11111111-1111-1111-1111-111111111103` | Retail order trace |
| Resolved shipment from `seed-shipment-out-caudat-hcm-q1-001` | `shipment_id` | `11111111-1111-1111-1111-111111111103` | Completed delivery trace |
| Resolved shipment from `seed-shipment-active-bmt-hcm-q7` | `shipment_id` | `11111111-1111-1111-1111-111111111104` | Active delivery trace |

Dashboard behavior expectation:

- `/dashboard/traceability` nên có shortcut/search suggestion cho `BATCH-CAUDAT-AR-MED-001`.
- Nếu chưa có UI list, default demo ID trong guide là `BATCH-CAUDAT-AR-MED-001`.

### 6.3 Elasticsearch

Index `coffee_traceability` cần mirror các `trace_documents` ở trên để public trace hoặc search không rỗng.

## 7. Finance Dashboard Expected State

Sau seed, `/dashboard/finance` phải có ít nhất 4 payments:

| Payment Ref | Store | Provider | Amount | Status |
| --- | --- | --- | --- | --- |
| `pi_seed_caudat_hcm_q1` | HCM Q1 | STRIPE | `1200000.00 VND` | `SUCCEEDED` |
| `vnpay_seed_bmt_hcm_q7_active` | HCM Q7 | VNPAY | `760000.00 VND` | `SUCCEEDED` |
| `pi_seed_failed_hn_caugiay` | Cau Giay | STRIPE | `430000.00 VND` | `FAILED` |
| `pi_seed_refunded_hn_hoankiem` | Hoan Kiem | STRIPE | `980000.00 VND` | `REFUNDED` |

Expected summary:

- Total payments: `4`.
- Successful count: `2`.
- Successful volume: `1,960,000 VND`.
- Table demonstrates webhook pass/fail/refund history.

## 8. Logistics Dashboard Expected State

Sau seed, logistics data phải có:

- 14 locations.
- 3 vehicles.
- 3 drivers.
- 1 completed inbound pickup shipment.
- 1 completed outbound retail delivery shipment.
- 1 active outbound retail delivery shipment.

Vehicle/driver final state:

| Driver | Vehicle | Status | Reason |
| --- | --- | --- | --- |
| Driver Alpha | `VEHICLE-DEMO-001` | `IDLE` | Available for next Hanoi demo |
| Driver Beta | `VEHICLE-DEMO-002` | `IDLE` | Available for next Hanoi demo |
| Driver Gamma | `VEHICLE-DEMO-003` | `BUSY` | Active HCM Q7 delivery |

## 9. Audit Data

Seed audit records for the main scenario only, enough to show immutable event history without overloading Cassandra/Postgres fallback.

Minimum audit topics:

- `farm.harvest.created`
- `warehouse.intake.created`
- `retail.order.created`
- `payment.simulated_completed`
- `warehouse.stock.reserved`
- `logistics.delivery.completed`

Each audit record:

- deterministic `message_id`.
- non-empty `partition_key`.
- valid `previous_hash` / `current_hash` chain within the seeded partition.
- `store_id` set for retail/payment/logistics events.

## 10. Bootstrap Execution Order

### 10.1 Before Admin Seed: SQL Seed Baseline

Khi DB migration hoặc reset infrastructure hoàn tất, SQL seed đảm bảo master data nền đã có sẵn:

1. Farms.
2. Retail stores.
3. Logistics locations.
4. Vehicles and driver shell records.
5. Warehouse inventory baseline.
6. Payment webhook demo key.

### 10.2 Admin Seed: Scenario Enrichment

Khi Admin bấm **Initialize Data**, admin/code seed chạy theo thứ tự:

1. Auth identities and Casbin grouping policies.
2. Farm managers mapped to farm ownership.
3. Retail stores preserved/upserted.
4. Logistics locations, vehicles, drivers preserved/upserted.
5. Warehouse inventory and processing records.
6. Retail orders.
7. Payment records.
8. Logistics shipments and current driver/vehicle state.
9. Trace events and trace documents.
10. Elasticsearch trace index.
11. Audit logs.

Không truncate trong endpoint seed bình thường. Reset dữ liệu demo thuộc trách nhiệm của `deployments/reset-demo-state.sh`.

Admin seed có thể sinh dynamic primary key, nhưng phải ghi lại reference qua deterministic business key để các bước sau resolve được đúng bản ghi.

## 11. Acceptance Criteria

Dataset được xem là đủ tốt cho demo khi:

- Không còn email kiểu `mgr.<random>@runtimeroasters.com`.
- `/v1/users` có `ADMIN`, `FARM_MANAGER`, `WAREHOUSE_MGR`, `STORE_MGR`, `DRIVER`.
- Store manager login chỉ thấy store scope của mình.
- Warehouse manager login chỉ thấy warehouse scope của mình.
- Driver login map được vào đúng `drivers.user_id`.
- `/dashboard/traceability` không còn empty state sau bootstrap; demo được bằng `BATCH-CAUDAT-AR-MED-001`.
- `/dashboard/finance` không còn empty state và có đủ success/failed/refunded.
- Logistics map có active shipment thật để show route, không show route mặc định khi không có shipment active.
- Chạy seed lại lần 2 không tăng count bất thường và không tạo bản ghi trùng.
