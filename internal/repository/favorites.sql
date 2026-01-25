-- name: CreateFavorite :one
INSERT INTO favorites (user_id, menu_item_id)
VALUES ($1, $2)
RETURNING *;
-- name: DeleteFavorite :exec
DELETE FROM favorites
WHERE user_id = $1 AND menu_item_id = $2;
-- name: GetFavorite :one
SELECT *
FROM favorites
WHERE user_id = $1 AND menu_item_id = $2;
-- name: ListFavoritesByUser :many
SELECT f.id, f.menu_item_id, m.name, m.price, mi.image_url, f.created_at
FROM favorites f
JOIN menu_items m ON m.id = f.menu_item_id
JOIN menu_item_images mi ON mi.menu_item_id = m.id
WHERE f.user_id = $1
ORDER BY f.created_at DESC;
