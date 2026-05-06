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
            Hydra[Ory Hydra Pod - OAuth2/OIDC]
        end

    subgraph NS_App [Namespace: production]
            DemoSvc[Demo Service Pod - Go Clean Arch]
            MonitorSvc[Monitor Service Pod - SSE/WS]
        end

        subgraph NS_Infra [Namespace: infra]
            Postgres[(PostgreSQL StatefulSet)]
            Redis[(Redis Cache)]
            Jaeger[SigNoz/Jaeger - Tracing & Metrics]
            Kafka[(Apache Kafka)]
        end
    end

    ClientApp -- HTTPS/JWT --> KrakenD
    KrakenD -- Auth Flow --> Hydra
    Hydra -- Identity Check --> Kratos
    KrakenD -- API Forwarding --> DemoSvc
    MonitorSvc -- SSE/WebSocket --> ClientApp
    MonitorSvc -- Listen --> Kafka
    ClientApp -- Query API --> Jaeger
    DemoSvc -- Fetch JWKS --> Hydra
    DemoSvc -- Query/Persist --> Postgres
    DemoSvc -- Cache --> Redis
    DemoSvc -- Publish --> Kafka

    classDef cluster fill:#f8fafc,stroke:#334155,stroke-width:2px;
    classDef namespace fill:#ffffff,stroke:#94a3b8,stroke-dasharray: 5 5;
    class K8s_Cluster cluster;
    class NS_Gateway,NS_Auth,NS_App,NS_Infra namespace;
```

> [!NOTE]
> **Visual View (Excalidraw):** ![High-Level Infrastructure](./assets/system-infrastructure.png)

## 2. Sprint 2: Distributed Security Flow
Quy trình xác thực JWT và phân quyền Casbin tại từng service.

```mermaid
sequenceDiagram
    autonumber
    participant Client as Client App
    participant Gateway as KrakenD Gateway
    participant Hydra as OIDC Provider (Hydra)
    participant Kratos as Identity Server (Kratos)
    participant Service as Demo Service (Go)

    Note over Client, Service: Khởi tạo: Service tải Public Keys (JWKS) từ Hydra về bộ nhớ

    Client->>Hydra: Yêu cầu Đăng nhập (OAuth2 Flow)
    Hydra->>Kratos: Redirect tới Login UI (Browser)
    Kratos->>Client: Hiển thị Form Login
    Client->>Kratos: Submit Credentials
    Kratos-->>Hydra: Xác thực thành công (Identity)
    
    Note over Client, Hydra: Auto-Accept Consent (Server Component)
    Hydra->>Client: Redirect tới /consent (Server-side Accept)
    Client-->>Hydra: Hoàn tất cấp quyền (Behind the scenes)

    Hydra-->>Client: Trả về JWT (Access Token)

    Client->>Gateway: Request API + JWT
    Gateway->>Gateway: CORS & Tracing
    Gateway->>Service: Forward Request

    rect rgb(240, 253, 244)
        Note right of Service: Backend Security Layer
        Service->>Service: 1. Verify Signature (JWKS offline)
        Service->>Service: 2. Casbin RBAC Check
    end

    Service-->>Client: Data Response (200 OK)
```
