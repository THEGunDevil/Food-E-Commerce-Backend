-- name: CreateDeliveryZone :one
INSERT INTO delivery_zones (
    zone_name,
    area_names,
    delivery_fee,
    min_delivery_time,
    max_delivery_time,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;
-- name: GetDeliveryZoneByID :one
SELECT *
FROM delivery_zones
WHERE id = $1;
-- name: ListDeliveryZones :many
SELECT *
FROM delivery_zones
ORDER BY created_at DESC;
-- name: ListActiveDeliveryZones :many
SELECT *
FROM delivery_zones
WHERE is_active = TRUE
ORDER BY zone_name ASC;
-- name: GetDeliveryZoneByArea :one
SELECT *
FROM delivery_zones
WHERE $1 = ANY(area_names)
  AND is_active = TRUE
LIMIT 1;
-- name: UpdateDeliveryZone :one
UPDATE delivery_zones
SET
    zone_name = $2,
    area_names = $3,
    delivery_fee = $4,
    min_delivery_time = $5,
    max_delivery_time = $6,
    is_active = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: ToggleDeliveryZoneStatus :one
UPDATE delivery_zones
SET
    is_active = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: DeleteDeliveryZone :exec
DELETE FROM delivery_zones
WHERE id = $1;
