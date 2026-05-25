# 🔗 Technical Reference: API Contracts & Messaging

Runtime Roasters uses a combination of **Synchronous gRPC** for command/query and **Asynchronous Kafka** for state propagation. This document specifies the contracts that bind the microservices together.

---

## 1. Internal Communication (gRPC)

All internal calls between services use gRPC with **Protobuf** serialization.

### Core Service Interfaces
| Service | Protobuf Definition | Key Methods |
| :--- | :--- | :--- |
| **Auth** | `api/runtime/auth/v1/auth.proto` | `GetFullSnapshot`, `VerifyToken`, `ListUsers` |
| **Farm** | `api/runtime/farm/v1/farm.proto` | `CreateHarvest`, `GetFarm`, `ListFarms` |
| **Warehouse**| `api/runtime/warehouse/v1/wh.proto` | `ReserveStock`, `ReleaseStock`, `GetInventory` |

### Rationale: Why gRPC?
- **Type Safety**: Strict contracts prevent runtime errors due to missing fields.
- **Performance**: High-speed binary protocol with lower overhead than JSON.
- **Code Generation**: Modern tooling (**Buf**) allows for automatic client/server stub generation.

---

## 2. Edge API (REST/JSON)

External clients (Frontend/Mobile) communicate via the **KrakenD API Gateway** using REST.

| Method | Endpoint | Internal gRPC Mapping | Auth Required |
| :--- | :--- | :--- | :--- |
| `POST` | `/v1/auth/login` | `AuthService.Login` | No |
| `GET` | `/v1/stores` | `RetailService.ListStores` | Yes |
| `POST` | `/v1/orders` | `RetailService.CreateOrder` | Yes |
| `GET` | `/v1/harvests` | `FarmService.ListHarvests` | Yes |

---

## 3. Asynchronous Messaging (Kafka)

Kafka is the "central nervous system" of the platform, used for Sagas and CQRS.

### Key Kafka Topics & Events
| Topic | Event Type | Description |
| :--- | :--- | :--- |
| `order.saga.events` | `OrderCreated` | Triggers stock reservation in Warehouse. |
| `auth.policy.changed`| `PolicyUpdated` | Notifies services to reload Casbin rules. |
| `farm.harvest.created`| `HarvestCreated`| Propagates data to Warehouse, Trace, and Audit. |
| `logistics.gps.updated` | `DriverLocation` | High-frequency GPS updates for tracking. |

### Event Schema Standard
Every event follows the **CloudEvents** specification to ensure metadata consistency across the ecosystem.

The detailed event contract is maintained in [`CLOUDEVENTS_CONTRACT.md`](./CLOUDEVENTS_CONTRACT.md).

CloudEvents are emitted in JSON format. Required metadata:
- `id`, `type`, `source`, `subject`, `time`, `datacontenttype=application/json`.
- Extensions: `correlationid`, optional `causationid`, `traceid` when an OTel span exists, and relevant business IDs such as `orderid`, `storeid`, `shipmentid`, `harvestid`, `batchid`, `farmid`, `warehouseid`, `driverid`, `vehicleid`.
- `data` contains the typed business payload only.

OpenTelemetry `trace_id` is not a business identifier. It is used for SigNoz/ClickHouse observability and request/event correlation only.

---

## 4. System Constants & Enums

To maintain consistency, all services share these core domain values:

| Enum Group | Values | Description |
| :--- | :--- | :--- |
| **Order Status** | `PENDING`, `PREPARING`, `SHIPPING`, `COMPLETED`, `REJECTED` | Lifecycle of a retail transaction. |
| **Outbox Status** | `PENDING`, `COMPLETED`, `FAILED` | Reliability state of a Kafka message. |
| **Coffee Type** | `ARABICA`, `ROBUSTA`, `CHERRY`, `CULI` | Core product catalog. |

---
*Technical reference for Runtime Roasters Integration Standards.*
