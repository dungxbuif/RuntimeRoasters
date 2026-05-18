# Technology Stack

**Analysis Date:** 2025-02-13

## Languages

**Primary:**
- Go 1.25.0 - Entire backend implementation (`src/`)

**Secondary:**
- Shell - Setup and seed scripts (`deployments/*.sh`, `Makefile`, `Taskfile.yml`)
- Protobuf - API contracts (`api/runtime/`)

## Runtime

**Environment:**
- Docker - Containerized development and deployment (`deployments/docker-compose.dev.yaml`)

**Package Manager:**
- Go Modules - Dependency management (`src/go.mod`)
- Lockfile: `src/go.sum` present

## Frameworks

**Core:**
- Gin v1.12.0 - HTTP Web Framework for REST endpoints (`src/pkg/base/app.go`)
- gRPC v1.80.0 - High-performance RPC framework for service-to-service communication
- gRPC-Gateway v2.29.0 - gRPC-to-JSON proxy for REST compatibility
- GORM v1.31.1 - ORM for database interactions (`src/pkg/database/postgres.go`)

**Testing:**
- Testify v1.11.1 - Assertion and mocking library

**Build/Dev:**
- Taskfile/Task - Task runner for common operations (`Taskfile.yml`)
- Buf - Protocol buffer management and generation (`api/buf.yaml`)

## Core Framework Packages (`src/pkg/`)

- `base/`: Modular service bootstrap for Gin and gRPC (`src/pkg/base/app.go`)
- `config/`: Centralized configuration loading using Viper (`src/pkg/config/config.go`)
- `database/`: GORM-based PostgreSQL client with transaction support (`src/pkg/database/postgres.go`)
- `valkey/`: Valkey/Redis client supporting single-node and Sentinel (`src/pkg/valkey/client.go`)
- `kafka/`: Consumer and Producer wrappers for segmentio/kafka-go (`src/pkg/kafka/`)
- `logger/`: Structured logging with Zap and middleware for Gin/gRPC (`src/pkg/logger/`)
- `telemetry/`: OpenTelemetry setup for tracing (`src/pkg/telemetry/`)
- `errs/`: Standardized error handling and Gin middleware (`src/pkg/errs/`)

## Key Dependencies

**Critical:**
- Casbin v3.10.0 - Authorization library for ABAC/RBAC (`src/pkg/base/casbin/`)
- Ory Kratos/Hydra SDKs - Identity and OAuth2 integration (`github.com/ory/client-go`, `github.com/ory/hydra-client-go/v2`)
- OpenTelemetry (OTEL) - Distributed tracing and observability (`src/pkg/telemetry/`)

**Infrastructure:**
- Viper v1.21.0 - Configuration management (`src/pkg/config/`)
- Zap v1.27.1 - Structured logging (`src/pkg/logger/`)
- Go-Redis v9.18.0 - Client for Valkey/Redis interaction (`src/pkg/valkey/`)
- Kafka-go v0.4.51 - Kafka client for event-driven architecture (`src/pkg/kafka/`)

## Configuration

**Environment:**
- Environment variables managed via `.env` files and `viper`
- `src/pkg/config/` provides base configuration structures

**Build:**
- `src/Makefile` and `Taskfile.yml` for building services

## Platform Requirements

**Development:**
- Go 1.25+
- Docker & Docker Compose
- Task (optional but recommended)
- Buf CLI

**Production:**
- Linux-based containers
- PostgreSQL 16+
- Valkey 7.2+
- Kafka 7.6+

---

*Stack analysis: 2025-02-13*
