-- Retail Service Schema
CREATE TABLE IF NOT EXISTS stores (
    id            UUID PRIMARY KEY,
    name          VARCHAR(120) NOT NULL,
    city          VARCHAR(80) NOT NULL,
    address       VARCHAR(255) NOT NULL,
    manager_id    VARCHAR(64),
    manager_email VARCHAR(160),
    status        VARCHAR(20) DEFAULT 'ACTIVE',
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_stores_city ON stores(city);

CREATE TABLE IF NOT EXISTS orders (
    id              UUID PRIMARY KEY,
    store_id        UUID NOT NULL REFERENCES stores(id),
    items           JSONB NOT NULL,
    total_amount    DECIMAL(12,2) NOT NULL,
    status          VARCHAR(20) DEFAULT 'PENDING',
    idempotency_key VARCHAR(128) UNIQUE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_orders_store_id ON orders(store_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

CREATE TABLE IF NOT EXISTS outbox_events (
    id           UUID PRIMARY KEY,
    event_type   VARCHAR(120) NOT NULL,
    topic        VARCHAR(160) NOT NULL,
    key          VARCHAR(128) NOT NULL,
    payload      JSONB NOT NULL,
    trace_parent VARCHAR(128),
    trace_state  VARCHAR(512),
    status       VARCHAR(20) DEFAULT 'PENDING',
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS inbox_events (
    id          UUID PRIMARY KEY,
    message_id  VARCHAR(255) UNIQUE NOT NULL,
    event_type  VARCHAR(120) NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);

INSERT INTO stores (id, name, city, address, manager_email, status) VALUES
('11111111-1111-1111-1111-111111111101', 'Hoan Kiem Store', 'Hanoi', '2 Ly Thai To', 'mgr.hn.hoankiem@runtimeroasters.com', 'ACTIVE'),
('11111111-1111-1111-1111-111111111102', 'Cau Giay Store', 'Hanoi', '102 Tran Thai Tong', 'mgr.hn.caugiay@runtimeroasters.com', 'ACTIVE'),
('11111111-1111-1111-1111-111111111103', 'District 1 Store', 'Ho Chi Minh City', '45 Le Thanh Ton', 'mgr.hcm.q1@runtimeroasters.com', 'ACTIVE'),
('11111111-1111-1111-1111-111111111104', 'District 7 Store', 'Ho Chi Minh City', 'Phu My Hung', 'mgr.hcm.q7@runtimeroasters.com', 'ACTIVE'),
('11111111-1111-1111-1111-111111111105', 'Hai Chau Store', 'Da Nang', '15 Bach Dang', 'mgr.dn.haichau@runtimeroasters.com', 'ACTIVE')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    city = EXCLUDED.city,
    address = EXCLUDED.address,
    manager_email = EXCLUDED.manager_email,
    status = EXCLUDED.status,
    updated_at = CURRENT_TIMESTAMP;
