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
    warehouse_id VARCHAR(100),
    driver_id    VARCHAR(100),
    vehicle_id   VARCHAR(100),
    trace_id     VARCHAR(64),
    flow_id      VARCHAR(100),
    node_id      VARCHAR(100),
    edge_id      VARCHAR(100),
    pattern      VARCHAR(100),
    source_service VARCHAR(100),
    visibility   VARCHAR(40) DEFAULT 'public',
    display_payload JSONB NOT NULL DEFAULT '{}',
    payload      JSONB NOT NULL DEFAULT '{}',
    occurred_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Compatibility for existing local/demo databases created before topology projection fields.
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS driver_id VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS vehicle_id VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS trace_id VARCHAR(64);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS flow_id VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS node_id VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS edge_id VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS pattern VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS source_service VARCHAR(100);
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS visibility VARCHAR(40) DEFAULT 'public';
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS display_payload JSONB NOT NULL DEFAULT '{}';
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS payload JSONB NOT NULL DEFAULT '{}';
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS occurred_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE trace_events ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_trace_events_topic ON trace_events(topic);
CREATE INDEX IF NOT EXISTS idx_trace_events_flow_id ON trace_events(flow_id);
CREATE INDEX IF NOT EXISTS idx_trace_events_visibility ON trace_events(visibility);
CREATE INDEX IF NOT EXISTS idx_trace_events_occurred_at ON trace_events(occurred_at);
CREATE INDEX IF NOT EXISTS idx_trace_events_trace_id ON trace_events(trace_id);

CREATE TABLE IF NOT EXISTS trace_documents (
    id          UUID PRIMARY KEY,
    entity_id   VARCHAR(120) UNIQUE NOT NULL,
    entity_type VARCHAR(40) NOT NULL,
    store_id    VARCHAR(100),
    trace_ids   JSONB NOT NULL DEFAULT '[]',
    document    JSONB NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
