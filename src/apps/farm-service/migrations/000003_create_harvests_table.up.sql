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
('Cầu Đất Premium', 'CAU_DAT', 11.895, 108.538, 50.5, 'ARABICA', '00000000-0000-0000-0000-000000000001'),
('BMT Highlands', 'BUON_MA_THUOT', 12.666, 108.038, 120.0, 'ROBUSTA', '00000000-0000-0000-0000-000000000002')
ON CONFLICT DO NOTHING;
