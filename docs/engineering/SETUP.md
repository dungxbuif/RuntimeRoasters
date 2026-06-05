---
artifact_type: engineering_setup
id: SETUP
status: active
owner: shared
updated: 2026-06-05
---

# Local Setup

## Prerequisites

- Go workspace support via `go.work`.
- Node/npm for `src/apps/client-app`.
- Docker for local infrastructure.
- Taskfile for common repo commands.

## First Run

1. Start local infrastructure from repo scripts or Taskfile.
2. Run migrations for all services.
3. Open the client app at `http://localhost:3000`.
4. Log in as `admin@runtimeroasters.com`.
5. If the bootstrap modal appears, initialize the DB.

## Demo Reset

Use the reset script for a clean local demo:

```bash
./deployments/reset-demo-state.sh
```

After reset, log in as admin and run the bootstrap flow if prompted.

## Deterministic Demo Accounts

All accounts use the local demo password from `.env`. The exact canonical emails, roles, assignments, enum values, and seed counts are defined in `docs/requirements/MASTER_DATA.md`.

## Notes

- Historical dashboard/read-model data may be backend-seeded for demo readiness.
- Live demo flows should still use real service APIs/events unless explicitly labeled as simulation in `docs/requirements/REQUIREMENTS.md`.
- `auth-service` owns the admin bootstrap trigger. The first-run UI calls the system seed endpoint, `auth-service` resolves Kratos IDs by email, then downstream services seed their own records idempotently.
