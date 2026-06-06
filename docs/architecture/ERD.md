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
| **Order** | Existing warehouse supply-order workflow. |
| **Menu** | Versionable chain-wide menu. |
| **MenuItem** | One sellable drink-size entry from SAMPLE_MENU.md. |
| **InventoryLot** | Store stock with upstream harvest/batch/warehouse lineage. |
| **Sale** | Public POS invoice, separate from supply orders. |
| **SaleItem** | One sold cup, UI-issued `product_id`, and QR identity. |
| **StockMovement** | Immutable received/sold/adjustment ledger. |
| **StoreMenuInventory** | Materialized available units by store and menu item. |

```mermaid
erDiagram
    RETAIL_STORE ||--o{ ORDER : "places supply order"
    MENU ||--o{ MENU_ITEM : "contains"
    RETAIL_STORE ||--o{ INVENTORY_LOT : "holds"
    RETAIL_STORE ||--o{ SALE : "issues"
    SALE ||--|{ SALE_ITEM : "contains sold cups"
    MENU_ITEM ||--o{ SALE_ITEM : "snapshotted by"
    INVENTORY_LOT ||--o{ SALE_ITEM : "provides lineage"
    INVENTORY_LOT ||--o{ STOCK_MOVEMENT : "has ledger"
    RETAIL_STORE ||--o{ STORE_MENU_INVENTORY : "materializes"
    MENU_ITEM ||--o{ STORE_MENU_INVENTORY : "available at"
```

### Logistics Service
| Entity | Description |
| :--- | :--- |
| **Shipment** | Tracking unit for moving goods (Farm->Warehouse or Warehouse->Store). |
| **Driver** | Personnel assigned to shipments. |
| **Vehicle** | Asset used for transport. |
