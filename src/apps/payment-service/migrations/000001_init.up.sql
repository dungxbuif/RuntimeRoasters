-- Payment Service Schema
CREATE TABLE IF NOT EXISTS payments (
    id             UUID PRIMARY KEY,
    order_id       UUID NOT NULL,
    amount         DECIMAL(12,2) NOT NULL,
    currency       VARCHAR(10) DEFAULT 'VND',
    status         VARCHAR(20) DEFAULT 'PENDING',
    payment_method VARCHAR(40),
    transaction_id VARCHAR(100),
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);

CREATE TABLE IF NOT EXISTS inbox_events (
    id           UUID PRIMARY KEY,
    message_id   VARCHAR(255) UNIQUE NOT NULL,
    event_type   VARCHAR(120) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS webhook_events (
    id           UUID PRIMARY KEY,
    provider     VARCHAR(40) NOT NULL,
    payload      JSONB NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
