-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    event_id,
    title,
    message,
    type,
    priority,
    metadata
) VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING *;
-- name: GetNotificationByID :one
SELECT *
FROM notifications
WHERE id = $1;
-- name: ListNotificationsByUser :many
SELECT *
FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC;
-- name: MarkNotificationRead :one
UPDATE notifications
SET is_read = TRUE, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1;
