# [RR-11] Backend Security Core - JWT Validation

- **Summary:** Xây dựng cơ chế xác thực JWT phân tán tại tầng `pkg/base`, bao gồm JWKS caching, Identity propagation qua Context, và Token Revocation check.
- **Priority:** `HIGH`
- **Type:** Infrastructure
- **ADR Reference:** `ADR-2026-05-04-AUTH` — Session 2, 3, 4, 10, 12

---

## 📖 User Story
> As a developer, I want my microservice to automatically verify incoming JWTs using cached public keys and propagate the caller's identity into the request context — so that I never write repetitive auth logic per endpoint, and UseCase/Repository code receives typed identity claims cleanly.

## 🔍 Acceptance Criteria

### Scenario 1: JWKS Bootstrapping with Shared Secret
- **Given:** A microservice starts up.
- **When:** It initializes the Auth Middleware.
- **Then:** It fetches JWKS from `http://identity:4434/.well-known/jwks.json` using the `X-Internal-Secret` header.
- **And:** Public keys are cached in memory.

### Scenario 2: JWKS Cache with Stale-While-Revalidate
- **Given:** The identity server rotates its signing key.
- **When:** The cached JWKS becomes stale.
- **Then:** The service continues validating with the stale key while asynchronously refreshing in the background.
- **And:** No request is dropped during the key rotation window.

### Scenario 3: Valid Token — Identity Propagated to Context
- **Given:** A request arrives with a valid RS256-signed JWT.
- **When:** The auth middleware processes it.
- **Then:** It extracts claims (`sub`, `role`, `org_id`) and stores them in the request context via `pkg/base/identity` type-safe wrapper (private context key — no magic strings).
- **And:** UseCase code can retrieve identity via `identity.FromContext(ctx)` without touching HTTP headers.

### Scenario 4: Unauthorized — Invalid or Expired Token
- **Given:** A request arrives with an expired or tampered JWT.
- **When:** The middleware validates the signature.
- **Then:** It returns `401 Unauthorized` with a Problem Details JSON body (RFC 9457).

### Scenario 5: gRPC Internal — JWT in Metadata
- **Given:** Service A calls Service B via gRPC.
- **When:** Service A propagates the caller's JWT.
- **Then:** Service B's gRPC Security Interceptor validates the JWT from the gRPC Metadata.
- **And:** Identity is injected into the context before the handler runs.
