# Identity & Frontend Deployment Checklist (RR-9, RR-10, RR-11)

Tài liệu này tổng hợp các điểm cần rà soát và cấu hình trước khi triển khai hệ thống Identity (Kratos, Hydra) và Frontend (client-app) lên môi trường Production.

## 1. Các thay đổi đã Refactor (Centralized Config)
Các giá trị này đã được chuyển vào hệ thống biến môi trường, cần cung cấp giá trị đúng khi deploy:
- [x] **`NEXT_PUBLIC_JAEGER_URL`**: URL của Jaeger UI (vd: `http://localhost:16686`).
- [x] **`NEXT_PUBLIC_SWAGGER_JSON_URL`**: URL của file swagger.json (vd: `http://localhost:8081/swagger/demo.swagger.json`).
- [x] **`src/constants/env.ts`**: Đã hỗ trợ lấy cấu hình từ `process.env`.
- [x] **`.env.example`**: Đã cập nhật đầy đủ các biến môi trường cần thiết tại root và `src/apps/client-app`.

## 2. Rà soát Hardcoded (Nguy cơ bảo mật & Sai lệch logic)
Cần xử lý các điểm này trong code trước khi lên Prod:
- [ ] **OAuth2 Claims (`api/auth/consent/accept/route.ts`)**: `role` và `org_id` đang bị hardcode là `admin` và `org-root-001`. Cần logic để map động từ identity traits của user.
- [ ] **Hydra Issuer URL**: Trong config của Hydra, `URLS_SELF_ISSUER` phải khớp với domain public (vd: `https://auth.runtimeroasters.com/`). Nếu sai, Backend sẽ reject JWT do không khớp `iss` claim.
- [ ] **Internal Secret**: `INTERNAL_SECRET` đang dùng giá trị mặc định. Đây là lớp bảo vệ JWKS endpoint, cần được thay đổi và bảo mật tuyệt đối.
- [ ] **OIDC Client Secret**: `OIDC_CLIENT_SECRET` đang dùng giá trị demo (`client-secret`).

## 3. Cấu hình Môi trường (Environment Variables)
Đảm bảo các biến sau được cấu hình chính xác trong hệ thống CI/CD hoặc Orchestrator:

### Frontend (Client App)
- `NEXT_PUBLIC_APP_URL`: Domain của frontend.
- `NEXT_PUBLIC_HYDRA_PUBLIC_URL`: Domain public của Hydra.
- `HYDRA_ADMIN_URL`: URL nội bộ trỏ tới Hydra Admin (không public).
- `INTERNAL_SECRET`: Secret dùng để fetch JWKS qua Identity Proxy.
- `OIDC_CLIENT_SECRET`: Secret của OAuth2 Client.

### Backend (Go Services)
- `EXPECTED_ISSUER`: Phải khớp chính xác với `iss` claim trong JWT (vd: `http://localhost:4444/`).
- `JWKS_URL`: URL nội bộ trỏ tới JWKS endpoint (vd: `http://identity:4434/.well-known/jwks.json`).
- `INTERNAL_SECRET`: Dùng cho header `X-Internal-Secret` khi gọi Identity Proxy.

## 4. Network & Security
- [ ] **CORS Settings**: Cấu hình Hydra/Kratos để cho phép các domain Prod truy cập.
- [ ] **Identity Proxy (Nginx)**: Đảm bảo Nginx chặn mọi truy cập tới JWKS nếu thiếu hoặc sai `X-Internal-Secret`.
- [ ] **Token Storage**: Cân nhắc chuyển từ LocalStorage sang HttpOnly Cookie để chống XSS nếu yêu cầu bảo mật cao hơn.

---
*Cập nhật lần cuối: 2026-05-06 bởi Gemini CLI*
