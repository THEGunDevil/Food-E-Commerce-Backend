-- name: CreatePromotion :one
INSERT INTO promotions (
    title,
    description,
    discount_type,
    discount_value,
    min_order_amount,
    valid_from,
    valid_until,
    max_uses,
    used_count,
    is_active
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10
)
RETURNING *;
-- name: GetPromotionByID :one
SELECT *
FROM promotions
WHERE id = $1;
-- name: ListPromotions :many
SELECT *
FROM promotions
ORDER BY created_at DESC;
-- name: ListActivePromotions :many
SELECT *
FROM promotions
WHERE is_active = TRUE
  AND valid_from <= NOW()
  AND valid_until >= NOW()
ORDER BY valid_until ASC;
-- name: UpdatePromotion :one
UPDATE promotions
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    discount_type = COALESCE($4, discount_type),
    discount_value = COALESCE($5, discount_value),
    min_order_amount = COALESCE($6, min_order_amount),
    valid_from = COALESCE($7, valid_from),
    valid_until = COALESCE($8, valid_until),
    max_uses = COALESCE($9, max_uses),
    used_count = COALESCE($10, used_count),
    is_active = COALESCE($11, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: DeletePromotion :exec
DELETE FROM promotions
WHERE id = $1;
-- name: IncrementPromotionUsage :one
UPDATE promotions
SET used_count = used_count + 1
WHERE id = $1 AND (max_uses IS NULL OR used_count < max_uses)
RETURNING *;
