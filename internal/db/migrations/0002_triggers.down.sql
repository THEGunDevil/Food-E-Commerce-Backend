-- +migrate down
-- ------------------------------
-- TRIGGERS
-- ------------------------------
DROP TRIGGER IF EXISTS update_menu_item_stats_trigger ON reviews;
DROP TRIGGER IF EXISTS generate_order_number_trigger ON orders;
DROP TRIGGER IF EXISTS update_promotions_updated_at ON promotions;
DROP TRIGGER IF EXISTS update_user_addresses_updated_at ON user_addresses;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
DROP TRIGGER IF EXISTS update_menu_items_updated_at ON menu_items;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
