ALTER Table cart_items
ADD COLUMN session_id UUID;

ALTER TABLE cart_items
ADD CONSTRAINT cart_items_owner_check
CHECK (
    (user_id IS NOT NULL AND session_id IS NULL)
 OR (user_id IS NULL AND session_id IS NOT NULL)
);