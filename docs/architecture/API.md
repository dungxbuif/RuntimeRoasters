---
artifact_type: api_spec
id: API_SPEC
status: active
owner: shared
---

# API Specification

This document details the API contracts across the Runtime Roasters platform, including the public API Gateway (KrakenD) endpoints, private internal-only endpoints, and service-to-service gRPC interfaces.

---

## 1. System-wide Conventions

### Error Schema
All public REST API endpoints follow a standardized error structure on failure:
```json
{
  "error": {
    "code": "UNAUTHENTICATED",
    "message": "Missing or invalid authorization token",
    "details": []
  }
}
```

### Common Error Codes
* `INVALID_ARGUMENT` (HTTP 400): Malformed payload, missing required fields.
* `UNAUTHENTICATED` (HTTP 401): Missing or invalid Bearer token.
* `PERMISSION_DENIED` (HTTP 403): User lacks required Casbin role policy.
* `NOT_FOUND` (HTTP 404): Resource does not exist.
* `ALREADY_EXISTS` (HTTP 409): Resource conflict (e.g. idempotency or duplicate unique keys).
* `INTERNAL` (HTTP 500): Server error (database, cache, or message broker outages).

---

## 2. Public Gateway (KrakenD) & REST Endpoints

### 2.1 Identity & Authentication (Auth Service)

#### POST `/v1/auth/login/accept`
* **Description**: Processes login acceptance from Kratos/Hydra bridge redirect.
* **Authentication**: None (Internal system integration).
* **Request**:
  ```json
  {
    "login_challenge": "challenge_123abc",
    "subject": "user-uuid-999",
    "remember": true,
    "remember_for": 3600
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "redirect_to": "http://localhost:4444/oauth2/auth?client_id=..."
  }
  ```

#### GET `/v1/auth/me`
* **Description**: Retrieves identity profile of the currently logged-in user.
* **Authentication**: Bearer JWT.
* **Request Headers**:
  * `Authorization: Bearer <token>`
* **Response (200 OK)**:
  ```json
  {
    "user": {
      "id": "usr_77239a0",
      "email": "manager@runtimeroasters.com",
      "name": "Hoan Kiem Manager",
      "role": "STORE_MGR"
    }
  }
  ```

#### POST `/v1/auth/logout`
* **Description**: Self-service immediate logout and revocation of all active access tokens for the authenticated subject.
* **Authentication**: Bearer JWT.
* **Request Body**:
  ```json
  {
    "reason": "user_logout_all"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "all subject sessions revoked successfully"
  }
  ```

#### POST `/v1/admin/users/{subject}/security/revoke-all`
* **Description**: Administrative revocation of all access tokens for a specific user.
* **Authentication**: Bearer JWT (requires `ADMIN` role).
* **Request Body**:
  ```json
  {
    "reason": "compromised_credentials"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "subject": "usr_77239a0",
    "tokens_valid_after": "2026-06-10T04:20:00Z",
    "version": 2
  }
  ```

---

### 2.2 User Management (Auth Service)

#### POST `/v1/users`
* **Description**: Registers a new system user inside Kratos and configures Casbin authorization policy.
* **Authentication**: Bearer JWT (requires `ADMIN` role).
* **Request**:
  ```json
  {
    "email": "wh_mgr@runtimeroasters.com",
    "password": "securepassword123",
    "name": "Warehouse Manager",
    "role": "WAREHOUSE_MGR"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "user": {
      "id": "usr_883a1b0",
      "email": "wh_mgr@runtimeroasters.com",
      "name": "Warehouse Manager",
      "role": "WAREHOUSE_MGR"
    }
  }
  ```

---

### 2.3 Farm & Harvest (Farm Service)

#### POST `/v1/farms`
* **Description**: Registers a new coffee farm.
* **Authentication**: Bearer JWT (requires `ADMIN` role).
* **Request**:
  ```json
  {
    "name": "Dalat Highlands Farm",
    "location": "Dalat, Lam Dong",
    "area": 25.4,
    "coffee_type": "Arabica Catimor"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "id": "frm_ HighlandsDalat",
    "name": "Dalat Highlands Farm",
    "location": "Dalat, Lam Dong",
    "area": 25.4,
    "coffee_type": "Arabica Catimor"
  }
  ```

#### POST `/v1/harvests`
* **Description**: Declares a coffee cherry harvest from an assigned farm. Automatically schedules logistics pickup shipments.
* **Authentication**: Bearer JWT (requires `FARM_MANAGER` role scoped to the farm).
* **Request**:
  ```json
  {
    "farm_id": "frm_ HighlandsDalat",
    "quantity": 1850.5,
    "harvest_date": "2026-06-10"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "id": "hrv_90123a0",
    "farm_id": "frm_ HighlandsDalat",
    "quantity": 1850.5,
    "harvest_date": "2026-06-10",
    "status": "PENDING_PICKUP"
  }
  ```

---

### 2.4 Retail & POS (Retail Service)

#### POST `/v1/orders`
* **Description**: Places a warehouse stock supply order for a store. This starts the async Saga transaction.
* **Authentication**: Bearer JWT (requires `STORE_MGR` role).
* **Request Headers**:
  - `Idempotency-Key: idemp_112233`
* **Request**:
  ```json
  {
    "store_id": "str_hoankiem",
    "warehouse_id": "wh_hanoi",
    "items": [
      {
        "sku": "sku_arabica_medium",
        "quantity": 50
      }
    ]
  }
  ```
* **Response (202 Accepted)**:
  ```json
  {
    "id": "ord_9901abc",
    "status": "PROCESSING",
    "store_id": "str_hoankiem",
    "warehouse_id": "wh_hanoi",
    "items": [
      {
        "sku": "sku_arabica_medium",
        "quantity": 50
      }
    ]
  }
  ```

#### POST `/v1/public/retail/sales`
* **Description**: Registers a POS cup sale at a store (public kiosk path). Allocates inventory in FIFO order and maps lineage.
* **Authentication**: None (Public Kiosk).
* **Request Headers**:
  - `Idempotency-Key: idemp_sale_8899`
* **Request**:
  ```json
  {
    "store_id": "str_hoankiem",
    "items": [
      {
        "menu_item_id": "mn_espresso_double",
        "product_id": "cup_plat_0610a"
      }
    ]
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "sale_id": "sal_88231c",
    "store_id": "str_hoankiem",
    "occurred_at": "2026-06-10T04:25:00Z",
    "items": [
      {
        "menu_item_id": "mn_espresso_double",
        "product_id": "cup_plat_0610a"
      }
    ]
  }
  ```

---

### 2.5 Logistics & Vehicle Tracking (Logistics Service)

#### POST `/v1/logistics/drivers/location`
* **Description**: Submits the driver's current coordinates for live route monitoring.
* **Authentication**: Bearer JWT (requires `DRIVER` role).
* **Request**:
  ```json
  {
    "driver_id": "drv_loc_456",
    "latitude": 21.0285,
    "longitude": 105.8542
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "success": true
  }
  ```

#### GET `/v1/logistics/drivers/location`
* **Description**: Retrieves current positions of all active drivers from Valkey GEO.
* **Authentication**: Bearer JWT (requires `WAREHOUSE_MGR`, `RETAIL_MGR`, or `ADMIN`).
* **Response (200 OK)**:
  ```json
  {
    "locations": [
      {
        "driver_id": "drv_loc_456",
        "latitude": 21.0285,
        "longitude": 105.8542
      }
    ]
  }
  ```

---

### 2.6 Real-time Event Streaming (Socket Service)

#### WS `/v1/realtime/stream`
* **Description**: Authenticated persistent WebSocket for role-scoped operational updates.
* **Authentication**: Real-time Ticket (obtained via `/v1/realtime/tickets`).
* **Features**: Receives `topology.event` notifications, GPS updates, and system alerts.

#### WS `/v1/realtime/public/topology/ws`
* **Description**: Public anonymous WebSocket for sanitized architecture visualization.
* **Authentication**: None.
* **Features**: Receives anonymized GPS updates (Truck Name/Plate only) and public domain events.

---

### 2.7 Warehouse & Inventory (Warehouse Service)

#### GET `/v1/warehouse/inventory`
* **Description**: Retrieves current in-stock quantities per SKU.
* **Authentication**: Bearer JWT (requires `WAREHOUSE_MGR` or `ADMIN`).
* **Response (200 OK)**:
  ```json
  [
    {
      "warehouse_id": "wh_hanoi",
      "sku": "sku_arabica_medium",
      "quantity": 1240
    },
    {
      "warehouse_id": "wh_hanoi",
      "sku": "sku_robusta_dark",
      "quantity": 890
    }
  ]
  ```

---

### 2.7 Public Provenance Trace (Trace Service)

#### GET `/v1/public/trace/{trace_code}`
* **Description**: Resolves the complete origin and transportation journey of a cup by its identifier.
* **Authentication**: None (Public Visitor).
* **Response (200 OK)**:
  ```json
  {
    "product_id": "cup_plat_0610a",
    "menu_item_name": "Double Espresso",
    "sold_at": "2026-06-10T04:25:00Z",
    "store_name": "Hoan Kiem Store",
    "journey": [
      {
        "stage": "Harvest",
        "date": "2026-06-01",
        "origin": "Dalat Highlands Farm",
        "details": "Arabica Catimor, quantity: 1850.5kg"
      },
      {
        "stage": "Processing & Roasting",
        "date": "2026-06-04",
        "facility": "Hanoi Central Warehouse",
        "details": "Batch: bat_roast_442"
      },
      {
        "stage": "Logistics Dispatch",
        "date": "2026-06-06",
        "shipment_id": "shp_wh_to_store_112",
        "driver": "Nguyen Van A"
      }
    ]
  }
  ```

---

## 3. Private Internal HTTP Endpoints
These endpoints are not exposed to the public internet via the KrakenD API Gateway and are protected by `X-Internal-Secret` matching.

### 3.1 GET `/v1/internal/retail/sold-items/{product_id}`
* **Description**: Retrieves a single cup's POS lineage record from `retail-service` database. Called synchronously by `trace-service` during composition.
* **Headers**:
  - `X-Internal-Secret: <shared_secret>`
* **Response (200 OK)**:
  ```json
  {
    "product_id": "cup_plat_0610a",
    "menu_item_id": "mn_espresso_double",
    "lot_id": "lot_dalat_arabica_0601",
    "quantity": 1,
    "sold_at": "2026-06-10T04:25:00Z"
  }
  ```

---

## 4. Internal gRPC Services
Services communicate internally over gRPC on ports `50052-50060`.

| Service | Port | Proto Definition |
| :--- | :--- | :--- |
| **Auth** | 50052 | `api/runtime/auth/v1/auth.proto` |
| **Farm** | 50053 | `api/runtime/farm/v1/farm.proto` |
| **Retail** | 50054 | - |
| **Logistics** | 50055 | - |
| **Warehouse** | 50059 | - |
| **Socket** | 50060 | - |

### 4.1 AuthService (Auth Service)

#### Method: `GetTokenValidityBoundary`
* **Proto Definition**:
  ```protobuf
  rpc GetTokenValidityBoundary(GetTokenValidityBoundaryRequest) returns (GetTokenValidityBoundaryResponse);
  ```
* **Description**: Returns the revocation boundary timestamp for a subject. Used by other microservices upon cache misses.
* **Request**:
  ```protobuf
  message GetTokenValidityBoundaryRequest {
    string subject = 1;
  }
  ```
* **Response**:
  ```protobuf
  message GetTokenValidityBoundaryResponse {
    string subject = 1;
    string tokens_valid_after = 2; // RFC3339 formatted string
  }
  ```
