# Internal Engineering Standard: Canonical Template

Tài liệu này định nghĩa các tiêu chuẩn code mà mọi Microservice trong RuntimeRoasters phải tuân thủ, lấy `demo-service` làm bản mẫu chuẩn.

## 1. Error Handling (gRPC & REST)
Tuyệt đối không trả về lỗi thô (raw errors) từ tầng UseCase ra Transport.

- **Kỹ thuật:** Sử dụng `pkg/errs`.
- **Thực thi (Handler):** 
    ```go
    if err != nil {
        return nil, errs.ToGRPCError(err) // Tự động map ErrNotFound -> codes.NotFound
    }
    ```

## 2. Context-Aware Logging
Mọi dòng log phải đính kèm TraceID để có thể quan sát (Observability).

- **Kỹ thuật:** Sử dụng `logger.FromContext(ctx)`.
- **Thực thi:**
    ```go
    log := logger.FromContext(ctx)
    log.Info("Doing something...", zap.String("key", value))
    ```

## 3. Domain Validation
Logic kiểm tra tính hợp lệ của dữ liệu phải nằm tại Entity (Domain layer).

- **Kỹ thuật:** Implement phương thức `Validate() error` cho Domain Struct.
- **Thực thi (UseCase):**
    ```go
    if err := entity.Validate(); err != nil {
        return nil, err // errs.ErrValidation
    }
    ```

## 4. Database Transactions (Unit of Work)
Khi một UseCase cần ghi dữ liệu vào nhiều bảng hoặc nhiều bản ghi, phải sử dụng Transaction.

- **Kỹ thuật:** Sử dụng `db.WithTx(ctx, func(tx *gorm.DB) error { ... })`.
- **Thực thi:** Đảm bảo Repository có thể nhận `*gorm.DB` từ bên ngoài hoặc sử dụng mẫu Transaction Decorator.

## 5. Dependency Propagation
Ép buộc truyền `context.Context` xuyên suốt từ Handler -> UseCase -> Repository. Việc này là bắt buộc để duy trì chuỗi Trace OTel và trích xuất Identity (User ID).
