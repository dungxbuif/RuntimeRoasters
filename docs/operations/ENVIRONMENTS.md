# Environments

| Environment | Runtime Shape | Notes |
| --- | --- | --- |
| Local | Docker Compose infrastructure plus local app processes or containers | Uses `deployments/docker-compose.dev.yaml`; demo seed/reset is allowed. |
| Staging | Kubernetes planned | Mirrors production topology but may include demo seed data and self-hosted dependencies. |
| Production | Kubernetes planned | Managed Postgres/Kafka/Valkey/Elasticsearch preferred; demo seed data disabled. |

Production and staging configs must not hardcode localhost. Use service DNS, public issuer URLs, and environment-specific CORS allowlists.
