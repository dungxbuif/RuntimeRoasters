-- Define Enums
DO $$ BEGIN
    CREATE TYPE harvest_status_enum AS ENUM ('NEW', 'PROCESSING', 'COMPLETED');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE coffee_type_enum AS ENUM ('ARABICA', 'ROBUSTA', 'CHERRY', 'CULI');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS harvests (
    id BIGSERIAL PRIMARY KEY,
    farm_id BIGINT NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    owner_id UUID NOT NULL,
    coffee_type coffee_type_enum NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    harvest_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status harvest_status_enum NOT NULL DEFAULT 'NEW',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_harvests_owner_id ON harvests(owner_id);

CREATE INDEX IF NOT EXISTS idx_harvests_farm_id ON harvests(farm_id);
-- Farm Service Seed Data
INSERT INTO farms (name, location, latitude, longitude, area, coffee_type, owner_id) VALUES
('K''Ho Coffee Farm', 'CAU_DAT', 11.9404, 108.4442, 15.5, 'ARABICA', '00000000-0000-0000-0000-000000000001'),
('Cau Dat Arabica', 'CAU_DAT', 11.8950, 108.5380, 45.0, 'ARABICA', '00000000-0000-0000-0000-000000000001'),
('Son Pacamara Farm', 'CAU_DAT', 11.8500, 108.5000, 12.0, 'ARABICA', '00000000-0000-0000-0000-000000000001'),
('Aeroco Coffee', 'BUON_MA_THUOT', 12.6660, 108.0380, 20.0, 'ROBUSTA', '00000000-0000-0000-0000-000000000002'),
('Trung Nguyen Village', 'BUON_MA_THUOT', 12.7000, 108.0500, 5.0, 'ROBUSTA', '00000000-0000-0000-0000-000000000002'),
('Chư Sê Estate', 'PLEIKU', 14.0100, 108.0400, 30.0, 'ROBUSTA', '00000000-0000-0000-0000-000000000003')
ON CONFLICT DO NOTHING;
