# Testing Patterns

**Analysis Date:** 2025-05-14

## Test Framework

**Runner:**
- Standard Go `testing` package.
- `testify` for assertions and mocking.

**Assertion Library:**
- `github.com/stretchr/testify/assert`
- `github.com/stretchr/testify/suite`

**Run Commands:**
```bash
go test ./...              # Run all tests
go test -v ./...           # Verbose mode
go test -cover ./...       # Coverage
make test                  # Shortcut via Makefile
```

## Test File Organization

**Location:**
- Unit tests are often in a `tests` or `__tests__` subdirectory within the component directory:
  - `src/apps/farm-service/internal/usecase/tests/`
  - `src/apps/farm-service/internal/infrastructure/repository/tests/`
- Integration and E2E tests are located in:
  - `src/apps/farm-service/__tests__/`
  - `src/apps/farm-service/__tests__/e2e/`

**Naming:**
- Standard Go pattern: `*_test.go`.

**Structure:**
```
[component]/
├── internal/
│   ├── usecase/
│   │   ├── farm_usecase.go
│   │   └── tests/
│   │       └── farm_usecase_test.go
├── __tests__/
│   ├── farm_test.go
│   └── e2e/
│       └── security_test.go
```

## Test Structure

**Suite Organization:**
```go
// Using testify suite for E2E tests
type SecurityTestSuite struct {
	suite.Suite
	// ... context
}

func (s *SecurityTestSuite) SetupSuite() { /* ... */ }
func (s *SecurityTestSuite) Test_Case() { /* ... */ }

func TestSecuritySuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}
```

**Patterns:**
- **Setup/Teardown:** `SetupSuite` and `TearDownSuite` in `testify/suite`.
- **Identity Mocking:** Injecting `identity.Claims` into `context.Context` using `identity.InjectContext`.
- **Table-driven tests:** Often used for unit tests with multiple edge cases.

## Mocking

**Framework:** `github.com/stretchr/testify/mock`

**Patterns:**
```go
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, farm *domain.Farm) error {
	args := m.Called(ctx, farm)
	return args.Error(0)
}

// In test:
repo := new(MockRepository)
repo.On("Create", ctx, mock.Anything).Return(nil)
```

**What to Mock:**
- External dependencies (Repositories, Event Publishers, External APIs).
- Identity/Authentication in context.

**What NOT to Mock:**
- Domain entities and their internal logic (`Validate()` methods).
- Pure utility functions.

## Fixtures and Factories

**Test Data:**
- Domain objects are manually instantiated in tests:
```go
farmReq := &domain.Farm{
    Name:       "Test Farm",
    Location:   "CAU_DAT",
    Area:       10.5,
    CoffeeType: "ARABICA",
}
```

**Location:**
- Inline in test files or shared in a `tests` package if used across multiple files.

## Coverage

**Requirements:** None explicitly enforced in the current configuration, but standard coverage reporting is available.

**View Coverage:**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Types

**Unit Tests:**
- Test individual components (Usecases, Entities) in isolation.
- Dependencies are mocked using `testify/mock`.
- Files located in `internal/.../tests/`.

**Integration Tests:**
- Test components with real dependencies (e.g., PostgreSQL, Casbin enforcer).
- Databases are often cleaned up and migrated per test run: `src/apps/farm-service/internal/infrastructure/repository/tests/farm_repository_test.go`.

**E2E Tests:**
- Test full flows through the API Gateway (KrakenD).
- Verify cross-cutting concerns like security (Two-Gate AuthZ).
- Files located in `__tests__/e2e/`.

## Common Patterns

**Async Testing:**
- Use `go func()` with channels or `waitgroups` for testing background processes like the Outbox Relay.

**Error Testing:**
- Asserting specific sentinel errors: `assert.ErrorIs(t, err, errs.ErrValidation)`.

---

*Testing analysis: 2025-05-14*
