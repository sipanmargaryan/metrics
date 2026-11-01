ALTER TABLE products ADD COLUMN quantity BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD CONSTRAINT products_quantity_nonnegative CHECK (quantity >= 0);
