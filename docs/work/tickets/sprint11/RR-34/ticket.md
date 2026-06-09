---
artifact_type: ticket
id: RR-34
status: draft
owner: human
priority: high
lane: high-risk
human_fields:
  - title
  - priority
  - acceptance_criteria
  - scope
  - approval
ai_fields:
  - impacted_areas
  - test_expectations
  - verification_results
  - docs_review
  - context_updates
shared_fields:
  - status
  - trace
  - small_task_exemption
trace:
  backlog_item: BL-004
  requirement: REQ-AUTH-001, NFR-SEC-001
  phase: docs/work/tickets/sprint11/main.md
  detail_design: docs/work/tickets/sprint11/RR-34/technical_design.md
  test_verification: docs/work/tickets/sprint11/RR-34/test_verification.md
  validation_matrix: docs/work/VALIDATION_MATRIX.md
  docs_review: docs/work/tickets/sprint11/RR-34/docs_review.md
  adrs:
    - "planned: timestamp-based token revocation"
  release_notes: docs/releases/CHANGELOG.md
---

# RR-34: Timestamp-Based Token Revocation & Security Audit

## Status

- ID: `RR-34`
- Status: `draft`
- Type: security feature
- Priority: `high`
- Sprint: Sprint 11
- Owner: Tech Lead
- Approval: pending

## Context

### Problem

Runtime Roasters validates self-contained JWTs locally. A signed token remains
usable until `exp` even after logout, role removal, credential compromise, or an
administrator security action. The repository contains a per-`jti` Valkey
blacklist primitive, but it is not wired into service startup and would create
state proportional to the number of revoked tokens.

The required model must not store raw JWTs, JWT hashes, or one row/key for every
access token. Revocation is instead based on comparing the JWT issue time with a
server-owned validity boundary for its subject.

### Scope

- Add user-wide immediate access-token revocation using `iat` and a persisted
  `tokens_valid_after` timestamp.
- Revoke Hydra login and consent sessions for the subject so an old refresh
  grant cannot mint a new access token after the boundary.
- Cache the subject validity boundary in Valkey without making Valkey the source
  of truth.
- Enforce the check consistently in shared HTTP and gRPC authentication.
- Add self-service revoke-all/logout-all and privileged administrator revoke-all
  operations.
- Record security audit events without storing token material.
- Audit JWT validation, revocation authorization, failure handling, sensitive
  logging, and negative security cases in the touched path.

### Out Of Scope

- Revoking exactly one access token by `jti`.
- Per-device or per-session revocation using `sid`.
- Replacing Hydra/Kratos or changing signing algorithms/key management.
- mTLS, rate limiting for unrelated endpoints, and a repository-wide penetration
  test.
- Storing raw access tokens, refresh tokens, ID tokens, or their hashes in the
  Runtime Roasters database or Valkey.

## Acceptance Criteria

- [ ] JWT parsing requires valid `sub`, `iat`, and `exp`; malformed or missing
  required claims are rejected.
- [ ] A protected request is accepted when `iat > tokens_valid_after` and all
  existing signature, issuer, expiry, RBAC, and data-scope checks pass.
- [ ] A protected HTTP or gRPC request is rejected with an unauthenticated result
  when `iat <= tokens_valid_after`.
- [ ] Self-service revoke-all can only revoke the authenticated subject.
- [ ] Administrator revoke-all requires the existing administrative policy and
  accepts a reason; unauthorized callers cannot discover or revoke another
  subject.
- [ ] Repeating revoke-all is idempotent and never moves
  `tokens_valid_after` backwards.
- [ ] Revoke-all invalidates the subject's Hydra login and consent sessions so
  old refresh grants cannot create usable access tokens.
- [ ] PostgreSQL is the authoritative validity-boundary store; Valkey is only a
  bounded cache and no JWT/token string is persisted.
- [ ] Cache miss reads PostgreSQL and repopulates Valkey. Cache/storage errors
  follow the approved fail-closed policy for protected requests.
- [ ] Revocation emits an audit record containing actor, subject, action, reason,
  boundary time, result, and trace ID, but no JWT, refresh token, or secret.
- [ ] Same-second behavior is deterministic: `iat <= tokens_valid_after` is
  rejected and a new login must receive `iat > tokens_valid_after`.
- [ ] Existing non-revoked authentication remains backward compatible after the
  migration and rollout.
- [ ] UAT proves an active token succeeds before revoke, fails immediately after
  revoke, and a fresh login succeeds afterward.

## Flow, Test, And Seed Audit Notes

Reviewed against the current repository on 2026-06-07. The following gaps must
be resolved during implementation; the existing system does not yet satisfy the
proposed RR-34 flow.

### Authentication Flow Mismatches

- Frontend logout currently clears `rr_access_token` first and then starts only
  a Kratos browser logout. RR-34 must call revoke-all while the access token is
  still available, then clear local state and finish browser-cookie cleanup.
- Revoking Hydra login and consent sessions is not sufficient by itself. The
  subject's Kratos sessions must also be deleted so an existing identity session
  cannot silently authenticate again.
- The callback currently discards Hydra's refresh token and stores only the
  access token. Refresh-grant revocation therefore requires a server-side OAuth
  integration harness; browser E2E cannot prove this with the current client.
- The strict rule `iat <= tokens_valid_after` creates an intentional same-second
  boundary. A fresh login may need to wait/retry until Hydra issues an `iat`
  strictly greater than the committed boundary.
- An old cached boundary must not survive a successful revoke response. The
  Valkey projection must be updated monotonically before success is returned;
  cache write failure is a retryable partial failure.

### Test Coverage Gaps

- Existing JWT unit fixtures generally include `exp` but omit `iat`; they must be
  updated before required-claim enforcement is enabled.
- Existing Playwright journey tests frequently use mocked tokens and mocked API
  routes. Those tests do not prove live timestamp revocation.
- Required live proof must include Next.js, KrakenD, auth-service, one protected
  HTTP service, one protected gRPC service, PostgreSQL, Valkey, Hydra, Kratos,
  Kafka, and audit-service.
- Gateway JWT validation alone is insufficient because KrakenD does not know the
  subject's `tokens_valid_after`; service middleware enforcement must be tested.
- Negative tests must cover missing/invalid/future `iat`, exact-boundary time,
  dependency outage, unauthorized cross-user revoke, stale cache, and token-free
  logs/audit payloads.

### Seed And Migration Prerequisites

- Infrastructure seed creates only `admin@runtimeroasters.com` and the Hydra
  OAuth client.
- Role identities such as
  `mgr.hn.hoankiem@runtimeroasters.com` exist only after authenticated
  `Initialize Data` runs.
- Kratos identity IDs are generated. Tests must resolve subjects by seeded email
  instead of hard-coding IDs.
- Auth-service currently has no migration directory or `auth:migrate` task.
  RR-34 must add both and include auth migration in `migrate:all`.
- Recommended UAT actors are the seeded ADMIN for privileged revoke and the
  Hoan Kiem store manager for self-revoke and unauthorized-admin negative cases.

## Small Task Exemption

- Small task exemption: no
- Reason: changes authentication semantics, public/protected APIs, persistence,
  shared middleware, Hydra runtime integration, and failure policy.
- Impact checked: API=yes, DB=yes, Security=yes, Runtime=yes, Standards=no

## Impacted Areas

- Code: `auth-service`, shared JWT claims, HTTP/gRPC auth middleware, service
  bootstrap, Valkey security cache, Hydra client adapter, audit event producer.
- Requirements: clarify immediate user-wide revocation under `REQ-AUTH-001` and
  fail-closed behavior under `NFR-SEC-001`.
- Architecture: authentication flow, security-state ownership, cache boundary,
  Hydra invalidation sequence, and failure behavior.
- API: self revoke-all and administrator subject revoke-all endpoints.
- ERD/data: new auth-owned `user_security_states` table.
- Decisions: ADR required for timestamp-based revocation and fail-closed cache/
  database policy.

## Detail Design

- Required: yes
- Link: [technical_design.md](./technical_design.md)
- Approval: pending

## Test Expectations

- Unit: claim parsing, boundary comparison, monotonic update, authorization,
  cache behavior, failure mapping, and secret-safe audit payloads.
- Integration: PostgreSQL migration/repository, live Valkey cache, Hydra adapter,
  HTTP/gRPC middleware, and concurrency/idempotency.
- E2E: login, protected call, revoke-all, rejected old token, rejected old
  refresh/session, fresh login, successful protected call.
- Security negative cases: missing/invalid `iat`, future `iat`, wrong subject,
  unauthorized admin call, cache poisoning shape, unavailable dependencies, and
  logs containing no token material.
- UAT: required.
- Docs review: requirements, API, ERD, SDD, ADR, context, backlog, validation
  matrix, and release notes.

## Verification Results

- Command: not run; design only
- Result: pending implementation approval
- Notes: no implementation is authorized by this draft.

## Fix/Test Attempt Log

- Same-path failure attempts: 0 / 3
- Total fix/test cycles: 0 / 5
- Blocked by loop guard: no
- Human/design input needed: approve or revise the linked detail design.

## UAT

- Required: yes
- Expected behavior: old credentials stop working immediately after revoke-all;
  a fresh authentication creates usable credentials.
- Verified behavior: pending
- Sign-off: pending

## Completion Checklist

- [ ] Design approved
- [ ] Implementation complete
- [ ] Tests run and recorded
- [ ] Fix/test loop guard respected
- [ ] Validation matrix updated
- [ ] UAT completed
- [ ] Master docs reconciled
- [ ] Docs review completed
- [ ] ADR created
- [ ] `docs/CONTEXT.md` updated
- [x] `docs/work/BACKLOG.md` updated
- [x] Trace links initialized
