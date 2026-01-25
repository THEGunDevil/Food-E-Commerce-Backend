-- name: CreateReview :one
INSERT INTO reviews (
    user_id,
    menu_item_id,
    rating,
    comment
) VALUES (
    $1,$2,$3,$4
)
RETURNING *;
-- name: GetReviewByID :one
SELECT *
FROM reviews
WHERE id = $1;
-- name: ListMenuItemReviewsByMenuItemID :many
SELECT *
FROM reviews
WHERE menu_item_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
-- name: ListApprovedReviewsByMenuItem :many
SELECT *
FROM reviews
WHERE menu_item_id = $1
  AND is_approved = TRUE
ORDER BY created_at DESC;
-- name: ListReviewsByUser :many
SELECT *
FROM reviews
WHERE user_id = $1
ORDER BY created_at DESC;
-- name: UpdateReview :one
UPDATE reviews
SET
    rating = COALESCE($2, rating),
    comment = COALESCE($3, comment),
    is_approved = COALESCE($4, is_approved),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: DeleteReview :exec
DELETE FROM reviews
WHERE id = $1;
-- name: SetReviewApproval :one
UPDATE reviews
SET is_approved = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: CountMenuItemReviewsByMenuItemID :one
SELECT COUNT(*)
FROM reviews
WHERE menu_item_id = $1;
