-- name: CreateCategory :one
INSERT INTO categories (
  name,
  slug,
  description,
  cat_image_url,
  cat_image_public_id,
  display_order,
  is_active
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
-- name: GetCategoryByID :one
SELECT * FROM categories
WHERE id = $1;
-- name: GetCategoryBySlug :one
SELECT * FROM categories
WHERE slug = $1;
-- name: ListCategories :many
SELECT *
FROM categories
ORDER BY display_order ASC, created_at DESC;
-- name: ListActiveCategories :many
SELECT *
FROM categories
WHERE is_active = true
ORDER BY display_order ASC;
-- name: UpdateCategory :one
UPDATE categories
SET
  name = COALESCE($2, name),
  slug = COALESCE($3, slug),
  description = COALESCE($4, description),
  cat_image_url = COALESCE($5, cat_image_url),
  cat_image_public_id = COALESCE($6, cat_image_public_id),
  display_order = COALESCE($7, display_order),
  is_active = COALESCE($8, is_active),
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;
-- name: DeactivateCategory :exec
UPDATE categories
SET is_active = false,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
-- name: UpdateCategoryOrder :exec
UPDATE categories
SET display_order = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;
