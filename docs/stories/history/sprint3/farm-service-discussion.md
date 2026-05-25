# Sprint 3 Discussion: Farm Service Business Logic (Vertical Slice)

**Trạng thái:** `✅ COMPLETED / SIGNED-OFF`
**Ngày chốt:** 2026-05-10
**Thành phần tham gia:** TechLead, Human Developer (User) & Thư ký (Agent)
**Mục tiêu:** Thống nhất cách triển khai nghiệp vụ (Business Logic) đầu tiên của hệ thống, đảm bảo tính kế thừa và an toàn dữ liệu.

---

## 🕒 Nhật ký thảo luận chi tiết (Detailed Discussion Log)

### Session 1: Mở đầu & Tầm nhìn Sprint 3
- **TechLead:** Nhấn mạnh Sprint 3 là bước ngoặt chuyển từ Infra sang Business. Farm Service sẽ là "bản thiết kế mẫu" cho toàn bộ hệ thống sau này.
- **Vấn đề cốt lõi:** Làm thế nào để cân bằng giữa cấu trúc Clean Architecture chặt chẽ và tốc độ phát triển tính năng (Vertical Slice).

### Session 2: Kiến trúc Vertical Slice
- **Quyết định:** Thống nhất sử dụng **Vertical Slice** bên trong các Layer. Tức là các file sẽ được đặt tên theo tính năng (vd: `create_farm.go`, `list_farms.go`) thay vì gộp chung vào 1 file `farm.go` khổng lồ.
- **Lý do:** Tăng khả năng bảo trì và dễ dàng test độc lập từng luồng nghiệp vụ.

### Session 3: Data Scoping (ABAC) tại Repository
- **Quyết định:** Ép buộc lọc dữ liệu theo **`owner_id`** trực tiếp tại tầng **Repository**.
- **Kỹ thuật:** Sử dụng `identity.FromContext(ctx)` để lấy User ID và đưa vào câu lệnh SQL `WHERE`. Việc này đảm bảo tính an toàn dữ liệu ngay cả khi tầng UseCase quên kiểm tra.

### Session 4: Hoãn triển khai Optimistic Locking
- **Quyết định:** **Loại bỏ** yêu cầu về trường `version` và xử lý xung đột `Optimistic Locking` trong Sprint này.
- **Lý do:** Cần một Sprint riêng để thiết kế mô hình High-level và Global cho toàn bộ hệ thống thay vì làm đơn lẻ tại Farm Service.

### Session 5: Casbin tại Query Level (Khả năng sẵn sàng)
- **TechLead Phân tích:** Dựa trên nguyên tắc **KISS**. Chỉ dùng Casbin Query Level cho logic cực kỳ phức tạp/hay thay đổi. Với các logic đơn giản như Ownership/Tenant, dùng Native GORM `WHERE` để đảm bảo hiệu năng và Type-safety.
- **Quyết định:** Sprint 3 dùng Native GORM Ownership check tại Repository.
- **Yêu cầu mở rộng:** Toàn bộ kiến trúc phải được design để **sẵn sàng (Ready)** tích hợp Casbin Query Level khi có nhu cầu trong tương lai.

### Session 6: Global Casbin Query Mapping (Review hệ thống)
- **TechLead Phân tích:** Casbin Query Level (ABAC-to-SQL) sẽ không dùng tràn lan mà tập trung vào các Service có logic phân quyền dữ liệu (Data Visibility) thay đổi thường xuyên hoặc phụ thuộc nhiều thuộc tính.

**Bản đồ ứng dụng dự kiến:**

1.  **Warehouse Service (Sprint 3+):**
    *   **Flow:** `ListInventory`, `SearchStock`.
    *   **Logic:** Quản kho vùng nào chỉ thấy hàng vùng đó; Quản kho tổng thấy hết.
    *   **Pattern:** **Regional Data Isolation**. Casbin sẽ sinh SQL `WHERE region_id IN (...)`.

2.  **Retail Service (Sprint 4):**
    *   **Flow:** `BrowseProducts`, `GetPricing`.
    *   **Logic:** Khách hàng hạng Gold thấy giá khác hạng Silver; Flash sale chỉ hiển thị cho một nhóm đối tượng.
    *   **Pattern:** **Dynamic Pricing & Visibility**. Casbin sinh SQL lọc sản phẩm theo `price_tier`.

3.  **Trace Service (Sprint 4):**
    *   **Flow:** `ViewAuditLog`, `TraceShipment`.
    *   **Logic:** Chỉ cán bộ thanh tra (Compliance Officer) mới thấy thông tin nhạy cảm về lô hàng.
    *   **Pattern:** **Field-Level Security / Masking**. Casbin sinh điều kiện để `SELECT` các cột tương ứng.

**Structural Pattern ứng dụng:**
- **GORM Scopes:** Tạo một bộ lọc dùng chung `CasbinScope(e *casbin.Enforcer, rvals...)` để nhúng vào mọi câu lệnh query tại Repository: `db.Scopes(CasbinScope(e, sub, obj, act)).Find(&results)`.

---

## 📝 Những thay đổi & Quyết định quan trọng (Decisions)
1. ✅ **Structure:** Sử dụng Vertical Slice trong Clean Architecture.
2. ✅ **Security:** Data Scoping (Ownership check) thực hiện tại Repository layer.
3. 🔴 **Deferred:** Optimistic Locking (Hẹn lại ở Sprint thiết kế global).
4. 🏗️ **Future-proof:** Design hệ thống đảm bảo khả năng mở rộng sang Casbin Query Level.
5. 🗺️ **System Mapping:** Xác định Warehouse và Retail là 2 service chủ đạo sẽ apply Casbin Query Level.

---

## 📅 Action Items cho Project Owner
- [ ] Xem lại toàn bộ kiến trúc design của toàn hệ thống (Global Architecture Review).

---

## 🏛️ Architectural Decision Record (ADR) - Sprint 3 Preview
...
