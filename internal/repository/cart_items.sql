-- name: AddCartItem :one
INSERT INTO cart_items (
    user_id,
    session_id,
    menu_item_id,
    quantity,
    special_instructions
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;
-- name: GetCartItemByIdentifierAndMenuItem :one
SELECT *
FROM cart_items
WHERE (user_id = $1 OR session_id = $1)
  AND menu_item_id = $2;

-- name: ListCartItemsByIdentifier :many
SELECT *
FROM cart_items
WHERE user_id = $1 OR session_id = $1
ORDER BY created_at ASC;

-- name: UpdateCartItem :one
UPDATE cart_items
SET
    quantity = $2,
    special_instructions = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: RemoveCartItem :exec
DELETE FROM cart_items
WHERE id = $1;
-- name: ClearCartByUser :exec
DELETE FROM cart_items
WHERE user_id = $1;
