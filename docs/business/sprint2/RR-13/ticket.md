# [RR-13] End-to-End Secure Integration

- **Summary:** Kiểm chứng luồng bảo mật toàn diện (Two-Gate Security Model) từ Client → KrakenD → Microservice, bao gồm cả gRPC internal calls và OTel trace correlation.
- **Priority:** `HIGH`
- **Type:** Verification
- **ADR Reference:** `ADR-2026-05-04-AUTH` — Session 8, 10, 12

---

## 📖 User Story
> As a project owner, I want to see the complete Two-Gate Security Model working end-to-end — scope validation at the gateway, role validation at the service, user identity visible in distributed traces — so that I can confirm the system is production-ready for business logic.

## 🔍 Acceptance Criteria

### Scenario 1: Unauthenticated Request Blocked (No Token)
- **Given:** No JWT is attached to the request.
- **When:** The request reaches KrakenD.
- **Then:** KrakenD returns `401 Unauthorized` before forwarding to any service.

### Scenario 2: Gate 1 — App Scope Blocked at Gateway
- **Given:** A JWT is present but the App's `scope` claim does not include the required scope for the endpoint.
- **When:** KrakenD's Native Validator processes the request.
- **Then:** KrakenD returns `403 Forbidden`. The request never reaches the microservice.

### Scenario 3: Gate 2 — Wrong User Role Blocked at Service
- **Given:** A JWT passes Gate 1 (valid scopes) but the user's `role` is insufficient for the action.
- **When:** The microservice's Casbin middleware evaluates the request.
- **Then:** The service returns `403 Forbidden`.

### Scenario 4: Fully Authorized Request Succeeds
- **Given:** A user with the correct role, a valid JWT, and correct App scopes.
- **When:** They call a protected endpoint (e.g., `GET /v1/demo`).
- **Then:** KrakenD passes the JWT through unchanged, the service validates offline, Casbin approves, and the response is `200 OK` with data.

### Scenario 5: Revoked Token Blocked at Service
- **Given:** A user has logged out (token JTI is in Redis Blacklist).
- **When:** They replay the old (still-signature-valid) JWT.
- **Then:** The service's auth middleware detects the blacklisted JTI and returns `401 Unauthorized`.

### Scenario 6: gRPC Internal Calls Secured
- **Given:** Service A makes an internal gRPC call to Service B on behalf of a user.
- **When:** Service A propagates the JWT in gRPC Metadata.
- **Then:** Service B's Security Interceptor validates the token and injects identity into context before the handler runs.
- **And:** Service B does NOT skip auth for internal callers.

### Scenario 7: User Identity Visible in Distributed Traces
- **Given:** A fully authorized request completes successfully.
- **When:** I inspect the trace in SigNoz.
- **Then:** The span for the microservice handler contains a `user.id` attribute matching the authenticated user's `sub` claim.
- **And:** The end-to-end trace spans KrakenD → Service with no gaps.

### Scenario 8: Key Rotation — No Request Dropped
- **Given:** The Identity Server rotates its RS256 signing key.
- **When:** Services are running with stale-while-revalidate JWKS cache.
- **Then:** No in-flight request fails during the rotation window.
- **And:** Within one cache refresh cycle, all services validate with the new key.
