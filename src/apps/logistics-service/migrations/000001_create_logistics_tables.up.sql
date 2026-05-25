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
-- Farms
('FARM-CAUDAT-001', 'Nông trại K''Ho (Đà Lạt)', 'FARM', 11.9404, 108.4442),
('FARM-CAUDAT-002', 'Nông trại Arabica Cầu Đất', 'FARM', 11.8950, 108.5380),
('FARM-CAUDAT-003', 'Nông trại Pacamara (Trạm Hành)', 'FARM', 11.8500, 108.5000),
('FARM-BMT-001', 'Nông trại Aeroco (Đắk Lắk)', 'FARM', 12.6660, 108.0380),
('FARM-BMT-002', 'Làng cà phê Trung Nguyên', 'FARM', 12.7000, 108.0500),
('FARM-PLEIKU-001', 'Nông trại Chư Sê (Gia Lai)', 'FARM', 14.0100, 108.0400),
-- Warehouses
('WAREHOUSE-HN-001', 'Xưởng rang Hà Nội', 'WAREHOUSE', 21.010, 105.530),
('WAREHOUSE-HCM-001', 'Xưởng rang Sóng Thần', 'WAREHOUSE', 10.880, 106.750),
('WAREHOUSE-DN-001', 'Xưởng rang Hòa Khánh', 'WAREHOUSE', 16.080, 108.150),
-- Retailers
('11111111-1111-1111-1111-111111111101', 'Cửa hàng Hoàn Kiếm', 'RETAILER', 21.028, 105.852),
('11111111-1111-1111-1111-111111111102', 'Cửa hàng Cầu Giấy', 'RETAILER', 21.036, 105.783),
('11111111-1111-1111-1111-111111111103', 'Cửa hàng Quận 1', 'RETAILER', 10.776, 106.700),
('11111111-1111-1111-1111-111111111104', 'Cửa hàng Quận 7', 'RETAILER', 10.729, 106.721),
('11111111-1111-1111-1111-111111111105', 'Cửa hàng Hải Châu', 'RETAILER', 16.066, 108.216)
ON CONFLICT (id) DO NOTHING;

INSERT INTO vehicles (id, plate_number, type, capacity_kg) VALUES
('VEHICLE-DEMO-001', '51C-12345', 'TRUCK', 5000),
('VEHICLE-DEMO-002', '29H-67890', 'TRUCK', 5000),
('VEHICLE-DEMO-003', '43A-11111', 'VAN', 1000)
ON CONFLICT (id) DO NOTHING;

INSERT INTO drivers (id, name, phone, vehicle_id, is_available) VALUES
('22222222-2222-2222-2222-222222222201', 'Driver Alpha', '0901234567', 'VEHICLE-DEMO-001', true),
('22222222-2222-2222-2222-222222222202', 'Driver Beta', '0907654321', 'VEHICLE-DEMO-002', true),
('22222222-2222-2222-2222-222222222203', 'Driver Gamma', '0900000000', 'VEHICLE-DEMO-003', true)
ON CONFLICT (id) DO NOTHING;
