# [RR-7] UI Shell - Next.js Dashboard Setup

- **Summary:** Khởi tạo dự án Frontend và dựng khung giao diện Dashboard.
- **Priority:** `HIGH`
- **Description:** Xây dựng nền móng cho Dashboard sử dụng Next.js 14+ với kiến trúc thành phần hiện đại.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Khởi tạo dự án Next.js
- **Given:** Máy dev đã có Node.js 20+.
- **When:** Tôi chạy lệnh khởi tạo dự án với TypeScript và App Router.
- **Then:** Dự án phải chạy thành công tại `localhost:3000`.

### Scenario 2: Xây dựng Dashboard Layout
- **Given:** CSS Tailwind đã được cài đặt.
- **When:** Tôi thiết kế Sidebar và Header.
- **Then:** Giao diện phải tương thích với các kích thước màn hình và có Sidebar hiển thị các mục "Farms", "Batches".

### Scenario 3: Cấu hình quản lý trạng thái (React Query)
- **Given:** Dự án đã có React Query.
- **When:** Tôi thực hiện cấu hình Provider tại root.
- **Then:** Toàn bộ ứng dụng phải sẵn sàng thực hiện các data fetching gọi API.

---

## 🛠️ Technical Notes
- Next.js 14+ (App Router).
- Tailwind CSS, Lucide Icons.
- TanStack Query (React Query).

## 📋 Sub-tasks
- [ ] Khởi tạo dự án trong `apps/dashboard-ui`.
- [ ] Thiết kế Layout & Theme.
- [ ] Setup Axios client với Interceptor hỗ trợ tracing header.
