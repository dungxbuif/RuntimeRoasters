# [RR-16] [BA] Quản lý Danh mục Nông trại (Create & List)

**User Story:**
Dưới vai trò là **Quản lý Nông nghiệp (Agriculture Manager)**, tôi muốn khai báo thông tin các nông trại trong hệ thống để tôi có thể bắt đầu quy trình theo dõi sản lượng và thu hoạch.

**Business Context:**
Dữ liệu nông trại là thông tin nền tảng. Nếu thông tin vùng trồng, độ cao hoặc loại hạt bị sai lệch, toàn bộ các chứng chỉ chất lượng và truy xuất nguồn gốc phía sau sẽ không còn giá trị.

---

## 🔄 Luồng Nghiệp vụ (Workflow)
1. **Truy cập Quản lý:** User vào mục "Farm Management".
2. **Khai báo Mới:** Nhấn "Add New Farm".
3. **Nhập liệu:** Nhập tên, địa điểm (Enum), diện tích, loại cà phê (Arabica/Robusta).
4. **Lưu trữ:** Hệ thống kiểm tra trùng lặp tên/vị trí và lưu vào Database.

---

## 🛠️ Quy tắc Nghiệp vụ (Business Rules)
- **Vùng địa lý:** Chỉ cho phép chọn các vùng đã định nghĩa trong `domain_enums.md` (Cầu Đất, Buôn Ma Thuột,...).
- **Phân quyền:** Chỉ người dùng có role `FARM_MANAGER` hoặc `ADMIN` mới được tạo nông trại mới theo policy và scope hiện hành.
- **Duy nhất:** Không được phép có 2 nông trại trùng cả Tên và Vùng địa lý.

---

## ✅ Acceptance Criteria (AC)
### Scenario 1: Khai báo nông trại thành công
- **Given:** Tôi nhập đầy đủ thông tin: "Nông trại X", "Cầu Đất", "50 hecta".
- **When:** Tôi nhấn "Save".
- **Then:** Hệ thống lưu dữ liệu và hiển thị nông trại mới trong danh sách.

### Scenario 2: Thiếu thông tin bắt buộc
- **When:** Tôi để trống mục "Location" và nhấn "Save".
- **Then:** Hệ thống báo lỗi "Vùng địa lý là bắt buộc" và không cho lưu.

---

## 📊 Yêu cầu Dữ liệu
- `Name` (Required, String)
- `Location_Code` (Required, Enum)
- `Acreage` (Số hecta)
- `Coffee_Type` (Enum)
