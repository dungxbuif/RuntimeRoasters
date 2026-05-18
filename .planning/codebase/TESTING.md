# Testing Patterns

**Analysis Date:** 2025-02-13

## Test Framework

**Runner:**
- Go: Standard `go test`
- Frontend: Playwright (`@playwright/test`)

**Assertion Library:**
- `testify/assert` and `testify/require` for Go.

**Run Commands:**
```bash
go test ./...              # Run all backend tests
cd src/apps/farm-service && go test ./... # Run tests for specific service
npm run test:e2e           # Run frontend E2E tests (Playwright)
```

## Test File Organization

**Location:**
- Backend: Co-located in `__tests__` subdirectories within the package they test (e.g., `internal/usecase/__tests__/`).
- Frontend: `tests/` or co-located with components.

**Naming:**
- Go: `[file]_test.go`
- TypeScript: `[file].spec.ts` or `[file].test.ts`

**Structure:**
```
src/apps/[service]/
├── internal/
│   ├── usecase/
│   │   ├── farm_usecase.go
│   │   └── __tests__/
│   │       └── farm_usecase_test.go
│   └── infrastructure/
│       └── repository/
│           ├── farm_repository.go
│           └── tests/
│               └── farm_repository_test.go
└── __tests__/
    └── e2e/
        └── security_test.go
```

## Test Structure

**Suite Organization (Go):**
```go
func TestCreateFarm(t *testing.T) {
    // Setup
    repo := new(MockRepository)
    u := usecase.NewFarmUsecase(repo)
    ctx := context.Background()

    // Mocking
    repo.On("Create", mock.Anything, mock.Anything).Return(nil)

    // Execute
    result, err := u.CreateFarm(ctx, farmReq)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    repo.AssertExpectations(t)
}
```

**Patterns:**
- **Setup pattern:** Initialize mocks and inject them into the service/usecase.
- **Teardown pattern:** Not frequently needed in unit tests, but used in `testify/suite` for E2E.
- **Assertion pattern:** Use `testify/assert` for non-fatal checks and `testify/require` for fatal ones.

## Mocking

**Framework:** `testify/mock`

**Patterns:**
```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, farm *domain.Farm) error {
    args := m.Called(ctx, farm)
    return args.Error(0)
}
```

**What to Mock:**
- Databases (Repositories)
- External Service Clients (Kratos, Hydra, etc.)
- Message Brokers (Kafka)

**What NOT to Mock:**
- Domain Entities
- Internal helper functions or pure logic

## Fixtures and Factories

**Test Data:**
```go
farmReq := &domain.Farm{
    Name:       "Test Farm",
    Location:   "CAU_DAT",
    Area:       10.5,
    CoffeeType: "ARABICA",
}
```

**Location:**
- Usually defined within the test file or a `shared_test.go` in the same `__tests__` directory.

## Coverage

**Requirements:**
- Target: >80% code coverage for core business logic (UseCases).

**View Coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Types

**Unit Tests:**
- Focus: Business logic in `usecase` layer.
- Isolation: Mocks all external dependencies.
- Location: `internal/usecase/__tests__/`.

**Integration Tests:**
- Focus: Database queries and repository logic.
- Isolation: Uses real (test) database, often via Docker/Testcontainers.
- Location: `internal/infrastructure/repository/tests/`.

**E2E Tests:**
- Focus: Full request-response cycle, security gates, and service-to-service communication.
- Setup: Requires full infrastructure (Postgres, Kafka, KrakenD, Kratos) to be running.
- Location: `src/apps/[service]/__tests__/e2e/`.
- Frontend: Playwright tests in `src/apps/client-app/`.

## Common Patterns

**Async Testing:**
- Use channels or `Eventually` from Gomega (though not widely used here, standard `testify` patterns are preferred).

**Error Testing:**
```go
_, err := u.CreateFarm(ctx, invalidReq)
assert.ErrorIs(t, err, errs.ErrValidation)
```

---

*Testing analysis: 2025-02-13*
