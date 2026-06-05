# [RR-8] Client App — Business UI & API Explorer

- **Goal:** Xây dựng apps/client-app/ làm giao diện chính cho người dùng.
- **Business Value:** Cung cấp giao diện trực quan cho người dùng cuối và công cụ API Explorer cho lập trình viên để dễ dàng khám phá và thử nghiệm API.
- **Priority:** `HIGH`

## 🔍 Acceptance Criteria

### Scenario 1: API Explorer Integration
- **Given:** `client-app` đang chạy tại `localhost:3000`.
- **When:** Tôi truy cập `/explorer`.
- **Then:** Swagger UI hiển thị và tự động load `demo.swagger.json` thông qua KrakenD Gateway.

### Scenario 2: Test endpoints via Gateway
- **Given:** API Explorer đã load contract.
- **When:** Tôi "Try it out" endpoint `GET /v1/demo/ping`.
- **Then:** Request được gửi tới `localhost:8081` (KrakenD) và trả về dữ liệu từ Demo Service.

### Scenario 3: Base Layout
- **Then:** App có SideBar/TopBar cơ bản.