# Auth Review

It cannot yet be said that it is "completely unified" according to the initial design. The core has been implemented in the code: JWT claims, `store_ids`, HMAC webhooks, replay protection, and record-level scoping. However, there is still a significant gap in the Casbin matcher, which creates a risk that route-level RBAC for roles like `STORE_MGR` may not function correctly as described in the policy.

### Current Auth Flow

**Gate 1 is JWT validation:**

- HTTP services use `security.NewHTTPGuards(...)`.
- Middleware verifies RS256 JWT via JWKS, issuer `http://localhost:4444/`.
- Claims are injected into the context for the service layer to read.

**Key Code:**
- `src/pkg/base/security/http.go`
- `src/pkg/base/auth/transport/http/middleware.go`
- `src/pkg/base/auth/token/jwt.go`

**Extended Claims:**
- `sub`
- `email`
- `role`
- `org_id`
- `store_ids`
- `jti`

The JWT parser can read both root claims and nested extensions.

### Identity / Hydra / Kratos

**Kratos schema now includes optional `store_ids`:**
- `deployments/kratos/identity.schema.json`

**Auth-service saves `org_id` and `store_ids` into traits when creating a user:**
- `src/apps/auth-service/internal/usecase/user_usecase.go`

**Consent accept route retrieves Kratos traits and includes them in the Hydra token session:**
- `src/apps/client-app/src/lib/ory/identity-admin.ts`
- `src/apps/client-app/src/app/api/auth/consent/accept/route.ts`

This point now matches the design: the "source of truth" for store access is the `store_ids` within the identity/JWT.

### Record-Level Authorization

**Common helper implemented:**
- `src/pkg/base/identity/scope.go`

**Rules:**
- `ADMIN`: Allowed to view everything via the service-level helper.
- `STORE_MGR`: Only allowed to operate on stores within `StoreIDs`.
- `STORE_MGR` without `store_ids`: Fail-closed.
- Record without `store_id`: Only `ADMIN` can view.

**Applied to:**
- **Retail**: list stores, create order, get order.
- **Payment**: list payments, get payment by order.
- **Logistics**: list shipments, deliver shipment, update driver location.
- **Trace**: trace events/documents by `store_id`.
- **Audit**: audit logs by `store_id`.

### Webhook Auth

Webhooks do not use Bearer JWT, aligning with the design for provider callbacks.

**Implemented:**
- **Stripe**: `Stripe-Signature: t=<unix>,v1=<hmac_sha256(secret, t + "." + raw_body)>`
- **VNPay**: `X-VNPAY-Timestamp` + `X-VNPAY-Signature`
- **Timestamp tolerance**: 5 minutes.
- Missing secret/signature/timestamp or incorrect signature will be rejected.
- Replay/idempotency protection via `WebhookEvent` unique constraint `(provider, event_id)`.

**Key Code:**
- `src/apps/payment-service/internal/provider/provider.go`
- `src/apps/payment-service/internal/usecase/service.go`
- `src/apps/payment-service/internal/domain/models.go`

**KrakenD** still keeps webhooks public at the route-auth layer and now forwards the `VNPay` timestamp header:
- `deployments/krakend/krakend.json`

### Critical Casbin Gap

This is the part that is not yet fully unified.

**Policy has been changed to `STORE_MGR`:**
- `src/apps/auth-service/internal/infrastructure/casbin/default_policies.csv`

**However, the current matcher is:**
`m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))`

In the HTTP middleware, `r.sub` is currently `claims.Role` (e.g., `STORE_MGR`), not the user ID. Because the matcher only uses `g(r.sub, p.sub)` without checking `r.sub == p.sub`, a policy like:
`p, STORE_MGR, /v1/stores, read`
may not match if there is no grouping `g, STORE_MGR, STORE_MGR`.

**Impact:**
- `STORE_MGR` route-level RBAC might be denied before reaching the service.
- `ADMIN` is not guaranteed full `/v1/*` access according to policy `p, ADMIN, /v1/*, .*`, because the matcher does not directly compare `ADMIN == ADMIN`.
- `ADMIN` only surely inherits roles declared via `g, ADMIN, STORE_MGR`, `g, ADMIN, FARM_MANAGER`, etc.

**Recommended Fix:**
- Modify the matcher to include direct role policy matching:
`m = r.sub == p.sub || g(r.sub, p.sub) || g(r.sub, "admin")`
while keeping `keyMatch` / `regexMatch` for object/action.

### Test Status

**Passed:**
- JWT claim extraction (including email, `store_ids`, nested extensions).
- Identity scope helper.
- Payment webhook valid/invalid/idempotent scenarios.
- Targeted backend packages for auth/payment/retail/logistics/trace/audit.

**Not fully passed `go test ./...` due to existing out-of-scope errors:**
- `apps/warehouse-service` copy has a space in the import path.
- `scripts` contains multiple `main` functions in the same package.

Frontend lint also fails due to existing errors outside the auth claim route logic.

### Final Assessment

Backend service-level authorization is aligned with the design. Webhook protection is also following a more practical direction than JWT. However, to say Auth is "completely unified," the Casbin matcher must be addressed first, as Gate 2 may currently block correct role policies. Once the matcher is fixed and integration tests are run via KrakenD with a real `STORE_MGR` JWT, the Auth mechanism will reach a unified end-to-end state.
