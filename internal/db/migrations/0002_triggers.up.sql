-- +migrate Up

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS
$$
BEGIN
    NEW.order_number := 'ORD-' || to_char(CURRENT_TIMESTAMP, 'YYMMDD') || '-' ||
                        LPAD((floor(random() * 9999))::text, 4, '0');
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER generate_order_number_trigger
BEFORE INSERT ON orders
FOR EACH ROW
EXECUTE FUNCTION generate_order_number();

-- 3️⃣ Function to update menu item average rating
CREATE OR REPLACE FUNCTION update_menu_item_stats()
RETURNS TRIGGER AS
$$
BEGIN
    UPDATE menu_items
    SET average_rating = (
        SELECT COALESCE(AVG(rating),0)
        FROM reviews
        WHERE menu_item_id = NEW.menu_item_id
        AND is_approved = TRUE
    )
    WHERE id = NEW.menu_item_id;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Triggers for tables that have updated_at
CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_items_updated_at
BEFORE UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_addresses_updated_at
BEFORE UPDATE ON user_addresses
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_promotions_updated_at
BEFORE UPDATE ON promotions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_item_stats_trigger
AFTER INSERT OR UPDATE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_menu_item_stats();
