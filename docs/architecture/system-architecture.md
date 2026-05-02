# System Architecture Diagrams

Tài liệu này chứa các sơ đồ kiến trúc và luồng dữ liệu của Runtime Roasters, được vẽ bằng Mermaid.

## 1. High-Level Infrastructure (Kubernetes)
Sơ đồ mô phỏng cấu trúc phân lớp trong K8s Cluster.

```mermaid
graph TB
    subgraph External_Internet [Public Internet]
        ClientApp[Client App - Next.js 15]
    end

    subgraph K8s_Cluster [Kubernetes Cluster]
        
        subgraph NS_Gateway [Namespace: gateway]
            KrakenD[KrakenD Pod - Ingress Controller]
        end

        subgraph NS_Auth [Namespace: auth]
            Kratos[Ory Kratos Pod - Identity Server]
        end

        subgraph NS_App [Namespace: production]
            DemoSvc[Demo Service Pod - Go Clean Arch]
        end

        subgraph NS_Infra [Namespace: infra]
            Postgres[(PostgreSQL StatefulSet)]
            Redis[(Redis Cache)]
            Jaeger[Jaeger Deployment - Tracing]
        end
    end

    ClientApp -- HTTPS/JWT --> KrakenD
    KrakenD -- Auth Check --> Kratos
    KrakenD -- API Forwarding --> DemoSvc
    DemoSvc -- Fetch JWKS --> Kratos
    DemoSvc -- Query/Persist --> Postgres
    DemoSvc -- Cache --> Redis
    KrakenD -. OTLP .-> Jaeger
    DemoSvc -. OTLP .-> Jaeger

    classDef cluster fill:#f8fafc,stroke:#334155,stroke-width:2px;
    classDef namespace fill:#ffffff,stroke:#94a3b8,stroke-dasharray: 5 5;
    class K8s_Cluster cluster;
    class NS_Gateway,NS_Auth,NS_App,NS_Infra namespace;
```

## 2. Sprint 2: Distributed Security Flow
Quy trình xác thực JWT và phân quyền Casbin tại từng service.

```mermaid
sequenceDiagram
    autonumber
    participant Client as Client App
    participant Gateway as KrakenD Gateway
    participant IDP as Identity Server (Kratos)
    participant Service as Demo Service (Go)

    Note over Client, Service: Khởi tạo: Service tải Public Keys (JWKS) về bộ nhớ

    Client->>IDP: Đăng nhập (Username/Password)
    IDP-->>Client: Trả về JWT (chứa sub, role, org_id)

    Client->>Gateway: Request API + JWT
    Gateway->>Gateway: CORS & Tracing
    Gateway->>Service: Forward Request

    rect rgb(240, 253, 244)
        Note right of Service: Backend Security Layer
        Service->>Service: 1. Verify Signature (In-memory)
        Service->>Service: 2. Casbin RBAC Check
    end

    Service-->>Client: Data Response (200 OK)
```
