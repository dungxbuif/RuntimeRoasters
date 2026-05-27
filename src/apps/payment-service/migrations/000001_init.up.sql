-- Payment Service Schema
CREATE TABLE IF NOT EXISTS payments (
    id           UUID PRIMARY KEY,
    order_id     UUID NOT NULL UNIQUE,
    store_id     UUID NOT NULL,
    provider     VARCHAR(40) NOT NULL,
    provider_ref VARCHAR(160) NOT NULL UNIQUE,
    refund_ref   VARCHAR(160),
    items        JSONB NOT NULL DEFAULT '[]'::jsonb,
    amount       DECIMAL(12,2) NOT NULL,
    currency     VARCHAR(10) DEFAULT 'USD',
    status       VARCHAR(24) DEFAULT 'PENDING',
    checkout_url VARCHAR(512),
    simulated    BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Compatibility for existing local/demo databases created before webhook-gated paid orders.
ALTER TABLE payments ADD COLUMN IF NOT EXISTS store_id UUID;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS provider VARCHAR(40);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS provider_ref VARCHAR(160);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS refund_ref VARCHAR(160);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS items JSONB NOT NULL DEFAULT '[]';
ALTER TABLE payments ADD COLUMN IF NOT EXISTS checkout_url VARCHAR(512);
ALTER TABLE payments ADD COLUMN IF NOT EXISTS simulated BOOLEAN NOT NULL DEFAULT true;
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_store_id ON payments(store_id);
CREATE INDEX IF NOT EXISTS idx_payments_provider ON payments(provider);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_refund_ref ON payments(refund_ref);
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_provider_ref ON payments(provider_ref);

CREATE TABLE IF NOT EXISTS inbox_events (
    id           UUID PRIMARY KEY,
    message_id   VARCHAR(255) UNIQUE NOT NULL,
    event_type   VARCHAR(120) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS webhook_events (
    id           UUID PRIMARY KEY,
    provider     VARCHAR(40) NOT NULL,
    event_id     VARCHAR(160) NOT NULL,
    provider_ref VARCHAR(160) NOT NULL,
    status       VARCHAR(40) NOT NULL,
    reason       VARCHAR(512),
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);
-- Compatibility for existing webhook_events table created with only provider/payload.
ALTER TABLE webhook_events ADD COLUMN IF NOT EXISTS event_id VARCHAR(160);
ALTER TABLE webhook_events ADD COLUMN IF NOT EXISTS provider_ref VARCHAR(160);
ALTER TABLE webhook_events ADD COLUMN IF NOT EXISTS status VARCHAR(40);
ALTER TABLE webhook_events ADD COLUMN IF NOT EXISTS reason VARCHAR(512);
CREATE UNIQUE INDEX IF NOT EXISTS idx_webhook_events_provider_event ON webhook_events(provider, event_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_provider_ref ON webhook_events(provider_ref);

CREATE TABLE IF NOT EXISTS payment_webhook_keys (
    id             UUID PRIMARY KEY,
    provider       VARCHAR(40) NOT NULL,
    key_name       VARCHAR(120) NOT NULL,
    signing_secret VARCHAR(255) NOT NULL,
    active         BOOLEAN NOT NULL DEFAULT true,
    demo           BOOLEAN NOT NULL DEFAULT true,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, key_name)
);

INSERT INTO payment_webhook_keys (id, provider, key_name, signing_secret, active, demo)
VALUES (
    '00000000-0000-0000-0000-000000000451',
    'STRIPE',
    'demo-stripe-local',
    'whsec_rr_demo_stripe_local',
    true,
    true
)
ON CONFLICT (provider, key_name) DO NOTHING;
