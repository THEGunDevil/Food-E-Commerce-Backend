-- name: AddCartItem :one
INSERT INTO cart_items (
    cart_id,
    menu_item_id,
    quantity,
    special_instructions
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetSubTotal :one
SELECT 
    COALESCE(
        SUM(
            CASE 
                WHEN mi.discount_price > 0 THEN mi.discount_price 
                ELSE mi.price 
            END * c.quantity
        ), 
        0
    )::DECIMAL(10, 2) AS subtotal
FROM cart_items c
JOIN menu_items mi ON c.menu_item_id = mi.id
WHERE c.cart_id = $1;

-- name: GetCartItemByCartAndMenuItem :one
SELECT *
FROM cart_items
WHERE cart_id = $1
  AND menu_item_id = $2;

-- name: ListCartItemsByCart :many
SELECT *
FROM cart_items
WHERE cart_id = $1
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

-- name: ClearCart :exec
DELETE FROM cart_items
WHERE cart_id = $1;

-- name: GetCartBySessionID :one
SELECT *
FROM carts
WHERE session_id = $1;
