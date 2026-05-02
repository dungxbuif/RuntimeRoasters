# [RR-9] Identity Server Infrastructure

- **Summary:** Thiết lập hạ tầng quản lý danh tính (Identity Provider) tập trung.
- **Priority:** `CRITICAL`
- **Type:** Infrastructure

---

## 📖 User Story
> As a system architect, I want a centralized Identity Server so that I can manage users, issue JWTs, and provide a secure source of truth for authentication across all microservices.

## 🔍 Acceptance Criteria
### Scenario 1: Identity Service Deployment
- **Given:** Docker environment is ready.
- **When:** I run `docker compose up identity`.
- **Then:** Identity server starts successfully and connects to the `identity_db`.

### Scenario 2: JWKS Accessibility
- **Given:** Identity server is running.
- **When:** I request `GET /.well-known/jwks.json`.
- **Then:** I receive a valid JSON Web Key Set containing the public keys needed for JWT verification.

## 🛠️ Technical Notes
- Component: Ory Kratos (recommended for lightweight identity) or a custom Go Identity service.
- Database: Postgres (separate schema or DB).
- JWT Claims required: `sub`, `exp`, `iat`, `role`, `org_id`.
