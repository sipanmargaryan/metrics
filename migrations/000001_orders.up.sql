CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    name TEXT NOT NULL,
    quantity BIGINT NOT NULL DEFAULT 1,
    warehouse TEXT NOT NULL DEFAULT ''
);

-- Adopt tables created by the previous AutoMigrate setup as well.
ALTER TABLE orders ADD COLUMN IF NOT EXISTS quantity BIGINT NOT NULL DEFAULT 1;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS warehouse TEXT NOT NULL DEFAULT '';
