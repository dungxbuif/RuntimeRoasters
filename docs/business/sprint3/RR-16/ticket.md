# [RR-16] Farm Repository & Database

- **Summary:** Triển khai tầng lưu trữ dữ liệu bền vững cho Nông trại với khả năng truy vấn linh hoạt.
- **Priority:** `HIGH`
- **Type:** Feature

---

## 📖 User Story
> As a farm owner, I want my farm data to be persisted reliably so that I can access historical records and search for specific farms based on various criteria.

## 💰 Business Value
Đảm bảo tính toàn vẹn và sẵn sàng của dữ liệu. Khả năng tìm kiếm giúp người dùng quản lý số lượng lớn nông trại một cách hiệu quả.

## 🔍 Acceptance Criteria

### Scenario 1: Persistent Storage
- **Given:** A farm creation request.
- **When:** The system processes the request.
- **Then:** The data must be stored in a relational database (PostgreSQL).

### Scenario 2: Data Retrieval
- **Given:** An existing farm ID.
- **When:** Searching by ID.
- **Then:** The system returns full details of that farm.

### Scenario 3: Filtering & Search
- **Given:** A list of farms.
- **When:** Filtering by Region or Farm Type.
- **Then:** The system returns only the matching records.

### Scenario 4: Audit Trails
- **Given:** Any modification to farm data.
- **When:** The record is saved.
- **Then:** `CreatedAt` and `UpdatedAt` timestamps must be automatically managed.
