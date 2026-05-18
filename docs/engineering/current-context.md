# Runtime Roasters Current Context

Last updated: 2026-05-15

This note captures implementation facts discovered during local Sprint 1-5 testing. Keep it aligned with `AGENT.md` when architecture or operating decisions change.

## Current Decisions

- `FARMER` is removed from the product role model. Farm-side operations use `FARM_MANAGER`; broad farm administration uses `FARM_ADMIN`; system-wide administration uses `ADMIN`.
- Demo/local identity data should not seed or document Farmer accounts. Current local users are `admin@runtimeroasters.com` (`ADMIN`) and `manager.caudat@runtimeroasters.com` (`FARM_MANAGER`).
- Farm-service authorization is a distributed-reader model:
  - `auth-service` is the centralized Casbin policy writer and owner of `auth_db.casbin_rule`.
  - farm-service bootstraps policies via gRPC `AuthService.GetFullSnapshot`.
  - farm-service listens to Kafka topic `auth.policy.changed` for live refresh.
  - polling is retained only as a fallback safety net.
- Architecture Topology is the root page (`/`) and is also reused at `/dashboard/topology-mesh`. It uses dashboard-native node cards and status styling. Do not restore the old `HighFidelityTopology` diagram for either route.

## Bugs Found And Fixed

- `POST /v1/auth/login/accept` returned HTTP `404` with a success body because Gin `NoRoute` had already set status before grpc-gateway wrote the response. `pkg/base.App.FinalizeRoutes` now sets status `200` before delegating `/v1` paths to the gateway mux.
- After OIDC callback, the app landed on `/?access_token=...` and did not enter the dashboard. The home page now stores the access token, strips it from the URL, refreshes auth state, and redirects to `/dashboard/users`.
- Logout followed by login could loop on Kratos `session_already_available`. The login page now restarts the Hydra authorization flow when Kratos already has a browser session but the app has no access token.
- `GET /v1/users` returned gRPC `Unimplemented` because auth-service did not implement `ListUsers` and `CreateUser` handlers. Both RPCs now map to the existing user use case.
- Farm authorization returned `403` because old policies used stale method names such as `/farm.v1.FarmService/*`; generated gRPC full methods are `/runtime.farm.v1.FarmService/...`.
- Farm-service originally treated every gRPC method as `write`, denying read calls such as `ListFarms`. The Casbin gRPC interceptor now maps `Get/List` to `read`, `Delete` to `delete`, and other calls to `write`.
- Farm repository scoping previously checked user subject/grouping instead of JWT role for global role permissions. `GormScoper` now evaluates the role while still scoping non-admin results by subject ownership.
- Frontend called unsupported `GET /v1/harvests`. Harvest listing now fans out through supported `GET /v1/farms/{id}/harvests`.
- farm-service crashed if auth-service was not ready during initial policy snapshot fetch. It now starts the resilient reader background bootstrap instead of failing process startup.
- The root page still rendered the old Topography/Flow control screen after the new topology was built. `/` now renders `ArchitectureTopology` while preserving the OIDC callback token capture and redirect behavior.

## Operational Notes

- Local service URLs:
  - client-app: `http://localhost:3000`
  - KrakenD: `http://localhost:8081`
  - auth-service HTTP/gRPC: `8082` / `50052`
  - farm-service HTTP/gRPC: `8083` / `50053`
  - Kratos public/admin through identity proxy: `4433` / `4434`
  - Hydra public/admin: `4444` / `4445`
  - Kafka UI: `http://localhost:8090`
- Auth policy change topic: `auth.policy.changed`.
- User event topic remains `auth.user.events`.
- Farm harvest event topic remains `farm.harvest.events`.

## Verification Snapshot

- `go test ./pkg/base/casbin/... ./apps/auth-service/internal/usecase ./apps/farm-service/internal/...`
- `npm run lint` in `src/apps/client-app`
- `GET http://127.0.0.1:8081/v1/users` should return only `ADMIN` and `FARM_MANAGER` demo users.
- farm-service startup log should include `Auth policies synchronized from snapshot` and `Consumer listening on topic: auth.policy.changed`.
