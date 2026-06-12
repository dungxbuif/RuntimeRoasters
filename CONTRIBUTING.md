# Contributing to Runtime Roasters

First off, thank you for considering contributing to Runtime Roasters! It's people like you that make open-source software such a great community.

## 🚀 Getting Started

To get started with local development, ensure you have the prerequisites installed:
- Go 1.24+
- Docker & Docker Compose
- Task (`go install github.com/go-task/task/v3/cmd/task@latest`)

### Local Development Workflow

1. **Start the Infrastructure**
   ```bash
   task infra
   ```
   This boots PostgreSQL, Kafka, Valkey, Elasticsearch, Cassandra, and Identity services.

2. **Seed the Database**
   ```bash
   task env:reset
   ```
   This command resets databases, applies migrations, and seeds test data.

3. **Run the Services**
   ```bash
   task dev
   ```
   This runs the backend microservices using `air` for hot-reloading.

## 📝 Code Standards

We adhere to Clean Architecture principles. When adding a new feature:
- Keep domain entities isolated in `internal/domain`.
- Place business logic in `internal/usecase`.
- Infrastructure implementations (DB, Kafka) go in `internal/infrastructure`.
- Ensure you run the linter before submitting a PR:
  ```bash
  golangci-lint run ./...
  ```

## 🤝 Submitting Pull Requests

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/amazing-feature`).
3. Make your changes and commit them (`git commit -m 'Add some amazing feature'`).
4. Push to the branch (`git push origin feature/amazing-feature`).
5. Open a Pull Request!

If you're opening a PR that changes architecture, please submit an Architecture Decision Record (ADR) in `docs/decisions/`.
