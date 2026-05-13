# Database Design: Farm Service (The Traceability Origin)

Tài liệu này định nghĩa cấu trúc dữ liệu tối thiểu nhưng đầy đủ để trình diễn (demo) các mẫu thiết kế Microservices trong dự án.

## 🏛️ 1. Sơ đồ Quan hệ (ERD)
- **Farms** (1) <--- (n) **Harvests**
- **Business Logic** <--- (1) **Outbox Events** (Technical Table)

---

## 📊 2. Chi tiết các Bảng (Tables)

### 2.1 Bảng `farms` (Danh mục Nông trại)
- **Vai trò**: Lưu trữ thông tin tài sản vật lý và quyền sở hữu.
- **Mẫu thiết kế**: Row-level Security (Ownership Scoping).

| Cột | Kiểu dữ liệu | Ràng buộc | Mô tả |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PK | Định danh duy nhất (Dùng UUID cho Microservices). |
| `name` | VARCHAR(255) | NOT NULL | Tên nông trại. |
| `location` | TEXT | | Địa chỉ/Vùng trồng. Giới hạn trong danh sách tĩnh (VD: Cầu Đất, Đà Lạt; Buôn Ma Thuột, Đắk Lắk; Pleiku, Gia Lai; Gia Nghĩa, Đắk Nông; Kon Tum). |
| `area` | DECIMAL(10,2) | NOT NULL | Diện tích (Hectares). Bắt buộc > 0. |
| `coffee_type` | VARCHAR(100) | | Giống cà phê chủ đạo. |
| `owner_id` | UUID | NOT NULL | **Khóa an ninh**: Subject ID từ Identity Server. |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Thời gian tạo. |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Thời gian cập nhật. |

---

### 2.2 Bảng `harvests` (Nhật ký Thu hoạch)
- **Vai trò**: Lưu trữ sản lượng đầu vào cho chuỗi cung ứng.
- **Mẫu thiết kế**: **Aggregate Root**. Harvest phụ thuộc hoàn toàn vào sự tồn tại và quyền sở hữu của Farm.

| Cột | Kiểu dữ liệu | Ràng buộc | Mô tả |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PK | Mã lô thu hoạch gốc. |
| `farm_id` | UUID | FK | Link tới `farms.id`. ON DELETE CASCADE. |
| `harvest_date` | DATE | NOT NULL | Ngày hái quả thực tế. |
| `quantity` | DECIMAL(12,2) | NOT NULL | Khối lượng (kg). Bắt buộc > 0. |
| `status` | VARCHAR(50) | DEFAULT 'NEW' | Trạng thái (NEW, PROCESSING, SHIPPED). |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | |

---

### 2.3 Bảng `outbox_events` (Sự kiện Phân tán)
- **Vai trò**: Đảm bảo tính nhất quán dữ liệu khi tích hợp với Kafka (Event-Driven) cho các luồng quan trọng.
- **Mẫu thiết kế**: **Transactional Outbox**. Áp dụng khi có nhu cầu đồng bộ trạng thái nguyên tử. Ghi sự kiện cùng lúc với dữ liệu nghiệp vụ.

| Cột | Kiểu dữ liệu | Ràng buộc | Mô tả |
| :--- | :--- | :--- | :--- |
| `id` | UUID | PK | Định danh sự kiện. |
| `event_type` | VARCHAR(100) | NOT NULL | Loại sự kiện (VD: `farm.harvest.created`). |
| `payload` | JSONB | NOT NULL | Dữ liệu chi tiết của sự kiện. |
| `retry_count` | INT | DEFAULT 0 | Số lần retry. |
| `processed_at` | TIMESTAMPTZ | DEFAULT NULL | Thời gian xử lý. |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Thời gian tạo. |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Thời gian cập nhật. |

---

