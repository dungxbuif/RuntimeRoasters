-- Retail Service development schema. Business data is owned by Admin master
-- seed and runtime demo scenarios, never by migrations.

CREATE TABLE stores (
    id UUID PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(120) NOT NULL,
    city VARCHAR(80) NOT NULL,
    address VARCHAR(255) NOT NULL,
    manager_id VARCHAR(64),
    manager_email VARCHAR(160),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_stores_city ON stores(city);

CREATE TABLE menus (
    id VARCHAR(64) PRIMARY KEY,
    code VARCHAR(64) UNIQUE NOT NULL,
    name VARCHAR(120) NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL CHECK (status IN ('DRAFT', 'ACTIVE', 'INACTIVE')),
    effective_from TIMESTAMPTZ,
    effective_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (effective_to IS NULL OR effective_from IS NULL OR effective_to > effective_from)
);

CREATE TABLE menu_items (
    id VARCHAR(64) PRIMARY KEY,
    menu_id VARCHAR(64) NOT NULL REFERENCES menus(id),
    code VARCHAR(64) UNIQUE NOT NULL,
    product_group_code VARCHAR(64) NOT NULL,
    category_code VARCHAR(32) NOT NULL,
    size VARCHAR(8) NOT NULL CHECK (size IN ('S', 'M', 'L')),
    name VARCHAR(120) NOT NULL,
    description TEXT,
    stock_sku VARCHAR(80) NOT NULL,
    coffee_type VARCHAR(30) NOT NULL CHECK (coffee_type IN ('ARABICA', 'ROBUSTA')),
    price DECIMAL(12,2) NOT NULL CHECK (price > 0),
    consumption_quantity DECIMAL(12,3) NOT NULL CHECK (consumption_quantity > 0),
    consumption_unit VARCHAR(16) NOT NULL DEFAULT 'GRAM' CHECK (consumption_unit = 'GRAM'),
    image_url TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (menu_id, product_group_code, size)
);
CREATE INDEX idx_menu_items_menu_active_order ON menu_items(menu_id, active, display_order);
CREATE INDEX idx_menu_items_stock_sku ON menu_items(stock_sku);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL REFERENCES stores(id),
    items JSONB NOT NULL,
    total_amount DECIMAL(12,2) NOT NULL CHECK (total_amount >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    idempotency_key VARCHAR(128) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_orders_store_id ON orders(store_id);
CREATE INDEX idx_orders_status ON orders(status);

CREATE TABLE inventory_lots (
    id VARCHAR(64) PRIMARY KEY,
    store_id UUID NOT NULL REFERENCES stores(id),
    stock_sku VARCHAR(80) NOT NULL,
    source_batch_id VARCHAR(128) NOT NULL,
    source_harvest_id VARCHAR(128) NOT NULL,
    source_warehouse_id VARCHAR(80) NOT NULL,
    source_order_id UUID REFERENCES orders(id),
    received_quantity DECIMAL(12,3) NOT NULL CHECK (received_quantity >= 0),
    available_quantity DECIMAL(12,3) NOT NULL CHECK (available_quantity >= 0),
    unit VARCHAR(16) NOT NULL DEFAULT 'GRAM' CHECK (unit = 'GRAM'),
    status VARCHAR(20) NOT NULL CHECK (status IN ('AVAILABLE', 'DEPLETED', 'BLOCKED')),
    received_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (available_quantity <= received_quantity),
    CHECK (expires_at IS NULL OR expires_at > received_at)
);
CREATE INDEX idx_inventory_lots_lookup ON inventory_lots(store_id, stock_sku, status, received_at);
CREATE INDEX idx_inventory_lots_lineage ON inventory_lots(source_batch_id, source_harvest_id, source_warehouse_id);

CREATE TABLE sales (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL REFERENCES stores(id),
    invoice_no VARCHAR(64) UNIQUE NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('COMPLETED', 'VOIDED')),
    subtotal DECIMAL(12,2) NOT NULL CHECK (subtotal >= 0),
    total_amount DECIMAL(12,2) NOT NULL CHECK (total_amount >= 0),
    idempotency_key VARCHAR(128) UNIQUE NOT NULL,
    sold_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sales_store_sold ON sales(store_id, sold_at DESC);

CREATE TABLE sale_items (
    id UUID PRIMARY KEY,
    sale_id UUID NOT NULL REFERENCES sales(id),
    store_id UUID NOT NULL REFERENCES stores(id),
    menu_item_id VARCHAR(64) NOT NULL REFERENCES menu_items(id),
    inventory_lot_id VARCHAR(64) NOT NULL REFERENCES inventory_lots(id),
    product_id VARCHAR(128) UNIQUE NOT NULL,
    trace_code VARCHAR(128) UNIQUE NOT NULL,
    product_name VARCHAR(120) NOT NULL,
    sku VARCHAR(80) NOT NULL,
    size VARCHAR(8) NOT NULL,
    unit_price DECIMAL(12,2) NOT NULL CHECK (unit_price > 0),
    consumed_quantity DECIMAL(12,3) NOT NULL CHECK (consumed_quantity > 0),
    consumed_unit VARCHAR(16) NOT NULL DEFAULT 'GRAM' CHECK (consumed_unit = 'GRAM'),
    public_url TEXT NOT NULL,
    sold_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (trace_code = product_id)
);
CREATE INDEX idx_sale_items_store_sold ON sale_items(store_id, sold_at DESC);

CREATE TABLE stock_movements (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL REFERENCES stores(id),
    inventory_lot_id VARCHAR(64) NOT NULL REFERENCES inventory_lots(id),
    stock_sku VARCHAR(80) NOT NULL,
    movement_type VARCHAR(20) NOT NULL CHECK (movement_type IN ('RECEIVED', 'SOLD', 'ADJUSTMENT')),
    quantity_delta DECIMAL(12,3) NOT NULL CHECK (quantity_delta <> 0),
    unit VARCHAR(16) NOT NULL DEFAULT 'GRAM' CHECK (unit = 'GRAM'),
    reference_type VARCHAR(30) NOT NULL CHECK (reference_type IN ('SUPPLY_ORDER', 'SALE_ITEM', 'DEMO_SCENARIO', 'ADJUSTMENT')),
    reference_id VARCHAR(128) NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (movement_type, reference_type, reference_id, inventory_lot_id)
);
CREATE INDEX idx_stock_movements_lot_time ON stock_movements(inventory_lot_id, occurred_at);

CREATE TABLE store_menu_inventories (
    store_id UUID NOT NULL REFERENCES stores(id),
    menu_item_id VARCHAR(64) NOT NULL REFERENCES menu_items(id),
    available_units BIGINT NOT NULL CHECK (available_units >= 0),
    active_lot_count BIGINT NOT NULL CHECK (active_lot_count >= 0),
    version BIGINT NOT NULL DEFAULT 1,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (store_id, menu_item_id)
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(120) NOT NULL,
    topic VARCHAR(160) NOT NULL,
    key VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    trace_parent VARCHAR(128),
    trace_state VARCHAR(512),
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_outbox_status_created ON outbox_events(status, created_at);

CREATE TABLE inbox_events (
    id UUID PRIMARY KEY,
    message_id VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(120) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL
);
