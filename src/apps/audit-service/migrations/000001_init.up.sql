-- Audit Service Schema
CREATE TABLE IF NOT EXISTS audit_logs (
    id            UUID PRIMARY KEY,
    partition_key VARCHAR(120) NOT NULL,
    message_id    VARCHAR(255) UNIQUE NOT NULL,
    topic         VARCHAR(160) NOT NULL,
    store_id      VARCHAR(100),
    payload       JSONB NOT NULL,
    previous_hash VARCHAR(64),
    current_hash  VARCHAR(64) NOT NULL,
    occurred_at   TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_partition_key ON audit_logs(partition_key);
CREATE INDEX IF NOT EXISTS idx_audit_logs_topic ON audit_logs(topic);
CREATE INDEX IF NOT EXISTS idx_audit_logs_store_id ON audit_logs(store_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_occurred_at ON audit_logs(occurred_at);
