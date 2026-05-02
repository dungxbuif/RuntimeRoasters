# [RR-8] Client App — Business UI & API Explorer

- **Summary:** Xây dựng `apps/client-app/` làm giao diện chính cho người dùng và tích hợp API Explorer.
- **Priority:** `HIGH`
- **Status:** `TO_DO`
- **Depends on:** RR-4

---

## User Story

> As a developer/user, I want a single client application where I can access the business features and an integrated API Explorer to test endpoints directly through the Gateway.

---

## Acceptance Criteria

### Scenario 1: API Explorer Integration
- **Given:** `client-app` đang chạy tại `localhost:3000`.
- **When:** Tôi truy cập `/explorer` (hoặc tích hợp trong layout).
- **Then:** Swagger UI hiển thị và tự động load `demo.swagger.json` thông qua KrakenD Gateway (`localhost:8081/swagger/demo.swagger.json`).

### Scenario 2: Test endpoints via Gateway
- **Given:** API Explorer đã load contract.
- **When:** Tôi "Try it out" endpoint `GET /v1/demo/ping`.
- **Then:** Request được gửi tới `localhost:8081` (KrakenD) và trả về dữ liệu từ Demo Service.

### Scenario 3: Base Layout
- **Then:** App có SideBar/TopBar cơ bản, cho phép chuyển đổi giữa Dashboard và API Explorer.

---

## 🛠️ Technical Notes
- Framework: Next.js 15 (App Router).
- Swagger UI: `swagger-ui-react`.
- Gateway Endpoint cho Swagger: KrakenD cần proxy `/swagger/*` tới Demo Service.
