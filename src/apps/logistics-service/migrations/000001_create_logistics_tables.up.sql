-- 000001_create_logistics_tables.up.sql

CREATE TABLE IF NOT EXISTS vehicles (
    id VARCHAR(64) PRIMARY KEY,
    plate_number VARCHAR(32) NOT NULL UNIQUE,
    type VARCHAR(32) NOT NULL,
    capacity_kg DECIMAL(10, 2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    vehicle_id VARCHAR(64) REFERENCES vehicles(id),
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS locations (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    driver_id UUID REFERENCES drivers(id),
    status VARCHAR(32) NOT NULL,
    destination_address TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    delivered_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS processed_kafka_messages (
    msg_key VARCHAR(255) PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_drivers_available ON drivers(is_available);
CREATE INDEX idx_locations_type ON locations(type);
-- Logistics Service Seed Data
INSERT INTO locations (id, name, type, lat, lng) VALUES
('FARM_CAU_DAT', 'Nông trại Cầu Đất (Arabica)', 'FARM', 11.895, 108.538),
('FARM_BMT', 'Nông trại Buôn Ma Thuột (Robusta)', 'FARM', 12.666, 108.038),
('FARM_PLEIKU', 'Nông trại Pleiku (Blend)', 'FARM', 14.010, 108.040),
('WH_SONG_THAN', 'Xưởng rang KCN Sóng Thần', 'WAREHOUSE', 10.880, 106.750),
('WH_HOA_LAC', 'Xưởng rang KCN Hòa Lạc', 'WAREHOUSE', 21.010, 105.530),
('WH_HOA_KHANH', 'Xưởng rang KCN Hòa Khánh', 'WAREHOUSE', 16.080, 108.150),
('RET_HCM', 'Cửa hàng Quận 1, TP.HCM', 'RETAILER', 10.776, 106.700),
('RET_HN', 'Cửa hàng Hoàn Kiếm, Hà Nội', 'RETAILER', 21.028, 105.852),
('RET_DN', 'Cửa hàng Hải Châu, Đà Nẵng', 'RETAILER', 16.066, 108.216)
ON CONFLICT (id) DO NOTHING;

INSERT INTO vehicles (id, plate_number, type, capacity_kg) VALUES
('VEH_TRUCK_001', '51C-12345', 'TRUCK', 5000),
('VEH_TRUCK_002', '29H-67890', 'TRUCK', 5000),
('VEH_VAN_001', '43A-11111', 'VAN', 1000)
ON CONFLICT (id) DO NOTHING;

INSERT INTO drivers (id, name, phone, vehicle_id, is_available) VALUES
('00000000-0000-0000-0000-000000000001', 'Driver Alpha', '0901234567', 'VEH_TRUCK_001', true),
('00000000-0000-0000-0000-000000000002', 'Driver Beta', '0907654321', 'VEH_TRUCK_002', true)
ON CONFLICT (id) DO NOTHING;
