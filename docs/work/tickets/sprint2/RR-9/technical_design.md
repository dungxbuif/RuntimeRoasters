# Dev Notes - [RR-9] Identity Server Infrastructure

## 🛠️ Technical Implementation Details
- **Component:** Ory Kratos (Identity Management) + Ory Hydra (OAuth2/OIDC Provider).
- **Database:** Postgres (dedicated `identity_db` and `hydra_db`).
- **Key Algorithm:** RS256 (Asymmetric Key Pair). Hydra manages Signing Keys; Public Keys distributed via JWKS.
- **S2S Security:** Hydra's JWKS endpoint is protected by `X-Internal-Secret` (Sprint 2 MVP) via an Nginx sidecar.
- **Token Revocation:** Valkey Distributed Blacklist + Kafka event bus wired up in this ticket; revocation logic consumed in RR-11.
- **Key Rotation:** Kratos handles key rotation; consumers implement stale-while-revalidate (RR-11 concern).

## 🏗️ Architecture Design

### 2.1 Components (Docker & Network)
- **`identity_db` & `hydra_db` (Postgres):** Separate databases for Kratos and Hydra.
- **`kratos`:** Manages Identity.
- **`hydra`:** Manages OAuth2 flows.
- **`identity` (Identity Proxy - Nginx Sidecar):** Wraps Hydra's Admin/Public ports to protect JWKS endpoint. Runs on port `4434`.
- **`kafka` & `valkey`:** Included in infrastructure.

### 2.2 JWT & OIDC Configuration
- **Signing Algorithm:** RS256 managed by Hydra.
- **Login & Consent Flow:**
  1. User accesses Client App -> Redirect to Hydra.
  2. Hydra redirects to Login UI (Client App).
  3. Client App calls Kratos to authenticate user.
  4. After success, Client App calls Hydra "Accept Login Request".
  5. Hydra redirects to Consent UI (Client App).
  6. Client App calls Hydra "Accept Consent Request" (with claims like `role`, `org_id`).
  7. Hydra issues JWT with custom claims.

### 2.3 JWKS Endpoint Protection
Nginx (container `identity`) proxies to Hydra's Public port:
- Location `/.well-known/jwks.json`: Checks `X-Internal-Secret` header.

## 🚀 Implementation Plan
- **Step 1:** Update `deployments/init-db.sql` with `CREATE DATABASE identity_db;` and `CREATE DATABASE hydra_db;`.
- **Step 2:** Configure Kratos and Hydra (DSN, Login/Consent Provider URLs, Issuer URL).
- **Step 3:** Setup Nginx Gatekeeper to proxy `hydra:4444` and protect JWKS.
- **Step 4:** Update Docker Compose with `kratos`, `hydra`, `identity` (nginx).
- **Step 5:** Local Validation of JWKS via Proxy with secret header.
