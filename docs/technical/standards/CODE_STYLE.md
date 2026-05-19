# Internal Engineering Standard: Canonical Template

This document defines the code standards that every Microservice in RuntimeRoasters must adhere to, using `demo-service` as the standard template.

## 1. Error Handling (gRPC & REST)
Raw errors must never be returned from the UseCase layer to the Transport layer.

- **Technique:** Use `pkg/errs`.
- **Implementation (Handler):** 
    ```go
    if err != nil {
        return nil, errs.ToGRPCError(err) // Automatically maps ErrNotFound -> codes.NotFound
    }
    ```

## 2. Context-Aware Logging
Every log line must include a TraceID for Observability.

- **Technique:** Use `logger.FromContext(ctx)`.
- **Implementation:**
    ```go
    log := logger.FromContext(ctx)
    log.Info("Doing something...", zap.String("key", value))
    ```

## 3. Domain Validation
Data validation logic must reside within the Entity (Domain layer).

- **Technique:** Implement a `Validate() error` method for the Domain Struct.
- **Implementation (UseCase):**
    ```go
    if err := entity.Validate(); err != nil {
        return nil, err // errs.ErrValidation
    }
    ```

## 4. Database Transactions (Unit of Work)
When a UseCase needs to write data to multiple tables or records, a Transaction must be used.

- **Technique:** Use `db.WithTx(ctx, func(tx *gorm.DB) error { ... })`.
- **Implementation:** Ensure the Repository can receive `*gorm.DB` from external sources or use a Transaction Decorator pattern.

## 5. Dependency Propagation
The `context.Context` must be propagated throughout from Handler -> UseCase -> Repository. This is mandatory to maintain the OTel Trace chain and extract Identity (User ID).
