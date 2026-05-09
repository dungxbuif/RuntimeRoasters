# [RR-10] Client-Side Auth & Login Flow

- **Summary:** Xây dựng luồng đăng nhập qua Ory Kratos, quản lý JWT bằng HttpOnly Cookie, và bảo vệ route với Next.js Middleware.
- **Priority:** `HIGH`
- **Type:** Feature
- **ADR Reference:** `ADR-2026-05-04-AUTH` — Session 8, 9, 13

---

## 📖 User Story
> As a user, I want to log in with my credentials through the Identity Server so that I receive a JWT stored securely, and every subsequent API request is automatically authorized through KrakenD's scope validation.

## 🔍 Acceptance Criteria

### Scenario 1: Login Page & Form
- **Given:** I navigate to `/login`.
- **Then:** I see a login form styled according to the Design System (Industrial Style).

### Scenario 2: Successful Login & Secure Token Storage
- **Given:** I enter valid credentials and submit.
- **When:** The Identity Server (Kratos) issues a JWT.
- **Then:** The token is stored in an `HttpOnly Cookie` (not localStorage).
- **And:** I am redirected to `/` (dashboard).

### Scenario 3: App Scope Validation at Gateway (Gate 1)
- **Given:** I am logged in and my Client App holds a valid JWT with required scopes.
- **When:** I call any protected API endpoint via KrakenD.
- **Then:** KrakenD's Native Validator checks the App's `scopes` claim.
- **And:** If the scope is missing, the request is rejected at the gateway with `403 Forbidden` before reaching any microservice.

### Scenario 4: Authenticated Request with Automatic Token Injection
- **Given:** I am logged in.
- **When:** I perform any action (e.g., Trigger Ping on the Control Plane).
- **Then:** The Axios interceptor automatically attaches the `Authorization: Bearer <token>` header.
- **And:** KrakenD passes the raw JWT through to the backend service unchanged (Pass-through mode).

### Scenario 5: Route Protection for Unauthenticated Users
- **Given:** I am not logged in.
- **When:** I navigate to any protected route (e.g., `/dashboard`).
- **Then:** Next.js Middleware intercepts the request and redirects me to `/login`.

### Scenario 6: Logout
- **Given:** I am logged in.
- **When:** I click logout.
- **Then:** The HttpOnly Cookie is cleared and I am redirected to `/login`.
