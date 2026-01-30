-- name: CreateCartWithUser :one
INSERT INTO carts (session_id, user_id)
VALUES ($1, $2)
RETURNING id, session_id, user_id, status, created_at, updated_at;

-- name: UpdateCartWithUser :exec
UPDATE carts
SET user_id = $1
WHERE session_id = $2 AND user_id IS NULL;


