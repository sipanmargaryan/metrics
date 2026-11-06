ALTER TABLE products DROP CONSTRAINT products_remaining_quantity_valid;
ALTER TABLE products DROP COLUMN initial_quantity;
ALTER TABLE products DROP COLUMN price;
ALTER TABLE products RENAME COLUMN remaining_quantity TO quantity;
ALTER TABLE products ADD CONSTRAINT products_quantity_nonnegative CHECK (quantity >= 0);
