# [RR-9] Identity Server Infrastructure

- **Summary:** Thiết lập Ory Kratos làm Identity Provider tập trung, phát hành JWT bằng RS256 và phân phối Public Key qua JWKS endpoint.
- **Priority:** `CRITICAL`
- **Type:** Infrastructure
- **ADR Reference:** `ADR-2026-05-04-AUTH` — Session 2, 4, 12, 13

---

## 📖 User Story
> As a system architect, I want a centralized Identity Server running in OIDC Provider mode so that I can manage users, issue asymmetric JWTs, and distribute public keys via JWKS — enabling all microservices to self-validate tokens without a round-trip to the Identity Server.

## 🔍 Acceptance Criteria

### Scenario 1: Identity Service Deployment
- **Given:** Docker environment is ready.
- **When:** I run `docker compose up identity`.
- **Then:** Ory Kratos starts, connects to `identity_db` (Postgres), and its health endpoint returns `200 OK`.

### Scenario 2: JWKS Endpoint Accessible (Internal Only)
- **Given:** Identity server is running.
- **When:** A microservice calls `GET http://identity:4434/.well-known/jwks.json` with a valid `X-Internal-Secret` header.
- **Then:** It receives a valid JSON Web Key Set containing at least one RS256 public key.
- **And:** Requests without the `X-Internal-Secret` header return `401 Unauthorized`.

### Scenario 3: JWT Issuance with Required Claims
- **Given:** A user authenticates successfully.
- **When:** Kratos issues a JWT.
- **Then:** The token payload contains all required claims: `sub`, `exp`, `iat`, `role`, `org_id`.
- **And:** The token is signed with RS256 (asymmetric).

### Scenario 4: Token Revocation Infrastructure
- **Given:** The auth infrastructure is provisioned.
- **When:** Redis and Kafka are running.
- **Then:** A Distributed Token Blacklist (Redis) is reachable by backend services for logout/lock operations.

## 🛠️ Technical Notes
- **Component:** Ory Kratos in OIDC Provider mode.
- **Database:** Postgres (dedicated `identity_db`).
- **Key Algorithm:** RS256 (Asymmetric Key Pair). Kratos holds Private Key; Public Keys distributed via JWKS.
- **S2S Security:** JWKS endpoint protected by `X-Internal-Secret` Shared Secret Header (Sprint 2 MVP). Upgrade to mTLS in Sprint 4.
- **Token Revocation:** Redis Distributed Blacklist + Kafka event bus wired up in this ticket; revocation logic consumed in RR-11.
- **Key Rotation:** Kratos handles key rotation; consumers implement stale-while-revalidate (RR-11 concern).
