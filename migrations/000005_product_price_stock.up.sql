ALTER TABLE products RENAME COLUMN quantity TO remaining_quantity;
ALTER TABLE products DROP CONSTRAINT products_quantity_nonnegative;
ALTER TABLE products ADD COLUMN price BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN initial_quantity BIGINT NOT NULL DEFAULT 0;

-- Stock was not decremented before this migration. Preserve existing quantities
-- as both the starting stock and the available stock; do not replay old orders.
UPDATE products SET initial_quantity = remaining_quantity;

ALTER TABLE products ADD CONSTRAINT products_price_nonnegative CHECK (price >= 0);
ALTER TABLE products ADD CONSTRAINT products_initial_quantity_nonnegative CHECK (initial_quantity >= 0);
ALTER TABLE products ADD CONSTRAINT products_remaining_quantity_valid
    CHECK (remaining_quantity >= 0 AND remaining_quantity <= initial_quantity);
