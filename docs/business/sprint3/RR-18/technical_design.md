# [RR-18] Technical Design: Farm Control Plane (UI)

**Status:** `DRAFT`
**Author:** Tech Lead

---

## 1. Context & Goal
Triển khai giao diện người dùng cho Farm Service trong ứng dụng `client-app` (Next.js).

---

## 2. Technical Stack
- **Framework:** Next.js (App Router).
- **Styling:** Tailwind CSS.
- **Data Fetching:** TanStack Query (React Query) hoặc SWR.
- **Components:** Headless UI hoặc Radix UI cho các thành phần interactive.

---

## 3. Implementation Details

### 3.1 Routing
- **List Page:** `/app/farms/page.tsx`
- **Create Modal/Page:** `/app/farms/create/page.tsx` hoặc sử dụng Dialog component.

### 3.2 Data Integration
- Gọi API qua KrakenD Gateway: `GET /api/v1/farms`.
- Header: Luôn đính kèm `Authorization: Bearer <token>` từ session hiện tại.

### 3.3 Components Structure
- `FarmTable`: Hiển thị danh sách farm với pagination.
- `FarmForm`: Reusable form cho Create/Update farm.
- `FarmTypeSelect`: Dropdown được seed dữ liệu từ config hoặc API.

---

## 4. Security
- **Client-side Guard:** Kiểm tra role trong JWT (phía client) để ẩn/hiện menu "Farms".
- **Middleware Guard:** Sử dụng Next.js Middleware để redirect người dùng không có quyền truy cập vào route `/app/farms`.
