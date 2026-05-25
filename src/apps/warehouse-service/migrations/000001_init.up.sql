-- Warehouse Service Schema
CREATE TABLE IF NOT EXISTS pick_up_requests (
    id                 VARCHAR(64) PRIMARY KEY,
    harvest_id         VARCHAR(64) UNIQUE NOT NULL,
    farm_id            VARCHAR(80),
    warehouse_id       VARCHAR(80) NOT NULL,
    origin_location_id VARCHAR(80),
    quantity           DECIMAL(10,2) NOT NULL,
    coffee_type        VARCHAR(20) NOT NULL,
    origin_code        VARCHAR(10) NOT NULL,
    status             VARCHAR(32) DEFAULT 'REQUESTED',
    notification_id    VARCHAR(64),
    shipment_id        VARCHAR(64),
    created_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    dispatched_at      TIMESTAMP WITH TIME ZONE,
    received_at        TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS production_batches (
    id                   VARCHAR(80) PRIMARY KEY,
    batch_id             VARCHAR(255) UNIQUE NOT NULL,
    warehouse_id         VARCHAR(80),
    status               VARCHAR(20) DEFAULT 'DRAFT',
    total_input_weight   DECIMAL(10,2) DEFAULT 0,
    total_output_weight  DECIMAL(10,2) DEFAULT 0,
    weight_loss_percent  DECIMAL(5,2) DEFAULT 0,
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS intakes (
    id           VARCHAR(64) PRIMARY KEY,
    harvest_id   VARCHAR(64) NOT NULL,
    pickup_id    VARCHAR(64),
    warehouse_id VARCHAR(80),
    coffee_type  VARCHAR(20) NOT NULL,
    origin_code  VARCHAR(10) NOT NULL,
    quantity     DECIMAL(10,2) NOT NULL,
    status       VARCHAR(20) DEFAULT 'UNASSIGNED',
    batch_id     VARCHAR(80) REFERENCES production_batches(id),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roast_runs (
    id            VARCHAR(64) PRIMARY KEY,
    batch_id      VARCHAR(80) NOT NULL REFERENCES production_batches(id),
    run_number    INT NOT NULL,
    input_weight  DECIMAL(10,2) NOT NULL,
    output_weight DECIMAL(10,2) NOT NULL,
    status        VARCHAR(20) DEFAULT 'COMPLETED',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inventories (
    id                 VARCHAR(64) PRIMARY KEY,
    coffee_type        VARCHAR(20) NOT NULL,
    origin_code        VARCHAR(10) NOT NULL,
    warehouse_id       VARCHAR(80),
    sku                VARCHAR(255) UNIQUE NOT NULL,
    available_quantity DECIMAL(12,2) DEFAULT 0,
    updated_at         TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inbox_events (
    id          VARCHAR(64) PRIMARY KEY,
    message_id  VARCHAR(255) UNIQUE NOT NULL,
    event_type  VARCHAR(160),
    processed_at TIMESTAMP WITH TIME ZONE
);
-- Warehouse Service Seed Data
INSERT INTO inventories (id, coffee_type, origin_code, warehouse_id, sku, available_quantity) VALUES
('INV_ARABICA_001', 'ARABICA', 'VN-LD', 'WAREHOUSE-HN-001', 'SKU-AR-VN-LD-001', 1000.00),
('INV_ROBUSTA_001', 'ROBUSTA', 'VN-DL', 'WAREHOUSE-HN-001', 'SKU-RB-VN-DL-001', 2000.00)
ON CONFLICT (id) DO NOTHING;
