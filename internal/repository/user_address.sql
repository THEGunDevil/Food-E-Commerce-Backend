-- name: CreateUserAddress :one
INSERT INTO user_addresses (
    user_id,
    label,
    address_line1,
    address_line2,
    area,
    city,
    postal_code,
    latitude,
    longitude,
    is_default
) VALUES (
    $1, $2, $3, sqlc.narg(address_line2), $4,
    COALESCE(sqlc.narg(city), 'Dhaka'), sqlc.narg(postal_code), sqlc.narg(latitude), sqlc.narg(longitude),
    COALESCE(sqlc.narg(is_default), false)
) RETURNING *;

-- name: GetUserAddressByID :one
SELECT * FROM user_addresses WHERE id = $1 LIMIT 1;

-- name: ListUserAddresses :many
SELECT * FROM user_addresses 
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateUserAddress :one
UPDATE user_addresses
SET
    label        = COALESCE(sqlc.narg(label), label),
    address_line1 = COALESCE(sqlc.narg(address_line1), address_line1),
    address_line2 = COALESCE(sqlc.narg(address_line2), address_line2),
    area         = COALESCE(sqlc.narg(area), area),
    city         = COALESCE(sqlc.narg(city), city),
    postal_code  = COALESCE(sqlc.narg(postal_code), postal_code),
    latitude     = COALESCE(sqlc.narg(latitude), latitude),
    longitude    = COALESCE(sqlc.narg(longitude), longitude),
    is_default   = COALESCE(sqlc.narg(is_default), is_default),
    updated_at   = CURRENT_TIMESTAMP
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUserAddress :exec
DELETE FROM user_addresses WHERE id = $1;

-- name: SetDefaultUserAddress :exec
UPDATE user_addresses
SET is_default = false
WHERE user_id = $1 AND id != $2;

UPDATE user_addresses
SET is_default = true
WHERE id = $2;

-- name: CountUserAddresses :one
SELECT COUNT(*) FROM user_addresses
WHERE user_id = $1;

-- name: SearchUserAddresses :many
SELECT * FROM user_addresses 
WHERE user_id = $1 AND (
    label ILIKE '%' || $2 || '%' OR
    address_line1 ILIKE '%' || $2 || '%' OR
    area ILIKE '%' || $2 || '%' OR
    city ILIKE '%' || $2 || '%'
)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;