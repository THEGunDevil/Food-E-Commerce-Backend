-- name: ListActiveWebhooksByEventType :many
SELECT *
FROM webhooks
WHERE event_type = $1 AND is_active = TRUE;
