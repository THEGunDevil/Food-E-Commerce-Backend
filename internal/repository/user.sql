-- name: CreateUser :one
INSERT INTO users (
    email,
    phone,
    full_name,
    password_hash,
    bio,
    avatar_url
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
    email        = COALESCE(sqlc.narg(email), email),
    phone        = COALESCE(sqlc.narg(phone), phone),
    full_name    = COALESCE(sqlc.narg(full_name), full_name),
    bio          = COALESCE(sqlc.narg(bio), bio),
    avatar_url   = COALESCE(sqlc.narg(avatar_url), avatar_url),
    avatar_public_id   = COALESCE(sqlc.narg(avatar_public_id), avatar_public_id),
    role         = COALESCE(sqlc.narg(role), role),
    is_active    = COALESCE(sqlc.narg(is_active), is_active),
    is_verified  = COALESCE(sqlc.narg(is_verified), is_verified),
    updated_at   = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING *;


-- name: UpdateUserPassword :one
UPDATE users
SET 
    password_hash = $2,
    token_version = token_version + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: BanUser :one
UPDATE users
SET 
    is_banned = true,
    ban_reason = $2,
    ban_until = $3,
    is_permanent_ban = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UnbanUser :one
UPDATE users
SET 
    is_banned = false,
    ban_reason = NULL,
    ban_until = NULL,
    is_permanent_ban = false,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateTokenVersion :exec
UPDATE users
SET token_version = token_version + 1
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: GetActiveUsers :many
SELECT * FROM users 
WHERE is_active = true AND is_banned = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchUsers :many
SELECT * FROM users 
WHERE 
    email ILIKE '%' || $1 || '%' OR
    full_name ILIKE '%' || $1 || '%' OR
    phone ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;