# [RR-4] Infrastructure Completion — Proto Toolchain, DI & API Gateway

- **Summary:** Hoàn thiện hạ tầng kỹ thuật nền tảng để các service có thể giao tiếp qua Gateway và expose API có tài liệu.
- **Priority:** `HIGH`
- **Type:** Infrastructure

---

## 📖 User Stories

**Story 1 — Proto Toolchain**
> As a developer, I want a single command to compile `.proto` files into Go code and Swagger docs, so that the API contract is always in sync with the implementation.

**Story 2 — HTTP-to-gRPC Transcoding**
> As a developer, I want the base framework to support mounting an HTTP gateway alongside the gRPC server, so that any service can expose REST endpoints without running a separate process.

**Story 3 — Dependency Injection**
> As a developer, I want dependencies to be wired automatically using Google Wire, so that `main.go` remains clean and each layer is independently testable.

**Story 4 — API Gateway**
> As a developer, I want an API Gateway in the local environment that routes external HTTP traffic to the correct internal service, so that the architecture reflects production from day one.

**Story 5 — Swagger UI**
> As a developer, I want a Swagger UI available in the local environment that loads the API contract from the running service, so that I can test endpoints without writing curl commands.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Proto compilation
- **Given:** A `.proto` file exists with at least one service definition.
- **When:** I run `make proto`.
- **Then:** Go stubs and a `swagger.json` file are generated with no errors.

### Scenario 2: Gateway routes to service
- **Given:** The API Gateway and Farm Service are running.
- **When:** I send `GET /v1/farms/demo` to the Gateway.
- **Then:** The request reaches the Farm Service gRPC handler and returns a valid JSON response.

### Scenario 3: Swagger UI reflects the contract
- **Given:** Farm Service is running and exposing its `swagger.json`.
- **When:** I open Swagger UI in the browser.
- **Then:** The endpoint `GET /v1/farms/demo` is visible and executable.

### Scenario 4: Dependencies are injected automatically
- **Given:** A new dependency is added to any layer (e.g. a repository).
- **When:** I run `wire gen`.
- **Then:** `wire_gen.go` is updated and `main.go` requires no manual changes.

---

## 🛠️ Technical Notes
- Proto toolchain: `buf` (lint, breaking change detection, codegen).
- HTTP-to-gRPC transcoding: `grpc-gateway`.
- Dependency Injection: `google/wire`.
- API Gateway: `KrakenD`.
- Swagger UI: `swaggerapi/swagger-ui` Docker image.

## 📋 Sub-tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-4-1](./RR-4-1/ticket.md) | Proto Toolchain — buf setup & farm.proto | 🕒 To Do |
| [RR-4-2](./RR-4-2/ticket.md) | `pkg/base` — RegisterGateway & ServeSwagger | 🕒 To Do |
| [RR-4-3](./RR-4-3/ticket.md) | Dependency Injection — Google Wire scaffold | 🕒 To Do |
| [RR-4-4](./RR-4-4/ticket.md) | Infrastructure — KrakenD & Swagger UI | 🕒 To Do |
| [RR-4-5](./RR-4-5/ticket.md) | Farm Service — GetDemoFarm handler | 🕒 To Do |
