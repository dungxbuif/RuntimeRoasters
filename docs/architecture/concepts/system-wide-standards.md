# System-Wide Architectural Standards: Runtime Roasters

Tài liệu này định nghĩa các tiêu chuẩn kỹ thuật, thư viện dùng chung và các mẫu thiết kế áp dụng xuyên suốt cho toàn bộ hệ thống microservices.

---

## 1. Tiêu chuẩn Thư viện Dùng chung (`pkg/`)

Các thư viện trong thư mục `src/pkg/` là nền tảng giúp giảm thiểu code lặp và đảm bảo tính đồng nhất.

| Package | Vai trò | Cơ chế / Công nghệ |
| :--- | :--- | :--- |
| `pkg/base/auth` | Authentication | **Three-Gate Model**, gRPC Metadata Propagator. |
| `pkg/base/casbin` | Authorization | **Distributed Enforcement**, gRPC-based Snapshot. |
| `pkg/database` | Persistence | **GORM Wrapper**, Unit of Work (TxManager). |
| `pkg/kafka` | Messaging | **CloudEvents Standard**, Sync/Async Producer. |
| `pkg/errs` | Error Handling | **RFC 9457 (Problem Details)**. |
| `pkg/valkey` | Caching | **Valkey 7.2**, High-performance Redis replacement. |

---

## 2. Các Mô hình Thiết kế Quan trọng (Core Patterns)

### 2.1. Internal gRPC-Only Mandate
Chúng ta áp dụng chính sách **Zero Internal HTTP**.
- Mọi giao tiếp giữa các service phải sử dụng gRPC.
- Gateway (KrakenD) thực hiện chuyển đổi HTTP sang gRPC.
- Lợi ích: Type-safety tuyệt đối, hiệu năng cao, giảm bề mặt tấn công.

### 2.2. Three-Gate Authorization
Hệ thống bảo vệ dữ liệu qua 3 lớp (Gate):
1.  **Gate 1 (Identity):** Kiểm tra JWT tại Gateway/Interceptor.
2.  **Gate 2 (Method-level):** Casbin kiểm tra quyền gọi RPC (VD: `can call ListFarms`).
3.  **Gate 3 (Data-level):** Sử dụng `GormScoper` tại tầng Repository để tự động thêm điều kiện lọc dữ liệu (VD: `WHERE owner_id = ?`) dựa trên Policy.

### 2.3. Selective Transactional Outbox
Chỉ áp dụng cho các luồng nghiệp vụ cần tính nhất quán cao (VD: Thu hoạch, Đặt hàng). Ghi dữ liệu nghiệp vụ và sự kiện vào DB trong cùng một Transaction để đảm bảo Kafka message không bao giờ bị mất.

---

## 3. Quy trình Tài liệu API (Swagger Flow)

Chúng ta không sử dụng tài liệu build từ Protobuf cho các đối tác bên ngoài.
- **Source of Truth:** Cấu hình Gateway (`krakend.json`).
- **Quy trình:** 
    1. Dev cập nhật endpoint trong KrakenD.
    2. Chạy `task swagger-export` để xuất OpenAPI spec từ KrakenD.
    3. Frontend tự động hiển thị tài liệu mới tại `/api-docs`.
- **Lý do:** Đảm bảo tài liệu phản ánh chính xác những gì Gateway đang thực sự mở ra.

---

## 4. Tiêu chuẩn Khởi tạo (Manual DI)

Chúng ta không sử dụng DI Framework (như Wire). Việc khởi tạo được thực hiện thủ công tại **Composition Root**:
- Đảm bảo tính minh bạch: "Mọi thứ được nối dây ở đâu?".
- Cấu trúc khởi tạo tập trung tại `internal/app/init.go`.

---
**Ký duyệt:** TechLead
*Cập nhật lần cuối: 2026-05-13*
