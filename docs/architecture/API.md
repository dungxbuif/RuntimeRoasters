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
- `GET /v1/public/retail/stores`: List active stores for the public kiosk.
- `GET /v1/public/retail/menu-items`: List the active chain-wide menu without authentication.
- `GET /v1/public/retail/stores/{store_id}/menu-items`: List the active menu with materialized store availability without authentication.
- `POST /v1/public/retail/sales`: Create a no-payment sold-cup sale without authentication. Requires `Idempotency-Key`; every item contains one `menu_item_id` and one client-issued `product_id`.
- `GET /v1/retail/sales`: List sales within the authenticated manager's store scope.
- `GET /v1/retail/sales/{sale_id}`: Get an in-scope sale and its sold cups.
- `GET /v1/retail/stores/{store_id}/sold-items`: List sold cups for an assigned store.
- `GET /v1/orders`: List orders.
- `POST /v1/orders`: Create a new order (requires Idempotency-Key).
- `GET /v1/orders/{id}`: Get order status.
- `POST /v1/orders/{id}/confirm`: Confirm order (Saga manual step).

Retail internal support route, not exposed through KrakenD:

- `GET /v1/internal/retail/sold-items?limit=25`: List recent sold cups for trace composition. Requires `X-Internal-Secret`.
- `GET /v1/internal/retail/sold-items/{product_id}`: Resolve the operational sold-cup and lineage record. Requires `X-Internal-Secret`.

`POST /v1/public/retail/sales` atomically writes the sale, one sale item per
cup, FIFO lot decrement, stock movement, affected materialized availability,
and a pending `retail.sale.completed` outbox CloudEvent. Validation errors use
`400`, missing active resources use `404`, and duplicate identities or
insufficient inventory use `409`.

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
- `GET /v1/public/trace-items`: List recently sold traceable cups without authentication.
- `GET /v1/public/trace/{trace_code}`: Resolve a sanitized sold-cup journey without authentication; `trace_code` equals Retail `product_id`.
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
