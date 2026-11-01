CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL CHECK (btrim(name) <> ''),
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    CONSTRAINT products_warehouse_name_key UNIQUE (warehouse_id, name)
);
CREATE INDEX idx_products_warehouse_id ON products(warehouse_id);

-- Preserve old order names as products in their original warehouses.
INSERT INTO products (name, warehouse_id)
SELECT DISTINCT name, warehouse_id FROM orders;

ALTER TABLE orders ADD COLUMN product_id BIGINT;
UPDATE orders AS o SET product_id = p.id
FROM products AS p WHERE p.name = o.name AND p.warehouse_id = o.warehouse_id;
ALTER TABLE orders ALTER COLUMN product_id SET NOT NULL;
ALTER TABLE orders ADD CONSTRAINT orders_product_id_fkey
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT;
CREATE INDEX idx_orders_product_id ON orders(product_id);
ALTER TABLE orders DROP COLUMN warehouse_id;
ALTER TABLE orders DROP COLUMN name;
