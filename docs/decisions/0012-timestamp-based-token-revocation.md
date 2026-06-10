---
artifact_type: adr
id: ADR-0012
status: accepted
owner: shared
decision_approval: approved_in_chat_2026-06-09
trace:
  requirement: REQ-AUTH-001, NFR-SEC-001
  tickets_or_bugs:
    - RR-34
  master_docs:
    - docs/requirements/REQUIREMENTS.md
    - docs/architecture/API.md
    - docs/architecture/ERD.md
    - docs/architecture/SDD/README.md
---

# ADR 0012: Timestamp-Based Token Revocation

## Context

Runtime Roasters validates signed JWT access tokens locally. A signed token remains usable until `exp` even after a logout, role removal, credential compromise, or admin security action. Maintaining a blacklist for every individual token (`jti`) creates unbounded state proportional to the number of revoked tokens and complicates microservice integration.

## Decision

- Implement token revocation by comparing the token's issue time (`iat`) with a server-controlled subject boundary (`tokens_valid_after`).
- Store the validity boundary authoritatively in PostgreSQL `auth_db` (`user_security_states` table) and cache it in Valkey for performance.
- Reject any token where `token.iat <= subject.tokens_valid_after`.
- Introduce a shared `TokenValidityChecker` client used by HTTP and gRPC auth middlewares/interceptors across all domain services.
- If a subject's token validity boundary is requested and missed in Valkey, domain services call `auth-service` via a private gRPC endpoint to fetch the boundary and update Valkey.
- When revocation is triggered (self-service or admin action), immediately update the PostgreSQL DB and write-through to Valkey, then invalidate the subject's active login and consent sessions in Ory Kratos and Hydra.
- Enforce a fail-closed policy: if PostgreSQL/Valkey lookup fails and validity cannot be verified, reject the protected request.

## Alternatives

- Blacklist individual token IDs (`jti`): rejected because storage grows with every token, and requires complex replication.
- User `token_version` integer: rejected because it requires embedding synchronized version integers inside JWTs, which is harder to coordinate with Ory Hydra custom claims.

## Consequences

- Positive: immediate user-wide revocation with minimal, constant storage per subject.
- Positive: central control over revocation with low latency via Valkey caching.
- Negative: protected requests fail closed if both Valkey and `auth-service` are unavailable.
- Negative: same-second updates require a strict comparison (`iat <= boundary` is invalid), meaning a new login must receive `iat > boundary` to succeed.
