# master_demo_data.md — Runtime Roasters Master Demo Dataset

Tài liệu này tổng hợp toàn bộ dữ liệu mẫu (Master Data) của hệ thống Runtime Roasters, phục vụ cho mục đích phát triển, kiểm thử (UAT) và trình diễn hệ thống.

---

## 📋 1. Danh sách Người dùng (Identities)

Tất cả mật khẩu mặc định: `Hello@123`

| Tên người dùng | Email (Username) | Role | Mục đích Demo |
| :--- | :--- | :--- | :--- |
| **System Admin** | `admin@runtimeroasters.com` | `ADMIN` | Quản trị toàn hệ thống, tạo user. |
| **Agri Admin** | `agri.admin@runtimeroasters.com` | `FARM_ADMIN` | Quản lý danh mục vùng trồng, gán sở hữu. |
| **Manager Cầu Đất** | `manager.caudat@runtimeroasters.com` | `FARM_MANAGER` | Quản lý vận hành tại vùng Đà Lạt. |
| **Manager Buôn Ma Thuột**| `manager.bmt@runtimeroasters.com` | `FARM_MANAGER` | Quản lý vận hành tại Đắk Lắk. |
| **Farmer K'Ho** | `farmer.kho@runtimeroasters.com` | `FARMER` | Ghi chép nhật ký thu hoạch thực địa. |
| **Roast Master** | `processor@runtimeroasters.com` | `PROCESSOR` | Tiếp nhận cà phê nhân và chế biến. |
| **Driver Alpha** | `driver@runtimeroasters.com` | `DRIVER` | Vận chuyển hàng hóa giữa các điểm. |

---

## ☕ 2. Danh sách Nông trại (Farms - Famous in Vietnam)

Dữ liệu này được thiết kế để khớp với các vùng trồng nổi tiếng thực tế tại Việt Nam.

| Tên Nông trại | Địa điểm (Location) | Diện tích | Giống Cà phê | Chủ sở hữu (Owner) |
| :--- | :--- | :--- | :--- | :--- |
| **K'Ho Coffee Farm** | Lạc Dương, Đà Lạt | 15.5 ha | Arabica Heirloom | `manager.caudat@...` |
| **Cau Dat Arabica** | Cầu Đất, Đà Lạt | 45.0 ha | Arabica Typica | `manager.caudat@...` |
| **Son Pacamara Farm** | Trạm Hành, Đà Lạt | 12.0 ha | Arabica Pacamara | `manager.caudat@...` |
| **Aeroco Coffee** | Ea Kao, Buôn Ma Thuột | 20.0 ha | Specialty Robusta | `manager.bmt@...` |
| **Trung Nguyen Village**| Tân Lợi, Buôn Ma Thuột | 5.0 ha | Robusta | `manager.bmt@...` |
| **Chư Sê Estate** | Chư Sê, Gia Lai | 30.0 ha | Robusta | `agri.admin@...` |

---

## 📦 3. Nhật ký Thu hoạch Mẫu (Harvests)

| Farm | Ngày thu hoạch | Sản lượng (kg) | Trạng thái |
| :--- | :--- | :--- | :--- |
| K'Ho Coffee Farm | 2026-05-10 | 500.00 | `NEW` |
| K'Ho Coffee Farm | 2026-05-11 | 350.00 | `PROCESSING` |
| Aeroco Coffee | 2026-05-09 | 1200.00 | `SHIPPED` |
| Cau Dat Arabica | 2026-05-12 | 800.00 | `NEW` |

---

## 🛠️ 4. Technical JSON (Dành cho Seeding/API Test)

### 4.1 JSON Tạo Identity (Kratos)
```json
{
  "traits": {
    "email": "manager.caudat@runtimeroasters.com",
    "name": "Manager Cầu Đất",
    "role": "FARM_MANAGER"
  },
  "credentials": { "password": { "config": { "password": "Hello@123" } } }
}
```

### 4.2 JSON Tạo Nông trại (Farm Service API)
```json
{
  "name": "K'Ho Coffee Farm",
  "location": "Lạc Dương, Đà Lạt",
  "area": 15.5,
  "coffee_type": "Arabica Heirloom",
  "owner_id": "[ID_CỦA_MANAGER_CAU_DAT]"
}
```

---

## 📖 5. Ghi chú Bảo mật & Vận hành
1. **Ownership**: Nông trại chỉ hiển thị cho người dùng có `id` khớp với `owner_id` (trừ `ADMIN` và `FARM_ADMIN` có quyền view-all).
2. **Role Mapping**: Hệ thống tự động đồng bộ Role từ Kratos sang Casbin qua cơ chế **Bootstrapping Sync** của Auth Service.
3. **Validation**: Khi tạo Farm mới, `area` phải > 0 và `location` nên chọn trong danh sách các tỉnh Tây Nguyên.
