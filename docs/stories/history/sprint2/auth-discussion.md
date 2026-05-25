# Security Design Discussion: Authorization Strategy (Comprehensive Log)

**Trạng thái:** `DRAFT / ĐANG THẢO LUẬN`
**Thành phần tham gia:** Techlead & User
**Mục tiêu:** Xây dựng hệ thống bảo mật Microservices chuẩn Production-grade cho Runtime Roasters.

---

## 🕒 Nhật ký thảo luận chi tiết (Detailed Discussion Log)

### Session 1: Chiến lược Phân quyền & Vị trí thực thi
- **Vấn đề:** Nên phân quyền tập trung tại Gateway hay phân tán tại từng Service?
- **Phân tích:** Nếu tập trung sẽ dẫn đến `Fat Gateway`, vi phạm tính **Isolation** (Độc lập). Nếu phân tán sẽ giúp services tự chủ nhưng khó quản lý chính sách.
- **Quyết định nháp:** Chọn **Decentralized Authorization** (Phân quyền phân tán). Mỗi service là một pháo đài tự bảo vệ mình theo nguyên lý **Zero Trust**.

### Session 2: Offline Validation & Asymmetric Keys
- **Giải pháp:** Sử dụng cặp **Asymmetric Key Pairs** (Public/Private Key) để thực hiện **Offline Validation**.
- **Cơ chế:** Identity Server dùng Private Key để ký JWT. Các Microservices dùng Public Key để tự verify token mà không cần gọi lại Identity Server, giúp giảm **Latency** và tránh **Bottleneck**.
- **JWKS (JSON Web Key Set):** Sử dụng endpoint `/.well-known/jwks.json` để phân phối Public Key tự động.

### Session 3: Identity Propagation trong Code Go
- **Vấn đề:** Truyền `user_id` và `role` vào UseCase/Repository sao cho sạch.
- **Quyết định:** Sử dụng `context.Context` kết hợp với **Type-safe Wrapper**.
- **Kỹ thuật:** Tạo package `pkg/base/identity` với **Private Context Key** để đảm bảo dữ liệu không bị ghi đè trái phép và code UseCase không bị "bẩn" bởi các `magic string`.

### Session 4: Service-to-Service Authentication
- **Vấn đề:** Bảo mật endpoint JWKS để tránh **Key Spoofing** (giả mạo Public Key).
- **Phương án:** Cân nhắc giữa **Network Isolation** (mạng nội bộ Docker) và **mTLS (Mutual TLS)**.
- **Lộ trình:** Sprint 2 dùng **Shared Secret Header** cho đơn giản; Sprint 4 sẽ nâng cấp lên **mTLS** để đạt chuẩn bảo mật cao nhất cho giao tiếp nội bộ.

### Session 5: Hybrid Authorization Model (RBAC + ABAC)
- **Vấn đề:** Casbin mạnh về **RBAC** (Role-based) nhưng khó xử lý **ABAC** (Attribute-based - ví dụ: quyền sở hữu dữ liệu).
- **Giải pháp:** Kết hợp **Middleware check Role** (Casbin) và **Repository check Ownership** (SQL filtering: `WHERE user_id = ?`). Đảm bảo cân bằng giữa hiệu năng và tính bảo trì.

### Session 6: Định nghĩa Roles & Least Privilege
- **Danh sách Roles:** `admin`, `farm_manager`, `processor`, `warehouse_mgr`, `driver`, `store_mgr`, `guest`.
- **Nguyên tắc:** Áp dụng **Least Privilege** (Quyền hạn tối thiểu). Tách biệt `processor` và `warehouse_mgr` dù logic có vẻ tương đồng để đảm bảo tính chuyên môn hóa của từng domain Microservice.

### Session 7: Advanced Casbin - Hierarchy & Data Scopes
- **Hierarchy:** Triển khai **Role Hierarchy** bằng `g(sub, role)` trong Casbin (ví dụ: Admin thừa kế quyền của tất cả các Role khác).
- **Data Scopes:** Sử dụng Casbin để định nghĩa các phạm vi dữ liệu (ví dụ: "vùng: Dak Lak") thay vì chỉ trả về Đúng/Sai. Repository sẽ lấy các "Scoping Claims" này để tự động generate SQL.

### Session 8: Two-Gate Security Model (Scopes vs Roles)
- **Gate 1 (Edge Validation):** Gateway (KrakenD) kiểm tra **Scopes** (quyền của App/Client). Nếu App không có scope cần thiết, request bị chặn ngay lập tức.
- **Gate 2 (Service Validation):** Microservice kiểm tra **Roles** (quyền của User) thông qua Casbin.
- **Ý nghĩa:** Bảo mật đa lớp, ngăn chặn rò rỉ dữ liệu kể cả khi User có quyền Admin nhưng dùng một ứng dụng không tin tưởng.

### Session 9: Khả năng mở rộng của Gateway (KrakenD)
- **Lựa chọn:** Ưu tiên **Native Validator** (Cấu hình JSON) để tối ưu hiệu năng.
- **Custom Logic:** Chỉ dùng **LUA Scripts** cho các logic đặc biệt. Tuyệt đối tránh **Go Plugins** tại Gateway do rủi ro về tương thích phiên bản (versioning mismatch).

### Session 10: Pure Distributed Authentication
- **Kiến trúc:** Để Gateway **Pass-through** nguyên vẹn JWT gốc từ Client xuống backend.
- **Hành động:** Từng service tự fetch JWKS và verify Token. 
- **Vai trò KrakenD:** Đóng vai trò **Traffic Coordinator**, lo các nhiệm vụ **Cross-cutting concerns** như **Rate Limiting**, **CORS**, và **OTel Tracing Injection**.

### Session 11: So chiếu KrakenD vs Library-based Gateway
- **Phân tích:** KrakenD tuy là binary engine nhưng hoàn toàn đáp ứng được tính linh hoạt nhờ cơ chế pipeline. Hiệu năng (throughput) của nó vượt trội hơn các giải pháp Gateway viết bằng code library truyền thống.

### Session 12: Production Hardening & Resilience
- **Token Revocation:** Đề xuất dùng **Distributed Blacklist (Valkey)** kết hợp Kafka để vô hiệu hóa Token ngay lập tức khi Logout/Lock.
- **Key Rotation:** Áp dụng cơ chế **Stale-while-revalidate** khi cache JWKS để hệ thống không bị gián đoạn khi Identity Server xoay vòng Key.
- **Security Interceptors:** Mọi cuộc gọi **gRPC Internal** phải được bảo vệ bằng Interceptors kiểm tra JWT trong Metadata.

### Session 13: Kratos & Hydra Integration (OIDC Compliance)
- **Vấn đề:** Ory Kratos tập trung vào Identity Management nhưng cần **Ory Hydra** để xử lý các luồng OAuth2/OIDC chuẩn chỉ và phát hành JWT (Access Token).
- **Giải pháp:** Sử dụng Kratos làm Identity Provider (IdP) và Hydra làm OAuth2 Provider. Client App sẽ handle Login/Consent UI.
- **Casbin-SQL:** Chốt hướng "Dùng Casbin triệt để" bằng cách dịch Policy thành SQL `WHERE` clause thông qua một bộ Helper tập trung.

---

## 🔑 Key Terms & Keywords
---

## 🏛️ Architectural Decision Record (ADR) - Sprint 2 Sign-off

**Mã quyết định:** `ADR-2026-05-04-AUTH`
**Chủ đề:** Thống nhất mô hình bảo mật Microservices cho Sprint 2.

### 1. Quyết định về Kiến trúc (Architecture)
- **Mô hình:** **Pure Distributed Authentication** kết hợp **Decentralized Authorization**.
- **Luồng đi:** Gateway (Pass-through JWT) -> Microservice (Self-validation).
- **Lý do:** Tuân thủ nguyên tắc **Zero Trust** và **Isolation**. Loại bỏ điểm nghẽn tập trung tại Gateway.

### 2. Quyết định về AuthN (Authentication)
- **Tiêu chuẩn:** JWT (Asymmetric RS256).
- **Phân phối Key:** Tự động qua **JWKS endpoint** từ Ory Hydra.
- **Xác thực:** Offline Validation tại từng service để tối ưu hiệu năng.

### 3. Quyết định về AuthZ (Authorization)
- **Mô hình 2 lớp (Two-Gate Security):**
    - **Lớp 1 (Gateway):** Check **Scopes** (App permissions) bằng Native Validator của KrakenD.
    - **Lớp 2 (Service):** Check **Roles** (User permissions) bằng Casbin Engine.
- **Data Security:** Sử dụng Casbin để trả về **Scoping Claims**, Repository thực hiện filter dữ liệu bằng SQL injection (`WHERE` clause).

### 4. Quyết định về Identity Propagation
- **Thực thi:** Sử dụng `context.Context` kèm **Type-safe Wrapper** (`pkg/base/identity`).
- **Bảo mật:** Dùng Private Key cho context để ngăn chặn giả mạo dữ liệu User trong nội bộ code Go.

### 5. Quyết định về Infrastructure & Resilience
- **Identity Server:** Ory Kratos (Identity Management) + Ory Hydra (OAuth2/OIDC Provider).
- **API Gateway:** KrakenD (Go-based, configuration-driven).
- **Thu hồi Token:** Sử dụng **Distributed Blacklist** (Valkey + Kafka).
- **S2S Security:** Shared Secret Header (MVP) và lộ trình nâng cấp lên mTLS (Sprint 4).

**== CHỐT PHƯƠNG ÁN & CHUYỂN SANG GIAI ĐOẠN THỰC THI ==**
