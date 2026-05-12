CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    metadata JSONB,
    retry_count INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'PENDING',
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_outbox_unprocessed ON outbox_events(created_at) WHERE status = 'PENDING';
