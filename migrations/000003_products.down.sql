ALTER TABLE orders ADD COLUMN name TEXT;
ALTER TABLE orders ADD COLUMN warehouse_id BIGINT;
UPDATE orders AS o SET name = p.name, warehouse_id = p.warehouse_id
FROM products AS p WHERE p.id = o.product_id;
ALTER TABLE orders ALTER COLUMN name SET NOT NULL;
ALTER TABLE orders ALTER COLUMN warehouse_id SET NOT NULL;
ALTER TABLE orders ADD CONSTRAINT orders_warehouse_id_fkey
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE RESTRICT;
CREATE INDEX idx_orders_warehouse_id ON orders(warehouse_id);
ALTER TABLE orders DROP COLUMN product_id;
DROP TABLE products;
