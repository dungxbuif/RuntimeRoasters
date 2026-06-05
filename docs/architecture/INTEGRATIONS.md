---
artifact_type: integration_spec
id: INTEGRATIONS_SPEC
status: active
owner: shared
---

# Integration & Event Topology

## Event Bus (Kafka)
The system uses Apache Kafka for asynchronous, reliable communication between microservices.

### Topic Registry

| Topic | Producer | Primary Consumer(s) | Description |
| :--- | :--- | :--- | :--- |
| `farm.harvest.created` | Farm | Warehouse | New harvest is ready for pickup. |
| `retail.order.created` | Retail | Warehouse, Payment | New customer order placed. |
| `payment.completed` | Payment | Retail, Warehouse | Payment successfully processed. |
| `warehouse.stock.reserved`| Warehouse | Retail, Logistics | Inventory locked for order. |
| `warehouse.pickup.requested`| Warehouse | Logistics | Request for farm-to-warehouse transport. |
| `logistics.delivery.assigned`| Logistics | Socket, Trace | Driver assigned to delivery. |
| `notification.created` | Any | Socket | Real-time notification for UI. |

## Service Dependencies

### Synchronous (gRPC)
- **Gateways -> Auth**: Token validation (KrakenD, Identity Proxy).
- **Retail -> Warehouse**: Immediate inventory check (Optional/Planned).
- **Any -> Trace**: Querying history.

### Asynchronous (Saga)
The system implements a **Choreographed Saga** for order fulfillment:
1. `retail.order.created`
2. **Payment** processes -> `payment.completed`
3. **Warehouse** reserves -> `warehouse.stock.reserved`
4. **Logistics** assigns -> `logistics.delivery.assigned`
5. **Retail** completes -> `OrderStatusCompleted`
