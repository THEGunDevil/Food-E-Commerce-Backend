-- +migrate down
-- Drop child table first
DROP TABLE IF EXISTS cart_items;

-- Then drop parent table
DROP TABLE IF EXISTS carts;