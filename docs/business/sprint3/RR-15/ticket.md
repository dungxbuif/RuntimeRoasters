# [RR-15] Farm Service Bootstrap & Domain

- **Summary:** Thiết lập dịch vụ Farm mới và định nghĩa các thực thể cốt lõi cho việc quản lý nông trại.
- **Priority:** `HIGH`
- **Type:** Feature

---

## 📖 User Story
> As a farm owner, I want a dedicated system to manage my farm data so that I can track production and resources effectively across multiple locations.

## 💰 Business Value
Cung cấp khả năng quản lý tập trung thông tin nông trại, tạo tiền đề cho việc truy xuất nguồn gốc (traceability) và tối ưu hóa vận hành trong các giai đoạn sau.

## 🔍 Acceptance Criteria

### Scenario 1: Farm Entity Definition
- **Given:** A need to store farm information.
- **When:** A farm record is created or viewed.
- **Then:** It must contain: Name, Location (Province/Region), Total Area (Hectares), and Farm Type (e.g., Arabica, Robusta).

### Scenario 2: Service Isolation
- **Given:** The system architecture.
- **When:** Farm management actions are performed.
- **Then:** They must be handled by a dedicated "Farm Service" to ensure scalability and independent deployment.

### Scenario 3: Initial Data Seed
- **Given:** A new installation of the Farm Service.
- **When:** The service starts for the first time.
- **Then:** It should optionally support loading initial reference data for regions and farm types.
