# Runtime Roasters K8s Readiness Plan

  ## Summary

  Chuẩn bị Runtime Roasters để deploy được lên một cụm Kubernetes có sẵn. Không plan phần dựng hạ tầng cluster, registry, ingress controller, storage class hay DNS; các phần đó được xem là external
  prerequisites.

  Mục tiêu là biến project hiện tại từ mô hình local/dev .env + docker compose infra + go run/air thành bộ artifact deployable trên K8s: container images, config contract, Helm manifests, migration/
  seeding jobs, health probes, và quy trình test local trước khi deploy.

  ## Key Changes

  - Container hóa app
      - Thêm Dockerfile production cho Go services: auth, farm, retail, warehouse, logistics, payment, trace, audit, socket.
      - Thêm Dockerfile production cho client-app dùng Next.js standalone build.
      - Chuẩn hóa image naming bằng biến:
          - IMAGE_REGISTRY
          - IMAGE_TAG
          - SERVICE_NAME
  - Chuẩn hóa config cho K8s
      - Tạo một env contract duy nhất cho từng service, tách rõ:
          - Config thường: ConfigMap.
          - Secret: Kubernetes Secret.
      - Không phụ thuộc vào file .env nằm trong repo khi chạy container.
      - Chuyển toàn bộ endpoint local sang DNS nội bộ K8s:
          - Postgres: postgres:5432
          - Kafka: kafka:9092
          - Valkey: valkey:6379
          - Elasticsearch: elasticsearch:9200
          - Cassandra: cassandra:9042
          - Auth gRPC: auth-service:50052
  - Helm chart
      - Tạo deployments/k8s/runtime-roasters/.
      - Chart quản lý:
          - Deployments.
          - Services.
          - ConfigMaps.
          - Secret templates.
          - Migration Jobs.
          - Seeder Jobs.
          - Ingress values, nhưng không hardcode domain.
      - Mỗi service có HTTP/gRPC ports, probes, resources, env, image, replica count.
  - KrakenD/Ory readiness
      - Tách KrakenD config cho K8s, thay host.docker.internal:* bằng service DNS.
      - Hydra/Kratos URL, issuer, callback, CORS origin chuyển thành Helm values.
      - Không hardcode localhost trong production/K8s config.
  - Database migration and seed
      - Tạo migration job chạy SQL migrations hiện có theo thứ tự:
          - auth_db, farm_db, retail_db, warehouse_db, logistics_db, payment_db, trace_db, audit_db.
      - Tạo Kratos/Hydra migration jobs.
      - Tạo seed job cho demo identities, OAuth client, policies, và dữ liệu bootstrap cần thiết.
      - Jobs phải idempotent hoặc fail rõ ràng nếu chạy lại không an toàn.

  ## Local Test Phase

  - Build toàn bộ images bằng script hoặc Make/Task command.
  - Chạy container app locally với env injection, không dùng .env service-local.
  - Có thể tiếp tục dùng compose infra hiện tại cho Postgres/Kafka/Ory/etc trong phase local.
  - Validate:
      - Container boot được.
      - Health endpoints pass:
          - /health/live
          - /health/ready
      - KrakenD route được tới các service container.
      - OIDC login/callback vẫn hoạt động.
      - Gateway protected routes trả 401 khi không có token.
      - Demo SAGA chạy qua retail, payment, warehouse, logistics, trace, audit.

  ## Deploy Readiness Artifacts

  - Add:
      - deployments/k8s/runtime-roasters/Chart.yaml
      - deployments/k8s/runtime-roasters/values.yaml
      - deployments/k8s/runtime-roasters/values.local.yaml
      - deployments/k8s/runtime-roasters/values.example.prod.yaml
      - deployments/k8s/runtime-roasters/templates/*.yaml
      - docs/deploy/k8s-readiness.md
      - docs/deploy/env-contract.md
  - Add commands:
      - task image:build
      - task image:push
      - task k8s:render
      - task k8s:deploy
      - task k8s:smoke

  ## Test Plan

  - Backend:
      - cd src && GOCACHE=/private/tmp/runtime-roasters-go-cache go test ./pkg/base/auth/... ./pkg/base/casbin/... ./pkg/base/security/... ./apps/auth-service/... ./apps/retail-service/... ./apps/
        payment-service/... ./apps/logistics-service/... ./apps/trace-service/... ./apps/audit-service/... ./apps/warehouse-service/...
  - Frontend:
      - cd src/apps/client-app && npm run lint && npm run build
  - Manifest validation:
      - helm lint deployments/k8s/runtime-roasters
      - helm template runtime-roasters deployments/k8s/runtime-roasters -f deployments/k8s/runtime-roasters/values.local.yaml
  - K8s smoke:
      - kubectl get pods
      - kubectl get svc
      - kubectl logs job/<migration-job>
      - curl -i <gateway-url>/v1/users returns 401 without token.
      - Browser login works through configured public URL.
      - Playwright live evidence test passes after deploy.

  ## Assumptions

  - Kubernetes cluster, registry, storage class, ingress controller, DNS, and TLS are provided separately.
  - This plan only makes Runtime Roasters deployable to K8s.
  - demo-service is excluded from the first K8s-ready scope.
  - Secrets are supplied by Kubernetes Secret values or an external secret manager later.
  - Helm is the deployment interface for the first production-shaped K8s release.