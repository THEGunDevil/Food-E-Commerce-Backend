-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    menu_item_id,
    menu_item_name,
    quantity,
    unit_price,
    total_price,
    special_instructions
) VALUES (
    $1,$2,$3,$4,$5,$6,$7
)
RETURNING *;
-- name: GetOrderItemByID :one
SELECT *
FROM order_items
WHERE id = $1;
-- name: ListOrderItemsByOrder :many
SELECT *
FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC;
-- name: UpdateOrderItem :one
UPDATE order_items
SET
    menu_item_id = COALESCE($2, menu_item_id),
    menu_item_name = COALESCE($3, menu_item_name),
    quantity = COALESCE($4, quantity),
    unit_price = COALESCE($5, unit_price),
    total_price = COALESCE($6, total_price),
    special_instructions = COALESCE($7, special_instructions)
WHERE id = $1
RETURNING *;
-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;
