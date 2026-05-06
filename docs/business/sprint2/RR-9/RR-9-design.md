# [RR-9] Identity Server Infrastructure (Kratos + Hydra) - Design & Analysis

Tài liệu này cung cấp bản phân tích yêu cầu chi tiết và thiết kế kỹ thuật cho Ticket **RR-9** (Task đầu tiên của Sprint 2), mở rộng với việc tích hợp Ory Hydra.

## 1. Phân tích Yêu cầu (Requirement Analysis)
Dựa theo yêu cầu mới nhất, hệ thống sẽ sử dụng bộ đôi **Ory Kratos** và **Ory Hydra** để đạt được khả năng phát hành JWT chuẩn OAuth2/OIDC.

### Các mục tiêu cốt lõi:
- **Identity Provider (IdP):** Ory Kratos quản lý người dùng, đăng ký, đăng nhập và hồ sơ (Traits).
- **OAuth2/OIDC Provider:** Ory Hydra quản lý các luồng OAuth2, phát hành Access Token (JWT) và ID Token.
- **Decentralized Authorization (Zero Trust):**
  - Hydra phát hành JWT ký bằng **RS256** (Asymmetric key).
  - Microservices tự xác thực Token offline thông qua Public Key lấy từ Endpoint JWKS của Hydra.
- **Custom Claims:** Bắt buộc Token payload phải chứa các thuộc tính bổ sung: `sub`, `exp`, `iat`, `role`, `org_id`.
- **Internal Security:** Endpoint JWKS của Hydra (`/.well-known/jwks.json`) phải được bảo vệ bởi `X-Internal-Secret` cho truy cập nội bộ.
- **Token Revocation Infrastructure:** Chuẩn bị sẵn Kafka và Redis phục vụ cho nghiệp vụ Distributed Blacklist ở task RR-11.

## 2. Thiết kế Kiến trúc (Architecture Design)

### 2.1 Các thành phần triển khai (Docker & Network)
- **`identity_db` & `hydra_db` (Postgres):** Hai database riêng biệt cho Kratos và Hydra.
- **`kratos`:** Quản lý Identity.
- **`hydra`:** Quản lý OAuth2 flows. Cần kết nối với Kratos thông qua một Login/Consent Provider (thường được tích hợp trong Client App hoặc một service riêng).
- **`identity` (Identity Proxy - Nginx Sidecar):** Bọc port Admin/Public của Hydra để bảo vệ JWKS endpoint. Chạy trên port `4434`.
- **`kafka` & `redis`:** Bổ sung Kafka vào hạ tầng.

### 2.2 JWT & OIDC Configuration
- **Signing Algorithm:** RS256 quản lý bởi Hydra.
- **Login & Consent Flow:**
  1. User truy cập Client App -> Redirect tới Hydra.
  2. Hydra redirect tới Login UI (Client App).
  3. Client App gọi Kratos để authenticate user.
  4. Sau khi login thành công, Client App gọi Hydra "Accept Login Request".
  5. Hydra redirect tới Consent UI (Client App).
  6. Client App gọi Hydra "Accept Consent Request" (kèm theo các claims như `role`, `org_id`).
  7. Hydra phát hành JWT chứa các custom claims.

### 2.3 Cơ chế Bảo vệ JWKS Endpoint
Nginx (container `identity`) sẽ proxy tới port Public của Hydra:
- Location `/.well-known/jwks.json`: Kiểm tra `X-Internal-Secret` header.

## 3. Kế hoạch Thực thi (Implementation Plan)

- **Bước 1: Cập nhật file Database Init (`deployments/init-db.sql`)**
  - Thêm `CREATE DATABASE identity_db;` và `CREATE DATABASE hydra_db;`.
- **Bước 2: Cấu hình Ory Kratos & Hydra**
  - Cấu hình Kratos như cũ.
  - Cấu hình Hydra: DSN, URL của Login/Consent Provider, Issuer URL.
- **Bước 3: Thiết lập Nginx Gatekeeper**
  - Proxy tới `hydra:4444` (Public) và bảo vệ JWKS.
- **Bước 4: Cập nhật Docker Compose**
  - Thêm `kratos`, `hydra`, `identity` (nginx).
- **Bước 5: Local Validation**
  - Verify JWKS từ Hydra qua Proxy với secret header.
