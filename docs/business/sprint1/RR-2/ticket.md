# [RR-2] Core Framework - Config & Logger

- **Goal:** Xây dựng hệ thống cấu hình và ghi log dùng chung (Shared Infrastructure).
- **Business Value:** Thống nhất cách cấu hình và ghi log cho toàn bộ hệ thống, giúp dễ dàng vận hành và debug.
- **Priority:** `HIGH`

## 📝 Description
Implement module load config động (Viper) và Logger có cấu trúc (Uber Zap) để thống nhất cách vận hành cho 10+ services.

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Load cấu hình động từ .env
- **Given:** File `.env` chứa biến `APP_PORT=8080`.
- **When:** Tôi khởi động service và gọi hàm `config.LoadConfig()`.
- **Then:** Giá trị `8080` phải được map chính xác vào struct cấu hình của Go.

### Scenario 2: Ghi log có cấu trúc trong môi trường Development
- **Given:** Service đang chạy ở chế độ `development`.
- **When:** Tôi gọi `logger.Info("Hello World")`.
- **Then:** Log phải hiển thị ở định dạng **Console** có màu sắc và thông tin dòng code (caller).

### Scenario 3: Ghi log có cấu trúc trong môi trường Production
- **Given:** Service đang chạy ở chế độ `production`.
- **When:** Tôi gọi `logger.Info("Hello World")`.
- **Then:** Log phải hiển thị ở định dạng **JSON** chuẩn để công cụ quản lý log có thể đọc được.