# Technical Design - [RR-16] Flow 1: Create & List Farm

Mục tiêu: Triển khai luồng tạo mới và hiển thị danh sách nông trại của người dùng.

## 🏗️ 1. Domain & Repository
- **File:** `internal/domain/farm.go`
- **Entity:**
    ```go
    type Farm struct {
        ID        string    `gorm:"primaryKey;type:uuid"`
        Name      string    `gorm:"not null"`
        Location  string    
        Area      float64   `gorm:"not null"`
        Type      string    
        OwnerID   string    `gorm:"not null;index"`
        CreatedAt time.Time
        UpdatedAt time.Time
    }
    ```
- **Interface:**
    ```go
    type FarmRepository interface {
        Create(ctx context.Context, farm *Farm) error
        ListByOwner(ctx context.Context, ownerID string) ([]*Farm, error)
    }
    ```

## 🛠️ 2. Infrastructure (Postgres)
- **File:** `internal/infrastructure/postgres/farm_repo.go` (Sử dụng Vertical Slice: có thể chia nhỏ thành `create.go`, `list.go`)
- **Logic:**
    - Hàm `ListByOwner`: Phải thực hiện lọc `db.Where("owner_id = ?", ownerID)`.

## 🧠 3. UseCase (Business Logic)
- **File:** `internal/usecase/create_farm.go` và `list_farms.go`
- **Logic `CreateFarm`:**
    1. Trích xuất `OwnerID` từ Context: `claims, _ := identity.FromContext(ctx)`.
    2. Kiểm tra `Area > 0`. Trả về `errs.ErrValidationError` nếu Area <= 0.
    3. Gán `farm.OwnerID = claims.Subject`.
    4. Gọi Repo `Create`.

## 📡 4. Delivery (gRPC)
- **File:** `internal/delivery/grpc/handler.go`
- **Hàm `CreateFarm`:** Map `CreateFarmRequest` sang Domain Entity, gọi UseCase, trả về `CreateFarmResponse`.
- **Hàm `ListFarms`:** Gọi UseCase, map slice entities sang slice Protobuf messages.

## 📋 Sub-tasks
- [ ] Định nghĩa Struct Entity trong Domain.
- [ ] Implement Repository với GORM.
- [ ] Viết UseCase `CreateFarm` và `ListFarms`.
- [ ] Hoàn thiện gRPC Handler.
