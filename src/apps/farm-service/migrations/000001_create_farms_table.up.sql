-- Define Enums
DO $$ BEGIN
    CREATE TYPE location_enum AS ENUM ('CAU_DAT', 'BUON_MA_THUOT', 'PLEIKU', 'GIA_NGHIA', 'KON_TUM');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE coffee_type_enum AS ENUM ('ARABICA', 'ROBUSTA', 'CHERRY', 'CULI');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS farms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location location_enum,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    area DECIMAL(10,2) NOT NULL CHECK (area > 0),
    coffee_type coffee_type_enum NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_farms_owner_id ON farms(owner_id);
