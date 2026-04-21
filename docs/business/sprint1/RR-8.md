# [RR-8] Client App — Business UI Shell

- **Summary:** Xây dựng layout và shell cho Business UI — khu vực `/app/*` trong `apps/client-app/`.
- **Priority:** `MEDIUM`
- **Depends on:** RR-7

---

## User Story

> As a user, I want a business operations dashboard with a clear layout and navigation, so that I can navigate to the core business features when they become available in Sprint 2.

---

## Acceptance Criteria

### Scenario 1: Layout chuẩn responsive
- **Given:** Client app đang chạy tại `localhost:3000`.
- **When:** Tôi truy cập `/app`.
- **Then:** Layout hiển thị Sidebar với navigation items (Farms, Batches, Logistics, Warehouse, Retail) và Main Content area — responsive trên desktop và tablet.

### Scenario 2: Dashboard landing page
- **Given:** Tôi đã đăng nhập.
- **When:** Tôi truy cập `/app`.
- **Then:** Trang hiển thị các placeholder sections: Supply Chain Overview, Recent Events, Quick Actions — đủ để demo UX mà chưa cần real data.

### Scenario 3: Navigation hoạt động
- **Given:** Tôi đang ở `/app`.
- **When:** Tôi click vào "Farms" trong Sidebar.
- **Then:** URL thay đổi thành `/app/farms` và hiển thị placeholder "Farm management — coming in Sprint 2".

### Scenario 4: Business UI call KrakenD trực tiếp
- **Given:** `NEXT_PUBLIC_GATEWAY_URL=http://localhost:8081` được set.
- **When:** Business UI cần gọi API.
- **Then:** Request đi thẳng từ browser đến KrakenD `:8081` — không qua BFF Route Handler. Không có `fetch` call nào trong Route Handler cho business data.
