# Technical Design: Modular Bootstrap Pattern (Reusability)

Goal: Eliminate redundancy (Boilerplate) in `apps/*/internal/app/app.go` and `cmd/main.go` by moving orchestration logic into the `pkg/base` package.

## 1. Current Problem
Currently, whenever a new service is created (such as the Farm Service), we have to copy-paste approximately 80% of the code from `demo-service/internal/app/app.go`. The repeating parts include:
- Initializing Casbin background sync.
- Registering Auth Middleware for HTTP.
- Setting up Readiness checks (ping DB/Valkey).
- Configuring Swagger/Gateway.
- Graceful Shutdown logic.

## 2. Solution: Orchestrator Pattern (Facade)
We will upgrade `pkg/base` so that it not only provides "registration" functions but also serves as a complete **Orchestrator**.

### Selected Design Pattern: **Template Method (via Composition) + Functional Hooks**.

### Proposed Structure in `pkg/base`:

#### A. Definition of `Dependencies` (Common Components)
Consolidate components used by all services into a centralized struct within `base`:
```go
type Dependencies struct {
    DB           *database.DB
    RDB          *valkeyclient.Client
    CasbinEngine casbin.Engine
    KeyProvider  provider.KeyProvider
}
```

#### B. Definition of `ServiceRegistrar` (Interface)
Each service only needs to provide its own identification information and handlers:
```go
type ServiceRegistrar interface {
    RegisterGRPC(server *grpc.Server)
    RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
    RegisterHTTP(router *gin.RouterGroup) // Optional
}
```

#### C. Centralized Bootstrap Function (`Launcher`)
This function will execute the entire "standard" flow that we are currently performing manually in `app.go`:
```go
func (a *App) Launch(deps Dependencies, registrar ServiceRegistrar) {
    // 1. Automatically run Casbin Sync if an Engine exists
    if deps.CasbinEngine != nil {
        if reader, ok := deps.CasbinEngine.(*casbin.ResilientReader); ok {
            reader.StartBackgroundSync(context.Background())
        }
    }

    // 2. Call the service-specific registration logic
    registrar.RegisterGRPC(a.grpcServer)
    a.RegisterGateway(registrar.RegisterGateway, a.config.GRPCPort)

    // 3. Automatically inject Auth Middleware into all requests via Gateway (if needed)
    // Or provide a group that is already wrapped with Auth for the Service
    
    // 4. Automatically register Readiness checks for DB & Valkey
    a.RegisterReadiness(func() error {
        if deps.DB != nil { /* ping */ }
        if deps.RDB != nil { /* ping */ }
        return nil
    })

    // 5. Run server
    a.Run(...)
}
```

## 3. Benefits After Implementation

### Code in `farm-service/internal/app/app.go` will be reduced to:
```go
func (a *App) Run() {
    a.Base.Launch(a.Deps, a.Handler) // Just a single line!
}
```

### Code in `cmd/main.go` (Using Manual DI):
Virtually unchanged, but the initialization logic will be cleaner as the passed parameters are encapsulated within `Dependencies`. All dependencies are initialized and passed manually at the Composition Root.

## 4. Implementation Roadmap (Post-Sprint 3)
1. **Refactor `pkg/base`**: Add the `Dependencies` struct and the `Launch` function.
2. **Standardize `Config`**: Ensure every service configuration is compatible so that `base` can read common parameters.
3. **Migration**: Convert `demo-service` to use `base.Launch` for verification. Subsequently, apply it to `farm-service`.

## 💡 Philosophy
"Base is the backbone, Service is the flesh." The backbone handles survival functions (Security, Health, Sync), while the flesh handles business functions (Business Logic).
