# Codebase Structure

**Analysis Date:** 2025-05-15

## Directory Layout

```
[project-root]/
├── api/                # gRPC/Protobuf contract definitions (Buf)
├── deployments/        # Infrastructure configuration (Docker, SQL, Ory, KrakenD)
├── docs/               # Technical and business documentation
├── src/                # Core Go source code (Monorepo module)
│   ├── apps/           # Microservice implementations
│   ├── pkg/            # Shared libraries and internal framework
│   ├── runtime/        # Generated gRPC code and local service stubs
│   └── scripts/        # Utility scripts (seeders, simulations)
├── go.work             # Go workspace configuration
└── Taskfile.yml        # Project-wide task runner (automation)
```

## Directory Purposes

**api/:**
- Purpose: Source of truth for all service interfaces and event schemas.
- Contains: `.proto` files organized by domain.
- Key files: `api/buf.yaml`, `api/runtime/farm/v1/farm.proto`.

**src/apps/:**
- Purpose: Isolated microservice implementations.
- Contains: Business logic, persistence adapters, and service bootstrap code.
- Key directories: `src/apps/farm-service`, `src/apps/retail-service`.

**src/pkg/:**
- Purpose: Reusable code shared across multiple services.
- Contains: Database helpers, logging, telemetry, and base security logic.
- Key packages: `src/pkg/base` (Auth/Casbin), `src/pkg/telemetry` (OTel), `src/pkg/logger`.

**deployments/:**
- Purpose: Environment setup and infrastructure as code.
- Contains: Docker Compose files, SQL seed scripts, and configuration for edge services (Ory, KrakenD).
- Key files: `deployments/docker-compose.dev.yaml`, `deployments/krakend/krakend.json`.

## Key File Locations

**Entry Points:**
- `src/apps/[service]/cmd/main.go`: Service startup logic.
- `deployments/seed.sh`: Main data seeding script.

**Configuration:**
- `src/apps/[service]/config/config.go`: Service-specific configuration structures.
- `.env.example`: Template for environment variables.

**Core Logic:**
- `src/apps/[service]/internal/usecase/service.go`: Primary business logic implementation.
- `src/apps/[service]/internal/domain/models.go`: Domain entities and interfaces.

**Testing:**
- `src/pkg/**/__tests__/`: Unit tests for shared packages.
- `src/apps/[service]/internal/**_test.go`: Service-level tests.

## Naming Conventions

**Files:**
- Go Source: `snake_case.go` (e.g., `init_db.go`).
- Protobuf: `snake_case.proto` (e.g., `farm_service.proto`).

**Directories:**
- Services: `kebab-case` (e.g., `farm-service`).
- Go Packages: `lowercase` (e.g., `telemetry`).

## Where to Add New Code

**New Microservice:**
1. Create directory in `src/apps/[new-service]`.
2. Define structure: `cmd/`, `config/`, `internal/app`, `internal/domain`, `internal/usecase`.
3. Register in `Taskfile.yml` if needed.

**New API Endpoint:**
1. Define in `api/runtime/[domain]/v1/[service].proto`.
2. Run `buf generate` to update generated code in `src/runtime/`.
3. Implement handler in `src/apps/[service]/internal/app/init.go`.

**New Shared Utility:**
1. Create new package in `src/pkg/[utility-name]`.
2. Ensure it doesn't depend on `src/apps/`.

## Special Directories

**src/runtime/:**
- Purpose: Contains code generated from Protobuf definitions.
- Generated: Yes (via Buf).
- Committed: Yes.

**.planning/:**
- Purpose: GSD-specific planning and codebase mapping documents.
- Generated: Yes.
- Committed: Yes.

---

*Structure analysis: 2025-05-15*
