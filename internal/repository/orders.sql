-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    delivery_address_id,
    delivery_type,
    delivery_address,
    customer_name,
    customer_phone,
    customer_email,
    subtotal,
    discount_amount,
    delivery_fee,
    vat_amount,
    total_amount,
    payment_method,
    payment_status,
    transaction_id,
    order_status,
    delivery_person_id,
    estimated_delivery,
    actual_delivery,
    special_instructions,
    cancelled_reason
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21
)
RETURNING *;
-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;
-- name: ListOrdersByUser :many
SELECT *
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC;
-- name: UpdateOrder :one
UPDATE orders
SET
    delivery_type = COALESCE($2, delivery_type),
    delivery_address = COALESCE($3, delivery_address),
    customer_name = COALESCE($4, customer_name),
    customer_phone = COALESCE($5, customer_phone),
    customer_email = COALESCE($6, customer_email),
    subtotal = COALESCE($7, subtotal),
    discount_amount = COALESCE($8, discount_amount),
    delivery_fee = COALESCE($9, delivery_fee),
    vat_amount = COALESCE($10, vat_amount),
    total_amount = COALESCE($11, total_amount),
    payment_method = COALESCE($12, payment_method),
    payment_status = COALESCE($13, payment_status),
    transaction_id = COALESCE($14, transaction_id),
    order_status = COALESCE($15, order_status),
    delivery_person_id = COALESCE($16, delivery_person_id),
    estimated_delivery = COALESCE($17, estimated_delivery),
    actual_delivery = COALESCE($18, actual_delivery),
    special_instructions = COALESCE($19, special_instructions),
    cancelled_reason = COALESCE($20, cancelled_reason),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;
