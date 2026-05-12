CREATE TABLE IF NOT EXISTS farms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location TEXT NOT NULL,
    area DECIMAL(10,2) NOT NULL CHECK (area > 0),
    coffee_type VARCHAR(100),
    owner_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_farms_owner_id ON farms(owner_id);
