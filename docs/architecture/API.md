---
artifact_type: api_spec
id: API_SPEC
status: active
owner: shared
---

# API Specification

## Public Gateway (KrakenD)
The system exposes a unified REST API through KrakenD on port `8081`.

### Auth Endpoints
- `POST /v1/auth/login/accept`: Accept Hydra login request.
- `GET /v1/auth/me`: Get current user identity.
- `GET /v1/auth/system/status`: System-wide health check.
- `POST /v1/auth/seed`: Trigger auth seeding.

### User Endpoints
- `GET /v1/users`: List users.
- `POST /v1/users`: Create a new user.

### Farm Endpoints
- `GET /v1/farms`: List farms.
- `POST /v1/farms`: Register a farm.
- `GET /v1/farms/{id}`: Get farm details.
- `PUT /v1/farms/{id}`: Update farm.
- `DELETE /v1/farms/{id}`: Delete farm.
- `GET /v1/harvests`: List all harvests.
- `POST /v1/harvests`: Record a harvest.
- `GET /v1/farms/{id}/harvests`: List harvests for a specific farm.

### Retail Endpoints
- `GET /v1/retail/stores`: List retail stores.
- `POST /v1/retail/stores`: Create a retail store.
- `GET /v1/orders`: List orders.
- `POST /v1/orders`: Create a new order (requires Idempotency-Key).
- `GET /v1/orders/{id}`: Get order status.
- `POST /v1/orders/{id}/confirm`: Confirm order (Saga manual step).

### Warehouse Endpoints
- `GET /v1/warehouse/warehouses`: List warehouses.
- `POST /v1/warehouse/warehouses`: Create a warehouse.
- `GET /v1/warehouse/inventory`: View global warehouse inventory.
- `GET /v1/warehouse/batches`: List production batches.
- `POST /v1/warehouse/batches`: Create production batch.

### Logistics Endpoints
- `GET /v1/logistics/shipments`: List shipments.
- `GET /v1/logistics/shipments/{id}`: Get shipment details.
- `POST /v1/logistics/shipments/{id}/assign`: Assign driver/vehicle.
- `POST /v1/logistics/drivers/location`: Update driver GPS.

### Trace Endpoints
- `GET /v1/trace/{id}`: Get trace topology for an entity.
- `GET /v1/trace/{id}/document`: Get detailed trace document.

## Internal gRPC Services
Services communicate over gRPC on ports `50052-50060`.

| Service | Port | Proto Definition |
| :--- | :--- | :--- |
| **Auth** | 50052 | `api/runtime/auth/v1/auth.proto` |
| **Farm** | 50053 | `api/runtime/farm/v1/farm.proto` |
| **Retail** | 50054 | - |
| **Logistics** | 50055 | - |
| **Warehouse** | 50059 | - |
| **Socket** | 50060 | - |
