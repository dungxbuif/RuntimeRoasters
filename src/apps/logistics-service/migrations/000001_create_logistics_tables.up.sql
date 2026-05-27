-- 000001_create_logistics_tables.up.sql

CREATE TABLE IF NOT EXISTS vehicles (
    id VARCHAR(64) PRIMARY KEY,
    plate_number VARCHAR(32) NOT NULL UNIQUE,
    type VARCHAR(32) NOT NULL,
    capacity_kg DECIMAL(10, 2),
    home_warehouse_id VARCHAR(80),
    status VARCHAR(24) NOT NULL DEFAULT 'IDLE',
    current_latitude DECIMAL(10, 8) DEFAULT 0,
    current_longitude DECIMAL(11, 8) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS home_warehouse_id VARCHAR(80);
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS status VARCHAR(24) NOT NULL DEFAULT 'IDLE';
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS current_latitude DECIMAL(10, 8) DEFAULT 0;
ALTER TABLE vehicles ADD COLUMN IF NOT EXISTS current_longitude DECIMAL(11, 8) DEFAULT 0;

CREATE TABLE IF NOT EXISTS drivers (
    id UUID PRIMARY KEY,
    user_id VARCHAR(160),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    vehicle_id VARCHAR(64) REFERENCES vehicles(id),
    status VARCHAR(24) NOT NULL DEFAULT 'IDLE',
    current_shipment_id VARCHAR(80),
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS user_id VARCHAR(160);
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS status VARCHAR(24) NOT NULL DEFAULT 'IDLE';
ALTER TABLE drivers ADD COLUMN IF NOT EXISTS current_shipment_id VARCHAR(80);

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
    type VARCHAR(32) NOT NULL DEFAULT 'RETAIL_DELIVERY',
    order_id VARCHAR(64) NOT NULL,
    driver_id UUID REFERENCES drivers(id),
    status VARCHAR(32) NOT NULL,
    current_leg VARCHAR(24),
    route_id VARCHAR(80),
    origin_location_id VARCHAR(80),
    destination_location_id VARCHAR(80),
    farm_id VARCHAR(80),
    harvest_id VARCHAR(80),
    warehouse_id VARCHAR(80),
    batch_id VARCHAR(80),
    origin_warehouse_id VARCHAR(80),
    destination_store_id VARCHAR(80),
    store_id VARCHAR(80),
    vehicle_id VARCHAR(80),
    destination_address TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    assigned_at TIMESTAMPTZ,
    departed_at TIMESTAMPTZ,
    arrived_at_origin_at TIMESTAMPTZ,
    loaded_at TIMESTAMPTZ,
    departed_origin_at TIMESTAMPTZ,
    arrived_destination_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    return_started_at TIMESTAMPTZ,
    returned_at TIMESTAMPTZ
);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS type VARCHAR(32) NOT NULL DEFAULT 'RETAIL_DELIVERY';
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS current_leg VARCHAR(24);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS route_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS origin_location_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS destination_location_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS farm_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS harvest_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS warehouse_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS batch_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS origin_warehouse_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS destination_store_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS store_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS vehicle_id VARCHAR(80);
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS assigned_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS departed_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS arrived_at_origin_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS loaded_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS departed_origin_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS arrived_destination_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS return_started_at TIMESTAMPTZ;
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS returned_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS processed_kafka_messages (
    msg_key VARCHAR(255) PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    processed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_shipments_order_id ON shipments(order_id);
CREATE INDEX IF NOT EXISTS idx_shipments_status ON shipments(status);
CREATE INDEX IF NOT EXISTS idx_shipments_destination_store_id ON shipments(destination_store_id);
CREATE INDEX IF NOT EXISTS idx_shipments_driver_id ON shipments(driver_id);
CREATE INDEX IF NOT EXISTS idx_drivers_available ON drivers(is_available);
CREATE INDEX IF NOT EXISTS idx_drivers_user_id ON drivers(user_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_home_warehouse_id ON vehicles(home_warehouse_id);
CREATE INDEX IF NOT EXISTS idx_locations_type ON locations(type);

INSERT INTO vehicles (id, plate_number, type, capacity_kg, home_warehouse_id, status) VALUES
('VEHICLE-DEMO-001', 'HN-51A-001', 'TRUCK', 5000, 'WAREHOUSE-HN-001', 'IDLE'),
('VEHICLE-DEMO-002', 'HN-51A-002', 'VAN', 1000, 'WAREHOUSE-HN-001', 'IDLE'),
('VEHICLE-DEMO-003', 'HCM-51A-003', 'TRUCK', 5000, 'WAREHOUSE-HCM-001', 'IDLE')
ON CONFLICT (id) DO UPDATE SET
    plate_number = EXCLUDED.plate_number,
    type = EXCLUDED.type,
    capacity_kg = EXCLUDED.capacity_kg,
    home_warehouse_id = EXCLUDED.home_warehouse_id,
    status = EXCLUDED.status,
    updated_at = NOW();

INSERT INTO drivers (id, user_id, name, phone, vehicle_id, status, is_available) VALUES
('22222222-2222-2222-2222-222222222201', 'driver@runtimeroasters.com', 'Driver Alpha', '0901234567', 'VEHICLE-DEMO-001', 'IDLE', true),
('22222222-2222-2222-2222-222222222202', 'driver.beta@runtimeroasters.com', 'Driver Beta', '0907654321', 'VEHICLE-DEMO-002', 'IDLE', true),
('22222222-2222-2222-2222-222222222203', 'driver.hcm@runtimeroasters.com', 'Driver Gamma', '0900000000', 'VEHICLE-DEMO-003', 'IDLE', true)
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    phone = EXCLUDED.phone,
    vehicle_id = EXCLUDED.vehicle_id,
    status = EXCLUDED.status,
    is_available = EXCLUDED.is_available,
    updated_at = NOW();

INSERT INTO locations (id, name, type, lat, lng) VALUES
('FARM-CAUDAT-001', 'KHo Coffee Farm', 'FARM', 11.9404, 108.4442),
('FARM-CAUDAT-002', 'Cau Dat Arabica Farm', 'FARM', 11.8950, 108.5380),
('FARM-CAUDAT-003', 'Son Pacamara Farm', 'FARM', 11.8500, 108.5000),
('FARM-BMT-001', 'Aeroco Coffee Farm', 'FARM', 12.6660, 108.0380),
('FARM-BMT-002', 'Trung Nguyen Village', 'FARM', 12.7000, 108.0500),
('FARM-PLEIKU-001', 'Chu Se Estate', 'FARM', 14.0100, 108.0400),
('WAREHOUSE-HN-001', 'Hanoi Roastery', 'WAREHOUSE', 21.010, 105.530),
('WAREHOUSE-HCM-001', 'Song Than Roastery', 'WAREHOUSE', 10.880, 106.750),
('WAREHOUSE-DN-001', 'Hoa Khanh Roastery', 'WAREHOUSE', 16.080, 108.150),
('11111111-1111-1111-1111-111111111101', 'Hoan Kiem Store', 'RETAILER', 21.028, 105.852),
('11111111-1111-1111-1111-111111111102', 'Cau Giay Store', 'RETAILER', 21.036, 105.783),
('11111111-1111-1111-1111-111111111103', 'District 1 Store', 'RETAILER', 10.776, 106.700),
('11111111-1111-1111-1111-111111111104', 'District 7 Store', 'RETAILER', 10.729, 106.721),
('11111111-1111-1111-1111-111111111105', 'Hai Chau Store', 'RETAILER', 16.066, 108.216)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    lat = EXCLUDED.lat,
    lng = EXCLUDED.lng,
    updated_at = NOW();
