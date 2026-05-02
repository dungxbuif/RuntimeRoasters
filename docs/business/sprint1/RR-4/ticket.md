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

**Story 5 — Swagger UI (Client App Integration)**
> As a developer, I want the Swagger contract to be accessible through the Gateway, so that the Client App can load and display the API Explorer.

---

## 🔍 Acceptance Criteria (BDD Specification)

### Scenario 1: Proto compilation
- **Given:** A `.proto` file exists with at least one service definition.
- **When:** I run `make proto`.
- **Then:** Go stubs and a `swagger.json` file are generated with no errors.

### Scenario 2: Gateway routes to service
- **Given:** The API Gateway and Demo Service are running.
- **When:** I send `GET /v1/demo/ping` to the Gateway.
- **Then:** The request reaches the Demo Service gRPC handler and returns a valid JSON response.

### Scenario 3: Swagger accessibility via Gateway
- **Given:** Demo Service is running and exposing its `swagger.json`.
- **When:** I call `GET localhost:8081/swagger/demo.swagger.json` via Gateway.
- **Then:** I receive the valid OpenAPI JSON content.

### Scenario 4: Dependencies are injected automatically
- **Given:** A new dependency is added to any layer (e.g. a repository).
- **When:** I run `wire gen`.
- **Then:** `wire_gen.go` is updated and `main.go` requires no manual changes.

---

## 🛠️ Technical Notes
- Proto toolchain: `buf`.
- HTTP-to-gRPC transcoding: `grpc-gateway`.
- Dependency Injection: `google/wire`.
- API Gateway: `KrakenD`.
- Client App: Next.js + Swagger UI React.

## 📋 Sub-tickets

| Ticket | Summary | Status |
| :--- | :--- | :--- |
| [RR-4-1](./RR-4-1/ticket.md) | Proto Toolchain — buf setup & demo.proto | ✅ Done |
| [RR-4-2](./RR-4-2/ticket.md) | `pkg/base` — RegisterGateway & ServeSwagger | ✅ Done |
| [RR-4-3](./RR-4-3/ticket.md) | Dependency Injection — Google Wire scaffold | ✅ Done |
| [RR-4-4](./RR-4-4/ticket.md) | Infrastructure — KrakenD & Client App Shell | 🚧 In Progress |
| [RR-4-5](./RR-4-5/ticket.md) | Demo Service — GetDemo handler | ✅ Done |
