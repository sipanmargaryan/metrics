ALTER TABLE orders ADD COLUMN warehouse TEXT NOT NULL DEFAULT '';
UPDATE orders AS o SET warehouse = w.name
FROM warehouses AS w WHERE o.warehouse_id = w.id;
ALTER TABLE orders DROP COLUMN warehouse_id;
DROP TABLE warehouses;
