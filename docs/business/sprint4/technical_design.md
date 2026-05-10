# Technical Design: Sprint 4 — Security Base (Refined)

Mục tiêu: Hoàn thiện nền tảng bảo mật vững chắc, đảm bảo xác thực và phân quyền xuyên suốt (End-to-End) và cơ chế thu hồi quyền hạn tức thì.

---

## 1. [RR-13] End-to-End Auth Integration

### 1.1. Chiến lược Scope đơn giản (Web/Mobile Unified)
- **Design**: Sử dụng một bộ scope duy nhất cho cả Web App và Mobile App để giảm tải cấu hình Gateway.
- **Khai báo**:
    - `farm:app`: Quyền truy cập các tính năng Nông trại.
    - `auth:app`: Quyền truy cập các tính năng định danh/admin.
    - `retail:app`: Quyền truy cập đặt hàng.
- **Gate 1 (KrakenD)**: Chỉ verify sự hiện diện của các scope này trong `scope` claim của JWT.

### 1.2. Bảo mật Liên dịch vụ (gRPC Metadata Propagation)
- **Kỹ thuật**: 
    - Khi Service A gọi Service B: Tự động trích xuất JWT từ incoming context và gắn vào outbound gRPC Metadata.
    - Nếu không có User JWT (System call): Sử dụng `Internal-Secret` đính kèm header `X-Internal-Token`.
- **Implementation**: Viết Interceptor dùng chung trong `pkg/base/auth` để tự động hóa luồng này.

### 1.3. Optional UserID Tracing
- **Design**: Trace-ID là bắt buộc, nhưng `user_id` là **Optional**.
- **Logic**:
    - Nếu JWT hợp lệ -> Trích xuất `sub` và gắn vào OTel span attribute `user.id`.
    - Nếu là System call -> Gắn attribute `user.id = "system"`.
    - Đảm bảo Trace luôn liền mạch trong SigNoz kể cả khi không có User context.

---

## 2. [RR-14] Token Revocation (Logout)

### 2.1. Distributed Blacklist với Redis
- **Vị trí**: Mọi request đi vào Microservice (HTTP/gRPC) đều phải check blacklist.
- **Key**: `blacklist:jti:{jti_id}`.
- **Fail-closed Policy**: Nếu cụm Redis gặp sự cố (Timeout/Conn Refused) -> Middleware PHẢI **Reject** request với lỗi `500 Internal Server Error`. Không cho phép "vượt rào" khi không kiểm tra được trạng thái thu hồi.

### 2.2. Luồng Logout & Password Change
- **Explicit Logout**: Người dùng nhấn nút -> `auth-service` gọi Hydra Revoke + Ghi Redis Blacklist.
- **Password Change (BA Requirement)**: Khi người dùng đổi mật khẩu thành công, `auth-service` sẽ thực hiện:
    1. Thu hồi phiên hiện tại.
    2. (Phase 2) Publish sự kiện `user.security.updated` để các thiết bị khác tự động log out.

### 2.3. Redis Optimization
- Dùng `EXISTS` command để kiểm tra nhanh.
- TTL của key trong Redis = `token_exp - current_time`.

---
*TechLead Signed-off: 2026-05-10*
