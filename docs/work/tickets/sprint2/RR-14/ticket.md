# [RR-14] Token Revocation & Global Logout

- **Summary:** Triển khai cơ chế thu hồi Token (Blacklist) sử dụng Valkey và lắng nghe Token Revocation events từ Identity Server để đảm bảo người dùng đã logout không thể sử dụng JWT cũ.
- **Priority:** `MEDIUM`
- **Type:** Security Enhancement
- **ADR Reference:** `ADR-2026-05-04-AUTH`

---

## 📖 User Story
> As a security administrator, I want to ensure that once a user logs out or their session is revoked, their existing JWTs are immediately invalidated across all microservices, even if the token has not yet expired.

## 🔍 Acceptance Criteria

### Scenario 1: Token Revocation Check in Middleware
- **Given:** A request arrives with a JWT that has a valid signature but has been revoked.
- **When:** The Auth Middleware checks the token's `jti` against the Valkey Distributed Blacklist.
- **Then:** It returns `401 Unauthorized`.

### Scenario 2: Logout Event Processing
- **Given:** A user performs a logout action.
- **When:** The Identity Server emits a revocation event or the logout handler is called.
- **Then:** The token's `jti` is added to Valkey with a TTL matching the token's remaining expiration time.

### Scenario 3: Global Session Revocation
- **Given:** An administrator revokes all sessions for a compromised account.
- **When:** The revocation events are propagated.
- **Then:** All active JWTs for that account are blacklisted.

## 🛠️ Technical Notes
- **Storage:** Valkey (Distributed Blacklist).
- **Key Pattern:** `blacklist:<jti>`.
- **TTL:** Phải khớp với thời gian còn lại của token (`exp` - `now`).
- **Middleware Integration:** Cập nhật `pkg/base/auth` middleware để thực hiện `EXISTS` check trong Valkey cho mỗi request (phải cực nhanh).
- **Event Source:** Kratos Webhooks hoặc Hydra Revocation events.
