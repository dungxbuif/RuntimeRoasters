---
artifact_type: erd_spec
id: ERD_SPEC
status: active
owner: shared
---

# Entity Relationship Diagram

## System Overview

```mermaid
erDiagram
    FARM ||--o{ HARVEST : "produces"
    HARVEST ||--o{ PICKUP_REQUEST : "triggers"
    PICKUP_REQUEST ||--|| SHIPMENT : "has"
    PICKUP_REQUEST ||--|| INTAKE : "results in"
    INTAKE }o--o| PRODUCTION_BATCH : "belongs to"
    PRODUCTION_BATCH ||--o{ ROAST_RUN : "contains"
    PRODUCTION_BATCH ||--o{ INVENTORY : "updates"
    
    ORDER ||--o{ RETAIL_STORE : "placed at"
    ORDER ||--o| DISPATCH_REQUEST : "triggers"
    DISPATCH_REQUEST ||--|| SHIPMENT : "has"
    
    SHIPMENT }o--|| DRIVER : "assigned to"
    SHIPMENT }o--|| VEHICLE : "uses"
    
    WAREHOUSE ||--o{ INTAKE : "hosts"
    WAREHOUSE ||--o{ PRODUCTION_BATCH : "manages"
    WAREHOUSE ||--o{ INVENTORY : "stores"
    WAREHOUSE ||--o{ DISPATCH_REQUEST : "fulfills"
```

## Core Entities

### Farm Service
| Entity | Description |
| :--- | :--- |
| **Farm** | Coffee plantation with location, area, and coffee type. |
| **Harvest** | A specific collection of coffee cherries from a farm. |

### Warehouse Service
| Entity | Description |
| :--- | :--- |
| **Warehouse** | Storage and processing facility. |
| **Intake** | Raw material entry into the warehouse. |
| **ProductionBatch** | Aggregation of intakes for roasting/processing. |
| **Inventory** | Stock levels per SKU in a warehouse. |
| **PickupRequest** | Request to collect harvest from a farm. |
| **DispatchRequest** | Request to send reserved stock to a store. |

### Retail Service
| Entity | Description |
| :--- | :--- |
| **RetailStore** | Coffee shop location. |
| **Order** | Customer order containing multiple SKUs. |

### Logistics Service
| Entity | Description |
| :--- | :--- |
| **Shipment** | Tracking unit for moving goods (Farm->Warehouse or Warehouse->Store). |
| **Driver** | Personnel assigned to shipments. |
| **Vehicle** | Asset used for transport. |
