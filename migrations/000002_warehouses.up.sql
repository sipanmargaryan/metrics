CREATE TABLE warehouses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL UNIQUE CHECK (btrim(name) <> '')
);

INSERT INTO warehouses (name)
SELECT DISTINCT COALESCE(NULLIF(btrim(warehouse), ''), 'Unassigned')
FROM orders;

ALTER TABLE orders ADD COLUMN warehouse_id BIGINT;
UPDATE orders AS o
SET warehouse_id = w.id
FROM warehouses AS w
WHERE w.name = COALESCE(NULLIF(btrim(o.warehouse), ''), 'Unassigned');

ALTER TABLE orders ALTER COLUMN warehouse_id SET NOT NULL;
ALTER TABLE orders ADD CONSTRAINT orders_warehouse_id_fkey
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE RESTRICT;
CREATE INDEX idx_orders_warehouse_id ON orders(warehouse_id);
ALTER TABLE orders DROP COLUMN warehouse;
