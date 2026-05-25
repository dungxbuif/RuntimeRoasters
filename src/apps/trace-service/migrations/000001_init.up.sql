-- Trace Service Schema
CREATE TABLE IF NOT EXISTS trace_events (
    id           UUID PRIMARY KEY,
    message_id   VARCHAR(255) UNIQUE NOT NULL,
    topic        VARCHAR(160) NOT NULL,
    batch_id     VARCHAR(100),
    order_id     VARCHAR(100),
    shipment_id  VARCHAR(100),
    store_id     VARCHAR(100),
    harvest_id   VARCHAR(100),
    farm_id      VARCHAR(100),
    warehouse_id VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS trace_documents (
    id          UUID PRIMARY KEY,
    entity_id   VARCHAR(120) UNIQUE NOT NULL,
    entity_type VARCHAR(40) NOT NULL,
    store_id    VARCHAR(100),
    trace_ids   JSONB NOT NULL DEFAULT '[]',
    document    JSONB NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
