-- name: CreateEvent :one
INSERT INTO events (
    event_type,
    payload
) VALUES ($1, $2)
RETURNING *;

-- name: ListUndeliveredEvents :many
SELECT *
FROM events
WHERE delivered = FALSE
ORDER BY created_at ASC;

-- name: MarkEventDelivered :one
UPDATE events
SET delivered = TRUE,
    delivered_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;
