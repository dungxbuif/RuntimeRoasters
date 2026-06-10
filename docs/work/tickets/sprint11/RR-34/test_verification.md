---
artifact_type: test_verification
id: VERIFY-RR-34
status: verified
owner: ai
trace:
  ticket: docs/work/tickets/sprint11/RR-34/ticket.md
  detail_design: docs/work/tickets/sprint11/RR-34/technical_design.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
---

# RR-34 Test Verification Plan

This document details the test strategy, specific test cases, and verification commands required to prove the correctness of the timestamp-based token revocation implementation.

## 1. Automated Test Plan

We have implemented unit and integration coverage for the JWT validation logic and interceptor/middleware validity checker checks.

### 1.1 Base Auth & Claims Unit Tests
* **Verification Command**: 
  ```bash
  GOCACHE=src/.go-cache GOTMPDIR=src/.go-tmp go test -v ./pkg/base/auth/__tests__
  ```
* **Test Cases Covered**:
  - `successful verification`: Ensures a valid JWT containing `sub`, `iss`, `iat`, and `exp` is parsed and returns populated claims.
  - `missing iat`: Ensures JWTs without the `iat` claim are rejected immediately.
  - `future iat beyond skew allowance`: Ensures JWTs with an `iat` timestamp more than 5 minutes in the future are rejected.
  - `missing sub`: Ensures JWTs without a `sub` claim are rejected.
  - `revoked token via TokenValidityChecker is rejected`: Verifies that gRPC unary interceptors reject tokens when the checker indicates the subject's token is invalid/revoked.
  - `valid token via TokenValidityChecker is accepted`: Verifies that gRPC unary interceptors allow requests when the checker validates the token's issue time.

### 1.2 Auth Service Integration Tests
* **Verification Command**:
  ```bash
  GOCACHE=src/.go-cache GOTMPDIR=src/.go-tmp go test -v ./apps/auth-service/...
  ```
* **Test Cases Covered**:
  - PostgreSQL schema and migration verification.
  - Repository `MonotonicUpdate` upserts and DB timestamp reads.
  - Valkey cache write-through and read-through validation.
  - Singleflight lookup coalescing on cache misses.

---

## 2. Manual Verification & UAT Plan

To prove E2E functionality across KrakenD, auth-service, Ory Hydra, Kratos, and the domain services:

### 2.1 Prerequisite: Seed & Bootstrap
1. Clean and reset the local docker stack:
   ```bash
   make docker-down && make docker-up
   ```
2. Run migrations:
   ```bash
   task auth:migrate
   ```
3. Initialize and seed system data (via Admin UI or seed task).

### 2.2 UAT Scenario 1: Self-Service Logout & Token Revocation
1. **Login**: Login as the Hoan Kiem store manager (`mgr.hn.hoankiem@runtimeroasters.com`) to receive an access token.
2. **Access Protected API**: Call a protected endpoint (e.g. `GET /v1/retail/sales`) and verify it returns `200 OK`.
3. **Logout/Revoke**: Trigger self-service logout:
   ```bash
   curl -X POST http://localhost:8081/v1/auth/logout \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{"reason": "user_logout"}'
   ```
4. **Verify Revocation**: Immediately call the protected endpoint `GET /v1/retail/sales` again using the same access token. Verify it returns `401 Unauthorized` / `Unauthenticated`.
5. **Verify Refresh Token Invalidation**: Attempt to use the OAuth refresh token to mint a new access token. Verify it is rejected because Hydra consent/login sessions were invalidated.
6. **Fresh Login**: Perform a fresh login. Verify the new access token is accepted.

### 2.3 UAT Scenario 2: Administrative Revoke (Privileged Action)
1. **Login as Admin**: Login as `admin@runtimeroasters.com` to receive an admin access token.
2. **Login as Target User**: Login as the Dalat farm manager to get a target user access token.
3. **Verify Dalat Access**: Verify the Dalat farm manager token can access Dalat farm endpoints successfully (`200 OK`).
4. **Trigger Admin Revoke**: Call the administrative revoke endpoint:
   ```bash
   curl -X POST http://localhost:8081/v1/admin/users/<dalat_subject_id>/security/revoke-all \
     -H "Authorization: Bearer <admin_access_token>" \
     -H "Content-Type: application/json" \
     -d '{"reason": "compromised_credentials"}'
   ```
5. **Verify Dalat Revoked**: Immediately call Dalat farm endpoints using the Dalat manager's token. Verify it is rejected with `401 Unauthorized`.
6. **Verify Non-Admin Access Denied**: Attempt to call the admin revoke endpoint using the Hoan Kiem manager's token. Verify it returns `403 Forbidden`.

---

## 3. Negative Security cases & Outage Plan

* **Valkey Outage**: Stop Valkey (`docker stop valkey`). Perform a protected request. Verify that the request still succeeds (falls back to PostgreSQL) but logs a degraded cache warning.
* **Database & Cache Outage**: Stop both Valkey and PostgreSQL. Perform a protected request. Verify it **fails closed** immediately returning a `500 Internal` error rather than allowing access.
